package templates

import (
	"fmt"
	"os"
)

// GetApprovalEmailTemplate generates account approval notification
func GetApprovalEmailTemplate(name string) EmailComponents {
	return EmailComponents{
		Title: "Account Approved",
		Content: fmt.Sprintf(`
			<p>Congratulations %s!</p>
			<p>Your iRankHub account has been approved. You now have full access to our platform features.</p>
			<p>You can now:</p>
			<ul>
				<li>View and participate in tournaments</li>
				<li>Access your personalized dashboard</li>
				<li>Connect with other members</li>
				<li>Track your progress and achievements</li>
			</ul>
		`, GetHighlightTemplate(name)),
		Buttons: []EmailButton{
			{
				Text: "Get Started",
				URL:  fmt.Sprintf("%s/dashboard", os.Getenv("FRONTEND_URL")),
			},
		},
	}
}

// GetRejectionEmailTemplate generates account rejection notification
func GetRejectionEmailTemplate(name string) EmailComponents {
	return EmailComponents{
		Title: "Account Application Status",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>After careful review of your application for an iRankHub account, we regret to inform you that we are unable to approve your account at this time.</p>
			<p>This decision may be due to one or more of the following reasons:</p>
			<ul>
				<li>Incomplete or incorrect information provided</li>
				<li>Unable to verify provided credentials</li>
				<li>Does not meet current eligibility criteria</li>
			</ul>
		`, GetHighlightTemplate(name)),
		Alerts: []string{
			"If you believe this decision was made in error, please contact our support team for assistance.",
		},
		Buttons: []EmailButton{
			{
				Text: "Contact Support",
				URL:  fmt.Sprintf("%s/support", os.Getenv("FRONTEND_URL")),
			},
		},
	}
}

// GetProfileUpdateEmailTemplate generates profile update notification
func GetProfileUpdateEmailTemplate(name string, changes map[string]string) EmailComponents {
	return EmailComponents{
		Title: "Profile Updated",
		Content: fmt.Sprintf(`
			<p>Hello %s,</p>
			<p>Your iRankHub profile has been successfully updated with the following changes:</p>
		`, GetHighlightTemplate(name)),
		Metadata: changes,
		Alerts: []string{
			"If you did not make these changes, please contact our support team immediately.",
		},
		Buttons: []EmailButton{
			{
				Text: "View Profile",
				URL:  fmt.Sprintf("%s/profile", os.Getenv("FRONTEND_URL")),
			},
		},
	}
}

// GetAccountDeletionEmailTemplate generates account deletion notification
func GetAccountDeletionEmailTemplate(name string) EmailComponents {
	return EmailComponents{
		Title: "Account Deleted",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>Your iRankHub account has been successfully deleted.</p>
			<p>All your personal data has been removed from our system in accordance with our data retention policy.</p>
			<p>We're sorry to see you go. If you change your mind, you're welcome to create a new account in the future.</p>
		`, GetHighlightTemplate(name)),
		Alerts: []string{
			"If you did not request this deletion, please contact our support team immediately.",
		},
	}
}

// GetAccountDeactivationEmailTemplate generates account deactivation notification
func GetAccountDeactivationEmailTemplate(name string) EmailComponents {
	return EmailComponents{
		Title: "Account Deactivated",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>Your iRankHub account has been deactivated as requested.</p>
			<p>During this period:</p>
			<ul>
				<li>Your profile will not be visible to other users</li>
				<li>You won't receive notifications or updates</li>
				<li>Your data will be preserved for future reactivation</li>
			</ul>
		`, GetHighlightTemplate(name)),
		Alerts: []string{
			"You can reactivate your account at any time by logging in.",
		},
		Buttons: []EmailButton{
			{
				Text: "Reactivate Account",
				URL:  fmt.Sprintf("%s/login", os.Getenv("FRONTEND_URL")),
			},
		},
	}
}

// GetAccountReactivationEmailTemplate generates account reactivation notification
func GetAccountReactivationEmailTemplate(name string) EmailComponents {
	return EmailComponents{
		Title: "Account Reactivated",
		Content: fmt.Sprintf(`
			<p>Welcome back, %s!</p>
			<p>Your iRankHub account has been successfully reactivated.</p>
			<p>You now have full access to:</p>
			<ul>
				<li>Your profile and settings</li>
				<li>Tournament participation</li>
				<li>Platform features and updates</li>
			</ul>
		`, GetHighlightTemplate(name)),
		Buttons: []EmailButton{
			{
				Text: "Go to Dashboard",
				URL:  fmt.Sprintf("%s/dashboard", os.Getenv("FRONTEND_URL")),
			},
		},
	}
}
