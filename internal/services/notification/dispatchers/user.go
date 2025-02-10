package dispatchers

import (
	"context"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/services/notification/models"
)

// UserDispatcher handles user-related notifications
type UserDispatcher struct {
	*BaseDispatcher
	options DispatcherOptions
}

// NewUserDispatcher creates a new user dispatcher
func NewUserDispatcher(base *BaseDispatcher, options DispatcherOptions) *UserDispatcher {
	return &UserDispatcher{
		BaseDispatcher: base,
		options:        options,
	}
}

// Dispatch implements the Dispatcher interface
func (d *UserDispatcher) Dispatch(ctx context.Context, n *models.Notification) error {
	// Set default expiration if not set
	if n.ExpiresAt.IsZero() {
		n.ExpiresAt = time.Now().AddDate(0, 0, d.options.ExpirationDays)
	}

	// Adjust delivery methods based on notification type
	switch n.Type {
	case models.ProfileUpdate:
		// Profile updates should be shown in all channels
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		n.Priority = models.MediumPriority

	case models.RoleAssignment:
		// Role changes are important
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		n.Priority = models.HighPriority

	case models.StatusChange:
		// Status changes are critical
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
			models.PushDelivery,
		}
		n.Priority = models.UrgentPriority
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
			n.UpdateDeliveryStatus(method, models.StatusFailed, err)
		}
	}

	return nil
}

// Helper methods for user-specific notifications

func (d *UserDispatcher) SendProfileUpdate(ctx context.Context, userID string, role models.UserRole, metadata models.UserMetadata) error {
	content := "Your profile has been updated with the following changes:\n"
	for field, value := range metadata.Changes {
		content += fmt.Sprintf("- %s: %s\n", field, value)
	}

	n := &models.Notification{
		Category: models.UserCategory,
		Type:     models.ProfileUpdate,
		UserID:   userID,
		UserRole: role,
		Title:    "Profile Updated",
		Content:  content,
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Profile",
				URL:   "/profile",
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *UserDispatcher) SendRoleAssignment(ctx context.Context, userID string, role models.UserRole, metadata models.UserMetadata) error {
	n := &models.Notification{
		Category: models.UserCategory,
		Type:     models.RoleAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    "Role Updated",
		Content:  fmt.Sprintf("Your role has been changed from %s to %s", metadata.PreviousRole, metadata.NewRole),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Permissions",
				URL:   "/profile/permissions",
			},
		},
		Priority: models.HighPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *UserDispatcher) SendStatusChange(ctx context.Context, userID string, role models.UserRole, metadata models.UserMetadata) error {
	n := &models.Notification{
		Category: models.UserCategory,
		Type:     models.StatusChange,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Account Status Changed"),
		Content:  fmt.Sprintf("Your account status has been changed. Reason: %s", metadata.Reason),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Contact Support",
				URL:   "/support",
			},
		},
		Priority: models.UrgentPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *UserDispatcher) SendRoleExpiration(ctx context.Context, userID string, role models.UserRole, metadata models.UserMetadata) error {
	n := &models.Notification{
		Category: models.UserCategory,
		Type:     models.RoleAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    "Role Expiration Notice",
		Content:  fmt.Sprintf("Your role as %s will expire on %s", metadata.NewRole, metadata.ExpirationDate.Format("January 2, 2006")),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Review Details",
				URL:   "/profile/role",
			},
		},
		Priority: models.HighPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}
