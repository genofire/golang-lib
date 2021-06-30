package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web"
	"dev.sum7.eu/genofire/golang-lib/web/webtest"
)

func TestAPIDeleteMyProfil(t *testing.T) {
	assert := assert.New(t)
	s, err := webtest.New()
	assert.NoError(err)
	defer s.Close()
	assert.NotNil(s)
	SetupMigration(s.DB)
	s.DB.MigrateTestdata()

	hErr := web.HTTPError{}
	// invalid
	err = s.Request(http.MethodDelete, "/api/v1/my/profil", nil, http.StatusUnauthorized, &hErr)
	assert.NoError(err)
	assert.Equal(APIErrorNoSession, hErr.Message)

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

	s.DB.ReRun("10-data-0008-01-user")
}
