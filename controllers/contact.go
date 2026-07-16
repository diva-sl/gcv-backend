package controllers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

type ContactRequest struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Company  string   `json:"company"`
	Phone    string   `json:"phone" binding:"required"`
	Services []string `json:"services"`
	Message  string   `json:"message" binding:"required"`
}

// sendEmail connects to the SMTP server and transmits an HTML email to the target recipient
func sendEmail(smtpHost, smtpPort, smtpUser, smtpPass, from, to, replyTo, subject, htmlBody string) error {
	fromHeader := fmt.Sprintf("From: GCV Digital Engineering <%s>\r\n", from)
	toHeader := fmt.Sprintf("To: %s\r\n", to)
	replyToHeader := ""
	if replyTo != "" {
		replyToHeader = fmt.Sprintf("Reply-To: %s\r\n", replyTo)
	}
	subjectHeader := fmt.Sprintf("Subject: %s\r\n", subject)
	mimeHeader := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	msg := []byte(fromHeader + toHeader + replyToHeader + subjectHeader + mimeHeader + htmlBody)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS Dial failed: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("creating SMTP client failed: %v", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP Auth failed: %v", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("sender validation failed: %v", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("recipient validation failed: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("opening data stream failed: %v", err)
	}

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("writing body data failed: %v", err)
	}

	return w.Close()
}

func SubmitContact(c *gin.Context) {
	var req ContactRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpFrom := os.Getenv("SMTP_FROM")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || smtpFrom == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "SMTP parameters are not fully configured in backend env."})
		return
	}

	var selectedServices = "None Specified"
	if len(req.Services) > 0 {
		selectedServices = ""
		for idx, service := range req.Services {
			if idx > 0 {
				selectedServices += ", "
			}
			selectedServices += service
		}
	}

	// 1. Formulate Admin Alert Email Body
	var adminSubject string
	if req.Company != "" {
		adminSubject = fmt.Sprintf("New GCV Inquiry from %s (%s)", req.Name, req.Company)
	} else {
		adminSubject = fmt.Sprintf("New GCV Inquiry from %s", req.Name)
	}

	adminBody := fmt.Sprintf(`
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #1e293b; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e2e8f0; border-radius: 8px;">
	<h2 style="color: #176df4; border-bottom: 2px solid #f1f5f9; padding-bottom: 10px;">New Consultation Request</h2>
	<p><strong>Full Name:</strong> %s</p>
	<p><strong>Email Address:</strong> %s</p>
	<p><strong>Company:</strong> %s</p>
	<p><strong>Mobile Number:</strong> %s</p>
	<p><strong>Capabilities Requested:</strong> %s</p>
	<div style="background-color: #f8fafc; border-left: 4px solid #176df4; padding: 15px; margin-top: 20px;">
		<h4 style="margin: 0 0 10px 0;">Project Scope & Details:</h4>
		<p style="margin: 0; white-space: pre-wrap;">%s</p>
	</div>
</body>
</html>
`, req.Name, req.Email, req.Company, req.Phone, selectedServices, req.Message)

	// Send Email 1: To yourself (admin alert)
	err := sendEmail(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom, smtpFrom, req.Email, adminSubject, adminBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("Failed sending admin notification: %v", err)})
		return
	}

	// 2. Formulate Client Auto-Reply Email Body
	clientSubject := "Inquiry Received - GCV Digital Engineering"
	clientBody := fmt.Sprintf(`
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #1e293b; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e2e8f0; border-radius: 8px;">
	<h2 style="color: #176df4; border-bottom: 2px solid #f1f5f9; padding-bottom: 10px;">Thank You for Contacting GCV</h2>
	<p>Dear %s,</p>
	<p>We have successfully received your inquiry regarding our digital engineering capabilities. A solutions architect is reviewing your project details and will follow up with you within 24 hours to schedule a consultation call.</p>
	<div style="background-color: #f8fafc; border-left: 4px solid #176df4; padding: 15px; margin-top: 20px; font-size: 13px; color: #64748b;">
		<p style="margin: 0;">This is an automated receipt for your records. Please do not reply directly to this email.</p>
	</div>
</body>
</html>
`, req.Name)

	// Send Email 2: To the customer (auto-reply receipt)
	_ = sendEmail(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom, req.Email, "", clientSubject, clientBody)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Inquiry transmitted successfully!"})
}

