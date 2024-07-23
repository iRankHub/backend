package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"

)

type InvitationService struct {
	db *sql.DB
}

func NewInvitationService(db *sql.DB) *InvitationService {
	return &InvitationService{db: db}
}

func (s *InvitationService) AcceptInvitation(ctx context.Context, req *tournament_management.AcceptInvitationRequest) (*tournament_management.AcceptInvitationResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	err := queries.UpdateInvitationStatus(ctx, models.UpdateInvitationStatusParams{
		Invitationid: req.GetInvitationId(),
		Status:       "accepted",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to accept invitation: %v", err)
	}

	return &tournament_management.AcceptInvitationResponse{
		Success: true,
		Message: "Invitation accepted successfully",
	}, nil
}

func (s *InvitationService) DeclineInvitation(ctx context.Context, req *tournament_management.DeclineInvitationRequest) (*tournament_management.DeclineInvitationResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	err := queries.UpdateInvitationStatus(ctx, models.UpdateInvitationStatusParams{
		Invitationid: req.GetInvitationId(),
		Status:       "declined",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to decline invitation: %v", err)
	}

	return &tournament_management.DeclineInvitationResponse{
		Success: true,
		Message: "Invitation declined successfully",
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

func (s *InvitationService) GetInvitationStatus(ctx context.Context, req *tournament_management.GetInvitationStatusRequest) (*tournament_management.GetInvitationStatusResponse, error) {
    if err := s.validateAuthentication(req.GetToken()); err != nil {
        return nil, err
    }

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

	return &tournament_management.GetInvitationStatusResponse{
		Status:          status.Status,
		RegisteredTeams: registeredTeams,
	}, nil
}

func (s *InvitationService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}