package v1

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) InitPositionsRoutes(api *gin.RouterGroup) {
	positions := api.Group("/positions")
	{
		positions.GET("/:id", h.getPosition)
	}
}

func (h *Handler) getPosition(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	reason := c.Query("reason")

	uuId, err := uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "empty id param")
		return
	}

	var position interface{}

	if reason != "" {
		position, err = h.services.Position.GetWithReasons(c, uuId)
		if err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to get operation")
			return
		}
	} else {
		position, err = h.services.Position.Get(c, uuId)
		if err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to get operation")
			return
		}
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: position})
}
