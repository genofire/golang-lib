package websocket

import (
	"net/http"

	"github.com/google/uuid"
)

// MessageHandleFunc for handling messages
type MessageHandleFunc func(msg *Message)

// WebsocketHandlerService to handle every Message on there Subject by Handlers
type WebsocketHandlerService struct {
	inputMSG        chan *Message
	server          *Server
	handlers        map[string]MessageHandleFunc
	FallbackHandler MessageHandleFunc
}

// NewWebsocketHandlerService with Websocket Server
func NewWebsocketHandlerService() *WebsocketHandlerService {
	ws := WebsocketHandlerService{
		handlers: make(map[string]MessageHandleFunc),
		inputMSG: make(chan *Message),
	}
	ws.server = NewServer(ws.inputMSG, NewSessionManager())
	return &ws
}

func (ws *WebsocketHandlerService) messageHandler() {
	for msg := range ws.inputMSG {
		if handler, ok := ws.handlers[msg.Subject]; ok {
			handler(msg)
		} else if ws.FallbackHandler != nil {
			ws.FallbackHandler(msg)
		}
	}
}

// SetHandler for a message type by subject
func (ws *WebsocketHandlerService) SetHandler(subject string, f MessageHandleFunc) {
	ws.handlers[subject] = f
}

// SendAll see Server.SendAll
func (ws *WebsocketHandlerService) SendAll(msg *Message) {
	if server := ws.server; server != nil {
		server.SendAll(msg)
	}
}

// SendSession see message to all connection of one session
func (ws *WebsocketHandlerService) SendSession(id uuid.UUID, msg *Message) {
	if server := ws.server; server != nil {
		if mgmt := server.sessionManager; mgmt != nil {
			mgmt.Send(id, msg)
		}
	}
}

// Listen on net/http server at `path` and start running handling
func (ws *WebsocketHandlerService) Listen(path string) {
	http.HandleFunc(path, ws.server.Handler)
	go ws.messageHandler()
}

// Close webserver
func (ws *WebsocketHandlerService) Close() {
	close(ws.inputMSG)
}
