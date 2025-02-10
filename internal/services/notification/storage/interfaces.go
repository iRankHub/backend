package storage

import (
	"context"

	"github.com/iRankHub/backend/internal/services/notification/models"
)

// NotificationStore defines the interface for notification storage operations
type NotificationStore interface {
	// Store adds a notification to the storage system
	Store(ctx context.Context, notification *models.Notification) error

	// Get retrieves a notification by ID
	Get(ctx context.Context, notificationID string) (*models.Notification, error)

	// GetUnread retrieves unread notifications for a user
	GetUnread(ctx context.Context, userID string) ([]*models.Notification, error)

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, notificationID string) error

	// UpdateDeliveryStatus updates the delivery status of a notification
	UpdateDeliveryStatus(ctx context.Context, notificationID string, method models.DeliveryMethod, status models.Status) error

	// DeleteExpired removes expired notifications
	DeleteExpired(ctx context.Context) error

	// GetForRetry retrieves failed notifications eligible for retry
	GetForRetry(ctx context.Context) ([]*models.Notification, error)

	// Close cleans up any resources
	Close() error
}

// MetadataStore defines the interface for notification metadata storage
type MetadataStore interface {
	// StoreMetadata stores notification metadata in the database
	StoreMetadata(ctx context.Context, notification *models.Notification) error

	// UpdateMetadata updates existing notification metadata
	UpdateMetadata(ctx context.Context, notificationID string, metadata map[string]interface{}) error

	// GetMetadata retrieves metadata for a notification
	GetMetadata(ctx context.Context, notificationID string) (map[string]interface{}, error)

	// DeleteExpiredMetadata removes metadata for expired notifications
	DeleteExpiredMetadata(ctx context.Context) error
}

// Storage combines notification and metadata storage
type Storage interface {
	NotificationStore
	MetadataStore
}
