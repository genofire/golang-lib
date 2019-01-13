package websocket

import (
	"github.com/google/uuid"
	"github.com/bdlm/log"

	"github.com/gorilla/websocket"
)

const channelBufSize = 1000

// Client of Websocket Server Connection
type Client struct {
	id        uuid.UUID
	server    *Server
	ws        *websocket.Conn
	out       chan *Message
	writeQuit chan bool
	readQuit  chan bool
}

// NewClient by websocket
func NewClient(s *Server, ws *websocket.Conn) *Client {
	if ws == nil {
		log.WithField("modul", "websocket").Panic("client cannot be created without websocket")
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

// GetID of Client ( UUID or Address to Client)
func (c *Client) GetID() string {
	if c.ws != nil {
		return c.ws.RemoteAddr().String()
	}
	return c.id.String()
}

// Write Message to Client
func (c *Client) Write(msg *Message) {
	select {
	case c.out <- msg:
	default:
		c.server.delClient(c)
		c.Close()
	}
}

// Close Client
func (c *Client) Close() {
	c.writeQuit <- true
	c.readQuit <- true
	log.WithField("modul", "websocket").Info("client disconnecting...", c.GetID())
}

// Listen write and read request via channel
func (c *Client) Listen() {
	go c.listenWrite()
	c.server.addClient(c)
	c.listenRead()
}

// handleInput manage session and valide message before send to server
func (c *Client) handleInput(msg *Message) {
	msg.From = c
	if sm := c.server.sessionManager; sm != nil && sm.HandleMessage(msg) {
		return
	}
	if ok, err := msg.Validate(); ok {
		msg.server = c.server
		c.server.msgChanIn <- msg
	} else {
		log.WithField("modul", "websocket").Println("no valid msg for:", c.GetID(), "error:", err, "\nmessage:", msg)
	}
}

// listenWrite request via channel
func (c *Client) listenWrite() {
	for {
		select {
		case msg := <-c.out:
			websocket.WriteJSON(c.ws, msg)

		case <-c.writeQuit:
			c.server.delClient(c)
			close(c.out)
			close(c.writeQuit)
			return
		}
	}
}

// listenRead request via channel
func (c *Client) listenRead() {
	for {
		select {

		case <-c.readQuit:
			c.server.delClient(c)
			close(c.readQuit)
			return

		default:
			var msg Message
			err := websocket.ReadJSON(c.ws, &msg)
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				return
			} else if err != nil {
				log.WithField("modul", "websocket").Warnf("error on reading %s: %s", c.GetID(), err)
				return
			} else {
				c.handleInput(&msg)
			}
		}
	}
}
