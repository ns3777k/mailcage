package storage

import "time"

type MIMEBody struct {
    Parts []*Content
}

type Path struct {
    Relays  []string
    Mailbox string
    Domain  string
    Params  string
}

type Content struct {
    Headers map[string][]string
    Body    string
    Size    int
    MIME    *MIMEBody
}

type RawMessage struct {
    From string
    To   []string
    Data string
    Helo string
}

type Message struct {
    ID string
    From    *Path
    To      []*Path
    CreatedAt time.Time
    Content *Content
    MIME    *MIMEBody // FIXME refactor to use Content.MIME
    Raw     *RawMessage
}

type Storage interface {
    Get(start int, limit int) ([]*Message, error)
    GetOne(id string) (*Message, error)
    Store(message *Message) (string, error)
    DeleteAll() error
    DeleteOne(id string) error
    Count() int
}
