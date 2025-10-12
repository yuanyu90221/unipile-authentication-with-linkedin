package user

import "time"

type UserEntity struct {
	ID             int64     `json:"id"`
	Account        string    `json:"account"`
	HashedPassword string    `json:"hashed_password,omitempty"`
	RefreshToken   *string   `json:"refresh_token,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateUserParam struct {
	Account        string `json:"account"`
	HashedPassword string `json:"hashed_password,omitempty"`
}
