package service

import (
	"context"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/auth"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
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

type Operation interface {
	Get(context.Context, uuid.UUID) ([]models.Operation, error)
	CreateFew(context.Context, []models.OperationDTO) error
	Update(context.Context, models.CompleteOperation) error
}

type Position interface {
	Get(context.Context, uuid.UUID) (models.Position, error)
	CreateFew(context.Context, map[string]uuid.UUID, [][]string) error
}

type Order interface {
	Parse(context.Context, *excelize.File) error
	Find(context.Context, string) ([]models.FindedOrder, error)
	GetWithPositions(ctx context.Context, id uuid.UUID) (order models.OrderWithPositions, err error)
	Create(context.Context, models.OrderDTO) (uuid.UUID, error)
	Update(context.Context, models.OrderDTO) error
	Delete(context.Context, models.OrderDTO) error
}

type Reason interface {
	Create(context.Context, models.ReasonDTO) (uuid.UUID, error)
}

type Services struct {
	RootOperation
	Operation
	Position
	Order
	Reason
}

type Deps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps Deps) *Services {
	rootOperation := NewRootOperationService(deps.Repos.RootOperation)
	reason := NewReasonService(deps.Repos.Reason)
	operation := NewOperationService(deps.Repos.Operation, reason)
	position := NewPositionService(deps.Repos.Position, operation, rootOperation)
	order := NewOrderService(deps.Repos.Order, position)

	return &Services{
		RootOperation: rootOperation,
		Reason:        reason,
		Operation:     operation,
		Position:      position,
		Order:         order,
	}
}
