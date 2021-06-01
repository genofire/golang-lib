package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dev.sum7.eu/genofire/golang-lib/web"
)

// MiddlewareLogin if user id in session for golang-gin
func MiddlewareLogin(ws *web.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := GetCurrentUserID(c)
		if !ok {
			c.Abort()
		}
	}
}

// MiddlewarePermissionParamUUID if user has access to obj, check access by uuid in golang-gin url param uuid
func MiddlewarePermissionParamUUID(ws *web.Service, obj HasPermission) gin.HandlerFunc {
	return MiddlewarePermissionParam(ws, obj, "uuid")
}

// MiddlewarePermissionParam if user has access to obj, check access in golang-gin url by param
func MiddlewarePermissionParam(ws *web.Service, obj HasPermission, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetCurrentUserID(c)
		if !ok {
			c.Abort()
		}
		objID, err := uuid.Parse(c.Params.ByName(param))
		if err != nil {
			c.JSON(http.StatusUnauthorized, web.HTTPError{
				Message: web.APIErrorInvalidRequestFormat,
				Error:   err.Error(),
			})
			c.Abort()
		}
		_, err = obj.HasPermission(ws.DB, userID, objID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, web.HTTPError{
				Message: http.StatusText(http.StatusUnauthorized),
				Error:   err.Error(),
			})
			c.Abort()
		}
	}
}
