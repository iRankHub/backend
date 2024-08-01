package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
	emails "github.com/iRankHub/backend/internal/utils/emails"
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	// Verify that the coordinator is a volunteer or admin
	coordinator, err := queries.GetUserByID(ctx, req.GetCoordinatorId())
	if err != nil {
		return nil, fmt.Errorf("failed to get coordinator: %v", err)
	}
	if coordinator.Userrole != "volunteer" && coordinator.Userrole != "admin" {
		return nil, fmt.Errorf("coordinator must be a volunteer or admin")
	}

	startDate, err := time.Parse("2006-01-02 15:04", req.GetStartDate())
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02 15:04", req.GetEndDate())
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	tournament, err := queries.CreateTournamentEntry(ctx, models.CreateTournamentEntryParams{
		Name:                       req.GetName(),
		Startdate:                  startDate,
		Enddate:                    endDate,
		Location:                   req.GetLocation(),
		Formatid:                   req.GetFormatId(),
		Leagueid:                   sql.NullInt32{Int32: req.GetLeagueId(), Valid: true},
		Coordinatorid:              req.GetCoordinatorId(),
		Numberofpreliminaryrounds:  int32(req.GetNumberOfPreliminaryRounds()),
		Numberofeliminationrounds:  int32(req.GetNumberOfEliminationRounds()),
		Judgesperdebatepreliminary: int32(req.GetJudgesPerDebatePreliminary()),
		Judgesperdebateelimination: int32(req.GetJudgesPerDebateElimination()),
		Tournamentfee:              fmt.Sprintf("%.2f", req.GetTournamentFee()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tournament: %v", err)
	}

	league, err := queries.GetLeagueByID(ctx, req.GetLeagueId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch league: %v", err)
	}

	format, err := queries.GetTournamentFormatByID(ctx, req.GetFormatId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tournament format: %v", err)
	}

	err = s.createInvitations(ctx, queries, tournament.Tournamentid, league)
	if err != nil {
		return nil, fmt.Errorf("failed to create invitations: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Send tournament invitations and confirmation emails asynchronously
	go func() {
		if err := emails.SendTournamentInvitations(context.Background(), tournament, league, format, models.New(s.db)); err != nil {
			log.Printf("Failed to send tournament invitations: %v\n", err)
		}
	}()

	go func() {
		if err := emails.SendTournamentCreationConfirmation(creatorEmail, creatorName, tournament.Name); err != nil {
			fmt.Printf("Failed to send tournament creation confirmation: %v\n", err)
		}
	}()

	// Send coordinator assignment email
	go func() {
		if err := emails.SendCoordinatorAssignmentEmail(coordinator, tournament, league, format); err != nil {
			fmt.Printf("Failed to send coordinator assignment email: %v\n", err)
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

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

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return tournamentModelToProto(updatedTournament), nil
}

func (s *TournamentService) DeleteTournament(ctx context.Context, req *tournament_management.DeleteTournamentRequest) (*tournament_management.DeleteTournamentResponse, error) {
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

	if err := queries.DeleteTournamentByID(ctx, req.GetTournamentId()); err != nil {
		return nil, fmt.Errorf("failed to delete tournament: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.DeleteTournamentResponse{
		Success: true,
		Message: "Tournament deleted successfully",
	}, nil
}

func (s *TournamentService) createInvitations(ctx context.Context, queries *models.Queries, tournamentID int32, league models.League) error {
    log.Printf("Creating invitations for tournament %d", tournamentID)

    // Invite schools based on league details
    var leagueDetails struct {
        Districts []string `json:"districts,omitempty"`
        Countries []string `json:"countries,omitempty"`
    }

    if err := json.Unmarshal(league.Details, &leagueDetails); err != nil {
        return fmt.Errorf("failed to unmarshal league details: %v", err)
    }

    var schools []models.School

    if league.Leaguetype == "local" {
        for _, district := range leagueDetails.Districts {
            schoolsBatch, err := queries.GetSchoolsByDistrict(ctx, sql.NullString{String: district, Valid: true})
            if err != nil {
                return fmt.Errorf("failed to fetch schools for district %s: %v", district, err)
            }
            schools = append(schools, schoolsBatch...)
        }
    } else if league.Leaguetype == "international" {
        for _, country := range leagueDetails.Countries {
            schoolsBatch, err := queries.GetSchoolsByCountry(ctx, sql.NullString{String: country, Valid: true})
            if err != nil {
                return fmt.Errorf("failed to fetch schools for country %s: %v", country, err)
            }
            schools = append(schools, schoolsBatch...)
        }
    } else {
        return fmt.Errorf("invalid league type: %s", league.Leaguetype)
    }

    for _, school := range schools {
        _, err := queries.CreateInvitation(ctx, models.CreateInvitationParams{
            Tournamentid: tournamentID,
            Schoolid:     sql.NullInt32{Int32: school.Schoolid, Valid: true},
            Volunteerid:  sql.NullInt32{},
            Studentid:    sql.NullInt32{},
            Status:       "pending",
        })
        if err != nil {
            log.Printf("Failed to create invitation for school %d: %v", school.Schoolid, err)
            return fmt.Errorf("failed to create invitation for school %d: %v", school.Schoolid, err)
        }
    }

    // Create invitations for volunteers
    volunteers, err := queries.GetAllVolunteers(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch volunteers: %v", err)
    }

    for _, volunteer := range volunteers {
        _, err := queries.CreateInvitation(ctx, models.CreateInvitationParams{
            Tournamentid: tournamentID,
            Schoolid:     sql.NullInt32{},
            Volunteerid:  sql.NullInt32{Int32: volunteer.Volunteerid, Valid: true},
            Studentid:    sql.NullInt32{},
            Status:       "pending",
        })
        if err != nil {
            log.Printf("Failed to create invitation for volunteer %d: %v", volunteer.Volunteerid, err)
            return fmt.Errorf("failed to create invitation for volunteer %d: %v", volunteer.Volunteerid, err)
        }
    }

    // For DAC league, also invite all students
    if league.Name == "DAC" {
        students, err := queries.GetAllStudents(ctx)
        if err != nil {
            return fmt.Errorf("failed to fetch all students: %v", err)
        }

        for _, student := range students {
            _, err := queries.CreateInvitation(ctx, models.CreateInvitationParams{
                Tournamentid: tournamentID,
                Schoolid:     sql.NullInt32{},
                Volunteerid:  sql.NullInt32{},
                Studentid:    sql.NullInt32{Int32: student.Studentid, Valid: true},
                Status:       "pending",
            })
            if err != nil {
                log.Printf("Failed to create invitation for student %d: %v", student.Studentid, err)
                return fmt.Errorf("failed to create invitation for student %d: %v", student.Studentid, err)
            }
        }
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

// Helper functions to convert between model and proto types
func tournamentModelToProto(t models.Tournament) *tournament_management.Tournament {
	return &tournament_management.Tournament{
		TournamentId:               t.Tournamentid,
		Name:                       t.Name,
		StartDate:                  t.Startdate.Format("2006-01-02 15:04"),
		EndDate:                    t.Enddate.Format("2006-01-02 15:04"),
		Location:                   t.Location,
		FormatId:                   t.Formatid,
		LeagueId:                   t.Leagueid.Int32,
		CoordinatorId:              t.Coordinatorid,
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
		CoordinatorId:              t.Coordinatorid,
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
