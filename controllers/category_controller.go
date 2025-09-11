package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golden-success-backend/database"
	"golden-success-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// handleImageUpload processes image upload and returns the image URL
func handleImageUpload(c *gin.Context, fieldName string) (string, error) {
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		// No file uploaded, return empty string (not an error)
		return "", nil
	}
	defer file.Close()



	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return "", fmt.Errorf("invalid file type. Only JPEG, PNG, GIF, and WebP are allowed")
	}

	// Validate file size (5MB max)
	const maxSize = 5 * 1024 * 1024 // 5MB
	if header.Size > maxSize {
		return "", fmt.Errorf("file size too large. Maximum size is 5MB")
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads/categories"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory")
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		// Try to get extension from content type
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}

	filename := fmt.Sprintf("category_%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file")
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file")
	}

	// Return the file URL
	fileURL := fmt.Sprintf("/uploads/categories/%s", filename)
	return fileURL, nil
}

// parseFormData extracts form data from multipart request
func parseFormData(c *gin.Context) (*models.CategoryCreateRequest, error) {
	req := &models.CategoryCreateRequest{}

	// Parse form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		return nil, fmt.Errorf("failed to parse multipart form")
	}

	// Extract text fields
	req.NameEn = c.PostForm("name_en")
	req.NameAr = c.PostForm("name_ar")
	req.DescriptionEn = c.PostForm("description_en")
	req.DescriptionAr = c.PostForm("description_ar")
	req.HomeDescriptionEn = c.PostForm("home_description_en")
	req.HomeDescriptionAr = c.PostForm("home_description_ar")
	req.Point1En = c.PostForm("point1_en")
	req.Point1Ar = c.PostForm("point1_ar")
	req.Point2En = c.PostForm("point2_en")
	req.Point2Ar = c.PostForm("point2_ar")
	req.Point3En = c.PostForm("point3_en")
	req.Point3Ar = c.PostForm("point3_ar")
	req.ImageURL = c.PostForm("image_url")

	// Parse priority
	if priorityStr := c.PostForm("priority"); priorityStr != "" {
		if priority, err := strconv.Atoi(priorityStr); err == nil {
			req.Priority = &priority
		}
	}

	// Parse is_active
	if isActiveStr := c.PostForm("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		req.IsActive = &isActive
	}

	return req, nil
}

// parseUpdateFormData extracts form data for update request
func parseUpdateFormData(c *gin.Context) (*models.CategoryUpdateRequest, error) {
	req := &models.CategoryUpdateRequest{}

	// Parse form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		return nil, fmt.Errorf("failed to parse multipart form")
	}

	// Extract text fields (only if provided)
	if nameEn := c.PostForm("name_en"); nameEn != "" {
		req.NameEn = &nameEn
	}
	if nameAr := c.PostForm("name_ar"); nameAr != "" {
		req.NameAr = &nameAr
	}
	if descEn := c.PostForm("description_en"); descEn != "" {
		req.DescriptionEn = &descEn
	}
	if descAr := c.PostForm("description_ar"); descAr != "" {
		req.DescriptionAr = &descAr
	}
	if homeDescEn := c.PostForm("home_description_en"); homeDescEn != "" {
		req.HomeDescriptionEn = &homeDescEn
	}
	if homeDescAr := c.PostForm("home_description_ar"); homeDescAr != "" {
		req.HomeDescriptionAr = &homeDescAr
	}
	if point1En := c.PostForm("point1_en"); point1En != "" {
		req.Point1En = &point1En
	}
	if point1Ar := c.PostForm("point1_ar"); point1Ar != "" {
		req.Point1Ar = &point1Ar
	}
	if point2En := c.PostForm("point2_en"); point2En != "" {
		req.Point2En = &point2En
	}
	if point2Ar := c.PostForm("point2_ar"); point2Ar != "" {
		req.Point2Ar = &point2Ar
	}
	if point3En := c.PostForm("point3_en"); point3En != "" {
		req.Point3En = &point3En
	}
	if point3Ar := c.PostForm("point3_ar"); point3Ar != "" {
		req.Point3Ar = &point3Ar
	}
	if imageURL := c.PostForm("image_url"); imageURL != "" {
		req.ImageURL = &imageURL
	}

	// Parse priority
	if priorityStr := c.PostForm("priority"); priorityStr != "" {
		if priority, err := strconv.Atoi(priorityStr); err == nil {
			req.Priority = &priority
		}
	}

	// Parse is_active
	if isActiveStr := c.PostForm("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		req.IsActive = &isActive
	}

	return req, nil
}

// GetCategories retrieves all categories ordered by priority
func GetCategories(c *gin.Context) {
	var categories []models.Category
	db := database.GetDB()

	// Query parameters for filtering
	isActive := c.Query("is_active")
	includeInactive := c.Query("include_inactive")

	query := db.Model(&models.Category{})

	// Apply filters
	if isActive == "true" || includeInactive != "true" {
		query = query.Where("is_active = ?", true)
	} else if isActive == "false" {
		query = query.Where("is_active = ?", false)
	}

	// Order by priority (ascending) then by English name
	if err := query.Order("priority ASC, name_en ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Set virtual fields for backward compatibility
	for i := range categories {
		categories[i].Name = categories[i].NameEn
		categories[i].Description = categories[i].DescriptionEn
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GetCategoriesWithProducts retrieves categories with their products
func GetCategoriesWithProducts(c *gin.Context) {
	var categories []models.Category
	db := database.GetDB()

	// Query parameters
	isActive := c.Query("is_active")
	includeProducts := c.Query("include_products") == "true"

	query := db.Model(&models.Category{})

	// Apply filters
	if isActive == "true" || isActive == "" {
		query = query.Where("is_active = ?", true)
	} else if isActive == "false" {
		query = query.Where("is_active = ?", false)
	}

	// Include products if requested
	if includeProducts {
		query = query.Preload("Products", "is_active = ?", true)
	}

	if err := query.Order("priority ASC, name_en ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Transform to response format
	var response []models.CategoryWithProductsResponse
	for _, category := range categories {
		// Count products in this category
		var productCount int64
		database.GetDB().Model(&models.Product{}).Where("category_id = ? AND is_active = ?", category.ID, true).Count(&productCount)

		categoryResponse := models.CategoryWithProductsResponse{
			ID:            category.ID,
			NameEn:        category.NameEn,
			NameAr:        category.NameAr,
			DescriptionEn: category.DescriptionEn,
			DescriptionAr: category.DescriptionAr,

			// Home screen description fields
			HomeDescriptionEn: category.HomeDescriptionEn,
			HomeDescriptionAr: category.HomeDescriptionAr,

			// Three key points fields
			Point1En: category.Point1En,
			Point1Ar: category.Point1Ar,
			Point2En: category.Point2En,
			Point2Ar: category.Point2Ar,
			Point3En: category.Point3En,
			Point3Ar: category.Point3Ar,

			// Image field
			ImageURL: category.ImageURL,

			Priority:      category.Priority,
			IsActive:      category.IsActive,
			CreatedAt:     category.CreatedAt,
			UpdatedAt:     category.UpdatedAt,
			ProductCount:  int(productCount),
			// Set virtual fields for backward compatibility
			Name:        category.NameEn,
			Description: category.DescriptionEn,
		}

		if includeProducts {
			categoryResponse.Products = category.Products
		}

		response = append(response, categoryResponse)
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetCategory retrieves a single category by ID
func GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	query := database.GetDB()
	
	// Check if we should include products
	if c.Query("include_products") == "true" {
		query = query.Preload("Products", "is_active = ?", true)
	}

	if err := query.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
		return
	}

	// Set virtual fields for backward compatibility
	category.Name = category.NameEn
	category.Description = category.DescriptionEn

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// CreateCategory creates a new category
func CreateCategory(c *gin.Context) {
	// Check content type to determine how to parse the request
	contentType := c.GetHeader("Content-Type")
	var req *models.CategoryCreateRequest
	var err error

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data (with potential image upload)
		req, err = parseFormData(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Handle JSON request (backward compatibility)
		var jsonReq models.CategoryCreateRequest
		if err := c.ShouldBindJSON(&jsonReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req = &jsonReq
	}

	// Check if English category name already exists
	var existingCategory models.Category
	if err := database.GetDB().Where("name_en = ?", req.NameEn).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category with this English name already exists"})
		return
	}

	// Handle image upload if present
	var imageURL string
	if strings.Contains(contentType, "multipart/form-data") {
		imageURL, err = handleImageUpload(c, "image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Image upload failed: %s", err.Error())})
			return
		}
	}

	category := models.Category{
		NameEn:        req.NameEn,
		NameAr:        req.NameAr,
		DescriptionEn: req.DescriptionEn,
		DescriptionAr: req.DescriptionAr,

		// Home screen description fields
		HomeDescriptionEn: req.HomeDescriptionEn,
		HomeDescriptionAr: req.HomeDescriptionAr,

		// Three key points fields
		Point1En: req.Point1En,
		Point1Ar: req.Point1Ar,
		Point2En: req.Point2En,
		Point2Ar: req.Point2Ar,
		Point3En: req.Point3En,
		Point3Ar: req.Point3Ar,

		// Image field - use uploaded image URL or provided URL
		ImageURL: imageURL,

		Priority:      0, // Default priority
		IsActive:      true,
	}

	// Override with provided values if present
	if req.ImageURL != "" && imageURL == "" {
		category.ImageURL = req.ImageURL
	}

	if req.Priority != nil {
		category.Priority = *req.Priority
	}

	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := database.GetDB().Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

// UpdateCategory updates an existing category
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := database.GetDB().First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
		return
	}

	// Check content type to determine how to parse the request
	contentType := c.GetHeader("Content-Type")
	var req *models.CategoryUpdateRequest
	var err error

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data (with potential image upload)
		req, err = parseUpdateFormData(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Handle JSON request (backward compatibility)
		var jsonReq models.CategoryUpdateRequest
		if err := c.ShouldBindJSON(&jsonReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req = &jsonReq
	}

	// Handle image upload if present
	var imageURL string
	if strings.Contains(contentType, "multipart/form-data") {
		imageURL, err = handleImageUpload(c, "image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Image upload failed: %s", err.Error())})
			return
		}
		// If a new image was uploaded, set it
		if imageURL != "" {
			req.ImageURL = &imageURL
		}
	}

	// Check if English name already exists (if being updated)
	if req.NameEn != nil && *req.NameEn != category.NameEn {
		var existingCategory models.Category
		if err := database.GetDB().Where("name_en = ? AND id != ?", *req.NameEn, category.ID).First(&existingCategory).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Category with this English name already exists"})
			return
		}
	}

	// Update fields if provided
	if req.NameEn != nil {
		category.NameEn = *req.NameEn
	}
	if req.NameAr != nil {
		category.NameAr = *req.NameAr
	}
	if req.DescriptionEn != nil {
		category.DescriptionEn = *req.DescriptionEn
	}
	if req.DescriptionAr != nil {
		category.DescriptionAr = *req.DescriptionAr
	}

	// Update home screen description fields
	if req.HomeDescriptionEn != nil {
		category.HomeDescriptionEn = *req.HomeDescriptionEn
	}
	if req.HomeDescriptionAr != nil {
		category.HomeDescriptionAr = *req.HomeDescriptionAr
	}

	// Update three key points fields
	if req.Point1En != nil {
		category.Point1En = *req.Point1En
	}
	if req.Point1Ar != nil {
		category.Point1Ar = *req.Point1Ar
	}
	if req.Point2En != nil {
		category.Point2En = *req.Point2En
	}
	if req.Point2Ar != nil {
		category.Point2Ar = *req.Point2Ar
	}
	if req.Point3En != nil {
		category.Point3En = *req.Point3En
	}
	if req.Point3Ar != nil {
		category.Point3Ar = *req.Point3Ar
	}

	// Update image field
	if req.ImageURL != nil {
		category.ImageURL = *req.ImageURL
	}

	if req.Priority != nil {
		category.Priority = *req.Priority
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := database.GetDB().Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory deletes a category (soft delete)
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := database.GetDB().First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
		return
	}

	// Check if category has products
	var productCount int64
	database.GetDB().Model(&models.Product{}).Where("category_id = ?", category.ID).Count(&productCount)
	
	if productCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category that has products assigned to it"})
		return
	}

	if err := database.GetDB().Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// UpdateCategoryPriorities updates multiple category priorities at once
func UpdateCategoryPriorities(c *gin.Context) {
	var req struct {
		Categories []struct {
			ID       uint `json:"id" validate:"required"`
			Priority int  `json:"priority" validate:"min=0"`
		} `json:"categories" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update priorities in a transaction
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, categoryUpdate := range req.Categories {
		if err := tx.Model(&models.Category{}).Where("id = ?", categoryUpdate.ID).Update("priority", categoryUpdate.Priority).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category priorities"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit priority updates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category priorities updated successfully"})
}

// GetCategoriesForHome retrieves categories for home page with limited fields
func GetCategoriesForHome(c *gin.Context) {
	var categories []models.Category
	db := database.GetDB()

	// Only get active categories ordered by priority
	if err := db.Where("is_active = ?", true).Order("priority ASC, name_en ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Transform to home page response format with only required fields
	var response []gin.H
	for _, category := range categories {
		response = append(response, gin.H{
			"id":                   category.ID,
			"name_en":              category.NameEn,
			"name_ar":              category.NameAr,
			"home_description_en":  category.HomeDescriptionEn,
			"home_description_ar":  category.HomeDescriptionAr,
			"image_url":            category.ImageURL,
			"priority":             category.Priority,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetCategoriesList retrieves a simple list of categories without pagination
func GetCategoriesList(c *gin.Context) {
	var categories []models.Category
	db := database.GetDB()

	// Query parameters for filtering
	isActive := c.Query("is_active")
	includeInactive := c.Query("include_inactive")

	query := db.Model(&models.Category{})

	// Apply filters
	if isActive == "true" || includeInactive != "true" {
		query = query.Where("is_active = ?", true)
	} else if isActive == "false" {
		query = query.Where("is_active = ?", false)
	}

	// Order by priority (ascending) then by English name
	if err := query.Order("priority ASC, name_en ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Transform to simple list response format with only required fields
	var response []gin.H
	for _, category := range categories {
		response = append(response, gin.H{
			"id":                   category.ID,
			"name_en":              category.NameEn,
			"name_ar":              category.NameAr,
			"home_description_en":  category.HomeDescriptionEn,
			"home_description_ar":  category.HomeDescriptionAr,
			"image_url":            category.ImageURL,
			"priority":             category.Priority,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
