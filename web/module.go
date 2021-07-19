package web

import (
	"github.com/bdlm/log"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// A ModuleRegisterFunc is a module.
type ModuleRegisterFunc func(*gin.Engine, *Service)

// ModuleRegister adds f to ws's list of modules.
func (ws *Service) ModuleRegister(f ModuleRegisterFunc) {
	ws.modules = append(ws.modules, f)
}

// Bind executes all of ws's modules with r.
func (ws *Service) Bind(r *gin.Engine) {
	for _, f := range ws.modules {
		f(r, ws)
	}

	log.Infof("loaded %d modules", len(ws.modules))
	r.Use(static.Serve("/", static.LocalFile(ws.Webroot, false)))
}
