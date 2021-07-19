package web

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	TestRunTLS = ""
)

func TestRun(t *testing.T) {
	assert := assert.New(t)

	s := &Service{AccessLog: true, Listen: "8.8.8.8:80"}
	s.ModuleRegister(func(_ *gin.Engine, _ *Service) {})
	// HTTP - failed
	err := s.Run()
	assert.Error(err)

	s.ACME.Enable = true
	// acme with listen port - panic
	assert.Panics(func() {
		s.Run()
	})

	if TestRunTLS == "false" {
		return
	}
	s.Listen = ""
	// httpS - failed
	err = s.Run()
	assert.Error(err)
}
