package web

import (
	"github.com/bdlm/log"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	modules []ModuleRegisterFunc
)

// ModuleRegisterFunc format of module which registered to WebService
type ModuleRegisterFunc func(*gin.Engine, *Service)

// ModuleRegister used on start of WebService
func ModuleRegister(f ModuleRegisterFunc) {
	modules = append(modules, f)
}

// Bind WebService to gin.Engine
func (ws *Service) Bind(r *gin.Engine) {
	for _, f := range modules {
		f(r, ws)
	}

	log.Infof("loaded %d modules", len(modules))
	r.Use(static.Serve("/", static.LocalFile(ws.Webroot, false)))
}
