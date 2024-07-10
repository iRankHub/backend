package utils

import (
	"fmt"
	"net/smtp"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}
	viper.AutomaticEnv()
}

func getAuthEmailTemplate(title, content string) string {
	logoURL := viper.GetString("LOGO_URL")
	if logoURL == "" {
		logoURL = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSy1c8yfmVvRgCThDUvkJTmpTrV92ANV7iSRQ&s"
	}

	return fmt.Sprintf(`
		<html>
		<head>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					padding: 20px;
					background-color: #ffffff;
				}
				h1 {
					color: #333333;
				}
				.logo {
					max-width: 200px;
					height: auto;
					margin-bottom: 20px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<img src="%s" alt="iRankHub Logo" class="logo">
				<h1>%s</h1>
				%s
			</div>
		</body>
		</html>
	`, logoURL, title, content)
}

func sendAuthEmail(to, subject, body string) error {
	from := viper.GetString("EMAIL_FROM")
	password := viper.GetString("EMAIL_PASSWORD")
	smtpHost := viper.GetString("SMTP_HOST")
	smtpPort := viper.GetString("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject = "Subject: " + subject + "\n"
	msg := []byte(subject + mime + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func SendWelcomeEmail(to, name string) error {
	subject := "Welcome to iRankHub"
	content := fmt.Sprintf(`
		<p>Welcome to iRankHub, %s!</p>
		<p>Thank you for signing up. Your account is currently pending approval.</p>
		<p>You will receive another email once your account has been reviewed.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name)
	body := getAuthEmailTemplate("Welcome to iRankHub", content)
	return sendAuthEmail(to, subject, body)
}

func SendAdminWelcomeEmail(to, name string) error {
	subject := "Welcome to iRankHub - Admin Account"
	content := fmt.Sprintf(`
		<p>Welcome to iRankHub, %s!</p>
		<p>Your admin account has been successfully created and is ready to use.</p>
		<p>You can now log in to the admin dashboard and start managing the platform.</p>
		<p>If you have any questions or need assistance, please don't hesitate to contact our support team.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name)
	body := getAuthEmailTemplate("Welcome to iRankHub - Admin Account", content)
	return sendAuthEmail(to, subject, body)
}

func SendPasswordResetEmail(to, resetToken string) error {
	subject := "Password Reset Request"
	content := fmt.Sprintf(`
		<p>We received a request to reset your password. If you didn't make this request, you can ignore this email.</p>
		<p>To reset your password, click the button below:</p>
		<p><a href="https://irankhub.com/reset-password?token=%s" style="background-color: #4CAF50; color: white; padding: 14px 20px; text-align: center; text-decoration: none; display: inline-block;">Reset Password</a></p>
		<p>This link will expire in 15 minutes.</p>
		<p>If you're having trouble, copy and paste the following URL into your web browser:</p>
		<p>https://irankhub.com/reset-password?token=%s</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, resetToken, resetToken)
	body := getAuthEmailTemplate("Password Reset Request", content)
	return sendAuthEmail(to, subject, body)
}

func SendForcedPasswordResetEmail(to, resetToken string) error {
	subject := "Security Alert: Forced Password Reset"
	content := fmt.Sprintf(`
		<p>We've detected multiple failed login attempts on your account. As a security measure, we've temporarily locked your account and are requiring a password reset.</p>
		<p>To reset your password and regain access to your account, click the button below:</p>
		<p><a href="https://irankhub.com/forced-reset-password?token=%s" style="background-color: #f44336; color: white; padding: 14px 20px; text-align: center; text-decoration: none; display: inline-block;">Reset Password Now</a></p>
		<p>This link will expire in 15 minutes.</p>
		<p>If you're having trouble, copy and paste the following URL into your web browser:</p>
		<p>https://irankhub.com/forced-reset-password?token=%s</p>
		<p>If you didn't attempt to log in recently, please contact our support team immediately as your account may be at risk.</p>
		<p>Best regards,<br>The iRankHub Security Team</p>
	`, resetToken, resetToken)
	body := getAuthEmailTemplate("Security Alert: Forced Password Reset", content)
	return sendAuthEmail(to, subject, body)
}