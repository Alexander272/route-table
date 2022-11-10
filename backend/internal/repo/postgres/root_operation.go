package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RootOperationRepo struct {
	db *sqlx.DB
}

func NewRootOperationRepo(db *sqlx.DB) *RootOperationRepo {
	return &RootOperationRepo{db: db}
}

func (r *RootOperationRepo) Get(ctx context.Context) (operations []models.RootOperation, err error) {
	query := fmt.Sprintf("SELECT id, title, gasket, step_number FROM %s", RootOperationTable)

	if err := r.db.Select(&operations, query); err != nil {
		return operations, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return operations, nil
}

func (r *RootOperationRepo) Create(ctx context.Context, operation models.RootOperationDTO) (id uuid.UUID, err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, title, gasket, step_number) VALUES ($1, $2, $3, $4)", RootOperationTable)

	id = uuid.New()

	_, err = r.db.Exec(query, id, operation.Title, operation.Gasket, operation.StepNumber)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, nil
}

func (r *RootOperationRepo) Update(ctx context.Context, operation models.RootOperationDTO) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1, gasket=$2, step_number=$3 WHERE id=$4", RootOperationTable)

	_, err := r.db.Exec(query, operation.Title, operation.Gasket, operation.StepNumber, operation.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *RootOperationRepo) Delete(ctx context.Context, operation models.RootOperationDTO) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", RootOperationTable)

	if _, err := r.db.Exec(query, operation.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
