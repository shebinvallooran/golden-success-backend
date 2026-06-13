package routes

import (
	"golden-success-backend/controllers"
	"golden-success-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "Server is running"})
	})

	// Serve uploaded files
	r.Static("/uploads", "./uploads")

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Authentication routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", controllers.Login)           // POST /api/v1/auth/login
			auth.GET("/profile", middleware.AuthMiddleware(), controllers.GetProfile) // GET /api/v1/auth/profile
		}

		// Public product routes (no auth required)
		products := v1.Group("/products")
		{
			products.GET("", controllers.GetProducts)       // GET /api/v1/products
			products.GET("/list", controllers.GetProductList) // GET /api/v1/products/list (full list without pagination)
			products.GET("/:id", controllers.GetProduct)    // GET /api/v1/products/:id
		}

		// Protected product routes (auth required)
		protectedProducts := v1.Group("/products")
		protectedProducts.Use(middleware.AuthMiddleware())
		{
			protectedProducts.POST("", controllers.CreateProduct)    // POST /api/v1/products
			protectedProducts.PUT("/:id", controllers.UpdateProduct) // PUT /api/v1/products/:id
			protectedProducts.DELETE("/:id", controllers.DeleteProduct) // DELETE /api/v1/products/:id
		}

		// Public category routes (no auth required)
		categories := v1.Group("/categories")
		{
			categories.GET("", controllers.GetCategories)                    // GET /api/v1/categories
			categories.GET("/with-products", controllers.GetCategoriesWithProducts) // GET /api/v1/categories/with-products
			categories.GET("/home", controllers.GetCategoriesForHome)        // GET /api/v1/categories/home (for home page)
			categories.GET("/list", controllers.GetCategoriesList)           // GET /api/v1/categories/list (simple list without pagination)
			categories.GET("/:id", controllers.GetCategory)                 // GET /api/v1/categories/:id
		}
		// Public quotes/enquiries route (no auth required)
		v1.POST("/quotes", controllers.CreateQuote) // POST /api/v1/quotes

		// Protected category routes (auth required)
		protectedCategories := v1.Group("/categories")
		protectedCategories.Use(middleware.AuthMiddleware())
		{
			protectedCategories.POST("", controllers.CreateCategory)                 // POST /api/v1/categories
			protectedCategories.PUT("/:id", controllers.UpdateCategory)              // PUT /api/v1/categories/:id
			protectedCategories.DELETE("/:id", controllers.DeleteCategory)           // DELETE /api/v1/categories/:id
			protectedCategories.PUT("/priorities", controllers.UpdateCategoryPriorities) // PUT /api/v1/categories/priorities
		}

		// Protected routes (auth required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Statistics routes
			protected.GET("/statistics", controllers.GetStatistics)     // GET /api/v1/statistics

			// Upload routes
			protected.POST("/upload", controllers.UploadImage)           // POST /api/v1/upload
			protected.DELETE("/uploads/:filename", controllers.DeleteImage) // DELETE /api/v1/uploads/:filename

			// Quotes/Enquiries routes
			protected.GET("/quotes", controllers.GetQuotes)                // GET /api/v1/quotes
			protected.PUT("/quotes/:id/status", controllers.UpdateQuoteStatus) // PUT /api/v1/quotes/:id/status
			protected.DELETE("/quotes/:id", controllers.DeleteQuote)       // DELETE /api/v1/quotes/:id

			// Notification settings routes
			protected.GET("/settings/notification", controllers.GetNotificationSettings) // GET /api/v1/settings/notification
			protected.PUT("/settings/notification", controllers.UpdateNotificationSettings) // PUT /api/v1/settings/notification
		}

		// Public upload route (no auth required for serving images)
		v1.GET("/uploads/:filename", controllers.ServeImage)  // GET /api/v1/uploads/:filename
	}
}
