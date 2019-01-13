package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/bdlm/log"
)

// Server of websocket
type Server struct {
	msgChanIn      chan *Message
	clients        map[string]*Client
	clientsMutex   sync.Mutex
	sessionManager *SessionManager
	upgrader       websocket.Upgrader
}

// NewServer to get a new Server for websockets
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

// Handler of websocket Server for binding on webserver
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithField("modul", "websocket").Warnf("error during upgrade to websocket: %s", err)
		return
	}
	client := NewClient(s, conn)
	defer client.Close()
	client.Listen()
}

func (s *Server) addClient(c *Client) {
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

func (s *Server) delClient(c *Client) {
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

// SendAll to Send a message on every Client
func (s *Server) SendAll(msg *Message) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	for _, c := range s.clients {
		c.Write(msg)
	}
}
