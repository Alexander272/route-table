package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/logger"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type OrderService struct {
	repo        repository.Order
	position    *PositionService
	urgencyHigh time.Duration
	urgencyMid  time.Duration
	ordersTerm  time.Duration
	queryDelay  time.Duration
	queryTime   time.Time
}

func NewOrderService(repo repository.Order, position *PositionService, urgencyHigh, urgencyMid, ordersTerm, queryDelay time.Duration) *OrderService {
	return &OrderService{
		repo:        repo,
		position:    position,
		urgencyHigh: urgencyHigh,
		urgencyMid:  urgencyMid,
		ordersTerm:  ordersTerm,
		queryDelay:  queryDelay,
		queryTime:   time.Now(),
	}
}

type AvailablePositions struct {
	Title string
	Type  string
}

// Прокладки которые будет вносится в базу
var availablePositions []AvailablePositions = []AvailablePositions{
	{Title: "СНП-Д", Type: "СНП"},
	{Title: "СНП-Г", Type: "СНП"},
	{Title: "СНП-В", Type: "СНП"},
	{Title: "СНП-Б", Type: "СНП"},
	{Title: "СНП-А", Type: "СНП"},
}

//* Загрузка заказа
// Пробегемся по всем позициям и добавляем их в массив, который после пердаем в функцию для создания позиций,
// если номера заказа нет мы его создаем
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
		if row[Template.Position] == "поз. и кол-во в ед.заказа" || row[Template.Position] == "" || row[Template.Count] == "" {
			continue
		}

		for _, v := range availablePositions {
			if strings.Contains(row[Template.Title], v.Title) {
				positions = append(positions, row)

				parts := strings.Split(row[Template.Order], " ")
				_, ok := orders[parts[2]]
				if !ok {
					deadline, err := time.Parse("02.01.2006", row[Template.Deadline])
					if err != nil {
						return fmt.Errorf("failed to parse date of deadline. error: %w", err)
					}
					id, err := s.Create(ctx, models.OrderDTO{Number: parts[2], Deadline: fmt.Sprintf("%d", deadline.Unix()), Date: parts[4]})
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

// Получеам список всех заказов, добавляем срочность и группируем по дате отгрузки
func (s *OrderService) GetAll(ctx context.Context) (orders []models.GroupedOrder, err error) {
	o, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders. error: %w", err)
	}

	if len(o) == 0 {
		return []models.GroupedOrder{}, nil
	}

	groupId := uuid.New()
	deadline, err := strconv.Atoi(o[0].Deadline)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date of deadline. error: %w", err)
	}
	date := time.Unix(int64(deadline), 0)

	var urgency string
	if time.Until(date) <= s.urgencyHigh {
		urgency = "Высокая"
	} else if time.Until(date) <= s.urgencyMid {
		urgency = "Средняя"
	} else {
		urgency = "Обычная"
	}

	orders = append(orders, models.GroupedOrder{
		Id:       groupId,
		Deadline: date.Format("02.01.2006"),
		Urgency:  urgency,
		Orders: []models.Order{{
			Id:       o[0].Id,
			Number:   strings.TrimLeft(o[0].Number, "0"),
			Done:     o[0].Done,
			Date:     o[0].Date,
			Progress: math.Round(o[0].Progress*1000) / 10,
		}},
	})

	for i, o := range o {
		if i == 0 {
			continue
		}

		deadline, err := strconv.Atoi(o.Deadline)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date of deadline. error: %w", err)
		}
		date := time.Unix(int64(deadline), 0)

		if date.Format("02.01.2006") == orders[len(orders)-1].Deadline {
			orders[len(orders)-1].Orders = append(orders[len(orders)-1].Orders, models.Order{
				Id:       o.Id,
				Number:   strings.TrimLeft(o.Number, "0"),
				Done:     o.Done,
				Date:     o.Date,
				Progress: math.Round(o.Progress*1000) / 10,
			})
		} else {
			groupId := uuid.New()

			if time.Until(date) <= s.urgencyHigh {
				urgency = "Высокая"
			} else if time.Until(date) <= s.urgencyMid {
				urgency = "Средняя"
			} else {
				urgency = "Обычная"
			}

			orders = append(orders, models.GroupedOrder{
				Id:       groupId,
				Deadline: date.Format("02.01.2006"),
				Urgency:  urgency,
				Orders: []models.Order{{
					Id:       o.Id,
					Number:   strings.TrimLeft(o.Number, "0"),
					Done:     o.Done,
					Date:     o.Date,
					Progress: math.Round(o.Progress*1000) / 10,
				}},
			})
		}
	}

	if time.Since(s.queryTime) >= s.queryDelay {
		s.queryTime = time.Now()
		if err := s.DeleteOld(ctx); err != nil {
			logger.Error(err)
		}
	}

	return orders, nil
}

// Группируем все заказы по срочности
func (s *OrderService) GetGrouped(ctx context.Context) (group models.UrgencyGroup, err error) {
	orders, err := s.GetAll(ctx)
	if err != nil {
		return models.UrgencyGroup{}, err
	}

	for _, o := range orders {
		if o.Urgency == "Высокая" {
			group.High = append(group.High, o)
		}
		if o.Urgency == "Средняя" {
			group.Middle = append(group.Middle, o)
		}
		if o.Urgency == "Обычная" {
			group.Low = append(group.Low, o)
		}
	}

	return group, nil
}

// Получение заказа с позициями
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

func (s *OrderService) DeleteOld(ctx context.Context) error {
	logger.Info("delete old order")

	term := time.Now().Add(-s.ordersTerm)
	if err := s.repo.DeleteOld(ctx, term); err != nil {
		return fmt.Errorf("failed to delete old order. error: %w", err)
	}
	return nil
}

func (s *OrderService) Delete(ctx context.Context, order models.OrderDTO) error {
	if err := s.repo.Delete(ctx, order); err != nil {
		return fmt.Errorf("failed to delete order. error: %w", err)
	}
	return nil
}
