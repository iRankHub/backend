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
	emails "github.com/iRankHub/backend/internal/utils/emails"
)

type InvitationService struct {
	db *sql.DB
}

func NewInvitationService(db *sql.DB) *InvitationService {
	return &InvitationService{db: db}
}

func (s *InvitationService) AcceptInvitation(ctx context.Context, req *tournament_management.AcceptInvitationRequest) (*tournament_management.AcceptInvitationResponse, error) {
	log.Printf("AcceptInvitation called with invitation ID: %d", req.GetInvitationId())

	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return nil, fmt.Errorf("authentication failed: %v", err)
	}
	log.Printf("Token validated successfully")

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}
	userID := sql.NullInt32{Int32: int32(userIDFloat), Valid: true}
	log.Printf("User ID from token: %d", userID.Int32)

	err = queries.UpdateInvitationStatusWithUserCheck(ctx, models.UpdateInvitationStatusWithUserCheckParams{
		Invitationid: req.GetInvitationId(),
		Status:       "accepted",
		Userid:       userID,
	})
	if err != nil {
		log.Printf("Failed to update invitation status: %v", err)
		return nil, fmt.Errorf("failed to accept invitation: %v", err)
	}
	log.Printf("Invitation status updated successfully")

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}
	log.Printf("Transaction committed successfully")

	log.Printf("Invitation %d accepted successfully", req.GetInvitationId())
	return &tournament_management.AcceptInvitationResponse{
		Success: true,
		Message: "Invitation accepted successfully",
	}, nil
}

func (s *InvitationService) DeclineInvitation(ctx context.Context, req *tournament_management.DeclineInvitationRequest) (*tournament_management.DeclineInvitationResponse, error) {
	log.Printf("DeclineInvitation called with invitation ID: %d", req.GetInvitationId())

	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return nil, fmt.Errorf("authentication failed: %v", err)
	}
	log.Printf("Token validated successfully")

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}
	userID := sql.NullInt32{Int32: int32(userIDFloat), Valid: true}

	err = queries.UpdateInvitationStatusWithUserCheck(ctx, models.UpdateInvitationStatusWithUserCheckParams{
		Invitationid: req.GetInvitationId(),
		Status:       "declined",
		Userid:       userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to decline invitation: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Invitation %d declined successfully", req.GetInvitationId())
	return &tournament_management.DeclineInvitationResponse{
		Success: true,
		Message: "Invitation declined successfully",
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

	tournament, err := queries.GetTournamentByID(ctx, invitation.Tournamentid)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament: %v", err)
	}

	league, err := queries.GetLeagueByID(ctx, tournament.Leagueid.Int32)
	if err != nil {
		return nil, fmt.Errorf("failed to get league: %v", err)
	}

	format, err := queries.GetTournamentFormatByID(ctx, tournament.Formatid)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament format: %v", err)
	}

	var email string
	var subject string
	var content string

	switch {
	case invitation.Schoolid.Valid:
		school, err := queries.GetSchoolByID(ctx, invitation.Schoolid.Int32)
		if err != nil {
			return nil, fmt.Errorf("failed to get school: %v", err)
		}
		email = school.Contactemail
		subject = fmt.Sprintf("Reminder: Invitation to %s Tournament", tournament.Name)
		content = emails.PrepareSchoolInvitationContent(school, convertToTournament(tournament), league, format)
	case invitation.Volunteerid.Valid:
		volunteer, err := queries.GetVolunteerByID(ctx, invitation.Volunteerid.Int32)
		if err != nil {
			return nil, fmt.Errorf("failed to get volunteer: %v", err)
		}
		user, err := queries.GetUserByID(ctx, volunteer.Userid)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %v", err)
		}
		email = user.Email
		subject = fmt.Sprintf("Reminder: Invitation to Judge at %s Tournament", tournament.Name)
		content = emails.PrepareVolunteerInvitationContent(volunteer, convertToTournament(tournament), league, format)
	case invitation.Studentid.Valid:
		student, err := queries.GetStudentByID(ctx, invitation.Studentid.Int32)
		if err != nil {
			return nil, fmt.Errorf("failed to get student: %v", err)
		}
		email = student.Email.String
		subject = fmt.Sprintf("Reminder: Invitation to Participate in %s Tournament", tournament.Name)
		content = emails.PrepareStudentInvitationContent(student, convertToTournament(tournament), league, format)
	default:
		return nil, fmt.Errorf("invalid invitation type")
	}

	err = emails.SendEmail(email, subject, content)
	if err != nil {
		return nil, fmt.Errorf("failed to send invitation email: %v", err)
	}

	_, err = queries.UpdateReminderSentAt(ctx, models.UpdateReminderSentAtParams{
		Invitationid:   invitation.Invitationid,
		Remindersentat: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update reminder sent timestamp: %v", err)
	}

	return &tournament_management.ResendInvitationResponse{
		Success: true,
		Message: "Invitation resent successfully",
	}, nil
}

func (s *InvitationService) GetInvitationStatus(ctx context.Context, req *tournament_management.GetInvitationStatusRequest) (*tournament_management.GetInvitationStatusResponse, error) {
	log.Printf("GetInvitationStatus called with invitation ID: %d", req.GetInvitationId())

	if err := s.validateAuthentication(req.GetToken()); err != nil {
		log.Printf("Authentication failed: %v", err)
		return nil, err
	}
	log.Printf("Token validated successfully")

	queries := models.New(s.db)
	status, err := queries.GetInvitationStatus(ctx, req.GetInvitationId())
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation status: %v", err)
	}

	teams, err := queries.GetTeamsByInvitation(ctx, sql.NullInt32{Int32: req.GetInvitationId(), Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get teams for invitation: %v", err)
	}

	var registeredTeams []*tournament_management.Team
	for _, team := range teams {
		registeredTeams = append(registeredTeams, &tournament_management.Team{
			TeamId:           team.Teamid,
			TeamName:         team.Name,
			NumberOfSpeakers: int32(team.NumberOfSpeakers),
		})
	}

	log.Printf("Invitation status retrieved successfully for invitation ID: %d", req.GetInvitationId())
	return &tournament_management.GetInvitationStatusResponse{
		Status:          status.Status,
		RegisteredTeams: registeredTeams,
	}, nil
}

func (s *InvitationService) RegisterTeam(ctx context.Context, req *tournament_management.RegisterTeamRequest) (*tournament_management.RegisterTeamResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)

	// Get the invitation details
	invitation, err := queries.GetInvitationStatus(ctx, req.GetInvitationId())
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation details: %v", err)
	}

	// Register the team
	team, err := queries.RegisterTeam(ctx, models.RegisterTeamParams{
		Name:         req.GetTeamName(),
		Schoolid:     invitation.Schoolid.Int32,
		Tournamentid: invitation.Tournamentid,
		Invitationid: sql.NullInt32{Int32: req.GetInvitationId(), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register team: %v", err)
	}

	return &tournament_management.RegisterTeamResponse{
		Success: true,
		Message: "Team registered successfully",
		TeamId:  team.Teamid,
	}, nil
}

func (s *InvitationService) AddTeamMember(ctx context.Context, req *tournament_management.AddTeamMemberRequest) (*tournament_management.AddTeamMemberResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	err := queries.AddTeamMember(ctx, models.AddTeamMemberParams{
		Teamid:    req.GetTeamId(),
		Studentid: req.GetStudentId(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add team member: %v", err)
	}

	return &tournament_management.AddTeamMemberResponse{
		Success: true,
		Message: "Team member added successfully",
	}, nil
}

// Helper function to convert GetTournamentByIDRow to Tournament
func convertToTournament(row models.GetTournamentByIDRow) models.Tournament {
	return models.Tournament{
		Tournamentid:               row.Tournamentid,
		Name:                       row.Name,
		Startdate:                  row.Startdate,
		Enddate:                    row.Enddate,
		Location:                   row.Location,
		Formatid:                   row.Formatid,
		Leagueid:                   row.Leagueid,
		Coordinatorid:              row.Coordinatorid,
		Numberofpreliminaryrounds:  row.Numberofpreliminaryrounds,
		Numberofeliminationrounds:  row.Numberofeliminationrounds,
		Judgesperdebatepreliminary: row.Judgesperdebatepreliminary,
		Judgesperdebateelimination: row.Judgesperdebateelimination,
		Tournamentfee:              row.Tournamentfee,
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