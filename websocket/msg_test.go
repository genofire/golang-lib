package websocket

import (
	"testing"

	"github.com/google/uuid"
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

func TestMSGAnswer(t *testing.T) {
	assert := assert.New(t)

	out := make(chan *Message, channelBufSize)
	client := &Client{
		id:        uuid.New(),
		out:       out,
		writeQuit: make(chan bool),
		readQuit:  make(chan bool),
	}

	conversationID := uuid.New()

	msg := &Message{
		From: client,
		ID:   conversationID,
	}

	go msg.Answer("hi", nil)
	msg = <-out

	assert.Equal(conversationID, msg.ID)
	assert.Equal(uuid.Nil, msg.Session)
	assert.Equal(client, msg.From)
	assert.Equal("hi", msg.Subject)
	assert.Nil(msg.Body)
}
