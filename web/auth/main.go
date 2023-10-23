package auth

import (
	"codeberg.org/genofire/golang-lib/web"
	"github.com/gin-gonic/gin"
)

// Register to WebService
func Register(r *gin.Engine, ws *web.Service) {
	apiLogin(r, ws)
	apiMyDelete(r, ws)
	apiMyPassword(r, ws)
	apiMyStatus(r, ws)
	apiPasswordCode(r, ws)
}
