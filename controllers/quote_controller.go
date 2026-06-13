package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"golden-success-backend/database"
	"golden-success-backend/models"
	"golden-success-backend/utils"

	"github.com/gin-gonic/gin"
)

// CreateQuote creates a new quote request/enquiry
func CreateQuote(c *gin.Context) {
	var quote models.Quote
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
		return
	}

	db := database.GetDB()
	quote.Status = "pending" // Default status

	if err := db.Create(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit enquiry"})
		return
	}

	// Trigger asynchronous email notification if settings enable it
	var settings models.NotificationSetting
	if err := db.First(&settings).Error; err == nil {
		if settings.EnableEmailNotifications {
			log.Println("Email notifications are enabled. Queueing notification email...")
			utils.SendEnquiryNotificationEmail(
				settings.SMTPHost,
				settings.SMTPPort,
				settings.SMTPUsername,
				settings.SMTPPassword,
				settings.SenderEmail,
				settings.NotificationEmail,
				settings.SMTPSecure,
				quote.Name,
				quote.Email,
				quote.Phone,
				quote.Company,
				quote.ProductName,
				quote.Message,
				quote.Quantity,
			)
		} else {
			log.Println("Email notifications are disabled in settings.")
		}
	} else {
		log.Printf("WARN: Could not fetch notification settings for email sending: %v", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Enquiry submitted successfully",
		"data":    quote,
	})
}

// GetQuotes retrieves a list of enquiries/quotes
func GetQuotes(c *gin.Context) {
	db := database.GetDB()
	var quotes []models.Quote

	// Order by most recent
	if err := db.Order("created_at desc").Find(&quotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enquiries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    quotes,
	})
}

// UpdateQuoteStatus updates the status of an enquiry/quote
func UpdateQuoteStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enquiry ID"})
		return
	}

	var req models.QuoteStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	db := database.GetDB()
	var quote models.Quote
	if err := db.First(&quote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enquiry not found"})
		return
	}

	quote.Status = req.Status
	if err := db.Save(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update enquiry status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Enquiry status updated successfully",
		"data":    quote,
	})
}

// DeleteQuote deletes an enquiry/quote by ID
func DeleteQuote(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enquiry ID"})
		return
	}

	db := database.GetDB()
	var quote models.Quote
	if err := db.First(&quote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enquiry not found"})
		return
	}

	if err := db.Delete(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete enquiry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Enquiry deleted successfully",
	})
}

// GetNotificationSettings fetches the current notification setup
func GetNotificationSettings(c *gin.Context) {
	db := database.GetDB()
	var settings models.NotificationSetting

	if err := db.First(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// UpdateNotificationSettings updates the notification setup
func UpdateNotificationSettings(c *gin.Context) {
	var req models.NotificationSetting
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings input: " + err.Error()})
		return
	}

	db := database.GetDB()
	var settings models.NotificationSetting
	if err := db.First(&settings).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settings not found"})
		return
	}

	// Update fields
	settings.NotificationEmail = req.NotificationEmail
	settings.EnableEmailNotifications = req.EnableEmailNotifications
	settings.SenderEmail = req.SenderEmail
	settings.SMTPHost = req.SMTPHost
	settings.SMTPPort = req.SMTPPort
	settings.SMTPUsername = req.SMTPUsername
	// Only update SMTP password if a new one is provided (not empty)
	if req.SMTPPassword != "" {
		settings.SMTPPassword = req.SMTPPassword
	}
	settings.SMTPSecure = req.SMTPSecure
	settings.UpdatedAt = time.Now()

	if err := db.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification settings"})
		return
	}

	// Return settings with password masked for security
	responseSettings := settings
	responseSettings.SMTPPassword = "••••••••"

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification settings updated successfully",
		"data":    responseSettings,
	})
}
