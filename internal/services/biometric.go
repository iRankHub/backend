package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/iRankHub/backend/internal/models"
)

type BiometricService struct {
	queries *models.Queries
}

func NewBiometricService(queries *models.Queries) *BiometricService {
	return &BiometricService{queries: queries}
}

func (s *BiometricService) GenerateBiometricToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *BiometricService) EnableBiometricLogin(ctx context.Context, userID int32) (string, error) {
    token, err := s.GenerateBiometricToken()
    if err != nil {
        return "", err
    }

    err = s.queries.SetBiometricToken(ctx, models.SetBiometricTokenParams{
        Userid:         userID,
        BiometricToken: sql.NullString{String: token, Valid: true},
    })
    if err != nil {
        return "", fmt.Errorf("failed to set biometric token: %v", err)
    }

    return token, nil
}

func (s *BiometricService) VerifyBiometricToken(ctx context.Context, token string) (*models.User, error) {
    user, err := s.queries.GetUserByBiometricToken(ctx, sql.NullString{String: token, Valid: true})
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("invalid biometric token")
        }
        return nil, fmt.Errorf("failed to verify biometric token: %v", err)
    }
    return &user, nil
}