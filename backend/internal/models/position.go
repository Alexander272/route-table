package models

import "github.com/google/uuid"

type PositionForOrder struct {
	Id        uuid.UUID `db:"id"`
	Position  int       `db:"position"`
	Count     int       `db:"count"`
	Title     string    `db:"title"`
	Ring      string    `db:"ring"`
	Deadline  string    `db:"deadline"`
	Connected uuid.UUID `db:"connected"`
	Done      bool      `db:"done"`
}

type Position struct {
	Id        uuid.UUID   `db:"id" json:"id"`
	Order     string      `db:"number" json:"order"`
	Position  int         `db:"position" json:"position"`
	Count     int         `db:"count" json:"count"`
	Title     string      `db:"title" json:"title"`
	Ring      string      `db:"ring" json:"ring"`
	Deadline  string      `db:"deadline" json:"deadline"`
	Connected uuid.UUID   `db:"connected" json:"connected,omitempty"`
	Done      bool        `db:"done" json:"done"`
	Operation []Operation `json:"operation"`
}

type PositionDTO struct {
	Id        uuid.UUID `json:"id"`
	OrderId   uuid.UUID `json:"orderId"`
	Position  int       `json:"position"`
	Count     int       `json:"count"`
	Title     string    `json:"title"`
	Ring      string    `json:"ring"`
	Deadline  string    `json:"deadline"`
	Connected uuid.UUID `json:"connected"`
	Done      bool      `json:"done"`
}
