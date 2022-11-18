package v1

import (
	"github.com/Alexander272/route-table/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
	// middleware *middleware.Middleware
}

// const CookieName = "session"

func NewHandler(services *service.Services) *Handler {
	// middleware.CookieName = CookieName

	return &Handler{
		services: services,
		// middleware: middleware,
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
		h.InitUsersRoutes(v1)
		v1.GET("/", h.notImplemented)
	}
}

func (h *Handler) notImplemented(c *gin.Context) {}
