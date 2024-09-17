package services

import (
	"context"
	"database/sql"
	"log"

	"github.com/robfig/cron/v3"
)

type TournamentCountsUpdateService struct {
	db   *sql.DB
	cron *cron.Cron
}

func NewTournamentCountsUpdateService(db *sql.DB) (*TournamentCountsUpdateService, error) {
	c := cron.New()
	return &TournamentCountsUpdateService{
		db:   db,
		cron: c,
	}, nil
}

func (s *TournamentCountsUpdateService) Start() {
	s.cron.AddFunc("0 0 * * *", s.UpdateTournamentCounts) // Run daily at midnight
	s.cron.Start()
}

func (s *TournamentCountsUpdateService) Stop() {
	s.cron.Stop()
}

func (s *TournamentCountsUpdateService) UpdateTournamentCounts() {
	ctx := context.Background()

	_, err := s.db.ExecContext(ctx, "SELECT update_tournament_counts()")
	if err != nil {
		log.Printf("Failed to update tournament counts: %v\n", err)
		return
	}

	log.Println("Tournament counts updated successfully")
}
