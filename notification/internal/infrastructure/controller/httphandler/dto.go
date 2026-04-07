package httphandler

type WSConnectRequest struct {
	UserID string `form:"user_id" binding:"required,uuid"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}
