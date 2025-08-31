package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"school-auth/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN required")
	}

	memcacheAddr := os.Getenv("MEMCACHE_ADDR")
	if memcacheAddr == "" {
		memcacheAddr = "localhost:11211"
	}

	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}

	mc := memcache.New(memcacheAddr)

	r := gin.Default()

	// Teacher routes
	r.POST("/teachers/auth", handlers.TeacherLoginHandler(dbConn, mc))
	r.POST("/teachers/verify", handlers.VerifyTeacherOTPHandler(dbConn, mc))

	// Student routes
	r.POST("/students/auth", handlers.StudentLoginHandler(dbConn, mc))
	r.POST("/students/verify", handlers.VerifyStudentOTPHandler(dbConn, mc))

	log.Println("Server running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed:", err)
	}
}
