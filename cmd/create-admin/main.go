package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"raven.go.invoice-builder/internal/models"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	emailFlag := flag.String("email", "", "Admin email address")
	passwordFlag := flag.String("password", "", "Admin password")
	nameFlag := flag.String("name", "System Admin", "Admin user name")
	flag.Parse()

	email := *emailFlag
	password := *passwordFlag
	name := *nameFlag

	if email == "" {
		email = "admin@example.com"
	}
	if password == "" {
		password = "AdminPassword123!"
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "raven_db")
	dbName := getEnv("DB_NAME", "invoice-app")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()

	if err := db.AutoMigrate(&models.Role{}, &models.Permission{}, &models.User{}, &models.CreditPackage{}, &models.AuditLog{}, &models.WebhookLog{}, &models.SystemSetting{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	var adminRole models.Role
	err = db.WithContext(ctx).Where("name = ?", "Admin").First(&adminRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			adminRole = models.Role{
				ID:          uuid.New(),
				Name:        "Admin",
				Description: "System Administrator with elevated privileges",
			}
			if err := db.WithContext(ctx).Create(&adminRole).Error; err != nil {
				log.Fatalf("Failed to create Admin role: %v", err)
			}
			fmt.Println("Created 'Admin' role.")
		} else {
			log.Fatalf("Failed to query Admin role: %v", err)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	var existingUser models.User
	err = db.WithContext(ctx).Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		existingUser.RoleID = adminRole.ID
		existingUser.PasswordHash = string(hashedPassword)
		existingUser.IsActivated = true
		existingUser.IsProfileComplete = true
		if name != "" && name != "System Admin" {
			existingUser.Name = name
		}
		if err := db.WithContext(ctx).Save(&existingUser).Error; err != nil {
			log.Fatalf("Failed to promote user to Admin: %v", err)
		}
		fmt.Printf("Successfully promoted existing user '%s' (%s) to Admin role.\n", existingUser.Name, email)
		return
	}

	newUser := models.User{
		Name:              name,
		Email:             email,
		PasswordHash:      string(hashedPassword),
		RoleID:            adminRole.ID,
		IsActivated:       true,
		IsProfileComplete: true,
	}
	newUser.ID = uuid.New()

	if err := db.WithContext(ctx).Create(&newUser).Error; err != nil {
		log.Fatalf("Failed to create Admin user: %v", err)
	}

	fmt.Printf("Successfully created new Admin user:\n  Name:     %s\n  Email:    %s\n  Password: %s\n", newUser.Name, newUser.Email, password)
}
