package service

import (
	"context"

	repository "github.com/Alexander272/route-table/internal/repo"
)

type PositionService struct {
	repo      repository.Position
	operation *OperationService
}

func NewPositionService(repo repository.Position, operation *OperationService) *PositionService {
	return &PositionService{
		repo:      repo,
		operation: operation,
	}
}

// func (s *PositionService) Create(ctx context.Context, position models.PositionDTO) (id uuid.UUID, err error) {

// }

func (s *PositionService) CreateFew(ctx context.Context) error {
	return nil
}
