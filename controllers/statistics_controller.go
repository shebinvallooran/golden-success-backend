package controllers

import (
	"net/http"

	"golden-success-backend/database"
	"golden-success-backend/models"

	"github.com/gin-gonic/gin"
)

// StatisticsResponse represents the dashboard statistics
type StatisticsResponse struct {
	TotalCategories   int64 `json:"total_categories"`
	TotalProducts     int64 `json:"total_products"`
	TotalActiveProducts int64 `json:"total_active_products"`
}

// GetStatistics retrieves dashboard statistics
func GetStatistics(c *gin.Context) {
	db := database.GetDB()
	var stats StatisticsResponse

	// Count total categories
	if err := db.Model(&models.Category{}).Count(&stats.TotalCategories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count categories"})
		return
	}

	// Count total products
	if err := db.Model(&models.Product{}).Count(&stats.TotalProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}

	// Count active products
	if err := db.Model(&models.Product{}).Where("is_active = ?", true).Count(&stats.TotalActiveProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count active products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}
