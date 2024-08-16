package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"

)

type SchoolService struct {
	db *sql.DB
}

func NewSchoolsManagementService(db *sql.DB) *SchoolService {
	return &SchoolService{db: db}
}

func (s *SchoolService) GetSchools(ctx context.Context, token string, page, pageSize int32) ([]models.School, int32, error) {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid token: %v", err)
	}

	queries := models.New(s.db)
	schools, err := queries.GetSchoolsPaginated(ctx, models.GetSchoolsPaginatedParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch schools: %v", err)
	}

	totalCount, err := queries.GetTotalSchoolCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total school count: %v", err)
	}

	return schools, int32(totalCount), nil
}

func (s *SchoolService) GetSchoolsNoAuth(ctx context.Context, page, pageSize int32) ([]models.School, int32, error) {
    queries := models.New(s.db)
    schools, err := queries.GetSchoolsPaginated(ctx, models.GetSchoolsPaginatedParams{
        Limit:  pageSize,
        Offset: (page - 1) * pageSize,
    })
    if err != nil {
        return nil, 0, fmt.Errorf("failed to fetch schools: %v", err)
    }

    totalCount, err := queries.GetTotalSchoolCount(ctx)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get total school count: %v", err)
    }

    return schools, int32(totalCount), nil
}