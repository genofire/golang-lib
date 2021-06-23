package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupAndPing(t *testing.T) {
	assert := assert.New(t)

	mock, s := NewFakeServer()
	assert.NotNil(mock)
	// correct setup
	err := s.Setup()
	assert.NoError(err)
	mock.Close()

	s.SMTPPassword = "wrong"
	mock, s = newFakeServer(s)
	// wrong password
	err = s.Setup()
	assert.Error(err)
	mock.Close()
}
