package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/route-table/internal/config"
	"github.com/Alexander272/route-table/internal/models"
)

type UrgencyService struct {
	urgency *config.UrgencyConfig
}

func NewUrgencyService(urgency *config.UrgencyConfig) *UrgencyService {
	return &UrgencyService{urgency: urgency}
}

func (s *UrgencyService) Get(ctx context.Context) models.Urgency {
	urgency := models.Urgency{
		High:   s.urgency.High.Hours(),
		Middle: s.urgency.Middle.Hours(),
	}

	return urgency
}

// изменение конфига программы (срочность заказов)
func (s *UrgencyService) Change(ctx context.Context, urgency models.Urgency) error {
	high, err := time.ParseDuration(fmt.Sprintf("%.1fh", urgency.High))
	if err != nil {
		return fmt.Errorf("failed to parse high urgency. error: %w", err)
	}
	middle, err := time.ParseDuration(fmt.Sprintf("%.1fh", urgency.Middle))
	if err != nil {
		return fmt.Errorf("failed to parse middle urgency. error: %w", err)
	}

	err = s.urgency.ChangeUrgency(high, middle)
	if err != nil {
		return fmt.Errorf("failed to change urgency. error: %w", err)
	}

	return nil
}
