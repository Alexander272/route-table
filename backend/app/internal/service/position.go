package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PositionService struct {
	repo          repository.Position
	operation     *OperationService
	rootOperation *RootOperationService
}

func NewPositionService(repo repository.Position, operation *OperationService, root *RootOperationService) *PositionService {
	return &PositionService{
		repo:          repo,
		operation:     operation,
		rootOperation: root,
	}
}

var gasket map[string][]string = map[string][]string{
	"СНП-Д": {"внутрен.", "наружное"},
	"СНП-Г": {"наружное"},
	"СНП-В": {"внутрен."},
	"СНП-Б": {""},
	"СНП-А": {""},
}

// Создание позиций
// каждая позиция привязывается к заказу, а также некоторые позици имеют два экземпляра и поэтому они связываются
// между собой через поле Connected
func (s *PositionService) CreateFew(ctx context.Context, orders map[string]uuid.UUID, positions [][]string) error {
	operations := make([]models.OperationDTO, 0, 200)
	pos := make([]models.PositionDTO, 0, len(positions))

	root, err := s.rootOperation.Get(ctx)
	if err != nil {
		return err
	}

	for _, row := range positions {
		rings := gasket[row[0]]
		for _, v := range rings {
			id := uuid.New()
			position, err := strconv.Atoi(row[Template.Position])
			if err != nil {
				return fmt.Errorf("failed to convert position. error: %w", err)
			}
			count, err := strconv.Atoi(row[Template.Count])
			if err != nil {
				return fmt.Errorf("failed to convert count. error: %w", err)
			}

			pos = append(pos, models.PositionDTO{
				Id:       id,
				OrderId:  orders[row[Template.Order]],
				Position: position,
				Count:    count,
				Title:    row[Template.Title],
				Ring:     v,
			})

			for _, r := range root {
				if r.Gasket == row[1] {
					opId := uuid.New()
					operations = append(operations, models.OperationDTO{
						Id:          opId,
						PositionId:  id,
						OperationId: r.Id,
						Remainder:   count,
					})
				}
			}
		}
		if len(rings) == 2 {
			pos[len(pos)-1].Connected = pos[len(pos)-2].Id
			pos[len(pos)-2].Connected = pos[len(pos)-1].Id
		}
	}

	if err := s.repo.CreateFew(ctx, pos); err != nil {
		return fmt.Errorf("failed to create few positions. error: %w", err)
	}
	if err := s.operation.CreateFew(ctx, operations); err != nil {
		return err
	}

	return nil
}

// получение позиций с операциями
func (s *PositionService) Get(ctx context.Context, positionId uuid.UUID, enabled pq.StringArray) (position models.Position, err error) {
	position, err = s.repo.Get(ctx, positionId)
	if err != nil {
		return models.Position{}, fmt.Errorf("failed to get position. error: %w", err)
	}

	operation, err := s.operation.Get(ctx, positionId, enabled)
	if err != nil {
		return models.Position{}, err
	}
	position.Operation = operation

	return position, nil
}

// получение позиций с опрерациями и причинами
func (s *PositionService) GetWithReasons(ctx context.Context, positionId uuid.UUID) (position models.PositionWithReason, err error) {
	pos, err := s.repo.Get(ctx, positionId)
	if err != nil {
		return models.PositionWithReason{}, fmt.Errorf("failed to get position. error: %w", err)
	}
	position = models.PositionWithReason{
		Id:        pos.Id,
		Order:     pos.Order,
		Position:  pos.Position,
		Count:     pos.Count,
		Title:     pos.Title,
		Ring:      pos.Ring,
		Done:      pos.Done,
		Deadline:  pos.Deadline,
		Connected: pos.Connected,
	}

	operation, err := s.operation.GetWithReasons(ctx, positionId)
	if err != nil {
		return models.PositionWithReason{}, err
	}
	position.Operation = operation

	return position, nil
}

// получение списка позиций для заказа (готовые позиции отсутствуют)
func (s *PositionService) GetForOrder(ctx context.Context, orderId uuid.UUID, enabled pq.StringArray) (positions []models.PositionForOrder, err error) {
	positions, err = s.repo.GetForOrder(ctx, orderId, enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions for order. error: %w", err)
	}

	if len(positions) == 0 {
		return positions, nil
	}

	t, err := strconv.Atoi(positions[0].Deadline)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deadline. error: %w", err)
	}
	deadline := time.Unix(int64(t), 0).Format("02.01.2006")
	for i := range positions {
		positions[i].Deadline = deadline
	}

	return positions, nil
}

// получение списка позиций для заказа вместе с готовыми позициями
func (s *PositionService) GetFullForOrder(ctx context.Context, orderId uuid.UUID) (positions []models.PositionForOrder, err error) {
	positions, err = s.repo.GetFullForOrder(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get full positions for order. error: %w", err)
	}

	t, err := strconv.Atoi(positions[0].Deadline)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deadline. error: %w", err)
	}
	deadline := time.Unix(int64(t), 0).Format("02.01.2006")
	for i := range positions {
		positions[i].Deadline = deadline
	}

	return positions, nil
}

// закрытие (выполнение) позиции
func (s *PositionService) Complete(ctx context.Context, position models.PositionDTO) error {
	if err := s.repo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to complete position. error: %w", err)
	}
	return nil
}

// закрытие операций и позиции, если операция финальная
func (s *PositionService) Update(ctx context.Context, position models.CompletePosition) error {
	if position.IsFinish {
		var op1, op2 models.OperationDTO
		if position.Connected != uuid.Nil {
			var err error
			op1, op2, err = s.operation.Check(ctx, position.Id, position.Connected, position.Operation.Done, position.Operation.Remainder)
			if err != nil {
				return err
			}

			if err := s.operation.Complete(ctx, op2); err != nil {
				return err
			}
			if op2.Done {
				comlited := fmt.Sprintf("%d", time.Now().Unix())
				if err := s.Complete(ctx, models.PositionDTO{Id: position.Connected, Done: op2.Done, Complited: comlited}); err != nil {
					return err
				}
			}
		} else {
			op1 = models.OperationDTO{
				Id:         position.Operation.Id,
				PositionId: position.Id,
				Done:       position.Operation.Done,
				Remainder:  position.Operation.Remainder,
				Date:       time.Now().Format("02.01.2006 15:04"),
			}
		}

		if err := s.operation.Complete(ctx, op1); err != nil {
			return err
		}
		if op1.Done {
			comlited := fmt.Sprintf("%d", time.Now().Unix())
			if err := s.Complete(ctx, models.PositionDTO{Id: position.Id, Done: op1.Done, Complited: comlited}); err != nil {
				return err
			}
		}

		return nil
	}

	connected, err := s.operation.GetConnected(ctx, position.Id, position.Operation.Id)
	if err != nil {
		return err
	}

	if err := s.operation.Update(ctx, position.Operation); err != nil {
		return err
	}
	if len(connected) > 0 {
		for _, o := range connected {
			remainder := o.Remainder - position.Operation.Count
			operation := models.CompleteOperation{
				Id:        o.Id,
				Done:      remainder == 0,
				Remainder: remainder,
				Count:     position.Operation.Count,
			}
			if err := s.operation.Update(ctx, operation); err != nil {
				return err
			}
		}
	}

	// if err := s.operation.DeleteSkipped(ctx, position.Id, position.Operation.Id, position.Count); err != nil {
	// 	return err
	// }

	return nil
}
