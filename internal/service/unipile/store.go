package unipile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
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

func (s UnipileStore) ListUnipileUserFederalByUserID(ctx context.Context,
	listFederaParam ListFederaParam,
) ([]UnipileUserFederalEntity, error) {
	// original sql
	queryBuilder := sq.Select("id", "account_id", "provider", "user_id", "status", "created_at", "updated_at").
		From("unipile_user_federals").PlaceholderFormat(sq.Dollar)
	queryBuilder = queryBuilder.Where(sq.Eq{"user_id": listFederaParam.UserID})
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to use query builder: %w", err)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []UnipileUserFederalEntity{}, nil
		}
		return nil, fmt.Errorf("failed to executed %w", err)
	}
	result := make([]UnipileUserFederalEntity, 0, 100)
	for rows.Next() {
		var resultUnipileUserFederal UnipileUserFedera
		err = rows.Scan(&resultUnipileUserFederal.ID,
			&resultUnipileUserFederal.AccountID,
			&resultUnipileUserFederal.Provider,
			&resultUnipileUserFederal.UserID,
			&resultUnipileUserFederal.Status,
			&resultUnipileUserFederal.CreatedAt,
			&resultUnipileUserFederal.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, ConvertToUnipileUserFederalEntity(resultUnipileUserFederal))
	}
	return result, nil
}

// ConvertToUnipileUserFederalEntity 把 model 轉換成 entity
func ConvertToUnipileUserFederalEntity(unipileUserFedera UnipileUserFedera) UnipileUserFederalEntity {
	return UnipileUserFederalEntity(unipileUserFedera)
}
