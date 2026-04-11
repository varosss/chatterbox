package httphandler

import "time"

type CreateChatRequest struct {
	DisplayName  string   `json:"display_name" example:"Group Chat"`
	Participants []string `json:"participants"`
}

type CreateChatResponse struct {
	ID string `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
}

type CreateMessageRequest struct {
	ChatID string `json:"chat_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Text   string `json:"text" example:"Hello, World!"`
}

type CreateMessageResponse struct {
	ID string `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
}

type ChatResponseData struct {
	ID             string   `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ParticipantIDs []string `json:"participant_ids"`
	DisplayName    string   `json:"display_name" example:"Group Chat 1"`
}

type MessageResponseData struct {
	ID        string    `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	SenderID  string    `json:"sender_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ChatID    string    `json:"chat_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Text      string    `json:"text" example:"Hello, world!"`
	CreatedAt time.Time `json:"created_at"`
}

// type ListChatsRequest struct {
// 	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
// 	Offset int `form:"offset" binding:"omitempty,min=0"`
// }

type ListChatsResponse struct {
	Chats []ChatResponseData `json:"chats"`
}

type ListMessagesRequest struct {
	ChatID string `form:"chat_id" binding:"required,uuid"`
}

type ListMessagesResponse struct {
	Messages []MessageResponseData `json:"messages"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}
