package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
}

func NewRepo(db *sqlx.DB, redis redis.Cmdable) *Repositories {
	return &Repositories{}
}
