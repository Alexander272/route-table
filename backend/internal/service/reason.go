package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type ReasonService struct {
	repo repository.Reason
}

func NewReasonService(repo repository.Reason) *ReasonService {
	return &ReasonService{repo: repo}
}

func (s *ReasonService) Get(ctx context.Context) (reasons []models.PosWithReason, err error) {
	reasons, err = s.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get reasons. error: %w", err)
	}

	return reasons, nil
}

func (s *ReasonService) GetFile(ctx context.Context) (file *excelize.File, err error) {
	file = excelize.NewFile()

	err = file.SetCellValue("Sheet1", "A1", "Наименование")
	if err != nil {
		return nil, err
	}
	err = file.SetCellValue("Sheet1", "B1", "Операция")
	if err != nil {
		return nil, err
	}
	err = file.SetCellValue("Sheet1", "C1", "Причина")
	if err != nil {
		return nil, err
	}
	err = file.SetCellValue("Sheet1", "D1", "Дата")
	if err != nil {
		return nil, err
	}

	index := 2

	reasons, err := s.Get(ctx)
	if err != nil {
		return nil, err
	}

	for i, r := range reasons {
		err = file.SetCellValue("Sheet1", fmt.Sprintf("A%d", index+i), r.PosTitle)
		if err != nil {
			return nil, err
		}
		err = file.SetCellValue("Sheet1", fmt.Sprintf("B%d", index+i), r.OpTitle)
		if err != nil {
			return nil, err
		}
		err = file.SetCellValue("Sheet1", fmt.Sprintf("C%d", index+i), r.Value)
		if err != nil {
			return nil, err
		}
		err = file.SetCellValue("Sheet1", fmt.Sprintf("D%d", index+i), r.Date)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func (s *ReasonService) Create(ctx context.Context, reason models.ReasonDTO) (id uuid.UUID, err error) {
	id, err = s.repo.Create(ctx, reason)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create reason. error: %w", err)
	}
	return id, nil
}
