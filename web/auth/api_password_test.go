package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web"
	"dev.sum7.eu/genofire/golang-lib/web/webtest"
)

func TestAPIPassword(t *testing.T) {
	assert := assert.New(t)
	s := webtest.New(assert)
	defer s.Close()
	assert.NotNil(s)
	SetupMigration(s.DB)
	s.DB.MigrateTestdata()

	passwordCurrent := "CHANGEME"
	passwordNew := "test"

	hErr := web.HTTPError{}
	// no auth
	err := s.Request(http.MethodPost, "/api/v1/my/auth/password", &passwordNew, http.StatusUnauthorized, &hErr)
	assert.NoError(err)
	assert.Equal(APIErrorNoSession, hErr.Message)

	err = s.TestLogin()
	assert.NoError(err)

	hErr = web.HTTPError{}
	// invalid
	err = s.Request(http.MethodPost, "/api/v1/my/auth/password", nil, http.StatusBadRequest, &hErr)
	assert.NoError(err)
	assert.Equal(web.APIErrorInvalidRequestFormat, hErr.Message)

	res := false
	// set new password
	err = s.Request(http.MethodPost, "/api/v1/my/auth/password", &passwordNew, http.StatusOK, &res)
	assert.NoError(err)
	assert.True(res)

	res = false
	// set old password
	err = s.Request(http.MethodPost, "/api/v1/my/auth/password", &passwordCurrent, http.StatusOK, &res)
	assert.NoError(err)
	assert.True(res)
}
