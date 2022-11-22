package service

import (
	"context"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/auth"
	"github.com/Alexander272/route-table/pkg/hasher"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

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
	Get(context.Context, uuid.UUID, pq.StringArray) ([]models.Operation, error)
	GetConnected(ctx context.Context, positionId, operationId uuid.UUID) ([]models.Operation, error)
	GetWithReasons(context.Context, uuid.UUID) ([]models.OperationWithReason, error)
	CreateFew(context.Context, []models.OperationDTO) error
	Update(context.Context, models.CompleteOperation) error
}

type Position interface {
	Get(context.Context, uuid.UUID, pq.StringArray) (models.Position, error)
	GetWithReasons(context.Context, uuid.UUID) (models.PositionWithReason, error)
	CreateFew(context.Context, map[string]uuid.UUID, [][]string) error
	Update(context.Context, models.CompletePosition) error
}

type Order interface {
	Parse(context.Context, *excelize.File) error
	Find(context.Context, string) ([]models.FindedOrder, error)
	GetAll(context.Context) ([]models.GroupedOrder, error)
	GetWithPositions(context.Context, uuid.UUID) (models.OrderWithPositions, error)
	Create(context.Context, models.OrderDTO) (uuid.UUID, error)
	Update(context.Context, models.OrderDTO) error
	Delete(context.Context, models.OrderDTO) error
}

type Reason interface {
	Create(context.Context, models.ReasonDTO) (uuid.UUID, error)
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
	SingOut(ctx context.Context, userId string) error
	Refresh(ctx context.Context, user models.UserWithRole) (models.User, string, error)
	CheckSession(ctx context.Context, u models.UserWithRole, token string) (bool, error)
	TokenParse(token string) (user models.UserWithRole, err error)
}

type Services struct {
	RootOperation
	Operation
	Position
	Order
	Reason
	Role
	User
	Session
}

type Deps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	Hasher          hasher.PasswordHasher
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	UrgencyHigh     time.Duration
	UrgencyMid      time.Duration
}

func NewServices(deps Deps) *Services {
	rootOperation := NewRootOperationService(deps.Repos.RootOperation)
	reason := NewReasonService(deps.Repos.Reason)
	operation := NewOperationService(deps.Repos.Operation, reason)
	position := NewPositionService(deps.Repos.Position, operation, rootOperation)
	order := NewOrderService(deps.Repos.Order, position, deps.UrgencyHigh, deps.UrgencyMid)
	role := NewRoleService(deps.Repos.Role)
	user := NewUserService(deps.Repos.User, deps.Hasher, role)
	session := NewSessionService(deps.Repos.Session, user, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)

	return &Services{
		RootOperation: rootOperation,
		Reason:        reason,
		Operation:     operation,
		Position:      position,
		Order:         order,
		Role:          role,
		User:          user,
		Session:       session,
	}
}
