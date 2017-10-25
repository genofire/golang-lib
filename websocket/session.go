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
	} else {
		id := msg.From.GetID()
		msg.Session = s.clientToSession[id]
	}

	return false
}
func (s *SessionManager) Remove(c *Client) {
	if c == nil {
		return
	}
	if id := c.GetID(); id != "" {
		session := s.clientToSession[id]
		if session != uuid.Nil {
			s.Lock()
			defer s.Unlock()
			list := s.sessionToClient[session]
			delete(list, id)
			if len(list) > 0 {
				s.sessionToClient[session] = list
			} else {
				delete(s.sessionToClient, session)
			}
		}
		delete(s.clientToSession, id)
	}

}

func (s *SessionManager) Send(id uuid.UUID, msg *Message) {
	session := s.sessionToClient[id]
	for _, c := range session {
		c.Write(msg)
	}
}
