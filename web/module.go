package web

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	ws.log.Info("bind modules", zap.Int("count", len(ws.modules)))
	if ws.Webroot != "" {
		ws.WebrootFS = static.LocalFile(ws.Webroot, false)
	}
	r.Use(func(c *gin.Context) {
		if !ws.WebrootIndexDisable {
			_, err := ws.WebrootFS.Open(c.Request.URL.Path)
			if err != nil {
				c.FileFromFS("/", ws.WebrootFS)
				return
			}
		}
		c.FileFromFS(c.Request.URL.Path, ws.WebrootFS)
	})
}
