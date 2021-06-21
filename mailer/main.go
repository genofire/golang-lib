package mailer

import (
	"github.com/go-mail/mail"
)

// Service to send mail
type Service struct {
	SMTPHost     string       `toml:"smtp_host"`
	SMTPPort     int          `toml:"smtp_port"`
	SMTPUsername string       `toml:"smtp_username"`
	SMTPPassword string       `toml:"smtp_password"`
	SMTPSSL      bool         `toml:"smtp_ssl"`
	Dailer       *mail.Dialer `toml:"-"`
	From         string       `toml:"from"`
}

// Ping mailer
func (m *Service) Ping() error {
	conn, err := m.Dailer.Dial()
	if err != nil {
		return err
	}
	// do not run in timeout or reconnect errors
	return conn.Close()
}

// Setup dailer (and ping)
func (m *Service) Setup() error {
	m.Dailer = mail.NewDialer(m.SMTPHost, m.SMTPPort, m.SMTPUsername, m.SMTPPassword)
	m.Dailer.SSL = m.SMTPSSL
	return m.Ping()
}
