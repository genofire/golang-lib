package websocket

import (
	"sync"

	"github.com/google/uuid"
)

// SessionMessageInit subject in messages
const SessionMessageInit = "session_init"

// SessionManager to handle reconnected websocket
type SessionManager struct {
	sessionToClient map[uuid.UUID]map[string]*Client
	clientToSession map[string]uuid.UUID
	sync.Mutex
}

// NewSessionManager to get a new SessionManager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessionToClient: make(map[uuid.UUID]map[string]*Client),
		clientToSession: make(map[string]uuid.UUID),
	}
}

// Init Session for given Client
func (s *SessionManager) Init(c *Client) {
	c.Write(&Message{
		From:    c,
		Subject: SessionMessageInit,
	})
}

// HandleMessage of client for Session
func (s *SessionManager) HandleMessage(msg *Message) bool {
	if msg == nil {
		return false
	}
	if msg.ID != uuid.Nil && msg.Subject == SessionMessageInit && msg.From != nil {
		s.Lock()
		defer s.Unlock()
		list := s.sessionToClient[msg.ID]
		if list == nil {
			list = make(map[string]*Client)
		}
		id := msg.From.GetID()
		list[id] = msg.From
		s.clientToSession[id] = msg.ID
		s.sessionToClient[msg.ID] = list
		return true
	}
	if msg.From != nil {
		id := msg.From.GetID()
		msg.Session = s.clientToSession[id]
	}

	return false
}

// Remove clients from SessionManagerer
// - 1. result: clients removed from session manager
// - 2. result: session closed and all clients removed
func (s *SessionManager) Remove(c *Client) (client bool, session bool) {
	if c == nil {
		return false, false
	}
	if id := c.GetID(); id != "" {
		session := s.clientToSession[id]
		defer delete(s.clientToSession, id)
		if session != uuid.Nil {
			s.Lock()
			defer s.Unlock()
			clients := s.sessionToClient[session]
			delete(clients, id)
			if len(clients) > 0 {
				s.sessionToClient[session] = clients
				return true, false
			}
			delete(s.sessionToClient, session)
			return true, true
		}
	}
	return false, false
}

// Send a message to a specific Session (and all his Websocket clients)
func (s *SessionManager) Send(id uuid.UUID, msg *Message) {
	s.Lock()
	defer s.Unlock()
	clients := s.sessionToClient[id]
	for _, c := range clients {
		c.Write(msg)
	}
}
