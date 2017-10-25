package websocket

import (
	"errors"

	"github.com/google/uuid"
)

type Message struct {
	ID      uuid.UUID   `json:"id,omitempty"`
	Session uuid.UUID   `json:"-"`
	From    *Client     `json:"-"`
	Subject string      `json:"subject,omitempty"`
	Body    interface{} `json:"body,omitempty"`
}

func (msg *Message) Validate() (bool, error) {
	if msg.Subject == "" {
		return false, errors.New("no subject definied")
	}
	if msg.From == nil {
		return false, errors.New("no sender definied")
	}
	return true, nil
}

func (msg *Message) Answer(subject string, body interface{}) {
	msg.From.Write(&Message{
		ID:      msg.ID,
		Session: msg.Session,
		From:    msg.From,
		Subject: subject,
		Body:    body,
	})
}
