package v1

import (
	"net/http"
	"os"
	"strings"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

func (h *Handler) InitOrderRoutes(api *gin.RouterGroup) {
	orders := api.Group("/orders", h.middleware.UserIdentity)
	{
		orders.POST("/parse", h.middleware.AccessForMaster, h.ordersParse)
		orders.GET("/", h.middleware.AccessForDisplay, h.getOrders)
		orders.GET("/group", h.middleware.AccessForDisplay, h.getGroup)
		orders.GET("/:id", h.getOrder)
		orders.GET("/number/:number", h.findOrders)
		orders.PUT("/:id", h.updateOrder)
		orders.DELETE("/:id", h.deleteOrder)

		orders.GET("/analytics", h.getAnalytics)
	}
}

// Загрузка заказа
func (h *Handler) ordersParse(c *gin.Context) {
	fileHeader, err := c.FormFile("order")
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "error while opening file")
		return
	}

	if fileHeader.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
		!strings.Contains(fileHeader.Filename, "xls") {
		response.NewErrorResponse(c, http.StatusInternalServerError, "invalid type file", "invalid file type")
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

// Получение заказов
func (h *Handler) getOrders(c *gin.Context) {
	orders, err := h.services.Order.GetAll(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: orders, Count: len(orders)})
}

// Получение заказов в виде группы по срочности
func (h *Handler) getGroup(c *gin.Context) {
	orders, err := h.services.Order.GetGrouped(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: orders})
}

// Получение заказа по id
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

	role, _ := c.Get(h.middleware.RoleCtx)
	enabledOperations, _ := c.Get(h.middleware.EnabledOperationsCtx)

	order, err := h.services.Order.GetWithPositions(c, Id, role.(string), enabledOperations.(pq.StringArray))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: order})
}

// получение заказов по номеру (лимит 5 штук)
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

// обновление заказа (пока только дата отгрузки)
func (h *Handler) updateOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	var dto models.OrderDTO
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

	if err := h.services.Order.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "Order updated successfully"})
}

func (h *Handler) deleteOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	orderId, err := uuid.Parse(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "empty id param")
		return
	}

	if err := h.services.Order.Delete(c, models.OrderDTO{Id: orderId}); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Id: id, Message: "Order deleted successfully"})
}

func (h *Handler) getAnalytics(c *gin.Context) {
	file, err := h.services.Order.GetForAnalytics(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	fileName := "analytics.xlsx"
	file.SaveAs(fileName)
	defer os.Remove(fileName)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.File(fileName)
}
