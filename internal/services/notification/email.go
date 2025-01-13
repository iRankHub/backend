package notification

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
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

	// Log configuration (excluding password)
	log.Printf("Initializing SMTP sender with config - Host: %s, Port: %s, From: %s",
		host, port, from)

	return &SMTPEmailSender{
		host:     host,
		port:     port,
		from:     from,
		password: password,
		timeout:  30 * time.Second, // Increased timeout
	}, nil
}

func (s *SMTPEmailSender) SendEmail(to, subject, content string) error {
	maxAttempts := 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Printf("Starting email send attempt %d of %d to %s", attempt, maxAttempts, to)

		err := s.sendWithTimeout(to, subject, content)
		if err == nil {
			log.Printf("Successfully sent email to %s on attempt %d", to, attempt)
			return nil
		}

		log.Printf("Attempt %d failed: %v", attempt, err)
		if attempt < maxAttempts {
			backoff := time.Duration(attempt) * time.Second
			log.Printf("Waiting %v before retry...", backoff)
			time.Sleep(backoff)
			continue
		}
		return fmt.Errorf("failed to send email after %d attempts: %w", maxAttempts, err)
	}
	return nil
}

func (s *SMTPEmailSender) sendWithTimeout(to, subject, content string) error {
	done := make(chan error, 1)

	go func() {
		// Step 1: DNS Resolution
		log.Printf("Attempting DNS lookup for %s", s.host)
		addrs, err := net.LookupHost(s.host)
		if err != nil {
			log.Printf("DNS lookup failed for %s: %v", s.host, err)
			done <- fmt.Errorf("DNS lookup failed: %w", err)
			return
		}
		log.Printf("Successfully resolved %s to addresses: %v", s.host, addrs)

		// Step 2: Test TCP Connection
		addr := fmt.Sprintf("%s:%s", s.host, s.port)
		dialer := &net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}

		log.Printf("Testing TCP connection to %s", addr)
		conn, err := dialer.Dial("tcp", addr)
		if err != nil {
			log.Printf("TCP connection failed: %v", err)
			done <- fmt.Errorf("TCP connection failed: %w", err)
			return
		}
		conn.Close()
		log.Printf("TCP connection test successful")

		// Step 3: SMTP Connection
		log.Printf("Initiating SMTP connection to %s", addr)
		c, err := smtp.Dial(addr)
		if err != nil {
			log.Printf("SMTP dial failed: %v", err)
			done <- fmt.Errorf("failed to connect to SMTP server: %w", err)
			return
		}
		defer func() {
			if err := c.Close(); err != nil {
				log.Printf("Warning: Failed to close SMTP connection: %v", err)
			}
		}()
		log.Printf("SMTP connection established")

		// Step 4: TLS Setup
		tlsConfig := &tls.Config{
			ServerName:         s.host,
			InsecureSkipVerify: false,
		}

		log.Printf("Starting TLS handshake")
		if err = c.StartTLS(tlsConfig); err != nil {
			log.Printf("TLS handshake failed: %v", err)
			done <- fmt.Errorf("failed to start TLS: %w", err)
			return
		}
		log.Printf("TLS handshake successful")

		// Step 5: Authentication
		log.Printf("Attempting SMTP authentication")
		auth := smtp.PlainAuth("", s.from, s.password, s.host)
		if err = c.Auth(auth); err != nil {
			log.Printf("SMTP authentication failed: %v", err)
			done <- fmt.Errorf("failed to authenticate: %w", err)
			return
		}
		log.Printf("SMTP authentication successful")

		// Step 6: Set Sender and Recipient
		log.Printf("Setting sender address: %s", s.from)
		if err = c.Mail(s.from); err != nil {
			log.Printf("Failed to set sender: %v", err)
			done <- fmt.Errorf("failed to set sender: %w", err)
			return
		}

		log.Printf("Setting recipient address: %s", to)
		if err = c.Rcpt(to); err != nil {
			log.Printf("Failed to set recipient: %v", err)
			done <- fmt.Errorf("failed to set recipient: %w", err)
			return
		}

		// Step 7: Send Email Content
		log.Printf("Opening data writer")
		wc, err := c.Data()
		if err != nil {
			log.Printf("Failed to create data writer: %v", err)
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

		log.Printf("Writing email content (length: %d bytes)", len(msg))
		if _, err = fmt.Fprint(wc, msg); err != nil {
			log.Printf("Failed to write email content: %v", err)
			done <- fmt.Errorf("failed to write email body: %w", err)
			return
		}

		log.Printf("Email content written successfully")
		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Email sending failed: %v", err)
		} else {
			log.Printf("Email sending completed successfully")
		}
		return err
	case <-time.After(s.timeout):
		log.Printf("Email sending timed out after %v", s.timeout)
		return fmt.Errorf("email sending timed out after %v", s.timeout)
	}
}
