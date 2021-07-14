package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web"
	"dev.sum7.eu/genofire/golang-lib/web/webtest"
)

func TestAPIPasswordCode(t *testing.T) {
	assert := assert.New(t)
	s, err := webtest.NewWithDBSetup(SetupMigration)
	assert.NoError(err)
	defer s.Close()
	assert.NotNil(s)

	forgetCode := uuid.New()
	passwordCurrent := "CHANGEME"
	passwordNew := "test"

	s.DB.DB.Model(&User{ID: TestUser1ID}).Update("forget_code", forgetCode)

	hErr := web.HTTPError{}
	// invalid
	err = s.Request(http.MethodPost, "/api/v1/auth/password/code", &passwordNew, http.StatusBadRequest, &hErr)
	assert.NoError(err)
	assert.Equal(web.APIErrorInvalidRequestFormat, hErr.Message)

	res := ""
	// set new password
	err = s.Request(http.MethodPost, "/api/v1/auth/password/code", &PasswordWithForgetCode{
		ForgetCode: forgetCode,
		Password:   passwordNew,
	}, http.StatusOK, &res)
	assert.NoError(err)
	assert.Equal("admin", res)

	hErr = web.HTTPError{}
	// set password without code
	err = s.Request(http.MethodPost, "/api/v1/auth/password/code", &PasswordWithForgetCode{
		ForgetCode: forgetCode,
		Password:   passwordCurrent,
	}, http.StatusBadRequest, &hErr)
	assert.NoError(err)
	assert.Equal(APIErrorUserNotFound, hErr.Message)

	forgetCode = uuid.New()
	s.DB.DB.Model(&User{ID: TestUser1ID}).Update("forget_code", forgetCode)

	res = ""
	// set old password
	err = s.Request(http.MethodPost, "/api/v1/auth/password/code", &PasswordWithForgetCode{
		ForgetCode: forgetCode,
		Password:   passwordCurrent,
	}, http.StatusOK, &res)
	assert.NoError(err)
	assert.Equal("admin", res)
}
