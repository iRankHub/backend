package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Action represents a possible user action for a notification
type Action struct {
	Type        ActionType      `json:"type"`
	Label       string          `json:"label"`
	URL         string          `json:"url,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
	Completed   bool            `json:"completed"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

// DeliveryStatus tracks the delivery status for each method
type DeliveryStatus struct {
	Status      Status     `json:"status"`
	Attempts    int        `json:"attempts"`
	LastAttempt *time.Time `json:"last_attempt,omitempty"`
	Error       string     `json:"error,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
}

// Notification represents a complete notification entity
type Notification struct {
	ID              string                            `json:"id"`
	Category        Category                          `json:"category"`
	Type            Type                              `json:"type"`
	UserID          string                            `json:"user_id"`
	UserRole        UserRole                          `json:"user_role"`
	Title           string                            `json:"title"`
	Content         string                            `json:"content"`
	DeliveryMethods []DeliveryMethod                  `json:"delivery_methods"`
	Priority        Priority                          `json:"priority"`
	Actions         []Action                          `json:"actions,omitempty"`
	Metadata        json.RawMessage                   `json:"metadata,omitempty"`
	DeliveryStatus  map[DeliveryMethod]DeliveryStatus `json:"delivery_status"`
	Status          Status                            `json:"status"`
	IsRead          bool                              `json:"is_read"`
	ReadAt          *time.Time                        `json:"read_at,omitempty"`
	CreatedAt       time.Time                         `json:"created_at"`
	UpdatedAt       time.Time                         `json:"updated_at"`
	ExpiresAt       time.Time                         `json:"expires_at"`
}

// Validate performs basic validation on the notification
func (n *Notification) Validate() error {
	if n.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if n.Category == "" {
		return fmt.Errorf("category is required")
	}
	if n.Type == "" {
		return fmt.Errorf("type is required")
	}
	if len(n.DeliveryMethods) == 0 {
		return fmt.Errorf("at least one delivery method is required")
	}
	if n.Title == "" {
		return fmt.Errorf("title is required")
	}
	if n.Content == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}

// SetMetadata sets the metadata for the notification
func (n *Notification) SetMetadata(metadata interface{}) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	n.Metadata = data
	return nil
}

// GetMetadata unmarshals the metadata into the provided interface
func (n *Notification) GetMetadata(v interface{}) error {
	if n.Metadata == nil {
		return nil
	}
	return json.Unmarshal(n.Metadata, v)
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	if !n.IsRead {
		now := time.Now()
		n.IsRead = true
		n.ReadAt = &now
		n.UpdatedAt = now
	}
}

// UpdateDeliveryStatus updates the delivery status for a specific method
func (n *Notification) UpdateDeliveryStatus(method DeliveryMethod, status Status, err error) {
	now := time.Now()
	if n.DeliveryStatus == nil {
		n.DeliveryStatus = make(map[DeliveryMethod]DeliveryStatus)
	}

	ds := n.DeliveryStatus[method]
	ds.Status = status
	ds.LastAttempt = &now
	ds.Attempts++

	if err != nil {
		ds.Error = err.Error()
	}

	if status == StatusDelivered {
		ds.DeliveredAt = &now
		ds.Error = ""
	}

	n.DeliveryStatus[method] = ds
	n.UpdatedAt = now
}

// IsExpired checks if the notification has expired
func (n *Notification) IsExpired() bool {
	return !n.ExpiresAt.IsZero() && time.Now().After(n.ExpiresAt)
}

// GetDeliveryAttempts returns the number of delivery attempts for a method
func (n *Notification) GetDeliveryAttempts(method DeliveryMethod) int {
	if status, exists := n.DeliveryStatus[method]; exists {
		return status.Attempts
	}
	return 0
}

// CanRetry determines if another delivery attempt should be made
func (n *Notification) CanRetry(method DeliveryMethod) bool {
	if status, exists := n.DeliveryStatus[method]; exists {
		if status.Status == StatusDelivered {
			return false
		}

		// For email, implement progressive retry delay
		if method == EmailDelivery {
			if status.LastAttempt == nil {
				return true
			}

			delay := time.Duration(0)
			switch status.Attempts {
			case 1:
				delay = 30 * time.Minute
			case 2:
				delay = 1 * time.Hour
			case 3:
				delay = 2 * time.Hour
			default:
				return false
			}

			return time.Since(*status.LastAttempt) >= delay
		}

		// For other methods, retry up to 3 times with 5-minute delays
		return status.Attempts < 3 && (status.LastAttempt == nil ||
			time.Since(*status.LastAttempt) >= 5*time.Minute)
	}
	return true
}
