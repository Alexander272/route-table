package middleware

import (
	"net/http"

	"github.com/Alexander272/route-table/internal/models/response"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) UserIdentity(c *gin.Context) {
	token, err := c.Cookie(m.CookieName)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	user, err := m.services.Session.TokenParse(token)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	isRefresh, err := m.services.Session.CheckSession(c, token)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	if isRefresh {
		_, token, err := m.services.Session.Refresh(c, user)
		if err != nil {
			response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "failed to refresh session")
			return
		}

		c.SetCookie(m.CookieName, token, int(m.auth.RefreshTokenTTL.Seconds()), "/", m.auth.Domain, m.auth.Secure, true)
	}

	c.Set(m.UserIdCtx, user.Id)
	c.Set(m.RoleCtx, user.Role.Role)
	c.Set(m.EnabledOperationsCtx, user.Role.Operations)
}

func (m *Middleware) AccessForMaster(c *gin.Context) {
	role, _ := c.Get(m.RoleCtx)
	if role != "master" {
		response.NewErrorResponse(c, http.StatusForbidden, "role "+role.(string)+" access is denied", "access is denied")
		return
	}
}

func (m *Middleware) AccessForDisplay(c *gin.Context) {
	role, _ := c.Get(m.RoleCtx)
	if role != "display" && role != "master" {
		response.NewErrorResponse(c, http.StatusForbidden, "role "+role.(string)+" access is denied", "access is denied")
		return
	}
}
