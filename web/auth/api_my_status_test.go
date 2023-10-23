package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"codeberg.org/genofire/golang-lib/web"
	"codeberg.org/genofire/golang-lib/web/webtest"
)

func TestAPIMyStatus(t *testing.T) {
	assert := assert.New(t)
	s, err := webtest.NewWithDBSetup(Register, SetupMigration)
	assert.NoError(err)
	defer s.Close()
	assert.NotNil(s)

	hErr := web.HTTPError{}
	// invalid
	err = s.Request(http.MethodGet, "/api/v1/my/auth/status", nil, http.StatusUnauthorized, &hErr)
	assert.NoError(err)
	assert.Equal(ErrAPINoSession.Error(), hErr.Message)

	err = s.TestLogin()
	assert.NoError(err)

	obj := User{}
	// invalid - user
	err = s.Request(http.MethodGet, "/api/v1/my/auth/status", nil, http.StatusOK, &obj)
	assert.NoError(err)
	assert.Equal("admin", obj.Username)

}
