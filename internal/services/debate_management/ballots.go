package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/iRankHub/backend/internal/grpc/proto/debate_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type BallotService struct {
	db *sql.DB
}

func NewBallotService(db *sql.DB) *BallotService {
	return &BallotService{db: db}
}

func (s *BallotService) GetBallots(ctx context.Context, req *debate_management.GetBallotsRequest) ([]*debate_management.Ballot, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	ballots, err := queries.GetBallotsByTournamentAndRound(ctx, models.GetBallotsByTournamentAndRoundParams{
		Tournamentid:       req.GetTournamentId(),
		Roundnumber:        req.GetRoundNumber(),
		Iseliminationround: req.GetIsElimination(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get ballots: %v", err)
	}

	return convertBallots(ballots), nil
}

func (s *BallotService) GetBallot(ctx context.Context, req *debate_management.GetBallotRequest) (*debate_management.Ballot, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	ballot, err := queries.GetBallotByID(ctx, req.GetBallotId())
	if err != nil {
		return nil, fmt.Errorf("failed to get ballot: %v", err)
	}

	return convertBallot(ballot), nil
}

func (s *BallotService) UpdateBallot(ctx context.Context, req *debate_management.UpdateBallotRequest) (*debate_management.Ballot, error) {
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

	// Update main ballot information
	err = queries.UpdateBallot(ctx, models.UpdateBallotParams{
		Ballotid: req.GetBallot().GetBallotId(),
		Team1totalscore: sql.NullString{
			String: fmt.Sprintf("%.2f", req.GetBallot().GetTeam1().GetTotalPoints()),
			Valid:  true,
		},
		Team2totalscore: sql.NullString{
			String: fmt.Sprintf("%.2f", req.GetBallot().GetTeam2().GetTotalPoints()),
			Valid:  true,
		},
		Recordingstatus: req.GetBallot().GetRecordingStatus(),
		Verdict:         req.GetBallot().GetVerdict(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update ballot: %v", err)
	}

	// Update speaker scores
	for _, speaker := range req.GetBallot().GetTeam1().GetSpeakers() {
		err = updateSpeakerScore(ctx, queries, speaker)
		if err != nil {
			return nil, err
		}
	}
	for _, speaker := range req.GetBallot().GetTeam2().GetSpeakers() {
		err = updateSpeakerScore(ctx, queries, speaker)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Fetch the updated ballot
	updatedBallot, err := queries.GetBallotByID(ctx, req.GetBallot().GetBallotId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated ballot: %v", err)
	}

	return convertBallot(updatedBallot), nil
}

func updateSpeakerScore(ctx context.Context, queries *models.Queries, speaker *debate_management.Speaker) error {
	err := queries.UpdateSpeakerScore(ctx, models.UpdateSpeakerScoreParams{
		Scoreid:       speaker.GetScoreId(),
		Speakerrank:   int32(speaker.GetRank()),
		Speakerpoints: fmt.Sprintf("%.2f", speaker.GetPoints()),
		Feedback:      sql.NullString{String: speaker.GetFeedback(), Valid: speaker.GetFeedback() != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to update speaker score: %v", err)
	}
	return nil
}

func convertBallots(dbBallots []models.GetBallotsByTournamentAndRoundRow) []*debate_management.Ballot {
	ballots := make([]*debate_management.Ballot, len(dbBallots))
	for i, dbBallot := range dbBallots {
		ballots[i] = &debate_management.Ballot{
			BallotId:      dbBallot.Ballotid,
			RoundNumber:   dbBallot.Roundnumber,
			IsElimination: dbBallot.Iseliminationround,
			RoomId:        dbBallot.Roomid,
			RoomName:      dbBallot.Roomname.String,
			Judges: []*debate_management.Judge{
				{
					JudgeId: dbBallot.Judgeid,
					Name:    dbBallot.Judgename,
				},
			},
			Team1: &debate_management.Team{
				TeamId: dbBallot.Team1id,
				Name:   dbBallot.Team1name,
			},
			Team2: &debate_management.Team{
				TeamId: dbBallot.Team2id,
				Name:   dbBallot.Team2name,
			},
			RecordingStatus: dbBallot.Recordingstatus,
			Verdict:         dbBallot.Verdict,
		}
	}
	return ballots
}

func convertBallot(dbBallot models.GetBallotByIDRow) *debate_management.Ballot {
	team1TotalPoints, _ := strconv.ParseFloat(dbBallot.Team1totalscore.String, 64)
	team2TotalPoints, _ := strconv.ParseFloat(dbBallot.Team2totalscore.String, 64)

	return &debate_management.Ballot{
		BallotId:      dbBallot.Ballotid,
		RoundNumber:   dbBallot.Roundnumber,
		IsElimination: dbBallot.Iseliminationround,
		RoomId:        dbBallot.Roomid,
		RoomName:      dbBallot.Roomname.String,
		Judges: []*debate_management.Judge{
			{
				JudgeId: dbBallot.Judgeid,
				Name:    dbBallot.Judgename,
			},
		},
		Team1: &debate_management.Team{
			TeamId:      dbBallot.Team1id,
			Name:        dbBallot.Team1name,
			TotalPoints: team1TotalPoints,
		},
		Team2: &debate_management.Team{
			TeamId:      dbBallot.Team2id,
			Name:        dbBallot.Team2name,
			TotalPoints: team2TotalPoints,
		},
		RecordingStatus: dbBallot.Recordingstatus,
		Verdict:         dbBallot.Verdict,
	}
}

func (s *BallotService) GetSpeakerScores(ctx context.Context, ballotID int32) ([]*debate_management.Speaker, error) {
	queries := models.New(s.db)
	scores, err := queries.GetSpeakerScoresByBallot(ctx, ballotID)
	if err != nil {
		return nil, fmt.Errorf("failed to get speaker scores: %v", err)
	}

	speakers := make([]*debate_management.Speaker, len(scores))
	for i, score := range scores {
		points, _ := strconv.ParseFloat(score.Speakerpoints, 64)
		speakers[i] = &debate_management.Speaker{
			SpeakerId: score.Speakerid,
			Name:      score.Firstname + " " + score.Lastname,
			ScoreId:   score.Scoreid,
			Rank:      int32(score.Speakerrank),
			Points:    points,
			Feedback:  score.Feedback.String,
		}
	}

	return speakers, nil
}
func (s *BallotService) validateAdminRole(token string) (map[string]interface{}, error) {
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

func (s *BallotService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}
