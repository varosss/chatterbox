package httphandler

import (
	"chatterbox/chat/internal/application/usecase"
	"chatterbox/chat/internal/domain/valueobject"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	CreateChatUC *usecase.CreateChatUseCase
}

func NewChatHandler(
	createChatUC *usecase.CreateChatUseCase,
) *ChatHandler {
	return &ChatHandler{
		CreateChatUC: createChatUC,
	}
}

func (h *ChatHandler) Create(c *gin.Context) {
	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	participantsIDs := make([]valueobject.UserID, len(req.Participants))
	for i, userID := range req.Participants {
		parsedUserID, err := valueobject.ParseUserID(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		participantsIDs[i] = parsedUserID
	}

	res, err := h.CreateChatUC.Execute(c.Request.Context(), usecase.CreateChatCommand{
		ParticipantIDs: participantsIDs,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateChatResponse{ID: res.ChatID.String()})
}
