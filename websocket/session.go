package websocket

import (
	"sync"

	"github.com/google/uuid"
)

const SessionMessageInit = "session_init"

type SessionManager struct {
	sessionToClient map[uuid.UUID]map[string]*Client
	clientToSession map[string]uuid.UUID
	sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessionToClient: make(map[uuid.UUID]map[string]*Client),
		clientToSession: make(map[string]uuid.UUID),
	}
}

func (s *SessionManager) Init(c *Client) {
	c.Write(&Message{Subject: SessionMessageInit})
}
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
	} else if msg.From != nil {
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
			} else {
				delete(s.sessionToClient, session)
				return true, true
			}
		}
	}
	return false, false
}

func (s *SessionManager) Send(id uuid.UUID, msg *Message) {
	clients := s.sessionToClient[id]
	for _, c := range clients {
		c.Write(msg)
	}
}
