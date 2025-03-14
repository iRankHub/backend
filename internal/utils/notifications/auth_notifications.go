package utils

import (
	"fmt"
	"os"

	"github.com/iRankHub/backend/internal/services/notification"
)

func getAuthEmailTemplate(content string) string {
	logoURL := os.Getenv("LOGO_URL")
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
				%s
			</div>
		</body>
		</html>
	`, logoURL, content)
}

func SendWelcomeEmail(notificationService *notification.NotificationService, to, name string) error {
	subject := "Welcome to iRankHub"
	content := fmt.Sprintf(`
        <p>Hello, %s!</p>
        <p>Thank you for signing up. Your account is currently pending approval.</p>
        <p>You will receive another email once your account has been reviewed.</p>
        <p>Best regards,<br>The iRankHub Team</p>
    `, name)
	body := getAuthEmailTemplate(content)
	return SendNotification(notificationService, notification.EmailNotification, to, subject, body)
}

func SendAdminWelcomeEmail(notificationService *notification.NotificationService, to, name string, userID int32) error {
	subject := "Welcome to iRankHub - Admin Account"
	content := fmt.Sprintf(`
        <p>Hello, %s!</p>
        <p>Your admin account has been successfully created and is ready to use.</p>
        <p>You can now log in to the admin dashboard and start managing the platform.</p>
        <p>If you have any questions or need assistance, please don't hesitate to contact our support team.</p>
        <p>Best regards,<br>The iRankHub Team</p>
    `, name)
	body := getAuthEmailTemplate(content)

	// Send email notification
	if err := SendNotification(notificationService, notification.EmailNotification, to, subject, body); err != nil {
		return err
	}

	// Send in-app notification
	inAppContent := fmt.Sprintf("Welcome, %s! Your admin account is ready to use.", name)
	if err := SendNotification(notificationService, notification.InAppNotification, fmt.Sprintf("%d", userID), subject, inAppContent); err != nil {
		return err
	}

	return nil
}

func SendPasswordResetEmail(notificationService *notification.NotificationService, to, resetToken string) error {
	subject := "Password Reset Request"
	content := fmt.Sprintf(`
        <p>We received a request to reset your password. If you didn't make this request, you can ignore this email.</p>
        <p>To reset your password, click the button below:</p>
        <p><a href="%s/reset-password?token=%s" style="background-color: #4CAF50; color: white; padding: 14px 20px; text-align: center; text-decoration: none; display: inline-block;">Reset Password</a></p>
        <p>This link will expire in 15 minutes.</p>
        <p>If you're having trouble, copy and paste the following URL into your web browser:</p>
        <p>%s/reset-password?token=%s</p>
        <p>Best regards,<br>The iRankHub Team</p>
    `, os.Getenv("FRONTEND_URL"), resetToken, os.Getenv("FRONTEND_URL"), resetToken)
	body := getAuthEmailTemplate(content)
	return SendNotification(notificationService, notification.EmailNotification, to, subject, body)
}

func SendForcedPasswordResetEmail(notificationService *notification.NotificationService, to, resetToken string) error {
	subject := "Security Alert: Forced Password Reset"
	content := fmt.Sprintf(`
        <p>We've detected multiple failed login attempts on your account. As a security measure, we've temporarily locked your account and are requiring a password reset.</p>
        <p>To reset your password and regain access to your account, click the button below:</p>
        <p><a href="%s/forced-reset-password?token=%s" style="background-color: #f44336; color: white; padding: 14px 20px; text-align: center; text-decoration: none; display: inline-block;">Reset Password Now</a></p>
        <p>This link will expire in 15 minutes.</p>
        <p>If you're having trouble, copy and paste the following URL into your web browser:</p>
        <p>%s/forced-reset-password?token=%s</p>
        <p>If you didn't attempt to log in recently, please contact our support team immediately as your account may be at risk.</p>
        <p>Best regards,<br>The iRankHub Security Team</p>
    `, os.Getenv("FRONTEND_URL"), resetToken, os.Getenv("FRONTEND_URL"), resetToken)
	body := getAuthEmailTemplate(content)
	return SendNotification(notificationService, notification.EmailNotification, to, subject, body)
}

func SendTwoFactorOTPEmail(notificationService *notification.NotificationService, to, otp string) error {
	subject := "Security Verification: Action Required"
	content := fmt.Sprintf(`
        <p>Hello,</p>
        <p>We hope this email finds you well. We're reaching out because we've detected some unusual activity on your iRankHub account. As a precautionary measure, we need to verify your identity before proceeding.</p>
        <p>Please use the following One-Time Password (OTP) to confirm it's really you:</p>
        <h2 style="font-size: 24px; color: #4CAF50; text-align: center;">%s</h2>
        <p>This OTP will be valid for the next 15 minutes.</p>
        <p>Your account security is our top priority.</p>
        <p>Thank you for your cooperation in keeping your account safe.</p>
        <p>Best regards,<br>The iRankHub Security Team</p>
    `, otp)
	body := getAuthEmailTemplate(content)
	return SendNotification(notificationService, notification.EmailNotification, to, subject, body)
}

func SendTemporaryPasswordEmail(notificationService *notification.NotificationService, to, firstName, temporaryPassword string) error {
	subject := "Welcome to iRankHub - Your Temporary Password"
	content := fmt.Sprintf(`
        <p>Hello, %s!</p>
        <p>Welcome to iRankHub! Your account has been created as part of a batch import process.</p>
        <p>Your temporary password is: <strong>%s</strong></p>
        <p>Please log in and change your password immediately for security reasons.</p>
        <p>If you have any questions or concerns, please contact our support team.</p>
        <p>Best regards,<br>The iRankHub Team</p>
    `, firstName, temporaryPassword)
	body := getAuthEmailTemplate(content)
	return SendNotification(notificationService, notification.EmailNotification, to, subject, body)
}
