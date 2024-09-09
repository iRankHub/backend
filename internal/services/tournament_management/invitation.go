package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
	notification "github.com/iRankHub/backend/internal/utils/notifications"
	notifications "github.com/iRankHub/backend/internal/services/notification"
)

type InvitationService struct {
	db *sql.DB
}

func NewInvitationService(db *sql.DB) *InvitationService {
	return &InvitationService{db: db}
}

func (s *InvitationService) GetInvitationsByUser(ctx context.Context, req *tournament_management.GetInvitationsByUserRequest) (*tournament_management.GetInvitationsByUserResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}
	userID := int32(userIDFloat)

	queries := models.New(s.db)
	invitations, err := queries.GetInvitationsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitations: %v", err)
	}

	invitationInfos := make([]*tournament_management.InvitationInfo, len(invitations))
	for i, inv := range invitations {
		invitationInfos[i] = &tournament_management.InvitationInfo{
			InvitationId: inv.Invitationid,
			Status:       inv.Status,
			IdebateId:    inv.Inviteeid,
			InviteeRole:  inv.Inviteerole,
			CreatedAt:    inv.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:    inv.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &tournament_management.GetInvitationsByUserResponse{
		Invitations: invitationInfos,
	}, nil
}

func (s *InvitationService) GetInvitationsByTournament(ctx context.Context, req *tournament_management.GetInvitationsByTournamentRequest) (*tournament_management.GetInvitationsByTournamentResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	invitations, err := queries.GetInvitationsByTournament(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get invitations: %v", err)
	}

	invitationInfos := make([]*tournament_management.InvitationInfo, len(invitations))
	for i, inv := range invitations {
		invitationInfos[i] = &tournament_management.InvitationInfo{
			InvitationId: inv.Invitationid,
			Status:       inv.Status,
			IdebateId:    inv.Inviteeid,
			InviteeName:  inv.Inviteename.(string),
			InviteeRole:  inv.Inviteerole,
			CreatedAt:    inv.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:    inv.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return &tournament_management.GetInvitationsByTournamentResponse{
		Invitations: invitationInfos,
	}, nil
}

func (s *InvitationService) UpdateInvitationStatus(ctx context.Context, req *tournament_management.UpdateInvitationStatusRequest) (*tournament_management.UpdateInvitationStatusResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	updatedInvitation, err := queries.UpdateInvitationStatus(ctx, models.UpdateInvitationStatusParams{
		Invitationid: req.GetInvitationId(),
		Status:       req.GetNewStatus(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update invitation status: %v", err)
	}

	return &tournament_management.UpdateInvitationStatusResponse{
		Success: true,
		Message: fmt.Sprintf("Invitation status updated to %s", updatedInvitation.Status),
	}, nil
}

func (s *InvitationService) BulkUpdateInvitationStatus(ctx context.Context, req *tournament_management.BulkUpdateInvitationStatusRequest) (*tournament_management.BulkUpdateInvitationStatusResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	updatedInvitations, err := queries.BulkUpdateInvitationStatus(ctx, models.BulkUpdateInvitationStatusParams{
		Column1: req.GetInvitationIds(),
		Status:  req.GetNewStatus(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to bulk update invitation statuses: %v", err)
	}

	updatedIDs := make([]int32, len(updatedInvitations))
	for i, inv := range updatedInvitations {
		updatedIDs[i] = inv.Invitationid
	}

	return &tournament_management.BulkUpdateInvitationStatusResponse{
		Success:              true,
		Message:              fmt.Sprintf("%d invitations updated to status %s", len(updatedInvitations), req.GetNewStatus()),
		UpdatedInvitationIds: updatedIDs,
	}, nil
}

func (s *InvitationService) ResendInvitation(ctx context.Context, req *tournament_management.ResendInvitationRequest) (*tournament_management.ResendInvitationResponse, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	invitation, err := queries.GetInvitationByID(ctx, req.GetInvitationId())
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation: %v", err)
	}

	_, err = queries.UpdateReminderSentAt(ctx, models.UpdateReminderSentAtParams{
		Invitationid:   invitation.Invitationid,
		Remindersentat: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update reminder sent timestamp: %v", err)
	}

	// Send email asynchronously
	go s.sendInvitationEmailAsync(invitation.Invitationid)

	return &tournament_management.ResendInvitationResponse{
		Success: true,
		Message: "Invitation resend process started",
	}, nil
}

func (s *InvitationService) BulkResendInvitations(ctx context.Context, req *tournament_management.BulkResendInvitationsRequest) (*tournament_management.BulkResendInvitationsResponse, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	queries := models.New(s.db)

	for _, invitationID := range req.GetInvitationIds() {
		_, err := queries.GetInvitationByID(ctx, invitationID)
		if err != nil {
			log.Printf("Failed to get invitation %d: %v", invitationID, err)
			continue
		}

		_, err = queries.UpdateReminderSentAt(ctx, models.UpdateReminderSentAtParams{
			Invitationid:   invitationID,
			Remindersentat: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			log.Printf("Failed to update reminder sent timestamp for invitation %d: %v", invitationID, err)
			continue
		}

		// Send email asynchronously
		go s.sendInvitationEmailAsync(invitationID)
	}
	return &tournament_management.BulkResendInvitationsResponse{
		Success: true,
		Message: fmt.Sprintf("Resend process started for %d invitations", len(req.GetInvitationIds())),
	}, nil
}

func (s *InvitationService) sendInvitationEmailAsync(invitationID int32) {
	ctx := context.Background()
	queries := models.New(s.db)

	invitation, err := queries.GetInvitationByID(ctx, invitationID)
	if err != nil {
		log.Printf("Failed to get invitation for ID %d: %v", invitationID, err)
		return
	}

	tournament, err := queries.GetTournamentByID(ctx, invitation.Tournamentid)
	if err != nil {
		log.Printf("Failed to get tournament for invitation %d: %v", invitationID, err)
		return
	}

	league, err := queries.GetLeagueByID(ctx, tournament.Leagueid.Int32)
	if err != nil {
		log.Printf("Failed to get league for tournament %d: %v", tournament.Tournamentid, err)
		return
	}

	format, err := queries.GetTournamentFormatByID(ctx, tournament.Formatid)
	if err != nil {
		log.Printf("Failed to get tournament format for tournament %d: %v", tournament.Tournamentid, err)
		return
	}

	tournamentModel := models.Tournament{
		Tournamentid:               tournament.Tournamentid,
		Name:                       tournament.Name,
		Startdate:                  tournament.Startdate,
		Enddate:                    tournament.Enddate,
		Location:                   tournament.Location,
		Formatid:                   tournament.Formatid,
		Leagueid:                   tournament.Leagueid,
		Coordinatorid:              tournament.Coordinatorid,
		Numberofpreliminaryrounds:  tournament.Numberofpreliminaryrounds,
		Numberofeliminationrounds:  tournament.Numberofeliminationrounds,
		Judgesperdebatepreliminary: tournament.Judgesperdebatepreliminary,
		Judgesperdebateelimination: tournament.Judgesperdebateelimination,
		Tournamentfee:              tournament.Tournamentfee,
	}

	var subject, content string
	var email string

	switch invitation.Inviteerole {
	case "student":
		student, err := queries.GetStudentByIDebateID(ctx, sql.NullString{String: invitation.Inviteeid, Valid: invitation.Inviteeid != ""})
		if err != nil {
			log.Printf("Failed to get student details: %v", err)
			return
		}
		subject = fmt.Sprintf("Reminder: Invitation to Participate in %s Tournament", tournament.Name)
		content = notification.PrepareStudentInvitationContent(student, tournamentModel, league, format)
		email = student.Email.String
	case "school":
		school, err := queries.GetSchoolByIDebateID(ctx, sql.NullString{String: invitation.Inviteeid, Valid: invitation.Inviteeid != ""})
		if err != nil {
			log.Printf("Failed to get school details: %v", err)
			return
		}
		subject = fmt.Sprintf("Reminder: Invitation to %s Tournament", tournament.Name)
		content = notification.PrepareSchoolInvitationContent(school, tournamentModel, league, format)
		email = school.Contactemail
	case "volunteer":
		volunteer, err := queries.GetVolunteerByIDebateID(ctx, sql.NullString{String: invitation.Inviteeid, Valid: invitation.Inviteeid != ""})
		if err != nil {
			log.Printf("Failed to get volunteer details: %v", err)
			return
		}
		user, err := queries.GetUserByID(ctx, volunteer.Userid)
		if err != nil {
			log.Printf("Failed to get user details: %v", err)
			return
		}
		subject = fmt.Sprintf("Reminder: Invitation to Judge at %s Tournament", tournament.Name)
		content = notification.PrepareVolunteerInvitationContent(volunteer, tournamentModel, league, format)
		email = user.Email
	default:
		log.Printf("Unknown invitee role: %s", invitation.Inviteerole)
		return
	}

	err = notification.SendNotification(notifications.EmailNotification, email, subject, content)
	if err != nil {
		log.Printf("Failed to send invitation email to %s for invitation ID %d: %v", email, invitationID, err)
	} else {
		log.Printf("Successfully sent invitation email to %s for invitation ID %d", email, invitationID)
	}
}

func (s *InvitationService) validateAdminRole(token string) (map[string]interface{}, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can perform this action")
	}

	return claims, nil
}

func (s *InvitationService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}
