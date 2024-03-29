package status

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"codeberg.org/genofire/golang-lib/web/webtest"
)

func TestAPIStatus(t *testing.T) {
	assert := assert.New(t)
	s, err := webtest.New(Register)
	assert.NoError(err)
	defer s.Close()
	assert.NotNil(s)

	obj := Status{}
	// GET
	err = s.Request(http.MethodGet, "/api/status", nil, http.StatusOK, &obj)
	assert.NoError(err)
	assert.Equal(VERSION, obj.Version)
	assert.Equal(EXTRAS, obj.Extras)
	assert.True(obj.Up)

	UP = func() bool { return false }

	obj = Status{}
	// GET - failed status
	err = s.Request(http.MethodGet, "/api/status", nil, http.StatusInternalServerError, &obj)
	assert.NoError(err)
	assert.Equal(VERSION, obj.Version)
	assert.Equal(EXTRAS, obj.Extras)
	assert.False(obj.Up)
}
