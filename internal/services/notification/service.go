package notification

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	_ "time"

	"github.com/iRankHub/backend/internal/services/notification/dispatchers"
	"github.com/iRankHub/backend/internal/services/notification/models"
	"github.com/iRankHub/backend/internal/services/notification/senders"
	"github.com/iRankHub/backend/internal/services/notification/storage"
)

type Service struct {
	storage     *storage.CombinedStorage
	dispatchers *dispatchers.DispatcherFactory
	subscribers map[string]chan *models.Notification
	mu          sync.RWMutex
}

type ServiceConfig struct {
	RabbitMQURL    string
	EmailConfig    senders.EmailSenderConfig
	DispatcherOpts dispatchers.DispatcherOptions
}

func NewNotificationService(ctx context.Context, config ServiceConfig, dbConn *sql.DB) (*Service, error) {
	// Initialize email sender
	emailSender, err := senders.NewEmailSender(config.EmailConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create email sender: %w", err)
	}

	// Initialize in-app sender
	inAppSender := senders.NewInAppSender()

	// Initialize RabbitMQ sender
	rabbitmqConfig := senders.RabbitMQConfig{
		URL:          config.RabbitMQURL,
		Exchange:     "notifications",
		ExchangeType: "topic",
		TTL:          30 * 24 * 60 * 60 * 1000, // 30 days in milliseconds
	}

	rabbitmqSender, err := senders.NewRabbitMQSender(rabbitmqConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ sender: %w", err)
	}

	// Create base dispatcher
	baseDispatcher := dispatchers.NewBaseDispatcher(emailSender, inAppSender, rabbitmqSender)

	// Initialize dispatcher factory
	dispatcherFactory := dispatchers.NewDispatcherFactory(baseDispatcher, config.DispatcherOpts)

	// Initialize storage
	metadataStorage := storage.NewMetadataStorage(dbConn)
	rabbitmqStorage, _ := storage.NewRabbitMQStorage(config.RabbitMQURL, metadataStorage)
	combinedStorage := storage.NewCombinedStorage(rabbitmqStorage, metadataStorage)

	return &Service{
		storage:     combinedStorage,
		dispatchers: dispatcherFactory,
		subscribers: make(map[string]chan *models.Notification),
	}, nil
}

// SendNotification sends a notification using the appropriate dispatcher
func (s *Service) SendNotification(ctx context.Context, notification *models.Notification) error {
	dispatcher, err := s.dispatchers.GetDispatcher(notification.Category)
	if err != nil {
		return fmt.Errorf("failed to get dispatcher: %w", err)
	}

	// Store notification before dispatching
	if err := s.storage.Store(ctx, notification); err != nil {
		return fmt.Errorf("failed to store notification: %w", err)
	}

	// Dispatch notification
	if err := dispatcher.Dispatch(ctx, notification); err != nil {
		// Update storage with failure status
		_ = s.storage.UpdateDeliveryStatus(ctx, notification.ID, models.EmailDelivery, models.StatusFailed)
		return fmt.Errorf("failed to dispatch notification: %w", err)
	}

	// Notify subscribers
	s.notifySubscribers(notification)

	return nil
}

// GetUnreadNotifications retrieves unread notifications for a user
func (s *Service) GetUnreadNotifications(ctx context.Context, userID string) ([]*models.Notification, error) {
	return s.storage.GetUnread(ctx, userID)
}

// MarkAsRead marks notifications as read
func (s *Service) MarkAsRead(ctx context.Context, userID string, notificationIDs []string) error {
	for _, id := range notificationIDs {
		if err := s.storage.MarkAsRead(ctx, id); err != nil {
			return fmt.Errorf("failed to mark notification %s as read: %w", id, err)
		}
	}
	return nil
}

// Subscribe allows a client to subscribe to notifications
func (s *Service) Subscribe(ctx context.Context, userID string) (<-chan *models.Notification, func(), error) {
	s.mu.Lock()
	ch := make(chan *models.Notification, 100)
	s.subscribers[userID] = ch
	s.mu.Unlock()

	// Create unsubscribe function
	unsubscribe := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if ch, ok := s.subscribers[userID]; ok {
			close(ch)
			delete(s.subscribers, userID)
		}
	}

	// Start subscription cleanup when context is done
	go func() {
		<-ctx.Done()
		unsubscribe()
	}()

	return ch, unsubscribe, nil
}

// notifySubscribers sends notification to all subscribed clients
func (s *Service) notifySubscribers(notification *models.Notification) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if ch, ok := s.subscribers[notification.UserID]; ok {
		select {
		case ch <- notification:
		default:
			// Channel is full, skip notification
		}
	}
}

// Helper methods for common notification scenarios

func (s *Service) SendDebateAssignment(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata) error {
	// Marshal metadata to json.RawMessage
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal debate metadata: %w", err)
	}

	notification := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.RoundAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Round %d Assignment", metadata.RoundNumber),
		Content:  fmt.Sprintf("You have been assigned to debate in room %s", metadata.Room),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Details",
				URL:   fmt.Sprintf("/debates/%s", metadata.DebateID),
			},
		},
		Priority: models.HighPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendAccountCreation(ctx context.Context, userID string, role models.UserRole, metadata models.AuthMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal debate metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.AuthCategory,
		Type:     models.AccountCreation,
		UserID:   userID,
		UserRole: role,
		Title:    "Welcome to iRankHub",
		Content:  "Your account has been created successfully. Please complete your profile to get started.",
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Complete Profile",
				URL:   "/profile/edit",
			},
		},
		Priority: models.MediumPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendPasswordReset(ctx context.Context, userID string, role models.UserRole, resetToken string, metadata models.AuthMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal password metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.AuthCategory,
		Type:     models.PasswordReset,
		UserID:   userID,
		UserRole: role,
		Title:    "Password Reset Request",
		Content:  "A password reset has been requested for your account. Click the link below to reset your password.",
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Reset Password",
				URL:   fmt.Sprintf("/auth/reset-password?token=%s", resetToken),
			},
		},
		Priority: models.HighPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendSecurityAlert(ctx context.Context, userID string, role models.UserRole, alertMessage string, metadata models.AuthMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal security metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.AuthCategory,
		Type:     models.SecurityAlert,
		UserID:   userID,
		UserRole: role,
		Title:    "Security Alert",
		Content:  alertMessage,
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Review Activity",
				URL:   "/settings/security",
			},
		},
		Priority: models.UrgentPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendTwoFactorCode(ctx context.Context, userID string, role models.UserRole, code string, metadata models.AuthMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.AuthCategory,
		Type:     models.TwoFactorAuth,
		UserID:   userID,
		UserRole: role,
		Title:    "Two-Factor Authentication Code",
		Content:  fmt.Sprintf("Your verification code is: %s", code),
		Priority: models.UrgentPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendAccountApproval(ctx context.Context, userID string, role models.UserRole, metadata models.AuthMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.AuthCategory,
		Type:     models.AccountApproval,
		UserID:   userID,
		UserRole: role,
		Title:    "Account Approved",
		Content:  "Your account has been approved. You can now access all features of iRankHub.",
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Get Started",
				URL:   "/dashboard",
			},
		},
		Priority: models.MediumPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendProfileUpdate(ctx context.Context, userID string, role models.UserRole, metadata models.UserMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	content := "Your profile has been updated with the following changes:\n"
	for field, value := range metadata.Changes {
		content += fmt.Sprintf("- %s: %s\n", field, value)
	}

	notification := &models.Notification{
		Category: models.UserCategory,
		Type:     models.ProfileUpdate,
		UserID:   userID,
		UserRole: role,
		Title:    "Profile Updated",
		Content:  content,
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Profile",
				URL:   "/profile",
			},
		},
		Priority: models.MediumPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendRoleAssignment(ctx context.Context, userID string, role models.UserRole, metadata models.UserMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	notification := &models.Notification{
		Category: models.UserCategory,
		Type:     models.RoleAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    "Role Updated",
		Content:  fmt.Sprintf("Your role has been changed from %s to %s", metadata.PreviousRole, metadata.NewRole),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Permissions",
				URL:   "/profile/permissions",
			},
		},
		Priority: models.HighPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendStatusChange(ctx context.Context, userID string, role models.UserRole, metadata models.UserMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	notification := &models.Notification{
		Category: models.UserCategory,
		Type:     models.StatusChange,
		UserID:   userID,
		UserRole: role,
		Title:    "Account Status Changed",
		Content:  fmt.Sprintf("Your account status has been changed. Reason: %s", metadata.Reason),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Contact Support",
				URL:   "/support",
			},
		},
		Priority: models.UrgentPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendRoleExpiration(ctx context.Context, userID string, role models.UserRole, metadata models.UserMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.UserCategory,
		Type:     models.RoleAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    "Role Expiration Notice",
		Content:  fmt.Sprintf("Your role as %s will expire on %s", metadata.NewRole, metadata.ExpirationDate.Format("January 2, 2006")),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Review Details",
				URL:   "/profile/role",
			},
		},
		Priority: models.HighPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendTournamentInvitation(ctx context.Context, tournamentID string, inviteeID string, inviteeRole models.UserRole, metadata models.TournamentMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	notification := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.TournamentInvite,
		UserID:   inviteeID,
		UserRole: inviteeRole,
		Title:    fmt.Sprintf("Tournament Invitation: %s", metadata.TournamentName),
		Content: fmt.Sprintf("You have been invited to participate in the %s tournament. Tournament dates: %s to %s",
			metadata.TournamentName,
			metadata.StartDate.Format("Jan 2"),
			metadata.EndDate.Format("Jan 2, 2006")),
		Actions: []models.Action{
			{
				Type:  models.ActionAccept,
				Label: "Accept Invitation",
				URL:   fmt.Sprintf("/tournaments/%s/accept", tournamentID),
			},
			{
				Type:  models.ActionReject,
				Label: "Decline",
				URL:   fmt.Sprintf("/tournaments/%s/decline", tournamentID),
			},
			{
				Type:  models.ActionView,
				Label: "View Details",
				URL:   fmt.Sprintf("/tournaments/%s", tournamentID),
			},
		},
		Priority: models.MediumPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendTournamentRegistrationConfirmation(ctx context.Context, userID string, role models.UserRole, metadata models.TournamentMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.TournamentRegistration,
		UserID:   userID,
		UserRole: role,
		Title:    "Tournament Registration Confirmed",
		Content: fmt.Sprintf("Your registration for %s has been confirmed. Tournament dates: %s to %s",
			metadata.TournamentName,
			metadata.StartDate.Format("Jan 2"),
			metadata.EndDate.Format("Jan 2, 2006")),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Registration",
				URL:   fmt.Sprintf("/tournaments/%s/registration", metadata.TournamentID),
			},
			{
				Type:  models.ActionView,
				Label: "View Schedule",
				URL:   fmt.Sprintf("/tournaments/%s/schedule", metadata.TournamentID),
			},
		},
		Priority: models.MediumPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendTournamentPaymentReminder(ctx context.Context, userID string, role models.UserRole, metadata models.TournamentMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.TournamentPayment,
		UserID:   userID,
		UserRole: role,
		Title:    "Tournament Payment Due",
		Content: fmt.Sprintf("Payment of %s %v for %s tournament registration is due. Please complete your payment to confirm your participation.",
			metadata.Currency,
			metadata.Fee,
			metadata.TournamentName),
		Actions: []models.Action{
			{
				Type:  models.ActionSubmit,
				Label: "Complete Payment",
				URL:   fmt.Sprintf("/tournaments/%s/payment", metadata.TournamentID),
			},
		},
		Priority: models.HighPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendTournamentScheduleUpdate(ctx context.Context, userID string, role models.UserRole, metadata models.TournamentMetadata, changes string) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.TournamentSchedule,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Schedule Update: %s", metadata.TournamentName),
		Content:  fmt.Sprintf("The tournament schedule has been updated: %s", changes),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Updated Schedule",
				URL:   fmt.Sprintf("/tournaments/%s/schedule", metadata.TournamentID),
			},
		},
		Priority: models.HighPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendCoordinatorAssignment(ctx context.Context, userID string, role models.UserRole, metadata models.TournamentMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.TournamentCategory,
		Type:     models.CoordinatorAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    "Tournament Coordinator Assignment",
		Content: fmt.Sprintf("You have been assigned as coordinator for %s tournament (%s - %s)",
			metadata.TournamentName,
			metadata.StartDate.Format("Jan 2"),
			metadata.EndDate.Format("Jan 2, 2006")),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Tournament Details",
				URL:   fmt.Sprintf("/tournaments/%s", metadata.TournamentID),
			},
			{
				Type:  models.ActionView,
				Label: "Coordinator Dashboard",
				URL:   fmt.Sprintf("/tournaments/%s/coordinator", metadata.TournamentID),
			},
		},
		Priority: models.MediumPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendRoundAssignment(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.RoundAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Round %d Assignment", metadata.RoundNumber),
		Content: fmt.Sprintf("You have been assigned to debate in room %s.\nTime: %s\nTeams: %s vs %s",
			metadata.Room,
			metadata.StartTime.Format("3:04 PM"),
			metadata.Team1,
			metadata.Team2),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Details",
				URL:   fmt.Sprintf("/debates/%s", metadata.DebateID),
			},
			{
				Type:  models.ActionView,
				Label: "View Motion",
				URL:   fmt.Sprintf("/debates/%s/motion", metadata.DebateID),
			},
		},
		Priority: models.HighPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendJudgeAssignment(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata) error {
	title := "Judge Assignment"
	if metadata.HeadJudge == userID {
		title = "Head Judge Assignment"
	}

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	notification := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.JudgeAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    title,
		Content: fmt.Sprintf("You have been assigned as %s for Round %d in room %s.\nTime: %s\nTeams: %s vs %s",
			title == "Judge",
			metadata.RoundNumber,
			metadata.Room,
			metadata.StartTime.Format("3:04 PM"),
			metadata.Team1,
			metadata.Team2),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Assignment",
				URL:   fmt.Sprintf("/debates/%s", metadata.DebateID),
			},
			{
				Type:  models.ActionSubmit,
				Label: "Submit Ballot",
				URL:   fmt.Sprintf("/debates/%s/ballot", metadata.DebateID),
			},
		},
		Priority: models.HighPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendBallotReminder(ctx context.Context, userID string, metadata models.DebateMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.BallotSubmission,
		UserID:   userID,
		UserRole: models.VolunteerRole,
		Title:    "Ballot Submission Reminder",
		Content: fmt.Sprintf("Please submit your ballot for Round %d (Room %s)\nDebate: %s vs %s",
			metadata.RoundNumber,
			metadata.Room,
			metadata.Team1,
			metadata.Team2),
		Actions: []models.Action{
			{
				Type:  models.ActionSubmit,
				Label: "Submit Ballot",
				URL:   fmt.Sprintf("/debates/%s/ballot", metadata.DebateID),
			},
		},
		Priority: models.UrgentPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendDebateResults(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata, winner string, score string) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.DebateResults,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Round %d Results", metadata.RoundNumber),
		Content: fmt.Sprintf("Results for debate in room %s:\n%s vs %s\nWinner: %s\nScore: %s",
			metadata.Room,
			metadata.Team1,
			metadata.Team2,
			winner,
			score),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Results",
				URL:   fmt.Sprintf("/debates/%s/results", metadata.DebateID),
			},
		},
		Priority: models.MediumPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendRoomChange(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata, oldRoom string) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	notification := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.RoomChange,
		UserID:   userID,
		UserRole: role,
		Title:    "Room Change Alert",
		Content: fmt.Sprintf("Your Round %d debate has been moved from Room %s to Room %s",
			metadata.RoundNumber,
			oldRoom,
			metadata.Room),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Updated Details",
				URL:   fmt.Sprintf("/debates/%s", metadata.DebateID),
			},
		},
		Priority: models.UrgentPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendReportGenerated(ctx context.Context, userID string, role models.UserRole, metadata models.ReportMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal report metadata: %w", err)
	}

	notification := &models.Notification{
		Category: models.ReportCategory,
		Type:     models.ReportGeneration,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("%s Report Ready", metadata.ReportType),
		Content: fmt.Sprintf("Your requested %s report for %s has been generated and is ready for viewing. Report size: %s",
			metadata.ReportType,
			metadata.Period,
			metadata.FileSize),
		Actions: []models.Action{
			{
				Type:  models.ActionDownload,
				Label: "Download Report",
				URL:   metadata.DownloadURL,
			},
			{
				Type:  models.ActionView,
				Label: "View Online",
				URL:   fmt.Sprintf("/reports/%s", metadata.ReportID),
			},
		},
		Priority: models.LowPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendPerformanceReport(ctx context.Context, userID string, role models.UserRole, metadata models.ReportMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal performance report metadata: %w", err)
	}

	content := fmt.Sprintf("Your performance report for %s is now available. ", metadata.Period)
	if len(metadata.KeyMetrics) > 0 {
		content += "\n\nKey Metrics:\n"
		for metric, value := range metadata.KeyMetrics {
			content += fmt.Sprintf("- %s: %v\n", metric, value)
		}
	}

	notification := &models.Notification{
		Category: models.ReportCategory,
		Type:     models.PerformanceReport,
		UserID:   userID,
		UserRole: role,
		Title:    "Performance Report Available",
		Content:  content,
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Report",
				URL:   fmt.Sprintf("/reports/%s", metadata.ReportID),
			},
			{
				Type:  models.ActionDownload,
				Label: "Download Report",
				URL:   metadata.DownloadURL,
			},
		},
		Priority: models.MediumPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendAnalyticsReport(ctx context.Context, userID string, role models.UserRole, metadata models.ReportMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal analytics report metadata: %w", err)
	}

	content := fmt.Sprintf("The analytics report for %s has been generated. ", metadata.Period)
	if len(metadata.Summary) > 0 {
		content += "\n\nHighlights:\n"
		for key, value := range metadata.Summary {
			content += fmt.Sprintf("- %s: %s\n", key, value)
		}
	}

	notification := &models.Notification{
		Category: models.ReportCategory,
		Type:     models.AnalyticsReport,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("%s Analytics Report", metadata.Period),
		Content:  content,
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Analytics",
				URL:   fmt.Sprintf("/reports/%s", metadata.ReportID),
			},
			{
				Type:  models.ActionDownload,
				Label: "Download Full Report",
				URL:   metadata.DownloadURL,
			},
		},
		Priority: models.MediumPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

func (s *Service) SendAuditReport(ctx context.Context, userID string, role models.UserRole, metadata models.ReportMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal audit report metadata: %w", err)
	}

	notification := &models.Notification{
		Category: models.ReportCategory,
		Type:     models.AuditReport,
		UserID:   userID,
		UserRole: role,
		Title:    "Audit Report Ready",
		Content: fmt.Sprintf("An audit report covering %s has been generated for your review. Please review it at your earliest convenience.",
			metadata.Period),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Review Audit",
				URL:   fmt.Sprintf("/reports/%s", metadata.ReportID),
			},
			{
				Type:  models.ActionDownload,
				Label: "Download Report",
				URL:   metadata.DownloadURL,
			},
		},
		Priority: models.HighPriority,
		Metadata: json.RawMessage(metadataBytes),
	}

	return s.SendNotification(ctx, notification)
}

// Close gracefully shuts down the service
func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close all subscriber channels
	for _, ch := range s.subscribers {
		close(ch)
	}
	s.subscribers = make(map[string]chan *models.Notification)

	// Close storage
	return s.storage.Close()
}
