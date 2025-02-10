package dispatchers

import (
	"context"
	"github.com/iRankHub/backend/internal/services/notification/models"
	_ "github.com/iRankHub/backend/internal/services/notification/senders"
)

// Dispatcher defines the interface for notification dispatchers
type Dispatcher interface {
	// Dispatch handles the notification sending logic for a specific category
	Dispatch(ctx context.Context, notification *models.Notification) error
}

// Factory defines a notification dispatcher factory
type Factory interface {
	// GetDispatcher returns the appropriate dispatcher for a notification category
	GetDispatcher(category models.Category) (Dispatcher, error)
}

// BaseDispatcher provides common functionality for all dispatchers
type BaseDispatcher struct {
	emailSender    models.Sender
	inAppSender    models.Sender
	rabbitmqSender models.Queue
}

// NewBaseDispatcher creates a new base dispatcher
func NewBaseDispatcher(
	emailSender models.Sender,
	inAppSender models.Sender,
	rabbitmqSender models.Queue,
) *BaseDispatcher {
	return &BaseDispatcher{
		emailSender:    emailSender,
		inAppSender:    inAppSender,
		rabbitmqSender: rabbitmqSender,
	}
}

// DispatcherOptions contains configuration for dispatchers
type DispatcherOptions struct {
	EmailEnabled   bool
	InAppEnabled   bool
	PushEnabled    bool
	Priority       models.Priority
	ExpirationDays int
}

// DefaultDispatcherOptions returns default dispatcher options
func DefaultDispatcherOptions() DispatcherOptions {
	return DispatcherOptions{
		EmailEnabled:   true,
		InAppEnabled:   true,
		PushEnabled:    false,
		Priority:       models.MediumPriority,
		ExpirationDays: 30,
	}
}
