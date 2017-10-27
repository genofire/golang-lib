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
	sessionManager *SessionManager
	upgrader       websocket.Upgrader
}

func NewServer(msgChanIn chan *Message, sessionManager *SessionManager) *Server {
	return &Server{
		clients:        make(map[string]*Client),
		msgChanIn:      msgChanIn,
		sessionManager: sessionManager,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info(err)
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
		s.clients[id] = c
		s.clientsMutex.Unlock()
		if s.sessionManager != nil {
			s.sessionManager.Init(c)
		}
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
		if s.sessionManager != nil {
			s.sessionManager.Remove(c)
		}
	}
}
