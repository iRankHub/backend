package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iRankHub/backend/internal/models"
)

type JudgeAssignmentService struct {
	db *sql.DB
}

// ValidateHeadJudgeParams matches the CheckHeadJudgeExistsParams from models
type ValidateHeadJudgeParams struct {
	Tournamentid  int32
	Roomid        int32
	Roundnumber   int32
	Iselimination bool
}

func NewJudgeAssignmentService(db *sql.DB) *JudgeAssignmentService {
	return &JudgeAssignmentService{db: db}
}

type JudgeUpdateOperation struct {
	TournamentID   int32
	JudgeID        int32
	RoundNumber    int32
	IsElimination  bool
	OldRoomID      int32
	NewRoomID      int32
	WasHeadJudge   bool
	NewIsHeadJudge bool
}

func (s *JudgeAssignmentService) UpdateJudgeAssignment(ctx context.Context, op JudgeUpdateOperation) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	// 1. Get current judge assignment details
	currentAssignment, err := queries.GetJudgeAssignment(ctx, models.GetJudgeAssignmentParams{
		Tournamentid:  op.TournamentID,
		Judgeid:       op.JudgeID,
		Roundnumber:   op.RoundNumber,
		Iselimination: op.IsElimination,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to get current assignment: %v", err)
	}

	// 2. If moving a head judge, ensure replacement in old room
	if currentAssignment.Isheadjudge {
		// Find another judge in the old room to promote to head judge
		newHeadJudgeID, err := queries.GetEligibleHeadJudge(ctx, models.GetEligibleHeadJudgeParams{
			Tournamentid:  op.TournamentID,
			Roundnumber:   op.RoundNumber,
			Roomid:        op.OldRoomID,
			Iselimination: op.IsElimination,
			Judgeid:       op.JudgeID, // This is the ExcludeJudgeID field
		})
		if err != nil {
			return fmt.Errorf("cannot move head judge - no replacement available in old room: %v", err)
		}

		// Update the new head judge
		err = queries.UpdateJudgeToHeadJudge(ctx, models.UpdateJudgeToHeadJudgeParams{
			Judgeid:       newHeadJudgeID,
			Tournamentid:  op.TournamentID,
			Roundnumber:   op.RoundNumber,
			Roomid:        op.OldRoomID,
			Iselimination: op.IsElimination,
		})
		if err != nil {
			return fmt.Errorf("failed to update new head judge: %v", err)
		}

		// Update related ballot records
		err = queries.TransferBallotOwnership(ctx, models.TransferBallotOwnershipParams{
			Judgeid:   op.JudgeID,     // This is OldJudgeID
			Judgeid_2: newHeadJudgeID, // This is NewJudgeID
			Debateid:  currentAssignment.Debateid,
		})
		if err != nil {
			return fmt.Errorf("failed to transfer ballot ownership: %v", err)
		}
	}

	// 3. Check head judge status in the new room
	if op.NewIsHeadJudge {
		// Demote current head judge in the new room if exists
		err = queries.DemoteCurrentHeadJudge(ctx, models.DemoteCurrentHeadJudgeParams{
			Tournamentid:  op.TournamentID,
			Roundnumber:   op.RoundNumber,
			Roomid:        op.NewRoomID,
			Iselimination: op.IsElimination,
		})
		if err != nil {
			return fmt.Errorf("failed to demote current head judge: %v", err)
		}
	}

	// 4. Update the judge's assignment
	err = queries.UpdateJudgeAssignment(ctx, models.UpdateJudgeAssignmentParams{
		Judgeid:            op.JudgeID,
		Tournamentid:       op.TournamentID,
		Roundnumber:        op.RoundNumber,
		Roomid:             op.NewRoomID, // Changed from NewRoomID
		Isheadjudge:        op.NewIsHeadJudge,
		Iseliminationround: op.IsElimination,
	})
	if err != nil {
		return fmt.Errorf("failed to update judge assignment: %v", err)
	}

	// 5. Validate head judge presence in both rooms
	err = s.validateHeadJudgePresence(ctx, queries, models.CheckHeadJudgeExistsParams{
		Tournamentid:  op.TournamentID,
		Roomid:        op.OldRoomID,
		Roundnumber:   op.RoundNumber,
		Iselimination: op.IsElimination,
	})
	if err != nil {
		return fmt.Errorf("old room validation failed: %v", err)
	}

	err = s.validateHeadJudgePresence(ctx, queries, models.CheckHeadJudgeExistsParams{
		Tournamentid:  op.TournamentID,
		Roomid:        op.NewRoomID,
		Roundnumber:   op.RoundNumber,
		Iselimination: op.IsElimination,
	})
	if err != nil {
		return fmt.Errorf("new room validation failed: %v", err)
	}

	return tx.Commit()
}

func (s *JudgeAssignmentService) validateHeadJudgePresence(ctx context.Context, queries *models.Queries, params models.CheckHeadJudgeExistsParams) error {
	hasHeadJudge, err := queries.CheckHeadJudgeExists(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to check head judge presence: %v", err)
	}

	if !hasHeadJudge {
		return fmt.Errorf("room %d does not have a head judge assigned", params.Roomid)
	}

	return nil
}
