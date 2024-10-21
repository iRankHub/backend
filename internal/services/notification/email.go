package notification

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type SMTPEmailSender struct{}

func NewSMTPEmailSender() *SMTPEmailSender {
	return &SMTPEmailSender{}
}

func (s *SMTPEmailSender) SendEmail(to, subject, body string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject = "Subject: " + subject + "\n"
	msg := []byte(subject + mime + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
