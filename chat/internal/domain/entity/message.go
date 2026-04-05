package entity

import (
	"chatterbox/chat/internal/domain/event"
	"chatterbox/chat/internal/domain/valueobject"
	"errors"
	"time"
)

type Message struct {
	id        valueobject.MessageID
	senderID  valueobject.UserID
	chatID    valueobject.ChatID
	text      string
	createdAt time.Time

	events []event.Event
}

func NewMessage(
	senderID valueobject.UserID,
	chatID valueobject.ChatID,
	receiverIDs []valueobject.UserID,
	text string,
) (*Message, error) {
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}

	m := &Message{
		valueobject.NewMessageID(),
		senderID,
		chatID,
		text,
		time.Now(),
		[]event.Event{},
	}

	m.recordEvent(event.NewMessageCreatedEvent(
		senderID,
		receiverIDs,
		chatID,
		m.ID(),
		text,
	))

	return m, nil
}

func (m *Message) recordEvent(event event.Event) {
	m.events = append(m.events, event)
}

func (m *Message) PullEvents() []event.Event {
	defer func() {
		m.events = nil
	}()

	return m.events
}

func (m *Message) ID() valueobject.MessageID {
	return m.id
}

func (m *Message) SenderID() valueobject.UserID {
	return m.senderID
}

func (m *Message) ChatID() valueobject.ChatID {
	return m.chatID
}

func (m *Message) Text() string {
	return m.text
}

func (m *Message) CreatedAt() time.Time {
	return m.createdAt
}
