package service

import (
	"context"
	"time"

	"github.com/Alexander272/route-table/internal/config"
	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/auth"
	"github.com/Alexander272/route-table/pkg/hasher"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

// индексы таблицы заказов
var Template models.Template = models.Template{
	Order:    18,
	Title:    3,
	Position: 6,
	Marking:  7,
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
	Get(context.Context, uuid.UUID, pq.StringArray) ([]models.Operation, error)
	GetConnected(ctx context.Context, positionId, operationId uuid.UUID) ([]models.Operation, error)
	GetWithReasons(context.Context, uuid.UUID) ([]models.OperationWithReason, error)
	CreateFew(context.Context, []models.OperationDTO) error
	Update(ctx context.Context, positionId, groupId uuid.UUID, operation models.CompleteOperation) error
	UpdateCount(context.Context, models.UpdateCount) error
	Rollback(context.Context, uuid.UUID, []string) error
}

type CompletedOperation interface {
	Get(context.Context, uuid.UUID) ([]models.CompletedOperation, error)
	Create(context.Context, models.CompleteOperation) (uuid.UUID, error)
	Delete(context.Context, models.CompletedOperation) error
}

type Position interface {
	Get(context.Context, uuid.UUID, pq.StringArray) (models.Position, error)
	GetWithReasons(context.Context, uuid.UUID) (models.PositionWithReason, error)
	CreateFew(context.Context, map[string]uuid.UUID, [][]string) error
	Action(context.Context, models.CompletePosition) error
	UpdateCount(context.Context, models.UpdateCount) error
	Rollback(context.Context, models.RollbackPosition) error
}

type Order interface {
	Parse(context.Context, *excelize.File) error
	Find(context.Context, string) ([]models.FindedOrder, error)
	GetAll(context.Context) ([]models.GroupedOrder, error)
	GetGrouped(context.Context) (models.UrgencyGroup, error)
	GetWithPositions(context.Context, uuid.UUID, string, pq.StringArray) (models.OrderWithPositions, error)
	Create(context.Context, models.OrderDTO) (uuid.UUID, error)
	Update(context.Context, models.OrderDTO) error
	Delete(context.Context, models.OrderDTO) error

	GetForAnalytics(context.Context) (*excelize.File, error)
}

type Reason interface {
	Get(context.Context) ([]models.PosWithReason, error)
	GetFile(context.Context) (*excelize.File, error)
	Create(context.Context, models.ReasonDTO) (uuid.UUID, error)
	DeleteFew(context.Context, []string) error
}

type Role interface {
	Get(context.Context, uuid.UUID) (models.Role, error)
	GetAll(context.Context) ([]models.Role, error)
	Create(context.Context, models.RoleDTO) (uuid.UUID, error)
	Update(context.Context, models.RoleDTO) error
	Delete(context.Context, models.RoleDTO) error
}

type User interface {
	GetAll(context.Context) ([]models.User, error)
	Create(context.Context, models.UserDTO) (uuid.UUID, error)
	Update(context.Context, models.UserDTO) error
	Delete(context.Context, models.UserDTO) error
}

type Session interface {
	SignIn(ctx context.Context, u models.SignIn) (models.User, string, error)
	SingOut(ctx context.Context, token string) error
	Refresh(ctx context.Context, user models.UserWithRole) (models.User, string, error)
	CheckSession(ctx context.Context, token string) (bool, error)
	TokenParse(token string) (user models.UserWithRole, err error)
}

type Urgency interface {
	Get(ctx context.Context) models.Urgency
	Change(ctx context.Context, urgency models.Urgency) error
}

type Services struct {
	RootOperation
	Operation
	CompletedOperation
	Position
	Order
	Reason
	Role
	User
	Session
	Urgency
}

type Deps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	Hasher          hasher.PasswordHasher
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Urgency         *config.UrgencyConfig
	QueryDelay      time.Duration
	OrdersTerm      time.Duration
}

func NewServices(deps Deps) *Services {
	rootOperation := NewRootOperationService(deps.Repos.RootOperation)
	reason := NewReasonService(deps.Repos.Reason)
	completed := NewCompletedOperationService(deps.Repos.CompletedOperation)
	operation := NewOperationService(deps.Repos.Operation, reason, completed)
	position := NewPositionService(deps.Repos.Position, operation, rootOperation)
	order := NewOrderService(deps.Repos.Order, position, deps.Urgency, deps.OrdersTerm, deps.QueryDelay)
	role := NewRoleService(deps.Repos.Role)
	user := NewUserService(deps.Repos.User, deps.Hasher, role)
	session := NewSessionService(deps.Repos.Session, user, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)
	urgency := NewUrgencyService(deps.Urgency)

	return &Services{
		RootOperation: rootOperation,
		Reason:        reason,
		Operation:     operation,
		Position:      position,
		Order:         order,
		Role:          role,
		User:          user,
		Session:       session,
		Urgency:       urgency,
	}
}
