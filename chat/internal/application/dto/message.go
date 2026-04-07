package dto

import (
	"time"
)

type Message struct {
	ID        string
	SenderID  string
	ChatID    string
	Text      string
	CreatedAt time.Time
}
