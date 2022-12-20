package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Get(ctx context.Context, id uuid.UUID) (order models.OrderWithPositions, err error) {
	query := fmt.Sprintf("SELECT id, number, done, deadline FROM %s WHERE id=$1", OrdersTable)

	if err := r.db.Get(&order, query, id); err != nil {
		return models.OrderWithPositions{}, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return order, nil
}

func (r *OrderRepo) GetAll(ctx context.Context) (orders []models.Order, err error) {
	query := fmt.Sprintf(`SELECT id, number, done, date, deadline,
		(SELECT count(case when done = true then done end)/count(done)::real FROM %s WHERE order_id=%s.id) as progress
		FROM %s WHERE done=false ORDER BY deadline, number`, PositionsTable, OrdersTable, OrdersTable)

	if err := r.db.Select(&orders, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return orders, nil
}

func (r *OrderRepo) GetForAnalytics(ctx context.Context) (analytics []models.Analytics, err error) {
	query := fmt.Sprintf(`SELECT number, %s.date as order_start, %s.complited as order_end, %s.position, %s.title as pos_title, 
	%s.ring, %s.complited as pos_end, %s.title as oper_title, %s.date as oper_end
	FROM %s	INNER JOIN %s ON order_id=%s.id	INNER JOIN %s ON position_id=%s.id INNER JOIN %s ON operation_id=%s.id
	WHERE %s.done=true ORDER BY number, %s.position, %s.ring, %s.step_number`,
		OrdersTable, OrdersTable, PositionsTable, PositionsTable,
		PositionsTable, PositionsTable, RootOperationTable, OperationsTable,
		OrdersTable, PositionsTable, OrdersTable, OperationsTable, PositionsTable, RootOperationTable, RootOperationTable,
		OrdersTable, PositionsTable, PositionsTable, RootOperationTable,
	)

	if err := r.db.Select(&analytics, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return analytics, nil
}

func (r *OrderRepo) Find(ctx context.Context, number string) (orders []models.FindedOrder, err error) {
	query := fmt.Sprintf("SELECT id, number, done FROM %s WHERE number LIKE $1 ORDER BY done LIMIT 5", OrdersTable)

	if err := r.db.Select(&orders, query, "%"+number+"%"); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return orders, nil
}

func (r *OrderRepo) Create(ctx context.Context, order models.OrderDTO) (id uuid.UUID, err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, number, done, deadline, date) VALUES ($1, $2, $3, $4, $5) RETURNING id", OrdersTable)
	id = uuid.New()

	_, err = r.db.Exec(query, id, order.Number, order.Done, order.Deadline, order.Date)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, nil
}

func (r *OrderRepo) Update(ctx context.Context, order models.OrderDTO) error {
	query := fmt.Sprintf("UPDATE %s SET deadline=$1 WHERE id=$2", OrdersTable)

	_, err := r.db.Exec(query, order.Deadline, order.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *OrderRepo) Complete(ctx context.Context, order models.OrderDTO) error {
	query := fmt.Sprintf("UPDATE %s SET done=$1 WHERE id=$2", OrdersTable)

	_, err := r.db.Exec(query, order.Done, order.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *OrderRepo) DeleteOld(ctx context.Context, term time.Time) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE done=true AND complited<$1 AND complited != ''", OrdersTable)

	if _, err := r.db.Exec(query, term.Unix()); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Delete(ctx context.Context, order models.OrderDTO) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", OrdersTable)

	if _, err := r.db.Exec(query, order.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
