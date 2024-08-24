package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/models"
)

type LogoutService struct {
	db *sql.DB
}

func NewLogoutService(db *sql.DB) *LogoutService {
	return &LogoutService{
		db: db,
	}
}

func (s *LogoutService) Logout(ctx context.Context, userID int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Update last logout time
	err = queries.UpdateLastLogout(ctx, models.UpdateLastLogoutParams{
		Userid:     userID,
		LastLogout: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update last logout: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
