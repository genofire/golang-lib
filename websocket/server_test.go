package websocket

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	assert := assert.New(t)

	srv := NewServer(nil, NewSessionManager())
	assert.NotNil(srv)

	req, _ := http.NewRequest("GET", "url", nil)
	w := httptest.NewRecorder()
	srv.Handler(w, req)

	out := make(chan *Message)
	c := &Client{
		out:    out,
		server: srv,
	}
	srv.AddClient(nil)
	go srv.AddClient(c)
	msg := <-out
	assert.Equal(SessionMessageInit, msg.Subject)

	srv.DelClient(nil)
	srv.DelClient(c)
}

func TestServerSendAll(t *testing.T) {
	assert := assert.New(t)
	srv := NewServer(nil, nil)
	assert.NotNil(srv)

	out1 := make(chan *Message)
	c1 := &Client{
		id:     uuid.New(),
		out:    out1,
		server: srv,
	}

	out2 := make(chan *Message)
	c2 := &Client{
		id:     uuid.New(),
		out:    out2,
		server: srv,
	}
	srv.AddClient(c1)
	srv.AddClient(c2)

	go func() {

		msg := <-out1
		assert.Equal("hi", msg.Subject)

	}()
	go func() {

		msg := <-out2
		assert.Equal("hi", msg.Subject)

	}()

	srv.SendAll(&Message{
		Subject: "hi",
	})

	srv.DelClient(c2)
	srv.DelClient(c1)
}
