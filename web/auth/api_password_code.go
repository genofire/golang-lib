package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"dev.sum7.eu/genofire/golang-lib/web"
)

// PasswordWithForgetCode - JSON Request to set password without login
type PasswordWithForgetCode struct {
	ForgetCode uuid.UUID `json:"forget_code"`
	Password   string    `json:"password"`
}

// @Summary Change Password with ForgetCode
// @Description Change Password of any user by generated forget code
// @Tags auth
// @Accept json
// @Produce  json
// @Success 200 {object} string "username of changed password (e.g. `"admin"`)"
// @Failure 400 {object} web.HTTPError
// @Failure 401 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /api/v1/auth/password/code [post]
// @Param body body PasswordWithForgetCode false "new password and forget code"
func apiPasswordCode(r *gin.Engine, ws *web.Service) {
	r.POST("/api/v1/auth/password/code", func(c *gin.Context) {
		var req PasswordWithForgetCode
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, web.HTTPError{
				Message: web.APIErrorInvalidRequestFormat,
				Error:   err.Error(),
			})
			return
		}
		d := User{}
		if err := ws.DB.Where("forget_code", req.ForgetCode).First(&d).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, web.HTTPError{
					Message: APIErrorUserNotFound,
					Error:   err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, web.HTTPError{
				Message: APIErrroCreatePassword,
				Error:   err.Error(),
			})
			return
		}
		if err := d.SetPassword(req.Password); err != nil {
			c.JSON(http.StatusInternalServerError, web.HTTPError{
				Message: APIErrroCreatePassword,
				Error:   err.Error(),
			})
			return
		}
		d.ForgetCode = nil

		if err := ws.DB.Save(&d).Error; err != nil {
			c.JSON(http.StatusInternalServerError, web.HTTPError{
				Message: web.APIErrorInternalDatabase,
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, d.Username)
	})
}
