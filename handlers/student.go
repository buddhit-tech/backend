package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"school-auth/internal/models"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

// Request payloads
type StudentLoginRequest struct {
	StudentID string `json:"student_id"`
}

type VerifyStudentOTPRequest struct {
	UID string `json:"uid"`
	OTP string `json:"otp"`
}

// StudentLoginHandler generates a unique OTP for students
func StudentLoginHandler(db *sql.DB, mc *memcache.Client, otpTTLSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StudentLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.StudentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student_id required"})
			return
		}

		var s models.Student
		err := db.QueryRow(`
			SELECT id, student_id, full_name, email, phone, school, teacher_id, dob, image, class
			FROM students WHERE student_id=$1
		`, req.StudentID).Scan(&s.ID, &s.StudentID, &s.FullName, &s.Email, &s.Phone, &s.School, &s.TeacherID, &s.DOB, &s.Image, &s.Class)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}

		uid := StudentGenerateUID()
		otp := StudentGenerateOTP()

		// Store OTP in Memcached
		mc.Set(&memcache.Item{
			Key:        "student_otp_" + uid,
			Value:      []byte(otp),
			Expiration: int32(otpTTLSeconds),
		})

		c.JSON(http.StatusOK, gin.H{
			"uid":     uid,
			"otp":     otp, // for testing
			"message": "OTP sent successfully",
			"student": s,
		})
	}
}

// VerifyStudentOTPHandler verifies OTP and returns session token
func VerifyStudentOTPHandler(mc *memcache.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyStudentOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.UID == "" || req.OTP == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uid & otp required"})
			return
		}

		item, err := mc.Get("student_otp_" + req.UID)
		if err != nil || string(item.Value) != req.OTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid otp"})
			return
		}

		sessionToken := StudentGenerateSessionToken()
		mc.Set(&memcache.Item{
			Key:        "student_session_" + req.UID,
			Value:      []byte(sessionToken),
			Expiration: int32(time.Hour.Seconds()),
		})

		c.JSON(http.StatusOK, gin.H{"session_token": sessionToken})
	}
}

// ------------------
// Helpers (student-specific)
// ------------------
func StudentGenerateUID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func StudentGenerateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func StudentGenerateOTP() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)[:6]
}
