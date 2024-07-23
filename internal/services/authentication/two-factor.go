package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/iRankHub/backend/internal/models"
	emails "github.com/iRankHub/backend/internal/utils/emails"
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

func (s *TwoFactorService) GenerateOTP(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now().Add(15*time.Minute))
}

func (s *TwoFactorService) GenerateTwoFactorOTP(ctx context.Context, email string) error {
	queries := models.New(s.db)

	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if !user.TwoFactorSecret.Valid {
		return fmt.Errorf("two-factor authentication not enabled for this user")
	}

	otp, err := s.GenerateOTP(user.TwoFactorSecret.String)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %v", err)
	}

	err = emails.SendTwoFactorOTPEmail(user.Email, otp)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %v", err)
	}

	return nil
}

func (s *TwoFactorService) VerifyTwoFactor(ctx context.Context, email, code string) (bool, error) {
	queries := models.New(s.db)

	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %v", err)
	}

	if !user.TwoFactorSecret.Valid {
		return false, fmt.Errorf("two-factor authentication not enabled for this user")
	}

	valid := totp.Validate(code, user.TwoFactorSecret.String)

	return valid, nil
}

func (s *TwoFactorService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}