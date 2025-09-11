package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golden-success-backend/database"
	"golden-success-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JSONResponse sends a JSON response without HTML escaping
func JSONResponse(c *gin.Context, code int, obj interface{}) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // This prevents HTML escaping
	if err := encoder.Encode(obj); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Data(code, "application/json; charset=utf-8", buf.Bytes())
}

// GetProducts retrieves all products with optional filtering
func GetProducts(c *gin.Context) {
	var products []models.Product
	db := database.GetDB()

	// Query parameters for filtering
	categoryID := c.Query("category_id")
	isActive := c.Query("is_active")
	search := c.Query("search")

	query := db.Model(&models.Product{}).Preload("CategoryInfo")

	// Apply filters
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if isActive != "" {
		if isActive == "true" {
			query = query.Where("is_active = ?", true)
		} else if isActive == "false" {
			query = query.Where("is_active = ?", false)
		}
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// Populate virtual category field for backward compatibility
	for i := range products {
		if products[i].CategoryInfo != nil {
			products[i].Category = products[i].CategoryInfo.NameEn
		}
	}

	JSONResponse(c, http.StatusOK, gin.H{
		"data": products,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ProductListResponse represents the simplified product response for the product list API
type ProductListResponse struct {
	ID            uint   `json:"id"`
	NameEn        string `json:"name_en"`
	NameAr        string `json:"name_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	ImageURL      string `json:"image_url"`
	CategoryEn    string `json:"category_en"`
	CategoryAr    string `json:"category_ar"`
}

// GetProductList retrieves all products without pagination with only specific fields
func GetProductList(c *gin.Context) {
	var products []models.Product
	db := database.GetDB()

	// Get all active products with category information
	query := db.Model(&models.Product{}).Preload("CategoryInfo").Where("is_active = ?", true)

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// Transform to simplified response format
	var response []ProductListResponse
	for _, product := range products {
		productResponse := ProductListResponse{
			ID:            product.ID,
			NameEn:        product.NameEn,
			NameAr:        product.NameAr,
			DescriptionEn: product.DescriptionEn,
			DescriptionAr: product.DescriptionAr,
			ImageURL:      product.ImageURL,
		}

		// Add category names if category info is available
		if product.CategoryInfo != nil {
			productResponse.CategoryEn = product.CategoryInfo.NameEn
			productResponse.CategoryAr = product.CategoryInfo.NameAr
		}

		response = append(response, productResponse)
	}

	JSONResponse(c, http.StatusOK, gin.H{
		"data": response,
		"total": len(response),
	})
}

// GetProduct retrieves a single product by ID
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.GetDB().Preload("CategoryInfo").First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	// Populate virtual category field for backward compatibility
	if product.CategoryInfo != nil {
		product.Category = product.CategoryInfo.NameEn
	}

	JSONResponse(c, http.StatusOK, gin.H{"data": product})
}

// parseProductFormData parses multipart form data for product requests
func parseProductFormData(c *gin.Context) (*models.ProductCreateRequest, error) {
	req := &models.ProductCreateRequest{}

	// Debug: Print all form values
	fmt.Println("=== Form Data Received ===")
	if err := c.Request.ParseMultipartForm(32 << 20); err == nil {
		for key, values := range c.Request.MultipartForm.Value {
			fmt.Printf("%s: %v\n", key, values)
		}
	}
	fmt.Println("=========================")

	// Extract text fields
	req.NameEn = c.PostForm("name_en")
	req.NameAr = c.PostForm("name_ar")
	req.DescriptionEn = c.PostForm("description_en")
	req.DescriptionAr = c.PostForm("description_ar")
	req.ImageURL = c.PostForm("image_url")

	fmt.Printf("Extracted fields: NameEn=%s, NameAr=%s, DescriptionEn=%s, DescriptionAr=%s, ImageURL=%s\n",
		req.NameEn, req.NameAr, req.DescriptionEn, req.DescriptionAr, req.ImageURL)

	// Parse category_id
	categoryIDStr := c.PostForm("category_id")
	fmt.Printf("Category ID string: '%s'\n", categoryIDStr)
	if categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			req.CategoryID = uint(categoryID)
			fmt.Printf("Parsed category ID: %d\n", req.CategoryID)
		} else {
			fmt.Printf("Error parsing category_id '%s': %v\n", categoryIDStr, err)
			return nil, fmt.Errorf("invalid category_id format: %s", categoryIDStr)
		}
	}

	// Parse is_active
	if isActiveStr := c.PostForm("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		} else {
			return nil, fmt.Errorf("invalid is_active format: %s", isActiveStr)
		}
	}

	return req, nil
}

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	// Check content type to determine how to parse the request
	contentType := c.GetHeader("Content-Type")
	fmt.Printf("=== CreateProduct Debug ===\n")
	fmt.Printf("Content-Type: '%s'\n", contentType)
	fmt.Printf("Method: %s\n", c.Request.Method)
	fmt.Printf("URL: %s\n", c.Request.URL.String())

	var req *models.ProductCreateRequest
	var err error

	if strings.Contains(contentType, "multipart/form-data") {
		fmt.Println("Using multipart form data parsing")
		// Handle multipart form data (with potential image upload)
		req, err = parseProductFormData(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Handle JSON request (backward compatibility)
		fmt.Println("Using JSON parsing")
		var jsonReq models.ProductCreateRequest
		if err := c.ShouldBindJSON(&jsonReq); err != nil {
			fmt.Printf("JSON parsing error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req = &jsonReq
	}

	// Verify category exists
	var category models.Category
	if err := database.GetDB().First(&category, req.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
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

	product := models.Product{
		NameEn:        req.NameEn,
		NameAr:        req.NameAr,
		DescriptionEn: req.DescriptionEn,
		DescriptionAr: req.DescriptionAr,
		CategoryID:    req.CategoryID,
		ImageURL:      imageURL,
		IsActive:      true,
	}

	// Override with provided values if present
	if req.ImageURL != "" && imageURL == "" {
		product.ImageURL = req.ImageURL
	}

	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := database.GetDB().Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	JSONResponse(c, http.StatusCreated, gin.H{"data": product})
}

// parseProductUpdateFormData parses multipart form data for product update requests
func parseProductUpdateFormData(c *gin.Context) (*models.ProductUpdateRequest, error) {
	req := &models.ProductUpdateRequest{}

	// Extract text fields
	if nameEn := c.PostForm("name_en"); nameEn != "" {
		req.NameEn = &nameEn
	}
	if nameAr := c.PostForm("name_ar"); nameAr != "" {
		req.NameAr = &nameAr
	}
	if descriptionEn := c.PostForm("description_en"); descriptionEn != "" {
		req.DescriptionEn = &descriptionEn
	}
	if descriptionAr := c.PostForm("description_ar"); descriptionAr != "" {
		req.DescriptionAr = &descriptionAr
	}
	if imageURL := c.PostForm("image_url"); imageURL != "" {
		req.ImageURL = &imageURL
	}

	// Parse category_id
	if categoryIDStr := c.PostForm("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryIDUint := uint(categoryID)
			req.CategoryID = &categoryIDUint
		}
	}

	// Parse is_active
	if isActiveStr := c.PostForm("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		}
	}

	return req, nil
}

// UpdateProduct updates an existing product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.GetDB().First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	// Check content type to determine how to parse the request
	contentType := c.GetHeader("Content-Type")
	var req *models.ProductUpdateRequest
	var err error

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data (with potential image upload)
		req, err = parseProductUpdateFormData(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Handle JSON request (backward compatibility)
		var jsonReq models.ProductUpdateRequest
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

	// Verify category exists (if being updated)
	if req.CategoryID != nil {
		var category models.Category
		if err := database.GetDB().First(&category, *req.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
	}

	// Update fields if provided
	if req.NameEn != nil {
		product.NameEn = *req.NameEn
	}
	if req.NameAr != nil {
		product.NameAr = *req.NameAr
	}
	if req.DescriptionEn != nil {
		product.DescriptionEn = *req.DescriptionEn
	}
	if req.DescriptionAr != nil {
		product.DescriptionAr = *req.DescriptionAr
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := database.GetDB().Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	JSONResponse(c, http.StatusOK, gin.H{"data": product})
}

// DeleteProduct deletes a product (soft delete)
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.GetDB().First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	if err := database.GetDB().Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	JSONResponse(c, http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
