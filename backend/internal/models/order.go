package models

import "github.com/google/uuid"

type FindedOrder struct {
	Id     uuid.UUID `db:"id" json:"id"`
	Number string    `db:"number" json:"number"`
	Done   bool      `db:"done" json:"done"`
}

type Order struct {
	Id       uuid.UUID `db:"id" json:"id"`
	Number   string    `db:"number" json:"number"`
	Done     bool      `db:"done" json:"done"`
	Date     string    `db:"date" json:"date"`
	Deadline string    `db:"deadline" json:"deadline"`
}

type OrderWithPositions struct {
	Id        uuid.UUID          `db:"id" json:"id"`
	Number    string             `db:"number" json:"number"`
	Done      bool               `db:"done" json:"done"`
	Positions []PositionForOrder `json:"positions"`
}

type OrderDTO struct {
	Id       uuid.UUID `json:"id"`
	Number   string    `json:"number"`
	Done     bool      `json:"done"`
	Date     string    `json:"date"`
	Deadline string    `json:"deadline"`
}
