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
	assert.NotNil(s)
	SetupMigration(s.DB)
	s.DB.MigrateTestdata()

	passwordCurrent := "CHANGEME"
	passwordNew := "test"

	hErr := web.HTTPError{}
	// no auth
	s.Request(http.MethodPost, "/api/v1/my/auth/password", &passwordNew, http.StatusUnauthorized, &hErr)
	assert.Equal(APIErrorNoSession, hErr.Message)

	s.TestLogin()

	hErr = web.HTTPError{}
	// invalid
	s.Request(http.MethodPost, "/api/v1/my/auth/password", nil, http.StatusBadRequest, &hErr)
	assert.Equal(web.APIErrorInvalidRequestFormat, hErr.Message)

	res := false
	// set new password
	s.Request(http.MethodPost, "/api/v1/my/auth/password", &passwordNew, http.StatusOK, &res)
	assert.True(res)

	res = false
	// set old password
	s.Request(http.MethodPost, "/api/v1/my/auth/password", &passwordCurrent, http.StatusOK, &res)
	assert.True(res)
}
