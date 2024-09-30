package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/models"
	notifications "github.com/iRankHub/backend/internal/services/notification"
	"github.com/iRankHub/backend/internal/utils"
	notification "github.com/iRankHub/backend/internal/utils/notifications"
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

	creatorID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("failed to get creator ID from token")
	}

	// Convert creatorID to int32
	creatorIDInt32 := int32(creatorID)

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
		Imageurl:                   sql.NullString{String: req.GetImageUrl(), Valid: req.GetImageUrl() != ""},
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

	// Send notification asynchronously
	go func() {
		// Use a new context for background operations
		bgCtx := context.Background()

		// Create a new NotificationService using the existing database connection
		notificationService, err := notifications.NewNotificationService(s.db)
		if err != nil {
			log.Printf("Failed to create notification service: %v", err)
			return
		}

		// Use the existing queries with the main database connection
		bgQueries := models.New(s.db)

		if err := notification.SendTournamentInvitations(bgCtx, notificationService, tournament, league, format, bgQueries); err != nil {
			log.Printf("Failed to send tournament invitations: %v", err)
		}

		if err := notification.SendTournamentCreationConfirmation(notificationService, creatorEmail, creatorName, tournament.Name, creatorIDInt32); err != nil {
			log.Printf("Failed to send tournament creation confirmation: %v", err)
		}

		if err := notification.SendCoordinatorAssignmentEmail(notificationService, coordinator, tournament, league, format); err != nil {
			log.Printf("Failed to send coordinator assignment email: %v", err)
		}
	}()

	createdTournament := tournamentModelToProto(tournament)
	createdTournament.CoordinatorName = coordinator.Name

	return createdTournament, nil
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

func (s *TournamentService) GetTournamentStats(ctx context.Context, req *tournament_management.GetTournamentStatsRequest) (*tournament_management.GetTournamentStatsResponse, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	stats, err := queries.GetTournamentStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament stats: %v", err)
	}

	totalPercentageChange := calculatePercentageChange(int64(stats.YesterdayTotalCount.Int32), stats.TotalTournaments)
	upcomingPercentageChange := calculatePercentageChange(int64(stats.YesterdayUpcomingCount.Int32), stats.UpcomingTournaments)

	return &tournament_management.GetTournamentStatsResponse{
		TotalTournaments:         int32(stats.TotalTournaments),
		UpcomingTournaments:      int32(stats.UpcomingTournaments),
		TotalPercentageChange:    totalPercentageChange,
		UpcomingPercentageChange: upcomingPercentageChange,
	}, nil
}

func calculatePercentageChange(oldValue, newValue int64) string {
	if oldValue == 0 && newValue == 0 {
		return "0.00%"
	}
	if oldValue == 0 {
		return "0.00%"
	}
	change := float64(newValue-oldValue) / float64(oldValue) * 100
	sign := "+"
	if change < 0 {
		sign = "-"
		change = math.Abs(change)
	}
	return fmt.Sprintf("%s%.2f%%", sign, change)
}

func (s *TournamentService) GetTournamentRegistrations(ctx context.Context, req *tournament_management.GetTournamentRegistrationsRequest) (*tournament_management.GetTournamentRegistrationsResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	registrations, err := queries.GetTournamentRegistrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament registrations: %v", err)
	}

	response := &tournament_management.GetTournamentRegistrationsResponse{
		Registrations: make([]*tournament_management.DailyRegistration, len(registrations)),
	}

	for i, reg := range registrations {
		response.Registrations[i] = &tournament_management.DailyRegistration{
			Date:  reg.RegistrationDate.Format("2006-01-02"),
			Count: int32(reg.RegistrationCount),
		}
	}

	return response, nil
}

func tournamentPaginatedRowToProto(t models.ListTournamentsPaginatedRow) *tournament_management.Tournament {
	return &tournament_management.Tournament{
		TournamentId:               t.Tournamentid,
		Name:                       t.Name,
		ImageUrl:                   t.Imageurl.String,
		StartDate:                  t.Startdate.Format("2006-01-02 15:04"),
		EndDate:                    t.Enddate.Format("2006-01-02 15:04"),
		Location:                   t.Location,
		FormatId:                   t.Formatid,
		LeagueId:                   t.Leagueid.Int32,
		CoordinatorId:              t.Coordinatorid,
		CoordinatorName:            t.Coordinatorname,
		NumberOfPreliminaryRounds:  int32(t.Numberofpreliminaryrounds),
		NumberOfEliminationRounds:  int32(t.Numberofeliminationrounds),
		JudgesPerDebatePreliminary: int32(t.Judgesperdebatepreliminary),
		JudgesPerDebateElimination: int32(t.Judgesperdebateelimination),
		TournamentFee:              parseFloat64(t.Tournamentfee),
		NumberOfSchools:            int32(t.Acceptedschoolscount),
		NumberOfTeams:              int32(t.Teamscount),
		LeagueName:                 t.Leaguename,
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

	result, err := queries.UpdateTournamentDetails(ctx, models.UpdateTournamentDetailsParams{
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
		Imageurl:                   sql.NullString{String: req.GetImageUrl(), Valid: req.GetImageUrl() != ""},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pairings have already been made")
		}
		return nil, fmt.Errorf("failed to update tournament details: %v", err)
	}

	if result.ErrorMessage.(bool) && result.ErrorMessage.(string) != "" {
		return nil, fmt.Errorf(result.ErrorMessage.(string))
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	updatedTournament := models.Tournament{
		Tournamentid:               result.Tournamentid,
		Name:                       result.Name,
		Imageurl:                   result.Imageurl,
		Startdate:                  result.Startdate,
		Enddate:                    result.Enddate,
		Location:                   result.Location,
		Formatid:                   result.Formatid,
		Leagueid:                   result.Leagueid,
		Numberofpreliminaryrounds:  result.Numberofpreliminaryrounds,
		Numberofeliminationrounds:  result.Numberofeliminationrounds,
		Judgesperdebatepreliminary: result.Judgesperdebatepreliminary,
		Judgesperdebateelimination: result.Judgesperdebateelimination,
		Tournamentfee:              result.Tournamentfee,
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

	// Create invitations for schools
	for _, school := range schools {
		invitation, err := queries.CreateInvitation(ctx, models.CreateInvitationParams{
			Tournamentid: tournamentID,
			Inviteeid:    school.Idebateschoolid.String,
			Inviteerole:  "school",
			Status:       "pending",
		})
		if err != nil {
			log.Printf("Failed to create invitation for school %d: %v", school.Schoolid, err)
			return fmt.Errorf("failed to create invitation for school %d: %v", school.Schoolid, err)
		}
		log.Printf("Created invitation %d for school %d", invitation.Invitationid, school.Schoolid)
	}

	// Create invitations for volunteers
	volunteers, err := queries.GetAllVolunteers(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch volunteers: %v", err)
	}

	for _, volunteer := range volunteers {
		invitation, err := queries.CreateInvitation(ctx, models.CreateInvitationParams{
			Tournamentid: tournamentID,
			Inviteeid:    volunteer.Idebatevolunteerid.String,
			Inviteerole:  "volunteer",
			Status:       "pending",
		})
		if err != nil {
			log.Printf("Failed to create invitation for volunteer %d: %v", volunteer.Volunteerid, err)
			return fmt.Errorf("failed to create invitation for volunteer %d: %v", volunteer.Volunteerid, err)
		}
		log.Printf("Created invitation %d for volunteer %d", invitation.Invitationid, volunteer.Volunteerid)
	}

	// For DAC league, also invite all students
	if league.Name == "DAC" {
		students, err := queries.GetAllStudents(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch all students: %v", err)
		}

		for _, student := range students {
			invitation, err := queries.CreateInvitation(ctx, models.CreateInvitationParams{
				Tournamentid: tournamentID,
				Inviteeid:    student.Idebatestudentid.String,
				Inviteerole:  "student",
				Status:       "pending",
			})
			if err != nil {
				log.Printf("Failed to create invitation for student %d: %v", student.Studentid, err)
				return fmt.Errorf("failed to create invitation for student %d: %v", student.Studentid, err)
			}
			log.Printf("Created invitation %d for student %d", invitation.Invitationid, student.Studentid)
		}
	}

	log.Printf("Finished creating invitations for tournament %d", tournamentID)
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
		ImageUrl:                   t.Imageurl.String,
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
		ImageUrl:                   t.Imageurl.String,
		StartDate:                  t.Startdate.Format("2006-01-02 15:04"),
		EndDate:                    t.Enddate.Format("2006-01-02 15:04"),
		Location:                   t.Location,
		FormatId:                   t.Formatid,
		LeagueId:                   t.Leagueid.Int32,
		CoordinatorId:              t.Coordinatorid,
		CoordinatorName:            t.Coordinatorname,
		NumberOfPreliminaryRounds:  int32(t.Numberofpreliminaryrounds),
		NumberOfEliminationRounds:  int32(t.Numberofeliminationrounds),
		JudgesPerDebatePreliminary: int32(t.Judgesperdebatepreliminary),
		JudgesPerDebateElimination: int32(t.Judgesperdebateelimination),
		TournamentFee:              parseFloat64(t.Tournamentfee),
		LeagueName:                 t.Leaguename,
	}
}

func parseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
