package notification

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

type EmailSender interface {
	SendEmail(to, subject, content string) error
}

type SMTPEmailSender struct {
	host     string
	port     string
	from     string
	password string
	timeout  time.Duration
}

func NewSMTPEmailSender() (EmailSender, error) {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	if host == "" || port == "" || from == "" || password == "" {
		return nil, fmt.Errorf("missing email configuration. Required: EMAIL_FROM, EMAIL_PASSWORD, SMTP_HOST, SMTP_PORT")
	}

	return &SMTPEmailSender{
		host:     host,
		port:     port,
		from:     from,
		password: password,
		timeout:  10 * time.Second,
	}, nil
}

func (s *SMTPEmailSender) SendEmail(to, subject, content string) error {
	maxAttempts := 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := s.sendWithTimeout(to, subject, content)
		if err == nil {
			return nil
		}
		log.Printf("Attempt %d failed: %v", attempt, err)
		if attempt < maxAttempts {
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}
		return fmt.Errorf("failed to send email after %d attempts: %w", maxAttempts, err)
	}
	return nil
}

func (s *SMTPEmailSender) sendWithTimeout(to, subject, content string) error {
	done := make(chan error, 1)

	go func() {
		// Configure TLS
		tlsConfig := &tls.Config{
			ServerName:         s.host,
			InsecureSkipVerify: false,
		}

		// Connect to the SMTP Server
		addr := fmt.Sprintf("%s:%s", s.host, s.port)
		log.Printf("Attempting to connect to SMTP server at %s", addr)

		c, err := smtp.Dial(addr)
		if err != nil {
			done <- fmt.Errorf("failed to connect to SMTP server: %w", err)
			return
		}
		defer c.Close()

		// Start TLS
		if err = c.StartTLS(tlsConfig); err != nil {
			done <- fmt.Errorf("failed to start TLS: %w", err)
			return
		}

		// Auth
		auth := smtp.PlainAuth("", s.from, s.password, s.host)
		if err = c.Auth(auth); err != nil {
			done <- fmt.Errorf("failed to authenticate: %w", err)
			return
		}

		// Set the sender and recipient
		if err = c.Mail(s.from); err != nil {
			done <- fmt.Errorf("failed to set sender: %w", err)
			return
		}
		if err = c.Rcpt(to); err != nil {
			done <- fmt.Errorf("failed to set recipient: %w", err)
			return
		}

		// Send the email body
		wc, err := c.Data()
		if err != nil {
			done <- fmt.Errorf("failed to create data writer: %w", err)
			return
		}
		defer wc.Close()

		msg := fmt.Sprintf("From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n", s.from, to, subject, content)

		if _, err = fmt.Fprint(wc, msg); err != nil {
			done <- fmt.Errorf("failed to write email body: %w", err)
			return
		}

		log.Printf("Successfully sent email to %s", to)
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(s.timeout):
		return fmt.Errorf("email sending timed out after %v", s.timeout)
	}
}
