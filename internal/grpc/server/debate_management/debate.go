package server

import (
	"context"
	"database/sql"
	"strings"

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
	teamService    *services.TeamService
	rankingService *services.RankingService
}

func NewDebateServer(db *sql.DB) (debate_management.DebateServiceServer, error) {
	return &debateServer{
		roomService:    services.NewRoomService(db),
		judgeService:   services.NewJudgeService(db),
		pairingService: services.NewPairingService(db),
		ballotService:  services.NewBallotService(db),
		teamService:    services.NewTeamService(db),
		rankingService: services.NewRankingService(db),
	}, nil
}

// Room operations
func (s *debateServer) GetRooms(ctx context.Context, req *debate_management.GetRoomsRequest) (*debate_management.GetRoomsResponse, error) {
	response, err := s.roomService.GetRooms(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get rooms: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetRoom(ctx context.Context, req *debate_management.GetRoomRequest) (*debate_management.GetRoomResponse, error) {
	response, err := s.roomService.GetRoom(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get room: %v", err)
	}
	return response, nil
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
	return judge, nil
}

func (s *debateServer) UpdateJudge(ctx context.Context, req *debate_management.UpdateJudgeRequest) (*debate_management.UpdateJudgeResponse, error) {
	response, err := s.judgeService.UpdateJudge(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update judge: %v", err)
	}
	return response, nil
}

// Pairing operations
func (s *debateServer) GetPairings(ctx context.Context, req *debate_management.GetPairingsRequest) (*debate_management.GetPairingsResponse, error) {
	response, err := s.pairingService.GetPairings(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get pairings: %v", err)
	}
	return response, nil
}

func (s *debateServer) UpdatePairings(ctx context.Context, req *debate_management.UpdatePairingsRequest) (*debate_management.UpdatePairingsResponse, error) {
	response, err := s.pairingService.UpdatePairings(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update pairings: %v", err)
	}
	return response, nil
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

func (s *debateServer) GetBallotByJudgeID(ctx context.Context, req *debate_management.GetBallotByJudgeIDRequest) (*debate_management.GetBallotByJudgeIDResponse, error) {
	ballot, err := s.ballotService.GetBallotByJudgeID(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get ballot by judge ID: %v", err)
	}
	return &debate_management.GetBallotByJudgeIDResponse{Ballot: ballot}, nil
}

func (s *debateServer) UpdateBallot(ctx context.Context, req *debate_management.UpdateBallotRequest) (*debate_management.UpdateBallotResponse, error) {
	ballot, err := s.ballotService.UpdateBallot(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update ballot: %v", err)
	}
	return &debate_management.UpdateBallotResponse{Ballot: ballot}, nil
}

// Team operations
func (s *debateServer) CreateTeam(ctx context.Context, req *debate_management.CreateTeamRequest) (*debate_management.Team, error) {
	team, err := s.teamService.CreateTeam(ctx, req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "speaker already in team"):
			return nil, status.Error(codes.InvalidArgument, "One or more speakers are already assigned to a team in this tournament.")
		case strings.Contains(err.Error(), "invalid speaker count"):
			return nil, status.Error(codes.InvalidArgument, "The number of speakers doesn't match the league requirements.")
		case strings.Contains(err.Error(), "database error"):
			return nil, status.Error(codes.Internal, "An unexpected error occurred. Please try again later.")
		default:
			return nil, status.Error(codes.Internal, "Failed to create team. Please try again.")
		}
	}
	return team, nil
}

func (s *debateServer) GetTeam(ctx context.Context, req *debate_management.GetTeamRequest) (*debate_management.Team, error) {
	team, err := s.teamService.GetTeam(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get team: %v", err)
	}
	return team, nil
}

func (s *debateServer) UpdateTeam(ctx context.Context, req *debate_management.UpdateTeamRequest) (*debate_management.Team, error) {
	team, err := s.teamService.UpdateTeam(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update team: %v", err)
	}
	return team, nil
}

func (s *debateServer) GetTeamsByTournament(ctx context.Context, req *debate_management.GetTeamsByTournamentRequest) (*debate_management.GetTeamsByTournamentResponse, error) {
	teams, err := s.teamService.GetTeamsByTournament(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get teams: %v", err)
	}
	return &debate_management.GetTeamsByTournamentResponse{Teams: teams}, nil
}

func (s *debateServer) DeleteTeam(ctx context.Context, req *debate_management.DeleteTeamRequest) (*debate_management.DeleteTeamResponse, error) {
	success, message, err := s.teamService.DeleteTeam(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete team: %v", err)
	}
	return &debate_management.DeleteTeamResponse{
		Success: success,
		Message: message,
	}, nil
}

// Algorithm integration
func (s *debateServer) GeneratePreliminaryPairings(ctx context.Context, req *debate_management.GeneratePreliminaryPairingsRequest) (*debate_management.GeneratePairingsResponse, error) {
	response, err := s.pairingService.GeneratePreliminaryPairings(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate preliminary pairings: %v", err)
	}
	return response, nil
}

func (s *debateServer) GenerateEliminationPairings(ctx context.Context, req *debate_management.GenerateEliminationPairingsRequest) (*debate_management.GeneratePairingsResponse, error) {
	response, err := s.pairingService.GenerateEliminationPairings(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate elimination pairings: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetTournamentStudentRanking(ctx context.Context, req *debate_management.TournamentRankingRequest) (*debate_management.TournamentRankingResponse, error) {
	response, err := s.rankingService.GetTournamentStudentRanking(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get tournament student ranking: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetOverallStudentRanking(ctx context.Context, req *debate_management.OverallRankingRequest) (*debate_management.OverallRankingResponse, error) {
	response, err := s.rankingService.GetOverallStudentRanking(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get overall student ranking: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetStudentOverallPerformance(ctx context.Context, req *debate_management.PerformanceRequest) (*debate_management.PerformanceResponse, error) {
	response, err := s.rankingService.GetStudentOverallPerformance(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get student overall performance: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetTournamentTeamsRanking(ctx context.Context, req *debate_management.TournamentTeamsRankingRequest) (*debate_management.TournamentTeamsRankingResponse, error) {
	response, err := s.rankingService.GetTournamentTeamsRanking(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get tournament teams ranking: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetTournamentSchoolRanking(ctx context.Context, req *debate_management.TournamentSchoolRankingRequest) (*debate_management.TournamentSchoolRankingResponse, error) {
	response, err := s.rankingService.GetTournamentSchoolRanking(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get tournament school ranking: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetOverallSchoolRanking(ctx context.Context, req *debate_management.OverallSchoolRankingRequest) (*debate_management.OverallSchoolRankingResponse, error) {
	response, err := s.rankingService.GetOverallSchoolRanking(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get overall school ranking: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetSchoolOverallPerformance(ctx context.Context, req *debate_management.SchoolPerformanceRequest) (*debate_management.SchoolPerformanceResponse, error) {
	response, err := s.rankingService.GetSchoolOverallPerformance(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get school overall performance: %v", err)
	}
	return response, nil
}
