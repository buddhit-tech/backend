package routes

import (
	"database/sql"

	"school-auth/internal/config"
	"school-auth/internal/handlers"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB, mc *memcache.Client, cfg *config.Config) {
    // Student routes
    r.POST("/student/login", handlers.StudentLogin(db, mc, cfg.DevMode, cfg.OTPTTL))
    r.POST("/student/otp/verify", handlers.StudentVerifyOTP(db, mc, []byte(cfg.JWTSecret)))

    // Teacher routes
    r.POST("/teacher/login", handlers.TeacherLogin(db, mc, cfg.DevMode, cfg.OTPTTL))
    r.POST("/teacher/otp/verify", handlers.TeacherVerifyOTP(db, mc, []byte(cfg.JWTSecret)))

    // Health check
    r.GET("/healthz", func(c *gin.Context) {
        c.String(200, "ok")
    })
}

