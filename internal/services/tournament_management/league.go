package services

import (
	"context"
	"database/sql"
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	league, err := queries.CreateLeague(ctx, models.CreateLeagueParams{
		Name:       req.GetName(),
		Leaguetype: req.GetLeagueType().String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create league: %v", err)
	}

	switch req.GetLeagueDetails().(type) {
	case *tournament_management.CreateLeagueRequest_LocalDetails:
		localDetails := req.GetLocalDetails()
		err = queries.CreateLocalLeagueDetails(ctx, models.CreateLocalLeagueDetailsParams{
			Leagueid: league.Leagueid,
			Province: sql.NullString{String: localDetails.GetProvince(), Valid: localDetails.GetProvince() != ""},
			District: sql.NullString{String: localDetails.GetDistrict(), Valid: localDetails.GetDistrict() != ""},
		})
	case *tournament_management.CreateLeagueRequest_InternationalDetails:
		internationalDetails := req.GetInternationalDetails()
		err = queries.CreateInternationalLeagueDetails(ctx, models.CreateInternationalLeagueDetailsParams{
			Leagueid:  league.Leagueid,
			Continent: sql.NullString{String: internationalDetails.GetContinent(), Valid: internationalDetails.GetContinent() != ""},
			Country:   sql.NullString{String: internationalDetails.GetCountry(), Valid: internationalDetails.GetCountry() != ""},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create league details: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.League{
		LeagueId:   league.Leagueid,
		Name:       league.Name,
		LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[league.Leaguetype]),
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

	leagueResponse := &tournament_management.League{
		LeagueId:   league.Leagueid,
		Name:       league.Name,
		LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[league.Leaguetype]),
	}

	if league.Leaguetype == "LEAGUE_TYPE_LOCAL" {
		leagueResponse.LeagueDetails = &tournament_management.League_LocalDetails{
			LocalDetails: &tournament_management.LocalLeagueDetails{
				Province: league.Detail1.String,
				District: league.Detail2.String,
			},
		}
	} else if league.Leaguetype == "LEAGUE_TYPE_INTERNATIONAL" {
		leagueResponse.LeagueDetails = &tournament_management.League_InternationalDetails{
			InternationalDetails: &tournament_management.InternationalLeagueDetails{
				Continent: league.Detail1.String,
				Country:   league.Detail2.String,
			},
		}
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
		leagueResponse := &tournament_management.League{
			LeagueId:   league.Leagueid,
			Name:       league.Name,
			LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[league.Leaguetype]),
		}

		if league.Leaguetype == "LEAGUE_TYPE_LOCAL" {
			leagueResponse.LeagueDetails = &tournament_management.League_LocalDetails{
				LocalDetails: &tournament_management.LocalLeagueDetails{
					Province: league.Detail1.String,
					District: league.Detail2.String,
				},
			}
		} else if league.Leaguetype == "LEAGUE_TYPE_INTERNATIONAL" {
			leagueResponse.LeagueDetails = &tournament_management.League_InternationalDetails{
				InternationalDetails: &tournament_management.InternationalLeagueDetails{
					Continent: league.Detail1.String,
					Country:   league.Detail2.String,
				},
			}
		}

		leagueResponses[i] = leagueResponse
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	updatedLeague, err := queries.UpdateLeagueDetails(ctx, models.UpdateLeagueDetailsParams{
		Leagueid:   req.GetLeagueId(),
		Name:       req.GetName(),
		Leaguetype: req.GetLeagueType().String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update league details: %v", err)
	}

	switch req.GetLeagueDetails().(type) {
	case *tournament_management.UpdateLeagueRequest_LocalDetails:
		localDetails := req.GetLocalDetails()
		err = queries.UpdateLocalLeagueDetailsInfo(ctx, models.UpdateLocalLeagueDetailsInfoParams{
			Leagueid: req.GetLeagueId(),
			Province: sql.NullString{String: localDetails.GetProvince(), Valid: localDetails.GetProvince() != ""},
			District: sql.NullString{String: localDetails.GetDistrict(), Valid: localDetails.GetDistrict() != ""},
		})
	case *tournament_management.UpdateLeagueRequest_InternationalDetails:
		internationalDetails := req.GetInternationalDetails()
		err = queries.UpdateInternationalLeagueDetailsInfo(ctx, models.UpdateInternationalLeagueDetailsInfoParams{
			Leagueid:  req.GetLeagueId(),
			Continent: sql.NullString{String: internationalDetails.GetContinent(), Valid: internationalDetails.GetContinent() != ""},
			Country:   sql.NullString{String: internationalDetails.GetCountry(), Valid: internationalDetails.GetCountry() != ""},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update league details: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.League{
		LeagueId:   updatedLeague.Leagueid,
		Name:       updatedLeague.Name,
		LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[updatedLeague.Leaguetype]),
	}, nil
}

func (s *LeagueService) DeleteLeague(ctx context.Context, req *tournament_management.DeleteLeagueRequest) (bool, error) {
	if err := s.validateAdminRole(req.GetToken()); err != nil {
		return false, err
	}

	queries := models.New(s.db)

	err := queries.DeleteLeagueByID(ctx, req.GetLeagueId())
	if err != nil {
		return false, fmt.Errorf("failed to delete league: %v", err)
	}

	return true, nil
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