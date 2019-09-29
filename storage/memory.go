package storage

import (
    "github.com/pkg/errors"
    "sync"
)

var (
    ErrMessageNotFound = errors.New("message not found")
)

type MemoryStorage struct {
    sync.RWMutex
    messages []*Message
    idIndex map[string]int
}

func CreateMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        messages: make([]*Message, 0),
        idIndex: make(map[string]int),
    }
}

func (m *MemoryStorage) Store(message *Message) (string, error) {
    m.Lock()
    defer m.Unlock()

    m.messages = append(m.messages, message)
    m.idIndex[message.ID] = len(m.messages) - 1

    return message.ID, nil
}

func (m *MemoryStorage) Count() int {
    m.RLock()
    defer m.RUnlock()

    return len(m.messages)
}

func (m *MemoryStorage) Get(start int, limit int) ([]*Message, error) {
    messages := make([]*Message, 0)
    messagesAmount := m.Count()

    if messagesAmount == 0 || start > messagesAmount {
        return messages, nil
    }

    if start+limit > messagesAmount {
        limit = messagesAmount - start
    }

    start = messagesAmount - start - 1
    end := start - limit

    if start < 0 {
        start = 0
    }

    if end < -1 {
        end = -1
    }

    m.RLock()
    for i := start; i > end; i-- {
        messages = append(messages, m.messages[i])
    }
    m.RUnlock()

    return messages, nil
}

func (m *MemoryStorage) DeleteAll() error {
    m.Lock()
    defer m.Unlock()

    m.messages = make([]*Message, 0)
    m.idIndex = make(map[string]int)

    return nil
}

func (m *MemoryStorage) DeleteOne(id string) error {
    m.Lock()
    defer m.Unlock()

    index, ok := m.idIndex[id]
    if !ok {
        return ErrMessageNotFound
    }

    delete(m.idIndex, id)

    for k, v := range m.idIndex {
        if v > index {
            m.idIndex[k] = v - 1
        }
    }

    m.messages = append(m.messages[:index], m.messages[index+1:]...)

    return nil
}

func (m *MemoryStorage) GetOne(id string) (*Message, error) {
    m.RLock()
    defer m.RUnlock()

    index, ok := m.idIndex[id]
    if !ok {
        return nil, ErrMessageNotFound
    }

    return m.messages[index], nil
}
