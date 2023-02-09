package v1

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (h *Handler) InitPositionsRoutes(api *gin.RouterGroup) {
	positions := api.Group("/positions", h.middleware.UserIdentity)
	{
		positions.GET("/:id", h.getPosition)
		positions.PUT("/:id/count", h.middleware.AccessForMaster, h.updatePositionCount)
	}
}

func (h *Handler) getPosition(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	role, _ := c.Get(h.middleware.RoleCtx)
	enabledOperations, _ := c.Get(h.middleware.EnabledOperationsCtx)

	uuId, err := uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "empty id param")
		return
	}

	var position interface{}

	if role == "master" || role == "manager" {
		position, err = h.services.Position.GetWithReasons(c, uuId)
		if err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to get operation")
			return
		}
	} else {
		position, err = h.services.Position.Get(c, uuId, enabledOperations.(pq.StringArray))
		if err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to get operation")
			return
		}
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: position})
}

func (h *Handler) updatePositionCount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.UpdateCount
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	var err error
	dto.Id, err = uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "empty id param")
		return
	}

	if err := h.services.Position.UpdateCount(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Position count update successfully"})
}
