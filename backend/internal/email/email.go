package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type EmailService struct {
	config    *EmailConfig
	templates map[string]*template.Template
}

func NewEmailService(config *EmailConfig) *EmailService {
	service := &EmailService{
		config:    config,
		templates: make(map[string]*template.Template),
	}
	service.loadTemplates()
	return service
}

func (s *EmailService) loadTemplates() {
	// Issue created template
	s.templates["issue_created"] = template.Must(template.New("issue_created").Parse(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%); color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f8f9fa; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; padding: 12px 24px; background: #2563eb; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🗺️ New Issue Reported</h1>
        </div>
        <div class="content">
            <h2>{{.Title}}</h2>
            <p><strong>Category:</strong> {{.Category}}</p>
            <p><strong>Location:</strong> {{.Location}}</p>
            <p>{{.Description}}</p>
            <a href="{{.IssueURL}}" class="button">View Issue</a>
        </div>
        <div class="footer">
            <p>Community Issue Mapper - Making our community better, together.</p>
        </div>
    </div>
</body>
</html>
`))

	// Issue updated template
	s.templates["issue_updated"] = template.Must(template.New("issue_updated").Parse(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #10b981 0%, #059669 100%); color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f8f9fa; padding: 30px; border-radius: 0 0 8px 8px; }
        .status { display: inline-block; padding: 6px 12px; background: #10b981; color: white; border-radius: 4px; font-weight: bold; }
        .button { display: inline-block; padding: 12px 24px; background: #10b981; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Issue Updated</h1>
        </div>
        <div class="content">
            <h2>{{.Title}}</h2>
            <p><strong>New Status:</strong> <span class="status">{{.NewStatus}}</span></p>
            {{if .UpdateMessage}}
            <p><strong>Update:</strong> {{.UpdateMessage}}</p>
            {{end}}
            <a href="{{.IssueURL}}" class="button">View Issue</a>
        </div>
        <div class="footer">
            <p>Community Issue Mapper</p>
        </div>
    </div>
</body>
</html>
`))

	// New comment template
	s.templates["new_comment"] = template.Must(template.New("new_comment").Parse(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%); color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f8f9fa; padding: 30px; border-radius: 0 0 8px 8px; }
        .comment { background: white; padding: 15px; border-left: 4px solid #8b5cf6; margin: 20px 0; }
        .button { display: inline-block; padding: 12px 24px; background: #8b5cf6; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>💬 New Comment</h1>
        </div>
        <div class="content">
            <h2>{{.IssueTitle}}</h2>
            <p><strong>{{.CommenterName}}</strong> commented:</p>
            <div class="comment">
                {{.Comment}}
            </div>
            <a href="{{.IssueURL}}" class="button">View Discussion</a>
        </div>
        <div class="footer">
            <p>Community Issue Mapper</p>
        </div>
    </div>
</body>
</html>
`))
}

type EmailData map[string]interface{}

func (s *EmailService) SendEmail(to []string, subject string, templateName string, data EmailData) error {
	tmpl, exists := s.templates[templateName]
	if !exists {
		return fmt.Errorf("template %s not found", templateName)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.send(to, subject, body.String())
}

func (s *EmailService) send(to []string, subject string, htmlBody string) error {
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	return smtp.SendMail(addr, auth, s.config.FromEmail, to, []byte(message))
}

// Convenience methods
func (s *EmailService) SendIssueCreatedEmail(to []string, issueData EmailData) error {
	return s.SendEmail(to, "New Issue Reported", "issue_created", issueData)
}

func (s *EmailService) SendIssueUpdatedEmail(to []string, issueData EmailData) error {
	return s.SendEmail(to, "Issue Updated", "issue_updated", issueData)
}

func (s *EmailService) SendNewCommentEmail(to []string, commentData EmailData) error {
	return s.SendEmail(to, "New Comment on Issue", "new_comment", commentData)
}
