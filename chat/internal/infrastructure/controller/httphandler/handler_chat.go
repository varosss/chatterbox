package httphandler

import (
	"chatterbox/chat/internal/application/usecase"
	"chatterbox/chat/internal/domain/valueobject"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	CreateChatUC *usecase.CreateChatUseCase
	ListChatsUC  *usecase.ListChatsUseCase
}

func NewChatHandler(
	createChatUC *usecase.CreateChatUseCase,
	listChatsUC *usecase.ListChatsUseCase,
) *ChatHandler {
	return &ChatHandler{
		CreateChatUC: createChatUC,
		ListChatsUC:  listChatsUC,
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

func (h *ChatHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	parsedUserID, err := valueobject.ParseUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "user_id is invalid"})
		return
	}

	res, err := h.ListChatsUC.Execute(c.Request.Context(), usecase.ListChatsCommand{
		UserID: parsedUserID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	chatsResp := make([]ChatResponseData, len(res.Chats))
	for i, chat := range res.Chats {
		chatsResp[i] = ChatResponseData{
			ID:             chat.ID,
			ParticipantIDs: chat.ParticipantIDs,
		}
	}

	c.JSON(http.StatusOK, ListChatsResponse{Chats: chatsResp})
}
