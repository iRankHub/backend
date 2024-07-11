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
	if err := s.validateAdminRole(ctx); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)
	format, err := queries.CreateTournamentFormat(ctx, models.CreateTournamentFormatParams{
		Formatname:      req.GetFormatName(),
		Description:     sql.NullString{String: req.GetDescription(), Valid: req.GetDescription() != ""},
		Speakersperteam: req.GetSpeakersPerTeam(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tournament format: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.TournamentFormat{
		FormatId:        int32(format.Formatid),
		FormatName:      format.Formatname,
		Description:     format.Description.String,
		SpeakersPerTeam: format.Speakersperteam,
	}, nil
}

func (s *FormatService) GetTournamentFormat(ctx context.Context, req *tournament_management.GetTournamentFormatRequest) (*tournament_management.TournamentFormat, error) {
	queries := models.New(s.db)
	format, err := queries.GetTournamentFormatByID(ctx, int32(req.GetFormatId()))
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament format: %v", err)
	}

	return &tournament_management.TournamentFormat{
		FormatId:        int32(format.Formatid),
		FormatName:      format.Formatname,
		Description:     format.Description.String,
		SpeakersPerTeam: format.Speakersperteam,
	}, nil
}

func (s *FormatService) ListTournamentFormats(ctx context.Context, req *tournament_management.ListTournamentFormatsRequest) (*tournament_management.ListTournamentFormatsResponse, error) {
	queries := models.New(s.db)
	formats, err := queries.ListTournamentFormatsPaginated(ctx, models.ListTournamentFormatsPaginatedParams{
		Limit:  req.GetPageSize(),
		Offset: req.GetPageToken(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tournament formats: %v", err)
	}

	formatResponses := make([]*tournament_management.TournamentFormat, len(formats))
	for i, format := range formats {
		formatResponses[i] = &tournament_management.TournamentFormat{
			FormatId:        int32(format.Formatid),
			FormatName:      format.Formatname,
			Description:     format.Description.String,
			SpeakersPerTeam: format.Speakersperteam,
		}
	}

	return &tournament_management.ListTournamentFormatsResponse{
		Formats:       formatResponses,
		NextPageToken: req.GetPageToken() + req.GetPageSize(),
	}, nil
}

func (s *FormatService) UpdateTournamentFormat(ctx context.Context, req *tournament_management.UpdateTournamentFormatRequest) (*tournament_management.TournamentFormat, error) {
	if err := s.validateAdminRole(ctx); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)
	updatedFormat, err := queries.UpdateTournamentFormatDetails(ctx, models.UpdateTournamentFormatDetailsParams{
		Formatid:        int32(req.GetFormatId()),
		Formatname:      req.GetFormatName(),
		Description:     sql.NullString{String: req.GetDescription(), Valid: req.GetDescription() != ""},
		Speakersperteam: req.GetSpeakersPerTeam(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update tournament format details: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.TournamentFormat{
		FormatId:        int32(updatedFormat.Formatid),
		FormatName:      updatedFormat.Formatname,
		Description:     updatedFormat.Description.String,
		SpeakersPerTeam: updatedFormat.Speakersperteam,
	}, nil
}

func (s *FormatService) DeleteTournamentFormat(ctx context.Context, req *tournament_management.DeleteTournamentFormatRequest) error {
	if err := s.validateAdminRole(ctx); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)
	err = queries.DeleteTournamentFormatByID(ctx, int32(req.GetFormatId()))
	if err != nil {
		return fmt.Errorf("failed to delete tournament format: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *FormatService) validateAdminRole(ctx context.Context) error {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("failed to validate token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return fmt.Errorf("unauthorized: only admins can perform this action")
	}

	return nil
}