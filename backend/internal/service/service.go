package service

import (
	"time"

	repository "github.com/Alexander272/route-table/internal/repo"
	"github.com/Alexander272/route-table/pkg/auth"
)

type Services struct {
}

type Deps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps Deps) *Services {
	return &Services{}
}
