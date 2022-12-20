package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CompleteOperationRepo struct {
	db *sqlx.DB
}

func NewCompletedOperationRepo(db *sqlx.DB) *CompleteOperationRepo {
	return &CompleteOperationRepo{db: db}
}

func (r *CompleteOperationRepo) Get(ctx context.Context, operationId uuid.UUID) (operations []models.CompletedOperation, err error) {
	query := fmt.Sprintf(`
		SELECT id, operation_id, group_id, remainder, count FROM %s
		WHERE group_id=(SELECT group_id FROM %s WHERE id=$1);
	`, CompletedOperTable, CompletedOperTable)

	if err := r.db.Select(&operations, query, operationId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return operations, nil
}

func (r *CompleteOperationRepo) Create(ctx context.Context, operation models.CompletedOperation) (id uuid.UUID, err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, operation_id, group_id, remainder, count) VALUES ($1, $2, $3, $4, $5)", CompletedOperTable)

	_, err = r.db.Exec(query, operation.Id, operation.Id, operation.GroupId, operation.Remainder, operation.Count)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, nil
}

func (r *CompleteOperationRepo) CreateFew(ctx context.Context, operations []models.CompletedOperation) error {
	query := fmt.Sprintf("INSERT INTO %s (id, operation_id, group_id, remainder, count) VALUES ", CompletedOperTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(operations))

	c := 5
	for i, p := range operations {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5))
		args = append(args, p.Id, p.OperationId, p.GroupId, p.Remainder, p.Count)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *CompleteOperationRepo) Delete(ctx context.Context, operation models.CompletedOperation) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE group_id=$1", CompletedOperTable)

	if _, err := r.db.Exec(query, operation.GroupId); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
