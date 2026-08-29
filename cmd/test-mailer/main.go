package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"raven.go.invoice-builder/internal/mailer"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading config from environment")
	}

	toFlag := flag.String("to", "mrbright63@gmail.com", "Recipient email address")
	flag.Parse()

	host := getEnv("SMTP_HOST", "mail.teckstyle.com")
	port := getEnvInt("SMTP_PORT", 465)
	username := getEnv("SMTP_USERNAME", "no-reply@teckstyle.com")
	password := getEnv("SMTP_PASSWORD", "")
	sender := getEnv("SMTP_SENDER", "Teks-Invoice <no-reply@teckstyle.com>")

	fmt.Printf("Testing SMTP Mailer with Configuration:\n")
	fmt.Printf("  Host:     %s\n", host)
	fmt.Printf("  Port:     %d\n", port)
	fmt.Printf("  Username: %s\n", username)
	fmt.Printf("  Sender:   %s\n", sender)
	fmt.Printf("  To:       %s\n\n", *toFlag)

	m := mailer.NewMailer(host, port, username, password, sender)

	err := m.SendMail(*toFlag, "user_welcome.tmpl", map[string]any{
		"Name":            "Test User",
		"ActivationToken": "test-token-123456",
	})
	if err != nil {
		log.Fatalf("FAILED to send test email: %v\n", err)
	}

	fmt.Println("SUCCESS! Test email dispatched successfully.")
}
