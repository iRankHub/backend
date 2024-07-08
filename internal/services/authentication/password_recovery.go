package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type RecoveryService struct {
	db *sql.DB
}

func NewRecoveryService(db *sql.DB) *RecoveryService {
	return &RecoveryService{db: db}
}

func (s *RecoveryService) GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *RecoveryService) RequestPasswordReset(ctx context.Context, email string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %v", err)
	}

	token, err := s.GenerateResetToken()
	if err != nil {
		return err
	}

	expires := time.Now().Add(15 * time.Minute)
	err = queries.SetResetToken(ctx, models.SetResetTokenParams{
		Userid:            user.Userid,
		ResetToken:        sql.NullString{String: token, Valid: true},
		ResetTokenExpires: sql.NullTime{Time: expires, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to set reset token: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Send password reset email
	err = utils.SendPasswordResetEmail(email, token)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %v", err)
	}

	return nil
}

func (s *RecoveryService) ForcedPasswordReset(ctx context.Context, email string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %v", err)
	}

	token, err := s.GenerateResetToken()
	if err != nil {
		return err
	}

	expires := time.Now().Add(15 * time.Minute)
	err = queries.SetResetToken(ctx, models.SetResetTokenParams{
		Userid:            user.Userid,
		ResetToken:        sql.NullString{String: token, Valid: true},
		ResetTokenExpires: sql.NullTime{Time: expires, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to set reset token: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Send forced password reset email
	err = utils.SendForcedPasswordResetEmail(email, token)
	if err != nil {
		return fmt.Errorf("failed to send forced password reset email: %v", err)
	}

	return nil
}

func (s *RecoveryService) ResetPassword(ctx context.Context, token, newPassword string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	user, err := queries.GetUserByResetToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid or expired reset token")
		}
		return fmt.Errorf("failed to verify reset token: %v", err)
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	err = queries.UpdateUserPassword(ctx, models.UpdateUserPasswordParams{
		Userid:   user.Userid,
		Password: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to update user password: %v", err)
	}

	err = queries.ClearResetToken(ctx, user.Userid)
	if err != nil {
		return fmt.Errorf("failed to clear reset token: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}