package models

import "github.com/google/uuid"

type OrderDTO struct {
	Id       uuid.UUID `json:"id"`
	Number   string    `json:"number"`
	Done     bool      `json:"done"`
	Deadline string    `json:"deadline"`
}
