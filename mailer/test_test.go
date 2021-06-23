package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeServer(t *testing.T) {
	assert := assert.New(t)

	s := &Service{
		SMTPHost:     "127.0.0.1",
		SMTPPort:     -2,
		SMTPUsername: "user",
		SMTPPassword: "password",
		SMTPSSL:      false,
	}

	// Port
	assert.Panics(func() {
		mock, _ := newFakeServer(s)
		mock.Close()
	})
}
