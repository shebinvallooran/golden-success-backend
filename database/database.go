package database

import (
	"fmt"
	"log"

	"golden-success-backend/config"
	"golden-success-backend/models"

	"gorm.io/driver/postgres"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.GetConfig()

	var err error

	// Check if we should use SQLite (for development/testing)
	if cfg.DBName == "sqlite" || cfg.DBHost == "sqlite" {
		// Use SQLite for easier setup with pure Go driver
		DB, err = gorm.Open(sqlite.Open("golden_success.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to SQLite database:", err)
		}
		log.Println("SQLite database connected successfully")
	} else {
		// Use PostgreSQL
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Failed to connect to PostgreSQL: %v", err)
			log.Println("Falling back to SQLite...")

			// Fallback to SQLite
			DB, err = gorm.Open(sqlite.Open("golden_success.db"), &gorm.Config{})
			if err != nil {
				log.Fatal("Failed to connect to any database:", err)
			}
			log.Println("SQLite database connected successfully (fallback)")
		} else {
			log.Println("PostgreSQL database connected successfully")
		}
	}

	// Handle category migration carefully
	migrateCategories()

	// Handle user migration
	migrateUsers()

	// Create default categories if none exist
	var categoryCount int64
	DB.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		defaultCategories := []models.Category{
			{NameEn: "Electronics", DescriptionEn: "Electronic devices and gadgets", Priority: 1, IsActive: true},
			{NameEn: "Clothing", DescriptionEn: "Apparel and fashion items", Priority: 2, IsActive: true},
			{NameEn: "Books", DescriptionEn: "Books and educational materials", Priority: 3, IsActive: true},
			{NameEn: "Home & Garden", DescriptionEn: "Home improvement and garden items", Priority: 4, IsActive: true},
			{NameEn: "Sports", DescriptionEn: "Sports and fitness equipment", Priority: 5, IsActive: true},
		}

		for _, category := range defaultCategories {
			DB.Create(&category)
		}
		log.Println("Default categories created")
	}

	// Create default admin user if none exist
	createDefaultUser()

	// Handle product migration carefully
	migrateProducts()

	log.Println("Database migration completed")
}

func migrateCategories() {
	// Check if categories table exists
	var tableExists int
	DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='categories'").Scan(&tableExists)

	if tableExists > 0 {
		// Check if old 'name' column exists
		var nameColumnExists int
		DB.Raw("SELECT COUNT(*) FROM pragma_table_info('categories') WHERE name='name'").Scan(&nameColumnExists)

		if nameColumnExists > 0 {
			log.Println("Migrating categories from old schema...")
			
			// Add new columns if they don't exist
			var nameEnExists int
			DB.Raw("SELECT COUNT(*) FROM pragma_table_info('categories') WHERE name='name_en'").Scan(&nameEnExists)
			if nameEnExists == 0 {
				DB.Exec("ALTER TABLE categories ADD COLUMN name_en TEXT")
				DB.Exec("ALTER TABLE categories ADD COLUMN name_ar TEXT")
				DB.Exec("ALTER TABLE categories ADD COLUMN description_en TEXT")
				DB.Exec("ALTER TABLE categories ADD COLUMN description_ar TEXT")
				
				// Migrate data from old columns to new columns
				DB.Exec("UPDATE categories SET name_en = name WHERE name_en IS NULL OR name_en = ''")
				DB.Exec("UPDATE categories SET description_en = description WHERE description_en IS NULL OR description_en = ''")
				log.Println("Category data migrated successfully")
			}
		}
	}

	// Now run the full migration
	err := DB.AutoMigrate(&models.Category{})
	if err != nil {
		log.Printf("Migration warning: %v", err)
		// Don't fail completely, just log the warning
	}
}

func migrateProducts() {
	// Check if products table exists and if it has category_id column
	var hasTable bool
	DB.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name='products'").Scan(&hasTable)

	if hasTable {
		// Check if category_id column exists
		var columnExists bool
		DB.Raw("SELECT COUNT(*) FROM pragma_table_info('products') WHERE name='category_id'").Scan(&columnExists)

		if !columnExists {
			log.Println("Adding category_id column to existing products...")

			// Add the column as nullable first
			DB.Exec("ALTER TABLE products ADD COLUMN category_id INTEGER")

			// Get the first category ID to use as default
			var firstCategoryID uint
			DB.Model(&models.Category{}).Select("id").First(&firstCategoryID)

			if firstCategoryID > 0 {
				// Update all existing products to use the first category
				DB.Exec("UPDATE products SET category_id = ? WHERE category_id IS NULL", firstCategoryID)
				log.Printf("Updated existing products to use category ID %d", firstCategoryID)
			}
		}
	}

	// Now run the full migration
	err := DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate products:", err)
	}

	log.Println("Database migration completed")
}

func migrateUsers() {
	// Auto migrate users table
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("User migration warning: %v", err)
	}
}

func createDefaultUser() {
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		// Create default admin user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			return
		}

		defaultUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			IsActive: true,
		}

		if err := DB.Create(&defaultUser).Error; err != nil {
			log.Printf("Failed to create default user: %v", err)
		} else {
			log.Println("Default admin user created (username: admin, password: admin123)")
		}
	}
}

func GetDB() *gorm.DB {
	return DB
}
