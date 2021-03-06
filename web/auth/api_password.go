package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dev.sum7.eu/genofire/golang-lib/web"
)

// @Summary Change Password
// @Description Change Password of current login user
// @Accept json
// @Produce  json
// @Success 200 {object} boolean "if password was saved (e.g. `true`)"
// @Failure 400 {object} web.HTTPError
// @Failure 401 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /api/v1/my/auth/password [post]
// @Security ApiKeyAuth
// @Param body body string false "new password"
func init() {
	web.ModuleRegister(func(r *gin.Engine, ws *web.Service) {
		r.POST("/api/v1/my/auth/password", MiddlewareLogin(ws), func(c *gin.Context) {
			d, ok := GetCurrentUser(c, ws)
			if !ok {
				return
			}
			var password string
			if err := c.BindJSON(&password); err != nil {
				c.JSON(http.StatusBadRequest, web.HTTPError{
					Message: web.APIErrorInvalidRequestFormat,
					Error:   err.Error(),
				})
				return
			}
			if err := d.SetPassword(password); err != nil {
				c.JSON(http.StatusInternalServerError, web.HTTPError{
					Message: APIErrroCreatePassword,
					Error:   err.Error(),
				})
				return
			}

			if err := ws.DB.Save(&d).Error; err != nil {
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
