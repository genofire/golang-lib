package websocket

import (
	"testing"

	"github.com/gorilla/websocket"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSessionManager(t *testing.T) {
	assert := assert.New(t)

	session := NewSessionManager()
	assert.NotNil(session)

	out := make(chan *Message, channelBufSize)
	client := &Client{
		out:       out,
		writeQuit: make(chan bool),
		readQuit:  make(chan bool),
		ws:        &websocket.Conn{},
	}

	session.Init(client)
	msg := <-out
	assert.Equal(SessionMessageInit, msg.Subject)

	result := session.HandleMessage(nil)
	assert.False(result)

	msgFillSession := &Message{}
	result = session.HandleMessage(msgFillSession)
	assert.False(result)

	result = session.HandleMessage(&Message{
		ID:      uuid.New(),
		From:    client,
		Subject: SessionMessageInit,
	})
	assert.True(result)
}
