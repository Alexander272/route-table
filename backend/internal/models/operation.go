package models

import "github.com/google/uuid"

type Operation struct {
	Id         uuid.UUID `db:"id" json:"id"`
	Title      string    `db:"title" json:"title"`
	Done       bool      `db:"done" json:"done"`
	Remainder  int       `db:"remainder" json:"remainder"`
	StepNumber int       `db:"step_number" json:"step_number,omitempty"`
}

type OperationDTO struct {
	Id          uuid.UUID `json:"id"`
	OperationId uuid.UUID `json:"operationId"`
	PositionId  uuid.UUID `json:"positionId"`
	Done        bool      `json:"done"`
	Remainder   int       `json:"remainder"`
	Date        string    `json:"date"`
}

type CompleteOperation struct {
	Id        uuid.UUID `json:"id"`
	Done      bool      `json:"done"`
	Remainder int       `json:"remainder"`
	Reason    string    `json:"reason"`
}
