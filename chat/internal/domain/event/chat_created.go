package event

import (
	"chatterbox/chat/internal/domain/valueobject"
)

type ChatCreated struct {
	BaseEvent
	ChatID         valueobject.ChatID
	ParticipantIDs []valueobject.UserID
}

func NewChatCreatedEvent(
	chatID valueobject.ChatID,
	participantIDs []valueobject.UserID,
) ChatCreated {
	return ChatCreated{
		BaseEvent:      NewBaseEvent(),
		ChatID:         chatID,
		ParticipantIDs: participantIDs,
	}
}

func (e ChatCreated) ParticipantIDsAsUUIDs() []string {
	uuids := make([]string, len(e.ParticipantIDs))
	for i, id := range e.ParticipantIDs {
		uuids[i] = id.String()
	}
	return uuids
}

func (e ChatCreated) Name() string {
	return "chat.created"
}
