package v1

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) InitOperationRoutes(api *gin.RouterGroup) {
	operations := api.Group("/operations")
	{
		operations.PATCH("/:id", h.completeOperation)
	}
}

func (h *Handler) completeOperation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "empty id param")
		return
	}

	var dto models.CompleteOperation
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}
	dto.Id = uuId

	if err := h.services.Operation.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to get operation")
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "Operation updated successfully"})
}
