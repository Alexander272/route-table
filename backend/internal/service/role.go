package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
)

type RoleService struct {
	repo repository.Role
}

func NewRoleService(repo repository.Role) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) Get(ctx context.Context, roleId uuid.UUID) (role models.Role, err error) {
	role, err = s.repo.Get(ctx, roleId)
	if err != nil {
		return role, fmt.Errorf("failed to get role. error: %w", err)
	}

	return role, nil
}

func (s *RoleService) GetAll(ctx context.Context) (roles []models.Role, err error) {
	roles, err = s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles. error: %w", err)
	}

	return roles, nil
}

func (s *RoleService) Create(ctx context.Context, role models.RoleDTO) (id uuid.UUID, err error) {
	id, err = s.repo.Create(ctx, role)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create role. error: %w", err)
	}

	return id, nil
}

func (s *RoleService) Update(ctx context.Context, role models.RoleDTO) error {
	if err := s.repo.Update(ctx, role); err != nil {
		return fmt.Errorf("failed to update role. error: %w", err)
	}
	return nil
}

func (s *RoleService) Delete(ctx context.Context, role models.RoleDTO) error {
	if err := s.repo.Delete(ctx, role); err != nil {
		return fmt.Errorf("failed to delete role. error: %w", err)
	}
	return nil
}
