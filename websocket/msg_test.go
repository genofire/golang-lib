package websocket

import (
	"sync"
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

func TestMSGReplay(t *testing.T) {
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
		ID:      conversationID,
		Subject: "lola",
	}
	err := msg.Replay(nil)
	assert.Error(err)

	msg.From = client

	done := make(chan bool)
	defer close(done)

	go func() {
		err := msg.Replay("hi")
		assert.NoError(err)
		done <- true
	}()
	msg = <-out
	<-done

	assert.Equal(conversationID, msg.ID)
	assert.Equal(uuid.Nil, msg.Session)
	assert.Equal(client, msg.From)
	assert.Equal("lola", msg.Subject)
	assert.Equal("hi", msg.Body)
}

func TestMSGSession(t *testing.T) {
	assert := assert.New(t)

	srv := NewServer(nil, nil)
	assert.NotNil(srv)

	sessionID := uuid.New()

	conversationID := uuid.New()
	msg := &Message{
		Session: sessionID,
		ID:      conversationID,
		Subject: "lola",
	}

	err := msg.ReplaySession("error")
	assert.Error(err)

	msg.server = srv
	err = msg.ReplaySession("error")
	assert.Error(err)

	srv.sessionManager = NewSessionManager()

	out1 := make(chan *Message, 3)
	c1 := &Client{
		id:     uuid.New(),
		out:    out1,
		server: srv,
	}

	out2 := make(chan *Message, 3)
	c2 := &Client{
		id:     uuid.New(),
		out:    out2,
		server: srv,
	}
	srv.addClient(c1)
	srv.addClient(c2)

	wgSession := sync.WaitGroup{}
	wg := sync.WaitGroup{}
	client := func(out chan *Message) {
		for msg := range out {
			if msg.Subject == SessionMessageInit {
				msg.ID = sessionID
				msg.From.handleInput(msg)
				wgSession.Done()
			} else {
				assert.Equal("lola", msg.Subject)
				assert.Equal("hi", msg.Body)
				assert.Equal(conversationID, msg.ID)
				assert.Equal(sessionID, msg.Session)
				wg.Done()
			}
		}
	}
	wg.Add(2)
	wgSession.Add(2)
	go client(out1)
	go client(out2)
	wgSession.Wait()

	err = msg.ReplaySession("hi")
	assert.NoError(err)
	wg.Wait()

	srv.delClient(c2)
	srv.delClient(c1)
}

func TestMSGEverbody(t *testing.T) {
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

	conversationID := uuid.New()
	msg := &Message{
		ID:      conversationID,
		Subject: "lola",
	}
	err := msg.ReplayEverybody("error")
	assert.Error(err)

	client := func(out chan *Message) {
		msg := <-out
		assert.Equal("lola", msg.Subject)
		assert.Equal("hi", msg.Body)
		assert.Equal(conversationID, msg.ID)
		assert.Equal(uuid.Nil, msg.Session)
		wg.Done()
	}
	wg.Add(2)
	go client(out1)
	go client(out2)

	msg.server = srv
	err = msg.ReplayEverybody("hi")
	assert.NoError(err)
	wg.Wait()

	srv.delClient(c2)
	srv.delClient(c1)
}
