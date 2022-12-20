package v1

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitUrgencyRoutes(api *gin.RouterGroup) {
	urgency := api.Group("/urgency", h.middleware.UserIdentity, h.middleware.AccessForMaster)
	{
		urgency.GET("/", h.getUrgency)
		urgency.PUT("/", h.changeUrgency)
	}
}

func (h *Handler) getUrgency(c *gin.Context) {
	urgency := h.services.Urgency.Get(c)

	c.JSON(http.StatusOK, response.DataResponse{Data: urgency})
}

func (h *Handler) changeUrgency(c *gin.Context) {
	var dto models.Urgency
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	if err := h.services.Urgency.Change(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Urgency updated successfully"})
}
