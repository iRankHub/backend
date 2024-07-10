package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type TournamentService struct {
	db *sql.DB
}

func NewTournamentService(db *sql.DB) *TournamentService {
	return &TournamentService{db: db}
}

func (s *TournamentService) CreateTournament(ctx context.Context, req *tournament_management.CreateTournamentRequest) (*tournament_management.Tournament, error) {
	// Check if the user is an admin
	claims, err := utils.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	userRole := claims["user_role"].(string)
	if userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can create tournaments")
	}

	// Create the tournament
	tournament, err := s.db.CreateTournamentEntry(ctx, models.CreateTournamentEntryParams{
		Name:                       req.GetName(),
		StartDate:                  req.GetStartDate().AsTime(),
		EndDate:                    req.GetEndDate().AsTime(),
		Location:                   req.GetLocation(),
		FormatID:                   req.GetFormatId(),
		LeagueID:                   req.GetLeagueId(),
		NumberOfPreliminaryRounds:  int32(req.GetNumberOfPreliminaryRounds()),
		NumberOfEliminationRounds:  int32(req.GetNumberOfEliminationRounds()),
		JudgesPerDebatePreliminary: int32(req.GetJudgesPerDebatePreliminary()),
		JudgesPerDebateElimination: int32(req.GetJudgesPerDebateElimination()),
		TournamentFee:              req.GetTournamentFee(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tournament: %v", err)
	}

	return &tournament_management.Tournament{
		TournamentId:               tournament.TournamentID,
		Name:                       tournament.Name,
		StartDate:                  &timestamppb.Timestamp{Seconds: tournament.StartDate.Unix()},
		EndDate:                    &timestamppb.Timestamp{Seconds: tournament.EndDate.Unix()},
		Location:                   tournament.Location,
		FormatId:                   tournament.FormatID,
		LeagueId:                   tournament.LeagueID,
		NumberOfPreliminaryRounds:  int32(tournament.NumberOfPreliminaryRounds),
		NumberOfEliminationRounds:  int32(tournament.NumberOfEliminationRounds),
		JudgesPerDebatePreliminary: int32(tournament.JudgesPerDebatePreliminary),
		JudgesPerDebateElimination: int32(tournament.JudgesPerDebateElimination),
		TournamentFee:              tournament.TournamentFee,
	}, nil
}
func (s *TournamentService) GetTournament(ctx context.Context, req *tournament_management.GetTournamentRequest) (*tournament_management.Tournament, error) {
	// Get the tournament by ID
	tournament, err := s.db.GetTournamentByID(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament: %v", err)
	}

	return &tournament_management.Tournament{
		TournamentId:               tournament.TournamentID,
		Name:                       tournament.Name,
		StartDate:                  &timestamppb.Timestamp{Seconds: tournament.StartDate.Unix()},
		EndDate:                    &timestamppb.Timestamp{Seconds: tournament.EndDate.Unix()},
		Location:                   tournament.Location,
		FormatId:                   tournament.FormatID,
		LeagueId:                   tournament.LeagueID,
		NumberOfPreliminaryRounds:  int32(tournament.NumberOfPreliminaryRounds),
		NumberOfEliminationRounds:  int32(tournament.NumberOfEliminationRounds),
		JudgesPerDebatePreliminary: int32(tournament.JudgesPerDebatePreliminary),
		JudgesPerDebateElimination: int32(tournament.JudgesPerDebateElimination),
		TournamentFee:              tournament.TournamentFee,
	}, nil

}
func (s *TournamentService) ListTournaments(ctx context.Context, req *tournament_management.ListTournamentsRequest) (*tournament_management.ListTournamentsResponse, error) {
	// List tournaments with pagination
	tournaments, err := s.db.ListTournamentsPaginated(ctx, models.ListTournamentsPaginatedParams{
		Limit:  int32(req.GetPageSize()),
		Offset: int32(req.GetPageToken()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tournaments: %v", err)
	}
	// Construct the Tournament responses
	tournamentResponses := make([]*tournament_management.Tournament, len(tournaments))
	for i, tournament := range tournaments {
		tournamentResponses[i] = &tournament_management.Tournament{
			TournamentId:               tournament.TournamentID,
			Name:                       tournament.Name,
			StartDate:                  &timestamppb.Timestamp{Seconds: tournament.StartDate.Unix()},
			EndDate:                    &timestamppb.Timestamp{Seconds: tournament.EndDate.Unix()},
			Location:                   tournament.Location,
			FormatId:                   tournament.FormatID,
			LeagueId:                   tournament.LeagueID,
			NumberOfPreliminaryRounds:  int32(tournament.NumberOfPreliminaryRounds),
			NumberOfEliminationRounds:  int32(tournament.NumberOfEliminationRounds),
			JudgesPerDebatePreliminary: int32(tournament.JudgesPerDebatePreliminary),
			JudgesPerDebateElimination: int32(tournament.JudgesPerDebateElimination),
			TournamentFee:              tournament.TournamentFee,
		}
	}

	return &tournament_management.ListTournamentsResponse{
		Tournaments:   tournamentResponses,
		NextPageToken: int32(req.GetPageToken()) + int32(req.GetPageSize()),
	}, nil

}
func (s *TournamentService) UpdateTournament(ctx context.Context, req *tournament_management.UpdateTournamentRequest) (*tournament_management.Tournament, error) {
	// Check if the user is an admin
	claims, err := utils.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	userRole := claims["user_role"].(string)
	if userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can update tournaments")
	}
	// Update the tournament details
	updatedTournament, err := s.db.UpdateTournamentDetails(ctx, models.UpdateTournamentDetailsParams{
		TournamentID:               req.GetTournamentId(),
		Name:                       req.GetName(),
		StartDate:                  req.GetStartDate().AsTime(),
		EndDate:                    req.GetEndDate().AsTime(),
		Location:                   req.GetLocation(),
		FormatID:                   req.GetFormatId(),
		LeagueID:                   req.GetLeagueId(),
		NumberOfPreliminaryRounds:  int32(req.GetNumberOfPreliminaryRounds()),
		NumberOfEliminationRounds:  int32(req.GetNumberOfEliminationRounds()),
		JudgesPerDebatePreliminary: int32(req.GetJudgesPerDebatePreliminary()),
		JudgesPerDebateElimination: int32(req.GetJudgesPerDebateElimination()),
		TournamentFee:              req.GetTournamentFee(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update tournament details: %v", err)
	}

	return &tournament_management.Tournament{
		TournamentId:               updatedTournament.TournamentID,
		Name:                       updatedTournament.Name,
		StartDate:                  &timestamppb.Timestamp{Seconds: updatedTournament.StartDate.Unix()},
		EndDate:                    &timestamppb.Timestamp{Seconds: updatedTournament.EndDate.Unix()},
		Location:                   updatedTournament.Location,
		FormatId:                   updatedTournament.FormatID,
		LeagueId:                   updatedTournament.LeagueID,
		NumberOfPreliminaryRounds:  int32(updatedTournament.NumberOfPreliminaryRounds),
		NumberOfEliminationRounds:  int32(updatedTournament.NumberOfEliminationRounds),
		JudgesPerDebatePreliminary: int32(updatedTournament.JudgesPerDebatePreliminary),
		JudgesPerDebateElimination: int32(updatedTournament.JudgesPerDebateElimination),
		TournamentFee:              updatedTournament.TournamentFee,
	}, nil
}
func (s *TournamentService) DeleteTournament(ctx context.Context, req *tournament_management.DeleteTournamentRequest) (*tournament_management.Empty, error) {
	// Check if the user is an admin
	claims, err := utils.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	userRole := claims["user_role"].(string)
	if userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can delete tournaments")
	}
	// Delete the tournament by ID
	err = s.db.DeleteTournamentByID(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete tournament: %v", err)
	}

	return &tournament_management.Empty{}, nil
}
