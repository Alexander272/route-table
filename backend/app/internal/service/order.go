package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Alexander272/route-table/internal/config"
	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/logger"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

type OrderService struct {
	repo       repository.Order
	position   *PositionService
	urgency    *config.UrgencyConfig
	ordersTerm time.Duration
	queryDelay time.Duration
	queryTime  time.Time
}

func NewOrderService(repo repository.Order, position *PositionService, urgency *config.UrgencyConfig, ordersTerm, queryDelay time.Duration) *OrderService {
	return &OrderService{
		repo:       repo,
		position:   position,
		urgency:    urgency,
		ordersTerm: ordersTerm,
		queryDelay: queryDelay,
		queryTime:  time.Now(),
	}
}

type AvailablePositions struct {
	Title string
	Type  string
}

// Прокладки которые будет вносится в базу
var availablePositions []AvailablePositions = []AvailablePositions{
	{Title: "СНП-Д", Type: "СНП-Д"},
	{Title: "СНП-Г", Type: "СНП-Г"},
	{Title: "СНП-В", Type: "СНП-В"},
	{Title: "СНП-Б", Type: "СНП-Б"},
	{Title: "СНП-А", Type: "СНП-А"},
}

// * Загрузка заказа
// Пробегаемся по всем позициям и добавляем их в массив, который после передаем в функцию для создания позиций,
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

		// сравниваем название текущей позиции со списком доступных позиций
		for _, v := range availablePositions {
			if strings.Contains(row[Template.Title], v.Title) {
				parts := strings.Split(row[Template.Order], " ")
				// смотрим создавали ли мы такой заказ, если нет то создаем и записываем его id
				_, ok := orders[parts[2]]
				if !ok {
					deadline, err := time.Parse("02.01.2006", row[Template.Deadline])
					if err != nil {
						return fmt.Errorf("failed to parse date of deadline. error: %w", err)
					}
					id, err := s.Create(ctx, models.OrderDTO{
						Number:   parts[2],
						Deadline: fmt.Sprintf("%d", deadline.Unix()),
						Date:     fmt.Sprintf("%s %s", parts[4], parts[5]),
					})
					if err != nil {
						return err
					}
					orders[parts[2]] = id
				}
				row[Template.Order] = parts[2]
				row[0] = v.Title
				row[1] = v.Type

				positions = append(positions, row)
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

// Выборка по части номера до 5 заказов
func (s *OrderService) Find(ctx context.Context, number string) (orders []models.FindedOrder, err error) {
	orders, err = s.repo.Find(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to find orders. error: %w", err)
	}

	return orders, nil
}

// Получаем список всех заказов, добавляем срочность и группируем по дате отгрузки
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
	if time.Until(date) <= s.urgency.High {
		urgency = "Высокая"
	} else if time.Until(date) <= s.urgency.Middle {
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
			Progress: math.Round(o[0].Progress * 100),
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
				Progress: math.Round(o.Progress * 100),
			})
		} else {
			groupId := uuid.New()

			if time.Until(date) <= s.urgency.High {
				urgency = "Высокая"
			} else if time.Until(date) <= s.urgency.Middle {
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
					Progress: math.Round(o.Progress * 100),
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
func (s *OrderService) GetWithPositions(ctx context.Context, id uuid.UUID, role string, enabled pq.StringArray) (order models.OrderWithPositions, err error) {
	order, err = s.repo.Get(ctx, id)
	if err != nil {
		return models.OrderWithPositions{}, fmt.Errorf("failed to get order. error: %w", err)
	}

	var positions []models.PositionForOrder
	if role != "master" {
		positions, err = s.position.GetForOrder(ctx, id, enabled)
		if err != nil {
			return models.OrderWithPositions{}, err
		}
	} else {
		positions, err = s.position.GetFullForOrder(ctx, id)
		if err != nil {
			return models.OrderWithPositions{}, err
		}
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

// Удаление старых выполненных заказов
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

// Получение файла excel со статистикой по заказам
func (s *OrderService) GetForAnalytics(ctx context.Context) (file *excelize.File, err error) {
	analytics, err := s.repo.GetForAnalytics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics. error: %w", err)
	}

	file = excelize.NewFile()
	sheet := "Sheet1"
	columnNames := []string{"№ Заказа", "Позиция", "Наименование", "Кольцо", "Операция", "Дата заказа", "Дата выполнения", "Всего дней изготовления"}
	columnAxis := map[string]string{
		"number":    "A",
		"position":  "B",
		"title":     "c",
		"ring":      "D",
		"operation": "E",
		"dateStart": "F",
		"dateEnd":   "G",
		"term":      "H",
	}

	numberStyle, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "ffffff",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order style. error: %s", err)
	}
	numberStyle1, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "dffffb",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"dffffb"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order style. error: %s", err)
	}

	orderStyle, err := file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"aaff96"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order style. error: %s", err)
	}

	posStyle, err := file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"dffffb"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create position style. error: %s", err)
	}

	err = file.SetSheetRow(sheet, "A1", &columnNames)
	if err != nil {
		return nil, fmt.Errorf("failed set title row. error: %w", err)
	}

	curNum := 2
	for i, a := range analytics {
		if i == 0 {
			file.SetRowStyle(sheet, curNum, curNum, orderStyle)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), a.Number)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateStart"], curNum), a.OrderStart)
			orderEnd := ""
			term := ""
			if a.OrderEnd != "" {
				t, err := strconv.Atoi(a.OrderEnd)
				if err != nil {
					return nil, fmt.Errorf("failed to parse order date. error: %w", err)
				}
				end := time.Unix(int64(t), 0)
				orderEnd = end.Format("02.01.2006 15:04")
				start, err := time.Parse("02.01.2006 15:04:05", a.OrderStart)
				if err != nil {
					logger.Error("failed to parse order start. error: %w", err)
				}
				if err == nil {
					term = fmt.Sprintf("%.2f", (end.Sub(start).Hours() / 24))
				}
			}
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateEnd"], curNum), orderEnd)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["term"], curNum), term)
			curNum++

			file.SetRowStyle(sheet, curNum, curNum, posStyle)
			file.SetCellStyle(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), fmt.Sprintf("%s%d", columnAxis["number"], curNum), numberStyle1)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), a.Number)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["position"], curNum), a.Position)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["title"], curNum), a.PosTitle)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["ring"], curNum), a.Ring)
			posEnd := ""
			if a.PosEnd != "" {
				t, err := strconv.Atoi(a.PosEnd)
				if err != nil {
					return nil, fmt.Errorf("failed to parse pos date. error: %w", err)
				}
				posEnd = time.Unix(int64(t), 0).Format("02.01.2006 15:04")
			}
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateEnd"], curNum), posEnd)
			curNum++

			file.SetCellStyle(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), fmt.Sprintf("%s%d", columnAxis["number"], curNum), numberStyle)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), a.Number)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["operation"], curNum), a.OperTitle)
			file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateEnd"], curNum), a.OperEnd)
			curNum++
		} else {
			if a.Number == analytics[i-1].Number {
				if a.Position == analytics[i-1].Position && a.PosTitle == analytics[i-1].PosTitle && a.Ring == analytics[i-1].Ring {
					file.SetCellStyle(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), fmt.Sprintf("%s%d", columnAxis["number"], curNum), numberStyle)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), a.Number)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["operation"], curNum), a.OperTitle)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateEnd"], curNum), a.OperEnd)
					curNum++
				} else {
					file.SetRowStyle(sheet, curNum, curNum, posStyle)
					file.SetCellStyle(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), fmt.Sprintf("%s%d", columnAxis["number"], curNum), numberStyle1)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), a.Number)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["position"], curNum), a.Position)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["title"], curNum), a.PosTitle)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["ring"], curNum), a.Ring)
					posEnd := ""
					if a.PosEnd != "" {
						t, err := strconv.Atoi(a.PosEnd)
						if err != nil {
							return nil, fmt.Errorf("failed to parse pos date. error: %w", err)
						}
						posEnd = time.Unix(int64(t), 0).Format("02.01.2006 15:04")
					}
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateEnd"], curNum), posEnd)
					curNum++

					file.SetCellStyle(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), fmt.Sprintf("%s%d", columnAxis["number"], curNum), numberStyle)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), a.Number)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["operation"], curNum), a.OperTitle)
					file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateEnd"], curNum), a.OperEnd)
					curNum++
				}
			} else {
				file.SetRowStyle(sheet, curNum, curNum, orderStyle)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), a.Number)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateStart"], curNum), a.OrderStart)
				orderEnd := ""
				term := ""
				if a.OrderEnd != "" {
					t, err := strconv.Atoi(a.OrderEnd)
					if err != nil {
						return nil, fmt.Errorf("failed to parse order date. error: %w", err)
					}
					end := time.Unix(int64(t), 0)
					orderEnd = end.Format("02.01.2006 15:04")
					start, err := time.Parse("02.01.2006 15:04:05", a.OrderStart)
					if err != nil {
						logger.Error("failed to parse order start. error: %w", err)
					}
					if err == nil {
						term = fmt.Sprintf("%.2f", (end.Sub(start).Hours() / 24))
					}
				}
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateEnd"], curNum), orderEnd)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["term"], curNum), term)
				curNum++

				file.SetRowStyle(sheet, curNum, curNum, posStyle)
				file.SetCellStyle(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), fmt.Sprintf("%s%d", columnAxis["number"], curNum), numberStyle1)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), a.Number)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["position"], curNum), a.Position)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["title"], curNum), a.PosTitle)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["ring"], curNum), a.Ring)
				posEnd := ""
				if a.PosEnd != "" {
					t, err := strconv.Atoi(a.PosEnd)
					if err != nil {
						return nil, fmt.Errorf("failed to parse pos date. error: %w", err)
					}
					posEnd = time.Unix(int64(t), 0).Format("02.01.2006 15:04")
				}
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateEnd"], curNum), posEnd)
				curNum++

				file.SetCellStyle(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), fmt.Sprintf("%s%d", columnAxis["number"], curNum), numberStyle)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["number"], curNum), a.Number)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["operation"], curNum), a.OperTitle)
				file.SetCellValue(sheet, fmt.Sprintf("%s%d", columnAxis["dateEnd"], curNum), a.OperEnd)
				curNum++
			}
		}
	}

	return file, nil
}
