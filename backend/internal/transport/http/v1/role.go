package v1

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) InitRoleRoutes(api *gin.RouterGroup) {
	roles := api.Group("/roles")
	{
		roles.GET("/", h.getRoles)
		roles.POST("/", h.createRole)
		roles.PUT("/:id", h.updateRole)
		roles.DELETE("/:id", h.deleteRole)
	}
}

func (h *Handler) getRoles(c *gin.Context) {
	roles, err := h.services.Role.Get(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to get roles")
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: roles, Count: len(roles)})
}

func (h *Handler) createRole(c *gin.Context) {
	var dto models.RoleDTO
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.services.Role.Create(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, response.IdResponse{Id: id.String(), Message: "Role added successfully"})
}

func (h *Handler) updateRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.RoleDTO
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

	if err := h.services.Role.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "Role updated successfully"})
}

func (h *Handler) deleteRole(c *gin.Context) {
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
	dto := models.RoleDTO{Id: Id}

	if err := h.services.Role.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "Role removed successfully"})
}
