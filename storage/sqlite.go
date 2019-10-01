package storage

import (
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

func (s *SQLiteStorage) GetEvents() chan *Event {
	return s.events
}

func (s *SQLiteStorage) createMessagesTable() error {
	stmt := `
    create table if not exists messages(
        ID string,
        "From" json,
        "To" json,
        CreatedAt timestamp,
        Content json,
        MIME json,
        Raw json
    );`
	_, err := s.db.Exec(stmt)

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

	stmt := `
    insert into messages (
        ID, "From", "To", CreatedAt, Content, MIME, Raw
    ) VALUES (
        :ID, :From, :To, :CreatedAt, :Content, :MIME, :Raw
    )`

	_, err := s.db.NamedExec(stmt, message)
	if err == nil {
		s.events <- addStoredMessageEvent(message)
	}

	return message.ID, err
}

func (s *SQLiteStorage) Get(start int, limit int) ([]*Message, error) {
	stmt := `select * from messages order by CreatedAt asc limit $1 offset $2`
	messages := make([]*Message, 0)
	err := s.db.Select(&messages, stmt, limit, start)

	return messages, err
}

func (s *SQLiteStorage) GetOne(id string) (*Message, error) {
	var message Message
	err := s.db.Get(&message, `select * from messages where id = $1`, id)

	return &message, err
}

func (s *SQLiteStorage) DeleteAll() error {
	s.Lock()
	defer s.Unlock()

	_, err := s.db.Exec(`delete from messages`)
	if err == nil {
		s.events <- addDeletedMessagesEvent()
	}

	return err
}

func (s *SQLiteStorage) DeleteOne(id string) error {
	s.Lock()
	defer s.Unlock()

	_, err := s.db.Exec(`delete from messages where ID = $1`, id)
	if err == nil {
		s.events <- addDeletedMessageEvent(id)
	}

	return err
}

func (s *SQLiteStorage) Count() (int, error) {
	var cnt int
	err := s.db.Get(&cnt, `select count(*) from messages`)

	return cnt, err
}
