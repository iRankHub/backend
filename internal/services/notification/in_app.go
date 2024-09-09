package notification

import (
	"context"
	"fmt"
	"strconv"

	"github.com/iRankHub/backend/internal/models"
)

type InAppSender interface {
	SendInAppNotification(to, content string) error
}

type DBInAppSender struct {
	queries *models.Queries
}

func NewDBInAppSender(queries *models.Queries) *DBInAppSender {
	return &DBInAppSender{queries: queries}
}

func (s *DBInAppSender) SendInAppNotification(to, content string) error {
	ctx := context.Background()
	userID, err := strconv.Atoi(to)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	_, err = s.queries.CreateNotification(ctx, models.CreateNotificationParams{
		Userid:  int32(userID),
		Type:    "in_app",
		Message: content,
	})
	if err != nil {
		return fmt.Errorf("failed to insert in-app notification: %w", err)
	}
	return nil
}