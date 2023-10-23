package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"codeberg.org/genofire/golang-lib/web"
	"codeberg.org/genofire/golang-lib/web/webtest"
)

func TestAPIDeleteMyProfil(t *testing.T) {
	assert := assert.New(t)
	s, err := webtest.NewWithDBSetup(Register, SetupMigration)
	assert.NoError(err)
	defer s.Close()
	assert.NotNil(s)

	hErr := web.HTTPError{}
	// invalid
	err = s.Request(http.MethodDelete, "/api/v1/my/profil", nil, http.StatusUnauthorized, &hErr)
	assert.NoError(err)
	assert.Equal(ErrAPINoSession.Error(), hErr.Message)

	err = s.Login(webtest.Login{
		Username: "admin",
		Password: "CHANGEME",
	})
	assert.NoError(err)

	res := false
	// company
	err = s.Request(http.MethodDelete, "/api/v1/my/profil", nil, http.StatusOK, &res)
	assert.NoError(err)
	assert.True(true)

	s.DB.ReMigrate("10-data-0008-01-user")
}
