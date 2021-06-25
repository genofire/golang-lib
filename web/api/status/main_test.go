package status

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web/webtest"
)

func TestAPIStatus(t *testing.T) {
	assert := assert.New(t)
	s := webtest.New(assert)
	defer s.Close()
	assert.NotNil(s)

	obj := Status{}
	// GET
	s.Request(http.MethodGet, "/api/status", nil, http.StatusOK, &obj)
	assert.Equal(VERSION, obj.Version)
	assert.Equal(EXTRAS, obj.Extras)
	assert.True(obj.Up)

	UP = func() bool { return false }

	obj = Status{}
	// GET - failed status
	s.Request(http.MethodGet, "/api/status", nil, http.StatusInternalServerError, &obj)
	assert.Equal(VERSION, obj.Version)
	assert.Equal(EXTRAS, obj.Extras)
	assert.False(obj.Up)
}
