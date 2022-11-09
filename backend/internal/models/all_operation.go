package models

import "github.com/google/uuid"

type AllOperation struct {
	Id         uuid.UUID `db:"id"`
	Title      string    `db:"title"`
	Gasket     string    `db:"gasket"`
	StepNumber int       `db:"step_number"`
}

type AllOperationDTO struct {
	Title      string `json:"title" binding:"required"`
	Gasket     string `json:"gasket" binding:"required,min=3"`
	StepNumber int    `json:"stepNumber"`
}
