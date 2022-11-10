package service

import repository "github.com/Alexander272/route-table/internal/repo"

type OperationService struct {
	repo repository.Operation
}

func NewOperationService(repo repository.Operation) *OperationService {
	return &OperationService{repo: repo}
}
