package v1

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func (h *Handler) InitOrderRoutes(api *gin.RouterGroup) {
	orders := api.Group("/orders")
	{
		orders.POST("/parse", h.ordersParse)
		orders.GET("/:id", h.getOrder)
		orders.GET("/number/:number", h.findOrders)
	}
}

func (h *Handler) ordersParse(c *gin.Context) {
	fileHeader, err := c.FormFile("order")
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "error while opening file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "error while opening file")
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "error while reading file")
		return
	}

	err = h.services.Order.Parse(c, f)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Order added successfully"})
}

func (h *Handler) getOrder(c *gin.Context) {
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

	order, err := h.services.Order.GetWithPositions(c, Id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: order})
}

func (h *Handler) findOrders(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty number", "empty number param")
		return
	}

	orders, err := h.services.Order.Find(c, number)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: orders, Count: len(orders)})
}
