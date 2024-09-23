package utils

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/iRankHub/backend/internal/services/notification"
)

func GetEmailTemplate(content string) string {
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
				%s
			</div>
		</body>
		</html>
	`, logoURL, content)
}

func SendNotification(notificationService *notification.NotificationService, notificationType notification.NotificationType, to, subject, content string) error {
	err := notificationService.SendNotification(context.Background(), notification.Notification{
		Type:    notificationType,
		To:      to,
		Subject: subject,
		Content: content,
	})

	if err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}

	return nil
}
