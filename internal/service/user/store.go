package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
)

var (
	ErrorForDuplicateKey = errors.New("error for duplicate key")
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s UserStore) CreateUser(ctx context.Context, createUserParam CreateUserParam) (UserEntity, error) {
	queryBuilder, err :=
		s.db.Prepare("INSERT INTO users(account, hashed_password) VALUES($1,$2) RETURNING *;")
	if err != nil {
		return UserEntity{}, fmt.Errorf("prepare statement users: %w", err)
	}
	defer queryBuilder.Close()
	var resultUser User
	err = queryBuilder.QueryRowContext(ctx, createUserParam.Account, createUserParam.HashedPassword).Scan(
		&resultUser.ID,
		&resultUser.Account,
		&resultUser.HashedPassword,
		&resultUser.RefreshToken,
		&resultUser.CreatedAt,
		&resultUser.UpdatedAt,
	)
	log := logger.FromContext(ctx)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				log.Error(fmt.Sprintf("Duplicate key error %s", pgErr.Detail))
				return UserEntity{}, ErrorForDuplicateKey
			} else {
				log.Error(fmt.Sprintf("PostgreSQL error: %s", err))
			}
		}
		return UserEntity{}, fmt.Errorf("failed to insert users: %w", err)
	}
	return ConvertToUserEntity(resultUser), nil
}

// ConvertToUserEntity - 轉換 User 為 Entity
func ConvertToUserEntity(user User) UserEntity {
	return UserEntity(user)
}
