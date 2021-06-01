package ws

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/bdlm/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// WebsocketEndpoint to handle Request
type WebsocketEndpoint struct {
	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	publishLimiter *rate.Limiter

	subscribersMu sync.Mutex
	Subscribers   map[*Subscriber]struct{}

	// Message Handler
	handlers map[string]MessageHandleFunc
	// DefaultMessageHandler if no other handler for MessageType found
	DefaultMessageHandler MessageHandleFunc
	// Run Function on open connection by subscriper
	OnOpen SubscriberEventFunc
	// Run Function on close connection to subscriper
	OnClose SubscriberEventFunc
}

// Subscriber of websocket endpoint
type Subscriber struct {
	out       chan *Message
	closeSlow func()
}

// SubscriberEventFunc for handling connection state of Subsriber
type SubscriberEventFunc func(s *Subscriber, msg chan<- *Message)

// Message on websocket
type Message struct {
	Type       string                 `json:"type"`
	ID         *uuid.UUID             `json:"id,omitempty"`
	ReplyID    *uuid.UUID             `json:"reply_id,omitempty"`
	Body       map[string]interface{} `json:"body"`
	Subscriber *Subscriber            `json:"-"`
}

// Reply to Message
func (m *Message) Reply(msg *Message) {
	if m == nil || m.Subscriber == nil {
		return
	}
	if m.ID != nil {
		msg.ReplyID = m.ID
		if msg.ID == nil {
			id := uuid.New()
			msg.ID = &id
		}
	}
	m.Subscriber.out <- msg
}

// MessageHandleFunc for handling messages
type MessageHandleFunc func(ctx context.Context, msg *Message)

// NewEndpoint - create an empty websocket
func NewEndpoint() *WebsocketEndpoint {
	return &WebsocketEndpoint{
		publishLimiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
		Subscribers:    make(map[*Subscriber]struct{}),
		handlers:       make(map[string]MessageHandleFunc),
	}
}

// Broadcast Message to all subscriber (exclude sender of Message)
func (we *WebsocketEndpoint) Broadcast(msg *Message) {
	we.subscribersMu.Lock()
	defer we.subscribersMu.Unlock()

	we.publishLimiter.Wait(context.Background())

	for s := range we.Subscribers {
		if s == msg.Subscriber {
			continue
		}
		select {
		case s.out <- msg:
		default:
			go s.closeSlow()
		}
	}
}

// AddMessageHandler - add websocket message handler
func (we *WebsocketEndpoint) AddMessageHandler(typ string, f MessageHandleFunc) {
	we.handlers[typ] = f
}

// Handler - to register in gin webservice
func (we *WebsocketEndpoint) Handler(ctx *gin.Context) {
	c, err := websocket.Accept(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, false)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	err = we.addSubscriber(ctx, c)

	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	log.Errorf("subscriber stopped: %s", err)
}

// addSubscriber and startup of websocket endpoint
func (we *WebsocketEndpoint) addSubscriber(ctxGin *gin.Context, c *websocket.Conn) error {
	s := &Subscriber{
		out: make(chan *Message, 10),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}

	we.subscribersMu.Lock()
	we.Subscribers[s] = struct{}{}
	we.subscribersMu.Unlock()
	defer func() {
		we.subscribersMu.Lock()
		delete(we.Subscribers, s)
		we.subscribersMu.Unlock()
		if we.OnClose != nil {
			we.OnClose(s, s.out)
		}
		log.Debug("websocket closed")
	}()

	if we.OnOpen != nil {
		we.OnOpen(s, s.out)
	}

	ctx := ctxGin.Request.Context()

	go func() {
		err := we.readWorker(ctx, c, s)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
			websocket.CloseStatus(err) == websocket.StatusGoingAway {
			return
		}
		log.Errorf("websocket reading error: %s", err)
	}()

	log.Debug("websocket started")

	for {
		select {
		case msg := <-s.out:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// readWorker of subscriber
func (we *WebsocketEndpoint) readWorker(ctx context.Context, c *websocket.Conn, s *Subscriber) error {
	for {
		var msg Message
		err := wsjson.Read(ctx, c, &msg)
		if err != nil {
			return err
		}
		log.WithField("type", msg.Type).Debug("receive")
		msg.Subscriber = s
		if handler, ok := we.handlers[msg.Type]; ok {
			handler(ctx, &msg)
		} else if we.DefaultMessageHandler != nil {
			we.DefaultMessageHandler(ctx, &msg)
		}
	}
}

// writeTimeout send message to subscriber with timeout
func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg *Message) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return wsjson.Write(ctx, c, msg)
}
