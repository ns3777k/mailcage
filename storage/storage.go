package storage

import (
    "database/sql/driver"
    "encoding/json"
    "github.com/pkg/errors"
    "time"
)

var (
    ErrMessageNotFound = errors.New("message not found")
)

type MIMEBody struct {
    Parts []*Content
}

// Value implements the driver.Valuer interface
func (m MIMEBody) Value() (driver.Value, error) {
    b, err := json.Marshal(m)
    if err != nil {
        return nil, err
    }

    return string(b[:]), err
}

func (m *MIMEBody) Scan(src interface{}) error {
    if src == nil {
        *m = MIMEBody{}
    } else {
        return json.Unmarshal([]byte(src.(string)), m)
    }

    return nil
}

type Path struct {
    Relays  []string
    Mailbox string
    Domain  string
    Params  string
}

// Value implements the driver.Valuer interface
func (p Path) Value() (driver.Value, error) {
    b, err := json.Marshal(p)
    if err != nil {
        return nil, err
    }

    return string(b[:]), err
}

// Scan implements the sql.Scanner interface
func (p *Path) Scan(src interface{}) error {
   if src == nil {
       *p = Path{}
   } else {
       return json.Unmarshal([]byte(src.(string)), p)
   }

   return nil
}

type Content struct {
    Headers map[string][]string
    Body    string
    Size    int
    MIME    *MIMEBody
}

// Value implements the driver.Valuer interface
func (c Content) Value() (driver.Value, error) {
    b, err := json.Marshal(c)
    if err != nil {
        return nil, err
    }

    return string(b[:]), err
}

func (c *Content) Scan(src interface{}) error {
    if src == nil {
        *c = Content{}
    } else {
        return json.Unmarshal([]byte(src.(string)), c)
    }

    return nil
}

type RawMessage struct {
    From string
    To   []string
    Data string
    Helo string
}

// Value implements the driver.Valuer interface
func (r RawMessage) Value() (driver.Value, error) {
    b, err := json.Marshal(r)
    if err != nil {
        return nil, err
    }

    return string(b[:]), err
}

func (r *RawMessage) Scan(src interface{}) error {
    if src == nil {
        *r = RawMessage{}
    } else {
        return json.Unmarshal([]byte(src.(string)), r)
    }

    return nil
}

// Needed to serialize []*Path
type Paths []*Path

func (p Paths) Value() (driver.Value, error) {
    b, err := json.Marshal(p)
    if err != nil {
        return nil, err
    }

    return string(b[:]), err
}

// Scan implements the sql.Scanner interface
func (p *Paths) Scan(src interface{}) error {
    if src == nil {
        *p = Paths{}
    } else {
        return json.Unmarshal([]byte(src.(string)), p)
    }

    return nil
}

type Message struct {
    ID string `db:"ID"`
    From    *Path `db:"From"`
    To      Paths `db:"To"`
    CreatedAt time.Time `db:"CreatedAt"`
    Content *Content `db:"Content"`
    MIME    *MIMEBody `db:"MIME"` // FIXME refactor to use Content.MIME
    Raw     *RawMessage `db:"Raw"`
}

type Storage interface {
    Get(start int, limit int) ([]*Message, error)
    GetOne(id string) (*Message, error)
    Store(message *Message) (string, error)
    DeleteAll() error
    DeleteOne(id string) error
    Count() (int, error)
    Shutdown() error
    Init() error
}
