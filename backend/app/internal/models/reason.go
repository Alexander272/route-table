package models

import "github.com/google/uuid"

type ReasonDTO struct {
	Id          uuid.UUID `json:"id"`
	OperationId uuid.UUID `json:"operationId"`
	Date        string    `json:"date"`
	Value       string    `json:"value"`
}

type Reason struct {
	Id    uuid.UUID `json:"id"`
	Date  string    `json:"date"`
	Value string    `json:"value"`
}

type PosWithReason struct {
	Number   string `db:"number"`
	PosTitle string `db:"pos_title"`
	OpTitle  string `db:"op_title"`
	Date     string `db:"date"`
	Value    string `db:"value"`
}
