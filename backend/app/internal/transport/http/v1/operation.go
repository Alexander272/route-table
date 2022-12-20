package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) InitOperationRoutes(api *gin.RouterGroup) {
	operations := api.Group("/operations", h.middleware.UserIdentity)
	{
		operations.PUT("/:id", h.completeOperation)
		operations.POST("/rollback/:id", h.rollbackOperation)
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

	var dto models.CompletePosition
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}
	dto.Id = uuId

	if err := h.services.Position.Update(c, dto); err != nil {
		if strings.Contains(err.Error(), "connected operation not completed") {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "connected operation not completed")
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "Operation updated successfully"})
}

func (h *Handler) rollbackOperation(c *gin.Context) {
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
	var dto models.RollbackPosition
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}
	dto.OperationId = uuId

	if err := h.services.Position.Rollback(c, dto); err != nil {
		if errors.Is(err, models.ErrOperationNotFound) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Operation rollback successfully"})
}
