package templates

import (
	"fmt"
	"os"
)

// GetWelcomeEmailTemplate generates welcome email for new users
func GetWelcomeEmailTemplate(name string) EmailComponents {
	return EmailComponents{
		Title: "Welcome to iRankHub",
		Content: fmt.Sprintf(`
			<p>Hello %s,</p>
			<p>Welcome to iRankHub! Your account is currently pending approval.</p>
			<p>Our team will review your application and notify you once your account has been approved.</p>
			<p>In the meantime, you can prepare by:</p>
			<ul>
				<li>Reading our platform guidelines</li>
				<li>Reviewing the tournament formats</li>
				<li>Exploring available resources</li>
			</ul>
		`, GetHighlightTemplate(name)),
		Alerts: []string{
			"You will receive another email once your account has been reviewed.",
		},
		Buttons: []EmailButton{
			{
				Text: "View Guidelines",
				URL:  fmt.Sprintf("%s/docs", os.Getenv("FRONTEND_URL")),
			},
		},
	}
}

// GetAdminWelcomeEmailTemplate generates welcome email for new admin users
func GetAdminWelcomeEmailTemplate(name string) EmailComponents {
	return EmailComponents{
		Title: "Welcome to iRankHub - Admin Account",
		Content: fmt.Sprintf(`
			<p>Hello %s,</p>
			<p>Your admin account has been successfully created and is ready to use.</p>
			<p>As an admin, you have access to:</p>
			<ul>
				<li>User management dashboard</li>
				<li>Tournament administration tools</li>
				<li>System configuration settings</li>
				<li>Analytics and reporting</li>
			</ul>
		`, GetHighlightTemplate(name)),
		Buttons: []EmailButton{
			{
				Text: "Access Dashboard",
				URL:  fmt.Sprintf("%s/admin", os.Getenv("FRONTEND_URL")),
			},
		},
	}
}

// GetPasswordResetEmailTemplate generates password reset email
func GetPasswordResetEmailTemplate(resetToken string) EmailComponents {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), resetToken)
	return EmailComponents{
		Title: "Password Reset Request",
		Content: `
			<p>We received a request to reset your password.</p>
			<p>If you didn't make this request, you can safely ignore this email.</p>
		`,
		Alerts: []string{
			"This password reset link will expire in 15 minutes.",
		},
		Buttons: []EmailButton{
			{
				Text: "Reset Password",
				URL:  resetURL,
			},
		},
		Metadata: map[string]string{
			"Manual Link": resetURL,
		},
	}
}

// GetForcedPasswordResetEmailTemplate generates forced password reset email for security issues
func GetForcedPasswordResetEmailTemplate(resetToken string) EmailComponents {
	resetURL := fmt.Sprintf("%s/forced-reset-password?token=%s", os.Getenv("FRONTEND_URL"), resetToken)
	return EmailComponents{
		Title: "Security Alert: Required Password Reset",
		Content: `
			<p>We've detected multiple failed login attempts on your account.</p>
			<p>As a security measure, we've temporarily locked your account and are requiring a password reset.</p>
		`,
		Alerts: []string{
			"For your security, please reset your password immediately.",
			"If you didn't attempt to log in recently, please contact our support team as your account may be at risk.",
		},
		Buttons: []EmailButton{
			{
				Text: "Reset Password Now",
				URL:  resetURL,
			},
		},
		Metadata: map[string]string{
			"Manual Link": resetURL,
		},
	}
}

// GetTwoFactorOTPEmailTemplate generates email with 2FA code
func GetTwoFactorOTPEmailTemplate(otp string) EmailComponents {
	return EmailComponents{
		Title: "Security Verification Code",
		Content: fmt.Sprintf(`
			<p>A sign-in attempt requires further verification.</p>
			<p>To complete the sign-in, enter the verification code below:</p>
			<div style="text-align: center; padding: 20px;">
				<div style="font-size: 32px; letter-spacing: 5px; font-weight: bold; color: %s;">
					%s
				</div>
			</div>
		`, primaryColor, otp),
		Alerts: []string{
			"This code will expire in 15 minutes.",
			"If you didn't request this code, please secure your account by changing your password.",
		},
	}
}

// GetTemporaryPasswordEmailTemplate generates email with temporary password for imported users
func GetTemporaryPasswordEmailTemplate(firstName, temporaryPassword string) EmailComponents {
	return EmailComponents{
		Title: "Your Temporary Password",
		Content: fmt.Sprintf(`
			<p>Hello %s,</p>
			<p>Your iRankHub account has been created as part of a batch import process.</p>
			<p>Your temporary password is:</p>
			<div style="text-align: center; padding: 20px;">
				<div style="font-size: 24px; font-weight: bold; color: %s;">
					%s
				</div>
			</div>
		`, GetHighlightTemplate(firstName), primaryColor, temporaryPassword),
		Alerts: []string{
			"For security reasons, please change your password immediately after logging in.",
		},
		Buttons: []EmailButton{
			{
				Text: "Log In Now",
				URL:  fmt.Sprintf("%s/login", os.Getenv("FRONTEND_URL")),
			},
		},
	}
}
