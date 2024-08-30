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

func getUserEmailTemplate(title, content string) string {
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

func sendUserEmail(to, subject, body string) error {
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

func SendApprovalNotification(to, name string) error {
	subject := "Your iRankHub Account Has Been Approved"
	content := fmt.Sprintf(`
		<p>Congratulations, %s!</p>
		<p>Your iRankHub account has been approved. You can now log in and start using all the features of iRankHub.</p>
		<p>If you have any questions or need assistance, please don't hesitate to contact our support team.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name)
	body := getUserEmailTemplate("Account Approved", content)
	return sendUserEmail(to, subject, body)
}

func SendRejectionNotification(to, name string) error {
	subject := "iRankHub Account Application Status"
	content := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>We regret to inform you that your application for an iRankHub account has been rejected.</p>
		<p>If you believe this decision was made in error or if you would like to appeal this decision, please contact our support team for further assistance.</p>
		<p>Thank you for your interest in iRankHub.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name)
	body := getUserEmailTemplate("Account Application Status", content)
	return sendUserEmail(to, subject, body)
}

func SendProfileUpdateNotification(to, name string) error {
	subject := "Your iRankHub Profile Has Been Updated"
	content := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Your iRankHub profile has been successfully updated.</p>
		<p>If you did not make these changes or if you have any questions, please contact our support team immediately.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name)
	body := getUserEmailTemplate("Profile Updated", content)
	return sendUserEmail(to, subject, body)
}

func SendAccountDeletionNotification(to, name string) error {
	subject := "Your iRankHub Account Has Been Deleted"
	content := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>We're sorry to see you go. Your iRankHub account has been successfully deleted.</p>
		<p>If you did not request this action or if you have any questions, please contact our support team immediately.</p>
		<p>Thank you for being a part of iRankHub.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name)
	body := getUserEmailTemplate("Account Deletion", content)
	return sendUserEmail(to, subject, body)
}

func SendAccountDeactivationNotification(to, name string) error {
	subject := "Your iRankHub Account Has Been Deactivated"
	content := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>Your iRankHub account has been deactivated as per your request.</p>
		<p>If you wish to reactivate your account, please log in to your account and follow the reactivation instructions.</p>
		<p>If you did not request this action, please contact our support team immediately.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name)
	body := getUserEmailTemplate("Account Deactivation", content)
	return sendUserEmail(to, subject, body)
}

func SendAccountReactivationNotification(to, name string) error {
	subject := "Your iRankHub Account Has Been Reactivated"
	content := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>Your iRankHub account has been successfully reactivated.</p>
		<p>You can now log in and access all features of the platform.</p>
		<p>If you did not request this action, please contact our support team immediately.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name)
	body := getUserEmailTemplate("Account Reactivation", content)
	return sendUserEmail(to, subject, body)
}

func SendPasswordUpdateVerificationEmail(to, name, verificationCode string) error {
	subject := "iRankHub Password Update Verification"
	content := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>We received a request to update your iRankHub account password. To complete this process, please use the following verification code:</p>
		<h2>%s</h2>
		<p>This code will expire in 15 minutes. If you did not request a password update, please ignore this email and contact our support team immediately.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name, verificationCode)
	body := getUserEmailTemplate("Password Update Verification", content)
	return sendUserEmail(to, subject, body)
}

func SendPasswordUpdateConfirmationEmail(to, name string) error {
	subject := "iRankHub Password Update Confirmation"
	content := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>Your iRankHub account password has been successfully updated.</p>
		<p>If you did not make this change, please contact our support team immediately.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, name)
	body := getUserEmailTemplate("Password Update Confirmation", content)
	return sendUserEmail(to, subject, body)
}