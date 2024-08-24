package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/grpc/proto/debate_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type JudgeService struct {
	db *sql.DB
}

func NewJudgeService(db *sql.DB) *JudgeService {
	return &JudgeService{db: db}
}

func (s *JudgeService) GetJudges(ctx context.Context, req *debate_management.GetJudgesRequest) ([]*debate_management.Judge, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	judges, err := queries.GetJudgesByTournamentAndRound(ctx, models.GetJudgesByTournamentAndRoundParams{
		Tournamentid:  req.GetTournamentId(),
		Roundnumber:   req.GetRoundNumber(),
		Iselimination: req.GetIsElimination(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get judges: %v", err)
	}

	return convertJudges(judges), nil
}

func (s *JudgeService) GetJudge(ctx context.Context, req *debate_management.GetJudgeRequest) (*debate_management.Judge, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	judge, err := queries.GetJudgeByID(ctx, req.GetJudgeId())
	if err != nil {
		return nil, fmt.Errorf("failed to get judge: %v", err)
	}

	return convertJudge(judge), nil
}

func (s *JudgeService) AssignJudges(ctx context.Context, req *debate_management.AssignJudgesRequest) ([]*debate_management.Pairing, error) {
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

	// Get tournament details
	tournament, err := queries.GetTournamentByID(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament: %v", err)
	}

	// Get available judges
	judges, err := queries.GetAvailableJudges(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get available judges: %v", err)
	}

	// Get pairings for the round
	pairings, err := queries.GetPairingsByTournamentAndRound(ctx, models.GetPairingsByTournamentAndRoundParams{
		Tournamentid:       req.GetTournamentId(),
		Roundnumber:        req.GetRoundNumber(),
		Iseliminationround: req.GetIsElimination(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pairings: %v", err)
	}

	// Determine number of judges per debate
	judgesPerDebate := tournament.Judgesperdebatepreliminary
	if req.GetIsElimination() {
		judgesPerDebate = tournament.Judgesperdebateelimination
	}

	// Assign judges to debates
	assignedPairings := make([]*debate_management.Pairing, len(pairings))
    for i, pairing := range pairings {
        assignedJudges := make([]*debate_management.Judge, int(judgesPerDebate))
		for j := int32(0); j < judgesPerDebate; j++ {
			if len(judges) == 0 {
				return nil, fmt.Errorf("not enough judges for all debates")
			}

			judgeIndex := 0 // TODO: Implement a more sophisticated judge selection algorithm
			judge := judges[judgeIndex]

			isHeadJudge := j == 0

			err := queries.AssignJudgeToDebate(ctx, models.AssignJudgeToDebateParams{
				Tournamentid:  req.GetTournamentId(),
				Judgeid:       judge.Userid,
				Debateid:      pairing.Debateid,
				Roundnumber:   req.GetRoundNumber(),
				Iselimination: req.GetIsElimination(),
				Isheadjudge:   isHeadJudge,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to assign judge to debate: %v", err)
			}

			assignedJudges[j] = &debate_management.Judge{
				JudgeId:     judge.Userid,
				Name:        judge.Name,
				IsHeadJudge: isHeadJudge,
			}

			// Remove the assigned judge from the available judges
			judges = append(judges[:judgeIndex], judges[judgeIndex+1:]...)
		}

		assignedPairings[i] = convertJudgePairing(pairing)
		assignedPairings[i] = &debate_management.Pairing{
            PairingId:          assignedPairings[i].PairingId,
            RoundNumber:        assignedPairings[i].RoundNumber,
            IsEliminationRound: assignedPairings[i].IsEliminationRound,
            RoomId:             assignedPairings[i].RoomId,
            Team1:              assignedPairings[i].Team1,
            Team2:              assignedPairings[i].Team2,
            Judges:             assignedJudges,
        }
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return assignedPairings, nil
}

func (s *JudgeService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}

func (s *JudgeService) validateAdminRole(token string) (map[string]interface{}, error) {
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

func convertJudges(dbJudges []models.GetJudgesByTournamentAndRoundRow) []*debate_management.Judge {
	judges := make([]*debate_management.Judge, len(dbJudges))
	for i, dbJudge := range dbJudges {
		judges[i] = &debate_management.Judge{
			JudgeId:     dbJudge.Userid,
			Name:        dbJudge.Name,
			Email:       dbJudge.Email,
			IsHeadJudge: dbJudge.Isheadjudge,
		}
	}
	return judges
}

func convertJudge(dbJudge models.GetJudgeByIDRow) *debate_management.Judge {
	return &debate_management.Judge{
		JudgeId: dbJudge.Userid,
		Name:    dbJudge.Name,
		Email:   dbJudge.Email,
	}
}
func convertJudgePairing(dbPairing models.GetPairingsByTournamentAndRoundRow) *debate_management.Pairing {
	return &debate_management.Pairing{
		PairingId:          dbPairing.Debateid,
		RoundNumber:        dbPairing.Roundnumber,
		IsEliminationRound: dbPairing.Iseliminationround,
		RoomId:             dbPairing.Roomid,
		Team1: &debate_management.Team{
			TeamId: dbPairing.Team1id,
			Name:   dbPairing.Team1name,
		},
		Team2: &debate_management.Team{
			TeamId: dbPairing.Team2id,
			Name:   dbPairing.Team2name,
		},
	}
}
