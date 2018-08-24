package websocket

import (
	"net/http"
	"net/http/httptest"
	"sync"
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
	srv.addClient(nil)
	go srv.addClient(c)
	msg := <-out
	assert.Equal(SessionMessageInit, msg.Subject)

	srv.delClient(nil)
	srv.delClient(c)
}

func TestServerSendAll(t *testing.T) {
	assert := assert.New(t)
	srv := NewServer(nil, nil)
	assert.NotNil(srv)

	out1 := make(chan *Message, 2)
	c1 := &Client{
		id:     uuid.New(),
		out:    out1,
		server: srv,
	}

	out2 := make(chan *Message, 2)
	c2 := &Client{
		id:     uuid.New(),
		out:    out2,
		server: srv,
	}
	srv.addClient(c1)
	srv.addClient(c2)

	wg := sync.WaitGroup{}

	client := func(out chan *Message) {
		msg := <-out
		assert.Equal("hi", msg.Subject)
		wg.Done()
	}
	wg.Add(2)
	go client(out1)
	go client(out2)

	srv.SendAll(&Message{
		Subject: "hi",
	})
	wg.Wait()

	srv.delClient(c2)
	srv.delClient(c1)
}
