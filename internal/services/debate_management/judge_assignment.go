package services

import (
	"context"
	"database/sql"
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
	TargetJudgeID  int32
}

func (s *JudgeAssignmentService) UpdateJudgeAssignment(ctx context.Context, op JudgeUpdateOperation) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	// Ensure "Unassigned" room exists
	unassignedRoom, err := queries.EnsureUnassignedRoomExists(ctx, sql.NullInt32{
		Int32: op.TournamentID,
		Valid: true,
	})
	if err != nil {
		return fmt.Errorf("failed to ensure unassigned room exists: %v", err)
	}

	// Get current assignment for the judge being moved (if any)
	judgeA, err := queries.GetJudgeAssignment(ctx, models.GetJudgeAssignmentParams{
		Tournamentid:  op.TournamentID,
		Judgeid:       op.JudgeID,
		Roundnumber:   op.RoundNumber,
		Iselimination: op.IsElimination,
	})
	hasJudgeA := err == nil

	// Variables to store target judge info
	var targetJudgeID int32
	var targetJudgeIsHeadJudge bool
	var targetDebateID int32
	var hasTargetJudge bool

	// Find judge in the target room (if any)
	if op.NewRoomID > 0 && op.NewRoomID != unassignedRoom {
		// Try to find a judge in the target room
		roomJudge, err := queries.GetJudgeInRoom(ctx, models.GetJudgeInRoomParams{
			Tournamentid:  op.TournamentID,
			Roomid:        op.NewRoomID,
			Roundnumber:   op.RoundNumber,
			Iselimination: op.IsElimination,
		})
		if err == nil {
			// We found a judge in the target room
			targetJudgeID = roomJudge.Judgeid
			targetJudgeIsHeadJudge = roomJudge.Isheadjudge
			targetDebateID = roomJudge.Debateid
			hasTargetJudge = true
		}
	} else if op.TargetJudgeID > 0 {
		// We have a specific target judge to swap with
		targetJudge, err := queries.GetJudgeAssignment(ctx, models.GetJudgeAssignmentParams{
			Tournamentid:  op.TournamentID,
			Judgeid:       op.TargetJudgeID,
			Roundnumber:   op.RoundNumber,
			Iselimination: op.IsElimination,
		})
		if err == nil {
			// We found the specified target judge
			targetJudgeID = targetJudge.Judgeid
			targetJudgeIsHeadJudge = targetJudge.Isheadjudge
			targetDebateID = targetJudge.Debateid
			hasTargetJudge = true
		}
	}

	// Handle the different scenarios
	if hasJudgeA && hasTargetJudge {
		// Scenario 1/2/3: Both judges are assigned - swap them completely
		err = queries.SwapJudges(ctx, models.SwapJudgesParams{
			Judgeid:       op.JudgeID,
			Judgeid_2:     targetJudgeID,
			Tournamentid:  op.TournamentID,
			Roundnumber:   op.RoundNumber,
			Iselimination: op.IsElimination,
		})
		if err != nil {
			return fmt.Errorf("failed to swap judges: %v", err)
		}

		// Transfer ballot ownership if either judge was a head judge
		if judgeA.Isheadjudge || targetJudgeIsHeadJudge {
			err = queries.TransferBallotOwnership(ctx, models.TransferBallotOwnershipParams{
				Judgeid:   op.JudgeID,
				Judgeid_2: targetJudgeID,
			})
			if err != nil {
				return fmt.Errorf("failed to transfer ballot ownership: %v", err)
			}
		}
	} else if hasJudgeA && !hasTargetJudge {
		// Scenario 4: Moving judge to "Unassigned" room
		err = queries.UnassignJudge(ctx, models.UnassignJudgeParams{
			Judgeid:       op.JudgeID,
			Tournamentid:  op.TournamentID,
			Roundnumber:   op.RoundNumber,
			Iselimination: op.IsElimination,
		})
		if err != nil {
			return fmt.Errorf("failed to unassign judge: %v", err)
		}
	} else if !hasJudgeA && hasTargetJudge {
		// Scenario 4 reversed: Assigning unassigned judge to a debate
		// Get the debate that target judge is assigned to
		_, err := queries.GetDebateByRoomAndRound(ctx, models.GetDebateByRoomAndRoundParams{
			Tournamentid:       op.TournamentID,
			Roomid:             op.NewRoomID,
			Roundnumber:        op.RoundNumber,
			Iseliminationround: op.IsElimination,
		})
		if err != nil {
			return fmt.Errorf("failed to get debate: %v", err)
		}

		// Assign judgeA to the debate that target judge was assigned to
		err = queries.AssignJudgeToDebate(ctx, models.AssignJudgeToDebateParams{
			Tournamentid:  op.TournamentID,
			Judgeid:       op.JudgeID,
			Debateid:      targetDebateID,
			Roundnumber:   op.RoundNumber,
			Iselimination: op.IsElimination,
			Isheadjudge:   targetJudgeIsHeadJudge, // Take on target judge's head judge status
		})
		if err != nil {
			return fmt.Errorf("failed to assign judge to debate: %v", err)
		}

		// If taking over head judge position, update ballot ownership
		if targetJudgeIsHeadJudge {
			err = queries.TransferBallotOwnership(ctx, models.TransferBallotOwnershipParams{
				Judgeid:   targetJudgeID,
				Judgeid_2: op.JudgeID,
			})
			if err != nil {
				return fmt.Errorf("failed to transfer ballot ownership: %v", err)
			}
		}

		// Unassign target judge
		err = queries.UnassignJudge(ctx, models.UnassignJudgeParams{
			Judgeid:       targetJudgeID,
			Tournamentid:  op.TournamentID,
			Roundnumber:   op.RoundNumber,
			Iselimination: op.IsElimination,
		})
		if err != nil {
			return fmt.Errorf("failed to unassign target judge: %v", err)
		}
	}
	// Scenario 5: Both unassigned - nothing happens

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
