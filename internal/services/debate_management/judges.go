package services

import (
	"context"
	"database/sql"
	"errors"
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
			JudgeId:            int32(j.Judgeid),
			Name:               j.Name,
			IdebateId:          j.Idebatevolunteerid.String,
			PreliminaryDebates: int32(preliminaryDebates),
			EliminationDebates: int32(eliminationDebates),
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

	assignmentService := NewJudgeAssignmentService(s.db)
	queries := models.New(s.db)

	// Process preliminary rounds
	for roundNumber, roomInfo := range req.GetPreliminary() {
		// Get current assignment to know the old room
		currentAssignment, err := queries.GetJudgeAssignment(ctx, models.GetJudgeAssignmentParams{
			Tournamentid:  req.GetTournamentId(),
			Judgeid:       int32(req.GetJudgeId()),
			Roundnumber:   int32(roundNumber),
			Iselimination: false,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get current assignment for round %d: %v", roundNumber, err)
		}

		var oldRoomID int32
		var wasHeadJudge bool
		if err == nil { // Assignment exists
			oldRoomID = currentAssignment.Roomid
			wasHeadJudge = currentAssignment.Isheadjudge
		}

		err = assignmentService.UpdateJudgeAssignment(ctx, JudgeUpdateOperation{
			TournamentID:   req.GetTournamentId(),
			JudgeID:        int32(req.GetJudgeId()),
			RoundNumber:    int32(roundNumber),
			IsElimination:  false,
			OldRoomID:      oldRoomID,
			NewRoomID:      int32(roomInfo.GetRoomId()),
			WasHeadJudge:   wasHeadJudge,
			NewIsHeadJudge: roomInfo.GetIsHeadJudge(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update preliminary round %d: %v", roundNumber, err)
		}
	}

	// Process elimination rounds
	for roundNumber, roomInfo := range req.GetElimination() {
		// Get current assignment to know the old room
		currentAssignment, err := queries.GetJudgeAssignment(ctx, models.GetJudgeAssignmentParams{
			Tournamentid:  req.GetTournamentId(),
			Judgeid:       int32(req.GetJudgeId()),
			Roundnumber:   int32(roundNumber),
			Iselimination: true,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get current assignment for round %d: %v", roundNumber, err)
		}

		var oldRoomID int32
		var wasHeadJudge bool
		if err == nil { // Assignment exists
			oldRoomID = currentAssignment.Roomid
			wasHeadJudge = currentAssignment.Isheadjudge
		}

		err = assignmentService.UpdateJudgeAssignment(ctx, JudgeUpdateOperation{
			TournamentID:   req.GetTournamentId(),
			JudgeID:        int32(req.GetJudgeId()),
			RoundNumber:    int32(roundNumber),
			IsElimination:  true,
			OldRoomID:      oldRoomID,
			NewRoomID:      int32(roomInfo.GetRoomId()),
			WasHeadJudge:   wasHeadJudge,
			NewIsHeadJudge: roomInfo.GetIsHeadJudge(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update elimination round %d: %v", roundNumber, err)
		}
	}

	return &debate_management.UpdateJudgeResponse{
		Success: true,
		Message: "Judge assignments updated successfully",
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
