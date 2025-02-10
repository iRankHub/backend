package templates

import (
	"fmt"
	"os"
	"time"
)

// Brand colors
const (
	primaryColor    = "#f06000"
	secondaryColor  = "#f27522"
	accentColor     = "#f8b98e"
	backgroundColor = "#f4f4f4"
	textColor       = "#333333"
)

// GetBaseTemplate returns the base HTML template for all emails
func GetBaseTemplate(title, content string) string {
	logoURL := os.Getenv("LOGO_URL")
	if logoURL == "" {
		logoURL = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSy1c8yfmVvRgCThDUvkJTmpTrV92ANV7iSRQ&s"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        /* Reset styles */
        body, table, td, div, p {
            margin: 0;
            padding: 0;
            font-family: Arial, sans-serif;
            line-height: 1.6;
        }
        
        /* Base styles */
        body {
            background-color: %s;
            -webkit-font-smoothing: antialiased;
            height: 100%%;
            margin: 0;
            padding: 0;
            width: 100%%;
        }
        
        /* Container styles */
        .email-container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }
        
        /* Header styles */
        .header {
            background: linear-gradient(135deg, %s 0%%, %s 100%%);
            padding: 20px;
            text-align: center;
        }
        
        .logo {
            max-width: 150px;
            height: auto;
            margin-bottom: 15px;
        }
        
        .title {
            color: #ffffff;
            font-size: 24px;
            font-weight: bold;
            margin: 10px 0;
            text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
        }
        
        /* Content styles */
        .content {
            padding: 30px;
            background-color: #ffffff;
            color: %s;
        }
        
        /* Button styles */
        .button {
            display: inline-block;
            padding: 12px 24px;
            margin: 20px 0;
            background-color: %s;
            color: #ffffff;
            text-decoration: none;
            border-radius: 4px;
            font-weight: bold;
            text-align: center;
            transition: background-color 0.3s ease;
        }
        
        .button:hover {
            background-color: %s;
        }
        
        /* Footer styles */
        .footer {
            padding: 20px;
            background-color: #ffffff;
            border-top: 1px solid #eeeeee;
            text-align: center;
            font-size: 12px;
            color: #666666;
        }
        
        /* Helper classes */
        .alert {
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
            border-left: 4px solid %s;
            background-color: #fff3e0;
        }
        
        .highlight {
            color: %s;
            font-weight: bold;
        }
        
        .text-center {
            text-align: center;
        }
        
        /* Responsive styles */
        @media screen and (max-width: 600px) {
            .email-container {
                margin: 10px;
                width: auto !important;
            }
            
            .content {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <img src="%s" alt="iRankHub Logo" class="logo">
            <div class="title">%s</div>
        </div>
        <div class="content">
            %s
        </div>
        <div class="footer">
            <p>© %d iRankHub. All rights reserved.</p>
            <p>If you have any questions, please contact our support team.</p>
        </div>
    </div>
</body>
</html>
`, title, backgroundColor, primaryColor, secondaryColor, textColor, primaryColor, secondaryColor, primaryColor, primaryColor, logoURL, title, content, getCurrentYear())
}

// GetButtonTemplate returns a styled button template
func GetButtonTemplate(text, url string) string {
	return fmt.Sprintf(`<a href="%s" class="button" target="_blank">%s</a>`, url, text)
}

// GetAlertTemplate returns a styled alert box template
func GetAlertTemplate(content string) string {
	return fmt.Sprintf(`<div class="alert">%s</div>`, content)
}

// GetHighlightTemplate returns highlighted text
func GetHighlightTemplate(text string) string {
	return fmt.Sprintf(`<span class="highlight">%s</span>`, text)
}

// getCurrentYear returns the current year for the footer copyright
func getCurrentYear() int {
	return time.Now().Year()
}

// Common layout components that can be used across different email types
type EmailComponents struct {
	Title    string
	Content  string
	Buttons  []EmailButton
	Alerts   []string
	Metadata map[string]string
	Content2 string
}

type EmailButton struct {
	Text string
	URL  string
}

// BuildEmail constructs an email using the components
func BuildEmail(components EmailComponents) string {
	var content string

	// Add metadata if exists
	if len(components.Metadata) > 0 {
		content += "<div style='margin-bottom: 20px;'>"
		for key, value := range components.Metadata {
			content += fmt.Sprintf("<p><strong>%s:</strong> %s</p>", key, value)
		}
		content += "</div>"
	}

	// Add main content
	content += components.Content

	// Add alerts if any
	for _, alert := range components.Alerts {
		content += GetAlertTemplate(alert)
	}

	// Add buttons if any
	if len(components.Buttons) > 0 {
		content += "<div class='text-center'>"
		for _, button := range components.Buttons {
			content += GetButtonTemplate(button.Text, button.URL)
		}
		content += "</div>"
	}

	return GetBaseTemplate(components.Title, content)
}
