package templates

import (
	"encoding/json"
	"fmt"
	"github.com/iRankHub/backend/internal/utils"
	"os"
	"time"
)

// GetTournamentCreationTemplate generates notification for tournament creation
func GetTournamentCreationTemplate(coordinatorName, tournamentName string, tournamentID int32) EmailComponents {
	return EmailComponents{
		Title: "Tournament Created Successfully",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>The tournament "%s" has been successfully created in iRankHub.</p>
			<p>Next steps:</p>
			<ul>
				<li>Review tournament details</li>
				<li>Set up rounds and schedules</li>
				<li>Manage team registrations</li>
				<li>Assign judges and volunteers</li>
			</ul>
		`, GetHighlightTemplate(coordinatorName), GetHighlightTemplate(tournamentName)),
		Buttons: []EmailButton{
			{
				Text: "Manage Tournament",
				URL:  fmt.Sprintf("%s/tournaments/%d", os.Getenv("FRONTEND_URL"), tournamentID),
			},
		},
	}
}

// GetCoordinatorAssignmentTemplate generates notification for tournament coordinator assignment
func GetCoordinatorAssignmentTemplate(coordinatorName string, tournament TournamentInfo) EmailComponents {
	return EmailComponents{
		Title: fmt.Sprintf("Tournament Coordinator Assignment: %s", tournament.Name),
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>You have been assigned as the coordinator for the following tournament:</p>
		`, GetHighlightTemplate(coordinatorName)),
		Metadata: map[string]string{
			"Tournament": tournament.Name,
			"League":     tournament.League,
			"Format":     tournament.Format,
			"Location":   tournament.Location,
			"Date":       formatDateRange(tournament.StartDate, tournament.EndDate),
		},
		Content2: `
			<p>Your responsibilities include:</p>
			<ul>
				<li>Managing tournament organization and operations</li>
				<li>Overseeing participant registrations</li>
				<li>Coordinating with judges and volunteers</li>
				<li>Ensuring fair play and rule compliance</li>
				<li>Handling tournament-related issues</li>
			</ul>
		`,
		Buttons: []EmailButton{
			{
				Text: "Access Coordinator Dashboard",
				URL:  fmt.Sprintf("%s/tournaments/%d/manage", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

// GetSchoolInvitationTemplate generates tournament invitation for schools
func GetSchoolInvitationTemplate(schoolName string, tournament TournamentInfo) EmailComponents {
	// Initialize schedule calculator from utils
	calculator := utils.NewScheduleCalculator(
		tournament.StartDate,
		tournament.PreliminaryRounds,
		tournament.EliminationRounds,
	)

	// Calculate fees for different team counts
	baseFee := tournament.Fee
	teamFees := make(map[int]float64)
	for i := 1; i <= 3; i++ {
		teamFees[i] = baseFee + float64(i-1)*15000
	}

	// Parse motions
	var motionsData struct {
		Preliminary []struct {
			Text        string `json:"text"`
			RoundNumber int    `json:"roundNumber"`
		} `json:"preliminary"`
		Elimination []struct {
			Text        string `json:"text"`
			RoundNumber int    `json:"roundNumber"`
		} `json:"elimination"`
	}
	json.Unmarshal(tournament.Motions, &motionsData)

	// Format motions string
	motionsContent := "<h3>Preliminary Rounds</h3><ul>"
	for i := 0; i < tournament.PreliminaryRounds-1; i++ {
		motionsContent += fmt.Sprintf("<li>Round %d: %s</li>", i+1, motionsData.Preliminary[i].Text)
	}
	motionsContent += fmt.Sprintf("<li>Round %d: Impromptu Motion (15 min prep)</li>", tournament.PreliminaryRounds)
	motionsContent += "</ul>"

	return EmailComponents{
		Title: fmt.Sprintf("Tournament Invitation: %s", tournament.Name),
		Content: fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>iDebate Rwanda officially invites you to the %s tournament on %s at %s. Each school is requested to come with a maximum of three (3) teams, each with a maximum of three debaters with their teacher.</p>

            <h3>Schedule</h3>
            %s

            <h3>Registration Fees</h3>
            <ul>
                <li>One team: %s</li>
                <li>Two teams: %s</li>
                <li>Three teams: %s</li>
                <li>Extra participant: %s 3,500</li>
            </ul>

            %s

            <h3>Payment Details</h3>
            <p>Bank: Bank of Kigali (BK)<br>
            Account: 00044-00492526-06<br>
            Mobile Money: *182*8*1*022788#</p>
        `,
			GetHighlightTemplate(schoolName),
			tournament.Name,
			formatDate(tournament.StartDate),
			tournament.Location,
			calculator.FormatSchedule(),
			formatCurrency(teamFees[1], tournament.Currency),
			formatCurrency(teamFees[2], tournament.Currency),
			formatCurrency(teamFees[3], tournament.Currency),
			tournament.Currency,
			motionsContent,
		),
		Metadata: map[string]string{
			"Tournament": tournament.Name,
			"League":     tournament.League,
			"Format":     tournament.Format,
			"Location":   tournament.Location,
			"Date":       formatDateRange(tournament.StartDate, tournament.EndDate),
			"Fee":        formatCurrency(tournament.Fee, tournament.Currency),
		},
		Alerts: []string{
			fmt.Sprintf("Registration deadline: %s", formatDeadline(tournament.StartDate, -3)),
			fmt.Sprintf("Team changes deadline: %s", formatDeadline(tournament.StartDate, -2)),
			fmt.Sprintf("Check-in time: %s - %s",
				calculator.GetCheckInTime().Format("3:04 PM"),
				calculator.GetFirstDebateTime().Format("3:04 PM")),
			"Any team that arrives after 8:30 AM shall have to forfeit the first round.",
		},
		Buttons: []EmailButton{
			{
				Text: "Register Now",
				URL:  fmt.Sprintf("%s/tournaments/%d/register", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
			{
				Text: "View Details",
				URL:  fmt.Sprintf("%s/tournaments/%d", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

// GetVolunteerInvitationTemplate generates tournament invitation for judges/volunteers
func GetVolunteerInvitationTemplate(volunteerName string, tournament TournamentInfo) EmailComponents {
	return EmailComponents{
		Title: fmt.Sprintf("Judge Invitation: %s", tournament.Name),
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>You are invited to judge at the upcoming tournament:</p>
		`, GetHighlightTemplate(volunteerName)),
		Metadata: map[string]string{
			"Tournament": tournament.Name,
			"League":     tournament.League,
			"Format":     tournament.Format,
			"Location":   tournament.Location,
			"Date":       formatDateRange(tournament.StartDate, tournament.EndDate),
		},
		Alerts: []string{
			fmt.Sprintf("Response deadline: %s", formatDeadline(tournament.StartDate, -2)),
		},
		Buttons: []EmailButton{
			{
				Text: "Accept Invitation",
				URL:  fmt.Sprintf("%s/tournaments/%d/accept", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
			{
				Text: "View Details",
				URL:  fmt.Sprintf("%s/tournaments/%d", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

// GetStudentInvitationTemplate generates tournament invitation for DAC students
func GetStudentInvitationTemplate(studentName string, tournament TournamentInfo) EmailComponents {
	return EmailComponents{
		Title: fmt.Sprintf("DAC Tournament Invitation: %s", tournament.Name),
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>You are invited to participate in the upcoming DAC tournament:</p>
		`, GetHighlightTemplate(studentName)),
		Metadata: map[string]string{
			"Tournament": tournament.Name,
			"Format":     tournament.Format,
			"Location":   tournament.Location,
			"Date":       formatDateRange(tournament.StartDate, tournament.EndDate),
		},
		Alerts: []string{
			fmt.Sprintf("Response deadline: %s", formatDeadline(tournament.StartDate, -2)),
		},
		Buttons: []EmailButton{
			{
				Text: "Accept Invitation",
				URL:  fmt.Sprintf("%s/tournaments/%d/accept", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
			{
				Text: "View Details",
				URL:  fmt.Sprintf("%s/tournaments/%d", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

// GetRegistrationConfirmationTemplate generates registration confirmation
func GetRegistrationConfirmationTemplate(recipientName string, tournament TournamentInfo, registrationDetails RegistrationInfo) EmailComponents {
	return EmailComponents{
		Title: "Tournament Registration Confirmed",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>Your registration for %s has been confirmed.</p>
		`, GetHighlightTemplate(recipientName), GetHighlightTemplate(tournament.Name)),
		Metadata: map[string]string{
			"Tournament":   tournament.Name,
			"Date":         formatDateRange(tournament.StartDate, tournament.EndDate),
			"Teams":        fmt.Sprintf("%d", registrationDetails.TeamCount),
			"Total Amount": formatCurrency(registrationDetails.TotalAmount, tournament.Currency),
			"Status":       registrationDetails.PaymentStatus,
		},
		Alerts: getPaymentAlerts(registrationDetails),
		Buttons: []EmailButton{
			{
				Text: "View Registration",
				URL:  fmt.Sprintf("%s/tournaments/%d/registration", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

func GetTournamentReminderTemplate(recipientName string, tournament TournamentInfo, daysToGo int) EmailComponents {
	calculator := utils.NewScheduleCalculator(
		tournament.StartDate,
		tournament.PreliminaryRounds,
		tournament.EliminationRounds,
	)

	var title string
	if daysToGo == 7 {
		title = "1 WEEK TO GO!"
	} else if daysToGo == 1 {
		title = "TOMORROW"
	} else {
		title = fmt.Sprintf("%d DAYS TO GO!", daysToGo)
	}

	return EmailComponents{
		Title: title,
		Content: fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>The %s is getting closer and our anticipation is building! We can't wait for an exciting day of debate and civil discourse.</p>

            <h3>Key Information</h3>
            <p><strong>Date:</strong> %s</p>
            <p><strong>Venue:</strong> %s</p>
            <p><strong>Check-in Time:</strong> %s</p>

            <h3>Schedule Reminder</h3>
            %s

            <p>Don't forget to:</p>
            <ul>
                <li>Bring your team registration documents</li>
                <li>Have proof of payment ready</li>
                <li>Arrive on time - teams arriving after 8:30 AM will forfeit their first round</li>
            </ul>
        `,
			GetHighlightTemplate(recipientName),
			tournament.Name,
			formatDate(tournament.StartDate),
			tournament.Location,
			calculator.GetCheckInTime().Format("3:04 PM"),
			calculator.FormatSchedule(),
		),
		Alerts: []string{
			fmt.Sprintf("Check-in begins at %s", calculator.GetCheckInTime().Format("3:04 PM")),
			"Please ensure all payments are completed before the tournament.",
		},
		Buttons: []EmailButton{
			{
				Text: "View Tournament Details",
				URL:  fmt.Sprintf("%s/tournaments/%d", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
			{
				Text: "View Registration Status",
				URL:  fmt.Sprintf("%s/tournaments/%d/registration", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

// TeamResult contains team performance details
type TeamResult struct {
	TeamName    string
	Rank        int
	TotalPoints float64
	IsWinner    bool
}

// GetTournamentCongratulationsTemplate generates congratulatory email for tournament winners
func GetTournamentCongratulationsTemplate(recipient string, tournament TournamentInfo, result TeamResult) EmailComponents {
	isDAC := tournament.League == "DAC"

	var content string
	if isDAC {
		// For DAC tournaments, congratulate team directly
		content = fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>CONGRATULATIONS on your incredible victory at the %s!</p>
            <p>Your team demonstrated exceptional skill, dedication, and sportsmanship throughout the tournament. 
            Your performance has set a high standard for debate excellence.</p>
            <p>Final Results:</p>
            <ul>
                <li>Overall Rank: %d</li>
                <li>Total Points: %.2f</li>
            </ul>
            <p>We look forward to seeing your continued growth and success in future tournaments!</p>
        `,
			GetHighlightTemplate(recipient),
			tournament.Name,
			result.Rank,
			result.TotalPoints,
		)
	} else {
		// For regular tournaments, congratulate both school and team
		content = fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>CONGRATULATIONS to %s on their outstanding achievement at the %s!</p>
            <p>Your school and team have demonstrated exceptional dedication to debate excellence. 
            This victory reflects the hard work of your debaters, the support of your school, 
            and the commitment to fostering strong debate skills.</p>
            <p>Final Results:</p>
            <ul>
                <li>Team: %s</li>
                <li>Overall Rank: %d</li>
                <li>Total Points: %.2f</li>
            </ul>
            <p>We hope this achievement inspires continued excellence in debate at your school!</p>
        `,
			GetHighlightTemplate(recipient),
			result.TeamName,
			tournament.Name,
			result.TeamName,
			result.Rank,
			result.TotalPoints,
		)
	}

	return EmailComponents{
		Title:   fmt.Sprintf("Congratulations on Your %s Victory!", tournament.Name),
		Content: content,
		Metadata: map[string]string{
			"Tournament": tournament.Name,
			"Team":       result.TeamName,
			"Rank":       fmt.Sprintf("%d", result.Rank),
			"Points":     fmt.Sprintf("%.2f", result.TotalPoints),
		},
		Buttons: []EmailButton{
			{
				Text: "View Tournament Results",
				URL:  fmt.Sprintf("%s/tournaments/%d/results", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
			{
				Text: "View Team Performance",
				URL:  fmt.Sprintf("%s/tournaments/%d/teams/%s", os.Getenv("FRONTEND_URL"), tournament.ID, result.TeamName),
			},
		},
	}
}

// GetUpdateForConfirmedTeamsTemplate generates update email for confirmed teams
func GetUpdateForConfirmedTeamsTemplate(recipientName string, tournament TournamentInfo) EmailComponents {
	calculator := utils.NewScheduleCalculator(
		tournament.StartDate,
		tournament.PreliminaryRounds,
		tournament.EliminationRounds,
	)

	return EmailComponents{
		Title: fmt.Sprintf("Important Updates: %s This Saturday", tournament.Name),
		Content: fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>We are thrilled your school is participating in the %s this Saturday! Here are the key timing highlights and logistics for your attending teams:</p>

            <h3>Tournament Day Schedule</h3>
            %s

            <p>Additional Notes:</p>
            <ul>
                <li>Breakfast and refreshments will be provided</li>
                <li>Lunch will be served during the break</li>
                <li>Please bring any necessary debate materials</li>
            </ul>
        `,
			GetHighlightTemplate(recipientName),
			tournament.Name,
			calculator.FormatSchedule(),
		),
		Alerts: []string{
			"Don't forget to bring your team registration documents",
			"Have your proof of payment ready at registration",
			fmt.Sprintf("Registration closes at %s sharp", calculator.GetFirstDebateTime().Format("3:04 PM")),
		},
		Buttons: []EmailButton{
			{
				Text: "View Full Schedule",
				URL:  fmt.Sprintf("%s/tournaments/%d/schedule", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
			{
				Text: "Tournament Hub",
				URL:  fmt.Sprintf("%s/tournaments/%d", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

// GetScheduleUpdateTemplate generates schedule update notification
func GetScheduleUpdateTemplate(recipientName string, tournament TournamentInfo, update ScheduleUpdate) EmailComponents {
	return EmailComponents{
		Title: "Tournament Schedule Update",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>There has been an update to the schedule for %s:</p>
		`, GetHighlightTemplate(recipientName), GetHighlightTemplate(tournament.Name)),
		Metadata: map[string]string{
			"Type":     update.Type,
			"Previous": formatDateTime(update.PreviousTime),
			"New":      formatDateTime(update.NewTime),
			"Affected": update.AffectedItems,
		},
		Alerts: []string{
			"Please adjust your schedule accordingly.",
		},
		Buttons: []EmailButton{
			{
				Text: "View Updated Schedule",
				URL:  fmt.Sprintf("%s/tournaments/%d/schedule", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

// GetPaymentReminderTemplate generates payment reminder notification
func GetPaymentReminderTemplate(schoolName string, tournament TournamentInfo, payment PaymentInfo) EmailComponents {
	return EmailComponents{
		Title: "Tournament Payment Reminder",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>This is a reminder about the pending payment for %s:</p>
		`, GetHighlightTemplate(schoolName), GetHighlightTemplate(tournament.Name)),
		Metadata: map[string]string{
			"Amount Due":     formatCurrency(payment.AmountDue, tournament.Currency),
			"Due Date":       formatDate(payment.DueDate),
			"Payment Method": payment.PaymentMethod,
		},
		Alerts: []string{
			fmt.Sprintf("Payment deadline: %s", formatDate(payment.DueDate)),
			"Teams cannot participate without completed payment.",
		},
		Buttons: []EmailButton{
			{
				Text: "Complete Payment",
				URL:  fmt.Sprintf("%s/tournaments/%d/payment", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

type TournamentInfo struct {
	ID                    int32
	Name                  string
	League                string
	Format                string
	Motions               json.RawMessage
	Location              string
	StartDate             time.Time
	EndDate               time.Time
	Fee                   float64
	Currency              string
	PreliminaryRounds     int
	EliminationRounds     int
	JudgesPerDebatePrelim int
	JudgesPerDebateElim   int
}

type RegistrationInfo struct {
	TeamCount     int
	TotalAmount   float64
	PaymentStatus string
}

type ScheduleUpdate struct {
	Type          string
	PreviousTime  time.Time
	NewTime       time.Time
	AffectedItems string
}

type PaymentInfo struct {
	AmountDue     float64
	DueDate       time.Time
	PaymentMethod string
}

func formatDateRange(start, end time.Time) string {
	if start.Year() == end.Year() && start.Month() == end.Month() && start.Day() == end.Day() {
		return fmt.Sprintf("%s (%s - %s)",
			start.Format("Monday, January 2, 2006"),
			start.Format("3:04 PM"),
			end.Format("3:04 PM"))
	}
	return fmt.Sprintf("%s - %s",
		start.Format("Monday, January 2, 2006 3:04 PM"),
		end.Format("Monday, January 2, 2006 3:04 PM"))
}

func formatDeadline(date time.Time, daysBefore int) string {
	deadline := date.AddDate(0, 0, daysBefore)
	return deadline.Format("Monday, January 2, 2006 3:04 PM")
}

func formatDate(t time.Time) string {
	return t.Format("Monday, January 2, 2006")
}

func formatCurrency(amount float64, currency string) string {
	if currency == "RWF" {
		return fmt.Sprintf("RWF %.0f", amount)
	}
	return fmt.Sprintf("$%.2f", amount)
}

func getPaymentAlerts(registration RegistrationInfo) []string {
	var alerts []string
	switch registration.PaymentStatus {
	case "pending":
		alerts = append(alerts, "Payment is required to complete registration.")
	case "partial":
		alerts = append(alerts, "Remaining balance must be paid before the tournament.")
	case "paid":
		alerts = append(alerts, "Payment has been received. Thank you!")
	}
	return alerts
}
