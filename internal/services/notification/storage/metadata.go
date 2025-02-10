package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sqlc-dev/pqtype"

	"github.com/iRankHub/backend/internal/models"
	notificationModels "github.com/iRankHub/backend/internal/services/notification/models"
)

type MetadataStorage struct {
	queries *models.Queries
}

func NewMetadataStorage(dbConn *sql.DB) *MetadataStorage {
	return &MetadataStorage{
		queries: models.New(dbConn),
	}
}

func (s *MetadataStorage) StoreMetadata(ctx context.Context, notification *notificationModels.Notification) error {
	deliveryMethods, err := json.Marshal(notification.DeliveryMethods)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery methods: %w", err)
	}

	deliveryStatus, err := json.Marshal(notification.DeliveryStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery status: %w", err)
	}

	// Handle metadata
	metadataRaw := notification.Metadata
	var fileSize sql.NullString

	// Parse metadata to extract file size if it's a report notification
	if len(metadataRaw) > 0 {
		var reportMeta notificationModels.ReportMetadata
		if err := json.Unmarshal(metadataRaw, &reportMeta); err == nil {
			if reportMeta.FileSize != "" {
				fileSize = sql.NullString{
					String: reportMeta.FileSize,
					Valid:  true,
				}
			}
		}
	}

	params := models.CreateNotificationMetadataParams{
		NotificationID:  notification.ID,
		UserID:          notification.UserID,
		Category:        string(notification.Category),
		Type:            string(notification.Type),
		Status:          string(notification.Status),
		Priority:        string(notification.Priority),
		DeliveryMethods: deliveryMethods,
		DeliveryStatus:  deliveryStatus,
		Metadata:        pqtype.NullRawMessage{RawMessage: metadataRaw, Valid: len(metadataRaw) > 0},
		ExpiresAt:       notification.ExpiresAt,
		FileSize:        fileSize,
	}

	_, err = s.queries.CreateNotificationMetadata(ctx, params)
	return err
}

func (s *MetadataStorage) UpdateMetadata(ctx context.Context, notificationID string, updates map[string]interface{}) error {
	existingMetadata, err := s.queries.GetNotificationMetadata(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("failed to get existing metadata: %w", err)
	}

	if deliveryStatus, ok := updates["delivery_status"].(map[string]interface{}); ok {
		var currentStatus map[string]interface{}
		if err := json.Unmarshal(existingMetadata.DeliveryStatus, &currentStatus); err != nil {
			return fmt.Errorf("failed to unmarshal delivery status: %w", err)
		}

		for method, status := range deliveryStatus {
			currentStatus[method] = status
		}

		newStatus, err := json.Marshal(currentStatus)
		if err != nil {
			return fmt.Errorf("failed to marshal updated delivery status: %w", err)
		}

		if err := s.queries.UpdateNotificationDeliveryStatus(ctx, models.UpdateNotificationDeliveryStatusParams{
			NotificationID: notificationID,
			Column2:        newStatus, // Using the generated field name
		}); err != nil {
			return fmt.Errorf("failed to update delivery status: %w", err)
		}
	}

	if isRead, ok := updates["is_read"].(bool); ok && isRead {
		if err := s.queries.MarkNotificationAsRead(ctx, notificationID); err != nil {
			return fmt.Errorf("failed to mark notification as read: %w", err)
		}
	}

	if status, ok := updates["status"].(string); ok {
		deliveryStatusJSON, err := json.Marshal(updates["delivery_status"])
		if err != nil {
			return fmt.Errorf("failed to marshal delivery status: %w", err)
		}

		if err := s.queries.UpdateNotificationStatus(ctx, models.UpdateNotificationStatusParams{
			NotificationID: notificationID,
			Status:         status,
			DeliveryStatus: deliveryStatusJSON,
		}); err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
	}

	return nil
}

func (s *MetadataStorage) GetMetadata(ctx context.Context, notificationID string) (map[string]interface{}, error) {
	metadata, err := s.queries.GetNotificationMetadata(ctx, notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification metadata: %w", err)
	}

	result := map[string]interface{}{
		"notification_id": metadata.NotificationID,
		"user_id":         metadata.UserID,
		"category":        metadata.Category,
		"type":            metadata.Type,
		"status":          metadata.Status,
		"priority":        metadata.Priority,
		"is_read":         metadata.IsRead,
		"created_at":      metadata.CreatedAt,
		"updated_at":      metadata.UpdatedAt,
		"expires_at":      metadata.ExpiresAt,
	}

	var deliveryMethods []string
	if err := json.Unmarshal(metadata.DeliveryMethods, &deliveryMethods); err == nil {
		result["delivery_methods"] = deliveryMethods
	}

	var deliveryStatus map[string]interface{}
	if err := json.Unmarshal(metadata.DeliveryStatus, &deliveryStatus); err == nil {
		result["delivery_status"] = deliveryStatus
	}

	if metadata.Metadata.Valid {
		var metadataContent map[string]interface{}
		if err := json.Unmarshal([]byte(metadata.Metadata.RawMessage), &metadataContent); err == nil {
			result["metadata"] = metadataContent
		}
	}

	if metadata.FileSize.Valid {
		result["file_size"] = metadata.FileSize.String
	}

	return result, nil
}

func (s *MetadataStorage) GetNotificationsForRetry(ctx context.Context) ([]*notificationModels.Notification, error) {
	failedNotifications, err := s.queries.GetNotificationsToRetry(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications for retry: %w", err)
	}

	var results []*notificationModels.Notification
	for _, failed := range failedNotifications {
		notification := &notificationModels.Notification{
			ID:       failed.NotificationID,
			UserID:   failed.UserID,
			Category: notificationModels.Category(failed.Category),
			Type:     notificationModels.Type(failed.Type),
			Status:   notificationModels.Status(failed.Status),
			Priority: notificationModels.Priority(failed.Priority),
			IsRead:   failed.IsRead,
		}

		if err := json.Unmarshal(failed.DeliveryMethods, &notification.DeliveryMethods); err != nil {
			continue
		}

		if err := json.Unmarshal(failed.DeliveryStatus, &notification.DeliveryStatus); err != nil {
			continue
		}

		if failed.Metadata.Valid {
			var metadata interface{}
			metadataBytes := []byte(failed.Metadata.RawMessage)

			switch notification.Category {
			case notificationModels.ReportCategory:
				var m notificationModels.ReportMetadata
				if err := json.Unmarshal(metadataBytes, &m); err == nil {
					metadata = m
				}
			case notificationModels.DebateCategory:
				var m notificationModels.DebateMetadata
				if err := json.Unmarshal(metadataBytes, &m); err == nil {
					metadata = m
				}
			case notificationModels.TournamentCategory:
				var m notificationModels.TournamentMetadata
				if err := json.Unmarshal(metadataBytes, &m); err == nil {
					metadata = m
				}
			case notificationModels.UserCategory:
				var m notificationModels.UserMetadata
				if err := json.Unmarshal(metadataBytes, &m); err == nil {
					metadata = m
				}
			case notificationModels.AuthCategory:
				var m notificationModels.AuthMetadata
				if err := json.Unmarshal(metadataBytes, &m); err == nil {
					metadata = m
				}
			}

			if metadata != nil {
				metadataJSON, err := json.Marshal(metadata)
				if err == nil {
					notification.Metadata = json.RawMessage(metadataJSON)
				}
			}
		}

		results = append(results, notification)
	}

	return results, nil
}

func (s *MetadataStorage) DeleteExpiredMetadata(ctx context.Context) error {
	return s.queries.DeleteExpiredNotifications(ctx)
}

func (s *MetadataStorage) MarkAsRead(ctx context.Context, notificationID string) error {
	if err := s.queries.MarkNotificationAsRead(ctx, notificationID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}
