package models

import "github.com/google/uuid"

type ReasonDTO struct {
	Id          uuid.UUID `json:"id"`
	OperationId uuid.UUID `json:"operationId"`
	Date        string    `json:"date"`
	Value       string    `json:"value"`
}
