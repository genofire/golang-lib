package auth

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"dev.sum7.eu/genofire/golang-lib/web"
)

// IsLoginWithUserID get UserID of session in golang-gin
func IsLoginWithUserID(c *gin.Context) (uuid.UUID, bool) {
	session := sessions.Default(c)

	v := session.Get("user_id")
	if v == nil {
		return uuid.Nil, false
	}

	id := uuid.MustParse(v.(string))
	return id, true
}

// GetCurrentUserID get UserID of session in golang-gin
func GetCurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	id, ok := IsLoginWithUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.HTTPError{
			Message: ErrAPINoSession.Error(),
		})
	}
	return id, ok
}

// GetCurrentUser get User of session from database in golang-gin
func GetCurrentUser(c *gin.Context, ws *web.Service) (*User, bool) {
	id, ok := GetCurrentUserID(c)
	if !ok {
		return nil, false
	}
	d := &User{ID: id}
	if err := ws.DB.First(d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, web.HTTPError{
				Message: ErrAPIUserNotFound.Error(),
				Error:   err.Error(),
			})
			return nil, false
		}
		c.JSON(http.StatusInternalServerError, web.HTTPError{
			Message: web.ErrAPIInternalDatabase.Error(),
			Error:   err.Error(),
		})
		return nil, false
	}
	return d, true
}
