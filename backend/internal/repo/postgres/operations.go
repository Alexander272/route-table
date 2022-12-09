package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type OperationRepo struct {
	db *sqlx.DB
}

func NewOperationRepo(db *sqlx.DB) *OperationRepo {
	return &OperationRepo{db: db}
}

func (r *OperationRepo) Get(ctx context.Context, positionId uuid.UUID, enabled pq.StringArray) (operations []models.Operation, err error) {
	query := fmt.Sprintf(`SELECT %s.id, title, done, remainder, step_number, is_finish FROM %s INNER JOIN %s ON operation_id=%s.id 
		WHERE position_id=$1 AND array[%s.id] <@ $2 ORDER BY step_number`,
		OperationsTable, OperationsTable, RootOperationTable, RootOperationTable, RootOperationTable)

	if err := r.db.Select(&operations, query, positionId, enabled); err != nil {
		return operations, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return operations, nil
}

func (r *OperationRepo) GetAll(ctx context.Context, positionId uuid.UUID) (operations []models.Operation, err error) {
	query := fmt.Sprintf(`SELECT %s.id, title, done, remainder, step_number FROM %s INNER JOIN %s ON operation_id=%s.id 
		WHERE position_id=$1  ORDER BY step_number`, OperationsTable, OperationsTable, RootOperationTable, RootOperationTable)

	if err := r.db.Select(&operations, query, positionId); err != nil {
		return operations, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return operations, nil
}

func (r *OperationRepo) GetWithReasons(ctx context.Context, positionId uuid.UUID) (operations []models.OperationWithReason, err error) {
	query := fmt.Sprintf(`SELECT %s.id, title, done, remainder, step_number, is_finish, %s.id as reason_id, value, %s.date
  		FROM %s INNER JOIN %s ON %s.operation_id=%s.id LEFT JOIN %s ON %s.id=%s.operation_id 
		WHERE position_id=$1 ORDER BY step_number, %s.date`, OperationsTable, ReasonsTable, ReasonsTable,
		OperationsTable, RootOperationTable, OperationsTable, RootOperationTable,
		ReasonsTable, OperationsTable, ReasonsTable, ReasonsTable)

	if err := r.db.Select(&operations, query, positionId); err != nil {
		return operations, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return operations, nil
}

// Получение связанных операций
func (r *OperationRepo) GetConnected(ctx context.Context, positionId, operationId uuid.UUID) (operations []models.Operation, err error) {
	query := fmt.Sprintf(`SELECT id, remainder FROM %s WHERE position_id=$1 AND array[operation_id] <@ 
		(SELECT connected FROM %s INNER JOIN %s ON operation_id=%s.id WHERE %s.id=$2)`,
		OperationsTable, OperationsTable, RootOperationTable, RootOperationTable, OperationsTable)

	if err := r.db.Select(&operations, query, positionId, operationId); err != nil {
		return operations, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return operations, nil
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
		args = append(args, p.Id, p.OperationId, p.PositionId, p.Done, p.Remainder)
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

func (r *OperationRepo) DeleteFew(ctx context.Context, fewId []uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id IN (", OperationsTable)
	args := make([]interface{}, 0)
	values := make([]string, 0, len(fewId))

	for i, p := range fewId {
		values = append(values, fmt.Sprintf("$%d", i+1))
		args = append(args, p)
	}
	query += strings.Join(values, ", ")
	query += ")"

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
