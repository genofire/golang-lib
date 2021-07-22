package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web"
	"dev.sum7.eu/genofire/golang-lib/web/webtest"
)

func TestAPILogin(t *testing.T) {
	assert := assert.New(t)
	s, err := webtest.NewWithDBSetup(apiLogin, SetupMigration)
	assert.NoError(err)
	defer s.Close()
	assert.NotNil(s)

	hErr := web.HTTPError{}
	// invalid
	err = s.Request(http.MethodPost, "/api/v1/auth/login", 1, http.StatusBadRequest, &hErr)
	assert.NoError(err)
	assert.Equal(web.ErrAPIInvalidRequestFormat.Error(), hErr.Message)

	req := login{}
	hErr = web.HTTPError{}
	// invalid - user
	err = s.Request(http.MethodPost, "/api/v1/auth/login", &req, http.StatusUnauthorized, &hErr)
	assert.NoError(err)
	assert.Equal(ErrAPIUserNotFound.Error(), hErr.Message)

	req.Username = "admin"
	hErr = web.HTTPError{}
	// invalid - password
	err = s.Request(http.MethodPost, "/api/v1/auth/login", &req, http.StatusUnauthorized, &hErr)
	assert.NoError(err)
	assert.Equal(ErrAPIIncorrectPassword.Error(), hErr.Message)

	req.Password = "CHANGEME"
	obj := User{}
	// valid login
	err = s.Request(http.MethodPost, "/api/v1/auth/login", &req, http.StatusOK, &obj)
	assert.NoError(err)
	assert.Equal("admin", obj.Username)
}
