package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"golden-success-backend/database"
	"golden-success-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func main() {
	log.Println("Starting Golden Success Backend...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	log.Println("Environment variables loaded")

	// Initialize database
	log.Println("Initializing database...")
	database.InitDB()
	log.Println("Database initialized")

	// Create Gin router
	log.Println("Creating Gin router...")
	r := gin.Default()

	// Configure JSON encoder to not escape HTML
	r.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	// Configure CORS
	log.Println("Configuring CORS...")
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOriginFunc = func(origin string) bool {
		return true
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))


	// Setup routes
	log.Println("Setting up routes...")
	routes.SetupRoutes(r)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	log.Printf("API endpoints available at: http://localhost:%s/api/v1/", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
