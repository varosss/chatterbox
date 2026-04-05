package sender

import "time"

type Data interface{}

type Notification struct {
	ID     string `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	UserID string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Type   string `json:"type" example:"new_message"`
	Data   Data   `json:"data"`
}

type MessageData struct {
	ID         string    `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ChatID     string    `json:"chat_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	SenderID   string    `json:"sender_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Text       string    `json:"text" example:"Hello, world!"`
	OccurredAt time.Time `json:"occurred_at" example:"0001-01-01T00:00:00Z"`
}
