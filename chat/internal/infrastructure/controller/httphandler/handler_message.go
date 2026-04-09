package httphandler

import (
	"chatterbox/chat/internal/application/usecase"
	"chatterbox/chat/internal/domain/valueobject"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	CreateMessageUC *usecase.CreateMessageUseCase
	ListMessagesUC  *usecase.ListMessagesUseCase
}

func NewMessageHandler(
	createMessageUC *usecase.CreateMessageUseCase,
	listMessagesUC *usecase.ListMessagesUseCase,
) *MessageHandler {
	return &MessageHandler{
		CreateMessageUC: createMessageUC,
		ListMessagesUC:  listMessagesUC,
	}
}

func (h *MessageHandler) Create(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	senderID := c.MustGet("user_id").(string)

	parsedSenderID, err := valueobject.ParseUserID(senderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	parsedChatID, err := valueobject.ParseChatID(req.ChatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := h.CreateMessageUC.Execute(c.Request.Context(), usecase.CreateMessageCommand{
		SenderID: parsedSenderID,
		ChatID:   parsedChatID,
		Text:     req.Text,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateMessageResponse{ID: res.MessageID.String()})
}

func (h *MessageHandler) List(c *gin.Context) {
	var req ListMessagesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	parsedChatID, err := valueobject.ParseChatID(req.ChatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := h.ListMessagesUC.Execute(c.Request.Context(), usecase.ListMessagesCommand{
		ChatID: parsedChatID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	messagesResp := make([]MessageResponseData, len(res.Messages))
	for i, message := range res.Messages {
		messagesResp[i] = MessageResponseData{
			ID:        message.ID,
			SenderID:  message.SenderID,
			ChatID:    message.ChatID,
			Text:      message.Text,
			CreatedAt: message.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, ListMessagesResponse{Messages: messagesResp})
}
