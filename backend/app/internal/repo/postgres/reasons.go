package postgres

import (
	"context"
	"fmt"
	"strings"

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

func (r *ReasonRepo) Get(ctx context.Context) (reasons []models.PosWithReason, err error) {
	query := fmt.Sprintf(`SELECT number, %s.title AS pos_title, %s.title AS op_title, %s.date, value FROM %s
		INNER JOIN %s ON %s.operation_id = %s.id
		INNER JOIN %s ON %s.operation_id = %s.id
		INNER JOIN %s ON position_id = %s.id
		INNER JOIN %s ON order_id = %s.id`,
		PositionsTable, RootOperationTable, ReasonsTable, ReasonsTable,
		OperationsTable, ReasonsTable, OperationsTable,
		RootOperationTable, OperationsTable, RootOperationTable,
		PositionsTable, PositionsTable,
		OrdersTable, OrdersTable,
	)

	if err := r.db.Select(&reasons, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return reasons, nil
}

func (r *ReasonRepo) Create(ctx context.Context, reason models.ReasonDTO) (id uuid.UUID, err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, operation_id, date, value) VALUES ($1, $2, $3, $4)", ReasonsTable)
	id = uuid.New()

	_, err = r.db.Exec(query, id, reason.OperationId, reason.Date, reason.Value)
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

func (r *ReasonRepo) DeleteFew(ctx context.Context, reasons []string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id IN (", ReasonsTable)
	args := make([]interface{}, 0)
	values := make([]string, 0, len(reasons))

	for i, p := range reasons {
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
