package websocket

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	assert := assert.New(t)

	chanMsg := make(chan *Message)
	handlerService := NewWebsocketHandlerService()
	assert.NotNil(handlerService)

	handlerService.inputMSG = chanMsg
	handlerService.server.msgChanIn = chanMsg

	wg := sync.WaitGroup{}

	handlerService.SetHandler("dummy", func(msg *Message) {
		assert.Equal("expected", msg.Body)
		wg.Done()
	})
	wg.Add(1)

	handlerService.Listen("path")
	defer handlerService.Close()

	chanMsg <- &Message{Subject: "dummy", Body: "expected"}

	wg.Wait()

	wg.Add(1)
	handlerService.FallbackHandler = func(msg *Message) {
		assert.Equal("unexpected", msg.Body)
		wg.Done()
	}
	chanMsg <- &Message{Subject: "mist", Body: "unexpected"}
	wg.Wait()
}
