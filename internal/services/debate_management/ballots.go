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

	convertedBallot := convertBallot(ballot)

	// Fetch speaker scores
	speakerScores, err := queries.GetSpeakerScoresByBallot(ctx, req.GetBallotId())
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

	// Check if head judge has already submitted
	if !isAdmin && currentBallot.HeadJudgeSubmitted.Bool {
		return nil, fmt.Errorf("ballot can be submitted only once. head judge already submitted")
	}

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
		return nil, fmt.Errorf("failed to update ballot: %v", err)
	}

	// Get debate information
	debate, err := queries.GetDebateByBallotID(ctx, req.GetBallot().GetBallotId())
	if err != nil {
		return nil, fmt.Errorf("failed to get debate information: %v", err)
	}

	// Update TeamScores for Team1
	err = updateTeamScore(ctx, queries, debate.Team1id, req.GetBallot().GetTeam1().GetTotalPoints(), debate.Debateid, debate.Iseliminationround)
	if err != nil {
		return nil, err
	}

	// Update TeamScores for Team2
	err = updateTeamScore(ctx, queries, debate.Team2id, req.GetBallot().GetTeam2().GetTotalPoints(), debate.Debateid, debate.Iseliminationround)
	if err != nil {
		return nil, err
	}

    // Update speaker scores
    if req.GetBallot().GetTeam1().GetName() != "Public Speaking" {
        for _, speaker := range req.GetBallot().GetTeam1().GetSpeakers() {
            err = updateSpeakerScore(ctx, queries, speaker, req.GetBallot().GetBallotId())
            if err != nil {
                return nil, err
            }
        }
    }
    if req.GetBallot().GetTeam2().GetName() != "Public Speaking" {
        for _, speaker := range req.GetBallot().GetTeam2().GetSpeakers() {
            err = updateSpeakerScore(ctx, queries, speaker, req.GetBallot().GetBallotId())
            if err != nil {
                return nil, err
            }
        }
    }

	// Update team ranks
	err = updateTeamRanks(ctx, queries, debate.Team1id, req.GetBallot().GetBallotId(), debate.Debateid)
	if err != nil {
		return nil, err
	}
	err = updateTeamRanks(ctx, queries, debate.Team2id, req.GetBallot().GetBallotId(), debate.Debateid)
	if err != nil {
		return nil, err
	}

	    // Update team stats
    err = queries.UpdateTeamStats(ctx, models.UpdateTeamStatsParams{
        Teamid:       debate.Team1id,
        Tournamentid: debate.Tournamentid,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to update team 1 stats: %v", err)
    }
    err = queries.UpdateTeamStats(ctx, models.UpdateTeamStatsParams{
        Teamid:       debate.Team2id,
        Tournamentid: debate.Tournamentid,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to update team 2 stats: %v", err)
    }

	// Fetch the updated ballot before committing the transaction
	updatedBallot, err := queries.GetBallotByID(ctx, req.GetBallot().GetBallotId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated ballot: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return convertBallot(updatedBallot), nil
}

func updateTeamScore(ctx context.Context, queries *models.Queries, teamID int32, totalScore float64, debateID int32, isElimination bool) error {
    // First, try to update
    err := queries.UpdateTeamScore(ctx, models.UpdateTeamScoreParams{
		Teamid:        sql.NullInt32{Int32: teamID, Valid: true},
		Debateid:      sql.NullInt32{Int32: debateID, Valid: true},
		Totalscore:    sql.NullString{String: fmt.Sprintf("%.2f", totalScore), Valid: true},
		Iselimination: sql.NullBool{Bool: isElimination, Valid: true},
    })
    if err != nil {
        return fmt.Errorf("failed to update team score: %v", err)
    }

    // Then, try to insert if the update didn't affect any rows
    err = queries.InsertTeamScore(ctx, models.InsertTeamScoreParams{
		Teamid:        sql.NullInt32{Int32: teamID, Valid: true},
		Debateid:      sql.NullInt32{Int32: debateID, Valid: true},
		Totalscore:    sql.NullString{String: fmt.Sprintf("%.2f", totalScore), Valid: true},
		Iselimination: sql.NullBool{Bool: isElimination, Valid: true},
    })
    if err != nil {
        return fmt.Errorf("failed to insert team score: %v", err)
    }

    return nil
}

func updateSpeakerScore(ctx context.Context, queries *models.Queries, speaker *debate_management.Speaker, ballotID int32) error {
	err := queries.UpdateSpeakerScore(ctx, models.UpdateSpeakerScoreParams{
		Ballotid:      ballotID,
		Speakerid:     speaker.GetSpeakerId(),
		Speakerrank:   int32(speaker.GetRank()),
		Speakerpoints: fmt.Sprintf("%.2f", speaker.GetPoints()),
		Feedback:      sql.NullString{String: speaker.GetFeedback(), Valid: speaker.GetFeedback() != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to update speaker score: %v", err)
	}
	return nil
}

func convertSpeakerScore(score models.GetSpeakerScoresByBallotRow) *debate_management.Speaker {
	points, _ := strconv.ParseFloat(score.Speakerpoints, 64)
	return &debate_management.Speaker{
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

func updateTeamRanks(ctx context.Context, queries *models.Queries, teamID int32, ballotID int32, debateID int32) error {
	avgRank, err := queries.GetTeamAverageRank(ctx, models.GetTeamAverageRankParams{
		Teamid:   teamID,
		Ballotid: ballotID,
	})
	if err != nil {
		return fmt.Errorf("failed to get team average rank: %v", err)
	}

	// Convert the average rank to an integer
	intRank := int32(avgRank)

	err = queries.UpdateTeamScoreRank(ctx, models.UpdateTeamScoreRankParams{
		Teamid:   sql.NullInt32{Int32: teamID, Valid: true},
		Debateid: sql.NullInt32{Int32: debateID, Valid: true},
		Rank:     sql.NullInt32{Int32: intRank, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update team rank: %v", err)
	}
	return nil
}

func convertBallots(dbBallots []models.GetBallotsByTournamentAndRoundRow) []*debate_management.Ballot {
    ballots := make([]*debate_management.Ballot, len(dbBallots))
    for i, dbBallot := range dbBallots {
        ballots[i] = &debate_management.Ballot{
            BallotId:        dbBallot.Ballotid,
            RoundNumber:     dbBallot.Roundnumber,
            IsElimination:   dbBallot.Iseliminationround,
            RoomName:        dbBallot.Roomname,
            Judges: []*debate_management.Judge{
                {
                    Name: dbBallot.Headjudgename,
                },
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
	return convertBallot(models.GetBallotByIDRow(dbBallot))
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