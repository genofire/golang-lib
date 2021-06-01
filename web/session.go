package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// LoadSession module to start Session Handling in WebService
func (config *Service) LoadSession(r *gin.Engine) {
	store := cookie.NewStore([]byte(config.Session.Secret))
	r.Use(sessions.Sessions(config.Session.Name, store))
}
