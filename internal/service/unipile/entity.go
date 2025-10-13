package unipile

import "time"

type UnipileUserFederalEntity struct {
	ID        int64     `json:"id"`
	AccountID string    `json:"account_id"`
	Provider  string    `json:"provider"`
	UserID    int64     `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUnipileUserFederaParam struct {
	AccountID string `json:"account_id"`
	Provider  string `json:"provider"`
	UserID    int64  `json:"user_id"`
	Status    string `json:"status"`
}
