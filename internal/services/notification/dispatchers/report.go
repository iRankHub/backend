package dispatchers

import (
	"context"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/services/notification/models"
)

// ReportDispatcher handles report-related notifications
type ReportDispatcher struct {
	*BaseDispatcher
	options DispatcherOptions
}

// NewReportDispatcher creates a new report dispatcher
func NewReportDispatcher(base *BaseDispatcher, options DispatcherOptions) *ReportDispatcher {
	return &ReportDispatcher{
		BaseDispatcher: base,
		options:        options,
	}
}

// Dispatch implements the Dispatcher interface
func (d *ReportDispatcher) Dispatch(ctx context.Context, n *models.Notification) error {
	// Set default expiration if not set
	if n.ExpiresAt.IsZero() {
		n.ExpiresAt = time.Now().AddDate(0, 0, d.options.ExpirationDays)
	}

	// Adjust delivery methods based on notification type
	switch n.Type {
	case models.ReportGeneration:
		// Basic report generation notifications
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		n.Priority = models.LowPriority

	case models.PerformanceReport:
		// Performance reports are more important
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		n.Priority = models.MediumPriority

	case models.AnalyticsReport:
		// Analytics reports for admin/coordinators
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		n.Priority = models.MediumPriority

	case models.AuditReport:
		// Audit reports are important for admins
		n.DeliveryMethods = []models.DeliveryMethod{
			models.EmailDelivery,
			models.InAppDelivery,
		}
		n.Priority = models.HighPriority
	}

	// Store in RabbitMQ first
	if err := d.rabbitmqSender.Publish(ctx, n); err != nil {
		return err
	}

	// Send through each delivery method
	for _, method := range n.DeliveryMethods {
		var err error
		switch method {
		case models.EmailDelivery:
			if d.options.EmailEnabled {
				err = d.emailSender.Send(ctx, n)
			}
		case models.InAppDelivery:
			if d.options.InAppEnabled {
				err = d.inAppSender.Send(ctx, n)
			}
		}

		if err != nil {
			n.UpdateDeliveryStatus(method, models.StatusFailed, err)
		}
	}

	return nil
}

// Helper methods for report-specific notifications

func (d *ReportDispatcher) SendReportGenerated(ctx context.Context, userID string, role models.UserRole, metadata models.ReportMetadata) error {
	n := &models.Notification{
		Category: models.ReportCategory,
		Type:     models.ReportGeneration,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("%s Report Ready", metadata.ReportType),
		Content:  formatReportContent(metadata),
		Actions: []models.Action{
			{
				Type:  models.ActionDownload,
				Label: "Download Report",
				URL:   metadata.DownloadURL,
			},
			{
				Type:  models.ActionView,
				Label: "View Online",
				URL:   fmt.Sprintf("/reports/%s", metadata.ReportID),
			},
		},
		Priority: models.LowPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *ReportDispatcher) SendPerformanceReport(ctx context.Context, userID string, role models.UserRole, metadata models.ReportMetadata) error {
	n := &models.Notification{
		Category: models.ReportCategory,
		Type:     models.PerformanceReport,
		UserID:   userID,
		UserRole: role,
		Title:    "Performance Report Available",
		Content:  formatPerformanceContent(metadata),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Report",
				URL:   fmt.Sprintf("/reports/%s", metadata.ReportID),
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *ReportDispatcher) SendAnalyticsReport(ctx context.Context, userID string, role models.UserRole, metadata models.ReportMetadata) error {
	n := &models.Notification{
		Category: models.ReportCategory,
		Type:     models.AnalyticsReport,
		UserID:   userID,
		UserRole: role,
		Title:    fmt.Sprintf("%s Analytics Report", metadata.Period),
		Content:  formatAnalyticsContent(metadata),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "View Analytics",
				URL:   fmt.Sprintf("/reports/%s", metadata.ReportID),
			},
			{
				Type:  models.ActionDownload,
				Label: "Download Full Report",
				URL:   metadata.DownloadURL,
			},
		},
		Priority: models.MediumPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

func (d *ReportDispatcher) SendAuditReport(ctx context.Context, userID string, role models.UserRole, metadata models.ReportMetadata) error {
	n := &models.Notification{
		Category: models.ReportCategory,
		Type:     models.AuditReport,
		UserID:   userID,
		UserRole: role,
		Title:    "Audit Report Ready",
		Content:  formatAuditContent(metadata),
		Actions: []models.Action{
			{
				Type:  models.ActionView,
				Label: "Review Audit",
				URL:   fmt.Sprintf("/reports/%s", metadata.ReportID),
			},
		},
		Priority: models.HighPriority,
	}

	if err := n.SetMetadata(metadata); err != nil {
		return err
	}

	return d.Dispatch(ctx, n)
}

// Helper functions for formatting report content

func formatReportContent(metadata models.ReportMetadata) string {
	return fmt.Sprintf(
		"Your requested %s report for %s has been generated and is ready for viewing. "+
			"Report size: %s",
		metadata.ReportType,
		metadata.Period,
		metadata.FileSize,
	)
}

func formatPerformanceContent(metadata models.ReportMetadata) string {
	content := fmt.Sprintf(
		"Your performance report for %s is now available. ",
		metadata.Period,
	)

	if len(metadata.KeyMetrics) > 0 {
		content += "\n\nKey Metrics:\n"
		for metric, value := range metadata.KeyMetrics {
			content += fmt.Sprintf("- %s: %v\n", metric, value)
		}
	}

	return content
}

func formatAnalyticsContent(metadata models.ReportMetadata) string {
	content := fmt.Sprintf(
		"The analytics report for %s has been generated. ",
		metadata.Period,
	)

	if len(metadata.Summary) > 0 {
		content += "\n\nHighlights:\n"
		for key, value := range metadata.Summary {
			content += fmt.Sprintf("- %s: %s\n", key, value)
		}
	}

	return content
}

func formatAuditContent(metadata models.ReportMetadata) string {
	return fmt.Sprintf(
		"An audit report has been generated for your review. This report covers %s. "+
			"Please review it at your earliest convenience.",
		metadata.Period,
	)
}
