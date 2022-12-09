package v1

import (
	"net/http"
	"time"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/Alexander272/route-table/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) InitUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("/", h.getUsers)
		users.POST("/", h.createUser)
		users.PUT("/:id", h.updateUser)
		users.DELETE("/:id", h.deleteUser)
	}
}

func (h *Handler) getUsers(c *gin.Context) {
	users, err := h.services.User.GetAll(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to get users")
		return
	}

	logger.Debug(time.Now().Format("02.01.2006 15:04"))
	logger.Debug(time.Now().Location())
	logger.Debug(time.Now().Local().Format("02.01.2006 15:04"))
	logger.Debug(time.Now().UTC().Format("02.01.2006 15:04"))

	c.JSON(http.StatusOK, response.DataResponse{Data: users, Count: len(users)})
}

func (h *Handler) createUser(c *gin.Context) {
	var dto models.UserDTO
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	id, err := h.services.User.Create(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, response.IdResponse{Id: id.String(), Message: "User added successfully"})
}

func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.UserDTO
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

	if err := h.services.User.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "User updated successfully"})
}

func (h *Handler) deleteUser(c *gin.Context) {
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
	dto := models.UserDTO{Id: Id}

	if err := h.services.User.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "User removed successfully"})
}
