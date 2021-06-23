package mailer

import (
	"errors"
	"fmt"
	"net"
	"net/textproto"

	"github.com/bdlm/log"
)

type fakeServer struct {
	s *Service
	l net.Listener
}

// NewFakeServer - to get mocked Service for mail-service
func NewFakeServer() (*fakeServer, *Service) {
	s := &Service{
		SMTPHost:     "127.0.0.1",
		SMTPPort:     12025,
		SMTPUsername: "user",
		SMTPPassword: "password",
		SMTPSSL:      false,
	}
	return newFakeServer(s)
}

func newFakeServer(s *Service) (*fakeServer, *Service) {
	fs := &fakeServer{
		s: s,
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
		}
	}
}
