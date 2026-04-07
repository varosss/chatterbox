package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"chatterbox/notification/internal/infrastructure/adapter/websocket/hub"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ⚠️ TODO: заменить на нормальную проверку origin
	},
}

type WSHandler struct {
	hub *hub.InMemoryHub
}

func NewWSHandler(h *hub.InMemoryHub) *WSHandler {
	return &WSHandler{hub: h}
}

func (h *WSHandler) Handle(c *gin.Context) {
	var req WSConnectRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	h.hub.Register(req.UserID, conn)
}
