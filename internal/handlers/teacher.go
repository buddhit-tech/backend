package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"school-auth/internal/services"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

type TeacherLoginRequest struct {
	TeacherID string `json:"teacher_id"`
}

// TeacherLogin - POST /teacher/login
func TeacherLogin(db *sql.DB, mc *memcache.Client, devMode bool, otpTTL int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TeacherLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		var email string
		if err := db.QueryRow(`SELECT email FROM teachers WHERE teacher_id=$1`, req.TeacherID).Scan(&email); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "teacher not found"})
			return
		}

		otp := services.GenerateOTP()
		key := "teacher:" + req.TeacherID
		if err := services.StoreOTP(mc, key, otp, otpTTL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store otp"})
			return
		}
		services.SendOTP(devMode, email, otp)

		c.JSON(http.StatusOK, gin.H{"message": "OTP sent (or logged in dev mode)", "uid": req.TeacherID})
	}
}

// TeacherVerifyOTP - POST /teacher/otp/verify
func TeacherVerifyOTP(db *sql.DB, mc *memcache.Client, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		key := "teacher:" + req.UID
		if !services.VerifyOTP(mc, key, req.OTP) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired otp"})
			return
		}

		var name, email string
		if err := db.QueryRow(`SELECT full_name, email FROM teachers WHERE teacher_id=$1`, req.UID).Scan(&name, &email); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		token, err := services.CreateToken(jwtSecret, req.UID, "teacher", name, email, 24*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "OTP verified",
			"token":   token,
			"user":    gin.H{"uid": req.UID, "role": "teacher", "name": name, "email": email},
		})
	}
}
