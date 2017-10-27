package websocket

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)

	chanMsg := make(chan *Message)
	sm := NewSessionManager()

	srv := NewServer(chanMsg, sm)

	assert.Panics(func() {
		NewClient(srv, nil)
	})

	client := NewClient(srv, &websocket.Conn{})
	assert.NotNil(client)

	client = &Client{
		server:    srv,
		id:        uuid.New(),
		out:       make(chan *Message, channelBufSize),
		writeQuit: make(chan bool),
		readQuit:  make(chan bool),
	}

	client.handleInput(&Message{})

	go client.handleInput(&Message{Subject: "a"})
	msg := <-chanMsg
	assert.Equal("a", msg.Subject)

	// msg catched by sessionManager -> not read from chanMsg needed
	client.handleInput(&Message{
		ID:      uuid.New(),
		Subject: SessionMessageInit,
	})

}
