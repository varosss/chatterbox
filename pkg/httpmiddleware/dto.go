package httpmiddleware

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}
