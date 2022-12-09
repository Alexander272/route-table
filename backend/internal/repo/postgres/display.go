package postgres

import "github.com/jmoiron/sqlx"

type DisplayRepo struct {
	db *sqlx.DB
}

func NewDisplayRepo(db *sqlx.DB) *DisplayRepo {
	return &DisplayRepo{db: db}
}
