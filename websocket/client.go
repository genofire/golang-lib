package websocket

import (
	"io"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const channelBufSize = 100

type Client struct {
	id        uuid.UUID
	server    *Server
	ws        *websocket.Conn
	out       chan *Message
	writeQuit chan bool
	readQuit  chan bool
}

func NewClient(s *Server, ws *websocket.Conn) *Client {

	if ws == nil {
		log.Panic("ws cannot be nil")
	}

	return &Client{
		server:    s,
		ws:        ws,
		id:        uuid.New(), // fallback id (for testing)
		out:       make(chan *Message, channelBufSize),
		writeQuit: make(chan bool),
		readQuit:  make(chan bool),
	}
}

func (c *Client) GetID() string {
	if c.ws != nil {
		return c.ws.RemoteAddr().String()
	}
	return c.id.String()
}

func (c *Client) Write(msg *Message) {
	select {
	case c.out <- msg:
	default:
		c.server.DelClient(c)
		c.Close()
	}
}

func (c *Client) Close() {
	c.writeQuit <- true
	c.readQuit <- true
	log.Info("client disconnecting...", c.GetID())
}

// Listen Write and Read request via channel
func (c *Client) Listen() {
	go c.listenWrite()
	c.server.AddClient(c)
	c.listenRead()
}

func (c *Client) handleInput(msg *Message) {
	msg.From = c
	if sm := c.server.sessionManager; sm != nil && sm.HandleMessage(msg) {
		return
	}
	if ok, err := msg.Validate(); ok {
		c.server.msgChanIn <- msg
	} else {
		log.Println("no valid msg for:", c.GetID(), "error:", err, "\nmessage:", msg)
	}
}

// Listen write request via channel
func (c *Client) listenWrite() {
	for {
		select {
		case msg := <-c.out:
			websocket.WriteJSON(c.ws, msg)

		case <-c.writeQuit:
			c.server.DelClient(c)
			close(c.out)
			close(c.writeQuit)
			return
		}
	}
}

// Listen read request via channel
func (c *Client) listenRead() {
	for {
		select {

		case <-c.readQuit:
			c.server.DelClient(c)
			close(c.readQuit)
			return

		default:
			var msg Message
			err := websocket.ReadJSON(c.ws, &msg)
			if err == io.EOF {
				return
			} else if err != nil {
				log.Println(err, c.GetID())
			} else {
				c.handleInput(&msg)
			}
		}
	}
}
