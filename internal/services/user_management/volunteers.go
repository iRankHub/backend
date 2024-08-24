package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type VolunteerService struct {
	db *sql.DB
}

func NewVolunteersManagementService(db *sql.DB) *VolunteerService {
	return &VolunteerService{db: db}
}

func (s *VolunteerService) GetVolunteers(ctx context.Context, token string, page, pageSize int32) ([]models.Volunteer, int32, error) {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid token: %v", err)
	}

	queries := models.New(s.db)
	volunteers, err := queries.GetVolunteersPaginated(ctx, models.GetVolunteersPaginatedParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch volunteers: %v", err)
	}

	totalCount, err := queries.GetTotalVolunteerCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total volunteer count: %v", err)
	}

	return volunteers, int32(totalCount), nil
}
