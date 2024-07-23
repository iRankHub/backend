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

	for _, invitation := range invitations {
		if err := s.processInvitation(ctx, queries, invitation); err != nil {
			log.Printf("Failed to process invitation %d: %v\n", invitation.Invitationid, err)
		}
	}
}

func (s *ReminderService) processInvitation(ctx context.Context, queries *models.Queries, invitation models.TournamentInvitation) error {
	tournament, err := queries.GetTournamentByID(ctx, invitation.Tournamentid)
	if err != nil {
		return fmt.Errorf("failed to get tournament details: %v", err)
	}

	daysUntilTournament := int(tournament.Startdate.Sub(time.Now()).Hours() / 24)

	if !s.shouldSendReminder(daysUntilTournament) {
		return nil
	}

	recipient, recipientType, err := s.getRecipientInfo(ctx, queries, invitation)
	if err != nil {
		return err
	}

	if err := emails.SendReminderEmail(recipient, recipientType, tournament.Name, daysUntilTournament, invitation.Status, invitation.Invitationid); err != nil {
		return fmt.Errorf("failed to send reminder email: %v", err)
	}

	return s.updateReminderSentTimestamp(ctx, queries, invitation.Invitationid)
}

func (s *ReminderService) getRecipientInfo(ctx context.Context, queries *models.Queries, invitation models.TournamentInvitation) (string, string, error) {
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

func (s *ReminderService) shouldSendReminder(daysUntilTournament int) bool {
	reminderDays := []int{180, 150, 120, 90, 60, 30, 20, 15, 10, 5, 3, 0}
	for _, day := range reminderDays {
		if daysUntilTournament == day {
			return true
		}
	}
	return false
}
