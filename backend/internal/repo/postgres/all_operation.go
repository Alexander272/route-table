package postgres

import "github.com/jmoiron/sqlx"

type AllOperationRepo struct {
	db *sqlx.DB
}

func NewAllOperationRepo(db *sqlx.DB) *AllOperationRepo {
	return &AllOperationRepo{db: db}
}
