package common

import (
	"github.com/sampiiiii-dev/anvil_server/anvil/config"
	"github.com/sampiiiii-dev/anvil_server/anvil/logs"
	"go.uber.org/zap"
	"net"
	"net/smtp"
	"strconv"
	"sync"
	"time"
)

type SMTPClient interface {
	SendEmail(to, subject, body string, isHTML bool) error
}

var (
	once               sync.Once
	smtpClientInstance SMTPClient
)

func initializeSMTP(c *config.Config) {
	// Validate Configuration
	s := logs.HireScribe()
	if c.SMTP.Host == "" || c.SMTP.Username == "" || c.SMTP.Password == "" || c.SMTP.Port == 0 {
		s.Fatal("Invalid SMTP configuration")
	}

	// Test Connection
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(c.SMTP.Host, strconv.Itoa(c.SMTP.Port)), 10*time.Second)
	if err != nil {
		s.Fatal("Could not connect to SMTP server: %v", zap.Error(err))
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.Fatal("Could not close SMTP connection: %v", zap.Error(err))
		}
	}(conn)

	// Initialize SMTP Client
	smtpClientInstance = &realSMTPClient{
		host:     c.SMTP.Host,
		port:     c.SMTP.Port,
		username: c.SMTP.Username,
		password: c.SMTP.Password,
	}
}

func GetSMTPClient(c *config.Config) SMTPClient {
	once.Do(func() { initializeSMTP(c) })
	return smtpClientInstance
}

type realSMTPClient struct {
	host     string
	port     int
	username string
	password string
}

func (c *realSMTPClient) SendEmail(to, subject, body string, isHTML bool) error {
	auth := smtp.PlainAuth("", c.username, c.password, c.host)

	mimeType := "text/plain"
	if isHTML {
		mimeType = "text/html"
	}
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: " + mimeType + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := c.host + ":" + strconv.Itoa(c.port)
	return smtp.SendMail(addr, auth, c.username, []string{to}, msg)
}
