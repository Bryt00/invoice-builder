package main

import (
	"context"
	"encoding/gob"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"raven.go.invoice-builder/internal/jsonlog"
	"raven.go.invoice-builder/internal/mailer"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/paystack"
	"raven.go.invoice-builder/internal/services"
	"raven.go.invoice-builder/internal/tokens"
)

type config struct {
	port    int
	env     string
	frontendURL string
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	paystack struct {
		secretKey string
		publicKey string
	}
	jwt struct {
		secret string
	}
	admin struct {
		secretKey  string
		secretPath string
	}
}

type application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	jsonLogger *jsonlog.Logger
	jwtManager *tokens.JWTManager
	models     models.Models
	services   services.Services
	apiHandler *ApiHandler
	mailer     mailer.Mailer
	paystack   *paystack.Client
	config     config
}

// getEnv reads an environment variable or returns a fallback default.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvInt reads an environment variable as int or returns a fallback default.
func getEnvInt(key string, fallback int) int {
	s := getEnv(key, "")
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}

func main() {
	// Load .env file (non-fatal if missing, e.g. in production)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading config from environment")
	}

	gob.Register(uuid.UUID{})

	// Build config from environment variables
	var cfg config
	cfg.port = getEnvInt("APP_PORT", 4000)
	cfg.env = getEnv("APP_ENV", "development")
	cfg.frontendURL = getEnv("FRONTEND_URL", "http://localhost:5173")

	cfg.smtp.host = getEnv("SMTP_HOST", "mail.teckstyle.com")
	cfg.smtp.port = getEnvInt("SMTP_PORT", 465)
	cfg.smtp.username = getEnv("SMTP_USERNAME", "no-reply@teckstyle.com")
	cfg.smtp.password = getEnv("SMTP_PASSWORD", "")
	cfg.smtp.sender = getEnv("SMTP_SENDER", "Teks-Invoice <no-reply@teckstyle.com>")

	cfg.paystack.secretKey = getEnv("PAYSTACK_SECRET_KEY", "sk_test_mock_secret_key")
	cfg.paystack.publicKey = getEnv("PAYSTACK_PUBLIC_KEY", "pk_test_mock_public_key")

	cfg.admin.secretKey = getEnv("ADMIN_SECRET_KEY", "teks-admin-secret-key-2026")
	cfg.admin.secretPath = getEnv("ADMIN_SECRET_PATH", "admin")

	cfg.jwt.secret = getEnv("JWT_SECRET", "")
	if cfg.jwt.secret == "" {
		if cfg.env == "production" {
			log.Fatal("FATAL: JWT_SECRET environment variable must be set in production")
		}
		cfg.jwt.secret = "dev-only-insecure-fallback-key"
		log.Println("WARNING: Using insecure fallback JWT secret (development only)")
	}

	limiterEnabled := getEnv("LIMITER_ENABLED", "true")
	cfg.limiter.enabled = limiterEnabled == "true"
	rps, err := strconv.ParseFloat(getEnv("LIMITER_RPS", "2"), 64)
	if err != nil {
		rps = 2
	}
	cfg.limiter.rps = rps
	cfg.limiter.burst = getEnvInt("LIMITER_BURST", 4)

	// Logging Configuration
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("failed to create logs directory: %v", err)
	}

	infoLogFile, err := os.OpenFile("logs/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open info.log: %v", err)
	}
	defer infoLogFile.Close()

	errorLogFile, err := os.OpenFile("logs/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open error.log: %v", err)
	}
	defer errorLogFile.Close()

	apiLogFile, err := os.OpenFile("logs/api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open api.log: %v", err)
	}
	defer apiLogFile.Close()

	infoWriter := io.MultiWriter(os.Stdout, infoLogFile)
	errorWriter := io.MultiWriter(os.Stderr, errorLogFile)
	apiWriter := io.MultiWriter(os.Stdout, apiLogFile)

	infoLog := log.New(infoWriter, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(errorWriter, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
	jsonLogger := jsonlog.New(apiWriter, jsonlog.LevelInfo)
	jwtManager := tokens.NewJWTManager(cfg.jwt.secret, 3*24*time.Hour)

	// Database
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "raven_db")
	dbName := getEnv("DB_NAME", "invoice-app")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		errorLog.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		errorLog.Fatal(err)
		return
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// AutoMigrate all models
	err = db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.BusinessProfile{},
		&models.Client{},
		&models.Invoice{},
		&models.LineItem{},
		&models.Receipt{},
		&models.Payment{},
		&models.CreditTxn{},
		&models.CreditPackage{},
		&models.AuditLog{},
		&models.WebhookLog{},
		&models.SystemSetting{},
		&models.FinancialCategory{},
		&models.FinancialTransaction{},
		&models.RefreshToken{},
	)
	if err != nil {
		errorLog.Fatalf("AutoMigrate failed: %v", err)
	}
	infoLog.Println("Database schema migrated successfully")

	// Application setup
	addr := ":" + strconv.Itoa(cfg.port)
	appModels := models.NewModel(db)
	appMailer := mailer.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	appPaystack := paystack.NewClient(cfg.paystack.secretKey)
	
	appServices := services.NewServices(appModels, appMailer, appPaystack, jwtManager, cfg.frontendURL)
	apiHandler := NewApiHandler(appServices, appModels, jsonLogger, jwtManager, cfg.frontendURL)

	app := &application{
		errorLog:   errorLog,
		infoLog:    infoLog,
		jsonLogger: jsonLogger,
		jwtManager: jwtManager,
		models:     appModels,
		services:   appServices,
		apiHandler: apiHandler,
		config:     cfg,
		mailer:     appMailer,
		paystack:   appPaystack,
	}

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting API server on %s", addr)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errorLog.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	infoLog.Println("Shutting down API server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		errorLog.Fatalf("Server forced to shutdown: %v", err)
	}

	infoLog.Println("API server exiting")
}
