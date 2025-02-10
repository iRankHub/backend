package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/iRankHub/backend/internal/services/notification/models"
)

// CombinedStorage coordinates operations between RabbitMQ and database storage
type CombinedStorage struct {
	rabbitmq  *RabbitMQStorage
	metadata  *MetadataStorage
	retryInit sync.Once
}

func NewCombinedStorage(rabbitmq *RabbitMQStorage, metadata *MetadataStorage) *CombinedStorage {
	storage := &CombinedStorage{
		rabbitmq: rabbitmq,
		metadata: metadata,
	}

	// Initialize retry mechanism
	storage.retryInit.Do(func() {
		go storage.startRetryProcessor()
	})

	return storage
}

// Store stores a notification in both RabbitMQ and database
func (s *CombinedStorage) Store(ctx context.Context, notification *models.Notification) error {
	// First store metadata in database
	if err := s.metadata.StoreMetadata(ctx, notification); err != nil {
		return fmt.Errorf("failed to store metadata: %w", err)
	}

	// Then store in RabbitMQ
	if err := s.rabbitmq.Store(ctx, notification); err != nil {
		// If RabbitMQ storage fails, update metadata status to failed
		updateErr := s.metadata.UpdateMetadata(ctx, notification.ID, map[string]interface{}{
			"status": models.StatusFailed,
			"delivery_status": map[string]interface{}{
				"rabbitmq": map[string]interface{}{
					"status": models.StatusFailed,
					"error":  err.Error(),
				},
			},
		})
		if updateErr != nil {
			log.Printf("Failed to update metadata after RabbitMQ failure: %v", updateErr)
		}
		return fmt.Errorf("failed to store in RabbitMQ: %w", err)
	}

	return nil
}

// Get retrieves a notification from both storages and merges the data
func (s *CombinedStorage) Get(ctx context.Context, notificationID string) (*models.Notification, error) {
	// Get notification from RabbitMQ
	notification, err := s.rabbitmq.Get(ctx, notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification from RabbitMQ: %w", err)
	}

	// Get metadata from database
	metadata, err := s.metadata.GetMetadata(ctx, notificationID)
	if err != nil {
		log.Printf("Failed to get metadata for notification %s: %v", notificationID, err)
		// Continue with RabbitMQ data only
		return notification, nil
	}

	// Merge metadata into notification
	s.mergeMetadata(notification, metadata)

	return notification, nil
}

// GetUnread retrieves unread notifications from both storages
func (s *CombinedStorage) GetUnread(ctx context.Context, userID string) ([]*models.Notification, error) {
	notifications, err := s.rabbitmq.GetUnread(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread notifications from RabbitMQ: %w", err)
	}

	// Enrich notifications with metadata
	for _, notification := range notifications {
		metadata, err := s.metadata.GetMetadata(ctx, notification.ID)
		if err != nil {
			log.Printf("Failed to get metadata for notification %s: %v", notification.ID, err)
			continue
		}
		s.mergeMetadata(notification, metadata)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read in both storages
func (s *CombinedStorage) MarkAsRead(ctx context.Context, notificationID string) error {
	// Update metadata first
	if err := s.metadata.MarkAsRead(ctx, notificationID); err != nil {
		return fmt.Errorf("failed to mark metadata as read: %w", err)
	}

	// Update in RabbitMQ
	if err := s.rabbitmq.MarkAsRead(ctx, notificationID); err != nil {
		log.Printf("Failed to mark notification as read in RabbitMQ: %v", err)
		// Continue since metadata is updated
	}

	return nil
}

// UpdateDeliveryStatus updates delivery status in both storages
func (s *CombinedStorage) UpdateDeliveryStatus(ctx context.Context, notificationID string, method models.DeliveryMethod, status models.Status) error {
	// Update metadata first
	if err := s.metadata.UpdateMetadata(ctx, notificationID, map[string]interface{}{
		"delivery_status": map[string]interface{}{
			string(method): map[string]interface{}{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to update metadata delivery status: %w", err)
	}

	// Update in RabbitMQ
	if err := s.rabbitmq.UpdateDeliveryStatus(ctx, notificationID, method, status); err != nil {
		log.Printf("Failed to update delivery status in RabbitMQ: %v", err)
		// Continue since metadata is updated
	}

	return nil
}

// DeleteExpired removes expired notifications from both storages
func (s *CombinedStorage) DeleteExpired(ctx context.Context) error {
	// Delete from metadata first
	if err := s.metadata.DeleteExpiredMetadata(ctx); err != nil {
		return fmt.Errorf("failed to delete expired metadata: %w", err)
	}

	// Delete from RabbitMQ
	if err := s.rabbitmq.DeleteExpired(ctx); err != nil {
		log.Printf("Failed to delete expired notifications from RabbitMQ: %v", err)
		// Continue since metadata is cleaned up
	}

	return nil
}

// GetForRetry retrieves failed notifications that need retry
func (s *CombinedStorage) GetForRetry(ctx context.Context) ([]*models.Notification, error) {
	return s.metadata.GetNotificationsForRetry(ctx)
}

// startRetryProcessor starts a background process to handle retries
func (s *CombinedStorage) startRetryProcessor() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		if err := s.processRetries(ctx); err != nil {
			log.Printf("Error processing retries: %v", err)
		}
		cancel()
	}
}

// processRetries handles the retry logic for failed notifications
func (s *CombinedStorage) processRetries(ctx context.Context) error {
	notifications, err := s.GetForRetry(ctx)
	if err != nil {
		return fmt.Errorf("failed to get notifications for retry: %w", err)
	}

	for _, notification := range notifications {
		if !notification.CanRetry(models.EmailDelivery) {
			continue
		}

		// Attempt to store in RabbitMQ again
		if err := s.rabbitmq.Store(ctx, notification); err != nil {
			log.Printf("Retry failed for notification %s: %v", notification.ID, err)
			continue
		}

		// Update metadata status
		if err := s.metadata.UpdateMetadata(ctx, notification.ID, map[string]interface{}{
			"status": models.StatusPending,
			"delivery_status": map[string]interface{}{
				"rabbitmq": map[string]interface{}{
					"status":     models.StatusDelivered,
					"updated_at": time.Now(),
				},
			},
		}); err != nil {
			log.Printf("Failed to update metadata after successful retry: %v", err)
		}
	}

	return nil
}

// mergeMetadata merges database metadata into notification object
func (s *CombinedStorage) mergeMetadata(notification *models.Notification, metadata map[string]interface{}) {
	if isRead, ok := metadata["is_read"].(bool); ok {
		notification.IsRead = isRead
	}

	if readAt, ok := metadata["read_at"].(time.Time); ok {
		notification.ReadAt = &readAt
	}

	if deliveryStatus, ok := metadata["delivery_status"].(map[string]interface{}); ok {
		for method, status := range deliveryStatus {
			if statusMap, ok := status.(map[string]interface{}); ok {
				if statusStr, ok := statusMap["status"].(string); ok {
					if notification.DeliveryStatus == nil {
						notification.DeliveryStatus = make(map[models.DeliveryMethod]models.DeliveryStatus)
					}
					notification.DeliveryStatus[models.DeliveryMethod(method)] = models.DeliveryStatus{
						Status: models.Status(statusStr),
					}
				}
			}
		}
	}

	// In mergeMetadata function:
	if notificationMeta, ok := metadata["metadata"].(map[string]interface{}); ok {
		switch notification.Category {
		case models.ReportCategory:
			var reportMeta models.ReportMetadata
			if metaBytes, err := json.Marshal(notificationMeta); err == nil {
				if err := json.Unmarshal(metaBytes, &reportMeta); err == nil {
					// Convert back to json.RawMessage
					if metadataBytes, err := json.Marshal(reportMeta); err == nil {
						notification.Metadata = json.RawMessage(metadataBytes)
					}
				}
			}
		// Add other category cases similarly
		case models.DebateCategory:
			var debateMeta models.DebateMetadata
			if metaBytes, err := json.Marshal(notificationMeta); err == nil {
				if err := json.Unmarshal(metaBytes, &debateMeta); err == nil {
					if metadataBytes, err := json.Marshal(debateMeta); err == nil {
						notification.Metadata = json.RawMessage(metadataBytes)
					}
				}
			}
		case models.TournamentCategory:
			var tournamentMeta models.TournamentMetadata
			if metaBytes, err := json.Marshal(notificationMeta); err == nil {
				if err := json.Unmarshal(metaBytes, &tournamentMeta); err == nil {
					if metadataBytes, err := json.Marshal(tournamentMeta); err == nil {
						notification.Metadata = json.RawMessage(metadataBytes)
					}
				}
			}
		case models.UserCategory:
			var userMeta models.UserMetadata
			if metaBytes, err := json.Marshal(notificationMeta); err == nil {
				if err := json.Unmarshal(metaBytes, &userMeta); err == nil {
					if metadataBytes, err := json.Marshal(userMeta); err == nil {
						notification.Metadata = json.RawMessage(metadataBytes)
					}
				}
			}
		case models.AuthCategory:
			var authMeta models.AuthMetadata
			if metaBytes, err := json.Marshal(notificationMeta); err == nil {
				if err := json.Unmarshal(metaBytes, &authMeta); err == nil {
					if metadataBytes, err := json.Marshal(authMeta); err == nil {
						notification.Metadata = json.RawMessage(metadataBytes)
					}
				}
			}
		}
	}
}

// Close closes both storage systems
func (s *CombinedStorage) Close() error {
	var errs []error

	if err := s.rabbitmq.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close RabbitMQ storage: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing storage: %v", errs)
	}

	return nil
}
