package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"school-auth/internal/services"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

// Request payloads
type StudentLoginRequest struct {
	StudentID string `json:"student_id"`
}

type VerifyOTPRequest struct {
	UID string `json:"uid"`
	OTP string `json:"otp"`
}

// StudentLogin - POST /student/login
func StudentLogin(db *sql.DB, mc *memcache.Client, devMode bool, otpTTL int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StudentLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		var email string
		if err := db.QueryRow(`SELECT email FROM students WHERE student_id=$1`, req.StudentID).Scan(&email); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}

		otp := services.GenerateOTP()
		key := "student:" + req.StudentID
		if err := services.StoreOTP(mc, key, otp, otpTTL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store otp"})
			return
		}
		services.SendOTP(devMode, email, otp)

		c.JSON(http.StatusOK, gin.H{"message": "OTP sent (or logged in dev mode)", "uid": req.StudentID})
	}
}

// StudentVerifyOTP - POST /student/otp/verify
func StudentVerifyOTP(db *sql.DB, mc *memcache.Client, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		key := "student:" + req.UID
		if !services.VerifyOTP(mc, key, req.OTP) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired otp"})
			return
		}

		var name, email string
		if err := db.QueryRow(`SELECT full_name, email FROM students WHERE student_id=$1`, req.UID).Scan(&name, &email); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		token, err := services.CreateToken(jwtSecret, req.UID, "student", name, email, 24*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "OTP verified",
			"token":   token,
			"user":    gin.H{"uid": req.UID, "role": "student", "name": name, "email": email},
		})
	}
}
