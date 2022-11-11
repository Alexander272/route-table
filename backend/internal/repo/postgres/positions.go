package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PositionRepo struct {
	db *sqlx.DB
}

func NewPositionRepo(db *sqlx.DB) *PositionRepo {
	return &PositionRepo{db: db}
}

func (r *PositionRepo) GetForOrder(ctx context.Context, orderId uuid.UUID) (positions []models.PositionForOrder, err error) {
	query := fmt.Sprintf("SELECT id, position, count, title, ring, deadline, connected, done FROM %s WHERE order_id=$1", PositionsTable)

	if err := r.db.Select(&positions, query, orderId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return positions, nil
}

func (r *PositionRepo) Get(ctx context.Context, id uuid.UUID) (position models.Position, err error) {
	query := fmt.Sprintf(`SELECT %s.id, position, count, title, ring, %s.deadline, connected, %s.done, number FROM %s
		INNER JOIN %s ON order_id=%s.id WHERE %s.id=$1`,
		PositionsTable, PositionsTable, PositionsTable, PositionsTable, OrdersTable, OrdersTable, PositionsTable)

	if err := r.db.Get(&position, query, id); err != nil {
		return position, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return position, nil
}

func (r *PositionRepo) Create(ctx context.Context, position models.PositionDTO) (id uuid.UUID, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (id, order_id, position, count, title, ring, deadline, connected) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, PositionsTable)
	id = uuid.New()

	_, err = r.db.Exec(query, id, position.OrderId, position.Position, position.Count, position.Title, position.Ring, position.Deadline, position.Connected)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, nil
}

func (r *PositionRepo) CreateFew(ctx context.Context, positions []models.PositionDTO) error {
	query := fmt.Sprintf("INSERT INTO %s (id, order_id, position, count, title, ring, deadline, connected) VALUES ", PositionsTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(positions))

	c := 8
	for i, p := range positions {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5, i*c+6, i*c+7, i*c+8))
		args = append(args, p.Id, p.OrderId, p.Position, p.Count, p.Title, p.Ring, p.Deadline, p.Connected)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PositionRepo) Update(ctx context.Context, position models.PositionDTO) error {
	query := fmt.Sprintf("UPDATE %s SET done=$1 WHERE id=$2", PositionsTable)

	_, err := r.db.Exec(query, position.Done, position.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PositionRepo) Delete(ctx context.Context, position models.PositionDTO) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", PositionsTable)

	if _, err := r.db.Exec(query, position.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
