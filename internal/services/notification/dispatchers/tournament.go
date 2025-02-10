package dispatchers

import (
	"context"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/services/notification/models"
)

// TournamentDispatcher handles tournament-related notifications
type TournamentDispatcher struct {
	*BaseDispatcher
	options DispatcherOptions
}

// NewTournamentDispatcher creates a new tournament dispatcher
func NewTournamentDispatcher(base *BaseDispatcher, options DispatcherOptions) *TournamentDispatcher {
	return &TournamentDispatcher{
		BaseDispatcher: base,
		options:        options,
	}
}

// Dispatch implements the Dispatcher interface
func (d *TournamentDispatcher) Dispatch(ctx context.Context, n *models.Notification) error {
	// Set default expiration if not set
	if n.ExpiresAt.IsZero() {
		n.ExpiresAt = time.Now().AddDate(0, 0, d.options.ExpirationDays)
	}

	// Adjust delivery methods based on notification type
	switch n.Type {
	case models.TournamentInvite:
		// Invites are important but not urgent
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		n.Priority = models.MediumPriority

	case models.TournamentRegistration, models.TournamentPayment:
		// Registration and payment notifications need immediate attention
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
			models.PushDelivery,
		}
		n.Priority = models.HighPriority

	case models.TournamentSchedule:
		// Schedule updates need to reach everyone
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
			models.PushDelivery,
		}
		n.Priority = models.HighPriority

	case models.CoordinatorAssignment:
		// Coordinator assignments are important
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		n.Priority = models.MediumPriority
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

// Helper methods for tournament-specific notifications

func (d *TournamentDispatcher) SendTournamentInvite(ctx context.Context, userID string, role models.UserRole, metadata models.TournamentMetadata) error {
	n := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.TournamentInvite,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Tournament Invitation: %s", metadata.TournamentName),
		Content:  fmt.Sprintf("You have been invited to participate in the %s tournament.", metadata.TournamentName),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Accept Invitation",
				URL:   fmt.Sprintf("/tournaments/%s/accept", metadata.TournamentID),
			},
			{
				Type:  models.ActionView,
				Label: "View Details",
				URL:   fmt.Sprintf("/tournaments/%s", metadata.TournamentID),
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *TournamentDispatcher) SendPaymentReminder(ctx context.Context, userID string, role models.UserRole, metadata models.TournamentMetadata) error {
	n := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.TournamentPayment,
		UserID:   userID,
		UserRole: role,
		Title:    "Tournament Payment Reminder",
		Content:  fmt.Sprintf("Payment for %s tournament registration is due soon.", metadata.TournamentName),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Complete Payment",
				URL:   fmt.Sprintf("/tournaments/%s/payment", metadata.TournamentID),
			},
		},
		Priority: models.HighPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *TournamentDispatcher) SendScheduleUpdate(ctx context.Context, userID string, role models.UserRole, metadata models.TournamentMetadata, changes string) error {
	n := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.TournamentSchedule,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Schedule Update: %s", metadata.TournamentName),
		Content:  fmt.Sprintf("The tournament schedule has been updated: %s", changes),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Schedule",
				URL:   fmt.Sprintf("/tournaments/%s/schedule", metadata.TournamentID),
			},
		},
		Priority: models.HighPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *TournamentDispatcher) SendCoordinatorAssignment(ctx context.Context, userID string, role models.UserRole, metadata models.TournamentMetadata) error {
	n := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.CoordinatorAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    "Tournament Coordinator Assignment",
		Content:  fmt.Sprintf("You have been assigned as coordinator for %s tournament.", metadata.TournamentName),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Tournament",
				URL:   fmt.Sprintf("/tournaments/%s", metadata.TournamentID),
			},
			{
				Type:  models.ActionView,
				Label: "Coordinator Dashboard",
				URL:   fmt.Sprintf("/tournaments/%s/coordinator", metadata.TournamentID),
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *TournamentDispatcher) SendRegistrationConfirmation(ctx context.Context, userID string, role models.UserRole, metadata models.TournamentMetadata) error {
	n := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.TournamentRegistration,
		UserID:   userID,
		UserRole: role,
		Title:    "Tournament Registration Confirmed",
		Content:  fmt.Sprintf("Your registration for %s has been confirmed.", metadata.TournamentName),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Registration",
				URL:   fmt.Sprintf("/tournaments/%s/registration", metadata.TournamentID),
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}
