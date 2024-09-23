package services

import (
	"context"
	"database/sql"
	"fmt"
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
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	dbRankings, err := queries.GetTournamentStudentRanking(ctx, models.GetTournamentStudentRankingParams{
		Tournamentid: req.GetTournamentId(),
		Limit:        int32(req.GetPageSize()),
		Offset:       int32((req.GetPage() - 1) * req.GetPageSize()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament student ranking: %v", err)
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
			TotalWins:   int32(dbRanking.Totalwins),
			TotalPoints: totalPoints,
			AverageRank: float64(dbRanking.Averagerank),
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

func (s *RankingService) GetTournamentTeamsRanking(ctx context.Context, req *debate_management.TournamentTeamsRankingRequest) (*debate_management.TournamentTeamsRankingResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	dbRankings, err := queries.GetTournamentTeamsRanking(ctx, models.GetTournamentTeamsRankingParams{
		Tournamentid: req.GetTournamentId(),
		Limit:        int32(req.GetPageSize()),
		Offset:       int32((req.GetPage() - 1) * req.GetPageSize()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament teams ranking: %v", err)
	}

	totalCount, err := queries.GetTournamentTeamsRankingCount(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	rankings := make([]*debate_management.TeamRanking, len(dbRankings))
	for i, dbRanking := range dbRankings {
		// Handle Schoolnames
		schoolNamesStr, ok := dbRanking.Schoolnames.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for Schoolnames for team %s", dbRanking.Teamname)
		}
		schoolNames := strings.Split(schoolNamesStr, ",")

		// Handle Totalpoints
		totalPoints, err := convertToFloat64(dbRanking.Totalpoints)
		if err != nil {
			return nil, fmt.Errorf("failed to parse total points for team %s: %v", dbRanking.Teamname, err)
		}

		// Handle Averagerank
		averageRank, err := convertToFloat64(dbRanking.Averagerank)
		if err != nil {
			return nil, fmt.Errorf("failed to parse average rank for team %s: %v", dbRanking.Teamname, err)
		}

		rankings[i] = &debate_management.TeamRanking{
			TeamId:      dbRanking.Teamid,
			TeamName:    dbRanking.Teamname,
			SchoolNames: schoolNames,
			Wins:        int32(dbRanking.Wins),
			TotalPoints: totalPoints,
			AverageRank: averageRank,
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
	default:
		return 0, fmt.Errorf("unexpected type: %T", value)
	}
}

func (s *RankingService) GetTournamentSchoolRanking(ctx context.Context, req *debate_management.TournamentSchoolRankingRequest) (*debate_management.TournamentSchoolRankingResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	dbRankings, err := queries.GetTournamentSchoolRanking(ctx, models.GetTournamentSchoolRankingParams{
		Tournamentid: req.GetTournamentId(),
		Limit:        int32(req.GetPageSize()),
		Offset:       int32((req.GetPage() - 1) * req.GetPageSize()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament school ranking: %v", err)
	}

	totalCount, err := queries.GetTournamentSchoolRankingCount(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	rankings := make([]*debate_management.SchoolRanking, len(dbRankings))
	for i, dbRanking := range dbRankings {
		totalPoints, err := strconv.ParseFloat(dbRanking.Totalpoints, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse total points for school %s: %v", dbRanking.Schoolname, err)
		}

		rankings[i] = &debate_management.SchoolRanking{
			SchoolName:  dbRanking.Schoolname,
			TeamCount:   int32(dbRanking.Teamcount),
			TotalWins:   int32(dbRanking.Totalwins),
			AverageRank: dbRanking.Averagerank,
			TotalPoints: totalPoints,
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

	dbRankings, err := queries.GetOverallSchoolRanking(ctx)
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

func (s *RankingService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}

func (s *RankingService) getUserStudentID(ctx context.Context, userID int32) (int32, error) {
	queries := models.New(s.db)
	student, err := queries.GetStudentByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get student: %v", err)
	}
	return student.Studentid, nil
}
