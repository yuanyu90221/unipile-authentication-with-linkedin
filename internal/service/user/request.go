package user

// CreateUserRequest - 建立 user 的 request
type CreateUserRequest struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}
