package v1

import (
	"github.com/Alexander272/route-table/internal/config"
	"github.com/Alexander272/route-table/internal/service"
	"github.com/Alexander272/route-table/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

const CookieName = "session"

type Handler struct {
	services   *service.Services
	auth       *config.AuthConfig
	cookieName string
	middleware *middleware.Middleware
}

func NewHandler(services *service.Services, auth *config.AuthConfig, middleware *middleware.Middleware) *Handler {
	middleware.CookieName = CookieName

	return &Handler{
		services:   services,
		auth:       auth,
		cookieName: CookieName,
		middleware: middleware,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {

	v1 := api.Group("/v1")
	{
		h.InitRootOperationRoutes(v1)
		h.InitOrderRoutes(v1)
		h.InitPositionsRoutes(v1)
		h.InitOperationRoutes(v1)
		h.InitRoleRoutes(v1)
		h.InitReasonsRoutes(v1)
		h.InitUsersRoutes(v1)
		h.InitAuthRoutes(v1)
		h.InitUrgencyRoutes(v1)
		v1.GET("/", h.notImplemented)
	}
}

func (h *Handler) notImplemented(c *gin.Context) {}
