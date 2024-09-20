package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

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
	dbRankings, err := queries.GetTournamentStudentRanking(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament student ranking: %v", err)
	}

	rankings := make([]*debate_management.StudentRanking, len(dbRankings))
	for i, dbRanking := range dbRankings {
		rankings[i] = &debate_management.StudentRanking{
			StudentId:    dbRanking.Studentid,
			StudentName:  dbRanking.Studentname.(string),
			SchoolName:   dbRanking.Schoolname,
			TotalWins:    int32(dbRanking.Totalwins),
			TotalPoints:  float64(dbRanking.Totalpoints),
			AverageRank:  dbRanking.Averagerank,
		}
	}

	return &debate_management.TournamentRankingResponse{
		Rankings: rankings,
	}, nil
}

func (s *RankingService) GetOverallStudentRanking(ctx context.Context, req *debate_management.OverallRankingRequest) (*debate_management.OverallRankingResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	studentID, err := s.getUserStudentID(ctx, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to get student ID: %v", err)
	}

	dbRanking, err := queries.GetOverallStudentRanking(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overall student ranking: %v", err)
	}

	var topStudents []struct {
		Rank        int32   `json:"rank"`
		Name        string  `json:"name"`
		TotalPoints float64 `json:"totalPoints"`
		RankChange  int32   `json:"rankChange"`
	}
	err = json.Unmarshal([]byte(dbRanking.Topstudents), &topStudents)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal top students: %v", err)
	}

	response := &debate_management.OverallRankingResponse{
		StudentRank:   int32(dbRanking.Studentrank),
		TotalStudents: int32(dbRanking.Totalstudents),
		RankChange:    dbRanking.Rankchange,
		TopStudents:   make([]*debate_management.TopStudent, len(topStudents)),
		StudentInfo: &debate_management.StudentInfo{
			Name:        dbRanking.Studentname.(string),
			TotalPoints: float64(dbRanking.Totalpoints),
		},
	}

	for i, student := range topStudents {
		response.TopStudents[i] = &debate_management.TopStudent{
			Rank:        student.Rank,
			Name:        student.Name,
			TotalPoints: student.TotalPoints,
			RankChange:  student.RankChange,
		}
	}

	return response, nil
}

func (s *RankingService) GetStudentOverallPerformance(ctx context.Context, req *debate_management.PerformanceRequest) (*debate_management.PerformanceResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	studentID, err := s.getUserStudentID(ctx, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to get student ID: %v", err)
	}

	startDate := req.GetStartDate().AsTime()
	endDate := req.GetEndDate().AsTime()

	dbPerformance, err := queries.GetStudentOverallPerformance(ctx, models.GetStudentOverallPerformanceParams{
		Studentid: studentID,
		Startdate: startDate,
		Startdate_2:   endDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get student overall performance: %v", err)
	}

	performanceData := make([]*debate_management.PerformanceData, len(dbPerformance))
	for i, data := range dbPerformance {
		performanceData[i] = &debate_management.PerformanceData{
			TournamentDate: timestamppb.New(data.Startdate),
			StudentRank:    data.Studentrank,
			AverageRank:    data.Averagerank,
		}
	}

	return &debate_management.PerformanceResponse{
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