package templates

import (
	"fmt"
	"os"
	"time"
)

// ReportInfo contains information about generated reports
type ReportInfo struct {
	ReportID    int32
	ReportType  string
	GeneratedAt time.Time
	Period      string
	Size        string
	DownloadURL string
	ExpiresAt   time.Time
	Summary     map[string]string
	Insights    []string
}

// GetReportGeneratedTemplate generates notification for completed report generation
func GetReportGeneratedTemplate(recipientName string, report ReportInfo) EmailComponents {
	return EmailComponents{
		Title: fmt.Sprintf("%s Report Ready", report.ReportType),
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>Your requested report has been generated and is now available:</p>
		`, GetHighlightTemplate(recipientName)),
		Metadata: map[string]string{
			"Report Type": report.ReportType,
			"Period":      report.Period,
			"Generated":   formatDateTime(report.GeneratedAt),
			"Size":        report.Size,
			"Expires":     formatDateTime(report.ExpiresAt),
		},
		Content2: `
			<div style="margin-top: 20px;">
				<h3>Key Insights:</h3>
				<ul>
		` + formatInsights(report.Insights) + `
				</ul>
			</div>
		`,
		Alerts: []string{
			fmt.Sprintf("This report will be available for download until %s", formatDateTime(report.ExpiresAt)),
		},
		Buttons: []EmailButton{
			{
				Text: "Download Report",
				URL:  report.DownloadURL,
			},
			{
				Text: "View Online",
				URL:  fmt.Sprintf("%s/reports/%d", os.Getenv("FRONTEND_URL"), report.ReportID),
			},
		},
	}
}

// GetTournamentReportTemplate generates tournament performance report notification
func GetTournamentReportTemplate(recipientName string, tournament TournamentInfo, report ReportInfo) EmailComponents {
	return EmailComponents{
		Title: fmt.Sprintf("Tournament Report: %s", tournament.Name),
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>The tournament report for %s is now available:</p>
		`, GetHighlightTemplate(recipientName), GetHighlightTemplate(tournament.Name)),
		Metadata: map[string]string{
			"Tournament": tournament.Name,
			"Date":       formatDateRange(tournament.StartDate, tournament.EndDate),
			"Generated":  formatDateTime(report.GeneratedAt),
		},
		Content2: `
			<div style="margin-top: 20px;">
				<h3>Summary:</h3>
		` + formatSummary(report.Summary) + `
			</div>
		`,
		Buttons: []EmailButton{
			{
				Text: "View Full Report",
				URL:  fmt.Sprintf("%s/tournaments/%d/report", os.Getenv("FRONTEND_URL"), tournament.ID),
			},
		},
	}
}

// GetPerformanceReportTemplate generates individual/team performance report notification
func GetPerformanceReportTemplate(recipientName string, report ReportInfo) EmailComponents {
	return EmailComponents{
		Title: "Performance Report Available",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>Your performance report for %s is now available:</p>
		`, GetHighlightTemplate(recipientName), report.Period),
		Metadata: map[string]string{
			"Period":    report.Period,
			"Generated": formatDateTime(report.GeneratedAt),
		},
		Content2: `
			<div style="margin-top: 20px;">
				<h3>Key Statistics:</h3>
		` + formatSummary(report.Summary) + `
				<h3>Performance Insights:</h3>
				<ul>
		` + formatInsights(report.Insights) + `
				</ul>
			</div>
		`,
		Buttons: []EmailButton{
			{
				Text: "View Detailed Report",
				URL:  fmt.Sprintf("%s/reports/%d", os.Getenv("FRONTEND_URL"), report.ReportID),
			},
		},
	}
}

// GetAnalyticsReportTemplate generates analytics report notification
func GetAnalyticsReportTemplate(recipientName string, report ReportInfo) EmailComponents {
	return EmailComponents{
		Title: fmt.Sprintf("%s Analytics Report", report.Period),
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>The analytics report for %s has been generated:</p>
		`, GetHighlightTemplate(recipientName), report.Period),
		Metadata: map[string]string{
			"Report Type": report.ReportType,
			"Period":      report.Period,
			"Generated":   formatDateTime(report.GeneratedAt),
		},
		Content2: `
			<div style="margin-top: 20px;">
				<h3>Key Metrics:</h3>
		` + formatSummary(report.Summary) + `
				<h3>Notable Trends:</h3>
				<ul>
		` + formatInsights(report.Insights) + `
				</ul>
			</div>
		`,
		Buttons: []EmailButton{
			{
				Text: "View Analytics Dashboard",
				URL:  fmt.Sprintf("%s/analytics/dashboard", os.Getenv("FRONTEND_URL")),
			},
			{
				Text: "Download Full Report",
				URL:  report.DownloadURL,
			},
		},
	}
}

// GetAuditReportTemplate generates audit report notification
func GetAuditReportTemplate(recipientName string, report ReportInfo) EmailComponents {
	return EmailComponents{
		Title: "Audit Report Ready",
		Content: fmt.Sprintf(`
			<p>Dear %s,</p>
			<p>An audit report has been generated for your review:</p>
		`, GetHighlightTemplate(recipientName)),
		Metadata: map[string]string{
			"Audit Type": report.ReportType,
			"Period":     report.Period,
			"Generated":  formatDateTime(report.GeneratedAt),
		},
		Content2: `
			<div style="margin-top: 20px;">
				<h3>Audit Summary:</h3>
		` + formatSummary(report.Summary) + `
				<h3>Key Findings:</h3>
				<ul>
		` + formatInsights(report.Insights) + `
				</ul>
			</div>
		`,
		Alerts: []string{
			"Please review this report at your earliest convenience.",
			"Some findings may require immediate attention.",
		},
		Buttons: []EmailButton{
			{
				Text: "Review Audit Report",
				URL:  fmt.Sprintf("%s/reports/%d", os.Getenv("FRONTEND_URL"), report.ReportID),
			},
		},
	}
}

// Helper functions

func formatSummary(summary map[string]string) string {
	var content string
	if len(summary) == 0 {
		return "<p>No summary data available.</p>"
	}

	content = "<div class='summary'>"
	for key, value := range summary {
		content += fmt.Sprintf(`
			<div style="margin-bottom: 10px;">
				<strong>%s:</strong> %s
			</div>
		`, key, value)
	}
	content += "</div>"
	return content
}

func formatInsights(insights []string) string {
	if len(insights) == 0 {
		return "<li>No insights available.</li>"
	}

	var content string
	for _, insight := range insights {
		content += fmt.Sprintf("<li>%s</li>", insight)
	}
	return content
}
