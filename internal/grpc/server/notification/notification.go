package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/iRankHub/backend/internal/services/notification/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"strconv"
	"time"

	pb "github.com/iRankHub/backend/internal/grpc/proto/notification"
	"github.com/iRankHub/backend/internal/services/notification"
	"github.com/iRankHub/backend/internal/services/notification/dispatchers"
	"github.com/iRankHub/backend/internal/services/notification/senders"
)

type notificationServer struct {
	pb.UnimplementedNotificationServiceServer
	notificationService *notification.Service
}

func NewNotificationServer(ctx context.Context, db *sql.DB) (pb.NotificationServiceServer, error) {
	// Initialize email sender config
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %v", err)
	}

	emailConfig := senders.EmailSenderConfig{
		Host:        os.Getenv("SMTP_HOST"),
		Port:        smtpPort,
		Username:    os.Getenv("EMAIL_FROM"),
		Password:    os.Getenv("EMAIL_PASSWORD"),
		FromAddress: os.Getenv("EMAIL_FROM"),
		FromName:    "iRankHub",
		LogoURL:     os.Getenv("LOGO_URL"),
	}

	// Initialize dispatcher options
	dispatcherOpts := dispatchers.DispatcherOptions{
		EmailEnabled:   true,
		InAppEnabled:   true,
		PushEnabled:    false,
		ExpirationDays: 30,
	}

	// Create notification service
	service, err := notification.NewNotificationService(ctx, notification.ServiceConfig{
		RabbitMQURL:    os.Getenv("RABBITMQ_URL"),
		EmailConfig:    emailConfig,
		DispatcherOpts: dispatcherOpts,
	}, db)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification service: %w", err)
	}

	return &notificationServer{
		notificationService: service,
	}, nil
}

func (s *notificationServer) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	// Convert proto notification to service notification
	notification, err := protoToServiceNotification(req.GetNotification())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid notification: %v", err)
	}

	// Send notification
	if err := s.notificationService.SendNotification(ctx, notification); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send notification: %v", err)
	}

	return &pb.SendNotificationResponse{
		NotificationId: notification.ID,
		Status:         pb.Status_STATUS_DELIVERED,
	}, nil
}

func (s *notificationServer) GetUnreadNotifications(ctx context.Context, req *pb.GetUnreadNotificationsRequest) (*pb.GetUnreadNotificationsResponse, error) {
	notifications, err := s.notificationService.GetUnreadNotifications(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get unread notifications: %v", err)
	}

	protoNotifications := make([]*pb.Notification, len(notifications))
	for i, n := range notifications {
		protoNotification, err := serviceToProtoNotification(n)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert notification: %v", err)
		}
		protoNotifications[i] = protoNotification
	}

	return &pb.GetUnreadNotificationsResponse{
		Notifications: protoNotifications,
		Total:         int32(len(notifications)),
	}, nil
}

func (s *notificationServer) MarkNotificationsAsRead(ctx context.Context, req *pb.MarkNotificationsAsReadRequest) (*pb.MarkNotificationsAsReadResponse, error) {
	if err := s.notificationService.MarkAsRead(ctx, req.GetUserId(), req.GetNotificationIds()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark notifications as read: %v", err)
	}

	return &pb.MarkNotificationsAsReadResponse{
		Success: true,
	}, nil
}

func (s *notificationServer) SubscribeToNotifications(req *pb.SubscribeRequest, stream pb.NotificationService_SubscribeToNotificationsServer) error {
	ctx := stream.Context()

	// Convert subscription options
	opts := models.SubscriptionOptions{
		UserID:    req.GetUserId(),
		UserRole:  models.UserRole(req.GetUserRole().String()),
		Category:  models.Category(req.GetCategory().String()),
		Types:     make([]models.Type, len(req.GetTypes())),
		BatchSize: int(req.GetBatchSize()),
	}

	for i, t := range req.GetTypes() {
		opts.Types[i] = models.Type(t.String())
	}

	// Subscribe to notifications
	notifChan, cleanup, err := s.notificationService.Subscribe(ctx, req.GetUserId())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to subscribe to notifications: %v", err)
	}
	defer cleanup()

	// Stream notifications to client
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case notification, ok := <-notifChan:
			if !ok {
				return nil
			}

			protoNotif, err := serviceToProtoNotification(notification)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to convert notification: %v", err)
			}

			event := &pb.NotificationEvent{
				Event: &pb.NotificationEvent_Notification{
					Notification: protoNotif,
				},
			}

			if err := stream.Send(event); err != nil {
				return status.Errorf(codes.Internal, "failed to send notification event: %v", err)
			}
		}
	}
}

// Helper functions for type conversion

func protoToServiceNotification(proto *pb.Notification) (*models.Notification, error) {
	// Convert delivery methods
	deliveryMethods := make([]models.DeliveryMethod, len(proto.GetDeliveryMethods()))
	for i, dm := range proto.GetDeliveryMethods() {
		deliveryMethods[i] = models.DeliveryMethod(dm.String())
	}

	// Convert actions
	actions := make([]models.Action, len(proto.GetActions()))
	for i, a := range proto.GetActions() {
		completedAt, err := time.Parse(time.RFC3339, a.GetCompletedAt())
		if err != nil && a.GetCompletedAt() != "" {
			return nil, fmt.Errorf("invalid completed_at timestamp: %v", err)
		}

		actions[i] = models.Action{
			Type:        models.ActionType(a.GetType().String()),
			Label:       a.GetLabel(),
			URL:         a.GetUrl(),
			Data:        json.RawMessage(a.GetData()),
			Completed:   a.GetCompleted(),
			CompletedAt: &completedAt,
		}
	}

	// Parse timestamps
	readAt, _ := time.Parse(time.RFC3339, proto.GetReadAt())
	createdAt, _ := time.Parse(time.RFC3339, proto.GetCreatedAt())
	updatedAt, _ := time.Parse(time.RFC3339, proto.GetUpdatedAt())
	expiresAt, _ := time.Parse(time.RFC3339, proto.GetExpiresAt())

	// Handle metadata based on category
	var metadata json.RawMessage
	switch m := proto.GetMetadata().GetMetadata().(type) {
	case *pb.Metadata_Auth:
		metadata, _ = json.Marshal(protoToAuthMetadata(m.Auth))
	case *pb.Metadata_User:
		metadata, _ = json.Marshal(protoToUserMetadata(m.User))
	case *pb.Metadata_Tournament:
		metadata, _ = json.Marshal(protoToTournamentMetadata(m.Tournament))
	case *pb.Metadata_Debate:
		metadata, _ = json.Marshal(protoToDebateMetadata(m.Debate))
	case *pb.Metadata_Report:
		metadata, _ = json.Marshal(protoToReportMetadata(m.Report))
	}

	return &models.Notification{
		ID:              proto.GetId(),
		Category:        models.Category(proto.GetCategory().String()),
		Type:            models.Type(proto.GetType().String()),
		UserID:          proto.GetUserId(),
		UserRole:        models.UserRole(proto.GetUserRole().String()),
		Title:           proto.GetTitle(),
		Content:         proto.GetContent(),
		DeliveryMethods: deliveryMethods,
		Priority:        models.Priority(proto.GetPriority().String()),
		Actions:         actions,
		Metadata:        metadata,
		Status:          models.Status(proto.GetStatus().String()),
		IsRead:          proto.GetIsRead(),
		ReadAt:          &readAt,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		ExpiresAt:       expiresAt,
	}, nil
}

func serviceToProtoNotification(service *models.Notification) (*pb.Notification, error) {
	// Convert delivery methods
	deliveryMethods := make([]pb.DeliveryMethod, len(service.DeliveryMethods))
	for i, dm := range service.DeliveryMethods {
		val := pb.DeliveryMethod_value[string(dm)]
		deliveryMethods[i] = pb.DeliveryMethod(val)
	}

	// Convert actions
	actions := make([]*pb.Action, len(service.Actions))
	for i, a := range service.Actions {
		completedAt := ""
		if a.CompletedAt != nil {
			completedAt = a.CompletedAt.Format(time.RFC3339)
		}

		actions[i] = &pb.Action{
			Type:        pb.ActionType(pb.ActionType_value[string(a.Type)]),
			Label:       a.Label,
			Url:         a.URL,
			Data:        string(a.Data),
			Completed:   a.Completed,
			CompletedAt: completedAt,
		}
	}

	// Handle timestamps
	readAt := ""
	if service.ReadAt != nil {
		readAt = service.ReadAt.Format(time.RFC3339)
	}

	// Convert metadata based on category
	var metadata *pb.Metadata
	if service.Metadata != nil {
		switch service.Category {
		case models.AuthCategory:
			var authMeta models.AuthMetadata
			if err := json.Unmarshal(service.Metadata, &authMeta); err == nil {
				metadata = &pb.Metadata{
					Metadata: &pb.Metadata_Auth{
						Auth: authMetadataToProto(&authMeta),
					},
				}
			}
		case models.UserCategory:
			var userMeta models.UserMetadata
			if err := json.Unmarshal(service.Metadata, &userMeta); err == nil {
				metadata = &pb.Metadata{
					Metadata: &pb.Metadata_User{
						User: userMetadataToProto(&userMeta),
					},
				}
			}
		case models.TournamentCategory:
			var tournamentMeta models.TournamentMetadata
			if err := json.Unmarshal(service.Metadata, &tournamentMeta); err == nil {
				metadata = &pb.Metadata{
					Metadata: &pb.Metadata_Tournament{
						Tournament: tournamentMetadataToProto(&tournamentMeta),
					},
				}
			}
		case models.DebateCategory:
			var debateMeta models.DebateMetadata
			if err := json.Unmarshal(service.Metadata, &debateMeta); err == nil {
				metadata = &pb.Metadata{
					Metadata: &pb.Metadata_Debate{
						Debate: debateMetadataToProto(&debateMeta),
					},
				}
			}
		case models.ReportCategory:
			var reportMeta models.ReportMetadata
			if err := json.Unmarshal(service.Metadata, &reportMeta); err == nil {
				metadata = &pb.Metadata{
					Metadata: &pb.Metadata_Report{
						Report: reportMetadataToProto(&reportMeta),
					},
				}
			}
		}
	}

	return &pb.Notification{
		Id:              service.ID,
		Category:        pb.Category(pb.Category_value[string(service.Category)]),
		Type:            pb.NotificationType(pb.NotificationType_value[string(service.Type)]),
		UserId:          service.UserID,
		UserRole:        pb.UserRole(pb.UserRole_value[string(service.UserRole)]),
		Title:           service.Title,
		Content:         service.Content,
		DeliveryMethods: deliveryMethods,
		Priority:        pb.Priority(pb.Priority_value[string(service.Priority)]),
		Actions:         actions,
		Metadata:        metadata,
		Status:          pb.Status(pb.Status_value[string(service.Status)]),
		IsRead:          service.IsRead,
		ReadAt:          readAt,
		CreatedAt:       service.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       service.UpdatedAt.Format(time.RFC3339),
		ExpiresAt:       service.ExpiresAt.Format(time.RFC3339),
	}, nil
}
