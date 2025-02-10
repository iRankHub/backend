package senders

import (
	"context"
	"crypto/tls"
	"fmt"
	_ "log"
	"net"
	"net/smtp"
	"time"

	"github.com/iRankHub/backend/internal/services/notification/models"
	"github.com/iRankHub/backend/internal/services/notification/templates"
)

type EmailSender struct {
	config EmailSenderConfig
	retry  RetryConfig
}

func NewEmailSender(config EmailSenderConfig) (*EmailSender, error) {
	if err := validateEmailConfig(config); err != nil {
		return nil, err
	}

	return &EmailSender{
		config: config,
		retry: RetryConfig{
			MaxAttempts:  4,                // Initial + 3 retries
			InitialDelay: 30 * time.Minute, // 30 minutes
			MaxDelay:     2 * time.Hour,    // 2 hours
			Factor:       2,                // Double the delay each time
		},
	}, nil
}

func (s *EmailSender) Send(ctx context.Context, notification *models.Notification) error {
	// Check if this is a retry attempt
	attempts := notification.GetDeliveryAttempts(models.EmailDelivery)
	if attempts > 0 && !notification.CanRetry(models.EmailDelivery) {
		return fmt.Errorf("max retry attempts reached or too soon to retry")
	}

	// Convert notification to email template
	emailContent, err := s.prepareEmailContent(notification)
	if err != nil {
		notification.UpdateDeliveryStatus(models.EmailDelivery, models.StatusFailed, err)
		return fmt.Errorf("failed to prepare email content: %w", err)
	}

	// Set up email metadata
	to := notification.UserID // Assuming UserID is the email address
	subject := notification.Title

	// Attempt to send email with timeout
	done := make(chan error, 1)
	go func() {
		done <- s.sendEmail(to, subject, emailContent)
	}()

	select {
	case <-ctx.Done():
		notification.UpdateDeliveryStatus(models.EmailDelivery, models.StatusFailed, ctx.Err())
		return ctx.Err()
	case err := <-done:
		if err != nil {
			notification.UpdateDeliveryStatus(models.EmailDelivery, models.StatusFailed, err)
			return err
		}
		notification.UpdateDeliveryStatus(models.EmailDelivery, models.StatusDelivered, nil)
		return nil
	case <-time.After(30 * time.Second):
		err := fmt.Errorf("email sending timed out")
		notification.UpdateDeliveryStatus(models.EmailDelivery, models.StatusFailed, err)
		return err
	}
}

func (s *EmailSender) sendEmail(to, subject, content string) error {
	// Set up TLS config
	tlsConfig := &tls.Config{
		ServerName:         s.config.Host,
		InsecureSkipVerify: false,
	}

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer client.Close()

	// Start TLS
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("StartTLS failed: %w", err)
	}

	// Authenticate
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Set sender and recipient
	if err = client.Mail(s.config.FromAddress); err != nil {
		return fmt.Errorf("setting sender failed: %w", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("setting recipient failed: %w", err)
	}

	// Send the email body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("getting data writer failed: %w", err)
	}

	msg := fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n",
		s.config.FromName, s.config.FromAddress,
		to, subject, content)

	if _, err = writer.Write([]byte(msg)); err != nil {
		writer.Close()
		return fmt.Errorf("writing email body failed: %w", err)
	}

	if err = writer.Close(); err != nil {
		return fmt.Errorf("closing data writer failed: %w", err)
	}

	return nil
}

func (s *EmailSender) prepareEmailContent(notification *models.Notification) (string, error) {
	// Convert notification actions to email buttons
	var buttons []templates.EmailButton
	for _, action := range notification.Actions {
		buttons = append(buttons, templates.EmailButton{
			Text: action.Label,
			URL:  action.URL,
		})
	}

	// Get metadata map for email template
	var metadata map[string]string
	if notification.Metadata != nil {
		if err := notification.GetMetadata(&metadata); err != nil {
			return "", fmt.Errorf("failed to get metadata: %w", err)
		}
	}

	// Prepare email components
	components := templates.EmailComponents{
		Title:    notification.Title,
		Content:  notification.Content,
		Metadata: metadata,
		Buttons:  buttons,
	}

	// Build final email content using template
	return templates.BuildEmail(components), nil
}

func (s *EmailSender) Close() error {
	// No resources to clean up
	return nil
}

func validateEmailConfig(config EmailSenderConfig) error {
	if config.Host == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if config.Port == "" {
		return fmt.Errorf("SMTP port is required")
	}
	if config.Username == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if config.Password == "" {
		return fmt.Errorf("SMTP password is required")
	}
	if config.FromAddress == "" {
		return fmt.Errorf("from address is required")
	}
	if config.FromName == "" {
		return fmt.Errorf("from name is required")
	}
	return nil
}
