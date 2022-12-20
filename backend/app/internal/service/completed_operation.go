package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
)

type CompletedOperationService struct {
	repo repository.CompletedOperation
}

func NewCompletedOperationService(repo repository.CompletedOperation) *CompletedOperationService {
	return &CompletedOperationService{repo: repo}
}

func (s *CompletedOperationService) Get(ctx context.Context, operationId uuid.UUID) (operations []models.CompletedOperation, err error) {
	operations, err = s.repo.Get(ctx, operationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed operations. error: %w", err)
	}

	return operations, nil
}

func (s *CompletedOperationService) Create(ctx context.Context, operation models.CompletedOperation) (id uuid.UUID, err error) {
	id, err = s.repo.Create(ctx, operation)
	if err != nil {
		return id, fmt.Errorf("failed to create completed operation. error: %w", err)
	}

	return id, nil
}

func (s *CompletedOperationService) CreateFew(ctx context.Context, operation []models.CompletedOperation) error {
	if err := s.repo.CreateFew(ctx, operation); err != nil {
		return fmt.Errorf("failed to create few completed operations. error: %w", err)
	}

	return nil
}

func (s *CompletedOperationService) Delete(ctx context.Context, operation models.CompletedOperation) error {
	if err := s.repo.Delete(ctx, operation); err != nil {
		return fmt.Errorf("failed to delete completed operation. error: %w", err)
	}
	return nil
}
