package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OperationRepo struct {
	db *sqlx.DB
}

func NewOperationRepo(db *sqlx.DB) *OperationRepo {
	return &OperationRepo{db: db}
}

func (r *OperationRepo) Create(ctx context.Context, operation models.OperationDTO) (id uuid.UUID, err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, operation_id, position_id, done, remainder) VALUES ($1, $2, $3, $4, $5)", OperationsTable)
	id = uuid.New()

	_, err = r.db.Exec(query, id, operation.OperationId, operation.PositionId, operation.Done, operation.Remainder)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, nil
}

func (r *OperationRepo) CreateFew(ctx context.Context, operations []models.OperationDTO) error {
	query := fmt.Sprintf("INSERT INTO %s (id, operation_id, position_id, done, remainder) VALUES ", OperationsTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(operations))

	c := 5
	for i, p := range operations {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5))
		id := uuid.New()
		args = append(args, id, p.OperationId, p.PositionId, p.Done, p.Remainder)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OperationRepo) Update(ctx context.Context, operation models.OperationDTO) error {
	query := fmt.Sprintf("UPDATE %s SET done=$1, remainder=$2, date=$3 WHERE id=$4", OperationsTable)

	_, err := r.db.Exec(query, operation.Done, operation.Remainder, operation.Date, operation.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *OperationRepo) Delete(ctx context.Context, operation models.OperationDTO) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", OperationsTable)

	if _, err := r.db.Exec(query, operation.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
