package server

import (
	"context"
	"database/sql"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	services "github.com/iRankHub/backend/internal/services/tournament_management"

)

type tournamentServer struct {
	tournament_management.UnimplementedTournamentServiceServer
	leagueService     *services.LeagueService
	formatService     *services.FormatService
	tournamentService *services.TournamentService
	invitationService *services.InvitationService
}

func NewTournamentServer(db *sql.DB) (tournament_management.TournamentServiceServer, error) {
	leagueService := services.NewLeagueService(db)
	formatService := services.NewFormatService(db)
	tournamentService := services.NewTournamentService(db)
	invitationService := services.NewInvitationService(db)

	return &tournamentServer{
		leagueService:     leagueService,
		formatService:     formatService,
		tournamentService: tournamentService,
		invitationService: invitationService,
	}, nil
}

func (s *tournamentServer) CreateLeague(ctx context.Context, req *tournament_management.CreateLeagueRequest) (*tournament_management.CreateLeagueResponse, error) {
	league, err := s.leagueService.CreateLeague(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create league: %v", err)
	}
	return &tournament_management.CreateLeagueResponse{
		League: league,
	}, nil
}

func (s *tournamentServer) GetLeague(ctx context.Context, req *tournament_management.GetLeagueRequest) (*tournament_management.GetLeagueResponse, error) {
	league, err := s.leagueService.GetLeague(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get league: %v", err)
	}
	return &tournament_management.GetLeagueResponse{
		League: league,
	}, nil
}

func (s *tournamentServer) ListLeagues(ctx context.Context, req *tournament_management.ListLeaguesRequest) (*tournament_management.ListLeaguesResponse, error) {
	response, err := s.leagueService.ListLeagues(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list leagues: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) UpdateLeague(ctx context.Context, req *tournament_management.UpdateLeagueRequest) (*tournament_management.UpdateLeagueResponse, error) {
	league, err := s.leagueService.UpdateLeague(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update league: %v", err)
	}
	return &tournament_management.UpdateLeagueResponse{
		League: league,
	}, nil
}

func (s *tournamentServer) DeleteLeague(ctx context.Context, req *tournament_management.DeleteLeagueRequest) (*tournament_management.DeleteLeagueResponse, error) {
	success, err := s.leagueService.DeleteLeague(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete league: %v", err)
	}
	return &tournament_management.DeleteLeagueResponse{
		Success: success.Success,
		Message: "League deleted successfully",
	}, nil
}

func (s *tournamentServer) CreateTournamentFormat(ctx context.Context, req *tournament_management.CreateTournamentFormatRequest) (*tournament_management.CreateTournamentFormatResponse, error) {
	format, err := s.formatService.CreateTournamentFormat(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create tournament format: %v", err)
	}
	return &tournament_management.CreateTournamentFormatResponse{
		Format: format,
	}, nil
}

func (s *tournamentServer) GetTournamentFormat(ctx context.Context, req *tournament_management.GetTournamentFormatRequest) (*tournament_management.GetTournamentFormatResponse, error) {
	format, err := s.formatService.GetTournamentFormat(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get tournament format: %v", err)
	}
	return &tournament_management.GetTournamentFormatResponse{
		Format: format,
	}, nil
}

func (s *tournamentServer) ListTournamentFormats(ctx context.Context, req *tournament_management.ListTournamentFormatsRequest) (*tournament_management.ListTournamentFormatsResponse, error) {
	response, err := s.formatService.ListTournamentFormats(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list tournament formats: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) UpdateTournamentFormat(ctx context.Context, req *tournament_management.UpdateTournamentFormatRequest) (*tournament_management.UpdateTournamentFormatResponse, error) {
	format, err := s.formatService.UpdateTournamentFormat(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update tournament format: %v", err)
	}
	return &tournament_management.UpdateTournamentFormatResponse{
		Format: format,
	}, nil
}

func (s *tournamentServer) DeleteTournamentFormat(ctx context.Context, req *tournament_management.DeleteTournamentFormatRequest) (*tournament_management.DeleteTournamentFormatResponse, error) {
	err := s.formatService.DeleteTournamentFormat(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete tournament format: %v", err)
	}
	return &tournament_management.DeleteTournamentFormatResponse{
		Success: true,
		Message: "Tournament format deleted successfully",
	}, nil
}

func (s *tournamentServer) CreateTournament(ctx context.Context, req *tournament_management.CreateTournamentRequest) (*tournament_management.CreateTournamentResponse, error) {
    tournament, err := s.tournamentService.CreateTournament(ctx, req)
    if err != nil {
        log.Printf("Error creating tournament: %v", err)
        return nil, status.Errorf(codes.Internal, "Failed to create tournament: %v", err)
    }
    return &tournament_management.CreateTournamentResponse{
        Tournament: tournament,
    }, nil
}

func (s *tournamentServer) GetTournament(ctx context.Context, req *tournament_management.GetTournamentRequest) (*tournament_management.GetTournamentResponse, error) {
	tournament, err := s.tournamentService.GetTournament(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get tournament: %v", err)
	}
	return &tournament_management.GetTournamentResponse{
		Tournament: tournament,
	}, nil
}

func (s *tournamentServer) ListTournaments(ctx context.Context, req *tournament_management.ListTournamentsRequest) (*tournament_management.ListTournamentsResponse, error) {
	response, err := s.tournamentService.ListTournaments(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list tournaments: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) UpdateTournament(ctx context.Context, req *tournament_management.UpdateTournamentRequest) (*tournament_management.UpdateTournamentResponse, error) {
	tournament, err := s.tournamentService.UpdateTournament(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update tournament: %v", err)
	}
	return &tournament_management.UpdateTournamentResponse{
		Tournament: tournament,
	}, nil
}

func (s *tournamentServer) DeleteTournament(ctx context.Context, req *tournament_management.DeleteTournamentRequest) (*tournament_management.DeleteTournamentResponse, error) {
	success, err := s.tournamentService.DeleteTournament(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete tournament: %v", err)
	}
	return &tournament_management.DeleteTournamentResponse{
		Success: success.Success,
		Message: "Tournament deleted successfully",
	}, nil
}

func (s *tournamentServer) AcceptInvitation(ctx context.Context, req *tournament_management.AcceptInvitationRequest) (*tournament_management.AcceptInvitationResponse, error) {
	response, err := s.invitationService.AcceptInvitation(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to accept invitation: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) DeclineInvitation(ctx context.Context, req *tournament_management.DeclineInvitationRequest) (*tournament_management.DeclineInvitationResponse, error) {
	response, err := s.invitationService.DeclineInvitation(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to decline invitation: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) BulkAcceptInvitations(ctx context.Context, req *tournament_management.BulkAcceptInvitationsRequest) (*tournament_management.BulkAcceptInvitationsResponse, error) {
	response, err := s.invitationService.BulkAcceptInvitations(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to bulk accept invitations: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) BulkDeclineInvitations(ctx context.Context, req *tournament_management.BulkDeclineInvitationsRequest) (*tournament_management.BulkDeclineInvitationsResponse, error) {
	response, err := s.invitationService.BulkDeclineInvitations(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to bulk decline invitations: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) GetInvitationsByUser(ctx context.Context, req *tournament_management.GetInvitationsByUserRequest) (*tournament_management.GetInvitationsByUserResponse, error) {
    response, err := s.invitationService.GetInvitationsByUser(ctx, req)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to get invitations for user: %v", err)
    }
    return response, nil
}

func (s *tournamentServer) GetAllInvitations(ctx context.Context, req *tournament_management.GetAllInvitationsRequest) (*tournament_management.GetAllInvitationsResponse, error) {
	response, err := s.invitationService.GetAllInvitations(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get all invitations: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) ResendInvitation(ctx context.Context, req *tournament_management.ResendInvitationRequest) (*tournament_management.ResendInvitationResponse, error) {
	response, err := s.invitationService.ResendInvitation(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to resend invitation: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) BulkResendInvitations(ctx context.Context, req *tournament_management.BulkResendInvitationsRequest) (*tournament_management.BulkResendInvitationsResponse, error) {
	response, err := s.invitationService.BulkResendInvitations(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to bulk resend invitations: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) GetInvitationStatus(ctx context.Context, req *tournament_management.GetInvitationStatusRequest) (*tournament_management.GetInvitationStatusResponse, error) {
	log.Printf("GetInvitationStatus called with invitation ID: %d", req.GetInvitationId())
	response, err := s.invitationService.GetInvitationStatus(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get invitation status: %v", err)
	}
	return response, nil
}