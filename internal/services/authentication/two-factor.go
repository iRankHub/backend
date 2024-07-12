package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"fmt"

	"github.com/pquerna/otp/totp"

	"github.com/iRankHub/backend/internal/models"
)

type TwoFactorService struct {
	db *sql.DB
}

func NewTwoFactorService(db *sql.DB) *TwoFactorService {
	return &TwoFactorService{db: db}
}

func (s *TwoFactorService) GenerateSecret() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return base32.StdEncoding.EncodeToString(bytes), nil
}

func (s *TwoFactorService) EnableTwoFactor(ctx context.Context, userID int32) (string, string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	secret, err := s.GenerateSecret()
	if err != nil {
		return "", "", err
	}

	err = queries.UpdateUserTwoFactorSecret(ctx, models.UpdateUserTwoFactorSecretParams{
		Userid:          userID,
		TwoFactorSecret: sql.NullString{String: secret, Valid: true},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to update user two factor secret: %v", err)
	}

	err = queries.EnableTwoFactor(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to enable two-factor authentication: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "iRankHub",
		AccountName: fmt.Sprintf("user-%d", userID),
		Secret:      []byte(secret),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate TOTP key: %v", err)
	}

	return secret, key.URL(), nil
}

func (s *TwoFactorService) VerifyAndEnableTwoFactor(ctx context.Context, userID int32, code string) (bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %v", err)
	}

	if !user.TwoFactorSecret.Valid || !s.ValidateCode(user.TwoFactorSecret.String, code) {
		return false, fmt.Errorf("invalid two-factor code")
	}

	err = queries.EnableTwoFactor(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to enable two-factor authentication: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return true, nil
}

func (s *TwoFactorService) DisableTwoFactor(ctx context.Context, userID int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	err = queries.DisableTwoFactor(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to disable two-factor authentication: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *TwoFactorService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}