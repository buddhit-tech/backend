package main

import (
	"database/sql"
	"log"
	"os"
	"school-auth/internal/routes"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Env variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbSSL := os.Getenv("DB_SSLMODE")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	otpTTLSeconds := 300
	if val := os.Getenv("OTP_TTL_SECONDS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			otpTTLSeconds = parsed
		}
	}
	otpTTL := time.Duration(otpTTLSeconds) * time.Second

	// PostgreSQL connection
	psqlInfo := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " sslmode=" + dbSSL
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("DB Ping failed:", err)
	}
	log.Println("Connected to PostgreSQL")

	// Memcached connection
	memcacheAddr := os.Getenv("MEMCACHE_ADDR")
	if memcacheAddr == "" {
		memcacheAddr = "127.0.0.1:11211"
	}
	mc := memcache.New(memcacheAddr)
	if err := mc.Ping(); err != nil {
		log.Fatal("Memcached Ping failed:", err)
	}
	log.Println("Connected to Memcached")

	// Gin router
	router := gin.Default()

	// Register all routes
	routes.RegisterRoutes(router, db, mc, otpTTL)

	log.Println("Server running on port", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
