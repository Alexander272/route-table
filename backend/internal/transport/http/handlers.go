package http

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/config"
	"github.com/Alexander272/route-table/internal/service"
	httpV1 "github.com/Alexander272/route-table/internal/transport/http/v1"
	"github.com/Alexander272/route-table/pkg/limiter"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init(conf *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		limiter.Limit(conf.Limiter.RPS, conf.Limiter.Burst, conf.Limiter.TTL),
	)

	// docs.SwaggerInfo_swagger.Host = fmt.Sprintf("%s:%s", conf.Http.Host, conf.Http.Port)
	// if conf.Environment != "dev" {
	// 	docs.SwaggerInfo_swagger.Host = conf.Http.Host
	// }

	// if conf.Environment != "prod" {
	// 	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }

	// Init router
	router.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := httpV1.NewHandler(h.services)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
