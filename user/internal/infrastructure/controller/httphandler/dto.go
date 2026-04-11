package httphandler

type RegisterRequest struct {
	Email       string `json:"email" example:"john_doe@email.com"`
	Username    string `json:"username" example:"john_doe123"`
	DisplayName string `json:"display_name" example:"John Doe"`
	Password    string `json:"password" example:"s$*tv7bv1)"`
}

type RegisterResponse struct {
	ID string `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"john_doe@email.com"`
	Password string `json:"password" example:"s$*tv7bv1)"`
}

type UserResponseData struct {
	ID          string `json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Email       string `json:"email" example:"john_doe@email.com"`
	Username    string `json:"username" example:"john_doe123"`
	DisplayName string `json:"display_name" example:"John Doe"`
	Status      int    `json:"status" example:"active"`
}

type ListUsersRequest struct {
	UserIDs []string `form:"user_ids"`
	Limit   int      `form:"limit" example:"500"`
	Offset  int      `form:"offset" example:"500"`
}

type ListUsersResponse struct {
	Users []UserResponseData `json:"users"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}
