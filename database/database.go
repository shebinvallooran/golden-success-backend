package database

import (
	"log"
	"os"

	"golden-success-backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func InitDB() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "golden_success.db"
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to SQLite database: %v", err)
	}
	log.Printf("SQLite database connected successfully (file: %s)", dbPath)

	log.Println("Running database migrations...")
	err = DB.AutoMigrate(
		&models.Category{},
		&models.User{},
		&models.Product{},
		&models.Quote{},
		&models.NotificationSetting{},
	)
	if err != nil {
		log.Fatalf("FATAL: Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	createDefaultCategories()
	createDefaultUser()
	createDefaultSettings()
}

func createDefaultCategories() {
	var categoryCount int64
	DB.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		defaultCategories := []models.Category{}

		if len(defaultCategories) > 0 {
			if err := DB.Create(&defaultCategories).Error; err != nil {
				log.Printf("WARN: Could not create default categories: %v", err)
			} else {
				log.Println("Default categories created")
			}
		} else {
			log.Println("Skipping default category creation as none are defined.")
		}
	}
}

func createDefaultUser() {
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin1223"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("WARN: Failed to hash password for default user: %v", err)
			return
		}

		defaultUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			IsActive: true,
		}

		if err := DB.Create(&defaultUser).Error; err != nil {
			log.Printf("WARN: Failed to create default user: %v", err)
		} else {
			log.Println("Default admin user created (username: admin, password: admin1223)")
		}
	}
}

func createDefaultSettings() {
	var count int64
	DB.Model(&models.NotificationSetting{}).Count(&count)
	if count == 0 {
		defaultSettings := models.NotificationSetting{
			NotificationEmail:        "admin@goldensuccessksa.com",
			EnableEmailNotifications: false,
			SenderEmail:              "noreply@goldensuccessksa.com",
			SMTPHost:                 "smtp.gmail.com",
			SMTPPort:                 587,
			SMTPUsername:             "",
			SMTPPassword:             "",
			SMTPSecure:               false,
		}

		if err := DB.Create(&defaultSettings).Error; err != nil {
			log.Printf("WARN: Failed to create default settings: %v", err)
		} else {
			log.Println("Default notification settings created")
		}
	}
}

func GetDB() *gorm.DB {
	return DB
}

