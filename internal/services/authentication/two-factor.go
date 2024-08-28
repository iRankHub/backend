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
	return totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    900, // 15 minutes in seconds
		Skew:      1,   // Allow 1 period before and after
		Digits:    6,
	})
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

	go func() {
		err := emails.SendTwoFactorOTPEmail(user.Email, otp)
		if err != nil {
			fmt.Printf("failed to send OTP email: %v\n", err)
		}
	}()

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

	valid, err := totp.ValidateCustom(code, user.TwoFactorSecret.String, time.Now(), totp.ValidateOpts{
		Period:    900, // 15 minutes in seconds
		Skew:      1,   // Allow 1 period before and after
		Digits:    6,
	})

	if err != nil {
		return false, fmt.Errorf("error validating OTP: %v", err)
	}

	return valid, nil
}

func (s *TwoFactorService) EnableTwoFactor(ctx context.Context, userID int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	secret, err := s.GenerateSecret()
	if err != nil {
		return fmt.Errorf("failed to generate secret: %v", err)
	}

	err = queries.UpdateAndEnableTwoFactor(ctx, models.UpdateAndEnableTwoFactorParams{
		Userid:           userID,
		TwoFactorSecret:  sql.NullString{String: secret, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to enable two-factor authentication: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	go func() {
		otp, err := s.GenerateOTP(secret)
		if err != nil {
			return
		}

		err = emails.SendTwoFactorOTPEmail(user.Email, otp)
		if err != nil {
			fmt.Printf("failed to send OTP email: %v\n", err)
		}
	}()

	return nil
}

func (s *TwoFactorService) DisableTwoFactor(ctx context.Context, userID int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	err = queries.UpdateAndEnableTwoFactor(ctx, models.UpdateAndEnableTwoFactorParams{
		Userid:           userID,
		TwoFactorSecret:  sql.NullString{String: "", Valid: false},
	})
	if err != nil {
		return fmt.Errorf("failed to disable two-factor authentication: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *TwoFactorService) ValidateCode(secret, code string) (bool, error) {
	return totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    900, // 15 minutes in seconds
		Skew:      1,   // Allow 1 period before and after
		Digits:    6,
	})
}