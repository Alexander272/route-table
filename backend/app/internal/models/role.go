package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Role struct {
	Id         uuid.UUID      `db:"id" json:"id"`
	Title      string         `db:"title" json:"title"`
	Role       string         `db:"role" json:"role"`
	Operations pq.StringArray `db:"operations" json:"operations"`
}

type RoleDTO struct {
	Id         uuid.UUID   `json:"id"`
	Title      string      `json:"title"`
	Role       string      `json:"role"`
	Operations []uuid.UUID `json:"operations"`
}
