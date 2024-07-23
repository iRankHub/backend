package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/models"

)

func SendTournamentInvitations(ctx context.Context, tournament models.Tournament, league models.League, format models.Tournamentformat, queries *models.Queries) error {
	schools, err := fetchRelevantSchools(ctx, queries, league)
	if err != nil {
		return fmt.Errorf("failed to fetch relevant schools: %v", err)
	}

	volunteers, err := queries.GetAllVolunteers(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch volunteers: %v", err)
	}

	schoolSubject := fmt.Sprintf("Invitation to %s Tournament", tournament.Name)
	volunteerSubject := fmt.Sprintf("Invitation to Judge at %s Tournament", tournament.Name)

	batchSize := 50
	delay := 5 * time.Second

	var errors []error

	// Send invitations to schools in batches
	for i := 0; i < len(schools); i += batchSize {
		end := i + batchSize
		if end > len(schools) {
			end = len(schools)
		}

		batch := schools[i:end]
		for _, school := range batch {
			schoolEmailContent := prepareTournamentEmailContent(school, tournament, league, format)

			// Send to contact email
			err := SendEmail(school.Contactemail, schoolSubject, schoolEmailContent)
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to send invitation to contact email %s for school %s: %v", school.Contactemail, school.Schoolname, err))
			}

			// Send to school email
			err = SendEmail(school.Schoolemail, schoolSubject, schoolEmailContent)
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to send invitation to school email %s for school %s: %v", school.Schoolemail, school.Schoolname, err))
			}
		}

		time.Sleep(delay)
	}

	// Send invitations to volunteers in batches
	for i := 0; i < len(volunteers); i += batchSize {
		end := i + batchSize
		if end > len(volunteers) {
			end = len(volunteers)
		}

		batch := volunteers[i:end]
		for _, volunteer := range batch {
			volunteerEmailContent := prepareVolunteerEmailContent(volunteer, tournament, league, format)
			body := GetEmailTemplate(volunteerEmailContent)

			// The volunteer's email is stored in the Users table
			user, err := queries.GetUserByID(context.Background(), volunteer.Userid)
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to get email for volunteer ID %d: %v", volunteer.Volunteerid, err))
				continue
			}

			err = SendEmail(user.Email, volunteerSubject, body)
			if err != nil {
				volunteerID := "unknown"
				if volunteer.Idebatevolunteerid.Valid {
					volunteerID = volunteer.Idebatevolunteerid.String
				}
				errors = append(errors, fmt.Errorf("failed to send invitation to volunteer ID %d (iDebate ID: %s): %v", volunteer.Volunteerid, volunteerID, err))
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

func prepareTournamentEmailContent(school models.School, tournament models.Tournament, league models.League, format models.Tournamentformat) string {
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

func prepareVolunteerEmailContent(volunteer models.Volunteer, tournament models.Tournament, league models.League, format models.Tournamentformat) string {
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

func fetchRelevantSchools(ctx context.Context, queries *models.Queries, league models.League) ([]models.School, error) {
	var schools []models.School
	var err error

	var leagueDetails struct {
		Districts []string `json:"districts,omitempty"`
		Countries []string `json:"countries,omitempty"`
	}

	if len(league.Details) > 0 {
		err = json.Unmarshal(league.Details, &leagueDetails)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal league details: %v", err)
		}
	} else {
		return nil, fmt.Errorf("league details are empty")
	}

	var searchTerms []string
	if league.Leaguetype == "local" {
		searchTerms = append(searchTerms, leagueDetails.Districts...)
	} else if league.Leaguetype == "international" {
		searchTerms = append(searchTerms, leagueDetails.Countries...)
	}

	if len(searchTerms) == 0 {
		return nil, fmt.Errorf("no valid search terms found in league details")
	}

	for _, searchTerm := range searchTerms {
		var schoolsBatch []models.School
		nullSearchTerm := sql.NullString{String: searchTerm, Valid: true}
		if league.Leaguetype == "local" {
			schoolsBatch, err = queries.GetSchoolsByDistrict(ctx, nullSearchTerm)
		} else {
			schoolsBatch, err = queries.GetSchoolsByCountry(ctx, nullSearchTerm)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get schools: %v", err)
		}
		schools = append(schools, schoolsBatch...)
	}

	return schools, nil
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
