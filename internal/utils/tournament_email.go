package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/smtp"
	"time"

	"github.com/spf13/viper"

	"github.com/iRankHub/backend/internal/models"

)

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}
	viper.AutomaticEnv()
}

func getTournamentEmailTemplate(title, content string) string {
	logoURL := viper.GetString("LOGO_URL")
	if logoURL == "" {
		logoURL = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSy1c8yfmVvRgCThDUvkJTmpTrV92ANV7iSRQ&s"
	}

	return fmt.Sprintf(`
		<html>
		<head>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					padding: 20px;
					background-color: #ffffff;
				}
				h1 {
					color: #333333;
				}
				.logo {
					max-width: 200px;
					height: auto;
					margin-bottom: 20px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<img src="%s" alt="iRankHub Logo" class="logo">
				<h1>%s</h1>
				%s
			</div>
		</body>
		</html>
	`, logoURL, title, content)
}

func sendTournamentEmail(to, subject, body string) error {
	from := viper.GetString("EMAIL_FROM")
	password := viper.GetString("EMAIL_PASSWORD")
	smtpHost := viper.GetString("SMTP_HOST")
	smtpPort := viper.GetString("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject = "Subject: " + subject + "\n"
	msg := []byte(subject + mime + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func SendTournamentInvitations(ctx context.Context, tournament models.Tournament, league models.League, format models.Tournamentformat, queries *models.Queries) error {
    schools, err := fetchRelevantSchools(ctx, queries, league)
    if err != nil {
        return fmt.Errorf("failed to fetch relevant schools: %v", err)
    }

    subject := fmt.Sprintf("Invitation to %s Tournament", tournament.Name)
    emailContent := prepareTournamentEmailContent(tournament, league, format)

    batchSize := 50
    for i := 0; i < len(schools); i += batchSize {
        end := i + batchSize
        if end > len(schools) {
            end = len(schools)
        }

        batch := schools[i:end]
        for _, school := range batch {
            // Send to contact email
            err := sendTournamentEmail(school.Contactemail, subject, emailContent)
            if err != nil {
                fmt.Printf("Failed to send invitation to contact email %s: %v\n", school.Contactemail, err)
            }
        }

        time.Sleep(5 * time.Second)
    }

    return nil
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

func prepareTournamentEmailContent(tournament models.Tournament, league models.League, format models.Tournamentformat) string {
    content := fmt.Sprintf(`
        <p>Dear School Representative,</p>
        <p>We are excited to invite you to participate in the upcoming tournament:</p>
        <h2>%s</h2>
        <p><strong>League:</strong> %s</p>
        <p><strong>Format:</strong> %s</p>
        <p><strong>Location:</strong> %s</p>
        <p><strong>Dates:</strong> %s to %s</p>
        <p><strong>Tournament Fee:</strong> $%.s</p>
        <p>We look forward to your participation in this exciting event!</p>
        <p>For more information or to register, please log in to your iRankHub account.</p>
        <p>Best regards,<br>The iRankHub Team</p>
    `, tournament.Name, league.Name, format.Formatname, tournament.Location,
        tournament.Startdate.Format("Jan 2, 2006"),
        tournament.Enddate.Format("Jan 2, 2006"),
        tournament.Tournamentfee)

    return getTournamentEmailTemplate("Tournament Invitation", content)
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
	body := getTournamentEmailTemplate("Tournament Created", content)
	return sendTournamentEmail(to, subject, body)
}