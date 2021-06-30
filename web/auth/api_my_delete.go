package auth

import (
	"net/http"

	"dev.sum7.eu/genofire/golang-lib/web"
	"github.com/gin-gonic/gin"
)

// @Summary Delete own User
// @Description delete current loggedin user
// @Tags auth
// @Accept json
// @Produce  json
// @Success 200 {object} bool "true if deleted"
// @Failure 401 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /api/v1/my/profil [delete]
// @Security ApiKeyAuth
func init() {
	web.ModuleRegister(func(r *gin.Engine, ws *web.Service) {
		r.DELETE("/api/v1/my/profil", func(c *gin.Context) {
			id, ok := GetCurrentUserID(c)
			if !ok {
				return
			}
			if err := ws.DB.Delete(&User{ID: id}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, web.HTTPError{
					Message: web.APIErrorInternalDatabase,
					Error:   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, true)
		})
	})
}
