package dispatchers

import (
	"context"
	"time"

	"github.com/iRankHub/backend/internal/services/notification/models"
)

// AuthDispatcher handles authentication-related notifications
type AuthDispatcher struct {
	*BaseDispatcher
	options DispatcherOptions
}

// NewAuthDispatcher creates a new auth dispatcher
func NewAuthDispatcher(base *BaseDispatcher, options DispatcherOptions) *AuthDispatcher {
	if options.ExpirationDays == 0 {
		options.ExpirationDays = 30 // Default to 30 days
	}
	return &AuthDispatcher{
		BaseDispatcher: base,
		options:        options,
	}
}

// Dispatch implements the Dispatcher interface
func (d *AuthDispatcher) Dispatch(ctx context.Context, n *models.Notification) error {
	// Set default expiration if not set
	if n.ExpiresAt.IsZero() {
		n.ExpiresAt = time.Now().AddDate(0, 0, d.options.ExpirationDays)
	}

	// Adjust delivery methods based on notification type
	switch n.Type {
	case models.PasswordReset, models.TwoFactorAuth:
		// These are time-sensitive, use email only and shorter expiration
		n.DeliveryMethods = []models.DeliveryMethod{models.EmailDelivery}
		n.ExpiresAt = time.Now().Add(15 * time.Minute)
		n.Priority = models.HighPriority

	case models.SecurityAlert:
		// Security alerts should use all available channels
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
			models.PushDelivery,
		}
		n.Priority = models.UrgentPriority
		n.ExpiresAt = time.Now().AddDate(0, 0, 7) // Keep for 7 days

	case models.AccountCreation, models.AccountApproval:
		// Standard notifications using default channels
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		if n.Priority == "" {
			n.Priority = models.MediumPriority
		}
	}

	// Store in RabbitMQ first
	if err := d.rabbitmqSender.Publish(ctx, n); err != nil {
		return err
	}

	// Send through each delivery method
	for _, method := range n.DeliveryMethods {
		var err error
		switch method {
		case models.EmailDelivery:
			if d.options.EmailEnabled {
				err = d.emailSender.Send(ctx, n)
			}
		case models.InAppDelivery:
			if d.options.InAppEnabled {
				err = d.inAppSender.Send(ctx, n)
			}
		}

		if err != nil {
			// Log error but continue with other delivery methods
			// The notification is already in RabbitMQ for retry
			n.UpdateDeliveryStatus(method, models.StatusFailed, err)
		}
	}

	return nil
}

// Helper methods for auth-specific notifications

func (d *AuthDispatcher) SendPasswordReset(ctx context.Context, userID string, role models.UserRole, resetToken string, metadata models.AuthMetadata) error {
	n := &models.Notification{
		Category: models.AuthCategory,
		Type:     models.PasswordReset,
		UserID:   userID,
		UserRole: role,
		Title:    "Password Reset Request",
		Content:  "A password reset has been requested for your account. Click the link below to reset your password.",
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Reset Password",
				URL:   "/auth/reset-password?token=" + resetToken,
			},
		},
		Priority: models.HighPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *AuthDispatcher) SendTwoFactorCode(ctx context.Context, userID string, role models.UserRole, code string, metadata models.AuthMetadata) error {
	n := &models.Notification{
		Category: models.AuthCategory,
		Type:     models.TwoFactorAuth,
		UserID:   userID,
		UserRole: role,
		Title:    "Two-Factor Authentication Code",
		Content:  "Your verification code is: " + code,
		Priority: models.HighPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *AuthDispatcher) SendSecurityAlert(ctx context.Context, userID string, role models.UserRole, message string, metadata models.AuthMetadata) error {
	n := &models.Notification{
		Category: models.AuthCategory,
		Type:     models.SecurityAlert,
		UserID:   userID,
		UserRole: role,
		Title:    "Security Alert",
		Content:  message,
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Review Activity",
				URL:   "/settings/security",
			},
		},
		Priority: models.UrgentPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *AuthDispatcher) SendAccountCreation(ctx context.Context, userID string, role models.UserRole, metadata models.AuthMetadata) error {
	n := &models.Notification{
		Category: models.AuthCategory,
		Type:     models.AccountCreation,
		UserID:   userID,
		UserRole: role,
		Title:    "Welcome to iRankHub",
		Content:  "Your account has been created successfully. Please complete your profile to get started.",
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Complete Profile",
				URL:   "/profile/edit",
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}
