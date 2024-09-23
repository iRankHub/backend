package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type CountryService struct {
	db *sql.DB
}

func NewCountryManagementService(db *sql.DB) *CountryService {
	return &CountryService{db: db}
}

func (s *CountryService) GetCountries(ctx context.Context, token string) ([]models.Countrycode, error) {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return s.getCountriesInternal(ctx)
}

func (s *CountryService) GetCountriesNoAuth(ctx context.Context) ([]models.Countrycode, error) {
	return s.getCountriesInternal(ctx)
}

func (s *CountryService) getCountriesInternal(ctx context.Context) ([]models.Countrycode, error) {
	queries := models.New(s.db)
	return queries.GetAllCountries(ctx)
}
