package storage

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
)

type SQLiteStorage struct {
	sync.RWMutex
	directory string
	db        *sqlx.DB
	events    chan *Event
}

func CreateSQLiteStorage(directory string) *SQLiteStorage {
	return &SQLiteStorage{
		directory: strings.TrimRight(directory, "/"),
		events:    make(chan *Event, maxEventsQueueSize),
	}
}

func (s *SQLiteStorage) withTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), maxQueryTimeout)
}

func (s *SQLiteStorage) GetEvents() chan *Event {
	return s.events
}

func (s *SQLiteStorage) createMessagesTable() error {
	ctx, cancel := s.withTimeoutContext()
	defer cancel()

	stmt := `
    CREATE TABLE IF NOT EXISTS messages(
        ID STRING,
        "From" JSON,
        "To" JSON,
        CreatedAt TIMESTAMP,
        Content JSON,
        MIME JSON,
        Raw JSON,
        Unread BOOLEAN
    );`
	_, err := s.db.ExecContext(ctx, stmt)

	return err
}

func (s *SQLiteStorage) Init() error {
	s.Lock()
	defer s.Unlock()
	var err error

	s.db, err = sqlx.Connect("sqlite3", s.directory+"/mailcage.sqlite")
	if err != nil {
		return err
	}

	return s.createMessagesTable()
}

func (s *SQLiteStorage) Shutdown() error {
	return s.db.Close()
}

func (s *SQLiteStorage) Store(message *Message) (string, error) {
	s.Lock()
	defer s.Unlock()

	ctx, cancel := s.withTimeoutContext()
	defer cancel()

	stmt := `
    INSERT INTO messages (ID, "From", "To", CreatedAt, Content, MIME, Raw, Unread)
    VALUES (:ID, :From, :To, :CreatedAt, :Content, :MIME, :Raw, 1)`

	_, err := s.db.NamedExecContext(ctx, stmt, message)
	if err == nil {
		s.events <- addStoredMessageEvent(message)
	}

	return message.ID, err
}

func (s *SQLiteStorage) Update(message *Message) error {
	s.Lock()
	defer s.Unlock()

	ctx, cancel := s.withTimeoutContext()
	defer cancel()

	stmt := `
    UPDATE messages
    SET "From" = :From, "To" = :To, CreatedAt = :CreatedAt, Content = :Content, MIME = :MIME, Raw = :Raw, Unread = :Unread
    WHERE ID = :ID`

	_, err := s.db.NamedExecContext(ctx, stmt, message)
	if err == nil {
		s.events <- addStoredMessageEvent(message)
	}

	return err
}

func (s *SQLiteStorage) Get(start int, limit int) ([]*Message, error) {
	ctx, cancel := s.withTimeoutContext()
	defer cancel()

	stmt := `select * from messages order by CreatedAt desc limit $1 offset $2`
	messages := make([]*Message, 0)
	err := s.db.SelectContext(ctx, &messages, stmt, limit, start)

	return messages, err
}

func (s *SQLiteStorage) GetOne(id string) (*Message, error) {
	ctx, cancel := s.withTimeoutContext()
	defer cancel()

	var message Message
	err := s.db.GetContext(ctx, &message, `select * from messages where id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, ErrMessageNotFound
	}

	return &message, err
}

func (s *SQLiteStorage) DeleteAll() error {
	s.Lock()
	defer s.Unlock()

	ctx, cancel := s.withTimeoutContext()
	defer cancel()

	_, err := s.db.ExecContext(ctx, `delete from messages`)
	if err == nil {
		s.events <- addDeletedMessagesEvent()
	}

	return err
}

func (s *SQLiteStorage) DeleteOne(id string) error {
	s.Lock()
	defer s.Unlock()

	ctx, cancel := s.withTimeoutContext()
	defer cancel()

	_, err := s.db.ExecContext(ctx, `delete from messages where ID = $1`, id)
	if err == sql.ErrNoRows {
		return ErrMessageNotFound
	}

	if err == nil {
		s.events <- addDeletedMessageEvent(id)
	}

	return err
}

func (s *SQLiteStorage) Count() (int, error) {
	ctx, cancel := s.withTimeoutContext()
	defer cancel()

	var cnt int
	err := s.db.GetContext(ctx, &cnt, `select count(*) from messages`)

	return cnt, err
}
