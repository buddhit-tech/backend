package main

import (
	"log"
	"net/http"

	"school-auth/internal/config"
	"school-auth/internal/db"
	"school-auth/internal/routes"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	cfg := config.LoadConfig()

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer dbConn.Close()
	log.Println("✅ Connected to PostgreSQL successfully")

	mc := memcache.New(cfg.Memcache)
	log.Println("✅ Memcache client initialized")

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, dbConn, mc, cfg)

	log.Println("Server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
