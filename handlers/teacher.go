package handlers

import (
	"backend/models"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

type TeacherLoginRequest struct {
	TeacherID string `json:"teacher_id"`
}

type VerifyTeacherOTPRequest struct {
	UID string `json:"uid"`
	OTP string `json:"otp"`
}

// TeacherLoginHandler generates OTP
func TeacherLoginHandler(db *sql.DB, mc *memcache.Client, otpTTLSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TeacherLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.TeacherID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "teacher_id required"})
			return
		}

		var t models.Teacher
		err := db.QueryRow(`
			SELECT id, teacher_id, full_name, email, phone, school, dob, image
			FROM teachers WHERE teacher_id=$1
		`, req.TeacherID).Scan(&t.ID, &t.TeacherID, &t.FullName, &t.Email, &t.Phone, &t.School, &t.DOB, &t.Image)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "teacher not found"})
			return
		}

		uid := generateUID()
		otp := generateOTP()

		mc.Set(&memcache.Item{
			Key:        "teacher_otp_" + uid,
			Value:      []byte(otp),
			Expiration: int32(otpTTLSeconds),
		})

		c.JSON(http.StatusOK, gin.H{
			"uid":     uid,
			"otp":     otp,
			"message": "OTP sent successfully",
			"teacher": t,
		})
	}
}

// VerifyTeacherOTPHandler
func VerifyTeacherOTPHandler(mc *memcache.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyTeacherOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.UID == "" || req.OTP == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uid & otp required"})
			return
		}

		item, err := mc.Get("teacher_otp_" + req.UID)
		if err != nil || string(item.Value) != req.OTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid otp"})
			return
		}

		sessionToken := generateSessionToken()
		mc.Set(&memcache.Item{
			Key:        "teacher_session_" + req.UID,
			Value:      []byte(sessionToken),
			Expiration: int32(time.Hour.Seconds()),
		})

		c.JSON(http.StatusOK, gin.H{"session_token": sessionToken})
	}
}

// Helper functions (reuse from student.go)
func generateUID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateOTP() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)[:6]
}
