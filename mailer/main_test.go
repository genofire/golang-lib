package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-mail/mail"
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

func TestSend(t *testing.T) {
	assert := assert.New(t)

	mock, s := NewFakeServer()
	assert.NotNil(mock)
	// correct setup
	err := s.Setup()
	assert.NoError(err)

	m := mail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", "bob@example.com", "cora@example.com")
	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", "Hello Bob and Cora!")
	m.AddAlternative("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	err = s.Dailer.DialAndSend(m)
	assert.NoError(err)

	msg := <-mock.MSGS
	mock.Close()
	assert.Equal("alex@example.com", msg.Header["From"][0])
	assert.Contains(msg.Body, "Bob and Cora!")

}
