package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReasonRepo struct {
	db *sqlx.DB
}

func NewReasonRepo(db *sqlx.DB) *ReasonRepo {
	return &ReasonRepo{db: db}
}

func (r *ReasonRepo) Create(ctx context.Context, reason models.ReasonDTO) (id uuid.UUID, err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, operation_id, date, value) VALUES ($1, $2, $3, $4)", ReasonsTable)
	id = uuid.New()

	_, err = r.db.Exec(query, reason.OperationId, reason.Date, reason.Value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, nil
}

func (r *ReasonRepo) Delete(ctx context.Context, reason models.ReasonDTO) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", ReasonsTable)

	if _, err := r.db.Exec(query, reason.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
