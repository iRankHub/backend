package utils

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/models"

)

func SendTournamentInvitations(ctx context.Context, tournament models.Tournament, league models.League, format models.Tournamentformat, queries *models.Queries) error {
	invitations, err := queries.GetPendingInvitations(ctx, tournament.Tournamentid)
	if err != nil {
		return fmt.Errorf("failed to fetch pending invitations: %v", err)
	}

	batchSize := 50
	delay := 5 * time.Second

	var errors []error

	// Process invitations in batches
	for i := 0; i < len(invitations); i += batchSize {
		end := i + batchSize
		if end > len(invitations) {
			end = len(invitations)
		}

		batch := invitations[i:end]
		for _, invitation := range batch {
			var subject, content string
			var email string

			if invitation.Studentid.Valid {
				// This is a student invitation (for DAC league)
				subject = fmt.Sprintf("Invitation to Participate in %s Tournament", tournament.Name)
				student, err := queries.GetStudentByID(ctx, invitation.Studentid.Int32)
				if err != nil {
					errors = append(errors, fmt.Errorf("failed to get student details for invitation %d: %v", invitation.Invitationid, err))
					continue
				}
				content = PrepareStudentInvitationContent(student, tournament, league, format)
				email = student.Email.String
			} else if invitation.Schoolid.Valid {
				// This is a school invitation
				subject = fmt.Sprintf("Invitation to %s Tournament", tournament.Name)
				school, err := queries.GetSchoolByID(ctx, invitation.Schoolid.Int32)
				if err != nil {
					errors = append(errors, fmt.Errorf("failed to get school details for invitation %d: %v", invitation.Invitationid, err))
					continue
				}
				content = PrepareSchoolInvitationContent(school, tournament, league, format)
				email = school.Contactemail
			} else if invitation.Volunteerid.Valid {
				// This is a volunteer invitation
				subject = fmt.Sprintf("Invitation to Judge at %s Tournament", tournament.Name)
				volunteer, err := queries.GetVolunteerByID(ctx, invitation.Volunteerid.Int32)
				if err != nil {
					errors = append(errors, fmt.Errorf("failed to get volunteer details for invitation %d: %v", invitation.Invitationid, err))
					continue
				}
				user, err := queries.GetUserByID(ctx, volunteer.Userid)
				if err != nil {
					errors = append(errors, fmt.Errorf("failed to get user details for volunteer %d: %v", volunteer.Volunteerid, err))
					continue
				}
				content = PrepareVolunteerInvitationContent(volunteer, tournament, league, format)
				email = user.Email
			}

			err := SendEmail(email, subject, content)
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to send invitation to email %s for invitation ID %d: %v", email, invitation.Invitationid, err))
			}

			// Update the reminder sent timestamp
			_, err = queries.UpdateReminderSentAt(ctx, models.UpdateReminderSentAtParams{
				Invitationid:   invitation.Invitationid,
				Remindersentat: sql.NullTime{Time: time.Now(), Valid: true},
			})
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to update reminder sent timestamp for invitation %d: %v", invitation.Invitationid, err))
			}
		}

		time.Sleep(delay)
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered errors while sending tournament invitations: %v", errors)
	}

	return nil
}

func SendTournamentCreationConfirmation(to, name, tournamentName string) error {
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
	return SendEmail(to, subject, body)
}

func SendCoordinatorAssignmentEmail(coordinator models.User, tournament models.Tournament, league models.League, format models.Tournamentformat) error {
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
	return SendEmail(coordinator.Email, subject, body)
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
