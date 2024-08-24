package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type LeagueService struct {
	db *sql.DB
}

func NewLeagueService(db *sql.DB) *LeagueService {
	return &LeagueService{db: db}
}

func (s *LeagueService) CreateLeague(ctx context.Context, req *tournament_management.CreateLeagueRequest) (*tournament_management.League, error) {
	if err := s.validateAdminRole(req.GetToken()); err != nil {
		return nil, err
	}

	// Validate required fields
	if req.GetName() == "" {
		return nil, fmt.Errorf("league name is required")
	}
	if req.GetLeagueType() != tournament_management.LeagueType_local && req.GetLeagueType() != tournament_management.LeagueType_international {
		return nil, fmt.Errorf("invalid league type: must be either LOCAL or INTERNATIONAL")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	details := make(map[string]interface{})
	switch req.GetLeagueDetails().(type) {
	case *tournament_management.CreateLeagueRequest_LocalDetails:
		localDetails := req.GetLocalDetails()
		if len(localDetails.GetProvinces()) == 0 || len(localDetails.GetDistricts()) == 0 {
			return nil, fmt.Errorf("both province and district are required for local leagues")
		}
		details["provinces"] = localDetails.GetProvinces()
		details["districts"] = localDetails.GetDistricts()
	case *tournament_management.CreateLeagueRequest_InternationalDetails:
		internationalDetails := req.GetInternationalDetails()
		if len(internationalDetails.GetContinents()) == 0 || len(internationalDetails.GetCountries()) == 0 {
			return nil, fmt.Errorf("both continent and country are required for international leagues")
		}
		details["continents"] = internationalDetails.GetContinents()
		details["countries"] = internationalDetails.GetCountries()
	default:
		return nil, fmt.Errorf("league details are required")
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal details: %v", err)
	}

	league, err := queries.CreateLeague(ctx, models.CreateLeagueParams{
		Name:       req.GetName(),
		Leaguetype: req.GetLeagueType().String(),
		Details:    detailsJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create league: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.League{
		LeagueId:   league.Leagueid,
		Name:       league.Name,
		LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[league.Leaguetype]),
		Details:    string(league.Details),
	}, nil
}

func (s *LeagueService) GetLeague(ctx context.Context, req *tournament_management.GetLeagueRequest) (*tournament_management.League, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}
	queries := models.New(s.db)

	league, err := queries.GetLeagueByID(ctx, req.GetLeagueId())
	if err != nil {
		return nil, fmt.Errorf("failed to get league: %v", err)
	}

	var detailsMap map[string]interface{}
	err = json.Unmarshal(league.Details, &detailsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal league details: %v", err)
	}

	leagueResponse := &tournament_management.League{
		LeagueId:   league.Leagueid,
		Name:       league.Name,
		LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[league.Leaguetype]),
		Details:    string(league.Details),
	}

	return leagueResponse, nil
}

func (s *LeagueService) ListLeagues(ctx context.Context, req *tournament_management.ListLeaguesRequest) (*tournament_management.ListLeaguesResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}
	queries := models.New(s.db)

	leagues, err := queries.ListLeaguesPaginated(ctx, models.ListLeaguesPaginatedParams{
		Limit:  int32(req.GetPageSize()),
		Offset: int32(req.GetPageToken()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list leagues: %v", err)
	}

	leagueResponses := make([]*tournament_management.League, len(leagues))
	for i, league := range leagues {
		leagueResponses[i] = &tournament_management.League{
			LeagueId:   league.Leagueid,
			Name:       league.Name,
			LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[league.Leaguetype]),
			Details:    string(league.Details),
		}
	}

	return &tournament_management.ListLeaguesResponse{
		Leagues:       leagueResponses,
		NextPageToken: int32(req.GetPageToken()) + int32(req.GetPageSize()),
	}, nil
}

func (s *LeagueService) UpdateLeague(ctx context.Context, req *tournament_management.UpdateLeagueRequest) (*tournament_management.League, error) {
	if err := s.validateAdminRole(req.GetToken()); err != nil {
		return nil, err
	}

	// Validate required fields
	if req.GetLeagueId() == 0 {
		return nil, fmt.Errorf("league ID is required")
	}
	if req.GetName() == "" {
		return nil, fmt.Errorf("league name is required")
	}
	if req.GetLeagueType() != tournament_management.LeagueType_local && req.GetLeagueType() != tournament_management.LeagueType_international {
		return nil, fmt.Errorf("invalid league type: must be either LOCAL or INTERNATIONAL")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	details := make(map[string]interface{})
	switch req.GetLeagueDetails().(type) {
	case *tournament_management.UpdateLeagueRequest_LocalDetails:
		localDetails := req.GetLocalDetails()
		if len(localDetails.GetProvinces()) == 0 || len(localDetails.GetDistricts()) == 0 {
			return nil, fmt.Errorf("both province and district are required for local leagues")
		}
		details["provinces"] = localDetails.GetProvinces()
		details["districts"] = localDetails.GetDistricts()
	case *tournament_management.UpdateLeagueRequest_InternationalDetails:
		internationalDetails := req.GetInternationalDetails()
		if len(internationalDetails.GetContinents()) == 0 || len(internationalDetails.GetCountries()) == 0 {
			return nil, fmt.Errorf("both continent and country are required for international leagues")
		}
		details["continents"] = internationalDetails.GetContinents()
		details["countries"] = internationalDetails.GetCountries()
	default:
		return nil, fmt.Errorf("league details are required")
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal details: %v", err)
	}

	updatedLeague, err := queries.UpdateLeague(ctx, models.UpdateLeagueParams{
		Leagueid:   req.GetLeagueId(),
		Name:       req.GetName(),
		Leaguetype: req.GetLeagueType().String(),
		Details:    detailsJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update league: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.League{
		LeagueId:   updatedLeague.Leagueid,
		Name:       updatedLeague.Name,
		LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[updatedLeague.Leaguetype]),
		Details:    string(updatedLeague.Details),
	}, nil
}

func (s *LeagueService) DeleteLeague(ctx context.Context, req *tournament_management.DeleteLeagueRequest) (*tournament_management.DeleteLeagueResponse, error) {
	if err := s.validateAdminRole(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)

	err := queries.DeleteLeagueByID(ctx, req.GetLeagueId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete league: %v", err)
	}

	return &tournament_management.DeleteLeagueResponse{
		Success: true,
		Message: "League deleted successfully",
	}, nil
}

func (s *LeagueService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}

func (s *LeagueService) validateAdminRole(token string) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return fmt.Errorf("unauthorized: only admins can perform this action")
	}

	return nil
}
