package models

import "github.com/google/uuid"

type PositionForOrder struct {
	Id            uuid.UUID `db:"id" json:"id"`
	Position      int       `db:"position" json:"position"`
	Count         int       `db:"count" json:"count"`
	Title         string    `db:"title" json:"title"`
	Ring          string    `db:"ring" json:"ring"`
	Deadline      string    `db:"deadline" json:"deadline"`
	Connected     uuid.UUID `db:"connected" json:"connected"`
	Done          bool      `db:"done" json:"done"`
	LastOperation *string   `db:"last_operation" json:"lastOperation"`
	CurOperation  *string   `db:"cur_operation" json:"curOperation"`
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
	Operation []Operation `json:"operations"`
}

type PositionWithReason struct {
	Id        uuid.UUID             `db:"id" json:"id"`
	Order     string                `db:"number" json:"order"`
	Position  int                   `db:"position" json:"position"`
	Count     int                   `db:"count" json:"count"`
	Title     string                `db:"title" json:"title"`
	Ring      string                `db:"ring" json:"ring"`
	Deadline  string                `db:"deadline" json:"deadline"`
	Connected uuid.UUID             `db:"connected" json:"connected,omitempty"`
	Done      bool                  `db:"done" json:"done"`
	Operation []OperationWithReason `json:"operations"`
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
	Completed string    `json:"completed"`
}

type CompletePosition struct {
	Id        uuid.UUID         `json:"positionId"`
	Count     int               `json:"count"`
	IsFinish  bool              `json:"isFinish"`
	Connected uuid.UUID         `json:"connected"`
	Operation CompleteOperation `json:"operation"`
}

type RollbackPosition struct {
	Id                uuid.UUID `json:"id"`
	Connected         uuid.UUID `json:"connected"`
	Reasons           []string  `json:"reasons"`
	IsFinishOperation bool      `json:"isFinishOperation"`
	OperationId       uuid.UUID
}

type UpdateCount struct {
	Id    uuid.UUID `json:"id"`
	Count int       `json:"count"`
	Done  bool      `json:"done"`
}
