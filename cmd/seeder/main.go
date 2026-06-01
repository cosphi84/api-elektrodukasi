package main

import (
	"fmt"
	"log"
	"os"

	"elektrod/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Get database connection string
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=elektrodukasi port=5432 sslmode=disable"
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("✓ Connected to database")

	// Get password from environment
	password := os.Getenv("USER_PASSWORD")
	if password == "" {
		log.Fatal("USER_PASSWORD environment variable is not set")
	}

	// Seed user
	seedUser(db, password)

	// Seed category
	seedCategory(db)

	fmt.Println("\n✓ Seeding completed successfully")
}

func seedUser(db *gorm.DB, password string) {
	email := "risam1984@gmail.com"
	name := "Risam, S.T"

	// Check if user already exists
	var existingUser models.User
	result := db.Where("email = ? OR name = ?", email, name).First(&existingUser)

	if result.Error == nil {
		fmt.Printf("⊘ User '%s' already exists, skipping\n", email)
		return
	}

	if result.Error != gorm.ErrRecordNotFound {
		log.Fatalf("Error checking user existence: %v", result.Error)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create user
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		IsActive: true,
		Role:     "admin",
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("✓ User created: %s (%s) with role: admin\n", name, email)
}

func seedCategory(db *gorm.DB) {
	name := "Komponen Dasar"
	description := "Dasar Komponen elektronika"

	// Check if category already exists
	var existingCategory models.Category
	result := db.Where("name = ?", name).First(&existingCategory)

	if result.Error == nil {
		fmt.Printf("⊘ Category '%s' already exists, skipping\n", name)
		return
	}

	if result.Error != gorm.ErrRecordNotFound {
		log.Fatalf("Error checking category existence: %v", result.Error)
	}

	// Create category
	category := models.Category{
		Name:        name,
		Description: &description,
	}

	if err := db.Create(&category).Error; err != nil {
		log.Fatalf("Failed to create category: %v", err)
	}

	fmt.Printf("✓ Category created: %s\n", name)
}
