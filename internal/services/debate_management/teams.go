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

	// Validate token and get claims
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	// Extract user role and ID from claims
	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token: user_role not found")
	}

	// Check if user has permission (admin or school)
	if userRole != "admin" && userRole != "school" {
		log.Printf("Unauthorized: user role %s cannot create teams", userRole)
		return nil, fmt.Errorf("unauthorized: only admins and schools can create teams")
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
	// Validate token and get claims
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token: user_role not found")
	}

	// Check if user has permission (admin or school)
	if userRole != "admin" && userRole != "school" {
		log.Printf("Unauthorized: user role %s cannot view teams", userRole)
		return nil, fmt.Errorf("unauthorized: only admins and schools can view teams")
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
	// Validate token and get claims
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token: user_role not found")
	}

	// Check if user has permission (admin or school)
	if userRole != "admin" && userRole != "school" {
		return nil, fmt.Errorf("unauthorized: only admins and schools can update teams")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	// Check if team is part of any debate
	hasDebates, err := queries.TeamHasDebates(ctx, req.GetTeam().GetTeamId())
	if err != nil {
		return nil, fmt.Errorf("failed to check if team has debates: %v", err)
	}

	// Get current team members before updating
	currentMembers, err := queries.GetTeamMembers(ctx, req.GetTeam().GetTeamId())
	if err != nil {
		return nil, fmt.Errorf("failed to get current team members: %v", err)
	}

	// Create maps of current and new members for comparison
	currentMembersMap := make(map[int32]bool)
	for _, member := range currentMembers {
		currentMembersMap[member.Studentid] = true
	}

	newMembersMap := make(map[int32]bool)
	for _, speaker := range req.GetTeam().GetSpeakers() {
		newMembersMap[speaker.GetSpeakerId()] = true
	}

	// Check if team composition is changing and team has debates
	membersChanged := false
	if len(currentMembers) != len(req.GetTeam().GetSpeakers()) {
		membersChanged = true
	} else {
		for _, speaker := range req.GetTeam().GetSpeakers() {
			if !currentMembersMap[speaker.GetSpeakerId()] {
				membersChanged = true
				break
			}
		}
	}

	if hasDebates && membersChanged {
		if userRole != "admin" {
			return nil, fmt.Errorf("cannot modify team members: team is already part of debates - contact an admin")
		}

		// If user is admin, we'll allow the change but need to update speaker scores

		// Get all debates the team is part of
		debates, err := queries.GetDebatesByTeam(ctx, req.GetTeam().GetTeamId())
		if err != nil {
			return nil, fmt.Errorf("failed to get debates for team: %v", err)
		}

		// For each debate, update the speaker scores
		for _, debate := range debates {
			// Get ballots for this debate
			ballots, err := queries.GetBallotsByDebateID(ctx, debate.Debateid)
			if err != nil {
				return nil, fmt.Errorf("failed to get ballots for debate: %v", err)
			}

			for _, ballot := range ballots {
				// Delete existing speaker scores for this team in this ballot
				err = queries.DeleteSpeakerScoresByBallotAndTeam(ctx, models.DeleteSpeakerScoresByBallotAndTeamParams{
					Ballotid: ballot.Ballotid,
					Teamid:   req.GetTeam().GetTeamId(),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to delete speaker scores: %v", err)
				}

				// We'll recreate the speaker scores after updating team members
			}
		}
	}

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

	// If team members changed and team has debates, recreate speaker scores
	if hasDebates && membersChanged {
		// Get all debates the team is part of (again, since we need the updated info)
		debates, err := queries.GetDebatesByTeam(ctx, req.GetTeam().GetTeamId())
		if err != nil {
			return nil, fmt.Errorf("failed to get debates for team: %v", err)
		}

		// For each debate, create new speaker scores for the updated team members
		for _, debate := range debates {
			err = queries.CreateSpeakerScoresForTeam(ctx, models.CreateSpeakerScoresForTeamParams{
				Debateid: debate.Debateid,
				Teamid:   req.GetTeam().GetTeamId(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create new speaker scores: %v", err)
			}
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
	// Validate token and get claims
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return false, "", fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return false, "", fmt.Errorf("invalid token: user_role not found")
	}

	// Check if user has permission (admin or school)
	if userRole != "admin" && userRole != "school" {
		return false, "", fmt.Errorf("unauthorized: only admins and schools can delete teams")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	// Check if team is part of any debate
	hasDebates, err := queries.TeamHasDebates(ctx, req.GetTeamId())
	if err != nil {
		return false, "", fmt.Errorf("failed to check if team has debates: %v", err)
	}

	if hasDebates {
		return false, "Team cannot be deleted because it is part of one or more debates", nil
	}

	// Delete team members
	err = queries.DeleteTeamMembers(ctx, req.GetTeamId())
	if err != nil {
		return false, "", fmt.Errorf("failed to delete team members: %v", err)
	}

	// Delete the team
	err = queries.DeleteTeam(ctx, req.GetTeamId())
	if err != nil {
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
			Name:      dbSpeaker.Firstname + " " + dbSpeaker.Lastname,
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
