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
	cfg := config.LoadConfig()

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer dbConn.Close()

	mc := memcache.New(cfg.Memcache)

	r := gin.Default()

	// Register routes using *gin.Engine
	routes.RegisterRoutes(r, dbConn, mc, cfg)

	log.Println("Server running on port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
