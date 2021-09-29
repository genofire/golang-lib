package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/go-mail/mail"
)

func TestSetupAndPing(t *testing.T) {
	assert := assert.New(t)

	log := zap.L()

	mock, s := NewFakeServer(log)
	assert.NotNil(mock)
	// correct setup
	err := s.Setup()
	assert.NoError(err)
	mock.Close()

	s.SMTPPassword = "wrong"
	mock, s = newFakeServer(s, log)
	// wrong password
	err = s.Setup()
	assert.Error(err)
	mock.Close()
}

func TestSend(t *testing.T) {
	assert := assert.New(t)

	mock, s := NewFakeServer(zap.L())
	assert.NotNil(mock)
	// correct setup
	err := s.Setup()
	assert.NoError(err)

	m := mail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", "bob@example.com", "cora@example.com")
	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", "Hello Bob and Cora!")
	m.AddAlternative("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	err = s.Dailer.DialAndSend(m)
	assert.NoError(err)

	msg := <-mock.Mails
	mock.Close()
	assert.Equal(s.From, msg.Header["From"][0])
	assert.Contains(msg.Body, "Bob and Cora!")

}
