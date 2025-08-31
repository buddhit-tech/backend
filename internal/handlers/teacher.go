package handlers

import (
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

// Request payloads
type TeacherLoginRequest struct {
	TeacherID string `json:"teacher_id"`
}

type VerifyTeacherOTPRequest struct {
	UID string `json:"uid"`
	OTP string `json:"otp"`
}

// TeacherLoginHandler generates OTP for teachers
func TeacherLoginHandler(db any, mc *memcache.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TeacherLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.TeacherID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "teacher_id required"})
			return
		}

		uid := "teacher-" + req.TeacherID + "-" + time.Now().Format("20060102150405")
		otp := "123456" // demo OTP

		if err := mc.Set(&memcache.Item{
			Key:        uid,
			Value:      []byte(otp),
			Expiration: 300,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store otp"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"uid": uid, "msg": "OTP sent successfully"})
	}
}

// VerifyTeacherOTPHandler verifies OTP for teachers
func VerifyTeacherOTPHandler(db any, mc *memcache.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyTeacherOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.UID == "" || req.OTP == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uid & otp required"})
			return
		}

		item, err := mc.Get(req.UID)
		if err != nil || string(item.Value) != req.OTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid otp"})
			return
		}

		sessionToken := GenerateSessionToken()
		if err := mc.Set(&memcache.Item{
			Key:        "teacher_session_" + req.UID,
			Value:      []byte(sessionToken),
			Expiration: int32(time.Hour.Seconds()),
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"session_token": sessionToken})
	}
}
