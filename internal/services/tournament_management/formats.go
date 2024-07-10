package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type FormatService struct {
	db *sql.DB
}

func NewFormatService(db *sql.DB) *FormatService {
	return &FormatService{db: db}
}

func (s *FormatService) CreateTournamentFormat(ctx context.Context, req *tournament_management.CreateTournamentFormatRequest) (*tournament_management.TournamentFormat, error) {
	// Check if the user is an admin
	claims, err := utils.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	userRole := claims["user_role"].(string)
	if userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can create tournament formats")
	}

	// Create the tournament format
	format, err := s.db.CreateTournamentFormat(ctx, models.CreateTournamentFormatParams{
		FormatName:      req.GetFormatName(),
		Description:     req.GetDescription(),
		SpeakersPerTeam: int32(req.GetSpeakersPerTeam()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tournament format: %v", err)
	}

	return &tournament_management.TournamentFormat{
		FormatId:        format.FormatID,
		FormatName:      format.FormatName,
		Description:     format.Description,
		SpeakersPerTeam: int32(format.SpeakersPerTeam),
	}, nil
}

func (s *FormatService) GetTournamentFormat(ctx context.Context, req *tournament_management.GetTournamentFormatRequest) (*tournament_management.TournamentFormat, error) {
	// Get the tournament format by ID
	format, err := s.db.GetTournamentFormatByID(ctx, req.GetFormatId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament format: %v", err)
	}

	return &tournament_management.TournamentFormat{
		FormatId:        format.FormatID,
		FormatName:      format.FormatName,
		Description:     format.Description,
		SpeakersPerTeam: int32(format.SpeakersPerTeam),
	}, nil
}

func (s *FormatService) ListTournamentFormats(ctx context.Context, req *tournament_management.ListTournamentFormatsRequest) (*tournament_management.ListTournamentFormatsResponse, error) {
	// List tournament formats with pagination
	formats, err := s.db.ListTournamentFormatsPaginated(ctx, models.ListTournamentFormatsPaginatedParams{
		Limit:  int32(req.GetPageSize()),
		Offset: int32(req.GetPageToken()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tournament formats: %v", err)
	}

	// Construct the TournamentFormat responses
	formatResponses := make([]*tournament_management.TournamentFormat, len(formats))
	for i, format := range formats {
		formatResponses[i] = &tournament_management.TournamentFormat{
			FormatId:        format.FormatID,
			FormatName:      format.FormatName,
			Description:     format.Description,
			SpeakersPerTeam: int32(format.SpeakersPerTeam),
		}
	}

	return &tournament_management.ListTournamentFormatsResponse{
		Formats:       formatResponses,
		NextPageToken: int32(req.GetPageToken()) + int32(req.GetPageSize()),
	}, nil
}

func (s *FormatService) UpdateTournamentFormat(ctx context.Context, req *tournament_management.UpdateTournamentFormatRequest) (*tournament_management.TournamentFormat, error) {
	// Check if the user is an admin
	claims, err := utils.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	userRole := claims["user_role"].(string)
	if userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can update tournament formats")
	}

	// Update the tournament format details
	updatedFormat, err := s.db.UpdateTournamentFormatDetails(ctx, models.UpdateTournamentFormatDetailsParams{
		FormatID:        req.GetFormatId(),
		FormatName:      req.GetFormatName(),
		Description:     req.GetDescription(),
		SpeakersPerTeam: int32(req.GetSpeakersPerTeam()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update tournament format details: %v", err)
	}

	return &tournament_management.TournamentFormat{
		FormatId:        updatedFormat.FormatID,
		FormatName:      updatedFormat.FormatName,
		Description:     updatedFormat.Description,
		SpeakersPerTeam: int32(updatedFormat.SpeakersPerTeam),
	}, nil
}

func (s *FormatService) DeleteTournamentFormat(ctx context.Context, req *tournament_management.DeleteTournamentFormatRequest) (*tournament_management.Empty, error) {
	// Check if the user is an admin
	claims, err := utils.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	userRole := claims["user_role"].(string)
	if userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can delete tournament formats")
	}

	// Delete the tournament format by ID
	err = s.db.DeleteTournamentFormatByID(ctx, req.GetFormatId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete tournament format: %v", err)
	}

	return &tournament_management.Empty{}, nil
}
