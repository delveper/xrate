package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

// Config represents the configuration for SMTP client.
type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
}

// SMTP implements an SMTPClient interface.
type SMTP struct {
	cfg  Config
	conn net.Conn
}

func NewSMTPClient(cfg Config) *SMTP {
	return &SMTP{cfg: cfg}
}

// Connect connects to the SMTP server.
func (c *SMTP) Connect() error {
	addr := net.JoinHostPort(c.cfg.Host, c.cfg.Port)
	tlsCfg := tls.Config{InsecureSkipVerify: false, ServerName: c.cfg.Host}

	conn, err := tls.Dial("tcp", addr, &tlsCfg)
	if err != nil {
		return fmt.Errorf("connecting SMTP server: %v", err)
	}

	c.conn = conn

	return nil
}

// Auth authenticates with the SMTP server.
func (c *SMTP) Auth() error {
	auth := smtp.PlainAuth("", c.cfg.UserName, c.cfg.Password, c.cfg.Host)

	proto, resp, err := auth.Start(&smtp.ServerInfo{TLS: true})
	if err != nil {
		return fmt.Errorf("starting %s authentication: %v", proto, err)
	}

	if _, err := c.conn.Write(resp); err != nil {
		return fmt.Errorf("writing AUTH command: %v", err)
	}

	return nil
}

// Write writes the email msgTmpl to the SMTP server
func (c *SMTP) Write(msg []byte) error {
	if _, err := c.conn.Write(msg); err != nil {
		return fmt.Errorf("writing message: %v", err)
	}

	return nil
}

// Close closes the SMTP connection.
func (c *SMTP) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("closing SMTP connection: %v", err)
	}

	return nil
}
