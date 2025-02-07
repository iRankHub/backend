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
	roomService     *services.RoomService
	judgeService    *services.JudgeService
	pairingService  *services.PairingService
	ballotService   *services.BallotService
	teamService     *services.TeamService
	rankingService  *services.RankingService
	feedbackService *services.FeedbackService
}

func NewDebateServer(db *sql.DB) (debate_management.DebateServiceServer, error) {
	return &debateServer{
		roomService:     services.NewRoomService(db),
		judgeService:    services.NewJudgeService(db),
		pairingService:  services.NewPairingService(db),
		ballotService:   services.NewBallotService(db),
		teamService:     services.NewTeamService(db),
		rankingService:  services.NewRankingService(db),
		feedbackService: services.NewFeedbackService(db),
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

func (s *debateServer) GetStudentTournamentStats(ctx context.Context, req *debate_management.StudentTournamentStatsRequest) (*debate_management.StudentTournamentStatsResponse, error) {
	response, err := s.rankingService.GetStudentTournamentStats(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get student tournament stats: %v", err)
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

func (s *debateServer) GetVolunteerTournamentStats(ctx context.Context, req *debate_management.VolunteerTournamentStatsRequest) (*debate_management.VolunteerTournamentStatsResponse, error) {
	response, err := s.rankingService.GetVolunteerTournamentStats(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get volunteer tournament stats: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetStudentFeedback(ctx context.Context, req *debate_management.GetStudentFeedbackRequest) (*debate_management.GetStudentFeedbackResponse, error) {
	response, err := s.feedbackService.GetStudentFeedback(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get student feedback: %v", err)
	}
	return response, nil
}

func (s *debateServer) SubmitJudgeFeedback(ctx context.Context, req *debate_management.SubmitJudgeFeedbackRequest) (*debate_management.SubmitJudgeFeedbackResponse, error) {
	response, err := s.feedbackService.SubmitJudgeFeedback(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to submit judge feedback: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetJudgeFeedback(ctx context.Context, req *debate_management.GetJudgeFeedbackRequest) (*debate_management.GetJudgeFeedbackResponse, error) {
	response, err := s.feedbackService.GetJudgeFeedback(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get judge feedback: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetVolunteerRanking(ctx context.Context, req *debate_management.GetVolunteerRankingRequest) (*debate_management.GetVolunteerRankingResponse, error) {
	response, err := s.rankingService.GetVolunteerRanking(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get volunteer ranking: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetVolunteerPerformance(ctx context.Context, req *debate_management.GetVolunteerPerformanceRequest) (*debate_management.GetVolunteerPerformanceResponse, error) {
	response, err := s.rankingService.GetVolunteerPerformance(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get volunteer performance: %v", err)
	}
	return response, nil
}

func (s *debateServer) MarkStudentFeedbackAsRead(ctx context.Context, req *debate_management.MarkFeedbackAsReadRequest) (*debate_management.MarkFeedbackAsReadResponse, error) {
	response, err := s.feedbackService.MarkStudentFeedbackAsRead(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to mark student feedback as read: %v", err)
	}
	return response, nil
}

func (s *debateServer) MarkJudgeFeedbackAsRead(ctx context.Context, req *debate_management.MarkFeedbackAsReadRequest) (*debate_management.MarkFeedbackAsReadResponse, error) {
	response, err := s.feedbackService.MarkJudgeFeedbackAsRead(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to mark judge feedback as read: %v", err)
	}
	return response, nil
}

func (s *debateServer) GetTournamentVolunteerRanking(ctx context.Context, req *debate_management.TournamentVolunteerRankingRequest) (*debate_management.TournamentVolunteerRankingResponse, error) {
	response, err := s.rankingService.GetTournamentVolunteerRanking(ctx, req)
	if err != nil {
		// Check for specific error types
		switch {
		case strings.Contains(err.Error(), "invalid token"):
			return nil, status.Errorf(codes.Unauthenticated, "Authentication failed: %v", err)
		case strings.Contains(err.Error(), "tournament not found"):
			return nil, status.Errorf(codes.NotFound, "Tournament not found: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "Failed to get tournament volunteer ranking: %v", err)
		}
	}
	return response, nil
}

func (s *debateServer) SetRankingVisibility(ctx context.Context, req *debate_management.SetRankingVisibilityRequest) (*debate_management.SetRankingVisibilityResponse, error) {
	response, err := s.rankingService.SetRankingVisibility(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to set ranking visibility: %v", err)
	}
	return response, nil
}
