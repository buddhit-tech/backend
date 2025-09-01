package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"school-auth/models"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

// ------------------
// Request payloads
// ------------------
type StudentLoginRequest struct {
	StudentID string `json:"student_id"`
}

type VerifyStudentOTPRequest struct {
	UID string `json:"uid"`
	OTP string `json:"otp"`
}

// ------------------
// Handlers
// ------------------

// StudentLoginHandler generates a unique OTP for students
func StudentLoginHandler(db *sql.DB, mc *memcache.Client, otpTTLSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StudentLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.StudentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student_id required"})
			return
		}

		student, err := fetchStudentByID(db, req.StudentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}

		uid, otp := generateAndStoreStudentOTP(mc, otpTTLSeconds)

		c.JSON(http.StatusOK, gin.H{
			"uid":     uid,
			"otp":     otp, // for testing only
			"message": "OTP sent successfully",
			"student": student,
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

		if !verifyStudentOTP(mc, req.UID, req.OTP) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid otp"})
			return
		}

		sessionToken := generateAndStoreStudentSession(mc, req.UID)
		c.JSON(http.StatusOK, gin.H{"session_token": sessionToken})
	}
}

// ------------------
// DB Helper
// ------------------
func fetchStudentByID(db *sql.DB, studentID string) (models.Student, error) {
	var s models.Student
	err := db.QueryRow(`
		SELECT id, student_id, full_name, email, phone, school, teacher_id, dob, image, class
		FROM students WHERE student_id=$1
	`, studentID).Scan(&s.ID, &s.StudentID, &s.FullName, &s.Email, &s.Phone, &s.School, &s.TeacherID, &s.DOB, &s.Image, &s.Class)
	return s, err
}

// ------------------
// OTP & Session Helpers (student-specific)
// ------------------
func generateStudentUID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateStudentSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateStudentOTP() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)[:6]
}

func generateAndStoreStudentOTP(mc *memcache.Client, ttlSeconds int) (string, string) {
	uid := generateStudentUID()
	otp := generateStudentOTP()

	mc.Set(&memcache.Item{
		Key:        "student_otp_" + uid,
		Value:      []byte(otp),
		Expiration: int32(ttlSeconds),
	})

	return uid, otp
}

func verifyStudentOTP(mc *memcache.Client, uid, otp string) bool {
	item, err := mc.Get("student_otp_" + uid)
	if err != nil || string(item.Value) != otp {
		return false
	}
	return true
}

func generateAndStoreStudentSession(mc *memcache.Client, uid string) string {
	sessionToken := generateStudentSessionToken()
	mc.Set(&memcache.Item{
		Key:        "student_session_" + uid,
		Value:      []byte(sessionToken),
		Expiration: int32(time.Hour.Seconds()),
	})
	return sessionToken
}
