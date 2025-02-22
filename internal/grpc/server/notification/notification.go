package server

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iRankHub/backend/internal/grpc/proto/notification"
	services "github.com/iRankHub/backend/internal/services/notification"
)

type notificationServer struct {
	notification.UnimplementedNotificationServiceServer
	notificationService *services.NotificationService
}

func NewNotificationServer(db *sql.DB) (notification.NotificationServiceServer, error) {
	notificationService, err := services.NewNotificationService(db)
	if err != nil {
		return nil, err
	}

	return &notificationServer{
		notificationService: notificationService,
	}, nil
}

func (s *notificationServer) SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	err := s.notificationService.SendNotification(ctx, services.Notification{
		Type:    services.NotificationType(req.Type),
		To:      req.To,
		Subject: req.Subject,
		Content: req.Content,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to send notification: %v", err)
	}
	return &notification.SendNotificationResponse{Success: true}, nil
}

func (s *notificationServer) GetUnreadNotifications(ctx context.Context, req *notification.GetUnreadNotificationsRequest) (*notification.GetUnreadNotificationsResponse, error) {
	notifications, err := s.notificationService.GetUnreadNotifications(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get unread notifications: %v", err)
	}

	var protoNotifications []*notification.Notification
	for _, n := range notifications {
		protoNotifications = append(protoNotifications, &notification.Notification{
			Id:      int32(n.Notificationid),
			Type:    notification.NotificationType(notification.NotificationType_value[n.Type]),
			To:      n.Recipientemail.String,
			Subject: n.Subject.String,
			Content: n.Message,
		})
	}

	return &notification.GetUnreadNotificationsResponse{Notifications: protoNotifications}, nil
}

func (s *notificationServer) MarkNotificationsAsRead(ctx context.Context, req *notification.MarkNotificationsAsReadRequest) (*notification.MarkNotificationsAsReadResponse, error) {
	err := s.notificationService.MarkNotificationsAsRead(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to mark notifications as read: %v", err)
	}
	return &notification.MarkNotificationsAsReadResponse{Success: true}, nil
}

func (s *notificationServer) SubscribeToNotifications(req *notification.SubscribeRequest, stream notification.NotificationService_SubscribeToNotificationsServer) error {
	userID := req.GetUserId()

	notificationChan := make(chan *notification.Notification, 100) // Buffer size of 100, adjust as needed

	s.notificationService.RegisterNotificationChannel(userID, notificationChan)

	defer s.notificationService.UnregisterNotificationChannel(userID, notificationChan)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case notif := <-notificationChan:
			event := &notification.NotificationEvent{
				Notification: notif,
			}
			if err := stream.Send(event); err != nil {
				return status.Errorf(codes.Internal, "Failed to send notification event: %v", err)
			}
		}
	}
}
