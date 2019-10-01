package storage

type MessageEvent string

const (
	StoredMessageEvent   MessageEvent = "stored"
	DeletedMessageEvent  MessageEvent = "deleted_message"
	DeletedMessagesEvent MessageEvent = "deleted_messages"
)

type Event struct {
	Type    MessageEvent
	Payload interface{}
}

func addStoredMessageEvent(message *Message) *Event {
	return &Event{Type: StoredMessageEvent, Payload: message}
}

func addDeletedMessageEvent(id string) *Event {
	return &Event{Type: DeletedMessageEvent, Payload: id}
}

func addDeletedMessagesEvent() *Event {
	return &Event{Type: DeletedMessagesEvent}
}
