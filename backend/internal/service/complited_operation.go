package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
)

type ComplitedOperationService struct {
	repo repository.ComplitedOperation
}

func NewComplitedOperationService(repo repository.ComplitedOperation) *ComplitedOperationService {
	return &ComplitedOperationService{repo: repo}
}

func (s *ComplitedOperationService) Create(ctx context.Context, operation models.CompleteOperation) (id uuid.UUID, err error) {
	id, err = s.repo.Create(ctx, operation)
	if err != nil {
		return id, fmt.Errorf("failed to create complited operation. error: %w", err)
	}

	return id, nil
}

func (s *ComplitedOperationService) Delete(ctx context.Context, operation models.ComplitedOperation) error {
	if err := s.repo.Delete(ctx, operation); err != nil {
		return fmt.Errorf("failed to delete complited operation. error: %w", err)
	}
	return nil
}
