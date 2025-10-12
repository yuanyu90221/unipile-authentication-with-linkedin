package user

import "time"

type User struct {
	ID             int64     `json:"id" db:"id"`
	Account        string    `json:"account" db:"account"`
	HashedPassword string    `json:"hashed_password" db:"hashed_password"`
	RefreshToken   *string   `json:"refresh_token,omitempty" db:"refresh_token" `
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
