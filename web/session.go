package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// LoadSession starts session handling for s.
func (s *Service) LoadSession(r *gin.Engine) {
	store := cookie.NewStore([]byte(s.Session.Secret))
	r.Use(sessions.Sessions(s.Session.Name, store))
}
