package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type SMTPEmailSender struct {
	from       string
	password   string
	host       string
	port       string
	timeout    time.Duration
	maxRetries int
}

func NewSMTPEmailSender() (*SMTPEmailSender, error) {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	// Validate environment variables
	if from == "" || password == "" || host == "" || port == "" {
		return nil, fmt.Errorf("missing email configuration: EMAIL_FROM=%s, SMTP_HOST=%s, SMTP_PORT=%s", from, host, port)
	}

	return &SMTPEmailSender{
		from:       from,
		password:   password,
		host:       host,
		port:       port,
		timeout:    10 * time.Second, // Shorter timeout
		maxRetries: 3,                // Fewer retries
	}, nil
}

func (s *SMTPEmailSender) SendEmail(to, subject, body string) error {
	var lastErr error
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Retrying email send to %s (attempt %d/%d)", to, attempt+1, s.maxRetries)
			time.Sleep(time.Duration(attempt) * 2 * time.Second) // Exponential backoff
		}

		err := s.sendWithTimeout(to, subject, body)
		if err == nil {
			if attempt > 0 {
				log.Printf("Successfully sent email to %s after %d retries", to, attempt)
			}
			return nil
		}

		lastErr = err
		log.Printf("Failed to send email to %s: %v", to, err)
	}

	return fmt.Errorf("failed to send email after %d attempts: %v", s.maxRetries, lastErr)
}

func (s *SMTPEmailSender) sendWithTimeout(to, subject, body string) error {
	done := make(chan error, 1)

	go func() {
		auth := smtp.PlainAuth("", s.from, s.password, s.host)
		addr := fmt.Sprintf("%s:%s", s.host, s.port)

		headers := fmt.Sprintf("From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n",
			s.from, to, subject)

		msg := []byte(headers + body)

		done <- smtp.SendMail(addr, auth, s.from, []string{to}, msg)
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(s.timeout):
		return fmt.Errorf("email sending timed out after %v", s.timeout)
	}
}
