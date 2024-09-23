package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/iRankHub/backend/internal/models"
	notifications "github.com/iRankHub/backend/internal/services/notification"
	notification "github.com/iRankHub/backend/internal/utils/notifications"
)

type ReminderService struct {
	db                  *sql.DB
	cron                *cron.Cron
	notificationService *notifications.NotificationService
}

func NewReminderService(db *sql.DB) (*ReminderService, error) {
	c := cron.New()
	notificationService, err := notifications.NewNotificationService(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification service: %v", err)
	}
	return &ReminderService{
		db:                  db,
		cron:                c,
		notificationService: notificationService,
	}, nil
}

func (s *ReminderService) Start() {
	s.cron.AddFunc("0 0 * * *", s.SendReminders) // Run daily at midnight
	s.cron.Start()
}

func (s *ReminderService) Stop() {
	s.cron.Stop()
}

func (s *ReminderService) SendReminders() {
	ctx := context.Background()
	queries := models.New(s.db)

	// Get all active tournaments
	tournaments, err := queries.GetActiveTournaments(ctx)
	if err != nil {
		log.Printf("Failed to get active tournaments: %v\n", err)
		return
	}

	for _, tournament := range tournaments {
		invitations, err := queries.GetPendingInvitations(ctx, tournament.Tournamentid)
		if err != nil {
			log.Printf("Failed to get pending invitations for tournament %d: %v\n", tournament.Tournamentid, err)
			continue
		}

		if err := notification.SendReminderEmails(ctx, s.notificationService, invitations, queries); err != nil {
			log.Printf("Failed to send reminder notification for tournament %d: %v\n", tournament.Tournamentid, err)
		}
	}
}

func (s *ReminderService) getRecipientInfo(ctx context.Context, queries *models.Queries, invitation models.Tournamentinvitation) (string, string, error) {
	switch invitation.Inviteerole {
	case "school":
		school, err := queries.GetSchoolByIDebateID(ctx, sql.NullString{String: invitation.Inviteeid, Valid: invitation.Inviteeid != ""})
		if err != nil {
			return "", "", fmt.Errorf("failed to get school details: %v", err)
		}
		return school.Contactemail, "school", nil
	case "volunteer":
		volunteer, err := queries.GetVolunteerByIDebateID(ctx, sql.NullString{String: invitation.Inviteeid, Valid: invitation.Inviteeid != ""})
		if err != nil {
			return "", "", fmt.Errorf("failed to get volunteer details: %v", err)
		}
		user, err := queries.GetUserByID(ctx, volunteer.Userid)
		if err != nil {
			return "", "", fmt.Errorf("failed to get user details: %v", err)
		}
		return user.Email, "volunteer", nil
	case "student":
		student, err := queries.GetStudentByIDebateID(ctx, sql.NullString{String: invitation.Inviteeid, Valid: invitation.Inviteeid != ""})
		if err != nil {
			return "", "", fmt.Errorf("failed to get student details: %v", err)
		}
		return student.Email.String, "student", nil
	default:
		return "", "", fmt.Errorf("invalid invitation role: %s", invitation.Inviteerole)
	}
}

func (s *ReminderService) updateReminderSentTimestamp(ctx context.Context, queries *models.Queries, invitationID int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	if _, err := qtx.UpdateReminderSentAt(ctx, models.UpdateReminderSentAtParams{
		Invitationid:   invitationID,
		Remindersentat: sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		return fmt.Errorf("failed to update reminder sent timestamp: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *ReminderService) getShouldSendReminder(timeUntilTournament time.Duration, invitationStatus string, inviteeRole string) string {
	days := int(timeUntilTournament.Hours() / 24)
	hours := int(timeUntilTournament.Hours()) % 24

	if inviteeRole == "school" {
		if days == 3 && hours == 0 && invitationStatus == "pending" {
			return "school_accept"
		}
		if days == 2 && hours == 7 && invitationStatus == "accepted" {
			return "school_revoke"
		}
	} else if inviteeRole == "volunteer" || inviteeRole == "student" {
		if days == 2 && hours == 0 && invitationStatus == "pending" {
			return "accept"
		}
		if days == 2 && hours == 0 && invitationStatus == "accepted" {
			return "revoke"
		}
	}

	return "none"
}
