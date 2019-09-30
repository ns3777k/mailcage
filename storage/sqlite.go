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
}

func CreateSQLiteStorage(directory string) *SQLiteStorage {
	return &SQLiteStorage{
		directory: strings.TrimRight(directory, "/"),
	}
}
func (m *SQLiteStorage) createMessagesTable() error {
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
	_, err := m.db.Exec(stmt)

	return err
}

func (m *SQLiteStorage) Init() error {
	m.Lock()
	defer m.Unlock()
	var err error

	m.db, err = sqlx.Connect("sqlite3", m.directory+"/mailcage.sqlite")
	if err != nil {
		return err
	}

	return m.createMessagesTable()
}

func (m *SQLiteStorage) Shutdown() error {
	return m.db.Close()
}

func (m *SQLiteStorage) Store(message *Message) (string, error) {
	m.Lock()
	defer m.Unlock()

	stmt := `
    insert into messages (
        ID, "From", "To", CreatedAt, Content, MIME, Raw
    ) VALUES (
        :ID, :From, :To, :CreatedAt, :Content, :MIME, :Raw
    )`

	_, err := m.db.NamedExec(stmt, message)
	return message.ID, err
}

func (m *SQLiteStorage) Get(start int, limit int) ([]*Message, error) {
	stmt := `select * from messages order by CreatedAt asc limit $1 offset $2`
	messages := make([]*Message, 0)
	err := m.db.Select(&messages, stmt, limit, start)

	return messages, err
}

func (m *SQLiteStorage) GetOne(id string) (*Message, error) {
	var message Message
	err := m.db.Get(&message, `select * from messages where id = $1`, id)

	return &message, err
}

func (m *SQLiteStorage) DeleteAll() error {
	m.Lock()
	defer m.Unlock()

	_, err := m.db.Exec(`delete from messages`)
	return err
}

func (m *SQLiteStorage) DeleteOne(id string) error {
	m.Lock()
	defer m.Unlock()

	_, err := m.db.Exec(`delete from messages where ID = $1`, id)
	return err
}

func (m *SQLiteStorage) Count() (int, error) {
	var cnt int
	err := m.db.Get(&cnt, `select count(*) from messages`)

	return cnt, err
}
