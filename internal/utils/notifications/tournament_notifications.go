package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/services/notification"
)

func SendTournamentInvitations(ctx context.Context, notificationService *notification.NotificationService, tournament models.Tournament, league models.League, format models.Tournamentformat, queries *models.Queries) error {
	log.Printf("Starting to send invitations for tournament %d", tournament.Tournamentid)

	invitations, err := queries.GetPendingInvitations(ctx, tournament.Tournamentid)
	if err != nil {
		log.Printf("Error fetching pending invitations: %v", err)
		return fmt.Errorf("failed to fetch pending invitations: %v", err)
	}

	log.Printf("Found %d pending invitations for tournament %d", len(invitations), tournament.Tournamentid)

	if len(invitations) == 0 {
		log.Printf("No pending invitations found for tournament %d. This is unexpected.", tournament.Tournamentid)
		return nil
	}

	numWorkers := 5
	jobChan := make(chan models.GetPendingInvitationsRow, len(invitations))
	errChan := make(chan error, len(invitations))
	doneChan := make(chan bool)

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		go worker(ctx, notificationService, jobChan, errChan, doneChan, tournament, league, format, queries)
	}


	// Send jobs to workers
	for _, invitation := range invitations {
		jobChan <- invitation
	}
	close(jobChan)

	// Wait for all workers to finish
	for i := 0; i < numWorkers; i++ {
		<-doneChan
	}

	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		log.Printf("Encountered errors while sending tournament invitations: %v", errors)
		return fmt.Errorf("encountered errors while sending tournament invitations: %v", errors)
	}

	log.Printf("Successfully sent all invitations for tournament %d", tournament.Tournamentid)
	return nil
}

func worker(ctx context.Context, notificationService *notification.NotificationService, jobs <-chan models.GetPendingInvitationsRow, errors chan<- error, done chan<- bool, tournament models.Tournament, league models.League, format models.Tournamentformat, queries *models.Queries) {
	for invitation := range jobs {
		err := sendInvitation(ctx, notificationService, invitation, tournament, league, format, queries)
		errors <- err
	}
	done <- true
}

func sendInvitation(ctx context.Context, notificationService *notification.NotificationService, invitation models.GetPendingInvitationsRow, tournament models.Tournament, league models.League, format models.Tournamentformat, queries *models.Queries) error {
	var subject, content string
	var email string
	var userID int32

	switch invitation.Inviteerole {
	case "student":
		student, err := queries.GetStudentByIDebateID(ctx, sql.NullString{String: invitation.Inviteeid, Valid: invitation.Inviteeid != ""})
		if err != nil {
			return fmt.Errorf("failed to get student details: %v", err)
		}
		subject = fmt.Sprintf("Invitation to Participate in %s Tournament", tournament.Name)
		content = PrepareStudentInvitationContent(student, tournament, league, format)
		email = invitation.Inviteeemail.(string)
		userID = student.Userid
	case "school":
		school, err := queries.GetSchoolByIDebateID(ctx, sql.NullString{String: invitation.Inviteeid, Valid: invitation.Inviteeid != ""})
		if err != nil {
			return fmt.Errorf("failed to get school details: %v", err)
		}
		subject = fmt.Sprintf("Invitation to %s Tournament", tournament.Name)
		content = PrepareSchoolInvitationContent(school, tournament, league, format)
		email = invitation.Inviteeemail.(string)
		userID = school.Contactpersonid
	case "volunteer":
		volunteer, err := queries.GetVolunteerByIDebateID(ctx, sql.NullString{String: invitation.Inviteeid, Valid: invitation.Inviteeid != ""})
		if err != nil {
			return fmt.Errorf("failed to get volunteer details: %v", err)
		}
		subject = fmt.Sprintf("Invitation to Judge at %s Tournament", tournament.Name)
		content = PrepareVolunteerInvitationContent(volunteer, tournament, league, format)
		email = invitation.Inviteeemail.(string)
		userID = volunteer.Userid
	default:
		return fmt.Errorf("unknown invitee role: %s", invitation.Inviteerole)
	}

	log.Printf("Sending invitation email to %s for invitation ID %d", email, invitation.Invitationid)

	// Send email notification
	err := SendNotification(notificationService, notification.EmailNotification, email, subject, content)
	if err != nil {
		log.Printf("Failed to send invitation email to %s for invitation ID %d: %v", email, invitation.Invitationid, err)
		return fmt.Errorf("failed to send invitation to email %s for invitation ID %d: %v", email, invitation.Invitationid, err)
	}

	// Send in-app notification
	inAppContent := fmt.Sprintf("You've been invited to the %s Tournament", tournament.Name)
	err = SendNotification(notificationService, notification.InAppNotification, fmt.Sprintf("%d", userID), subject, inAppContent)
	if err != nil {
		log.Printf("Failed to send in-app notification to user ID %d for invitation ID %d: %v", userID, invitation.Invitationid, err)
		return fmt.Errorf("failed to send in-app notification to user ID %d for invitation ID %d: %v", userID, invitation.Invitationid, err)
	}

	// Update the reminder sent timestamp
	_, err = queries.UpdateReminderSentAt(ctx, models.UpdateReminderSentAtParams{
		Invitationid:   invitation.Invitationid,
		Remindersentat: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		log.Printf("Failed to update reminder sent timestamp for invitation ID %d: %v", invitation.Invitationid, err)
		return fmt.Errorf("failed to update reminder sent timestamp for invitation %d: %v", invitation.Invitationid, err)
	}

	log.Printf("Successfully sent invitation email, in-app notification, and updated timestamp for invitation ID %d", invitation.Invitationid)
	return nil
}

func SendTournamentCreationConfirmation(notificationService *notification.NotificationService, to, name, tournamentName string, userID int32) error {
    subject := "Tournament Created Successfully"
    content := fmt.Sprintf(`
        <p>Dear %s,</p>
        <p>We are pleased to inform you that the tournament "%s" has been successfully created in iRankHub.</p>
        <p>Invitations have been sent to eligible schools based on the league settings.</p>
        <p>You can now manage this tournament through your iRankHub dashboard.</p>
        <p>If you need to make any changes or have any questions, please don't hesitate to use the tournament management tools or contact our support team.</p>
        <p>Best regards,<br>The iRankHub Team</p>
    `, name, tournamentName)
    body := GetEmailTemplate(content)

    // Send email notification
    if err := SendNotification(notificationService, notification.EmailNotification, to, subject, body); err != nil {
        return err
    }

    // Send in-app notification
    inAppContent := fmt.Sprintf("Tournament '%s' has been successfully created", tournamentName)
    if err := SendNotification(notificationService, notification.InAppNotification, fmt.Sprintf("%d", userID), subject, inAppContent); err != nil {
        return err
    }

    return nil
}

func SendCoordinatorAssignmentEmail(notificationService *notification.NotificationService, coordinator models.User, tournament models.Tournament, league models.League, format models.Tournamentformat) error {
    subject := fmt.Sprintf("You've been assigned as coordinator for %s Tournament", tournament.Name)

    dateTimeInfo := formatDateTimeRange(tournament.Startdate, tournament.Enddate)

    content := fmt.Sprintf(`
        <p>Dear %s,</p>
        <p>We are pleased to inform you that you have been assigned as the coordinator for the following tournament:</p>
        <h2>%s</h2>
        <p><strong>League:</strong> %s</p>
        <p><strong>Format:</strong> %s</p>
        <p><strong>Location:</strong> %s</p>
        <p><strong>Date and Time:</strong> %s</p>
        <p><strong>Number of Preliminary Rounds:</strong> %d</p>
        <p><strong>Number of Elimination Rounds:</strong> %d</p>
        <p><strong>Judges per Debate (Preliminary):</strong> %d</p>
        <p><strong>Judges per Debate (Elimination):</strong> %d</p>
        <p><strong>Tournament Fee:</strong> %s</p>
        <p>As the tournament coordinator, your responsibilities include:</p>
        <ul>
            <li>Overseeing the tournament organization and ensuring smooth operations</li>
            <li>Managing participant registrations and team assignments</li>
            <li>Coordinating with judges and volunteers</li>
            <li>Handling any issues or concerns that arise during the tournament</li>
            <li>Ensuring fair play and adherence to tournament rules</li>
        </ul>
        <p>Please log in to your iRankHub account for more detailed information about the tournament and your coordinator dashboard.</p>
        <p>If you have any questions or need any assistance in your role as coordinator, please don't hesitate to contact our support team.</p>
        <p>Thank you for your commitment to making this tournament a success!</p>
        <p>Best regards,<br>The iRankHub Team</p>
    `, coordinator.Name, tournament.Name, league.Name, format.Formatname, tournament.Location,
        dateTimeInfo, tournament.Numberofpreliminaryrounds, tournament.Numberofeliminationrounds,
        tournament.Judgesperdebatepreliminary, tournament.Judgesperdebateelimination, tournament.Tournamentfee)

    body := GetEmailTemplate(content)

    // Send email notification
    if err := SendNotification(notificationService, notification.EmailNotification, coordinator.Email, subject, body); err != nil {
        return err
    }

    // Send in-app notification
    inAppContent := fmt.Sprintf("You've been assigned as coordinator for the %s Tournament", tournament.Name)
    if err := SendNotification(notificationService, notification.InAppNotification, fmt.Sprintf("%d", coordinator.Userid), subject, inAppContent); err != nil {
        return err
    }

    return nil
}

func PrepareSchoolInvitationContent(school models.School, tournament models.Tournament, league models.League, format models.Tournamentformat) string {
	var currencySymbol string
	if league.Leaguetype == "local" {
		currencySymbol = "RWF"
	} else {
		currencySymbol = "$"
	}

	dateTimeInfo := formatDateTimeRange(tournament.Startdate, tournament.Enddate)
	acceptanceDeadline := tournament.Startdate.Add(-3 * 24 * time.Hour)
	revokeDeadline := tournament.Startdate.Add(-2 * 24 * time.Hour).Add(-7 * time.Hour)

	content := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>We are excited to invite %s to participate in the upcoming tournament:</p>
		<h2>%s</h2>
		<p><strong>League:</strong> %s</p>
		<p><strong>Format:</strong> %s</p>
		<p><strong>Location:</strong> %s</p>
		<p><strong>Date and Time:</strong> %s</p>
		<p><strong>Tournament Fee:</strong> %s%s</p>
		<p>Important Deadlines:</p>
		<ul>
			<li>Deadline to accept invitation: %s</li>
			<li>Deadline to revoke acceptance or add/remove teams: %s</li>
		</ul>
		<p>We look forward to your participation in this exciting event!</p>
		<p>For more information or to register, please log in to your iRankHub account.</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, school.Schoolname, school.Schoolname, tournament.Name, league.Name, format.Formatname, tournament.Location,
		dateTimeInfo, currencySymbol, tournament.Tournamentfee,
		acceptanceDeadline.Format("Monday, January 2, 2006 at 11:59 PM"),
		revokeDeadline.Format("Monday, January 2, 2006 at 5:00 PM"))

	return GetEmailTemplate(content)
}

func PrepareVolunteerInvitationContent(volunteer models.Volunteer, tournament models.Tournament, league models.League, format models.Tournamentformat) string {
	dateTimeInfo := formatDateTimeRange(tournament.Startdate, tournament.Enddate)
	acceptanceDeadline := tournament.Startdate.Add(-2 * 24 * time.Hour)

	content := fmt.Sprintf(`
		<p>Dear %s %s,</p>
		<p>We are pleased to invite you to participate as a judge in the upcoming tournament:</p>
		<h2>%s</h2>
		<p><strong>League:</strong> %s</p>
		<p><strong>Format:</strong> %s</p>
		<p><strong>Location:</strong> %s</p>
		<p><strong>Date and Time:</strong> %s</p>
		<p>Important Deadline:</p>
		<ul>
			<li>Deadline to accept or decline invitation: %s</li>
		</ul>
		<p>Your participation as a judge is important to the success of this tournament. We value your commitment to fair play and your willingness to contribute to the debate community.</p>
		<p>Please log in to your iRankHub account to confirm your availability and see more details about the event. If you need any guidance or have questions about the judging process, we're here to help.</p>
		<p>Thank you for your dedication to supporting young debaters. Your involvement makes a real difference!</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, volunteer.Firstname, volunteer.Lastname, tournament.Name, league.Name, format.Formatname, tournament.Location,
		dateTimeInfo, acceptanceDeadline.Format("Monday, January 2, 2006 at 11:59 PM"))

	return content
}

func PrepareStudentInvitationContent(student models.Student, tournament models.Tournament, league models.League, format models.Tournamentformat) string {
	dateTimeInfo := formatDateTimeRange(tournament.Startdate, tournament.Enddate)
	acceptanceDeadline := tournament.Startdate.Add(-2 * 24 * time.Hour)

	content := fmt.Sprintf(`
		<p>Dear %s %s,</p>
		<p>We are excited to invite you to participate in the upcoming DAC tournament:</p>
		<h2>%s</h2>
		<p><strong>League:</strong> %s</p>
		<p><strong>Format:</strong> %s</p>
		<p><strong>Location:</strong> %s</p>
		<p><strong>Date and Time:</strong> %s</p>
		<p>Important Deadline:</p>
		<ul>
			<li>Deadline to accept or decline invitation: %s</li>
		</ul>
		<p>This is a great opportunity for you to showcase your debating skills in a competitive environment.</p>
		<p>Please log in to your iRankHub account to confirm your participation and see more details about the event. If you have any questions or need any assistance, please don't hesitate to reach out to our support team.</p>
		<p>We look forward to your participation in this exciting tournament!</p>
		<p>Best regards,<br>The iRankHub Team</p>
	`, student.Firstname, student.Lastname, tournament.Name, league.Name, format.Formatname, tournament.Location,
		dateTimeInfo, acceptanceDeadline.Format("Monday, January 2, 2006 at 11:59 PM"))

	return GetEmailTemplate(content)
}

func formatDateTimeRange(start, end time.Time) string {
	if start.Year() == end.Year() && start.Month() == end.Month() && start.Day() == end.Day() {
		// Same day
		return fmt.Sprintf("%s, %s from %s to %s",
			start.Weekday(),
			start.Format("January 2, 2006"),
			start.Format("15:04"),
			end.Format("15:04"))
	}
	// Different days
	return fmt.Sprintf("%s, %s to %s, %s",
		start.Weekday(),
		start.Format("January 2, 2006 at 15:04"),
		end.Weekday(),
		end.Format("January 2, 2006 at 15:04"))
}
