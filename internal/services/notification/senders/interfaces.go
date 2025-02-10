package senders

import (
	"context"
	"github.com/iRankHub/backend/internal/services/notification/models"
	"time"
)

// Sender defines the interface for notification delivery methods
type Sender interface {
	// Send delivers a notification through a specific channel
	Send(ctx context.Context, notification *models.Notification) error
	// Close cleans up any resources used by the sender
	Close() error
}

// EmailSenderConfig holds configuration for email delivery
type EmailSenderConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
	LogoURL     string
}

// RabbitMQConfig holds configuration for RabbitMQ
type RabbitMQConfig struct {
	URL          string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	QueueName    string
	TTL          int // Time to live in milliseconds
}

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Factor       float64 // For exponential backoff
}

// DeliveryResult represents the outcome of a delivery attempt
type DeliveryResult struct {
	Success    bool
	Error      error
	RetryAfter time.Duration
	Permanent  bool // Indicates if the failure is permanent
}

// HandlerFunc defines the signature for notification handlers
type HandlerFunc func(notification *models.Notification) error

// SubscriptionOptions defines options for notification subscriptions
type SubscriptionOptions struct {
	UserID    string
	UserRole  models.UserRole
	Category  models.Category
	Types     []models.Type
	BatchSize int
}

// Subscriber defines the interface for notification subscriptions
type Subscriber interface {
	// Subscribe starts listening for notifications
	Subscribe(ctx context.Context, opts SubscriptionOptions, handler HandlerFunc) error
	// Unsubscribe stops listening for notifications
	Unsubscribe(ctx context.Context) error
}

// Publisher defines the interface for publishing notifications
type Publisher interface {
	// Publish sends a notification to the message queue
	Publish(ctx context.Context, notification *models.Notification) error
}

// Queue defines the interface for queue operations
type Queue interface {
	Publisher
	Subscriber
	// DeleteExpired removes expired notifications
	DeleteExpired(ctx context.Context) error
}
