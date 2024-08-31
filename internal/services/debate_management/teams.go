package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/iRankHub/backend/internal/grpc/proto/debate_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type TeamService struct {
	db *sql.DB
}

func NewTeamService(db *sql.DB) *TeamService {
	return &TeamService{db: db}
}

func (s *TeamService) CreateTeam(ctx context.Context, req *debate_management.CreateTeamRequest) (*debate_management.Team, error) {
	log.Printf("CreateTeam called with name: %s, tournamentId: %d, speakers: %v", req.GetName(), req.GetTournamentId(), req.GetSpeakers())

	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		log.Printf("Admin role validation failed: %v", err)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		}
	}()

	queries := models.New(s.db).WithTx(tx)

	// Check if any of the speakers already belong to a team in this tournament
	for _, speaker := range req.GetSpeakers() {
		log.Printf("Checking existing team membership for speaker ID: %d", speaker.GetSpeakerId())
		hasTeam, err := queries.CheckExistingTeamMembership(ctx, models.CheckExistingTeamMembershipParams{
			Tournamentid: req.GetTournamentId(),
			Studentid:    speaker.GetSpeakerId(),
		})
		if err != nil {
			log.Printf("Failed to check existing team membership: %v", err)
			return nil, fmt.Errorf("failed to check existing team membership: %v", err)
		}
		if hasTeam {
			log.Printf("Speaker with ID %d already belongs to a team in this tournament", speaker.GetSpeakerId())
			return nil, fmt.Errorf("speaker with ID %d already belongs to a team in this tournament", speaker.GetSpeakerId())
		}
	}

	log.Printf("Creating team with name: %s, tournamentId: %d", req.GetName(), req.GetTournamentId())
	// Create the team
	team, err := queries.CreateTeam(ctx, models.CreateTeamParams{
		Name:         req.GetName(),
		Tournamentid: req.GetTournamentId(),
	})
	if err != nil {
		log.Printf("Failed to create team: %v", err)
		return nil, fmt.Errorf("failed to create team: %v", err)
	}
	log.Printf("Team created successfully with ID: %d", team.Teamid)

	// Add speakers to the team
	var speakers []*debate_management.Speaker
	for _, speaker := range req.GetSpeakers() {
		log.Printf("Adding speaker with ID %d to team %d", speaker.GetSpeakerId(), team.Teamid)
		_, err := queries.AddTeamMember(ctx, models.AddTeamMemberParams{
			Teamid:    team.Teamid,
			Studentid: speaker.GetSpeakerId(),
		})
		if err != nil {
			log.Printf("Failed to add team member: %v", err)
			return nil, fmt.Errorf("failed to add team member: %v", err)
		}
		speakers = append(speakers, &debate_management.Speaker{
			SpeakerId: speaker.GetSpeakerId(),
			Name:      speaker.GetName(),
		})
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}
	log.Printf("Transaction committed successfully")

	createdTeam := &debate_management.Team{
		TeamId:   team.Teamid,
		Name:     team.Name,
		Speakers: speakers,
	}
	log.Printf("Team created successfully: %+v", createdTeam)

	return createdTeam, nil
}


func (s *TeamService) GetTeam(ctx context.Context, req *debate_management.GetTeamRequest) (*debate_management.Team, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	team, err := queries.GetTeamByID(ctx, req.GetTeamId())
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %v", err)
	}

	speakers, err := queries.GetTeamMembers(ctx, req.GetTeamId())
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %v", err)
	}

	return convertTeam(team, speakers), nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, req *debate_management.UpdateTeamRequest) (*debate_management.Team, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	// Update team name
	err = queries.UpdateTeam(ctx, models.UpdateTeamParams{
		Teamid: req.GetTeam().GetTeamId(),
		Name:   req.GetTeam().GetName(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update team: %v", err)
	}

	// Remove existing team members
	err = queries.RemoveTeamMembers(ctx, req.GetTeam().GetTeamId())
	if err != nil {
		return nil, fmt.Errorf("failed to remove team members: %v", err)
	}

	// Add new team members
    for _, speaker := range req.GetTeam().GetSpeakers() {
        _, err := queries.AddTeamMember(ctx, models.AddTeamMemberParams{
            Teamid:    req.GetTeam().GetTeamId(),
            Studentid: speaker.GetSpeakerId(),
        })
        if err != nil {
            return nil, fmt.Errorf("failed to add team member: %v", err)
        }
    }

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return s.GetTeam(ctx, &debate_management.GetTeamRequest{TeamId: req.GetTeam().GetTeamId(), Token: req.GetToken()})
}

func (s *TeamService) GetTeamsByTournament(ctx context.Context, req *debate_management.GetTeamsByTournamentRequest) ([]*debate_management.Team, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	teams, err := queries.GetTeamsByTournament(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %v", err)
	}

	result := make([]*debate_management.Team, len(teams))
	for i, team := range teams {
		speakers, err := queries.GetTeamMembers(ctx, team.Teamid)
		if err != nil {
			return nil, fmt.Errorf("failed to get team members: %v", err)
		}
		result[i] = convertTeam(team, speakers)
	}

	return result, nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, req *debate_management.DeleteTeamRequest) (bool, string, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return false, "", err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	// Delete team members (if any)
	err = queries.DeleteTeamMembers(ctx, req.GetTeamId())
	if err != nil {
		return false, "", fmt.Errorf("failed to delete team members: %v", err)
	}

	// Try to delete the team
	err = queries.DeleteTeam(ctx, req.GetTeamId())
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "Team cannot be deleted because it is associated with one or more debates", nil
		}
		return false, "", fmt.Errorf("failed to delete team: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return false, "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	return true, "Team deleted successfully", nil
}

func convertTeam(dbTeam interface{}, dbSpeakers []models.GetTeamMembersRow) *debate_management.Team {
    var teamId int32
    var name string
    var leagueName string

    switch t := dbTeam.(type) {
    case models.GetTeamByIDRow:
        teamId = t.Teamid
        name = t.Name
        leagueName = "" // GetTeamByID doesn't return league name, so we leave it empty
    case models.GetTeamsByTournamentRow:
        teamId = t.Teamid
        name = t.Name
        leagueName = t.Leaguename
    default:
        // Handle unexpected type
        return nil
    }
    speakers := make([]*debate_management.Speaker, len(dbSpeakers))
    for i, dbSpeaker := range dbSpeakers {
        speakers[i] = &debate_management.Speaker{
            SpeakerId: dbSpeaker.Studentid,
            Name: dbSpeaker.Firstname + " " + dbSpeaker.Lastname,
        }
    }

    return &debate_management.Team{
        TeamId:     teamId,
        Name:       name,
        Speakers:   speakers,
        LeagueName: leagueName,
    }
}

func (s *TeamService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}

func (s *TeamService) validateAdminRole(token string) (map[string]interface{}, error) {
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
