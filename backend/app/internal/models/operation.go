package models

import "github.com/google/uuid"

type Operation struct {
	Id         uuid.UUID `db:"id" json:"id"`
	Title      string    `db:"title" json:"title"`
	Done       bool      `db:"done" json:"done"`
	Remainder  int       `db:"remainder" json:"remainder"`
	StepNumber int       `db:"step_number" json:"step_number,omitempty"`
	IsFinish   bool      `db:"is_finish" json:"isFinish"`
	Date       *string   `db:"date" json:"-"`
}

type OperationWithReason struct {
	Id         uuid.UUID `db:"id" json:"id"`
	Title      string    `db:"title" json:"title"`
	Done       bool      `db:"done" json:"done"`
	Remainder  int       `db:"remainder" json:"remainder"`
	StepNumber int       `db:"step_number" json:"step_number,omitempty"`
	IsFinish   bool      `db:"is_finish" json:"isFinish"`
	OpDate     *string   `db:"op_date" json:"date"`
	ReasonId   uuid.UUID `db:"reason_id" json:"-"`
	Value      *string   `db:"value" json:"-"`
	Date       *string   `db:"date" json:"-"`
	Reason     []Reason  `json:"reasons"`
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
	Count     int       `json:"count"`
	Reason    string    `json:"reason"`
}

type CompletedOperation struct {
	Id          uuid.UUID `json:"id" db:"id"`
	OperationId uuid.UUID `json:"operationId" db:"operation_id"`
	GroupId     uuid.UUID `json:"groupId" db:"group_id"`
	Remainder   int       `json:"remainder" db:"remainder"`
	Count       int       `json:"count" db:"count"`
	Date        int64     `db:"date"`
}
