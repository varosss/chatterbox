package httphandler

import (
	"chatterbox/user/internal/application/usecase"
	"chatterbox/user/internal/domain/valueobject"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	getUserUC   *usecase.GetUserUseCase
	listUsersUC *usecase.ListUsersUseCase
}

func NewUserHandler(
	getUserUC *usecase.GetUserUseCase,
	listUsersUC *usecase.ListUsersUseCase,
) *UserHandler {
	return &UserHandler{
		getUserUC:   getUserUC,
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
		Limit:   req.Limit,
		Offset:  req.Offset,
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

	c.JSON(http.StatusOK, ListUsersResponse{Users: usersRes})
}

func (h *UserHandler) Me(c *gin.Context) {
	parsedUserID, err := valueobject.ParseUserID(c.MustGet("user_id").(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := h.getUserUC.Execute(c.Request.Context(), usecase.GetUserCommand{
		UserID: parsedUserID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponseData{
		ID:          res.User.ID,
		Email:       res.User.Email,
		Username:    res.User.Username,
		DisplayName: res.User.DisplayName,
		Status:      res.User.Status,
	})
}
