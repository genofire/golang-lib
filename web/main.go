/*
Package web implements common functionality for web APIs using Gin and Gorm.

Modules

Modules provide functionality for a web server. A module is a function executed
before starting a server, accessing the Service and the Gin Engine. Each Service
maintains a list of modules. When it runs, it executes all of its modules.
*/
package web

import (
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
	Listen    string `toml:"listen"`
	AccessLog bool   `toml:"access_log"`
	Webroot   string `toml:"webroot"`
	ACME      struct {
		Enable  bool     `toml:"enable"`
		Domains []string `toml:"domains"`
		Cache   string   `toml:"cache"`
	} `toml:"acme"`
	Session struct {
		Name   string `toml:"name"`
		Secret string `toml:"secret"`
	} `toml:"session"`
	// internal
	DB     *gorm.DB        `toml:"-"`
	Mailer *mailer.Service `toml:"-"`

	log     *zap.Logger
	modules []ModuleRegisterFunc
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
