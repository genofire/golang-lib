package metrics

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"codeberg.org/genofire/golang-lib/web/webtest"
)

func TestMetricsLoaded(t *testing.T) {
	assert := assert.New(t)
	s, err := webtest.New(Register)
	assert.NoError(err)
	defer s.Close()
	assert.NotNil(s)

	// GET
	err = s.Request(http.MethodGet, "/metrics", nil, http.StatusOK, nil)
	assert.NoError(err)

	UP = func() bool { return false }

	// GET
	err = s.Request(http.MethodGet, "/metrics", nil, http.StatusOK, nil)
	assert.NoError(err)
}
