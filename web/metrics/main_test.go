package metrics

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web/webtest"
)

func TestMetricsLoaded(t *testing.T) {
	assert := assert.New(t)
	s := webtest.New(assert)
	assert.NotNil(s)

	// GET
	s.Request(http.MethodGet, "/metrics", nil, http.StatusOK, nil)

	UP = func() bool { return false }

	// GET
	s.Request(http.MethodGet, "/metrics", nil, http.StatusOK, nil)
}
