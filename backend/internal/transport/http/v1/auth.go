package v1

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-out", h.signOut)
		auth.POST("/refresh", h.refresh)
	}
}

func (h *Handler) signIn(c *gin.Context) {
	var dto models.SignIn
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	user, token, err := h.services.SignIn(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.SetCookie(h.cookieName, token, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.Host, h.auth.Secure, true)
	c.JSON(http.StatusOK, response.DataResponse{Data: user})
}

func (h *Handler) signOut(c *gin.Context) {
	token, err := c.Cookie(h.cookieName)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	user, err := h.services.Session.TokenParse(token)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	if err := h.services.Session.SingOut(c, user.Id.String()); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.SetCookie(h.cookieName, "", 0, "/", c.Request.Host, h.auth.Secure, true)
	c.JSON(http.StatusNoContent, response.IdResponse{Message: "Sign-out completed successfully"})
}

func (h *Handler) refresh(c *gin.Context) {
	token, err := c.Cookie(h.cookieName)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	user, err := h.services.Session.TokenParse(token)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	_, err = h.services.CheckSession(c, user, token)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	u, newToken, err := h.services.Refresh(c, user)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.SetCookie(h.cookieName, newToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.Host, h.auth.Secure, true)
	c.JSON(http.StatusOK, response.DataResponse{Data: u})
}
