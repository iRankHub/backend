package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/iRankHub/backend/internal/models"
	emails "github.com/iRankHub/backend/internal/utils/emails"

)

type ReminderService struct {
	db   *sql.DB
	cron *cron.Cron
}

func NewReminderService(db *sql.DB) *ReminderService {
	c := cron.New()
	return &ReminderService{db: db, cron: c}
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

	invitations, err := queries.GetPendingInvitations(ctx)
	if err != nil {
		log.Printf("Failed to get pending invitations: %v\n", err)
		return
	}

	if err := emails.SendReminderEmails(ctx, invitations, queries); err != nil {
		log.Printf("Failed to send reminder emails: %v\n", err)
	}
}


func (s *ReminderService) getRecipientInfo(ctx context.Context, queries *models.Queries, invitation models.Tournamentinvitation) (string, string, error) {
	if invitation.Schoolid.Valid {
		school, err := queries.GetSchoolByID(ctx, invitation.Schoolid.Int32)
		if err != nil {
			return "", "", fmt.Errorf("failed to get school details: %v", err)
		}
		return school.Contactemail, "school", nil
	} else if invitation.Volunteerid.Valid {
		volunteer, err := queries.GetVolunteerByID(ctx, invitation.Volunteerid.Int32)
		if err != nil {
			return "", "", fmt.Errorf("failed to get volunteer details: %v", err)
		}
		user, err := queries.GetUserByID(ctx, volunteer.Userid)
		if err != nil {
			return "", "", fmt.Errorf("failed to get user details: %v", err)
		}
		return user.Email, "volunteer", nil
	}
	return "", "", fmt.Errorf("invalid invitation: neither school nor volunteer ID is set")
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

func (s *ReminderService) getShouldSendReminder(timeUntilTournament time.Duration, invitationStatus string, isSchool bool) string {
	days := int(timeUntilTournament.Hours() / 24)
	hours := int(timeUntilTournament.Hours()) % 24

	if isSchool {
		if days == 3 && hours == 0 && invitationStatus == "pending" {
			return "school_accept"
		}
		if days == 2 && hours == 7 && invitationStatus == "accepted" {
			return "school_revoke"
		}
	} else {
		if days == 2 && hours == 0 && invitationStatus == "pending" {
			return "volunteer_accept"
		}
		if days == 2 && hours == 0 && invitationStatus == "accepted" {
			return "volunteer_revoke"
		}
	}

	return "none"
}