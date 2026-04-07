package entity

import (
	"chatterbox/chat/internal/domain/event"
	"chatterbox/chat/internal/domain/valueobject"
	"errors"
)

type Chat struct {
	id             valueobject.ChatID
	participantIDs []valueobject.UserID

	events []event.Event
}

func NewChat(
	participantsIDs []valueobject.UserID,
) (*Chat, error) {
	if len(participantsIDs) == 0 {
		return nil, errors.New("cannot create chat without participants")
	}

	chat := &Chat{
		id:             valueobject.NewChatID(),
		participantIDs: participantsIDs,
	}

	chat.recordEvent(event.NewChatCreatedEvent(chat.ID(), chat.ParticipantIDs()))

	return chat, nil
}

func ChatFromPrimitives(
	id valueobject.ChatID,
	participantIDs []valueobject.UserID,
) *Chat {
	return &Chat{
		id:             id,
		participantIDs: participantIDs,
	}
}

func (c *Chat) recordEvent(event event.Event) {
	c.events = append(c.events, event)
}

func (c *Chat) PullEvents() []event.Event {
	defer func() {
		c.events = nil
	}()

	return c.events
}

func (c *Chat) ID() valueobject.ChatID {
	return c.id
}

func (c *Chat) ParticipantIDs() []valueobject.UserID {
	return c.participantIDs
}

func (c *Chat) ParticipantIDsAsStrings() []string {
	uuids := make([]string, len(c.participantIDs))
	for i, id := range c.participantIDs {
		uuids[i] = id.String()
	}
	return uuids
}
