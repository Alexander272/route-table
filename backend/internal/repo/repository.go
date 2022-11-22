package repository

import (
	"context"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/repo/postgres"
	redisRepo "github.com/Alexander272/route-table/internal/repo/redis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RootOperation interface {
	Get(context.Context) ([]models.RootOperation, error)
	Create(context.Context, models.RootOperationDTO) (uuid.UUID, error)
	Update(context.Context, models.RootOperationDTO) error
	Delete(context.Context, models.RootOperationDTO) error
}

type Order interface {
	Get(context.Context, uuid.UUID) (models.OrderWithPositions, error)
	GetAll(ctx context.Context) ([]models.Order, error)
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
	Get(context.Context, uuid.UUID, pq.StringArray) ([]models.Operation, error)
	GetAll(context.Context, uuid.UUID) ([]models.Operation, error)
	GetWithReasons(context.Context, uuid.UUID) ([]models.OperationWithReason, error)
	GetConnected(ctx context.Context, positionId, operationId uuid.UUID) ([]models.Operation, error)
	Create(context.Context, models.OperationDTO) (uuid.UUID, error)
	CreateFew(context.Context, []models.OperationDTO) error
	Update(context.Context, models.OperationDTO) error
	Delete(context.Context, models.OperationDTO) error
	DeleteFew(context.Context, []uuid.UUID) error
}

type Reason interface {
	Create(context.Context, models.ReasonDTO) (uuid.UUID, error)
	Delete(context.Context, models.ReasonDTO) error
}

type Role interface {
	Get(context.Context, uuid.UUID) (models.Role, error)
	GetAll(context.Context) ([]models.Role, error)
	Create(context.Context, models.RoleDTO) (uuid.UUID, error)
	Update(context.Context, models.RoleDTO) error
	Delete(context.Context, models.RoleDTO) error
}

type User interface {
	Get(context.Context, models.SignIn) (models.UserWithRole, error)
	GetAll(context.Context) ([]models.User, error)
	Create(context.Context, models.UserDTO) (uuid.UUID, error)
	Update(context.Context, models.UserDTO) error
	Delete(context.Context, models.UserDTO) error
}

type Session interface {
	Create(ctx context.Context, sessionName string, data models.SessionData) error
	Get(ctx context.Context, sessionName string) (data models.SessionData, err error)
	GetDel(ctx context.Context, sessionName string) (data models.SessionData, err error)
	Remove(ctx context.Context, sessionName string) error
}

type Repositories struct {
	RootOperation
	Order
	Position
	Operation
	Reason
	Role
	User
	Session
}

func NewRepo(db *sqlx.DB, redis redis.Cmdable) *Repositories {
	return &Repositories{
		RootOperation: postgres.NewRootOperationRepo(db),
		Order:         postgres.NewOrderRepo(db),
		Position:      postgres.NewPositionRepo(db),
		Operation:     postgres.NewOperationRepo(db),
		Reason:        postgres.NewReasonRepo(db),
		Role:          postgres.NewRoleRepo(db),
		User:          postgres.NewUserRepo(db),
		Session:       redisRepo.NewSessionRepo(redis),
	}
}
