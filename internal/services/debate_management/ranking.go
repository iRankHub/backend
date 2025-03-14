package services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iRankHub/backend/internal/grpc/proto/debate_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type RankingService struct {
	db *sql.DB
}

func NewRankingService(db *sql.DB) *RankingService {
	return &RankingService{db: db}
}

func (s *RankingService) GetTournamentStudentRanking(ctx context.Context, req *debate_management.TournamentRankingRequest) (*debate_management.TournamentRankingResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user role in token")
	}

	// Check visibility for non-admin users
	if err := s.checkRankingVisibility(ctx, req.GetTournamentId(), "student", userRole); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	params := models.GetTournamentStudentRankingParams{
		Tournamentid: req.GetTournamentId(),
		Limit:        int32(req.GetPageSize()),
		Offset:       int32((req.GetPage() - 1) * req.GetPageSize()),
	}

	if req.GetSearch() != "" {
		params.Column4 = req.GetSearch()
	}

	dbRankings, err := queries.GetTournamentStudentRanking(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tournament student ranking: %v", err)
	}

	rankings := make([]*debate_management.StudentRanking, len(dbRankings))
	for i, dbRanking := range dbRankings {
		totalPoints, err := strconv.ParseFloat(dbRanking.Totalpoints, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse total points: %v", err)
		}

		rankings[i] = &debate_management.StudentRanking{
			StudentId:   dbRanking.Studentid,
			StudentName: dbRanking.Studentname.(string),
			SchoolName:  dbRanking.Schoolname,
			TotalWins:   int32(dbRanking.Wins),
			TotalPoints: totalPoints,
			AverageRank: float64(dbRanking.Averagerank),
			Place:       int32(dbRanking.Place),
		}
	}

	return &debate_management.TournamentRankingResponse{
		Rankings: rankings,
	}, nil
}

func (s *RankingService) GetOverallStudentRanking(ctx context.Context, req *debate_management.OverallRankingRequest) (*debate_management.OverallRankingResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user role in token")
	}

	if int32(userID) != req.GetUserId() && userRole != "admin" {
		return nil, status.Error(codes.PermissionDenied, "unauthorized access to student ranking")
	}

	queries := models.New(s.db)
	studentID, err := s.getUserStudentID(ctx, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to get student ID: %v", err)
	}

	dbRankings, err := queries.GetOverallStudentRanking(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get overall student ranking: %v", err)
	}

	var studentRanking *models.GetOverallStudentRankingRow
	var topStudents []*debate_management.TopStudent
	rankChanges := make(map[int32]int32)

	for i, ranking := range dbRankings {
		if ranking.Studentid == studentID {
			studentRanking = &dbRankings[i]
		}

		totalPoints, err := strconv.ParseFloat(ranking.Totalpoints, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse total points: %v", err)
		}

		if ranking.Currentrank <= 3 {
			topStudents = append(topStudents, &debate_management.TopStudent{
				Rank:        int32(ranking.Currentrank),
				Name:        ranking.Studentname.(string),
				TotalPoints: totalPoints,
			})
		}

		// Calculate rank changes
		if i > 0 {
			currentDate, ok := ranking.Lasttournamentdate.(time.Time)
			previousDate, prevOk := dbRankings[i-1].Lasttournamentdate.(time.Time)

			if ok && prevOk && currentDate.Before(previousDate) {
				rankChanges[ranking.Studentid] = int32(dbRankings[i-1].Currentrank - ranking.Currentrank)
			}
		}
	}

	if studentRanking == nil {
		return nil, fmt.Errorf("student not found in rankings")
	}

	studentTotalPoints, err := strconv.ParseFloat(studentRanking.Totalpoints, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse student total points: %v", err)
	}

	response := &debate_management.OverallRankingResponse{
		StudentRank:   int32(studentRanking.Currentrank),
		TotalStudents: int32(studentRanking.Totalstudents),
		RankChange:    rankChanges[studentID],
		TopStudents:   topStudents,
		StudentInfo: &debate_management.StudentInfo{
			Name:        studentRanking.Studentname.(string),
			TotalPoints: studentTotalPoints,
		},
	}

	return response, nil
}

func (s *RankingService) GetStudentOverallPerformance(ctx context.Context, req *debate_management.PerformanceRequest) (*debate_management.PerformanceResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user role in token")
	}

	if int32(userID) != req.GetUserId() && userRole != "admin" {
		return nil, status.Error(codes.PermissionDenied, "unauthorized access to student performance")
	}

	queries := models.New(s.db)
	studentID, err := s.getUserStudentID(ctx, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to get student ID: %v", err)
	}

	startDate, err := time.Parse("2006-01-02", req.GetStartDate())
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", req.GetEndDate())
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	dbPerformance, err := queries.GetStudentOverallPerformance(ctx, models.GetStudentOverallPerformanceParams{
		Studentid:   studentID,
		Startdate:   startDate,
		Startdate_2: endDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get student overall performance: %v", err)
	}

	performanceData := make([]*debate_management.PerformanceData, len(dbPerformance))
	for i, data := range dbPerformance {
		studentTotalPoints, err := strconv.ParseFloat(data.Studenttotalpoints, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse student total points: %v", err)
		}
		studentAveragePoints, err := strconv.ParseFloat(data.Studentaveragepoints, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse student average points: %v", err)
		}

		performanceData[i] = &debate_management.PerformanceData{
			TournamentDate:       data.Startdate.Format("2006-01-02"),
			StudentTotalPoints:   studentTotalPoints,
			StudentAveragePoints: studentAveragePoints,
			TournamentRank:       int32(data.Tournamentrank),
		}
	}

	return &debate_management.PerformanceResponse{
		PerformanceData: performanceData,
	}, nil
}

func (s *RankingService) GetStudentTournamentStats(ctx context.Context, req *debate_management.StudentTournamentStatsRequest) (*debate_management.StudentTournamentStatsResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user role in token")
	}

	if int32(userID) != req.GetStudentId() && userRole != "admin" {
		return nil, status.Error(codes.PermissionDenied, "unauthorized access to student tournament stats")
	}

	queries := models.New(s.db)
	stats, err := queries.GetStudentTournamentStats(ctx, req.GetStudentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get student tournament stats: %v", err)
	}

	totalPercentageChange := calculatePercentageChange(int64(stats.YesterdayTotalCount.Int32), int64(stats.TotalTournaments))
	upcomingPercentageChange := calculatePercentageChange(int64(stats.YesterdayUpcomingCount.Int32), int64(stats.UpcomingTournaments))

	attendedPercentageChange := "0.00%"
	if stats.TotalTournamentsLastYear > 0 {
		previousAttended := float64(stats.TotalTournamentsLastYear) - float64(stats.AttendedTournaments)
		if previousAttended > 0 {
			change := (float64(stats.AttendedTournaments) - previousAttended) / previousAttended * 100
			attendedPercentageChange = formatPercentageChange(change)
		}
	}

	return &debate_management.StudentTournamentStatsResponse{
		TotalTournaments:          int32(stats.TotalTournaments),
		TotalTournamentsChange:    totalPercentageChange,
		AttendedTournaments:       int32(stats.AttendedTournaments),
		AttendedTournamentsChange: attendedPercentageChange,
		UpcomingTournaments:       int32(stats.UpcomingTournaments),
		UpcomingTournamentsChange: upcomingPercentageChange,
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
	return formatPercentageChange(change)
}

func formatPercentageChange(change float64) string {
	sign := "+"
	if change < 0 {
		sign = "-"
		change = math.Abs(change)
	}
	return fmt.Sprintf("%s%.2f%%", sign, change)
}

func (s *RankingService) GetTournamentTeamsRanking(ctx context.Context, req *debate_management.TournamentTeamsRankingRequest) (*debate_management.TournamentTeamsRankingResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	// Get user role and check visibility
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user role in token")
	}

	if err := s.checkRankingVisibility(ctx, req.GetTournamentId(), "team", userRole); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	params := models.GetTournamentTeamsRankingParams{
		Tournamentid: req.GetTournamentId(),
		Limit:        int32(req.GetPageSize()),
		Offset:       int32((req.GetPage() - 1) * req.GetPageSize()),
	}

	if req.GetSearch() != "" {
		params.Column4 = req.GetSearch()
	}

	dbRankings, err := queries.GetTournamentTeamsRanking(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament teams ranking: %v", err)
	}

	totalCount, err := queries.GetTournamentTeamsRankingCount(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	rankings := make([]*debate_management.TeamRanking, len(dbRankings))
	for i, dbRanking := range dbRankings {
		// Proper handling of SchoolNames array
		var schoolNames []string

		switch v := dbRanking.Schoolnames.(type) {
		case []string:
			schoolNames = v
		case string:
			// Handle string format, might be "{school1,school2}" or "school1,school2"
			s := strings.TrimPrefix(v, "{")
			s = strings.TrimSuffix(s, "}")
			if s != "" {
				schoolNames = strings.Split(s, ",")
			}
		case []interface{}:
			// Handle array of interfaces
			for _, item := range v {
				if str, ok := item.(string); ok {
					schoolNames = append(schoolNames, str)
				}
			}
		default:
			// If we can't determine the type, log it and use empty array
			fmt.Printf("WARNING: Unknown SchoolNames type: %T for team %d\n",
				dbRanking.Schoolnames, dbRanking.Teamid)
			schoolNames = []string{}
		}

		// Clean up school names to remove any quotes or extra spaces
		for j, name := range schoolNames {
			schoolNames[j] = strings.Trim(name, "\" \t")
		}

		// Proper handling of numeric values
		totalPoints, err := convertToFloat64(dbRanking.Totalpoints)
		if err != nil {
			return nil, fmt.Errorf("failed to parse total points for team %d: %v",
				dbRanking.Teamid, err)
		}

		averageRank, err := convertToFloat64(dbRanking.Averagerank)
		if err != nil {
			return nil, fmt.Errorf("failed to parse average rank for team %d: %v",
				dbRanking.Teamid, err)
		}

		// Use sequential position (i+1) instead of the place from database
		sequentialPlace := int32(i + 1)

		rankings[i] = &debate_management.TeamRanking{
			TeamId:      dbRanking.Teamid,
			TeamName:    strings.Trim(dbRanking.Teamname, "\" \t"),
			SchoolNames: schoolNames,
			Wins:        int32(dbRanking.Wins),
			TotalPoints: totalPoints,
			AverageRank: averageRank,
			Place:       sequentialPlace,
		}
	}

	return &debate_management.TournamentTeamsRankingResponse{
		Rankings:   rankings,
		TotalCount: int32(totalCount),
	}, nil
}

// Helper function to convert interface{} to float64
func convertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case sql.NullString:
		if !v.Valid {
			return 0, nil
		}
		return strconv.ParseFloat(v.String, 64)
	case sql.NullFloat64:
		if !v.Valid {
			return 0, nil
		}
		return v.Float64, nil
	case sql.NullInt64:
		if !v.Valid {
			return 0, nil
		}
		return float64(v.Int64), nil
	default:
		return 0, fmt.Errorf("unexpected type: %T", value)
	}
}

func (s *RankingService) GetTournamentSchoolRanking(ctx context.Context, req *debate_management.TournamentSchoolRankingRequest) (*debate_management.TournamentSchoolRankingResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user role in token")
	}

	if err := s.checkRankingVisibility(ctx, req.GetTournamentId(), "school", userRole); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	params := models.GetTournamentSchoolRankingParams{
		Tournamentid: req.GetTournamentId(),
		Limit:        int32(req.GetPageSize()),
		Offset:       int32((req.GetPage() - 1) * req.GetPageSize()),
	}

	if req.GetSearch() != "" {
		params.Column4 = req.GetSearch()
	}

	dbRankings, err := queries.GetTournamentSchoolRanking(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament school ranking: %v", err)
	}

	totalCount, err := queries.GetTournamentSchoolRankingCount(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	rankings := make([]*debate_management.SchoolRanking, len(dbRankings))
	for i, dbRanking := range dbRankings {
		totalPoints, err := convertToFloat64(dbRanking.Totalpoints)
		if err != nil {
			return nil, fmt.Errorf("failed to parse total points: %v", err)
		}
		averageRank, err := convertToFloat64(dbRanking.Averagerank)
		if err != nil {
			return nil, fmt.Errorf("failed to parse average rank: %v", err)
		}

		rankings[i] = &debate_management.SchoolRanking{
			SchoolName:  dbRanking.Schoolname,
			TeamCount:   int32(dbRanking.Teamcount),
			TotalWins:   int32(dbRanking.Totalwins),
			AverageRank: averageRank,
			TotalPoints: totalPoints,
			Place:       int32(dbRanking.Place),
		}
	}

	return &debate_management.TournamentSchoolRankingResponse{
		Rankings:   rankings,
		TotalCount: int32(totalCount),
	}, nil
}

func (s *RankingService) GetOverallSchoolRanking(ctx context.Context, req *debate_management.OverallSchoolRankingRequest) (*debate_management.OverallSchoolRankingResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user role in token")
	}

	if int32(userID) != req.GetUserId() && userRole != "admin" {
		return nil, status.Error(codes.PermissionDenied, "unauthorized access to school ranking")
	}

	queries := models.New(s.db)
	schoolID, err := queries.GetSchoolIDByUserID(ctx, req.GetUserId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user is not associated with any school")
		}
		return nil, fmt.Errorf("failed to get school ID: %v", err)
	}

	dbRankings, err := queries.GetOverallSchoolRanking(ctx, schoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overall school ranking: %v", err)
	}

	var schoolRanking *models.GetOverallSchoolRankingRow
	var topSchools []*debate_management.TopSchool
	rankChanges := make(map[int32]int32)

	for i, ranking := range dbRankings {
		if ranking.Schoolid == schoolID {
			schoolRanking = &dbRankings[i]
		}

		totalPoints, err := strconv.ParseFloat(ranking.Totalpoints, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse total points: %v", err)
		}

		if ranking.Currentrank <= 3 {
			topSchools = append(topSchools, &debate_management.TopSchool{
				Rank:        int32(ranking.Currentrank),
				Name:        ranking.Schoolname,
				TotalPoints: totalPoints,
			})
		}

		// Calculate rank changes
		if i > 0 {
			currentDate, ok := ranking.Lasttournamentdate.(time.Time)
			previousDate, prevOk := dbRankings[i-1].Lasttournamentdate.(time.Time)

			if ok && prevOk && currentDate.Before(previousDate) {
				rankChanges[ranking.Schoolid] = int32(dbRankings[i-1].Currentrank - ranking.Currentrank)
			}
		}
	}

	if schoolRanking == nil {
		return nil, status.Errorf(codes.NotFound, "school not found in rankings")
	}

	schoolTotalPoints, err := strconv.ParseFloat(schoolRanking.Totalpoints, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse school total points: %v", err)
	}

	response := &debate_management.OverallSchoolRankingResponse{
		SchoolRank:   int32(schoolRanking.Currentrank),
		TotalSchools: int32(schoolRanking.Totalschools),
		RankChange:   rankChanges[schoolID],
		TopSchools:   topSchools,
		SchoolInfo: &debate_management.SchoolInfo{
			Name:        schoolRanking.Schoolname,
			TotalPoints: schoolTotalPoints,
		},
	}

	return response, nil
}

func (s *RankingService) GetSchoolOverallPerformance(ctx context.Context, req *debate_management.SchoolPerformanceRequest) (*debate_management.SchoolPerformanceResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user role in token")
	}

	if int32(userID) != req.GetUserId() && userRole != "admin" {
		return nil, status.Error(codes.PermissionDenied, "unauthorized access to school performance")
	}

	queries := models.New(s.db)
	schoolID, err := queries.GetSchoolIDByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to get school ID: %v", err)
	}

	startDate, err := time.Parse("2006-01-02", req.GetStartDate())
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", req.GetEndDate())
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	dbPerformance, err := queries.GetSchoolOverallPerformance(ctx, models.GetSchoolOverallPerformanceParams{
		Schoolid:    schoolID,
		Startdate:   startDate,
		Startdate_2: endDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get school overall performance: %v", err)
	}

	performanceData := make([]*debate_management.SchoolPerformanceData, len(dbPerformance))
	for i, data := range dbPerformance {
		schoolTotalPoints, err := convertToFloat64(data.Schooltotalpoints)
		if err != nil {
			return nil, fmt.Errorf("failed to convert school total points: %v", err)
		}

		schoolAveragePoints, err := convertToFloat64(data.Schoolaveragepoints)
		if err != nil {
			return nil, fmt.Errorf("failed to convert school average points: %v", err)
		}

		performanceData[i] = &debate_management.SchoolPerformanceData{
			TournamentDate:      data.Startdate.Format("2006-01-02"),
			SchoolTotalPoints:   schoolTotalPoints,
			SchoolAveragePoints: schoolAveragePoints,
			TournamentRank:      int32(data.Tournamentrank),
		}
	}

	return &debate_management.SchoolPerformanceResponse{
		PerformanceData: performanceData,
	}, nil
}

func (s *RankingService) GetVolunteerTournamentStats(ctx context.Context, req *debate_management.VolunteerTournamentStatsRequest) (*debate_management.VolunteerTournamentStatsResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	queries := models.New(s.db)
	stats, err := queries.GetVolunteerTournamentStats(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get volunteer tournament stats: %v", err)
	}

	roundsJudgedChange := calculatePercentageChange(
		int64(stats.YesterdayRoundsJudged),
		int64(stats.TotalRoundsJudged),
	)
	tournamentsAttendedChange := calculatePercentageChange(
		int64(stats.YesterdayTournamentsAttended),
		int64(stats.AttendedTournaments),
	)
	upcomingTournamentsChange := calculatePercentageChange(
		int64(stats.YesterdayUpcomingTournaments),
		int64(stats.UpcomingTournaments),
	)

	return &debate_management.VolunteerTournamentStatsResponse{
		TotalRoundsJudged:         int32(stats.TotalRoundsJudged),
		RoundsJudgedChange:        roundsJudgedChange,
		TournamentsAttended:       int32(stats.AttendedTournaments),
		TournamentsAttendedChange: tournamentsAttendedChange,
		UpcomingTournaments:       int32(stats.UpcomingTournaments),
		UpcomingTournamentsChange: upcomingTournamentsChange,
	}, nil
}

func (s *RankingService) GetVolunteerRanking(ctx context.Context, req *debate_management.GetVolunteerRankingRequest) (*debate_management.GetVolunteerRankingResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	queries := models.New(s.db)

	// First get volunteer ID from user ID
	volunteer, err := queries.GetVolunteerByUserID(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get volunteer: %v", err)
	}

	dbRankings, err := queries.GetOverallVolunteerRanking(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get volunteer ranking: %v", err)
	}

	var volunteerRanking *models.GetOverallVolunteerRankingRow
	var topVolunteers []*debate_management.TopVolunteer
	rankChanges := make(map[int32]int32)

	// Process rankings
	for i, ranking := range dbRankings {
		if ranking.Volunteerid == volunteer.Volunteerid {
			volunteerRanking = &dbRankings[i]
		}

		// Add top 3 volunteers
		if ranking.Currentrank <= 3 {
			topVolunteers = append(topVolunteers, &debate_management.TopVolunteer{
				Rank:          int32(ranking.Currentrank),
				Name:          ranking.Volunteername.(string),
				AverageRating: float64(ranking.Averagerating),
			})
		}

		// Calculate rank changes
		if i > 0 {
			currentDate, ok := ranking.Lasttournamentdate.(time.Time)
			previousDate, prevOk := dbRankings[i-1].Lasttournamentdate.(time.Time)

			if ok && prevOk && currentDate.Before(previousDate) {
				rankChanges[ranking.Volunteerid] = int32(dbRankings[i-1].Currentrank - ranking.Currentrank)
			}
		}
	}

	if volunteerRanking == nil {
		return nil, fmt.Errorf("volunteer not found in rankings")
	}

	return &debate_management.GetVolunteerRankingResponse{
		VolunteerRank:   int32(volunteerRanking.Currentrank),
		TotalVolunteers: int32(volunteerRanking.Totalvolunteers),
		RankChange:      rankChanges[volunteer.Volunteerid],
		TopVolunteers:   topVolunteers,
		VolunteerInfo: &debate_management.VolunteerInfo{
			Name:          volunteerRanking.Volunteername.(string),
			AverageRating: float64(volunteerRanking.Averagerating),
		},
	}, nil
}

func (s *RankingService) GetVolunteerPerformance(ctx context.Context, req *debate_management.GetVolunteerPerformanceRequest) (*debate_management.GetVolunteerPerformanceResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	startDate, err := time.Parse("2006-01-02", req.GetStartDate())
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", req.GetEndDate())
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	queries := models.New(s.db)
	dbPerformance, err := queries.GetVolunteerPerformance(ctx, models.GetVolunteerPerformanceParams{
		Userid:      int32(userID),
		Startdate:   startDate,
		Startdate_2: endDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get volunteer performance: %v", err)
	}

	performanceData := make([]*debate_management.VolunteerPerformanceData, len(dbPerformance))
	for i, data := range dbPerformance {
		performanceData[i] = &debate_management.VolunteerPerformanceData{
			TournamentDate:         data.Startdate.Format("2006-01-02"),
			VolunteerAverageRating: float64(data.Volunteeraveragerating),
			OverallAverageRating:   float64(data.Overallaveragerating),
			TournamentRank:         int32(data.Tournamentrank),
		}
	}

	return &debate_management.GetVolunteerPerformanceResponse{
		PerformanceData: performanceData,
	}, nil
}

func (s *RankingService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}

func (s *RankingService) GetTournamentVolunteerRanking(ctx context.Context, req *debate_management.TournamentVolunteerRankingRequest) (*debate_management.TournamentVolunteerRankingResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user role in token")
	}

	if err := s.checkRankingVisibility(ctx, req.GetTournamentId(), "volunteer", userRole); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	params := models.GetTournamentVolunteerRankingParams{
		Tournamentid: req.GetTournamentId(),
		Limit:        int32(req.GetPageSize()),
		Offset:       int32((req.GetPage() - 1) * req.GetPageSize()),
	}

	if req.GetSearch() != "" {
		params.Column4 = req.GetSearch()
	}

	dbRankings, err := queries.GetTournamentVolunteerRanking(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament volunteer ranking: %v", err)
	}

	totalCount, err := queries.GetTournamentVolunteerRankingCount(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	rankings := make([]*debate_management.VolunteerTournamentRank, len(dbRankings))
	for i, dbRanking := range dbRankings {
		// Convert average rating to float64
		averageRating, err := convertToFloat64(dbRanking.Averagerating)
		if err != nil {
			return nil, fmt.Errorf("failed to convert average rating for volunteer %v: %v", dbRanking.Volunteerid, err)
		}

		// Convert volunteer name to string
		volunteerName, ok := dbRanking.Volunteername.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert volunteer name for volunteer %v", dbRanking.Volunteerid)
		}

		rankings[i] = &debate_management.VolunteerTournamentRank{
			VolunteerId:       dbRanking.Volunteerid,
			VolunteerName:     volunteerName,
			AverageRating:     averageRating,
			PreliminaryRounds: int32(dbRanking.Preliminaryrounds),
			EliminationRounds: int32(dbRanking.Eliminationrounds),
			Rank:              int32(dbRanking.Place),
			Place:             int32(dbRanking.Place),
		}
	}

	return &debate_management.TournamentVolunteerRankingResponse{
		Rankings:   rankings,
		TotalCount: int32(totalCount),
	}, nil
}

func (s *RankingService) getUserStudentID(ctx context.Context, userID int32) (int32, error) {
	queries := models.New(s.db)
	student, err := queries.GetStudentByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get student: %v", err)
	}
	return student.Studentid, nil
}

func (s *RankingService) SetRankingVisibility(ctx context.Context, req *debate_management.SetRankingVisibilityRequest) (*debate_management.SetRankingVisibilityResponse, error) {
	// Verify admin role
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, status.Error(codes.PermissionDenied, "only admins can modify ranking visibility")
	}

	queries := models.New(s.db)
	err = queries.SetRankingVisibility(ctx, models.SetRankingVisibilityParams{
		Tournamentid: req.GetTournamentId(),
		Rankingtype:  req.GetRankingType(),
		Visibleto:    req.GetVisibleTo(),
		Isvisible:    req.GetIsVisible(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set ranking visibility: %v", err)
	}

	return &debate_management.SetRankingVisibilityResponse{
		Success: true,
		Message: "Ranking visibility updated successfully",
	}, nil
}

// Add this helper method to check visibility
func (s *RankingService) checkRankingVisibility(ctx context.Context, tournamentID int32, rankingType string, userRole string) error {
	if userRole == "admin" {
		return nil
	}

	queries := models.New(s.db)
	isVisible, err := queries.GetRankingVisibility(ctx, models.GetRankingVisibilityParams{
		Tournamentid: tournamentID,
		Rankingtype:  rankingType,
		Visibleto:    userRole,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return status.Error(codes.PermissionDenied, "ranking not visible for your role")
		}
		return status.Errorf(codes.Internal, "failed to check ranking visibility: %v", err)
	}

	if !isVisible {
		return status.Error(codes.PermissionDenied, "ranking not visible for your role")
	}

	return nil
}
