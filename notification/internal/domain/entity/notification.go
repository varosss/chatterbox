package entity

import (
	"chatterbox/notification/internal/domain/valueobject"
	"time"
)

type Notification struct {
	ID          valueobject.NotificationID
	RecepientID valueobject.UserID
	Type        valueobject.NotificationType
	Payload     interface{}
}

type MessagePayload struct {
	ID         string
	ChatID     string
	SenderID   string
	Text       string
	OccurredAt time.Time
}
