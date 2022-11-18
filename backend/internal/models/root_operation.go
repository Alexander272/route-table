package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RootOperation struct {
	Id         uuid.UUID      `db:"id"`
	Title      string         `db:"title"`
	Gasket     string         `db:"gasket"`
	StepNumber int            `db:"step_number"`
	Connected  pq.StringArray `db:"connected"`
}

type RootOperationDTO struct {
	Id         uuid.UUID `json:"id"`
	Title      string    `json:"title" binding:"required"`
	Gasket     string    `json:"gasket" binding:"required"`
	StepNumber int       `json:"stepNumber"`
}
