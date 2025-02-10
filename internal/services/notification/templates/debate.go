package templates

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// DebateInfo contains all necessary information about a debate
type DebateInfo struct {
	DebateID       int32
	TournamentName string
	RoundNumber    int
	IsElimination  bool
	Room           string
	StartTime      time.Time
	EndTime        time.Time
	Team1          string
	Team2          string
	Motion         string
	JudgeNames     []string
	HeadJudge      string
	BallotDeadline time.Time
}

// GetRoundAssignmentTemplate generates round assignment notifications
func GetRoundAssignmentTemplate(recipientName string, debate DebateInfo) EmailComponents {
	roundType := "Preliminary"
	if debate.IsElimination {
		roundType = "Elimination"
	}

	return EmailComponents{
		Title: fmt.Sprintf("%s Round %d Assignment", roundType, debate.RoundNumber),
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>You have been assigned to the following debate:</p>
		`, GetHighlightTemplate(recipientName)),
		Metadata: map[string]string{
			"Tournament": debate.TournamentName,
			"Round":      fmt.Sprintf("%s Round %d", roundType, debate.RoundNumber),
			"Room":       debate.Room,
			"Time":       formatDateTime(debate.StartTime),
			"Duration":   formatDuration(debate.StartTime, debate.EndTime),
			"Teams":      fmt.Sprintf("%s vs %s", debate.Team1, debate.Team2),
			"Motion":     debate.Motion,
		},
		Alerts: []string{
			"Please arrive 15 minutes before the debate starts.",
		},
		Buttons: []EmailButton{
			{
				Text: "View Debate Details",
				URL:  fmt.Sprintf("%s/debates/%d", os.Getenv("FRONTEND_URL"), debate.DebateID),
			},
		},
	}
}

// GetJudgeAssignmentTemplate generates judge assignment notifications
func GetJudgeAssignmentTemplate(judgeName string, debate DebateInfo, isHeadJudge bool) EmailComponents {
	role := "Judge"
	if isHeadJudge {
		role = "Head Judge"
	}

	return EmailComponents{
		Title: fmt.Sprintf("%s Assignment: Round %d", role, debate.RoundNumber),
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>You have been assigned as %s for the following debate:</p>
		`, GetHighlightTemplate(judgeName), role),
		Metadata: map[string]string{
			"Tournament": debate.TournamentName,
			"Round":      fmt.Sprintf("%d", debate.RoundNumber),
			"Room":       debate.Room,
			"Time":       formatDateTime(debate.StartTime),
			"Teams":      fmt.Sprintf("%s vs %s", debate.Team1, debate.Team2),
			"Panel":      strings.Join(debate.JudgeNames, ", "),
			"Motion":     debate.Motion,
		},
		Alerts: getJudgeAlerts(isHeadJudge, debate.BallotDeadline),
		Buttons: []EmailButton{
			{
				Text: "Submit Ballot",
				URL:  fmt.Sprintf("%s/debates/%d/ballot", os.Getenv("FRONTEND_URL"), debate.DebateID),
			},
			{
				Text: "View Debate",
				URL:  fmt.Sprintf("%s/debates/%d", os.Getenv("FRONTEND_URL"), debate.DebateID),
			},
		},
	}
}

// GetBallotReminderTemplate generates ballot submission reminder
func GetBallotReminderTemplate(judgeName string, debate DebateInfo) EmailComponents {
	return EmailComponents{
		Title: "Ballot Submission Reminder",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>This is a reminder to submit your ballot for the following debate:</p>
		`, GetHighlightTemplate(judgeName)),
		Metadata: map[string]string{
			"Tournament": debate.TournamentName,
			"Round":      fmt.Sprintf("%d", debate.RoundNumber),
			"Teams":      fmt.Sprintf("%s vs %s", debate.Team1, debate.Team2),
			"Deadline":   formatDateTime(debate.BallotDeadline),
		},
		Alerts: []string{
			fmt.Sprintf("Ballot submission deadline: %s", formatDateTime(debate.BallotDeadline)),
			"Please submit your ballot as soon as possible to avoid delays in the tournament.",
		},
		Buttons: []EmailButton{
			{
				Text: "Submit Ballot Now",
				URL:  fmt.Sprintf("%s/debates/%d/ballot", os.Getenv("FRONTEND_URL"), debate.DebateID),
			},
		},
	}
}

// GetDebateResultTemplate generates debate result notification
func GetDebateResultTemplate(recipientName string, debate DebateInfo, results DebateResults) EmailComponents {
	return EmailComponents{
		Title: fmt.Sprintf("Debate Results: Round %d", debate.RoundNumber),
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>The results for your debate have been finalized:</p>
		`, GetHighlightTemplate(recipientName)),
		Metadata: map[string]string{
			"Tournament": debate.TournamentName,
			"Round":      fmt.Sprintf("%d", debate.RoundNumber),
			"Winner":     results.WinningTeam,
			"Scores":     fmt.Sprintf("%s: %.2f | %s: %.2f", debate.Team1, results.Team1Score, debate.Team2, results.Team2Score),
		},
		Content2: fmt.Sprintf(`
			<div style="margin-top: 20px;">
				<h3>Judge Feedback:</h3>
				<p>%s</p>
			</div>
		`, results.Feedback),
		Buttons: []EmailButton{
			{
				Text: "View Full Results",
				URL:  fmt.Sprintf("%s/debates/%d/results", os.Getenv("FRONTEND_URL"), debate.DebateID),
			},
		},
	}
}

// GetRoomChangeTemplate generates room change notification
func GetRoomChangeTemplate(recipientName string, debate DebateInfo, oldRoom string) EmailComponents {
	return EmailComponents{
		Title: "Debate Room Change",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>There has been a room change for your upcoming debate:</p>
		`, GetHighlightTemplate(recipientName)),
		Metadata: map[string]string{
			"Tournament": debate.TournamentName,
			"Round":      fmt.Sprintf("%d", debate.RoundNumber),
			"Time":       formatDateTime(debate.StartTime),
			"Old Room":   oldRoom,
			"New Room":   debate.Room,
		},
		Alerts: []string{
			"Please note this room change and arrive at the new location on time.",
		},
		Buttons: []EmailButton{
			{
				Text: "View Updated Details",
				URL:  fmt.Sprintf("%s/debates/%d", os.Getenv("FRONTEND_URL"), debate.DebateID),
			},
		},
	}
}

// DebateResults contains the results of a debate
type DebateResults struct {
	WinningTeam string
	Team1Score  float64
	Team2Score  float64
	Feedback    string
}

// Helper functions

func getJudgeAlerts(isHeadJudge bool, deadline time.Time) []string {
	alerts := []string{
		"Please arrive 15 minutes before the debate starts.",
		fmt.Sprintf("Ballot submission deadline: %s", formatDateTime(deadline)),
	}

	if isHeadJudge {
		alerts = append(alerts,
			"As Head Judge, you are responsible for:",
			"- Coordinating with other judges",
			"- Leading the oral adjudication",
			"- Ensuring timely ballot submission from all judges",
		)
	}

	return alerts
}

func formatDuration(start, end time.Time) string {
	duration := end.Sub(start)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d hours %d minutes", hours, minutes)
	}
	return fmt.Sprintf("%d minutes", minutes)
}

func formatDateTime(t time.Time) string {
	return t.Format("Monday, January 2, 2006 at 3:04 PM")
}
