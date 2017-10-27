package websocket

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSessionManager(t *testing.T) {
	assert := assert.New(t)

	session := NewSessionManager()
	assert.NotNil(session)

	out := make(chan *Message, channelBufSize)
	client := &Client{
		id:        uuid.New(),
		out:       out,
		writeQuit: make(chan bool),
		readQuit:  make(chan bool),
	}

	go session.Init(client)
	msg := <-out
	assert.Equal(SessionMessageInit, msg.Subject)

	result := session.HandleMessage(nil)
	assert.False(result)

	msgFillSession := &Message{}
	result = session.HandleMessage(msgFillSession)
	assert.False(result)

	result = session.HandleMessage(&Message{
		ID:   uuid.New(),
		From: client,
	})
	assert.False(result)

	sessionID := uuid.New()
	result = session.HandleMessage(&Message{
		ID:      sessionID,
		From:    client,
		Subject: SessionMessageInit,
	})
	assert.True(result)

	go session.Send(sessionID, &Message{
		Subject: "some trash",
	})
	msg = <-out
	assert.Equal("some trash", msg.Subject)

	// a client need to disconnected
	c, s := session.Remove(nil)
	assert.False(c)
	assert.False(s)

	out2 := make(chan *Message, channelBufSize)
	client2 := &Client{
		id:        uuid.New(),
		out:       out2,
		writeQuit: make(chan bool),
		readQuit:  make(chan bool),
	}

	go session.Init(client2)
	msg = <-out2
	result = session.HandleMessage(&Message{
		ID:      sessionID,
		From:    client2,
		Subject: SessionMessageInit,
	})
	assert.True(result)

	// remove first client of session
	c, s = session.Remove(client)
	assert.True(c)
	assert.False(s)

	// remove last client of session
	c, s = session.Remove(client2)
	assert.True(c)
	assert.True(s)

	// all client disconnected already
	c, s = session.Remove(client2)
	assert.False(c)
	assert.False(s)
}
