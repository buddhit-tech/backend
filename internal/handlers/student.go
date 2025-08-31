package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

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

// StudentLoginHandler generates OTP for students
func StudentLoginHandler(db any, mc *memcache.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StudentLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.StudentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student_id required"})
			return
		}

		uid := generateUID()
		otp := "654321" // demo OTP

		if err := mc.Set(&memcache.Item{
			Key:        "student_otp_" + uid,
			Value:      []byte(otp),
			Expiration: 300,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store otp"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"uid": uid, "msg": "OTP sent (mocked)"})
	}
}

// VerifyStudentOTPHandler verifies OTP and returns session token
func VerifyStudentOTPHandler(db any, mc *memcache.Client) gin.HandlerFunc {
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

		sessionToken := GenerateSessionToken()
		if err := mc.Set(&memcache.Item{
			Key:        "student_session_" + req.UID,
			Value:      []byte(sessionToken),
			Expiration: int32(time.Hour.Seconds()),
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"session_token": sessionToken})
	}
}

// Helpers
func generateUID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
