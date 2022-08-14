/*
Package web implements common functionality for web APIs using Gin and Gorm.

Modules

Modules provide functionality for a web server. A module is a function executed
before starting a server, accessing the Service and the Gin Engine. Each Service
maintains a list of modules. When it runs, it executes all of its modules.
*/
package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	// acme
	"github.com/gin-gonic/autotls"
	"golang.org/x/crypto/acme/autocert"

	// internal
	"dev.sum7.eu/genofire/golang-lib/mailer"
	"gorm.io/gorm"
)

// A Service stores configuration of a server.
type Service struct {
	// config
	Listen              string          `config:"listen" toml:"listen"`
	AccessLog           bool            `config:"access_log" toml:"access_log"`
	WebrootIndexDisable bool            `config:"webroot_index_disable" toml:"webroot_index_disable"`
	Webroot             string          `config:"webroot" toml:"webroot"`
	WebrootFS           http.FileSystem `config:"-" toml:"-"`
	ACME                struct {
		Enable  bool     `config:"enable" toml:"enable"`
		Domains []string `config:"domains" toml:"domains"`
		Cache   string   `config: "cache" toml:"cache"`
	} `config:"acme" toml:"acme"`
	Session struct {
		Name   string `config:"name" toml:"name"`
		Secret string `config: "secret" toml:"secret"`
	} `config:"session" toml:"session"`
	// internal
	DB     *gorm.DB        `config:"-" toml:"-"`
	Mailer *mailer.Service `config:"-" toml:"-"`

	log     *zap.Logger
	modules []ModuleRegisterFunc
}

// SetLog - set new logger
func (s *Service) SetLog(l *zap.Logger) {
	s.log = l
}

// Log - get current logger
func (s *Service) Log() *zap.Logger {
	return s.log
}

// Run creates, configures, and runs a new gin.Engine using its registered
// modules.
func (s *Service) Run(log *zap.Logger) error {
	s.log = log
	gin.EnableJsonDecoderDisallowUnknownFields()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	// catch crashed
	r.Use(gin.Recovery())

	if s.AccessLog {
		r.Use(gin.Logger())
		s.log.Debug("request logging enabled")
	}
	s.LoadSession(r)
	s.Bind(r)

	if s.ACME.Enable {
		if s.Listen != "" {
			s.log.Panic("For ACME / Let's Encrypt it is not possible to set `listen`")
		}
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(s.ACME.Domains...),
			Cache:      autocert.DirCache(s.ACME.Cache),
		}
		return autotls.RunWithManager(r, &m)
	}
	return r.Run(s.Listen)
}
