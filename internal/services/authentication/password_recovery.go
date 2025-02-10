package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/services/notification"
	notificationModels "github.com/iRankHub/backend/internal/services/notification/models"
	"github.com/iRankHub/backend/internal/utils"
)

type RecoveryService struct {
	db                  *sql.DB
	notificationService *notification.Service
}

func NewRecoveryService(db *sql.DB, ns *notification.Service) *RecoveryService {
	return &RecoveryService{
		db:                  db,
		notificationService: ns,
	}
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

	// Get client metadata
	metadata := notificationModels.AuthMetadata{
		DeviceInfo:   utils.GetDeviceInfo(ctx),
		Location:     "Password Reset Request",
		LastAttempt:  time.Now(),
		AttemptCount: int(utils.GetClientMetadata(ctx).AttemptCount),
		IPAddress:    utils.GetIPAddress(ctx),
	}

	err = s.notificationService.SendPasswordReset(
		ctx,
		email,
		s.getUserRole(user.Userrole),
		token,
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to send password reset notification: %v", err)
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

	// Get client metadata and increment attempt count for security tracking
	metadata := notificationModels.AuthMetadata{
		DeviceInfo:   utils.GetDeviceInfo(ctx),
		Location:     "Password Reset Request",
		LastAttempt:  time.Now(),
		AttemptCount: int(utils.GetClientMetadata(ctx).AttemptCount),
		IPAddress:    utils.GetIPAddress(ctx),
	}

	// Send security alert first
	err = s.notificationService.SendSecurityAlert(
		ctx,
		email,
		s.getUserRole(user.Userrole),
		fmt.Sprintf("Multiple failed login attempts detected from IP: %s. For security reasons, you must reset your password.", metadata.IPAddress),
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to send security alert: %v", err)
	}

	// Then send password reset notification
	err = s.notificationService.SendPasswordReset(
		ctx,
		email,
		s.getUserRole(user.Userrole),
		token,
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to send password reset notification: %v", err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("invalid or expired reset token")
		}
		return fmt.Errorf("failed to verify reset token: %v", err)
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	err = queries.UpdatePasswordAndClearResetCode(ctx, models.UpdatePasswordAndClearResetCodeParams{
		Userid:   user.Userid,
		Password: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Get client metadata
	metadata := notificationModels.AuthMetadata{
		DeviceInfo:   utils.GetDeviceInfo(ctx),
		Location:     "Password Reset Request",
		LastAttempt:  time.Now(),
		AttemptCount: int(utils.GetClientMetadata(ctx).AttemptCount),
		IPAddress:    utils.GetIPAddress(ctx),
	}

	err = s.notificationService.SendSecurityAlert(
		ctx,
		user.Email,
		s.getUserRole(user.Userrole),
		fmt.Sprintf("Your password has been successfully reset from IP: %s. If you did not make this change, please contact support immediately.", metadata.IPAddress),
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to send password change confirmation: %v", err)
	}

	return nil
}

func (s *RecoveryService) getUserRole(role string) notificationModels.UserRole {
	switch role {
	case "admin":
		return notificationModels.AdminRole
	case "school":
		return notificationModels.SchoolRole
	case "student":
		return notificationModels.StudentRole
	case "volunteer":
		return notificationModels.VolunteerRole
	default:
		return notificationModels.UnspecifiedRole
	}
}
