package web

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	assert := assert.New(t)

	ModuleRegister(func(_ *gin.Engine, _ *Service) {
	})

	s := &Service{AccessLog: true, Listen: "8.8.8.8:80"}
	// HTTP - failed
	err := s.Run()
	assert.Error(err)

	s.ACME.Enable = true
	// acme with listen port - panic
	assert.Panics(func() {
		s.Run()
	})

	s.Listen = ""
	// httpS - failed
	err = s.Run()
	assert.Error(err)
}
