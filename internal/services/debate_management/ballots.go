package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

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

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	// Get the ballot
	ballot, err := queries.GetBallotByID(ctx, req.GetBallotId())
	if err != nil {
		return nil, fmt.Errorf("failed to get ballot: %v", err)
	}

	// Convert the ballot to proto format
	convertedBallot := convertBallot(ballot)

	// Get all speaker scores in a single query
	speakerScores, err := queries.GetSpeakerScoresByBallot(ctx, req.GetBallotId())
	if err != nil {
		return nil, fmt.Errorf("failed to get speaker scores: %v", err)
	}

	// Create maps to store speakers for each team
	team1Speakers := make([]*debate_management.Speaker, 0)
	team2Speakers := make([]*debate_management.Speaker, 0)
	team1Names := make([]string, 0)
	team2Names := make([]string, 0)

	// Sort speakers into their respective teams
	for _, score := range speakerScores {
		speaker := convertSpeakerScore(score)
		if score.Teamid == convertedBallot.Team1.TeamId {
			team1Speakers = append(team1Speakers, speaker)
			team1Names = append(team1Names, speaker.Name)
		} else if score.Teamid == convertedBallot.Team2.TeamId {
			team2Speakers = append(team2Speakers, speaker)
			team2Names = append(team2Names, speaker.Name)
		} else {
			log.Printf("Warning: Speaker %d doesn't match either team (Team1: %d, Team2: %d)",
				speaker.SpeakerId, convertedBallot.Team1.TeamId, convertedBallot.Team2.TeamId)
		}
	}

	// Sort speakers by rank within each team
	sortSpeakersByRank(team1Speakers)
	sortSpeakersByRank(team2Speakers)

	// Assign sorted speakers to teams
	convertedBallot.Team1.Speakers = team1Speakers
	convertedBallot.Team1.SpeakerNames = team1Names
	convertedBallot.Team2.Speakers = team2Speakers
	convertedBallot.Team2.SpeakerNames = team2Names

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return convertedBallot, nil
}

func sortSpeakersByRank(speakers []*debate_management.Speaker) {
	sort.Slice(speakers, func(i, j int) bool {
		return speakers[i].Rank < speakers[j].Rank
	})
}

func (s *BallotService) GetBallotByJudgeID(ctx context.Context, req *debate_management.GetBallotByJudgeIDRequest) (*debate_management.Ballot, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	ballot, err := queries.GetBallotByJudgeID(ctx, models.GetBallotByJudgeIDParams{
		Judgeid:      req.GetJudgeId(),
		Tournamentid: req.GetTournamentId(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get ballot: %v", err)
	}

	convertedBallot := convertJudgeBallot(ballot)

	// Fetch speaker scores
	speakerScores, err := queries.GetSpeakerScoresByBallot(ctx, ballot.Ballotid)
	if err != nil {
		return nil, fmt.Errorf("failed to get speaker scores: %v", err)
	}

	// Assign speaker scores to teams
	for _, score := range speakerScores {
		speaker := convertSpeakerScore(score)
		if score.Teamid == convertedBallot.Team1.TeamId {
			convertedBallot.Team1.Speakers = append(convertedBallot.Team1.Speakers, speaker)
			convertedBallot.Team1.SpeakerNames = append(convertedBallot.Team1.SpeakerNames, speaker.Name)
		} else if score.Teamid == convertedBallot.Team2.TeamId {
			convertedBallot.Team2.Speakers = append(convertedBallot.Team2.Speakers, speaker)
			convertedBallot.Team2.SpeakerNames = append(convertedBallot.Team2.SpeakerNames, speaker.Name)
		}
	}

	return convertedBallot, nil
}

func (s *BallotService) UpdateBallot(ctx context.Context, req *debate_management.UpdateBallotRequest) (*debate_management.Ballot, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user role in token")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	// Check if the user is admin or head judge
	isAdmin := userRole == "admin"
	isHeadJudge, err := queries.IsHeadJudgeForBallot(ctx, models.IsHeadJudgeForBallotParams{
		Ballotid: req.GetBallot().GetBallotId(),
		Judgeid:  int32(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check if user is head judge: %v", err)
	}

	if !isAdmin && !isHeadJudge {
		return nil, fmt.Errorf("unauthorized: only admins or the head judge can update this ballot")
	}

	// Get the current ballot state
	currentBallot, err := queries.GetBallotByID(ctx, req.GetBallot().GetBallotId())
	if err != nil {
		return nil, fmt.Errorf("failed to get current ballot state: %v", err)
	}
	log.Printf("Current ballot state: %+v\n", currentBallot)

	// Check if head judge has already submitted
	if !isAdmin && currentBallot.HeadJudgeSubmitted.Bool {
		log.Printf("Ballot %d has already been submitted by the head judge\n", req.GetBallot().GetBallotId())
		return nil, fmt.Errorf("ballot can be submitted only once. head judge already submitted")
	}

	// Update main ballot information
	log.Printf("Updating main ballot information: ballotID=%d", req.GetBallot().GetBallotId())
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
		Recordingstatus: "Recorded",
		Verdict:         req.GetBallot().GetVerdict(),
		Team1feedback:   sql.NullString{String: req.GetBallot().GetTeam1().GetFeedback(), Valid: true},
		Team2feedback:   sql.NullString{String: req.GetBallot().GetTeam2().GetFeedback(), Valid: true},
		LastUpdatedBy: sql.NullInt32{
			Int32: int32(userID),
			Valid: true,
		},
		HeadJudgeSubmitted: sql.NullBool{
			Bool:  isHeadJudge,
			Valid: true,
		},
	})
	if err != nil {
		log.Printf("Failed to update ballot: %v", err)
		return nil, fmt.Errorf("failed to update ballot: %v", err)
	}
	log.Printf("Successfully updated main ballot information for ballotID=%d", req.GetBallot().GetBallotId())

	// Update speaker scores
	log.Println("Updating speaker scores")
	err = updateSpeakerScores(ctx, queries, req.GetBallot())
	if err != nil {
		log.Printf("Failed to update speaker scores: %v", err)
		return nil, err
	}

	// Update team scores
	log.Println("Updating team scores")
	err = updateTeamScores(ctx, queries, req.GetBallot())
	if err != nil {
		log.Printf("Failed to update team scores: %v", err)
		return nil, err
	}

	// Fetch the updated ballot
	updatedBallot, err := queries.GetBallotByID(ctx, req.GetBallot().GetBallotId())
	if err != nil {
		log.Printf("Failed to fetch updated ballot: %v", err)
		return nil, fmt.Errorf("failed to fetch updated ballot: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully updated ballot: %+v", updatedBallot)
	return convertBallot(updatedBallot), nil
}

func updateSpeakerScores(ctx context.Context, queries *models.Queries, ballot *debate_management.Ballot) error {
	log.Println("Starting updateSpeakerScores")

	for _, team := range []*debate_management.Team{ballot.GetTeam1(), ballot.GetTeam2()} {
		for _, speaker := range team.GetSpeakers() {
			log.Printf("Updating speaker score: speakerID=%d, rank=%d, points=%.2f, feedback=%s",
				speaker.GetSpeakerId(), speaker.GetRank(), speaker.GetPoints(), speaker.GetFeedback())
			err := queries.UpdateSpeakerScore(ctx, models.UpdateSpeakerScoreParams{
				Ballotid:      ballot.GetBallotId(),
				Speakerid:     speaker.GetSpeakerId(),
				Speakerrank:   int32(speaker.GetRank()),
				Speakerpoints: fmt.Sprintf("%.2f", speaker.GetPoints()),
				Feedback:      sql.NullString{String: speaker.GetFeedback(), Valid: speaker.GetFeedback() != ""},
			})
			if err != nil {
				log.Printf("Failed to update speaker score: %v", err)
				return fmt.Errorf("failed to update speaker score: %v", err)
			}
		}
	}

	log.Println("Finished updateSpeakerScores successfully")
	return nil
}

func updateTeamScores(ctx context.Context, queries *models.Queries, ballot *debate_management.Ballot) error {
	log.Println("Starting updateTeamScores")

	// Update TeamScores for Team1
	log.Printf("Updating TeamScores for Team1: teamID=%d, totalScore=%.2f, debateID=%d",
		ballot.GetTeam1().GetTeamId(), ballot.GetTeam1().GetTotalPoints(), ballot.GetBallotId())
	err := updateTeamScore(ctx, queries, ballot.GetTeam1().GetTeamId(), ballot.GetTeam1().GetTotalPoints(), ballot.GetBallotId())
	if err != nil {
		log.Printf("Failed to update Team1 score: %v", err)
		return err
	}

	// Update TeamScores for Team2
	log.Printf("Updating TeamScores for Team2: teamID=%d, totalScore=%.2f, debateID=%d",
		ballot.GetTeam2().GetTeamId(), ballot.GetTeam2().GetTotalPoints(), ballot.GetBallotId())
	err = updateTeamScore(ctx, queries, ballot.GetTeam2().GetTeamId(), ballot.GetTeam2().GetTotalPoints(), ballot.GetBallotId())
	if err != nil {
		log.Printf("Failed to update Team2 score: %v", err)
		return err
	}

	log.Println("Finished updateTeamScores successfully")
	return nil
}

func updateTeamScore(ctx context.Context, queries *models.Queries, teamID int32, totalScore float64, debateID int32) error {
	log.Printf("Starting updateTeamScore: teamID=%d, totalScore=%.2f, debateID=%d",
		teamID, totalScore, debateID)

	err := queries.UpdateTeamScore(ctx, models.UpdateTeamScoreParams{
		Teamid:     sql.NullInt32{Int32: teamID, Valid: true},
		Debateid:   sql.NullInt32{Int32: debateID, Valid: true},
		Totalscore: sql.NullString{String: fmt.Sprintf("%.2f", totalScore), Valid: true},
	})
	if err != nil {
		log.Printf("Failed to update team score: %v", err)
		return fmt.Errorf("failed to update team score: %v", err)
	}

	log.Println("Finished updateTeamScore successfully")
	return nil
}

func convertSpeakerScore(score models.GetSpeakerScoresByBallotRow) *debate_management.Speaker {
	points, err := strconv.ParseFloat(score.Speakerpoints, 64)
	if err != nil {
		log.Printf("Warning: Failed to parse speaker points '%s': %v", score.Speakerpoints, err)
		points = 0
	}

	return &debate_management.Speaker{
		SpeakerId: score.Speakerid,
		Name:      strings.TrimSpace(score.Firstname + " " + score.Lastname),
		ScoreId:   score.Scoreid,
		Rank:      int32(score.Speakerrank),
		Points:    points,
		Feedback:  score.Feedback.String,
		TeamId:    score.Teamid,
		TeamName:  score.Teamname,
	}
}

func convertBallots(dbBallots []models.GetBallotsByTournamentAndRoundRow) []*debate_management.Ballot {
	ballots := make([]*debate_management.Ballot, len(dbBallots))
	for i, dbBallot := range dbBallots {
		ballots[i] = &debate_management.Ballot{
			BallotId:      dbBallot.Ballotid,
			RoundNumber:   dbBallot.Roundnumber,
			IsElimination: dbBallot.Iseliminationround,
			RoomName:      dbBallot.Roomname,
			Judges: []*debate_management.Judge{
				{
					Name: dbBallot.Headjudgename,
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
			Feedback:    dbBallot.Team1feedback.String,
		},
		Team2: &debate_management.Team{
			TeamId:      dbBallot.Team2id,
			Name:        dbBallot.Team2name,
			TotalPoints: team2TotalPoints,
			Feedback:    dbBallot.Team2feedback.String,
		},
		RecordingStatus:    dbBallot.Recordingstatus,
		Verdict:            dbBallot.Verdict,
		LastUpdatedBy:      dbBallot.LastUpdatedBy.Int32,
		LastUpdatedAt:      dbBallot.LastUpdatedAt.Time.String(),
		HeadJudgeSubmitted: dbBallot.HeadJudgeSubmitted.Bool,
	}
}

func convertJudgeBallot(dbBallot models.GetBallotByJudgeIDRow) *debate_management.Ballot {
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
			Feedback:    dbBallot.Team1feedback.String,
		},
		Team2: &debate_management.Team{
			TeamId:      dbBallot.Team2id,
			Name:        dbBallot.Team2name,
			TotalPoints: team2TotalPoints,
			Feedback:    dbBallot.Team2feedback.String,
		},
		RecordingStatus:    dbBallot.Recordingstatus,
		Verdict:            dbBallot.Verdict,
		LastUpdatedBy:      dbBallot.LastUpdatedBy.Int32,
		LastUpdatedAt:      dbBallot.LastUpdatedAt.Time.String(),
		HeadJudgeSubmitted: dbBallot.HeadJudgeSubmitted.Bool,
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
			TeamId:    score.Teamid,
			TeamName:  score.Teamname,
		}
	}

	return speakers, nil
}

func (s *BallotService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}
