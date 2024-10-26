package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iRankHub/backend/internal/grpc/proto/debate_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type FeedbackService struct {
	db *sql.DB
}

func NewFeedbackService(db *sql.DB) *FeedbackService {
	return &FeedbackService{db: db}
}

func (s *FeedbackService) GetStudentFeedback(ctx context.Context, req *debate_management.GetStudentFeedbackRequest) (*debate_management.GetStudentFeedbackResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	queries := models.New(s.db)
	studentID, err := queries.GetStudentByUserID(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get student ID: %v", err)
	}

	feedbacks, err := queries.GetStudentFeedback(ctx, models.GetStudentFeedbackParams{
		Studentid:    studentID.Studentid,
		Tournamentid: req.GetTournamentId(),
		Limit:        int32(req.GetPageSize()),
		Offset:       int32((req.GetPage() - 1) * req.GetPageSize()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get student feedback: %v", err)
	}

	totalCount, err := queries.GetStudentFeedbackCount(ctx, models.GetStudentFeedbackCountParams{
		Speakerid:    studentID.Studentid,
		Tournamentid: req.GetTournamentId(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	entries := make([]*debate_management.StudentFeedbackEntry, len(feedbacks))
	for i, f := range feedbacks {
		speakerPoints, err := strconv.ParseFloat(f.Speakerpoints, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse speaker points: %v", err)
		}

		// Parse judges info from JSON string
		var judgesData []map[string]interface{}
		err = json.Unmarshal([]byte(f.Judgesinfo), &judgesData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse judges info: %v", err)
		}

		judges := make([]*debate_management.JudgeInfo, len(judgesData))
		for j, judge := range judgesData {
			judges[j] = &debate_management.JudgeInfo{
				JudgeId:     int32(judge["judge_id"].(float64)),
				JudgeName:   judge["judge_name"].(string),
				IsHeadJudge: judge["is_head_judge"].(bool),
			}
		}

		entries[i] = &debate_management.StudentFeedbackEntry{
			RoundNumber:        f.Roundnumber,
			IsEliminationRound: f.Iseliminationround,
			DebateId:           f.Debateid,
			BallotId:           f.Ballotid,
			SpeakerPoints:      speakerPoints,
			Feedback:           f.Feedback.String,
			IsRead:             f.Isread.Bool,
			HeadJudgeName:      f.Headjudgename,
			RoomName:           f.Roomname,
			OpponentTeamName:   f.Opponentteamname,
			StudentTeamName:    f.Studentteamname,
			Judges:             judges,
		}
	}

	return &debate_management.GetStudentFeedbackResponse{
		FeedbackEntries: entries,
		TotalCount:      int32(totalCount),
	}, nil
}

func (s *FeedbackService) SubmitJudgeFeedback(ctx context.Context, req *debate_management.SubmitJudgeFeedbackRequest) (*debate_management.SubmitJudgeFeedbackResponse, error) {
	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	queries := models.New(s.db).WithTx(tx)
	studentID, err := queries.GetStudentByUserID(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get student ID: %v", err)
	}

	// Create judge feedback
	_, err = queries.CreateJudgeFeedback(ctx, models.CreateJudgeFeedbackParams{
		Judgeid:                sql.NullInt32{Int32: req.GetJudgeId(), Valid: true},
		Studentid:              sql.NullInt32{Int32: studentID.Studentid, Valid: true},
		Debateid:               sql.NullInt32{Int32: req.GetDebateId(), Valid: true},
		Clarityrating:          sql.NullFloat64{Float64: req.GetClarityRating(), Valid: true},
		Constructivenessrating: sql.NullFloat64{Float64: req.GetConstructivenessRating(), Valid: true},
		Timelinessrating:       sql.NullFloat64{Float64: req.GetTimelinessRating(), Valid: true},
		Fairnessrating:         sql.NullFloat64{Float64: req.GetFairnessRating(), Valid: true},
		Engagementrating:       sql.NullFloat64{Float64: req.GetEngagementRating(), Valid: true},
		Textfeedback:           sql.NullString{String: req.GetTextFeedback(), Valid: req.GetTextFeedback() != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create judge feedback: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &debate_management.SubmitJudgeFeedbackResponse{
		Success: true,
		Message: "Feedback submitted successfully",
	}, nil
}

func (s *FeedbackService) GetJudgeFeedback(ctx context.Context, req *debate_management.GetJudgeFeedbackRequest) (*debate_management.GetJudgeFeedbackResponse, error) {
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

    feedbacks, err := queries.GetJudgeFeedbackList(ctx, models.GetJudgeFeedbackListParams{
        Judgeid: sql.NullInt32{Int32: volunteer.Volunteerid, Valid: true},
        Limit:   int32(req.GetPageSize()),
        Offset:  int32((req.GetPage() - 1) * req.GetPageSize()),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get judge feedback: %v", err)
    }

    totalCount, err := queries.GetJudgeFeedbackCount(ctx, sql.NullInt32{Int32: volunteer.Volunteerid, Valid: true})
    if err != nil {
        return nil, fmt.Errorf("failed to get total count: %v", err)
    }

    entries := make([]*debate_management.JudgeFeedbackEntry, len(feedbacks))
    for i, f := range feedbacks {
        entries[i] = &debate_management.JudgeFeedbackEntry{
            StudentAlias:           generateAlias(f.Studentid.Int32),
            TournamentDate:         f.Tournamentdate.Format("2006-01-02"),
            IsRead:                 f.Isread.Bool,
            ClarityRating:          f.Clarityrating.Float64,
            ConstructivenessRating: f.Constructivenessrating.Float64,
            TimelinessRating:       f.Timelinessrating.Float64,
            FairnessRating:         f.Fairnessrating.Float64,
            EngagementRating:       f.Engagementrating.Float64,
            TextFeedback:           f.Textfeedback.String,
            RoundNumber:            f.Roundnumber,
            IsEliminationRound:     f.Iseliminationround,
            FeedbackId:            f.Feedbackid,  // Changed to use feedback_id
        }
    }

    return &debate_management.GetJudgeFeedbackResponse{
        FeedbackEntries: entries,
        TotalCount:      int32(totalCount),
    }, nil
}

func (s *FeedbackService) MarkStudentFeedbackAsRead(ctx context.Context, req *debate_management.MarkFeedbackAsReadRequest) (*debate_management.MarkFeedbackAsReadResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	queries := models.New(s.db)
	student, err := queries.GetStudentByUserID(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %v", err)
	}

	err = queries.MarkStudentFeedbackAsRead(ctx, models.MarkStudentFeedbackAsReadParams{
		Speakerid: student.Studentid,
		Ballotid:  req.GetFeedbackId(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to mark feedback as read: %v", err)
	}

	return &debate_management.MarkFeedbackAsReadResponse{
		Success: true,
		Message: "Feedback marked as read",
	}, nil
}

func (s *FeedbackService) MarkJudgeFeedbackAsRead(ctx context.Context, req *debate_management.MarkFeedbackAsReadRequest) (*debate_management.MarkFeedbackAsReadResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid user ID in token")
	}

	queries := models.New(s.db)
	volunteer, err := queries.GetVolunteerByUserID(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get volunteer: %v", err)
	}

	err = queries.MarkJudgeFeedbackAsRead(ctx, models.MarkJudgeFeedbackAsReadParams{
		Feedbackid: req.GetFeedbackId(),
		Judgeid:    sql.NullInt32{Int32: volunteer.Volunteerid, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to mark feedback as read: %v", err)
	}

	return &debate_management.MarkFeedbackAsReadResponse{
		Success: true,
		Message: "Feedback marked as read",
	}, nil
}

// generateAlias creates a consistent but anonymous name for a student
func generateAlias(studentID int32) string {
	rand.Seed(int64(studentID)) // Use student ID as seed for consistency
	adjectives := []string{"Swift", "Bright", "Clever", "Dynamic", "Eager"}
	nouns := []string{"Debater", "Speaker", "Thinker", "Scholar", "Mind"}
	return adjectives[rand.Intn(len(adjectives))] + " " + nouns[rand.Intn(len(nouns))]
}
