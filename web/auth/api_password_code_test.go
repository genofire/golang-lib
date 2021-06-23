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
	s := webtest.New(assert)
	assert.NotNil(s)
	SetupMigration(s.DB)
	s.DB.MigrateTestdata()

	forgetCode := uuid.New()
	passwordCurrent := "CHANGEME"
	passwordNew := "test"

	s.DB.DB.Model(&User{ID: TestUser1ID}).Update("forget_code", forgetCode)

	hErr := web.HTTPError{}
	// invalid
	s.Request(http.MethodPost, "/api/v1/auth/password/code", &passwordNew, http.StatusBadRequest, &hErr)
	assert.Equal(web.APIErrorInvalidRequestFormat, hErr.Message)

	res := ""
	// set new password
	s.Request(http.MethodPost, "/api/v1/auth/password/code", &PasswordWithForgetCode{
		ForgetCode: forgetCode,
		Password:   passwordNew,
	}, http.StatusOK, &res)
	assert.Equal("admin", res)

	hErr = web.HTTPError{}
	// set password without code
	s.Request(http.MethodPost, "/api/v1/auth/password/code", &PasswordWithForgetCode{
		ForgetCode: forgetCode,
		Password:   passwordCurrent,
	}, http.StatusBadRequest, &hErr)
	assert.Equal(APIErrorUserNotFound, hErr.Message)

	forgetCode = uuid.New()
	s.DB.DB.Model(&User{ID: TestUser1ID}).Update("forget_code", forgetCode)

	res = ""
	// set old password
	s.Request(http.MethodPost, "/api/v1/auth/password/code", &PasswordWithForgetCode{
		ForgetCode: forgetCode,
		Password:   passwordCurrent,
	}, http.StatusOK, &res)
	assert.Equal("admin", res)
}
