package httphandler

type RegisterRequest struct {
	Email       string `json:"email" example:"john_doe@email.com"`
	Username    string `json:"username" example:"john_doe123"`
	DisplayName string `json:"display_name" example:"John Doe"`
	Password    string `json:"password" example:"s$*tv7bv1)"`
}

type RegisterResponse struct {
	UserID string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"john_doe@email.com"`
	Password string `json:"password" example:"s$*tv7bv1)"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
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
}

type ListUsersResponse struct {
	Users []UserResponseData `json:"users"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}
