package event

import "time"

type MessageCreated struct {
	MessageID  string    `json:"message_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ChatID     string    `json:"chat_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	SenderID   string    `json:"sender_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Receivers  []string  `json:"receivers"`
	Text       string    `json:"text" example:"Hello, world!"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (e MessageCreated) Type() string {
	return "message.created"
}
