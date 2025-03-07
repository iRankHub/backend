package services

import (
	"context"
	"database/sql"
	_ "errors"
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

	// Ensure "Unassigned" room exists
	_, err := queries.EnsureUnassignedRoomExists(ctx, sql.NullInt32{
		Int32: req.GetTournamentId(),
		Valid: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ensure unassigned room exists: %v", err)
	}

	// Get all volunteers who accepted the tournament invitation
	availableJudges, err := queries.GetAvailableJudges(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get available judges: %v", err)
	}

	// Create a map to track unique judges and prevent duplicates
	uniqueJudges := make(map[int32]*debate_management.Judge)

	for _, j := range availableJudges {
		// Skip if we already have this judge in our map
		if _, exists := uniqueJudges[j.Userid]; exists {
			continue
		}

		preliminaryDebates, err := queries.CountJudgeDebates(ctx, models.CountJudgeDebatesParams{
			Judgeid:            j.Userid,
			Tournamentid:       req.GetTournamentId(),
			Iseliminationround: false,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to count preliminary debates: %v", err)
		}

		eliminationDebates, err := queries.CountJudgeDebates(ctx, models.CountJudgeDebatesParams{
			Judgeid:            j.Userid,
			Tournamentid:       req.GetTournamentId(),
			Iseliminationround: true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to count elimination debates: %v", err)
		}

		uniqueJudges[j.Userid] = &debate_management.Judge{
			JudgeId:            j.Userid,
			Name:               j.Name,
			IdebateId:          j.Idebatevolunteerid.String,
			PreliminaryDebates: int32(preliminaryDebates),
			EliminationDebates: int32(eliminationDebates),
		}
	}

	// Convert map to slice
	result := make([]*debate_management.Judge, 0, len(uniqueJudges))
	for _, judge := range uniqueJudges {
		result = append(result, judge)
	}

	return result, nil
}

// Update the GetJudge method in internal/services/judge.go
func (s *JudgeService) GetJudge(ctx context.Context, req *debate_management.GetJudgeRequest) (*debate_management.GetJudgeResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)

	// Ensure "Unassigned" room exists
	unassignedRoom, err := queries.EnsureUnassignedRoomExists(ctx, sql.NullInt32{
		Int32: req.GetTournamentId(),
		Valid: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ensure unassigned room exists: %v", err)
	}

	judge, err := queries.GetJudgeDetails(ctx, int32(req.GetJudgeId()))
	if err != nil {
		return nil, fmt.Errorf("failed to get judge details: %v", err)
	}

	// Get judge assignments for preliminary rounds
	preliminaryRooms, err := queries.GetJudgeRooms(ctx, models.GetJudgeRoomsParams{
		Judgeid:            int32(req.GetJudgeId()),
		Tournamentid:       req.GetTournamentId(),
		Iseliminationround: false,
	})
	if err != nil && !isNoRowsError(err) {
		return nil, fmt.Errorf("failed to get preliminary rooms: %v", err)
	}

	// Get judge assignments for elimination rounds
	eliminationRooms, err := queries.GetJudgeRooms(ctx, models.GetJudgeRoomsParams{
		Judgeid:            int32(req.GetJudgeId()),
		Tournamentid:       req.GetTournamentId(),
		Iseliminationround: true,
	})
	if err != nil && !isNoRowsError(err) {
		return nil, fmt.Errorf("failed to get elimination rooms: %v", err)
	}

	// Get tournament details
	tournament, err := queries.GetTournamentByID(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament: %v", err)
	}

	response := &debate_management.GetJudgeResponse{
		JudgeId:     int32(judge.Judgeid),
		Name:        judge.Name,
		IdebateId:   judge.Idebatevolunteerid.String,
		Preliminary: make(map[int32]*debate_management.RoomInfo),
		Elimination: make(map[int32]*debate_management.RoomInfo),
	}

	// Create maps for quick access to room assignments
	prelimMap := make(map[int32]models.GetJudgeRoomsRow)
	for _, room := range preliminaryRooms {
		prelimMap[room.Roundnumber] = room
	}

	elimMap := make(map[int32]models.GetJudgeRoomsRow)
	for _, room := range eliminationRooms {
		elimMap[room.Roundnumber] = room
	}

	// Add all preliminary rounds (including unassigned ones)
	for round := int32(1); round <= tournament.Numberofpreliminaryrounds; round++ {
		if room, ok := prelimMap[round]; ok {
			response.Preliminary[round] = &debate_management.RoomInfo{
				RoomId:      room.Roomid,
				RoomName:    room.Roomname,
				IsHeadJudge: room.Isheadjudge,
			}
		} else {
			response.Preliminary[round] = &debate_management.RoomInfo{
				RoomId:      unassignedRoom,
				RoomName:    "Unassigned",
				IsHeadJudge: false,
			}
		}
	}

	// Add all elimination rounds (including unassigned ones)
	for round := int32(1); round <= tournament.Numberofeliminationrounds; round++ {
		if room, ok := elimMap[round]; ok {
			response.Elimination[round] = &debate_management.RoomInfo{
				RoomId:      room.Roomid,
				RoomName:    room.Roomname,
				IsHeadJudge: room.Isheadjudge,
			}
		} else {
			response.Elimination[round] = &debate_management.RoomInfo{
				RoomId:      unassignedRoom,
				RoomName:    "Unassigned",
				IsHeadJudge: false,
			}
		}
	}

	return response, nil
}

// Helper function to check if error is "no rows" error
func isNoRowsError(err error) bool {
	return err == sql.ErrNoRows
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

		var oldRoomID int32
		var wasHeadJudge bool
		if err == nil { // Assignment exists
			oldRoomID = currentAssignment.Roomid
			wasHeadJudge = currentAssignment.Isheadjudge
		}

		// Update with the simplified swap logic
		err = assignmentService.UpdateJudgeAssignment(ctx, JudgeUpdateOperation{
			TournamentID:   req.GetTournamentId(),
			JudgeID:        req.GetJudgeId(),
			RoundNumber:    roundNumber,
			IsElimination:  false,
			OldRoomID:      oldRoomID,
			NewRoomID:      roomInfo.GetRoomId(),
			WasHeadJudge:   wasHeadJudge,
			NewIsHeadJudge: roomInfo.GetIsHeadJudge(),
			TargetJudgeID:  0, // No specific target judge, will find by room
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

		var oldRoomID int32
		var wasHeadJudge bool
		if err == nil { // Assignment exists
			oldRoomID = currentAssignment.Roomid
			wasHeadJudge = currentAssignment.Isheadjudge
		}

		// Update with the simplified swap logic
		err = assignmentService.UpdateJudgeAssignment(ctx, JudgeUpdateOperation{
			TournamentID:   req.GetTournamentId(),
			JudgeID:        int32(req.GetJudgeId()),
			RoundNumber:    int32(roundNumber),
			IsElimination:  true,
			OldRoomID:      oldRoomID,
			NewRoomID:      int32(roomInfo.GetRoomId()),
			WasHeadJudge:   wasHeadJudge,
			NewIsHeadJudge: roomInfo.GetIsHeadJudge(),
			TargetJudgeID:  0, // No specific target judge, will find by room
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
