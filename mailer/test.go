package mailer

import (
	"errors"
	"fmt"
	"net"
	"net/textproto"

	"github.com/bdlm/log"
)

var defaultStartupPort = 12025

type fakeServer struct {
	s     *Service
	l     net.Listener
	Mails chan *TestingMail
}

// TestingMail a mail in format from test server
type TestingMail struct {
	Header textproto.MIMEHeader
	Body   string
}

// NewFakeServer - to get mocked Service for mail-service
func NewFakeServer() (*fakeServer, *Service) {
	s := &Service{
		SMTPHost:     "127.0.0.1",
		SMTPPort:     defaultStartupPort,
		SMTPUsername: "user",
		SMTPPassword: "password",
		SMTPSSL:      false,
		From:         "golang-lib@example.org",
	}
	defaultStartupPort++
	return newFakeServer(s)
}

func newFakeServer(s *Service) (*fakeServer, *Service) {
	fs := &fakeServer{
		s:     s,
		Mails: make(chan *TestingMail),
	}
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", fs.s.SMTPHost, fs.s.SMTPPort))
	if err != nil {
		log.Panicf("Error listing: %s", err)
		return nil, nil
	}
	fs.l = l
	go fs.run()
	return fs, s
}

func (fs *fakeServer) Close() {
	fs.l.Close()
}

func (fs *fakeServer) run() {
	for {
		conn, err := fs.l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Panicf("Error accepting: %s", err)
		}
		go fs.handle(conn)
	}
}
func (fs *fakeServer) handle(conn net.Conn) {
	defer conn.Close()
	c := textproto.NewConn(conn)
	defer c.Close()

	c.Cmd("220 localhost.fake ESMTP Postfix")
	s, _ := c.ReadLine()
	if len(s) < 6 || s[:4] != "EHLO" {
		c.Cmd("221 Bye")
		return
	}
	c.Cmd("250-Hello %s", s[5:])
	c.Cmd("250-PIPELINIG")
	c.Cmd("250 AUTH PLAIN")
	s, _ = c.ReadLine()
	if s == "AUTH PLAIN AHVzZXIAcGFzc3dvcmQ=" {
		c.Cmd("235 Authentication successful")
	} else {
		c.Cmd("535 Authentication failed")
		c.Cmd("221 Bye")
		return
	}
	for {
		s, _ = c.ReadLine()
		switch s {
		case "QUIT":
			c.Cmd("221 Bye")
			return
		case "DATA":
			c.Cmd("354 End data with <CR><LF>.<CR><LF>")
			head, _ := c.ReadMIMEHeader()
			data := ""
		handleMsgData:
			for {
				s, _ := c.ReadLine()
				switch s {
				case ".":
					break handleMsgData
				default:
					data = fmt.Sprintf("%s%s\n", data, s)
					c.Cmd("250 Ok")
				}

			}
			fs.Mails <- &TestingMail{
				Header: head,
				Body:   data,
			}
		default:
			// fmt.Println(s)
			// TODO : MAIL FROM: and RCPT TO:
			c.Cmd("250 Ok")
		}
	}
}
