package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/logger"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type OrderService struct {
	repo     repository.Order
	position *PositionService
}

func NewOrderService(repo repository.Order, position *PositionService) *OrderService {
	return &OrderService{
		repo:     repo,
		position: position,
	}
}

type AvailablePositions struct {
	Title string
	Type  string
}

var availablePositions []AvailablePositions = []AvailablePositions{
	{Title: "СНП-Д", Type: "СНП"},
	{Title: "СНП-Г", Type: "СНП"},
	{Title: "СНП-В", Type: "СНП"},
}

func (s *OrderService) Parse(ctx context.Context, file *excelize.File) error {
	orders := make(map[string]uuid.UUID, 0)
	positions := make([][]string, 0, 200)

	sheetName := file.GetSheetName(file.GetActiveSheetIndex())

	rows, err := file.Rows(sheetName)
	if err != nil {
		return err
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			logger.Error(err)
		}

		if len(row) < 18 {
			continue
		}
		if row[Template.Count] == "поз. и кол-во в ед.заказа" || row[Template.Position] == "" {
			continue
		}

		for _, v := range availablePositions {
			if strings.Contains(row[Template.Title], v.Title) {
				positions = append(positions, row)

				parts := strings.Split(row[Template.Order], " ")
				_, ok := orders[parts[2]]
				if !ok {
					id, err := s.Create(ctx, models.OrderDTO{Number: parts[2], Deadline: row[Template.Deadline], Date: parts[4]})
					if err != nil {
						return err
					}
					orders[parts[2]] = id
				}
				row[Template.Order] = parts[2]
				row[0] = v.Title
				row[1] = v.Type
			}
		}
	}
	if err = rows.Close(); err != nil {
		logger.Error(err)
	}

	if err := s.position.CreateFew(ctx, orders, positions); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) Find(ctx context.Context, number string) (orders []models.FindedOrder, err error) {
	orders, err = s.repo.Find(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to find orders. error: %w", err)
	}

	return orders, nil
}

func (s *OrderService) GetWithPositions(ctx context.Context, id uuid.UUID) (order models.OrderWithPositions, err error) {
	order, err = s.repo.Get(ctx, id)
	if err != nil {
		return models.OrderWithPositions{}, fmt.Errorf("failed to get order. error: %w", err)
	}

	positions, err := s.position.GetForOrder(ctx, id)
	if err != nil {
		return models.OrderWithPositions{}, err
	}
	order.Positions = positions

	return order, nil
}

func (s *OrderService) Create(ctx context.Context, order models.OrderDTO) (id uuid.UUID, err error) {
	id, err = s.repo.Create(ctx, order)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create order. error: %w", err)
	}

	return id, nil
}

func (s *OrderService) Update(ctx context.Context, order models.OrderDTO) error {
	if err := s.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order. error: %w", err)
	}
	return nil
}

func (s *OrderService) Delete(ctx context.Context, order models.OrderDTO) error {
	if err := s.repo.Delete(ctx, order); err != nil {
		return fmt.Errorf("failed to delete order. error: %w", err)
	}
	return nil
}
