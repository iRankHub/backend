package utils

import (
    "context"
    "fmt"
    "time"

    "github.com/spf13/viper"

    "github.com/iRankHub/backend/internal/models"
    "github.com/iRankHub/backend/internal/services/notification"
)

func init() {
    viper.SetConfigFile(".env")
    if err := viper.ReadInConfig(); err != nil {
        fmt.Printf("Error reading config file: %s\n", err)
    }
    viper.AutomaticEnv()
}

func SendReminderEmails(ctx context.Context, notificationService *notification.NotificationService, invitations []models.GetPendingInvitationsRow, queries *models.Queries) error {
    var errors []error
    batchSize := 50
    delay := 5 * time.Second

    for i := 0; i < len(invitations); i += batchSize {
        end := i + batchSize
        if end > len(invitations) {
            end = len(invitations)
        }

        batch := invitations[i:end]
        for _, invitation := range batch {
            if err := sendSingleReminderEmail(ctx, notificationService, invitation, queries); err != nil {
                errors = append(errors, fmt.Errorf("failed to send reminder email for invitation %d: %v", invitation.Invitationid, err))
            }
        }

        time.Sleep(delay)
    }

    if len(errors) > 0 {
        return fmt.Errorf("encountered errors while sending reminder emails: %v", errors)
    }

    return nil
}

func sendSingleReminderEmail(ctx context.Context, notificationService *notification.NotificationService, invitation models.GetPendingInvitationsRow, queries *models.Queries) error {
    timeUntilTournament := time.Until(invitation.Tournamentstartdate)
    reminderType := getShouldSendReminder(timeUntilTournament, invitation.Status, invitation.Inviteerole == "school")

if reminderType == "none" {
        return nil
    }

    recipient, recipientType, err := getRecipientInfo(invitation)
    if err != nil {
        return err
    }

    // Get user ID for in-app notification
    user, err := queries.GetUserByEmail(ctx, recipient)
    if err != nil {
        return fmt.Errorf("failed to get user by email: %v", err)
    }

    subject := fmt.Sprintf("Reminder: %s Tournament", invitation.Tournamentname)
    content := prepareReminderEmailContent(recipientType, invitation.Tournamentname, timeUntilTournament, invitation.Invitationid, reminderType)
    body := GetEmailTemplate(content)

    // Send email notification
    if err := SendNotification(notificationService, notification.EmailNotification, recipient, subject, body); err != nil {
        return err
    }

    // Send in-app notification
    inAppContent := fmt.Sprintf("Reminder for %s Tournament: %s", invitation.Tournamentname, reminderType)
    if err := SendNotification(notificationService, notification.InAppNotification, fmt.Sprintf("%d", user.Userid), subject, inAppContent); err != nil {
        return err
    }

    return nil
}

func prepareReminderEmailContent(recipientType, tournamentName string, timeUntilTournament time.Duration, invitationID int32, reminderType string) string {
    actionURL := fmt.Sprintf("%s/invitation/%d", viper.GetString("FRONTEND_URL"), invitationID)
    tournamentStart := time.Now().Add(timeUntilTournament)

    var content string
    switch reminderType {
    case "school_accept":
        deadline := tournamentStart.Add(-3 * 24 * time.Hour)
        content = fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>This is a final reminder about your invitation to the %s Tournament.</p>
            <p>The deadline to accept invitations is on %s (3 days before the competition at 11:59 PM).</p>
            <p>Please take a moment to accept or decline the invitation by clicking the link below:</p>
            <p><a href="%s">Respond to Invitation</a></p>
            <p>If you have any questions or concerns, please don't hesitate to contact us.</p>
            <p>We hope to see you at the tournament!</p>
            <p>Best regards,<br>The iRankHub Team</p>
        `, recipientType, tournamentName, deadline.Format("Monday, January 2, 2006 at 3:04 PM"), actionURL)
    case "school_revoke":
        deadline := tournamentStart.Add(-2 * 24 * time.Hour).Add(-7 * time.Hour)
        content = fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>This is a reminder about the upcoming %s Tournament, which you have accepted to participate in.</p>
            <p>The deadline to revoke acceptances or add/remove teams is on %s (2 days before the competition at 5:00 PM).</p>
            <p>If you need to make any changes, please log in to your iRankHub account or contact us immediately.</p>
            <p>We look forward to seeing you at the tournament!</p>
            <p>Best regards,<br>The iRankHub Team</p>
        `, recipientType, tournamentName, deadline.Format("Monday, January 2, 2006 at 3:04 PM"))
    case "volunteer_accept":
        deadline := tournamentStart.Add(-2 * 24 * time.Hour)
        content = fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>This is a final reminder about your invitation to judge at the %s Tournament.</p>
            <p>The deadline to accept invitations is on %s (2 days before the competition at 11:59 PM).</p>
            <p>Please take a moment to accept or decline the invitation by clicking the link below:</p>
            <p><a href="%s">Respond to Invitation</a></p>
            <p>If you have any questions or concerns, please don't hesitate to contact us.</p>
            <p>We hope you'll join us for this exciting event!</p>
            <p>Best regards,<br>The iRankHub Team</p>
        `, recipientType, tournamentName, deadline.Format("Monday, January 2, 2006 at 3:04 PM"), actionURL)
    case "volunteer_revoke":
        deadline := tournamentStart.Add(-2 * 24 * time.Hour)
        content = fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>This is a reminder about the upcoming %s Tournament, which you have accepted to judge.</p>
            <p>The deadline to revoke acceptances is on %s (2 days before the competition at 11:59 PM).</p>
            <p>If you need to make any changes, please log in to your iRankHub account or contact us immediately.</p>
            <p>We look forward to your participation in the tournament!</p>
            <p>Best regards,<br>The iRankHub Team</p>
        `, recipientType, tournamentName, deadline.Format("Monday, January 2, 2006 at 3:04 PM"))
    default:
        content = fmt.Sprintf(`
            <p>Dear %s,</p>
            <p>This is a reminder about the upcoming %s Tournament.</p>
            <p>The tournament starts on %s.</p>
            <p>If you have any questions or need to update your participation status, please log in to your iRankHub account or contact us directly.</p>
            <p>Best regards,<br>The iRankHub Team</p>
        `, recipientType, tournamentName, tournamentStart.Format("Monday, January 2, 2006 at 3:04 PM"))
    }

    return content
}

func getRecipientInfo(invitation models.GetPendingInvitationsRow) (string, string, error) {
    switch invitation.Inviteerole {
    case "school":
        return invitation.Inviteeemail.(string), "school", nil
    case "volunteer":
        return invitation.Inviteeemail.(string), "volunteer", nil
    case "student":
        return invitation.Inviteeemail.(string), "student", nil
    default:
        return "", "", fmt.Errorf("invalid invitation role: %s", invitation.Inviteerole)
    }
}

func getShouldSendReminder(timeUntilTournament time.Duration, invitationStatus string, isSchool bool) string {
    days := int(timeUntilTournament.Hours() / 24)
    hours := int(timeUntilTournament.Hours()) % 24

    if isSchool {
        if days == 3 && hours == 0 && invitationStatus == "pending" {
            return "school_accept"
        }
        if days == 2 && hours == 7 && invitationStatus == "accepted" {
            return "school_revoke"
        }
    } else {
        if days == 2 && hours == 0 && invitationStatus == "pending" {
            return "volunteer_accept"
        }
        if days == 2 && hours == 0 && invitationStatus == "accepted" {
            return "volunteer_revoke"
        }
    }

    return "none"
}

// Utility functions

func ShouldSendReminder(daysUntilTournament int) bool {
    reminderDays := []int{180, 150, 120, 90, 60, 30, 20, 15, 10, 5, 3, 0}
    for _, day := range reminderDays {
        if daysUntilTournament == day {
            return true
        }
    }
    return false
}