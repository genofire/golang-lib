package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"codeberg.org/genofire/golang-lib/web"
)

// @Summary Login status
// @Description show user_id and username if logged in
// @Tags auth
// @Accept json
// @Produce  json
// @Success 200 {object} User
// @Failure 401 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /api/v1/my/auth/status [get]
// @Security ApiKeyAuth
func apiMyStatus(r *gin.Engine, ws *web.Service) {
	r.GET("/api/v1/my/auth/status", MiddlewareLogin(ws), func(c *gin.Context) {
		d, ok := GetCurrentUser(c, ws)
		if ok {
			c.JSON(http.StatusOK, d)
		}
	})
}
