package event

import (
	"chatterbox/chat/internal/domain/valueobject"
)

type MessageCreated struct {
	BaseEvent
	SenderID    valueobject.UserID
	ReceiverIDs []valueobject.UserID
	ChatID      valueobject.ChatID
	MessageID   valueobject.MessageID
	Text        string
}

func NewMessageCreatedEvent(
	senderID valueobject.UserID,
	receiverIDs []valueobject.UserID,
	chatID valueobject.ChatID,
	messageID valueobject.MessageID,
	text string,
) MessageCreated {
	return MessageCreated{
		BaseEvent:   NewBaseEvent(),
		SenderID:    senderID,
		ReceiverIDs: receiverIDs,
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        text,
	}
}

func (e MessageCreated) Name() string {
	return "message.created"
}
