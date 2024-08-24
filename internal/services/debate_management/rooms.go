package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/grpc/proto/debate_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type RoomService struct {
	db *sql.DB
}

func NewRoomService(db *sql.DB) *RoomService {
	return &RoomService{db: db}
}

func (s *RoomService) GetRooms(ctx context.Context, req *debate_management.GetRoomsRequest) ([]*debate_management.Room, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	rooms, err := queries.GetRoomsByTournamentAndRound(ctx, models.GetRoomsByTournamentAndRoundParams{
		Tournamentid:  req.GetTournamentId(),
		Roundnumber:   req.GetRoundNumber(),
		Iselimination: req.GetIsElimination(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %v", err)
	}

	return convertRooms(rooms), nil
}

func (s *RoomService) GetRoom(ctx context.Context, req *debate_management.GetRoomRequest) (*debate_management.Room, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	roomData, err := queries.GetRoomByID(ctx, req.GetRoomId())
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %v", err)
	}

	return convertRoomWithStatus(roomData), nil
}

func convertRoomWithStatus(dbRoomData []models.GetRoomByIDRow) *debate_management.Room {
	if len(dbRoomData) == 0 {
		return nil
	}

	room := &debate_management.Room{
		RoomId:      dbRoomData[0].Roomid,
		RoomName:    dbRoomData[0].Roomname,
		RoundStatus: make([]*debate_management.RoundStatus, len(dbRoomData)),
	}

	for i, data := range dbRoomData {
		room.RoundStatus[i] = &debate_management.RoundStatus{
			RoundNumber:   data.Roundnumber.Int32,
			IsElimination: data.Iselimination.Bool,
			IsOccupied:    data.Isoccupied.Bool,
		}
	}

	return room
}

func (s *RoomService) UpdateRoom(ctx context.Context, req *debate_management.UpdateRoomRequest) (*debate_management.Room, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	queries := models.New(s.db)

	// Update the room name
	_, err = queries.UpdateRoom(ctx, models.UpdateRoomParams{
		Roomid:   req.GetRoom().GetRoomId(),
		Roomname: req.GetRoom().GetRoomName(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update room: %v", err)
	}

	// Fetch the updated room data, including round status
	updatedRoomData, err := queries.GetRoomByID(ctx, req.GetRoom().GetRoomId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated room data: %v", err)
	}

	// Convert the updated room data to the debate_management.Room type
	return convertRoomWithStatus(updatedRoomData), nil
}
func (s *RoomService) AssignRoomsToDebates(ctx context.Context, tournamentID int32, roundNumber int32, isElimination bool) error {
	queries := models.New(s.db)

	// Get available rooms
	availableRooms, err := queries.GetAvailableRooms(ctx, models.GetAvailableRoomsParams{
		Tournamentid:  tournamentID,
		Roundnumber:   roundNumber,
		Iselimination: isElimination,
	})
	if err != nil {
		return fmt.Errorf("failed to get available rooms: %v", err)
	}

	// Get debates without rooms
	debatesWithoutRooms, err := queries.GetDebatesWithoutRooms(ctx, models.GetDebatesWithoutRoomsParams{
		Tournamentid:       tournamentID,
		Roundnumber:        roundNumber,
		Iseliminationround: isElimination,
	})
	if err != nil {
		return fmt.Errorf("failed to get debates without rooms: %v", err)
	}

	// Assign rooms to debates
	for i, debate := range debatesWithoutRooms {
		if i >= len(availableRooms) {
			return fmt.Errorf("not enough rooms for all debates")
		}

		err := queries.AssignRoomToDebate(ctx, models.AssignRoomToDebateParams{
			Debateid: debate.Debateid,
			Roomid:   availableRooms[i].Roomid,
		})
		if err != nil {
			return fmt.Errorf("failed to assign room to debate: %v", err)
		}
	}

	return nil
}

func (s *RoomService) validateAdminRole(token string) (map[string]interface{}, error) {
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

func (s *RoomService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}

func convertRooms(dbRooms []models.GetRoomsByTournamentAndRoundRow) []*debate_management.Room {
	rooms := make([]*debate_management.Room, len(dbRooms))
	for i, dbRoom := range dbRooms {
		rooms[i] = &debate_management.Room{
			RoomId:   dbRoom.Roomid,
			RoomName: dbRoom.Roomname,
			RoundStatus: []*debate_management.RoundStatus{
				{
					RoundNumber:   dbRoom.Roundnumber,
					IsElimination: dbRoom.Iselimination,
					IsOccupied:    dbRoom.Isoccupied,
				},
			},
		}
	}
	return rooms
}
