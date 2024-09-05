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
	judges, err := queries.GetJudgesForTournament(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get judges: %v", err)
	}

	var result []*debate_management.Judge
	for _, j := range judges {
		preliminaryDebates, err := queries.CountJudgeDebates(ctx, models.CountJudgeDebatesParams{
			Judgeid:            j.Judgeid,
			Tournamentid:       req.GetTournamentId(),
			Iseliminationround: false,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to count preliminary debates: %v", err)
		}

		eliminationDebates, err := queries.CountJudgeDebates(ctx, models.CountJudgeDebatesParams{
			Judgeid:            j.Judgeid,
			Tournamentid:       req.GetTournamentId(),
			Iseliminationround: true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to count elimination debates: %v", err)
		}

		result = append(result, &debate_management.Judge{
			JudgeId:             int32(j.Judgeid),
			Name:                j.Name,
			IdebateId:           j.Idebatevolunteerid.String,
			PreliminaryDebates:  int32(preliminaryDebates),
			EliminationDebates:  int32(eliminationDebates),
		})
	}

	return result, nil
}

func (s *JudgeService) GetJudge(ctx context.Context, req *debate_management.GetJudgeRequest) (*debate_management.GetJudgeResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	judge, err := queries.GetJudgeDetails(ctx, int32(req.GetJudgeId()))
	if err != nil {
		return nil, fmt.Errorf("failed to get judge details: %v", err)
	}

	preliminaryRooms, err := queries.GetJudgeRooms(ctx, models.GetJudgeRoomsParams{
		Judgeid:            int32(req.GetJudgeId()),
		Tournamentid:       req.GetTournamentId(),
		Iseliminationround: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get preliminary rooms: %v", err)
	}

	eliminationRooms, err := queries.GetJudgeRooms(ctx, models.GetJudgeRoomsParams{
		Judgeid:            int32(req.GetJudgeId()),
		Tournamentid:       req.GetTournamentId(),
		Iseliminationround: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get elimination rooms: %v", err)
	}

	response := &debate_management.GetJudgeResponse{
		JudgeId:     int32(judge.Judgeid),
		Name:        judge.Name,
		IdebateId:   judge.Idebatevolunteerid.String,
		Preliminary: make(map[int32]*debate_management.RoomInfo),
		Elimination: make(map[int32]*debate_management.RoomInfo),
	}

	for _, room := range preliminaryRooms {
		response.Preliminary[room.Roundnumber] = &debate_management.RoomInfo{
			RoomId:   int32(room.Roomid),
			RoomName: room.Roomname,
		}
	}

	for _, room := range eliminationRooms {
		response.Elimination[room.Roundnumber] = &debate_management.RoomInfo{
			RoomId:   int32(room.Roomid),
			RoomName: room.Roomname,
		}
	}

	return response, nil
}

func (s *JudgeService) UpdateJudge(ctx context.Context, req *debate_management.UpdateJudgeRequest) (*debate_management.UpdateJudgeResponse, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	for roundNumber, roomID := range req.GetRoomAssignments() {
		err := qtx.UpdateJudgeRoom(ctx, models.UpdateJudgeRoomParams{
			Judgeid:      int32(req.GetJudgeId()),
			Tournamentid: req.GetTournamentId(),
			Roundnumber:  int32(roundNumber),
			Roomid:       int32(roomID),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update judge room: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &debate_management.UpdateJudgeResponse{
		Success: true,
		Message: "Judge rooms updated successfully",
	}, nil
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