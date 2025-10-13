package unipile

import (
	"context"
	"database/sql"
	"fmt"
)

type UnipileStore struct {
	db *sql.DB
}

func NewUnipileStore(db *sql.DB) *UnipileStore {
	return &UnipileStore{
		db: db,
	}
}

func (s UnipileStore) CreateUnipileUserFederal(ctx context.Context,
	createUnipileUserFederalParam CreateUnipileUserFederaParam,
) (UnipileUserFederalEntity, error) {
	queryBuilder, err :=
		s.db.Prepare("INSERT INTO unipile_user_federals(account_id, provider, user_id, status) VALUES($1, $2, $3, $4) RETURNING *;")
	if err != nil {
		return UnipileUserFederalEntity{}, fmt.Errorf("prepare statement unipile_user_federals: %w", err)
	}
	defer queryBuilder.Close()
	var resultUnipileUserFederal UnipileUserFedera
	err = queryBuilder.QueryRowContext(ctx,
		createUnipileUserFederalParam.AccountID,
		createUnipileUserFederalParam.Provider,
		createUnipileUserFederalParam.UserID,
		createUnipileUserFederalParam.Status,
	).Scan(
		&resultUnipileUserFederal.ID,
		&resultUnipileUserFederal.AccountID,
		&resultUnipileUserFederal.Provider,
		&resultUnipileUserFederal.UserID,
		&resultUnipileUserFederal.Status,
		&resultUnipileUserFederal.CreatedAt,
		&resultUnipileUserFederal.UpdatedAt,
	)
	if err != nil {
		return UnipileUserFederalEntity{}, fmt.Errorf("failed to insert users: %w", err)
	}
	return ConvertToUnipileUserFederalEntity(resultUnipileUserFederal), nil
}

// ConvertToUnipileUserFederalEntity 把 model 轉換成 entity
func ConvertToUnipileUserFederalEntity(unipileUserFedera UnipileUserFedera) UnipileUserFederalEntity {
	return UnipileUserFederalEntity(unipileUserFedera)
}
