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
	// Check if the user is an admin
	claims, err := utils.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	userRole := claims["user_role"].(string)
	if userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can create leagues")
	}

	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Create the league
	league, err := tx.CreateLeague(ctx, models.CreateLeagueParams{
		Name:       req.GetName(),
		LeagueType: models.LeagueType(req.GetLeagueType().String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create league: %v", err)
	}

	// Create the league details based on the league type
	switch req.GetLeagueDetails().(type) {
	case *tournament_management.CreateLeagueRequest_LocalDetails:
		localDetails := req.GetLocalDetails()
		err = tx.CreateLocalLeagueDetails(ctx, models.CreateLocalLeagueDetailsParams{
			LeagueID: league.LeagueID,
			Province: localDetails.GetProvince(),
			District: localDetails.GetDistrict(),
		})
	case *tournament_management.CreateLeagueRequest_InternationalDetails:
		internationalDetails := req.GetInternationalDetails()
		err = tx.CreateInternationalLeagueDetails(ctx, models.CreateInternationalLeagueDetailsParams{
			LeagueID:  league.LeagueID,
			Continent: internationalDetails.GetContinent(),
			Country:   internationalDetails.GetCountry(),
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create league details: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.League{
		LeagueId:   league.LeagueID,
		Name:       league.Name,
		LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[league.LeagueType.String()]),
	}, nil
}

func (s *LeagueService) GetLeague(ctx context.Context, req *tournament_management.GetLeagueRequest) (*tournament_management.League, error) {
	// Get the league by ID
	league, err := s.db.GetLeagueByID(ctx, req.GetLeagueId())
	if err != nil {
		return nil, fmt.Errorf("failed to get league: %v", err)
	}

	// Construct the League response
	leagueResponse := &tournament_management.League{
		LeagueId:   league.LeagueID,
		Name:       league.Name,
		LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[league.LeagueType.String()]),
	}

	// Set the league details based on the league type
	if league.LeagueType == models.LeagueTypeLocal {
		leagueResponse.LeagueDetails = &tournament_management.League_LocalDetails{
			LocalDetails: &tournament_management.LocalLeagueDetails{
				Province: league.Detail1,
				District: league.Detail2,
			},
		}
	} else if league.LeagueType == models.LeagueTypeInternational {
		leagueResponse.LeagueDetails = &tournament_management.League_InternationalDetails{
			InternationalDetails: &tournament_management.InternationalLeagueDetails{
				Continent: league.Detail1,
				Country:   league.Detail2,
			},
		}
	}

	return leagueResponse, nil
}

func (s *LeagueService) ListLeagues(ctx context.Context, req *tournament_management.ListLeaguesRequest) (*tournament_management.ListLeaguesResponse, error) {
	// List leagues with pagination
	leagues, err := s.db.ListLeaguesPaginated(ctx, models.ListLeaguesPaginatedParams{
		Limit:  int32(req.GetPageSize()),
		Offset: int32(req.GetPageToken()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list leagues: %v", err)
	}

	// Construct the League responses
	leagueResponses := make([]*tournament_management.League, len(leagues))
	for i, league := range leagues {
		leagueResponse := &tournament_management.League{
			LeagueId:   league.LeagueID,
			Name:       league.Name,
			LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[league.LeagueType.String()]),
		}

		// Set the league details based on the league type
		if league.LeagueType == models.LeagueTypeLocal {
			leagueResponse.LeagueDetails = &tournament_management.League_LocalDetails{
				LocalDetails: &tournament_management.LocalLeagueDetails{
					Province: league.Detail1,
					District: league.Detail2,
				},
			}
		} else if league.LeagueType == models.LeagueTypeInternational {
			leagueResponse.LeagueDetails = &tournament_management.League_InternationalDetails{
				InternationalDetails: &tournament_management.InternationalLeagueDetails{
					Continent: league.Detail1,
					Country:   league.Detail2,
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
	// Check if the user is an admin
	claims, err := utils.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	userRole := claims["user_role"].(string)
	if userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can update leagues")
	}

	// Update the league details
	updatedLeague, err := s.db.UpdateLeagueDetails(ctx, models.UpdateLeagueDetailsParams{
		LeagueID:   req.GetLeagueId(),
		Name:       req.GetName(),
		LeagueType: models.LeagueType(req.GetLeagueType().String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update league details: %v", err)
	}

	// Update the league details based on the league type
	switch req.GetLeagueDetails().(type) {
	case *tournament_management.UpdateLeagueRequest_LocalDetails:
		localDetails := req.GetLocalDetails()
		err = s.db.UpdateLocalLeagueDetailsInfo(ctx, models.UpdateLocalLeagueDetailsInfoParams{
			LeagueID: req.GetLeagueId(),
			Province: localDetails.GetProvince(),
			District: localDetails.GetDistrict(),
		})
	case *tournament_management.UpdateLeagueRequest_InternationalDetails:
		internationalDetails := req.GetInternationalDetails()
		err = s.db.UpdateInternationalLeagueDetailsInfo(ctx, models.UpdateInternationalLeagueDetailsInfoParams{
			LeagueID:  req.GetLeagueId(),
			Continent: internationalDetails.GetContinent(),
			Country:   internationalDetails.GetCountry(),
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update league details: %v", err)
	}

	return &tournament_management.League{
		LeagueId:   updatedLeague.LeagueID,
		Name:       updatedLeague.Name,
		LeagueType: tournament_management.LeagueType(tournament_management.LeagueType_value[updatedLeague.LeagueType.String()]),
	}, nil
}

func (s *LeagueService) DeleteLeague(ctx context.Context, req *tournament_management.DeleteLeagueRequest) (*tournament_management.Empty, error) {
	// Check if the user is an admin
	claims, err := utils.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	userRole := claims["user_role"].(string)
	if userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can delete leagues")
	}

	// Delete the league by ID
	err = s.db.DeleteLeagueByID(ctx, req.GetLeagueId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete league: %v", err)
	}

	return &tournament_management.Empty{}, nil
}
