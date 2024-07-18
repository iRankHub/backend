package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type LoginService struct {
    db               *sql.DB
    twoFactorService *TwoFactorService
    recoveryService  *RecoveryService
}

func NewLoginService(db *sql.DB, twoFactorService *TwoFactorService, recoveryService *RecoveryService) *LoginService {
    return &LoginService{
        db:               db,
        twoFactorService: twoFactorService,
        recoveryService:  recoveryService,
    }
}

func (s *LoginService) Login(ctx context.Context, emailOrId, password string) (*models.User, error) {
    start := time.Now()
    defer func() {
        log.Printf("Total login time: %v", time.Since(start))
    }()

    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %v", err)
    }
    defer tx.Rollback()

    queries := models.New(tx)

    getUserStart := time.Now()
    userRow, err := queries.GetUserByEmailOrIDebateIDAndUpdateLoginAttempt(ctx, emailOrId)
    log.Printf("GetUserByEmailOrIDebateIDAndUpdateLoginAttempt time: %v", time.Since(getUserStart))
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("invalid email/ID or password")
        }
        return nil, fmt.Errorf("failed to retrieve user: %v", err)
    }
    // Convert GetUserByEmailOrIDebateIDRow to User
    user := &models.User{
        Userid:               userRow.Userid,
        Webauthnuserid:       userRow.Webauthnuserid,
        Name:                 userRow.Name,
        Email:                userRow.Email,
        Password:             userRow.Password,
        Userrole:             userRow.Userrole,
        Status:               userRow.Status,
        Verificationstatus:   userRow.Verificationstatus,
        Deactivatedat:        userRow.Deactivatedat,
        TwoFactorSecret:      userRow.TwoFactorSecret,
        TwoFactorEnabled:     userRow.TwoFactorEnabled,
        FailedLoginAttempts:  userRow.FailedLoginAttempts,
        LastLoginAttempt:     userRow.LastLoginAttempt,
        LastLogout:           userRow.LastLogout,
        ResetToken:           userRow.ResetToken,
        ResetTokenExpires:    userRow.ResetTokenExpires,
        CreatedAt:            userRow.CreatedAt,
        UpdatedAt:            userRow.UpdatedAt,
        DeletedAt:            userRow.DeletedAt,
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    compareStart := time.Now()
    err = utils.ComparePasswords(user.Password, password)
    log.Printf("ComparePasswords time: %v", time.Since(compareStart))
    if err != nil {
        handleErr := s.HandleFailedLoginAttempt(ctx, user)
        if handleErr != nil {
            return nil, handleErr
        }
        return nil, fmt.Errorf("invalid email/ID or password")
    }

    if user.TwoFactorEnabled.Valid && user.TwoFactorEnabled.Bool {
        return nil, fmt.Errorf("two factor authentication required")
    }

    err = s.HandleSuccessfulLogin(ctx, user.Userid)
    if err != nil {
        return nil, fmt.Errorf("failed to handle successful login: %v", err)
    }

    return user, nil
}

func (s *LoginService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %v", err)
    }
    defer tx.Rollback()

    queries := models.New(tx)

    user, err := queries.GetUserByEmail(ctx, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to retrieve user: %v", err)
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return &user, nil
}

func (s *LoginService) HandleFailedLoginAttempt(ctx context.Context, user *models.User) error {
    start := time.Now()
    defer func() {
        log.Printf("HandleFailedLoginAttempt time: %v", time.Since(start))
    }()

    queries := models.New(s.db)

    updatedUser, err := queries.IncrementAndGetFailedLoginAttempts(ctx, user.Userid)
    if err != nil {
        return fmt.Errorf("failed to update and get login attempts: %v", err)
    }

    if updatedUser.FailedLoginAttempts.Int32 >= 4 {
        if updatedUser.TwoFactorEnabled.Valid && updatedUser.TwoFactorEnabled.Bool {
            return fmt.Errorf("two factor authentication required")
        } else {
            go func() {
                err := s.recoveryService.ForcedPasswordReset(context.Background(), updatedUser.Email)
                if err != nil {
                    log.Printf("failed to initiate forced password reset: %v\n", err)
                }
            }()
            return fmt.Errorf("password reset required")
        }
    }

    return nil
}

func (s *LoginService) HandleSuccessfulLogin(ctx context.Context, userID int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	err = queries.ResetFailedLoginAttempts(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to reset login attempts: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}