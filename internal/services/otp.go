package services

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// GenerateOTP creates a 6-digit numeric OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// StoreOTP stores OTP in Memcached
func StoreOTP(mc *memcache.Client, key, otp string, ttl int) error {
	return mc.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(otp),
		Expiration: int32(ttl),
	})
}

// VerifyOTP checks OTP from Memcached
func VerifyOTP(mc *memcache.Client, key, otp string) bool {
	item, err := mc.Get(key)
	if err != nil {
		return false
	}
	return string(item.Value) == otp
}

// SendOTP sends OTP via email (or logs if devMode)
func SendOTP(devMode bool, email, otp string) {
	if devMode {
		log.Printf("[DEV MODE] OTP for %s: %s", email, otp)
		return
	}

	// Load SMTP info from environment variables
	smtpHost := getenv("SMTP_HOST", "")
	smtpPort := getenv("SMTP_PORT", "587")
	smtpUser := getenv("SMTP_USER", "")
	smtpPass := getenv("SMTP_PASS", "")
	from := getenv("SMTP_FROM", smtpUser)

	subject := "Your OTP Code"
	body := fmt.Sprintf("Hello,\n\nYour OTP code is: %s\nIt is valid for 5 minutes.", otp)
	msg := []byte("From: " + from + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, msg)
	if err != nil {
		log.Printf("Failed to send OTP to %s: %v", email, err)
		return
	}
	log.Printf("OTP sent to %s successfully", email)
}

// getenv helper
func getenv(key, fallback string) string {
	if v := getenvInternal(key); v != "" {
		return v
	}
	return fallback
}

// simple wrapper to get environment variable
func getenvInternal(key string) string {
	return getenvOS(key)
}

// can use os.Getenv
func getenvOS(key string) string {
	return os.Getenv(key)
}
