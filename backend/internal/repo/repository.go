package repository

import (
	"context"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/repo/postgres"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RootOperation interface {
	Get(context.Context) ([]models.RootOperation, error)
	Create(context.Context, models.RootOperationDTO) (uuid.UUID, error)
	Update(context.Context, models.RootOperationDTO) error
	Delete(context.Context, models.RootOperationDTO) error
}

type Order interface {
	Get(context.Context, uuid.UUID) (models.OrderWithPositions, error)
	Find(context.Context, string) ([]models.FindedOrder, error)
	Create(context.Context, models.OrderDTO) (uuid.UUID, error)
	Update(context.Context, models.OrderDTO) error
	Delete(context.Context, models.OrderDTO) error
}

type Position interface {
	GetForOrder(context.Context, uuid.UUID) ([]models.PositionForOrder, error)
	Get(context.Context, uuid.UUID) (models.Position, error)
	Create(context.Context, models.PositionDTO) (uuid.UUID, error)
	CreateFew(context.Context, []models.PositionDTO) error
	Update(context.Context, models.PositionDTO) error
	Delete(context.Context, models.PositionDTO) error
}

type Operation interface {
	Get(context.Context, uuid.UUID) ([]models.Operation, error)
	Create(context.Context, models.OperationDTO) (uuid.UUID, error)
	CreateFew(context.Context, []models.OperationDTO) error
	Update(context.Context, models.OperationDTO) error
	Delete(context.Context, models.OperationDTO) error
}

type Reason interface {
	Create(context.Context, models.ReasonDTO) (uuid.UUID, error)
	Delete(context.Context, models.ReasonDTO) error
}

type Repositories struct {
	RootOperation
	Order
	Position
	Operation
	Reason
}

func NewRepo(db *sqlx.DB, redis redis.Cmdable) *Repositories {
	return &Repositories{
		RootOperation: postgres.NewRootOperationRepo(db),
		Order:         postgres.NewOrderRepo(db),
		Position:      postgres.NewPositionRepo(db),
		Operation:     postgres.NewOperationRepo(db),
		Reason:        postgres.NewReasonRepo(db),
	}
}
