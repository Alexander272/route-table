package middleware

import (
	"github.com/Alexander272/route-table/internal/config"
	"github.com/Alexander272/route-table/internal/service"
)

type Middleware struct {
	CookieName           string
	services             *service.Services
	auth                 *config.AuthConfig
	UserIdCtx            string
	RoleCtx              string
	EnabledOperationsCtx string
}

func NewMiddleware(services *service.Services, auth *config.AuthConfig) *Middleware {
	return &Middleware{
		services:             services,
		auth:                 auth,
		UserIdCtx:            "userId",
		RoleCtx:              "role",
		EnabledOperationsCtx: "operations",
	}
}
