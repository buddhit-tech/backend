package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBHost       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	MemcacheAddr string
	JWTSecret    string
	DevMode      bool
	OTPTTL       int
	DBURL        string
}

// LoadConfig loads .env and environment variables
func LoadConfig() *Config {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	} else {
		log.Println("✅ .env loaded successfully")
	}

	// OTP TTL
	otpTTL := 300
	if v := os.Getenv("OTP_TTL_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			otpTTL = n
		}
	}

	// Dev mode
	devMode := os.Getenv("DEV_MODE") == "true"

	cfg := &Config{
		Port:         mustGetEnv("PORT"),
		DBHost:       mustGetEnv("DB_HOST"),
		DBUser:       mustGetEnv("POSTGRES_USER"),
		DBPassword:   mustGetEnv("POSTGRES_PASSWORD"),
		DBName:       mustGetEnv("POSTGRES_DB"),
		DBSSLMode:    mustGetEnv("DB_SSLMODE"),
		MemcacheAddr: mustGetEnv("MEMCACHE_ADDR"), // 👈 fixed
		JWTSecret:    mustGetEnv("JWT_SECRET"),
		DevMode:      devMode,
		OTPTTL:       otpTTL,
		DBURL:        os.Getenv("DATABASE_URL"), // optional
	}

	// DEBUG: print config (hide password)
	log.Printf("✅ DB Config -> host: %s user: %s dbname: %s sslmode: %s\n",
		cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBSSLMode)

	return cfg
}

// mustGetEnv returns env value or fatally exits if empty
func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return val
}
