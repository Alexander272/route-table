package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CompliteOperationRepo struct {
	db *sqlx.DB
}

func NewComplitedOperationRepo(db *sqlx.DB) *CompliteOperationRepo {
	return &CompliteOperationRepo{db: db}
}

func (r *CompliteOperationRepo) Create(ctx context.Context, operation models.CompleteOperation) (id uuid.UUID, err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, operation_id, remainder, count) VALUES ($1, $2, $3, $4)", ComplitedOperTable)
	id = uuid.New()

	_, err = r.db.Exec(query, id, operation.Id, operation.Remainder, operation.Count)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, nil
}

func (r *CompliteOperationRepo) Delete(ctx context.Context, operation models.ComplitedOperation) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", ComplitedOperTable)

	if _, err := r.db.Exec(query, operation.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
