package utils

import (
	"fmt"
	"net/smtp"

	"github.com/spf13/viper"
)

func SendWelcomeEmail(to, name string) error {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read .env file: %v", err)
	}

	from := viper.GetString("EMAIL_FROM")
	password := viper.GetString("EMAIL_PASSWORD")
	smtpHost := viper.GetString("SMTP_HOST")
	smtpPort := viper.GetString("SMTP_PORT")

	subject := "Welcome to iRankHub"

	// Create the HTML email template
	body := fmt.Sprintf(`
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
				.step {
					margin-bottom: 10px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSy1c8yfmVvRgCThDUvkJTmpTrV92ANV7iSRQ&s" alt="iDebate Logo" class="logo">
				<h1>Welcome to iRankHub, %s!</h1>
				<p>Thank you for signing up.</p>
				<p>Here are the next steps:</p>
				<ol>
					<li class="step">Complete your profile</li>
					<li class="step">Join or create a team</li>
					<li class="step">Explore upcoming tournaments</li>
				</ol>
				<p>If you have any questions, feel free to reach out to us.</p>
				<p>Best regards,<br>The iRankHub Team</p>
			</div>
		</body>
		</html>
	`, name)

	// Set the email headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Create the email message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Set up authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}