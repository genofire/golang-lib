package mailer

import (
	"github.com/go-mail/mail"
)

// Service to send mail
type Service struct {
	SMTPHost     string       `config:"smtp_host" toml:"smtp_host"`
	SMTPPort     int          `config:"smtp_port" toml:"smtp_port"`
	SMTPUsername string       `config:"smtp_username" toml:"smtp_username"`
	SMTPPassword string       `config:"smtp_password" toml:"smtp_password"`
	SMTPSSL      bool         `config:"smtp_ssl"  toml:"smtp_ssl"`
	Dailer       *mail.Dialer `config:"-" toml:"-"`
	From         string       `config:"from" toml:"from"`
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
