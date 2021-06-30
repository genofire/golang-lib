package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dev.sum7.eu/genofire/golang-lib/web"
)

// @Summary Login status
// @Description show user_id and username if logged in
// @Tags auth
// @Accept json
// @Produce  json
// @Success 200 {object} User
// @Failure 401 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /api/v1/auth/status [get]
// @Security ApiKeyAuth
func init() {
	web.ModuleRegister(func(r *gin.Engine, ws *web.Service) {
		r.GET("/api/v1/auth/status", MiddlewareLogin(ws), func(c *gin.Context) {
			d, ok := GetCurrentUser(c, ws)
			if ok {
				c.JSON(http.StatusOK, d)
			}
		})
	})
}
