package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
		mock, _ := newFakeServer(s, zap.L())
		mock.Close()
	})
}
