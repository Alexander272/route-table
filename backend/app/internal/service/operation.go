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
	repo      repository.Operation
	completed *CompletedOperationService
	reason    *ReasonService
}

func NewOperationService(repo repository.Operation, reason *ReasonService, completed *CompletedOperationService) *OperationService {
	return &OperationService{
		repo:      repo,
		completed: completed,
		reason:    reason,
	}
}

// Получение операций
func (s *OperationService) Get(ctx context.Context, positionId uuid.UUID, enabled pq.StringArray) (operations []models.Operation, err error) {
	operations, err = s.repo.Get(ctx, positionId, enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to get operation. error: %w", err)
	}

	return operations, nil
}

// получение связанных операций
func (s *OperationService) GetConnected(ctx context.Context, positionId, operationId uuid.UUID) (operations []models.Operation, err error) {
	operations, err = s.repo.GetConnected(ctx, positionId, operationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get connected operation. error: %w", err)
	}

	return operations, nil
}

// получение операций вместе с причинами (причины группируются)
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
				IsFinish:   owr.IsFinish,
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
					IsFinish:   owr.IsFinish,
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

// создание операций
func (s *OperationService) CreateFew(ctx context.Context, operations []models.OperationDTO) error {
	if err := s.repo.CreateFew(ctx, operations); err != nil {
		return fmt.Errorf("failed to create few operations. err: %w", err)
	}
	return nil
}

// Раньше была проверка на выполнение. Теперь просто получение связанных операций
func (s *OperationService) Check(ctx context.Context, posId1, posId2 uuid.UUID, done bool, remainder int,
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
		Id:         operations1[len(operations1)-1].Id,
		PositionId: posId1,
		Done:       done,
		Remainder:  remainder,
		Date:       time.Now().Format("02.01.2006 15:04"),
	}
	op2 = models.OperationDTO{
		Id:         operations2[len(operations2)-1].Id,
		PositionId: posId2,
		Done:       done,
		Remainder:  remainder,
		Date:       time.Now().Format("02.01.2006 15:04"),
	}

	return op1, op2, nil
}

// закрытие текущей операции и всех не закрытых для позиции (вызывается только когда выполняется финальная операция)
func (s *OperationService) Complete(ctx context.Context, operation models.OperationDTO) error {
	if err := s.repo.Update(ctx, operation); err != nil {
		return fmt.Errorf("failed to complete operation. error: %w", err)
	}
	if err := s.repo.CompleteSkipped(ctx, operation); err != nil {
		return fmt.Errorf("failed to complete skipped operations. error: %w", err)
	}
	return nil
}

// добавление причины, при наличии, обновление операций (связанных и всех предыдущих, если операция выполнена) и создание записей для отмены
func (s *OperationService) Update(ctx context.Context, positionId, groupId uuid.UUID, operation models.CompleteOperation) error {
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

	if operation.Id == uuid.Nil {
		newOperation, err := s.repo.GetLast(ctx, positionId)
		if err != nil {
			return fmt.Errorf("failed to get operation. error: %w", err)
		}
		operation.Id = newOperation.Id
	}

	connected, err := s.GetConnected(ctx, positionId, operation.Id)
	if err != nil {
		return err
	}
	var skipped []models.Operation
	if operation.Done {
		skipped, err = s.repo.GetSkipped(ctx, positionId, operation.Id)
		if err != nil {
			return fmt.Errorf("failed to get skipped operations. error: %w", err)
		}
	}

	connected = append(connected, models.Operation{
		Id:        operation.Id,
		Remainder: operation.Count,
	})

	var operations []models.OperationDTO
	var completed []models.CompletedOperation

	for i, o := range connected {
		operations = append(operations, models.OperationDTO{
			Id:        o.Id,
			Done:      operation.Done,
			Remainder: operation.Remainder,
		})
		if operation.Done {
			operations[i].Date = time.Now().Format("02.01.2006 15:04")
		}
		completed = append(completed, models.CompletedOperation{
			Id:          o.Id,
			OperationId: o.Id,
			GroupId:     groupId,
			Remainder:   operation.Remainder,
			Count:       o.Remainder,
		})
	}
	for _, o := range skipped {
		operations = append(operations, models.OperationDTO{
			Id:        o.Id,
			Done:      operation.Done,
			Remainder: operation.Remainder,
			Date:      time.Now().Format("02.01.2006 15:04"),
		})
		completed = append(completed, models.CompletedOperation{
			Id:          o.Id,
			OperationId: o.Id,
			GroupId:     groupId,
			Remainder:   operation.Remainder,
			Count:       o.Remainder,
		})
	}

	if err := s.repo.UpdateFew(ctx, operations); err != nil {
		return fmt.Errorf("failed to update few operations. error: %w", err)
	}
	if err := s.completed.CreateFew(ctx, completed); err != nil {
		return err
	}

	return nil
}

// обновление количества
func (s *OperationService) UpdateCount(ctx context.Context, position models.UpdateCount) error {
	operations, err := s.repo.GetAll(ctx, position.Id)
	if err != nil {
		return fmt.Errorf("failed to get operations. error: %w", err)
	}

	var updatedOperations []models.OperationDTO
	addRemainder := position.Count - operations[len(operations)-1].Remainder
	for _, o := range operations {
		updatedOperations = append(updatedOperations, models.OperationDTO{
			Id:        o.Id,
			Done:      false,
			Remainder: o.Remainder + addRemainder,
			Date:      "",
		})
	}

	if err := s.repo.UpdateFew(ctx, updatedOperations); err != nil {
		return fmt.Errorf("failed to update few operations. error: %w", err)
	}

	return nil
}

// откат операции или группы операций
func (s *OperationService) Rollback(ctx context.Context, operationId uuid.UUID, reasons []string) error {
	if len(reasons) > 0 {
		if err := s.reason.DeleteFew(ctx, reasons); err != nil {
			return err
		}
	}

	completed, err := s.completed.Get(ctx, operationId)
	if err != nil {
		return err
	}

	if len(completed) == 0 {
		return fmt.Errorf("completed operation not found")
	}

	var operations []models.OperationDTO
	for _, c := range completed {
		operations = append(operations, models.OperationDTO{
			Id:        c.OperationId,
			Done:      false,
			Remainder: c.Remainder + c.Count,
			Date:      "",
		})
	}

	if err := s.repo.UpdateFew(ctx, operations); err != nil {
		return fmt.Errorf("failed to update few operations. error: %w", err)
	}

	if err := s.completed.Delete(ctx, models.CompletedOperation{GroupId: completed[0].GroupId}); err != nil {
		return err
	}

	return nil
}
