package storage

import (
    _ "github.com/mattn/go-sqlite3"
    "github.com/jmoiron/sqlx"
    "strings"
    "sync"
)

type FsStorage struct {
    sync.RWMutex
    directory string
    db *sqlx.DB
}

func CreateFsStorage(directory string) *FsStorage {
    return &FsStorage{
        directory: strings.TrimRight(directory, "/"),
    }
}
func (m *FsStorage) createMessagesTable() error {
    stmt := `create table if not exists messages (ID string, "From" json, "To" json, CreatedAt timestamp, Content json, MIME json, Raw json);`
    _, err := m.db.Exec(stmt)

    return err
}

func (m *FsStorage) Init() error {
    m.Lock()
    defer m.Unlock()
    var err error

    m.db, err = sqlx.Connect("sqlite3", m.directory+"/mailcage.sqlite")
    if err != nil {
        return err
    }

    return m.createMessagesTable()
}

func (m *FsStorage) Shutdown() error {
    return m.db.Close()
}

func (m *FsStorage) Store(message *Message) (string, error) {
    m.Lock()
    defer m.Unlock()

    _, err := m.db.NamedExec(`insert into messages (ID, "From", "To", CreatedAt, Content, MIME, Raw) VALUES (:ID, :From, :To, :CreatedAt, :Content, :MIME, :Raw)`, message)
    return message.ID, err
}

func (m *FsStorage) Get(start int, limit int) ([]*Message, error) {
    messages := make([]*Message, 0)
    err := m.db.Select(&messages, `select * from messages order by CreatedAt asc limit $1 offset $2`, limit, start)

    return messages, err
}

func (m *FsStorage) GetOne(id string) (*Message, error) {
    var message Message
    err := m.db.Get(&message, `select * from messages where id = $1`, id)

    return &message, err
}

func (m *FsStorage) DeleteAll() error {
    m.Lock()
    defer m.Unlock()

    _, err := m.db.Exec(`delete from messages`)
    return err
}

func (m *FsStorage) DeleteOne(id string) error {
    m.Lock()
    defer m.Unlock()

    _, err := m.db.Exec(`delete from messages where ID = $1`, id)
    return err
}

func (m *FsStorage) Count() (int, error) {
    var cnt int
    err := m.db.Get(&cnt, `select count(*) from messages`)

    return cnt, err
}
