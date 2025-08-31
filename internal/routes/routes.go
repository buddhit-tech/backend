package routes

import (
	"database/sql"
	"net/http"

	"school-auth/internal/config"
	"school-auth/internal/handlers"

	"github.com/bradfitz/gomemcache/memcache"
)

// RegisterRoutes sets up all HTTP routes
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, mc *memcache.Client, cfg *config.Config) {
	// Student routes
	mux.HandleFunc("/student/login", handlers.StudentLogin(db, mc, cfg.DevMode, cfg.OTPTTL))
	mux.HandleFunc("/student/otp/verify", handlers.StudentVerifyOTP(db, mc, []byte(cfg.JWTSecret)))

	// Teacher routes
	mux.HandleFunc("/teacher/login", handlers.TeacherLogin(db, mc, cfg.DevMode, cfg.OTPTTL))
	mux.HandleFunc("/teacher/otp/verify", handlers.TeacherVerifyOTP(db, mc, []byte(cfg.JWTSecret)))

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}
