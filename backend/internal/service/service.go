package service

import (
	"context"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/auth"
	"github.com/google/uuid"
)

// var Template map[string]int = map[string]int{
// 	"Order":    18,
// 	"Title":    3,
// 	"Position": 6,
// 	"Count":    10,
// 	"Deadline": 5,
// }

var Template models.Template = models.Template{
	Order:    18,
	Title:    3,
	Position: 6,
	Count:    10,
	Deadline: 5,
}

type RootOperation interface {
	Get(context.Context) ([]models.RootOperation, error)
	Create(context.Context, models.RootOperationDTO) (uuid.UUID, error)
	Update(context.Context, models.RootOperationDTO) error
	Delete(context.Context, models.RootOperationDTO) error
}

type Operation interface{}

type Position interface{}

type Order interface {
	Create(ctx context.Context, order models.OrderDTO) (id uuid.UUID, err error)
	Update(ctx context.Context, order models.OrderDTO) error
	Delete(ctx context.Context, order models.OrderDTO) error
}

type Services struct {
	RootOperation
	Operation
	Position
	Order
}

type Deps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps Deps) *Services {
	rootOperation := NewRootOperationService(deps.Repos.RootOperation)
	operation := NewOperationService(deps.Repos.Operation)
	position := NewPositionService(deps.Repos.Position, operation)
	order := NewOrderService(deps.Repos.Order, position)

	return &Services{
		RootOperation: rootOperation,
		Operation:     operation,
		Position:      position,
		Order:         order,
	}
}
