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
	s, err := webtest.New()
	assert.NoError(err)
	defer s.Close()
	assert.NotNil(s)
	SetupMigration(s.DB)
	s.DB.MigrateTestdata()

	hErr := web.HTTPError{}
	// invalid
	err = s.Request(http.MethodGet, "/api/v1/auth/status", nil, http.StatusUnauthorized, &hErr)
	assert.NoError(err)
	assert.Equal(APIErrorNoSession, hErr.Message)

	err = s.TestLogin()
	assert.NoError(err)

	obj := User{}
	// invalid - user
	err = s.Request(http.MethodGet, "/api/v1/auth/status", nil, http.StatusOK, &obj)
	assert.NoError(err)
	assert.Equal("admin", obj.Username)

}
