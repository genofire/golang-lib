package websocket


import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMSGValidate(t *testing.T) {
	assert := assert.New(t)
	
	msg := &Message{}
	assert.False(msg.Validate())
	
	msg.Subject = "login"
	assert.False(msg.Validate())

	msg.From = &Client{}
	assert.True(msg.Validate())

	msg.Subject = ""
	assert.False(msg.Validate())
}
