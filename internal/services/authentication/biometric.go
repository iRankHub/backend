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
	db *sql.DB
}

func NewBiometricService(db *sql.DB) *BiometricService {
	return &BiometricService{db: db}
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
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return "", fmt.Errorf("failed to start transaction: %v", err)
    }
    defer tx.Rollback()

    queries := models.New(tx)

    token, err := s.GenerateBiometricToken()
    if err != nil {
        return "", err
    }

    err = queries.SetBiometricToken(ctx, models.SetBiometricTokenParams{
        Userid:         userID,
        BiometricToken: sql.NullString{String: token, Valid: true},
    })
    if err != nil {
        return "", fmt.Errorf("failed to set biometric token: %v", err)
    }

    if err := tx.Commit(); err != nil {
        return "", fmt.Errorf("failed to commit transaction: %v", err)
    }

    return token, nil
}

func (s *BiometricService) VerifyBiometricToken(ctx context.Context, token string) (*models.User, error) {
    tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %v", err)
    }
    defer tx.Rollback()

    queries := models.New(tx)

    user, err := queries.GetUserByBiometricToken(ctx, sql.NullString{String: token, Valid: true})
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("invalid biometric token")
        }
        return nil, fmt.Errorf("failed to verify biometric token: %v", err)
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return &user, nil
}