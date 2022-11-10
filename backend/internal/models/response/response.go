package response

import (
	"strings"

	"github.com/Alexander272/route-table/pkg/logger"
	"github.com/gin-gonic/gin"
)

type DataResponse struct {
	Data  interface{} `json:"data"`
	Count int         `json:"count,omitempty"`
}

type IdResponse struct {
	Id      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func NewErrorResponse(c *gin.Context, statusCode int, err, message string) {
	code := "U001"
	if strings.Contains(err, "execute query") {
		code = "MD001"
	} else if strings.Contains(err, "EOF") {
		code = "E001"
	}

	logger.Errorf("Url: %s | ClientIp: %s | ErrorResponse: %s", c.Request.URL, c.ClientIP(), err)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{Message: message, Code: code})
}
