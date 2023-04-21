package v1

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) InitRootOperationRoutes(api *gin.RouterGroup) {
	root := api.Group("/root-operations")
	{
		root.GET("/", h.getRootOperation)
		root.POST("/", h.createRootOperation)
		root.PUT("/:id", h.updateRootOperation)
		root.DELETE("/:id", h.deleteRootOperation)
	}
}

func (h *Handler) getRootOperation(c *gin.Context) {
	operation, err := h.services.RootOperation.Get(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to get operation")
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: operation, Count: len(operation)})
}

func (h *Handler) createRootOperation(c *gin.Context) {
	var dto models.RootOperationDTO
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.services.RootOperation.Create(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, response.IdResponse{Id: id.String(), Message: "Root operation added successfully"})
}

func (h *Handler) updateRootOperation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.RootOperationDTO
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

	if err := h.services.RootOperation.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "Root operation updated successfully"})
}

func (h *Handler) deleteRootOperation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	Id, err := uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "empty id param")
		return
	}

	dto := models.RootOperationDTO{Id: Id}

	if err := h.services.RootOperation.Delete(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "Root operation removed successfully"})
}
