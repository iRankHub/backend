package dispatchers

import (
	"context"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/services/notification/models"
)

// DebateDispatcher handles debate-related notifications
type DebateDispatcher struct {
	*BaseDispatcher
	options DispatcherOptions
}

// NewDebateDispatcher creates a new debate dispatcher
func NewDebateDispatcher(base *BaseDispatcher, options DispatcherOptions) *DebateDispatcher {
	return &DebateDispatcher{
		BaseDispatcher: base,
		options:        options,
	}
}

// Dispatch implements the Dispatcher interface
func (d *DebateDispatcher) Dispatch(ctx context.Context, n *models.Notification) error {
	// Set default expiration if not set
	if n.ExpiresAt.IsZero() {
		n.ExpiresAt = time.Now().AddDate(0, 0, d.options.ExpirationDays)
	}

	// Adjust delivery methods based on notification type
	switch n.Type {
	case models.RoundAssignment, models.JudgeAssignment:
		// Assignments need immediate attention
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
			models.PushDelivery,
		}
		n.Priority = models.HighPriority

	case models.BallotSubmission:
		// Ballot submissions are time-sensitive
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
			models.PushDelivery,
		}
		n.Priority = models.UrgentPriority

	case models.DebateResults:
		// Results should be accessible everywhere
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		n.Priority = models.MediumPriority

	case models.RoomChange:
		// Room changes are urgent
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
			models.PushDelivery,
		}
		n.Priority = models.UrgentPriority
	}

	// Store in RabbitMQ first
	if err := d.rabbitmqSender.Publish(ctx, n); err != nil {
		return err
	}

	// Send through each delivery method
	for _, method := range n.DeliveryMethods {
		var err error
		switch method {
		case models.EmailDelivery:
			if d.options.EmailEnabled {
				err = d.emailSender.Send(ctx, n)
			}
		case models.InAppDelivery:
			if d.options.InAppEnabled {
				err = d.inAppSender.Send(ctx, n)
			}
		}

		if err != nil {
			n.UpdateDeliveryStatus(method, models.StatusFailed, err)
		}
	}

	return nil
}

// Helper methods for debate-specific notifications

func (d *DebateDispatcher) SendRoundAssignment(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata) error {
	n := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.RoundAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Round %d Assignment", metadata.RoundNumber),
		Content:  fmt.Sprintf("You have been assigned to a debate in room %s", metadata.Room),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Details",
				URL:   fmt.Sprintf("/debates/%s", metadata.DebateID),
			},
		},
		Priority: models.HighPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *DebateDispatcher) SendJudgeAssignment(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata) error {
	title := "Judge Assignment"
	if metadata.HeadJudge == userID {
		title = "Head Judge Assignment"
	}

	n := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.JudgeAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    title,
		Content:  fmt.Sprintf("You have been assigned as judge for Round %d in room %s", metadata.RoundNumber, metadata.Room),
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
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *DebateDispatcher) SendBallotReminder(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata) error {
	n := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.BallotSubmission,
		UserID:   userID,
		UserRole: role,
		Title:    "Ballot Submission Reminder",
		Content:  fmt.Sprintf("Please submit your ballot for Round %d (Room %s)", metadata.RoundNumber, metadata.Room),
		Actions: []models.Action{
			{
				Type:  models.ActionSubmit,
				Label: "Submit Now",
				URL:   fmt.Sprintf("/debates/%s/ballot", metadata.DebateID),
			},
		},
		Priority: models.UrgentPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *DebateDispatcher) SendDebateResults(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata, winner string, score string) error {
	n := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.DebateResults,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Round %d Results", metadata.RoundNumber),
		Content:  fmt.Sprintf("Results for debate between %s vs %s: %s won (%s)", metadata.Team1, metadata.Team2, winner, score),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Results",
				URL:   fmt.Sprintf("/debates/%s/results", metadata.DebateID),
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *DebateDispatcher) SendRoomChange(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata, oldRoom string) error {
	n := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.RoomChange,
		UserID:   userID,
		UserRole: role,
		Title:    "Room Change Alert",
		Content:  fmt.Sprintf("Your Round %d debate has been moved from Room %s to Room %s", metadata.RoundNumber, oldRoom, metadata.Room),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Updated Details",
				URL:   fmt.Sprintf("/debates/%s", metadata.DebateID),
			},
		},
		Priority: models.UrgentPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *DebateDispatcher) SendBallotSubmissionConfirmation(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata) error {
	n := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.BallotSubmission,
		UserID:   userID,
		UserRole: role,
		Title:    "Ballot Submitted Successfully",
		Content:  fmt.Sprintf("Your ballot for Round %d (Room %s) has been submitted successfully.", metadata.RoundNumber, metadata.Room),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Ballot",
				URL:   fmt.Sprintf("/debates/%s/ballot", metadata.DebateID),
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *DebateDispatcher) SendMotionAnnouncement(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata) error {
	n := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.RoundAssignment,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("Round %d Motion", metadata.RoundNumber),
		Content:  fmt.Sprintf("The motion for your debate is: %s", metadata.Motion),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Details",
				URL:   fmt.Sprintf("/debates/%s", metadata.DebateID),
			},
		},
		Priority: models.HighPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *DebateDispatcher) SendJudgeFeedbackRequest(ctx context.Context, userID string, role models.UserRole, metadata models.DebateMetadata) error {
	n := &models.Notification{
		Category: models.DebateCategory,
		Type:     models.DebateResults,
		UserID:   userID,
		UserRole: role,
		Title:    "Judge Feedback Request",
		Content:  fmt.Sprintf("Please provide feedback for the judges from your Round %d debate.", metadata.RoundNumber),
		Actions: []models.Action{
			{
				Type:  models.ActionSubmit,
				Label: "Provide Feedback",
				URL:   fmt.Sprintf("/debates/%s/feedback", metadata.DebateID),
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}
