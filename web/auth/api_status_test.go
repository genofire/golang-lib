package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web"
	"dev.sum7.eu/genofire/golang-lib/web/webtest"
)

func TestAPIStatus(t *testing.T) {
	assert := assert.New(t)
	s := webtest.New(assert)
	assert.NotNil(s)
	SetupMigration(s.DB)
	s.DB.MigrateTestdata()

	hErr := web.HTTPError{}
	// invalid
	s.Request(http.MethodGet, "/api/v1/auth/status", nil, http.StatusUnauthorized, &hErr)
	assert.Equal(APIErrorNoSession, hErr.Message)

	s.TestLogin()

	obj := User{}
	// invalid - user
	s.Request(http.MethodGet, "/api/v1/auth/status", nil, http.StatusOK, &obj)
	assert.Equal("admin", obj.Username)

}
