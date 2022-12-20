package v1

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitReasonsRoutes(api *gin.RouterGroup) {
	reasons := api.Group("/reasons", h.middleware.UserIdentity)
	{
		reasons.GET("/", h.getReasons)
		reasons.GET("/file", h.getReasonsFile)
	}
}

func (h *Handler) getReasons(c *gin.Context) {
	reasons, err := h.services.Reason.Get(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: reasons})
}

func (h *Handler) getReasonsFile(c *gin.Context) {
	file, err := h.services.Reason.GetFile(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	fileName := fmt.Sprintf("%s.xlsx", time.Now().Format("02.01.2006 15:04"))
	file.SaveAs(fileName)
	defer os.Remove(fileName)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	// c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.File(fileName)
}
