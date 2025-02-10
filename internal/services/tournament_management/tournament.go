package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sqlc-dev/pqtype"
	"log"
	"math"
	"strconv"
	"strings"
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

	// Validate motions
	if err := validateMotions(req); err != nil {
		return nil, fmt.Errorf("invalid motions: %v", err)
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

	// Create motions JSON structure
	motions := map[string]interface{}{
		"preliminary": make([]map[string]interface{}, 0),
		"elimination": make([]map[string]interface{}, 0),
	}

	// Add preliminary motions
	for _, m := range req.GetMotions().GetPreliminaryMotions() {
		motions["preliminary"] = append(motions["preliminary"].([]map[string]interface{}), map[string]interface{}{
			"text":        m.GetText(),
			"infoSlide":   m.GetInfoSlide(),
			"roundNumber": m.GetRoundNumber(),
		})
	}

	// Add elimination motions
	for _, m := range req.GetMotions().GetEliminationMotions() {
		motions["elimination"] = append(motions["elimination"].([]map[string]interface{}), map[string]interface{}{
			"text":        m.GetText(),
			"infoSlide":   m.GetInfoSlide(),
			"roundNumber": m.GetRoundNumber(),
		})
	}

	motionsJSON, err := json.Marshal(motions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal motions: %v", err)
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
		Motions:                    pqtype.NullRawMessage{RawMessage: motionsJSON, Valid: true},
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

	// Send notifications asynchronously
	go func() {
		bgCtx := context.Background()
		notificationService, err := notifications.NewNotificationService(s.db)
		if err != nil {
			log.Printf("Failed to create notification service: %v", err)
			return
		}

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
	createdTournament.Motions = req.GetMotions()

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

	protoTournament := tournamentRowToProto(tournament)

	// Convert JSONB motions to proto format
	if tournament.Motions.Valid {
		motions, err := motionsJSONToProto(string(tournament.Motions.RawMessage))
		if err != nil {
			return nil, fmt.Errorf("failed to convert tournament motions: %v", err)
		}
		protoTournament.Motions = motions
	}

	// Generate presigned URL for the image
	if tournament.Imageurl.Valid && tournament.Imageurl.String != "" {
		s3Client, err := utils.NewS3Client()
		if err != nil {
			return nil, fmt.Errorf("failed to create S3 client: %v", err)
		}
		key := utils.ExtractKeyFromURL(tournament.Imageurl.String)
		presignedURL, err := s3Client.GetSignedURL(ctx, key, time.Hour)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL: %v", err)
		}
		protoTournament.ImageUrl = presignedURL
	}

	return protoTournament, nil
}

func (s *TournamentService) ListTournaments(ctx context.Context, req *tournament_management.ListTournamentsRequest) (*tournament_management.ListTournamentsResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}
	queries := models.New(s.db)

	tournaments, err := queries.ListTournamentsPaginated(ctx, models.ListTournamentsPaginatedParams{
		Limit:       int32(req.GetPageSize()),
		Offset:      int32(req.GetPageToken()),
		SearchQuery: req.GetSearchQuery(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tournaments: %v", err)
	}

	tournamentResponses := make([]*tournament_management.Tournament, len(tournaments))
	s3Client, err := utils.NewS3Client()
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %v", err)
	}

	for i, tournament := range tournaments {
		tournamentResponses[i] = tournamentPaginatedRowToProto(tournament)
		if tournament.Imageurl.Valid && tournament.Imageurl.String != "" {
			key := utils.ExtractKeyFromURL(tournament.Imageurl.String)
			presignedURL, err := s3Client.GetSignedURL(ctx, key, time.Hour)
			if err != nil {
				return nil, fmt.Errorf("failed to generate presigned URL: %v", err)
			}
			tournamentResponses[i].ImageUrl = presignedURL
		}
	}

	return &tournament_management.ListTournamentsResponse{
		Tournaments:   tournamentResponses,
		NextPageToken: int32(req.GetPageToken()) + int32(req.GetPageSize()),
	}, nil
}

func (s *TournamentService) GetTournamentStats(ctx context.Context, req *tournament_management.GetTournamentStatsRequest) (*tournament_management.GetTournamentStatsResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	// Extract user role and ID from token
	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, fmt.Errorf("user_role not found in token")
	}

	queries := models.New(s.db)
	var stats models.GetTournamentStatsRow

	params := models.GetTournamentStatsParams{
		UserRole: userRole,
	}

	if userRole == "admin" {
		// For admin, only pass the role
		params.SchoolID = sql.NullInt32{Valid: false}
	} else {
		// For school users, get their school ID
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("user_id not found in token")
		}

		schoolInfo, err := queries.GetSchoolByUserID(ctx, int32(userIDFloat))
		if err != nil {
			return nil, fmt.Errorf("failed to get school ID by user ID: %v", err)
		}

		params.SchoolID = sql.NullInt32{Int32: schoolInfo.Schoolid, Valid: true}
	}

	stats, err = queries.GetTournamentStats(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament stats: %v", err)
	}

	// Convert interface{} to int64 for calculations
	yesterdayActiveDebaters, ok := stats.YesterdayActiveDebatersCount.(int64)
	if !ok {
		// Handle the case where it might be int32 or another numeric type
		if val, ok := stats.YesterdayActiveDebatersCount.(int32); ok {
			yesterdayActiveDebaters = int64(val)
		} else {
			return nil, fmt.Errorf("unexpected type for YesterdayActiveDebatersCount")
		}
	}

	activeDebaters, ok := stats.ActiveDebaterCount.(int64)
	if !ok {
		// Handle the case where it might be int32 or another numeric type
		if val, ok := stats.ActiveDebaterCount.(int32); ok {
			activeDebaters = int64(val)
		} else {
			return nil, fmt.Errorf("unexpected type for ActiveDebaterCount")
		}
	}

	totalPercentageChange := calculatePercentageChange(int64(stats.YesterdayTotalCount.Int32), stats.TotalTournaments)
	upcomingPercentageChange := calculatePercentageChange(int64(stats.YesterdayUpcomingCount.Int32), stats.UpcomingTournaments)
	activeDebatersPercentageChange := calculatePercentageChange(yesterdayActiveDebaters, activeDebaters)

	return &tournament_management.GetTournamentStatsResponse{
		TotalTournaments:               int32(stats.TotalTournaments),
		UpcomingTournaments:            int32(stats.UpcomingTournaments),
		TotalPercentageChange:          totalPercentageChange,
		UpcomingPercentageChange:       upcomingPercentageChange,
		ActiveDebaters:                 int32(activeDebaters),
		ActiveDebatersPercentageChange: activeDebatersPercentageChange,
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

	// Validate motions
	if err := validateMotions(req); err != nil {
		return nil, fmt.Errorf("invalid motions: %v", err)
	}

	startDate, err := time.Parse("2006-01-02 15:04", req.GetStartDate())
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02 15:04", req.GetEndDate())
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	// Create motions JSON structure
	motions := map[string]interface{}{
		"preliminary": make([]map[string]interface{}, 0),
		"elimination": make([]map[string]interface{}, 0),
	}

	// Add preliminary motions
	for _, m := range req.GetMotions().GetPreliminaryMotions() {
		motions["preliminary"] = append(motions["preliminary"].([]map[string]interface{}), map[string]interface{}{
			"text":        m.GetText(),
			"infoSlide":   m.GetInfoSlide(),
			"roundNumber": m.GetRoundNumber(),
		})
	}

	// Add elimination motions
	for _, m := range req.GetMotions().GetEliminationMotions() {
		motions["elimination"] = append(motions["elimination"].([]map[string]interface{}), map[string]interface{}{
			"text":        m.GetText(),
			"infoSlide":   m.GetInfoSlide(),
			"roundNumber": m.GetRoundNumber(),
		})
	}

	motionsJSON, err := json.Marshal(motions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal motions: %v", err)
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
		Motions:                    pqtype.NullRawMessage{RawMessage: motionsJSON, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("pairings have already been made")
		}
		return nil, fmt.Errorf("failed to update tournament details: %v", err)
	}

	// Check if ErrorMessage is not nil and is a string
	if errMsg, ok := result.ErrorMessage.(string); ok && errMsg != "" {
		return nil, errors.New(errMsg)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	updatedTournament := tournamentModelToProto(models.Tournament{
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
		Motions:                    result.Motions,
	})

	// Parse the saved motions JSON back into the proto structure
	if result.Motions.Valid {
		var savedMotions map[string]interface{}
		if err := json.Unmarshal(result.Motions.RawMessage, &savedMotions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal saved motions: %v", err)
		}
		updatedTournament.Motions = req.GetMotions()
	}

	// Generate presigned URL for the updated image
	if result.Imageurl.Valid && result.Imageurl.String != "" {
		s3Client, err := utils.NewS3Client()
		if err != nil {
			return nil, fmt.Errorf("failed to create S3 client: %v", err)
		}
		key := utils.ExtractKeyFromURL(result.Imageurl.String)
		presignedURL, err := s3Client.GetSignedURL(ctx, key, time.Hour)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL: %v", err)
		}
		updatedTournament.ImageUrl = presignedURL
	}

	return updatedTournament, nil
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

func (s *TournamentService) SearchTournaments(ctx context.Context, req *tournament_management.SearchTournamentsRequest) (*tournament_management.SearchTournamentsResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	// Clean and prepare the search query
	searchQuery := strings.TrimSpace(req.GetQuery())
	if searchQuery == "" {
		return &tournament_management.SearchTournamentsResponse{
			Tournaments: []*tournament_management.TournamentSearchResult{},
		}, nil
	}

	// Add wildcards for partial matching
	searchPattern := "%" + searchQuery + "%"

	// Execute the search query
	queries := models.New(s.db)
	tournaments, err := queries.SearchTournaments(ctx, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search tournaments: %v", err)
	}

	// Convert the results to proto message
	results := make([]*tournament_management.TournamentSearchResult, len(tournaments))
	for i, t := range tournaments {
		results[i] = &tournament_management.TournamentSearchResult{
			TournamentId: t.Tournamentid,
			Name:         t.Name,
		}
	}

	return &tournament_management.SearchTournamentsResponse{
		Tournaments: results,
	}, nil
}

func (s *TournamentService) SendInvitations(ctx context.Context, req *tournament_management.SendInvitationsRequest) (*tournament_management.SendInvitationsResponse, error) {
	_, err := s.validateAdminRole(req.Token)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	tournament, err := queries.GetTournamentByID(ctx, req.TournamentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament: %v", err)
	}

	league, err := queries.GetLeagueByID(ctx, tournament.Leagueid.Int32)
	if err != nil {
		return nil, fmt.Errorf("failed to get league: %v", err)
	}

	failedUserIDs := []int32{}
	for _, userID := range req.UserIds {
		userDetails, err := queries.GetUserDetailsForInvitation(ctx, userID)
		if err != nil {
			failedUserIDs = append(failedUserIDs, userID)
			continue
		}

		isDAC := strings.ToUpper(league.Name) == "DAC"
		if isDAC {
			if userDetails.Userrole != "volunteer" && userDetails.Userrole != "student" {
				failedUserIDs = append(failedUserIDs, userID)
				continue
			}
		} else {
			if userDetails.Userrole != "volunteer" && userDetails.Userrole != "school" {
				failedUserIDs = append(failedUserIDs, userID)
				continue
			}
		}

		_, err = queries.CreateInvitation(ctx, models.CreateInvitationParams{
			Tournamentid: req.TournamentId,
			Inviteeid:    userDetails.Idebateid.(string),
			Inviteerole:  userDetails.Userrole,
			Status:       "pending",
		})
		if err != nil {
			failedUserIDs = append(failedUserIDs, userID)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Send notifications asynchronously using existing flow
	go func() {
		bgCtx := context.Background()
		notificationService, err := notifications.NewNotificationService(s.db)
		if err != nil {
			log.Printf("Failed to create notification service: %v", err)
			return
		}
		bgQueries := models.New(s.db)

		format, err := bgQueries.GetTournamentFormatByID(bgCtx, tournament.Formatid)
		if err != nil {
			log.Printf("Failed to get tournament format: %v", err)
			return
		}

		// Convert GetTournamentByIDRow to Tournament
		tournamentModel := models.Tournament{
			Tournamentid:               tournament.Tournamentid,
			Name:                       tournament.Name,
			Startdate:                  tournament.Startdate,
			Enddate:                    tournament.Enddate,
			Location:                   tournament.Location,
			Formatid:                   tournament.Formatid,
			Leagueid:                   tournament.Leagueid,
			Coordinatorid:              tournament.Coordinatorid,
			Numberofpreliminaryrounds:  tournament.Numberofpreliminaryrounds,
			Numberofeliminationrounds:  tournament.Numberofeliminationrounds,
			Judgesperdebatepreliminary: tournament.Judgesperdebatepreliminary,
			Judgesperdebateelimination: tournament.Judgesperdebateelimination,
			Tournamentfee:              tournament.Tournamentfee,
			Imageurl:                   tournament.Imageurl,
		}

		err = notification.SendTournamentInvitations(bgCtx, notificationService, tournamentModel, league, format, bgQueries)
		if err != nil {
			log.Printf("Failed to send tournament invitations: %v", err)
		}
	}()

	return &tournament_management.SendInvitationsResponse{
		Success:       len(failedUserIDs) < len(req.UserIds),
		Message:       fmt.Sprintf("Successfully sent %d invitations", len(req.UserIds)-len(failedUserIDs)),
		FailedUserIds: failedUserIDs,
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

func validateMotions(req interface{}) error {
	var preliminaryMotions []*tournament_management.Motion
	var eliminationMotions []*tournament_management.Motion
	var numPreliminaryRounds int32
	var numEliminationRounds int32

	switch r := req.(type) {
	case *tournament_management.CreateTournamentRequest:
		preliminaryMotions = r.GetMotions().GetPreliminaryMotions()
		eliminationMotions = r.GetMotions().GetEliminationMotions()
		numPreliminaryRounds = r.GetNumberOfPreliminaryRounds()
		numEliminationRounds = r.GetNumberOfEliminationRounds()
	case *tournament_management.UpdateTournamentRequest:
		preliminaryMotions = r.GetMotions().GetPreliminaryMotions()
		eliminationMotions = r.GetMotions().GetEliminationMotions()
		numPreliminaryRounds = r.GetNumberOfPreliminaryRounds()
		numEliminationRounds = r.GetNumberOfEliminationRounds()
	default:
		return fmt.Errorf("unsupported request type")
	}

	// Check if we have the correct number of motions
	if len(preliminaryMotions) != int(numPreliminaryRounds) {
		return fmt.Errorf("number of preliminary motions (%d) does not match number of preliminary rounds (%d)",
			len(preliminaryMotions), numPreliminaryRounds)
	}

	if len(eliminationMotions) != int(numEliminationRounds) {
		return fmt.Errorf("number of elimination motions (%d) does not match number of elimination rounds (%d)",
			len(eliminationMotions), numEliminationRounds)
	}

	// Validate each preliminary motion
	seenPrelimRounds := make(map[int32]bool)
	for i, motion := range preliminaryMotions {
		if motion.GetText() == "" {
			return fmt.Errorf("preliminary motion %d has empty text", i+1)
		}

		roundNum := motion.GetRoundNumber()
		if roundNum < 1 || roundNum > numPreliminaryRounds {
			return fmt.Errorf("invalid preliminary round number: %d (should be between 1 and %d)",
				roundNum, numPreliminaryRounds)
		}

		if seenPrelimRounds[roundNum] {
			return fmt.Errorf("duplicate preliminary round number: %d", roundNum)
		}
		seenPrelimRounds[roundNum] = true
	}

	// Validate each elimination motion
	seenElimRounds := make(map[int32]bool)
	for i, motion := range eliminationMotions {
		if motion.GetText() == "" {
			return fmt.Errorf("elimination motion %d has empty text", i+1)
		}

		roundNum := motion.GetRoundNumber()
		if roundNum < 1 || roundNum > numEliminationRounds {
			return fmt.Errorf("invalid elimination round number: %d (should be between 1 and %d)",
				roundNum, numEliminationRounds)
		}

		if seenElimRounds[roundNum] {
			return fmt.Errorf("duplicate elimination round number: %d", roundNum)
		}
		seenElimRounds[roundNum] = true

		// Optional: validate info slide (if it's required)
		if motion.GetInfoSlide() == "" {
			// Uncomment if info slide is required
			// return fmt.Errorf("elimination motion %d has empty info slide", i+1)
		}
	}

	return nil
}

// Helper function to convert JSONB motion data to proto
func motionsJSONToProto(motionsJSON string) (*tournament_management.TournamentMotions, error) {
	var rawMotions map[string]interface{}
	if err := json.Unmarshal([]byte(motionsJSON), &rawMotions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal motions JSON: %v", err)
	}

	protoMotions := &tournament_management.TournamentMotions{
		PreliminaryMotions: make([]*tournament_management.Motion, 0),
		EliminationMotions: make([]*tournament_management.Motion, 0),
	}

	// Convert preliminary motions
	if prelim, ok := rawMotions["preliminary"].([]interface{}); ok {
		for _, m := range prelim {
			motion, err := convertMotion(m)
			if err != nil {
				return nil, err
			}
			protoMotions.PreliminaryMotions = append(protoMotions.PreliminaryMotions, motion)
		}
	}

	// Convert elimination motions
	if elim, ok := rawMotions["elimination"].([]interface{}); ok {
		for _, m := range elim {
			motion, err := convertMotion(m)
			if err != nil {
				return nil, err
			}
			protoMotions.EliminationMotions = append(protoMotions.EliminationMotions, motion)
		}
	}

	return protoMotions, nil
}

// Helper function to convert a single motion from JSON to proto
func convertMotion(m interface{}) (*tournament_management.Motion, error) {
	motionMap, ok := m.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid motion format")
	}

	motion := &tournament_management.Motion{}

	if text, ok := motionMap["text"].(string); ok {
		motion.Text = text
	}
	if infoSlide, ok := motionMap["infoSlide"].(string); ok {
		motion.InfoSlide = infoSlide
	}
	if roundNum, ok := motionMap["roundNumber"].(float64); ok {
		motion.RoundNumber = int32(roundNum)
	}

	return motion, nil
}

func parseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
