package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	msgChanIn      chan *Message
	clients        map[string]*Client
	clientsMutex   sync.Mutex
	SessionManager *SessionManager
	upgrader       websocket.Upgrader
}

func NewServer(msgChanIn chan *Message) *Server {
	return &Server{
		clients:        make(map[string]*Client),
		msgChanIn:      msgChanIn,
		SessionManager: NewSessionManager(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(s, conn)
	defer client.Close()
	client.Listen()
}

func (s *Server) AddClient(c *Client) {
	if c == nil {
		return
	}
	if id := c.GetID(); id != "" {
		s.clientsMutex.Lock()
		defer s.clientsMutex.Unlock()
		s.clients[id] = c
	}
}

func (s *Server) DelClient(c *Client) {
	if c == nil {
		return
	}
	if id := c.GetID(); id != "" {
		s.clientsMutex.Lock()
		delete(s.clients, id)
		s.clientsMutex.Unlock()
		s.SessionManager.Remove(c)
	}
}
