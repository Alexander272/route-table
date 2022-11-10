package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
)

type RootOperationService struct {
	repo repository.RootOperation
}

func NewRootOperationService(repo repository.RootOperation) *RootOperationService {
	return &RootOperationService{repo: repo}
}

func (s *RootOperationService) Get(ctx context.Context) (operation []models.RootOperation, err error) {
	operation, err = s.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all root operation. error: %w", err)
	}

	return operation, nil
}

func (s *RootOperationService) Create(ctx context.Context, operation models.RootOperationDTO) (id uuid.UUID, err error) {
	id, err = s.repo.Create(ctx, operation)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create root operation. error: %w", err)
	}

	return id, nil
}

func (s *RootOperationService) Update(ctx context.Context, operation models.RootOperationDTO) error {
	if err := s.repo.Update(ctx, operation); err != nil {
		return fmt.Errorf("failed to update root operation. error: %w", err)
	}
	return nil
}

func (s *RootOperationService) Delete(ctx context.Context, operation models.RootOperationDTO) error {
	if err := s.repo.Delete(ctx, operation); err != nil {
		return fmt.Errorf("failed to delete root operation. error: %w", err)
	}
	return nil
}
