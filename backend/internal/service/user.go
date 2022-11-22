package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/hasher"
	"github.com/google/uuid"
)

type UserService struct {
	repo   repository.User
	hasher hasher.PasswordHasher
	role   *RoleService
}

func NewUserService(repo repository.User, hasher hasher.PasswordHasher, role *RoleService) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
		role:   role,
	}
}

func (s *UserService) Get(ctx context.Context, u models.SignIn) (user models.UserWithRole, err error) {
	user, err = s.repo.Get(ctx, u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, err
		}
		return user, fmt.Errorf("failed to get user. error: %w", err)
	}

	parts := strings.Split(user.Password, ".")
	salt := parts[1]
	password := parts[0]

	pass, err := s.hasher.Hash(u.Password, salt)
	if err != nil {
		return user, fmt.Errorf("failed to hash password. error: %w", err)
	}

	if pass != password {
		return user, models.ErrPassword
	}

	user.Role, err = s.role.Get(ctx, user.RoleId)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) GetAll(ctx context.Context) (users []models.User, err error) {
	users, err = s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users. error: %w", err)
	}

	return users, nil
}

func (s *UserService) Create(ctx context.Context, user models.UserDTO) (id uuid.UUID, err error) {
	salt, err := s.hasher.GenerateSalt()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create salt. error: %w", err)
	}
	pass, err := s.hasher.Hash(user.Password, salt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password. error: %w", err)
	}
	user.Password = fmt.Sprintf("%s.%s", pass, salt)

	id, err = s.repo.Create(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user. error: %w", err)
	}

	return id, nil
}

func (s *UserService) Update(ctx context.Context, user models.UserDTO) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user. error: %w", err)
	}
	return nil
}

func (s *UserService) Delete(ctx context.Context, user models.UserDTO) error {
	if err := s.repo.Delete(ctx, user); err != nil {
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return nil
}
