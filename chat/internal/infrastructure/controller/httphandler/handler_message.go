package httphandler

import (
	"chatterbox/chat/internal/application/usecase"
	"chatterbox/chat/internal/domain/valueobject"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	CreateMessageUC *usecase.CreateMessageUseCase
}

func NewMessageHandler(
	createMessageUC *usecase.CreateMessageUseCase,
) *MessageHandler {
	return &MessageHandler{
		CreateMessageUC: createMessageUC,
	}
}

func (h *MessageHandler) Get() {
	// TODO
}

func (h *MessageHandler) Create(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	senderID, err := valueobject.ParseUserID(req.SenderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	chatID, err := valueobject.ParseChatID(req.ChatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	res, err := h.CreateMessageUC.Execute(c.Request.Context(), usecase.CreateMessageCommand{
		SenderID: senderID,
		ChatID:   chatID,
		Text:     req.Text,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateMessageResponse{ID: res.MessageID.String()})
}
