package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

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
	claims, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	creatorEmail, ok := claims["user_email"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get creator email from token")
	}
	creatorName, ok := claims["user_name"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get creator name from token")
	}

	queries := models.New(s.db)

	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	// Create the tournament
	startDate, err := time.Parse("2006-01-02 15:04", req.GetStartDate())
    if err != nil {
        return nil, fmt.Errorf("invalid start date format: %v", err)
    }

    endDate, err := time.Parse("2006-01-02 15:04", req.GetEndDate())
    if err != nil {
        return nil, fmt.Errorf("invalid end date format: %v", err)
    }

    tournament, err := qtx.CreateTournamentEntry(ctx, models.CreateTournamentEntryParams{
        Name:                       req.GetName(),
        Startdate:                  startDate,
        Enddate:                    endDate,
		Location:                   req.GetLocation(),
		Formatid:                   req.GetFormatId(),
		Leagueid:                   sql.NullInt32{Int32: req.GetLeagueId(), Valid: true},
		Numberofpreliminaryrounds:  int32(req.GetNumberOfPreliminaryRounds()),
		Numberofeliminationrounds:  int32(req.GetNumberOfEliminationRounds()),
		Judgesperdebatepreliminary: int32(req.GetJudgesPerDebatePreliminary()),
		Judgesperdebateelimination: int32(req.GetJudgesPerDebateElimination()),
		Tournamentfee:              fmt.Sprintf("%.2f", req.GetTournamentFee()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tournament: %v", err)
	}

	// Fetch the associated league and format
	league, err := qtx.GetLeagueByID(ctx, req.GetLeagueId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch league: %v", err)
	}

	format, err := qtx.GetTournamentFormatByID(ctx, req.GetFormatId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tournament format: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

// Convert GetLeagueByIDRow to League
leagueForInvitation := models.League{
	Leagueid:   league.Leagueid,
	Name:       league.Name,
	Leaguetype: league.Leaguetype,
	Details:    league.Details,
}

// Send tournament invitations
go func() {
	err := utils.SendTournamentInvitations(context.Background(), tournament, leagueForInvitation, format, queries)
	if err != nil {
		fmt.Printf("Failed to send tournament invitations: %v\n", err)
	}
}()

	// Send confirmation email to the tournament creator
	go func() {
		err := utils.SendTournamentCreationConfirmation(creatorEmail, creatorName, tournament.Name)
		if err != nil {
			fmt.Printf("Failed to send tournament creation confirmation: %v\n", err)
		}
	}()

	return tournamentModelToProto(tournament), nil
}

func (s *TournamentService) GetTournament(ctx context.Context, req *tournament_management.GetTournamentRequest) (*tournament_management.Tournament, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}
	queries := models.New(s.db)
	tournament, err := queries.GetTournamentByID(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament: %v", err)
	}

	return tournamentRowToProto(tournament), nil
}

func (s *TournamentService) ListTournaments(ctx context.Context, req *tournament_management.ListTournamentsRequest) (*tournament_management.ListTournamentsResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}
	queries := models.New(s.db)
	tournaments, err := queries.ListTournamentsPaginated(ctx, models.ListTournamentsPaginatedParams{
		Limit:  int32(req.GetPageSize()),
		Offset: int32(req.GetPageToken()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tournaments: %v", err)
	}

	tournamentResponses := make([]*tournament_management.Tournament, len(tournaments))
	for i, tournament := range tournaments {
		tournamentResponses[i] = tournamentPaginatedRowToProto(tournament)
	}

	return &tournament_management.ListTournamentsResponse{
		Tournaments:   tournamentResponses,
		NextPageToken: int32(req.GetPageToken()) + int32(req.GetPageSize()),
	}, nil
}

func tournamentPaginatedRowToProto(t models.ListTournamentsPaginatedRow) *tournament_management.Tournament {
	return &tournament_management.Tournament{
		TournamentId:               t.Tournamentid,
		Name:                       t.Name,
		StartDate:                  t.Startdate.Format("2006-01-02 15:04"),
		EndDate:                    t.Enddate.Format("2006-01-02 15:04"),
		Location:                   t.Location,
		FormatId:                   t.Formatid,
		LeagueId:                   t.Leagueid.Int32,
		NumberOfPreliminaryRounds:  int32(t.Numberofpreliminaryrounds),
		NumberOfEliminationRounds:  int32(t.Numberofeliminationrounds),
		JudgesPerDebatePreliminary: int32(t.Judgesperdebatepreliminary),
		JudgesPerDebateElimination: int32(t.Judgesperdebateelimination),
		TournamentFee:              parseFloat64(t.Tournamentfee),
	}
}

func (s *TournamentService) UpdateTournament(ctx context.Context, req *tournament_management.UpdateTournamentRequest) (*tournament_management.Tournament, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	startDate, err := time.Parse("2006-01-02 15:04", req.GetStartDate())
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02 15:04", req.GetEndDate())
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	queries := models.New(s.db)
	updatedTournament, err := queries.UpdateTournamentDetails(ctx, models.UpdateTournamentDetailsParams{
		Tournamentid:               req.GetTournamentId(),
		Name:                       req.GetName(),
		Startdate:                  startDate,
		Enddate:                    endDate,
		Location:                   req.GetLocation(),
		Formatid:                   req.GetFormatId(),
		Leagueid:                   sql.NullInt32{Int32: req.GetLeagueId(), Valid: true},
		Numberofpreliminaryrounds:  int32(req.GetNumberOfPreliminaryRounds()),
		Numberofeliminationrounds:  int32(req.GetNumberOfEliminationRounds()),
		Judgesperdebatepreliminary: int32(req.GetJudgesPerDebatePreliminary()),
		Judgesperdebateelimination: int32(req.GetJudgesPerDebateElimination()),
		Tournamentfee:              fmt.Sprintf("%.2f", req.GetTournamentFee()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update tournament details: %v", err)
	}

	return tournamentModelToProto(updatedTournament), nil
}

func (s *TournamentService) DeleteTournament(ctx context.Context, req *tournament_management.DeleteTournamentRequest) error {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return err
	}


	queries := models.New(s.db)
	err = queries.DeleteTournamentByID(ctx, req.GetTournamentId())
	if err != nil {
		return fmt.Errorf("failed to delete tournament: %v", err)
	}

	return nil
}

func (s *TournamentService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}

func (s *TournamentService) validateAdminRole(token string) (map[string]interface{}, error) {
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

func tournamentModelToProto(t models.Tournament) *tournament_management.Tournament {
    return &tournament_management.Tournament{
        TournamentId:               t.Tournamentid,
        Name:                       t.Name,
        StartDate:                  t.Startdate.Format("2006-01-02 15:04"),
        EndDate:                    t.Enddate.Format("2006-01-02 15:04"),
		Location:                   t.Location,
		FormatId:                   t.Formatid,
		LeagueId:                   t.Leagueid.Int32,
		NumberOfPreliminaryRounds:  int32(t.Numberofpreliminaryrounds),
		NumberOfEliminationRounds:  int32(t.Numberofeliminationrounds),
		JudgesPerDebatePreliminary: int32(t.Judgesperdebatepreliminary),
		JudgesPerDebateElimination: int32(t.Judgesperdebateelimination),
		TournamentFee:              parseFloat64(t.Tournamentfee),
	}
}

func tournamentRowToProto(t models.GetTournamentByIDRow) *tournament_management.Tournament {
    return &tournament_management.Tournament{
        TournamentId:               t.Tournamentid,
        Name:                       t.Name,
        StartDate:                  t.Startdate.Format("2006-01-02 15:04"),
        EndDate:                    t.Enddate.Format("2006-01-02 15:04"),
		Location:                   t.Location,
		FormatId:                   t.Formatid,
		LeagueId:                   t.Leagueid.Int32,
		NumberOfPreliminaryRounds:  int32(t.Numberofpreliminaryrounds),
		NumberOfEliminationRounds:  int32(t.Numberofeliminationrounds),
		JudgesPerDebatePreliminary: int32(t.Judgesperdebatepreliminary),
		JudgesPerDebateElimination: int32(t.Judgesperdebateelimination),
		TournamentFee:              parseFloat64(t.Tournamentfee),
	}
}

func parseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}