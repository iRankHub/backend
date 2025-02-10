package models

import "context"

// Sender defines the interface for sending notifications
type Sender interface {
	// Send delivers a notification through a specific channel
	Send(ctx context.Context, notification *Notification) error
	// Close cleans up any resources used by the sender
	Close() error
}

// Queue defines the interface for notification queue operations
type Queue interface {
	// Publish sends a notification to the message queue
	Publish(ctx context.Context, notification *Notification) error
	// Subscribe starts listening for notifications
	Subscribe(ctx context.Context, opts SubscriptionOptions) (<-chan *Notification, func(), error)
	// Close cleans up the queue connection
	Close() error
}

// SubscriptionOptions defines options for notification subscriptions
type SubscriptionOptions struct {
	UserID    string
	UserRole  UserRole
	Category  Category
	Types     []Type
	BatchSize int
}
