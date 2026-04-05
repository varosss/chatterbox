package event

import "time"

type ChatCreated struct {
	ChatID       string    `json:"chat_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Participants []string  `json:"participants"`
	OccurredAt   time.Time `json:"occurred_at"`
}
