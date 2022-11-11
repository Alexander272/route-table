package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/google/uuid"
)

type ReasonService struct {
	repo repository.Reason
}

func NewReasonService(repo repository.Reason) *ReasonService {
	return &ReasonService{repo: repo}
}

func (s *ReasonService) Create(ctx context.Context, reason models.ReasonDTO) (id uuid.UUID, err error) {
	id, err = s.repo.Create(ctx, reason)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create reason. error: %w", err)
	}
	return id, nil
}
