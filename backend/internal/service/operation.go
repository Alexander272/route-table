package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func (s *OperationService) Get(ctx context.Context, positionId uuid.UUID, enabled pq.StringArray) (operations []models.Operation, err error) {
	operations, err = s.repo.Get(ctx, positionId, enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to get operation. error: %w", err)
	}

	return operations, nil
}

func (s *OperationService) GetConnected(ctx context.Context, positionId, operationId uuid.UUID) (operations []models.Operation, err error) {
	operations, err = s.repo.GetConnected(ctx, positionId, operationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get connected operation. error: %w", err)
	}

	return operations, nil
}

func (s *OperationService) GetWithReasons(ctx context.Context, positionId uuid.UUID) (operations []models.OperationWithReason, err error) {
	ops, err := s.repo.GetWithReasons(ctx, positionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get operation with reason. error: %w", err)
	}

	for i, owr := range ops {
		if i == 0 {
			operations = append(operations, models.OperationWithReason{
				Id:         owr.Id,
				Title:      owr.Title,
				Done:       owr.Done,
				Remainder:  owr.Remainder,
				StepNumber: owr.StepNumber,
			})
			if owr.Value != nil {
				operations[0].Reason = []models.Reason{{
					Id:    owr.ReasonId,
					Value: *owr.Value,
					Date:  *owr.Date,
				}}
			}
		} else {
			if operations[len(operations)-1].Id == owr.Id {
				operations[len(operations)-1].Reason = append(operations[len(operations)-1].Reason, models.Reason{
					Id:    owr.ReasonId,
					Value: *owr.Value,
					Date:  *owr.Date,
				})
			} else {
				operations = append(operations, models.OperationWithReason{
					Id:         owr.Id,
					Title:      owr.Title,
					Done:       owr.Done,
					Remainder:  owr.Remainder,
					StepNumber: owr.StepNumber,
				})
				if owr.Value != nil {
					operations[len(operations)-1].Reason = []models.Reason{{
						Id:    owr.ReasonId,
						Value: *owr.Value,
						Date:  *owr.Date,
					}}
				}
			}
		}
	}
	return operations, nil
}

func (s *OperationService) CreateFew(ctx context.Context, operations []models.OperationDTO) error {
	if err := s.repo.CreateFew(ctx, operations); err != nil {
		return fmt.Errorf("failed to create few operations. err: %w", err)
	}
	return nil
}

func (s *OperationService) Check(ctx context.Context, posId1, posId2 uuid.UUID, done bool, remaider int,
) (op1 models.OperationDTO, op2 models.OperationDTO, err error) {
	operations1, err := s.repo.GetAll(ctx, posId1)
	if err != nil {
		return op1, op2, fmt.Errorf("failed to get operation. error: %w", err)
	}
	operations2, err := s.repo.GetAll(ctx, posId2)
	if err != nil {
		return op1, op2, fmt.Errorf("failed to get operation. error: %w", err)
	}
	// if !operations1[len(operations1)-2].Done || !operations2[len(operations2)-2].Done {
	// 	return op1, op2, errors.New("connected operation not completed")
	// }
	op1 = models.OperationDTO{
		Id:        operations1[len(operations1)-1].Id,
		Done:      done,
		Remainder: remaider,
		Date:      time.Now().Format("02.01.2006 15:04"),
	}
	op2 = models.OperationDTO{
		Id:        operations2[len(operations2)-1].Id,
		Done:      done,
		Remainder: remaider,
		Date:      time.Now().Format("02.01.2006 15:04"),
	}

	return op1, op2, nil
}

func (s *OperationService) Complete(ctx context.Context, operation models.OperationDTO) error {
	if err := s.repo.Update(ctx, operation); err != nil {
		return fmt.Errorf("failed to complete operation. error: %w", err)
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

func (s *OperationService) DeleteSkipped(ctx context.Context, positionId, operationId uuid.UUID, count int) error {
	operations, err := s.repo.GetAll(ctx, positionId)
	if err != nil {
		return fmt.Errorf("failed to get operation. error: %w", err)
	}

	stepNumber := 0
	for _, o := range operations {
		if o.Id == operationId {
			stepNumber = o.StepNumber
		}
	}

	deleteId := make([]uuid.UUID, 0)
	for _, o := range operations {
		if o.StepNumber < stepNumber && !o.Done && o.Remainder == count {
			deleteId = append(deleteId, o.Id)
		}
	}

	if len(deleteId) > 0 {
		if err := s.repo.DeleteFew(ctx, deleteId); err != nil {
			return fmt.Errorf("failed to delete few. error: %w", err)
		}
	}

	return nil
}
