package models

import "github.com/google/uuid"

type OperationDTO struct {
	Id          uuid.UUID `json:"id"`
	OperationId uuid.UUID `json:"operationId"`
	PositionId  uuid.UUID `json:"positionId"`
	Done        bool      `json:"done"`
	Remainder   int       `json:"remainder"`
	Date        string    `json:"date"`
}
