package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type SubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func SubmitSubscribe(c *gin.Context) {
	var req SubscribeRequest

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

	// Send notification email to admin
	subject := "New GCV Insights Subscription"
	body := fmt.Sprintf(`
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #1e293b; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e2e8f0; border-radius: 8px;">
	<h2 style="color: #176df4; border-bottom: 2px solid #f1f5f9; padding-bottom: 10px;">New Insights Subscription</h2>
	<p>A new visitor has subscribed to receive GCV digital engineering insights:</p>
	<p><strong>Email Address:</strong> %s</p>
</body>
</html>
`, req.Email)

	err := sendEmail(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom, "contact@gcvdanta.com", "", subject, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("Failed sending subscription alert: %v", err)})
		return
	}

	// Send confirmation receipt to client
	clientSubject := "Subscribed to GCV Insights"
	clientBody := fmt.Sprintf(`
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #1e293b; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e2e8f0; border-radius: 8px;">
	<h2 style="color: #176df4; border-bottom: 2px solid #f1f5f9; padding-bottom: 10px;">Subscribed to GCV Insights</h2>
	<p>Thank you for subscribing to GCV digital engineering insights. We'll send you our latest thinking on design, technology, and platforms.</p>
	<div style="background-color: #f8fafc; border-left: 4px solid #176df4; padding: 15px; margin-top: 20px; font-size: 13px; color: #64748b;">
		<p style="margin: 0;">This is an automated receipt for your records. Please do not reply directly to this email.</p>
	</div>
</body>
</html>
`)

	_ = sendEmail(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom, req.Email, "", clientSubject, clientBody)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Subscribed successfully!"})
}
