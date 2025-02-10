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
	"github.com/iRankHub/backend/internal/services/notification"
	notificationModels "github.com/iRankHub/backend/internal/services/notification/models"
	"github.com/iRankHub/backend/internal/utils"
)

type TwoFactorService struct {
	db                  *sql.DB
	notificationService *notification.Service
}

func NewTwoFactorService(db *sql.DB, ns *notification.Service) *TwoFactorService {
	return &TwoFactorService{
		db:                  db,
		notificationService: ns,
	}
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
		Period: 900, // 15 minutes in seconds
		Skew:   1,   // Allow 1 period before and after
		Digits: 6,
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

	// Get client metadata and increment attempt count
	clientMeta := utils.FromContext(ctx)
	clientMeta.AttemptCount++

	// Update context with new attempt count
	ctx = utils.WithAttemptCount(ctx, clientMeta.AttemptCount)

	otp, err := s.GenerateOTP(user.TwoFactorSecret.String)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %v", err)
	}

	metadata := notificationModels.AuthMetadata{
		DeviceInfo:   clientMeta.DeviceInfo,
		Location:     "2FA Login",
		LastAttempt:  time.Now(),
		AttemptCount: int(clientMeta.AttemptCount),
		IPAddress:    clientMeta.IP,
	}

	err = s.notificationService.SendTwoFactorCode(
		ctx,
		email,
		s.getUserRole(user.Userrole),
		otp,
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to send OTP: %v", err)
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

	valid, err := totp.ValidateCustom(code, user.TwoFactorSecret.String, time.Now(), totp.ValidateOpts{
		Period: 900, // 15 minutes in seconds
		Skew:   1,   // Allow 1 period before and after
		Digits: 6,
	})

	// Get client metadata for logging attempts
	clientMeta := utils.FromContext(ctx)

	if err != nil || !valid {
		// Increment failed attempts
		clientMeta.AttemptCount++
		ctx = utils.WithAttemptCount(ctx, clientMeta.AttemptCount)

		// Send security alert if too many failed attempts
		if clientMeta.AttemptCount >= 3 {
			metadata := notificationModels.AuthMetadata{
				DeviceInfo:   clientMeta.DeviceInfo,
				Location:     "2FA Verification",
				LastAttempt:  time.Now(),
				AttemptCount: int(clientMeta.AttemptCount),
				IPAddress:    clientMeta.IP,
			}

			s.notificationService.SendSecurityAlert(
				ctx,
				email,
				s.getUserRole(user.Userrole),
				fmt.Sprintf("Multiple failed 2FA attempts detected from IP: %s", clientMeta.IP),
				metadata,
			)
		}

		if err != nil {
			return false, fmt.Errorf("error validating OTP: %v", err)
		}
		return false, nil
	}

	// Reset attempt count on successful verification
	ctx = utils.WithAttemptCount(ctx, 0)

	return true, nil
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
		Userid:          userID,
		TwoFactorSecret: sql.NullString{String: secret, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to enable two-factor authentication: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Get client metadata
	clientMeta := utils.FromContext(ctx)

	metadata := notificationModels.AuthMetadata{
		DeviceInfo:   clientMeta.DeviceInfo,
		Location:     "2FA Setup",
		LastAttempt:  time.Now(),
		AttemptCount: 0,
		IPAddress:    clientMeta.IP,
	}

	// Send initial setup OTP
	otp, err := s.GenerateOTP(secret)
	if err != nil {
		return fmt.Errorf("failed to generate setup OTP: %v", err)
	}

	err = s.notificationService.SendTwoFactorCode(
		ctx,
		user.Email,
		s.getUserRole(user.Userrole),
		otp,
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to send setup OTP: %v", err)
	}

	// Send confirmation of 2FA enablement
	err = s.notificationService.SendSecurityAlert(
		ctx,
		user.Email,
		s.getUserRole(user.Userrole),
		"Two-factor authentication has been enabled for your account.",
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to send 2FA enablement confirmation: %v", err)
	}

	return nil
}

func (s *TwoFactorService) DisableTwoFactor(ctx context.Context, userID int32) error {
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

	err = queries.UpdateAndEnableTwoFactor(ctx, models.UpdateAndEnableTwoFactorParams{
		Userid:          userID,
		TwoFactorSecret: sql.NullString{String: "", Valid: false},
	})
	if err != nil {
		return fmt.Errorf("failed to disable two-factor authentication: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Get client metadata
	clientMeta := utils.FromContext(ctx)

	metadata := notificationModels.AuthMetadata{
		DeviceInfo:   clientMeta.DeviceInfo,
		Location:     "2FA Disable",
		LastAttempt:  time.Now(),
		AttemptCount: 0,
		IPAddress:    clientMeta.IP,
	}

	// Send security alert about 2FA being disabled
	err = s.notificationService.SendSecurityAlert(
		ctx,
		user.Email,
		s.getUserRole(user.Userrole),
		"Two-factor authentication has been disabled for your account. If you did not make this change, please contact support immediately.",
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to send 2FA disablement alert: %v", err)
	}

	return nil
}

// getUserRole converts string role to notificationModels.UserRole
func (s *TwoFactorService) getUserRole(role string) notificationModels.UserRole {
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
