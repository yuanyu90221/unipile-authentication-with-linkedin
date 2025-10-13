package unipile

import "time"

type UnipileUserFedera struct {
	ID        int64     `json:"id" db:"id"`
	AccountID string    `json:"account_id" db:"account_id"`
	Provider  string    `json:"provider" db:"provider"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
