package auth

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"codeberg.org/genofire/golang-lib/web"
)

type login struct {
	Username string `json:"username" example:"kukoon"`
	Password string `json:"password" example:"super secret password"`
}

// @Summary Login
// @Description Login by username and password, you will get a cookie of current session
// @Tags auth
// @Accept json
// @Produce  json
// @Success 200 {object} User
// @Failure 400 {object} web.HTTPError
// @Failure 401 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /api/v1/auth/login [post]
// @Param body body login false "login"
func apiLogin(r *gin.Engine, ws *web.Service) {
	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		var data login
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, web.HTTPError{
				Message: web.ErrAPIInvalidRequestFormat.Error(),
				Error:   err.Error(),
			})
			return
		}

		d := &User{}
		if err := ws.DB.Where(map[string]interface{}{"username": data.Username}).First(d).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, web.HTTPError{
					Message: ErrAPIUserNotFound.Error(),
					Error:   err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, web.HTTPError{
				Message: web.ErrAPIInternalDatabase.Error(),
				Error:   err.Error(),
			})
			return
		}
		if !d.ValidatePassword(data.Password) {
			c.JSON(http.StatusUnauthorized, web.HTTPError{
				Message: ErrAPIIncorrectPassword.Error(),
			})
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", d.ID.String())
		if err := session.Save(); err != nil {
			c.JSON(http.StatusBadRequest, web.HTTPError{
				Message: ErrAPICreateSession.Error(),
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, d)
	})
}
