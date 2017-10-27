package websocket

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
