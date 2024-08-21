package server

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iRankHub/backend/internal/grpc/proto/debate_management"
	services "github.com/iRankHub/backend/internal/services/debate_management"

)

type debateServer struct {
	debate_management.UnimplementedDebateServiceServer
	roomService    *services.RoomService
	judgeService   *services.JudgeService
	pairingService *services.PairingService
	ballotService  *services.BallotService
}

func NewDebateServer(db *sql.DB) (debate_management.DebateServiceServer, error) {
	return &debateServer{
		roomService:    services.NewRoomService(db),
		judgeService:   services.NewJudgeService(db),
		pairingService: services.NewPairingService(db),
		ballotService:  services.NewBallotService(db),
	}, nil
}

// Room operations
func (s *debateServer) GetRooms(ctx context.Context, req *debate_management.GetRoomsRequest) (*debate_management.GetRoomsResponse, error) {
	rooms, err := s.roomService.GetRooms(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get rooms: %v", err)
	}
	return &debate_management.GetRoomsResponse{Rooms: rooms}, nil
}

func (s *debateServer) GetRoom(ctx context.Context, req *debate_management.GetRoomRequest) (*debate_management.GetRoomResponse, error) {
	room, err := s.roomService.GetRoom(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get room: %v", err)
	}
	return &debate_management.GetRoomResponse{Room: room}, nil
}

func (s *debateServer) UpdateRoom(ctx context.Context, req *debate_management.UpdateRoomRequest) (*debate_management.UpdateRoomResponse, error) {
	room, err := s.roomService.UpdateRoom(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update room: %v", err)
	}
	return &debate_management.UpdateRoomResponse{Room: room}, nil
}

// Judge operations
func (s *debateServer) GetJudges(ctx context.Context, req *debate_management.GetJudgesRequest) (*debate_management.GetJudgesResponse, error) {
	judges, err := s.judgeService.GetJudges(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get judges: %v", err)
	}
	return &debate_management.GetJudgesResponse{Judges: judges}, nil
}

func (s *debateServer) GetJudge(ctx context.Context, req *debate_management.GetJudgeRequest) (*debate_management.GetJudgeResponse, error) {
	judge, err := s.judgeService.GetJudge(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get judge: %v", err)
	}
	return &debate_management.GetJudgeResponse{Judge: judge}, nil
}

// Pairing operations
func (s *debateServer) GetPairings(ctx context.Context, req *debate_management.GetPairingsRequest) (*debate_management.GetPairingsResponse, error) {
	pairings, err := s.pairingService.GetPairings(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get pairings: %v", err)
	}
	return &debate_management.GetPairingsResponse{Pairings: pairings}, nil
}

func (s *debateServer) GetPairing(ctx context.Context, req *debate_management.GetPairingRequest) (*debate_management.GetPairingResponse, error) {
	pairing, err := s.pairingService.GetPairing(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get pairing: %v", err)
	}
	return &debate_management.GetPairingResponse{Pairing: pairing}, nil
}

func (s *debateServer) UpdatePairing(ctx context.Context, req *debate_management.UpdatePairingRequest) (*debate_management.UpdatePairingResponse, error) {
	pairing, err := s.pairingService.UpdatePairing(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update pairing: %v", err)
	}
	return &debate_management.UpdatePairingResponse{Pairing: pairing}, nil
}

// Ballot operations
func (s *debateServer) GetBallots(ctx context.Context, req *debate_management.GetBallotsRequest) (*debate_management.GetBallotsResponse, error) {
	ballots, err := s.ballotService.GetBallots(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get ballots: %v", err)
	}
	return &debate_management.GetBallotsResponse{Ballots: ballots}, nil
}

func (s *debateServer) GetBallot(ctx context.Context, req *debate_management.GetBallotRequest) (*debate_management.GetBallotResponse, error) {
	ballot, err := s.ballotService.GetBallot(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get ballot: %v", err)
	}
	return &debate_management.GetBallotResponse{Ballot: ballot}, nil
}

func (s *debateServer) UpdateBallot(ctx context.Context, req *debate_management.UpdateBallotRequest) (*debate_management.UpdateBallotResponse, error) {
	ballot, err := s.ballotService.UpdateBallot(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update ballot: %v", err)
	}
	return &debate_management.UpdateBallotResponse{Ballot: ballot}, nil
}

// Algorithm integration
func (s *debateServer) GeneratePairings(ctx context.Context, req *debate_management.GeneratePairingsRequest) (*debate_management.GeneratePairingsResponse, error) {
	pairings, err := s.pairingService.GeneratePairings(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate pairings: %v", err)
	}
	return &debate_management.GeneratePairingsResponse{Pairings: pairings}, nil
}

func (s *debateServer) AssignJudges(ctx context.Context, req *debate_management.AssignJudgesRequest) (*debate_management.AssignJudgesResponse, error) {
	pairings, err := s.judgeService.AssignJudges(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to assign judges: %v", err)
	}
	return &debate_management.AssignJudgesResponse{Pairings: pairings}, nil
}