package httphandler

type CreateChatRequest struct {
	Participants []string `json:"participants"`
}

type CreateChatResponse struct {
	ID string `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
}

type CreateMessageRequest struct {
	SenderID string `json:"sender_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ChatID   string `json:"chat_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Text     string `json:"text" example:"Hello, World!"`
}

type CreateMessageResponse struct {
	ID string `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}
