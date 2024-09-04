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

func (s *RoomService) GetRooms(ctx context.Context, req *debate_management.GetRoomsRequest) (*debate_management.GetRoomsResponse, error) {
    if err := s.validateAuthentication(req.GetToken()); err != nil {
        return nil, err
    }

    queries := models.New(s.db)
    dbRooms, err := queries.GetRoomsByTournament(ctx, sql.NullInt32{Int32: req.GetTournamentId(), Valid: true})
    if err != nil {
        return nil, fmt.Errorf("failed to get rooms: %v", err)
    }

    roomStatuses, err := s.convertRooms(ctx, req.GetTournamentId(), dbRooms)
    if err != nil {
        return nil, fmt.Errorf("failed to convert rooms: %v", err)
    }

    return &debate_management.GetRoomsResponse{
        Rooms: roomStatuses,
    }, nil
}

func (s *RoomService) GetRoom(ctx context.Context, req *debate_management.GetRoomRequest) (*debate_management.GetRoomResponse, error) {
	if err := s.validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)
	room, err := queries.GetRoomByID(ctx, req.GetRoomId())
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %v", err)
	}

	tournament, err := queries.GetTournamentByID(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament: %v", err)
	}

	preliminaryRounds := make([]*debate_management.RoundStatus, tournament.Numberofpreliminaryrounds)
	for i := 1; i <= int(tournament.Numberofpreliminaryrounds); i++ {
		status, err := s.getRoomStatusForRound(ctx, req.GetTournamentId(), room.Roomid, int32(i), false)
		if err != nil {
			return nil, fmt.Errorf("failed to get preliminary status for round %d: %v", i, err)
		}
		preliminaryRounds[i-1] = &debate_management.RoundStatus{
			Round:  int32(i),
			Status: status,
		}
	}

	eliminationRounds := make([]*debate_management.RoundStatus, tournament.Numberofeliminationrounds)
	for i := 1; i <= int(tournament.Numberofeliminationrounds); i++ {
		status, err := s.getRoomStatusForRound(ctx, req.GetTournamentId(), room.Roomid, int32(i), true)
		if err != nil {
			return nil, fmt.Errorf("failed to get elimination status for round %d: %v", i, err)
		}
		eliminationRounds[i-1] = &debate_management.RoundStatus{
			Round:  int32(i),
			Status: status,
		}
	}

	return convertSingleRoom(room, preliminaryRounds, eliminationRounds), nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, req *debate_management.UpdateRoomRequest) (*debate_management.Room, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	queries := models.New(s.db)

	updatedRoom, err := queries.UpdateRoom(ctx, models.UpdateRoomParams{
		Roomid:   req.GetRoom().GetRoomId(),
		Roomname: req.GetRoom().GetRoomName(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update room: %v", err)
	}

	return convertRoom(updatedRoom), nil
}

func (s *RoomService) getRoomStatus(ctx context.Context, tournamentID, roomID int32, isElimination bool) (string, error) {
	queries := models.New(s.db)
	debates, err := queries.GetDebatesByRoomAndTournament(ctx, models.GetDebatesByRoomAndTournamentParams{
		Tournamentid:       tournamentID,
		Roomid:             roomID,
		Iseliminationround: isElimination,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get debates: %v", err)
	}

	if len(debates) > 0 {
		return "occupied", nil
	}
	return "available", nil
}

func (s *RoomService) getRoomStatusForRound(ctx context.Context, tournamentID, roomID, roundNumber int32, isElimination bool) (string, error) {
	queries := models.New(s.db)
	debate, err := queries.GetDebateByRoomAndRound(ctx, models.GetDebateByRoomAndRoundParams{
		Tournamentid:       tournamentID,
		Roomid:             roomID,
		Roundnumber:        roundNumber,
		Iseliminationround: isElimination,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return "available", nil
		}
		return "", fmt.Errorf("failed to get debate: %v", err)
	}

	if debate.Debateid != 0 {
		return "occupied", nil
	}
	return "available", nil
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

func (s *RoomService) convertRooms(ctx context.Context, tournamentID int32, dbRooms []models.Room) ([]*debate_management.RoomStatus, error) {
    rooms := make([]*debate_management.RoomStatus, len(dbRooms))
    for i, dbRoom := range dbRooms {
        preliminaryStatus, err := s.getRoomStatus(ctx, tournamentID, dbRoom.Roomid, false)
        if err != nil {
            return nil, fmt.Errorf("failed to get preliminary status for room %d: %v", dbRoom.Roomid, err)
        }

        eliminationStatus, err := s.getRoomStatus(ctx, tournamentID, dbRoom.Roomid, true)
        if err != nil {
            return nil, fmt.Errorf("failed to get elimination status for room %d: %v", dbRoom.Roomid, err)
        }

        rooms[i] = &debate_management.RoomStatus{
            RoomId:      dbRoom.Roomid,
            Preliminary: preliminaryStatus,
            Elimination: eliminationStatus,
        }
    }
    return rooms, nil
}

func convertSingleRoom(dbRoom models.GetRoomByIDRow, preliminaryRounds, eliminationRounds []*debate_management.RoundStatus) *debate_management.GetRoomResponse {
    return &debate_management.GetRoomResponse{
        RoomId:      dbRoom.Roomid,
        Name:        dbRoom.Roomname,
        Preliminary: preliminaryRounds,
        Elimination: eliminationRounds,
    }
}

func convertRoom(dbRoom models.Room) *debate_management.Room {
    return &debate_management.Room{
        RoomId:   dbRoom.Roomid,
        RoomName: dbRoom.Roomname,
        Location: dbRoom.Location,
        Capacity: dbRoom.Capacity,
    }
}