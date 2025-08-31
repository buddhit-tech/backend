package main

import (
	"log"

	"school-auth/internal/config"
	"school-auth/internal/db"
	"school-auth/internal/routes"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to DB
	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer dbConn.Close()

	// Initialize Memcache client (expects address like "localhost:11211")
	mc := memcache.New(cfg.MemcacheAddr)
	if mc == nil {
		log.Fatal("Failed to initialize memcache client")
	}

	// Setup Gin router
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r, dbConn, mc, cfg)

	// Start server
	log.Println("✅ Server running on port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Server failed:", err)
	}
}
