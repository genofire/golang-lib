package websocket

import (
	"errors"

	"github.com/google/uuid"
)

// Message which send over websocket
type Message struct {
	server  *Server
	ID      uuid.UUID   `json:"id,omitempty"`
	Session uuid.UUID   `json:"-"`
	From    *Client     `json:"-"`
	Subject string      `json:"subject,omitempty"`
	Body    interface{} `json:"body,omitempty"`
}

// Validate is it valid message to forward
func (msg *Message) Validate() (bool, error) {
	if msg.Subject == "" {
		return false, errors.New("no subject definied")
	}
	if msg.From == nil {
		return false, errors.New("no sender definied")
	}
	return true, nil
}

// Replay to request
func (msg *Message) Replay(body interface{}) error {
	return msg.Answer(msg.Subject, body)
}

// Answer to replay at a request
func (msg *Message) Answer(subject string, body interface{}) error {
	if msg.From == nil {
		return errors.New("Message not received by a websocket Server")
	}
	msg.From.Write(&Message{
		ID:      msg.ID,
		Session: msg.Session,
		From:    msg.From,
		Subject: subject,
		Body:    body,
	})
	return nil
}

// ReplaySession to replay all of current Session
func (msg *Message) ReplaySession(body interface{}) error {
	return msg.AnswerSession(msg.Subject, body)
}

// AnswerSession to replay all of current Session
func (msg *Message) AnswerSession(subject string, body interface{}) error {
	if msg.server == nil {
		return errors.New("Message not received by a websocket Server")
	}
	if msg.server.sessionManager == nil {
		return errors.New("websocket Server run without SessionManager")
	}
	msg.server.sessionManager.Send(msg.Session, &Message{
		ID:      msg.ID,
		Session: msg.Session,
		From:    msg.From,
		Subject: subject,
		Body:    body,
	})
	return nil
}

// ReplayEverybody to replay all connection on Server
func (msg *Message) ReplayEverybody(body interface{}) error {
	return msg.AnswerEverybody(msg.Subject, body)
}

// AnswerEverybody to replay all connection on Server
func (msg *Message) AnswerEverybody(subject string, body interface{}) error {
	if msg.server == nil {
		return errors.New("Message not received by a websocket Server")
	}
	msg.server.SendAll(&Message{
		ID:      msg.ID,
		Session: msg.Session,
		From:    msg.From,
		Subject: subject,
		Body:    body,
	})
	return nil
}
