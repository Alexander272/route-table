package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
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
}

// func (s *PositionService) Create(ctx context.Context, position models.PositionDTO) (id uuid.UUID, err error) {

// }

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
				return err
			}
			count, err := strconv.Atoi(row[Template.Count])
			if err != nil {
				return err
			}

			pos = append(pos, models.PositionDTO{
				Id:       id,
				OrderId:  orders[row[Template.Order]],
				Position: position,
				Count:    count,
				Title:    row[Template.Title],
				Ring:     v,
				Deadline: row[Template.Deadline],
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

func (s *PositionService) Get(ctx context.Context, positionId uuid.UUID) (position models.Position, err error) {
	position, err = s.repo.Get(ctx, positionId)
	if err != nil {
		return models.Position{}, fmt.Errorf("failed to get position. error: %w", err)
	}

	operation, err := s.operation.Get(ctx, positionId)
	if err != nil {
		return models.Position{}, err
	}
	position.Operation = operation

	return position, nil
}

func (s *PositionService) GetForOrder(ctx context.Context, orderId uuid.UUID) (positions []models.PositionForOrder, err error) {
	positions, err = s.repo.GetForOrder(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions for order. error: %w", err)
	}

	return positions, nil
}
