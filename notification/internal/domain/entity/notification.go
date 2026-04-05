package entity

import "time"

type Data interface{}

type Notification struct {
	ID     string
	UserID string
	Type   string
	Data   Data
}

type MessageData struct {
	ID         string
	ChatID     string
	SenderID   string
	Text       string
	OccurredAt time.Time
}
