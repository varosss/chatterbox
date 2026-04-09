package httphandler

import (
	"chatterbox/user/internal/application/usecase"
	"chatterbox/user/internal/domain/valueobject"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	listUsersUC *usecase.ListUsersUseCase
}

func NewUserHandler(
	listUsersUC *usecase.ListUsersUseCase,
) *UserHandler {
	return &UserHandler{
		listUsersUC: listUsersUC,
	}
}

func (h *UserHandler) List(c *gin.Context) {
	var req ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userIDs := make([]valueobject.UserID, len(req.UserIDs))
	for i, rawUserID := range req.UserIDs {
		parsedUserID, err := valueobject.ParseUserID(rawUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		userIDs[i] = parsedUserID
	}

	res, err := h.listUsersUC.Execute(c.Request.Context(), usecase.ListUsersCommand{
		UserIDs: userIDs,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	usersRes := make([]UserResponseData, len(res.Users))
	for i, user := range res.Users {
		usersRes[i] = UserResponseData{
			ID:          user.ID,
			Email:       user.Email,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Status:      user.Status,
		}
	}

	c.JSON(http.StatusBadRequest, ListUsersResponse{Users: usersRes})
}
