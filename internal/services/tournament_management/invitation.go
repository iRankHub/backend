package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
	emails "github.com/iRankHub/backend/internal/utils/emails"
)

type InvitationService struct {
	db *sql.DB
}

func NewInvitationService(db *sql.DB) *InvitationService {
	return &InvitationService{db: db}
}

func (s *InvitationService) AcceptInvitation(ctx context.Context, req *tournament_management.AcceptInvitationRequest) (*tournament_management.AcceptInvitationResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}
	userID := sql.NullInt32{Int32: int32(userIDFloat), Valid: true}
	log.Printf("User ID from token: %d", userID.Int32)

	err = queries.UpdateInvitationStatusWithUserCheck(ctx, models.UpdateInvitationStatusWithUserCheckParams{
		Invitationid: req.GetInvitationId(),
		Status:       "accepted",
		Userid:       userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to accept invitation: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.AcceptInvitationResponse{
		Success: true,
		Message: "Invitation accepted successfully",
	}, nil
}

func (s *InvitationService) DeclineInvitation(ctx context.Context, req *tournament_management.DeclineInvitationRequest) (*tournament_management.DeclineInvitationResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}
	userID := sql.NullInt32{Int32: int32(userIDFloat), Valid: true}

	err = queries.UpdateInvitationStatusWithUserCheck(ctx, models.UpdateInvitationStatusWithUserCheckParams{
		Invitationid: req.GetInvitationId(),
		Status:       "declined",
		Userid:       userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to decline invitation: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.DeclineInvitationResponse{
		Success: true,
		Message: "Invitation declined successfully",
	}, nil
}

func (s *InvitationService) BulkAcceptInvitations(ctx context.Context, req *tournament_management.BulkAcceptInvitationsRequest) (*tournament_management.BulkAcceptInvitationsResponse, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	for _, invitationID := range req.GetInvitationIds() {
		err = queries.UpdateInvitationStatus(ctx, models.UpdateInvitationStatusParams{
			Invitationid: invitationID,
			Status:       "accepted",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to accept invitation %d: %v", invitationID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.BulkAcceptInvitationsResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully accepted %d invitations", len(req.GetInvitationIds())),
	}, nil
}

func (s *InvitationService) BulkDeclineInvitations(ctx context.Context, req *tournament_management.BulkDeclineInvitationsRequest) (*tournament_management.BulkDeclineInvitationsResponse, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(s.db).WithTx(tx)

	for _, invitationID := range req.GetInvitationIds() {
		err = queries.UpdateInvitationStatus(ctx, models.UpdateInvitationStatusParams{
			Invitationid: invitationID,
			Status:       "declined",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to decline invitation %d: %v", invitationID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &tournament_management.BulkDeclineInvitationsResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully declined %d invitations", len(req.GetInvitationIds())),
	}, nil
}

func (s *InvitationService) ResendInvitation(ctx context.Context, req *tournament_management.ResendInvitationRequest) (*tournament_management.ResendInvitationResponse, error) {
    _, err := s.validateAdminRole(req.GetToken())
    if err != nil {
        return nil, err
    }

    queries := models.New(s.db)

    invitation, err := queries.GetInvitationByID(ctx, req.GetInvitationId())
    if err != nil {
        return nil, fmt.Errorf("failed to get invitation: %v", err)
    }

    // Update the reminder sent timestamp
    _, err = queries.UpdateReminderSentAt(ctx, models.UpdateReminderSentAtParams{
        Invitationid:   invitation.Invitationid,
        Remindersentat: sql.NullTime{Time: time.Now(), Valid: true},
    })
    if err != nil {
        return nil, fmt.Errorf("failed to update reminder sent timestamp: %v", err)
    }

    // Send email asynchronously
    go s.sendInvitationEmailAsync(invitation.Tournamentid)

    return &tournament_management.ResendInvitationResponse{
        Success: true,
        Message: "Invitation resend process started",
    }, nil
}

func (s *InvitationService) BulkResendInvitations(ctx context.Context, req *tournament_management.BulkResendInvitationsRequest) (*tournament_management.BulkResendInvitationsResponse, error) {
    _, err := s.validateAdminRole(req.GetToken())
    if err != nil {
        return nil, err
    }

    queries := models.New(s.db)

    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %v", err)
    }
    defer tx.Rollback()

    for _, invitationID := range req.GetInvitationIds() {
        invitation, err := queries.GetInvitationByID(ctx, invitationID)
        if err != nil {
            return nil, fmt.Errorf("failed to get invitation %d: %v", invitationID, err)
        }

        _, err = queries.UpdateReminderSentAt(ctx, models.UpdateReminderSentAtParams{
            Invitationid:   invitationID,
            Remindersentat: sql.NullTime{Time: time.Now(), Valid: true},
        })
        if err != nil {
            return nil, fmt.Errorf("failed to update reminder sent timestamp for invitation %d: %v", invitationID, err)
        }

        // Send email asynchronously
        go s.sendInvitationEmailAsync(invitation.Tournamentid)
    }

    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return &tournament_management.BulkResendInvitationsResponse{
        Success: true,
        Message: fmt.Sprintf("Resend process started for %d invitations", len(req.GetInvitationIds())),
    }, nil
}

func (s *InvitationService) sendInvitationEmailAsync(invitationID int32) {
    ctx := context.Background()
    queries := models.New(s.db)

    invitations, err := queries.GetPendingInvitations(ctx, invitationID)
    if err != nil {
        log.Printf("Failed to get invitations for ID %d: %v", invitationID, err)
        return
    }

    if len(invitations) == 0 {
        log.Printf("No pending invitations found for ID %d", invitationID)
        return
    }

    // We'll process only the first invitation, as we're looking for a specific one
    invitation := invitations[0]

    tournamentRow, err := queries.GetTournamentByID(ctx, invitation.Tournamentid)
    if err != nil {
        log.Printf("Failed to get tournament for invitation %d: %v", invitationID, err)
        return
    }

    // Convert GetTournamentByIDRow to Tournament
    tournament := convertToTournament(tournamentRow)

    league, err := queries.GetLeagueByID(ctx, tournament.Leagueid.Int32)
    if err != nil {
        log.Printf("Failed to get league for tournament %d: %v", tournament.Tournamentid, err)
        return
    }

    format, err := queries.GetTournamentFormatByID(ctx, tournament.Formatid)
    if err != nil {
        log.Printf("Failed to get tournament format for tournament %d: %v", tournament.Tournamentid, err)
        return
    }

    var subject, content string
    var email string

    if invitation.Studentid.Valid {
        subject = fmt.Sprintf("Reminder: Invitation to Participate in %s Tournament", tournament.Name)
        student := models.Student{
            Studentid:  invitation.Studentid.Int32,
            Email:      invitation.Studentemail,
            Firstname:  invitation.Studentfirstname.String,
            Lastname:   invitation.Studentlastname.String,
        }
        content = emails.PrepareStudentInvitationContent(student, tournament, league, format)
        email = invitation.Studentemail.String
    } else if invitation.Schoolid.Valid {
        subject = fmt.Sprintf("Reminder: Invitation to %s Tournament", tournament.Name)
        school := models.School{
            Schoolid:     invitation.Schoolid.Int32,
            Schoolname:   invitation.Schoolname.String,
            Contactemail: invitation.Contactemail.String,
        }
        content = emails.PrepareSchoolInvitationContent(school, tournament, league, format)
        email = invitation.Contactemail.String
    } else if invitation.Volunteerid.Valid {
        subject = fmt.Sprintf("Reminder: Invitation to Judge at %s Tournament", tournament.Name)
        volunteer := models.Volunteer{
            Volunteerid: invitation.Volunteerid.Int32,
            Firstname:   invitation.Volunteerfirstname.String,
            Lastname:    invitation.Volunteerlastname.String,
        }
        content = emails.PrepareVolunteerInvitationContent(volunteer, tournament, league, format)
        email = invitation.Volunteeremail.String
    }

    err = emails.SendEmail(email, subject, content)
    if err != nil {
        log.Printf("Failed to send invitation email to %s for invitation ID %d: %v", email, invitationID, err)
    } else {
        log.Printf("Successfully sent invitation email to %s for invitation ID %d", email, invitationID)
    }
}

func (s *InvitationService) GetInvitationStatus(ctx context.Context, req *tournament_management.GetInvitationStatusRequest) (*tournament_management.GetInvitationStatusResponse, error) {
	log.Printf("GetInvitationStatus called with invitation ID: %d", req.GetInvitationId())

	if err := s.validateAuthentication(req.GetToken()); err != nil {
		log.Printf("Authentication failed: %v", err)
		return nil, err
	}
	log.Printf("Token validated successfully")

	queries := models.New(s.db)
	status, err := queries.GetInvitationStatus(ctx, req.GetInvitationId())
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation status: %v", err)
	}

	log.Printf("Invitation status retrieved successfully for invitation ID: %d", req.GetInvitationId())
	return &tournament_management.GetInvitationStatusResponse{
		Status: status,
	}, nil
}

func (s *InvitationService) GetInvitationsByUser(ctx context.Context, req *tournament_management.GetInvitationsByUserRequest) (*tournament_management.GetInvitationsByUserResponse, error) {
    claims, err := utils.ValidateToken(req.GetToken())
    if err != nil {
        return nil, fmt.Errorf("authentication failed: %v", err)
    }

    userIDFloat, ok := claims["user_id"].(float64)
    if !ok {
        return nil, fmt.Errorf("invalid user ID in token")
    }
    userID := sql.NullInt32{
        Int32: int32(userIDFloat),
        Valid: true,
    }

    queries := models.New(s.db)
    invitations, err := queries.GetInvitationsByUserID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get invitations: %v", err)
    }

    var protoInvitations []*tournament_management.Invitation
    for _, inv := range invitations {
        protoInvitations = append(protoInvitations, &tournament_management.Invitation{
            InvitationId: inv.Invitationid,
            TournamentId: inv.Tournamentid,
            Status:       inv.Status,
        })
    }

    return &tournament_management.GetInvitationsByUserResponse{
        Invitations: protoInvitations,
    }, nil
}

func (s *InvitationService) GetAllInvitations(ctx context.Context, req *tournament_management.GetAllInvitationsRequest) (*tournament_management.GetAllInvitationsResponse, error) {
	_, err := s.validateAdminRole(req.GetToken())
	if err != nil {
		return nil, err
	}

	queries := models.New(s.db)

	invitations, err := queries.GetAllInvitations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all invitations: %v", err)
	}

	var protoInvitations []*tournament_management.Invitation
	for _, inv := range invitations {
		protoInvitations = append(protoInvitations, &tournament_management.Invitation{
			InvitationId: inv.Invitationid,
			TournamentId: inv.Tournamentid,
			Status:       inv.Status,
		})
	}

	return &tournament_management.GetAllInvitationsResponse{
		Invitations: protoInvitations,
	}, nil
}

// Helper function to convert GetTournamentByIDRow to Tournament
func convertToTournament(row models.GetTournamentByIDRow) models.Tournament {
	return models.Tournament{
		Tournamentid:               row.Tournamentid,
		Name:                       row.Name,
		Startdate:                  row.Startdate,
		Enddate:                    row.Enddate,
		Location:                   row.Location,
		Formatid:                   row.Formatid,
		Leagueid:                   row.Leagueid,
		Coordinatorid:              row.Coordinatorid,
		Numberofpreliminaryrounds:  row.Numberofpreliminaryrounds,
		Numberofeliminationrounds:  row.Numberofeliminationrounds,
		Judgesperdebatepreliminary: row.Judgesperdebatepreliminary,
		Judgesperdebateelimination: row.Judgesperdebateelimination,
		Tournamentfee:              row.Tournamentfee,
	}
}

func (s *InvitationService) validateAdminRole(token string) (map[string]interface{}, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can perform this action")
	}

	return claims, nil
}

func (s *InvitationService) validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}