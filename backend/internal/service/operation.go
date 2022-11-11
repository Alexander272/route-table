package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
)

type OperationService struct {
	repo   repository.Operation
	reason *ReasonService
}

func NewOperationService(repo repository.Operation, reason *ReasonService) *OperationService {
	return &OperationService{
		repo:   repo,
		reason: reason,
	}
}

func (s *OperationService) Get(ctx context.Context, positionId uuid.UUID) (operatios []models.Operation, err error) {
	operatios, err = s.repo.Get(ctx, positionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get operation. error: %w", err)
	}

	return operatios, nil
}

func (s *OperationService) CreateFew(ctx context.Context, operations []models.OperationDTO) error {
	if err := s.repo.CreateFew(ctx, operations); err != nil {
		return fmt.Errorf("failed to create few operations. err: %w", err)
	}
	return nil
}

func (s *OperationService) Update(ctx context.Context, operation models.CompleteOperation) error {
	oper := models.OperationDTO{
		Id:        operation.Id,
		Done:      operation.Done,
		Remainder: operation.Remainder,
	}
	if operation.Done {
		oper.Date = time.Now().Format("02.01.2006 15:04")
	}

	if operation.Reason != "" {
		reason := models.ReasonDTO{
			OperationId: operation.Id,
			Value:       operation.Reason,
			Date:        time.Now().Format("02.01.2006 15:04"),
		}

		_, err := s.reason.Create(ctx, reason)
		if err != nil {
			return err
		}
	}

	if err := s.repo.Update(ctx, oper); err != nil {
		return fmt.Errorf("failed to updte operation. error: %w", err)
	}

	return nil
}
