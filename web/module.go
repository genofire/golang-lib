package web

import (
	"github.com/bdlm/log"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// ModuleRegisterFunc format of module which registered to WebService
type ModuleRegisterFunc func(*gin.Engine, *Service)

// ModuleRegister used on start of WebService
func (ws *Service) ModuleRegister(f ModuleRegisterFunc) {
	ws.modules = append(ws.modules, f)
}

// Bind WebService to gin.Engine
func (ws *Service) Bind(r *gin.Engine) {
	for _, f := range ws.modules {
		f(r, ws)
	}

	log.Infof("loaded %d modules", len(ws.modules))
	r.Use(static.Serve("/", static.LocalFile(ws.Webroot, false)))
}
