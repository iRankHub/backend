package server

import (
	"context"
	"database/sql"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	services "github.com/iRankHub/backend/internal/services/tournament_management"
	cronservices "github.com/iRankHub/backend/internal/services/tournament_management/cron_tasks"
)

type tournamentServer struct {
	tournament_management.UnimplementedTournamentServiceServer
	leagueService                 *services.LeagueService
	formatService                 *services.FormatService
	tournamentService             *services.TournamentService
	invitationService             *services.InvitationService
	billingService                *services.BillingService // Add this line
	reminderService               *cronservices.ReminderService
	tournamentCountsUpdateService *cronservices.TournamentCountsUpdateService
}

func NewTournamentServer(db *sql.DB) (tournament_management.TournamentServiceServer, error) {
	leagueService := services.NewLeagueService(db)
	formatService := services.NewFormatService(db)
	tournamentService := services.NewTournamentService(db)
	invitationService := services.NewInvitationService(db)
	billingService := services.NewBillingService(db) // Add this line

	reminderService, err := cronservices.NewReminderService(db)
	if err != nil {
		return nil, err
	}

	tournamentCountsUpdateService, err := cronservices.NewTournamentCountsUpdateService(db)
	if err != nil {
		return nil, err
	}

	server := &tournamentServer{
		leagueService:                 leagueService,
		formatService:                 formatService,
		tournamentService:             tournamentService,
		invitationService:             invitationService,
		billingService:                billingService, // Add this line
		reminderService:               reminderService,
		tournamentCountsUpdateService: tournamentCountsUpdateService,
	}

	// Start the cron services
	server.reminderService.Start()
	server.tournamentCountsUpdateService.Start()

	return server, nil
}

// a method to stop the cron services
func (s *tournamentServer) StopCronServices() {
	if s.reminderService != nil {
		s.reminderService.Stop()
	}
	if s.tournamentCountsUpdateService != nil {
		s.tournamentCountsUpdateService.Stop()
	}
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

func (s *tournamentServer) GetTournamentStats(ctx context.Context, req *tournament_management.GetTournamentStatsRequest) (*tournament_management.GetTournamentStatsResponse, error) {
	response, err := s.tournamentService.GetTournamentStats(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get tournament stats: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) GetTournamentRegistrations(ctx context.Context, req *tournament_management.GetTournamentRegistrationsRequest) (*tournament_management.GetTournamentRegistrationsResponse, error) {
	response, err := s.tournamentService.GetTournamentRegistrations(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get tournament registrations: %v", err)
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

func (s *tournamentServer) GetInvitationsByUser(ctx context.Context, req *tournament_management.GetInvitationsByUserRequest) (*tournament_management.GetInvitationsByUserResponse, error) {
	response, err := s.invitationService.GetInvitationsByUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get invitations for user: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) GetInvitationsByTournament(ctx context.Context, req *tournament_management.GetInvitationsByTournamentRequest) (*tournament_management.GetInvitationsByTournamentResponse, error) {
	response, err := s.invitationService.GetInvitationsByTournament(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get invitations for tournament: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) UpdateInvitationStatus(ctx context.Context, req *tournament_management.UpdateInvitationStatusRequest) (*tournament_management.UpdateInvitationStatusResponse, error) {
	response, err := s.invitationService.UpdateInvitationStatus(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update invitation status: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) BulkUpdateInvitationStatus(ctx context.Context, req *tournament_management.BulkUpdateInvitationStatusRequest) (*tournament_management.BulkUpdateInvitationStatusResponse, error) {
	response, err := s.invitationService.BulkUpdateInvitationStatus(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to bulk update invitation statuses: %v", err)
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

func (s *tournamentServer) CreateTournamentExpenses(ctx context.Context, req *tournament_management.CreateExpensesRequest) (*tournament_management.ExpensesResponse, error) {
	response, err := s.billingService.CreateTournamentExpenses(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create tournament expenses: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) UpdateTournamentExpenses(ctx context.Context, req *tournament_management.UpdateExpensesRequest) (*tournament_management.ExpensesResponse, error) {
	response, err := s.billingService.UpdateTournamentExpenses(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update tournament expenses: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) GetTournamentExpenses(ctx context.Context, req *tournament_management.GetExpensesRequest) (*tournament_management.ExpensesResponse, error) {
    response, err := s.billingService.GetTournamentExpenses(ctx, req)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to get tournament expenses: %v", err)
    }
    return response, nil
}

func (s *tournamentServer) CreateSchoolRegistration(ctx context.Context, req *tournament_management.CreateRegistrationRequest) (*tournament_management.RegistrationResponse, error) {
	response, err := s.billingService.CreateSchoolRegistration(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create school registration: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) UpdateSchoolRegistration(ctx context.Context, req *tournament_management.UpdateRegistrationRequest) (*tournament_management.RegistrationResponse, error) {
	response, err := s.billingService.UpdateSchoolRegistration(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update school registration: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) GetSchoolRegistration(ctx context.Context, req *tournament_management.GetRegistrationRequest) (*tournament_management.DetailedRegistrationResponse, error) {
	response, err := s.billingService.GetSchoolRegistration(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get school registration: %v", err)
	}
	return response, nil
}

func (s *tournamentServer) ListTournamentRegistrations(ctx context.Context, req *tournament_management.ListRegistrationsRequest) (*tournament_management.ListRegistrationsResponse, error) {
	response, err := s.billingService.ListTournamentRegistrations(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list tournament registrations: %v", err)
	}
	return response, nil
}
