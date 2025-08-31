package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"school-auth/internal/models"
	"school-auth/internal/services"

	"github.com/bradfitz/gomemcache/memcache"
)

type TeacherLoginRequest struct {
	TeacherID string `json:"teacher_id"`
}

// TeacherLogin - POST /teacher/login
func TeacherLogin(db *sql.DB, mc *memcache.Client, devMode bool, otpTTL int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TeacherLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		var email string
		if err := db.QueryRow(`SELECT email FROM teachers WHERE teacher_id=$1`, req.TeacherID).Scan(&email); err != nil {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}

		otp := services.GenerateOTP()
		key := "teacher:" + req.TeacherID

		if err := services.StoreOTP(mc, key, otp, otpTTL); err != nil {
			http.Error(w, "failed to store otp", http.StatusInternalServerError)
			return
		}

		services.SendOTP(devMode, email, otp)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "OTP sent (or logged in dev mode)",
			"uid":     req.TeacherID,
		})
	}
}

// TeacherVerifyOTP - POST /teacher/otp/verify
func TeacherVerifyOTP(db *sql.DB, mc *memcache.Client, jwtSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.VerifyOTPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		key := "teacher:" + req.UID
		if !services.VerifyOTP(mc, key, req.OTP) {
			http.Error(w, "invalid or expired otp", http.StatusUnauthorized)
			return
		}

		var name, email string
		if err := db.QueryRow(`SELECT full_name, email FROM teachers WHERE teacher_id=$1`, req.UID).Scan(&name, &email); err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		token, err := services.CreateToken(jwtSecret, req.UID, "teacher", name, email, 24*time.Hour)
		if err != nil {
			http.Error(w, "failed to create token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "OTP verified",
			"token":   token,
			"user": map[string]string{
				"uid":   req.UID,
				"role":  "teacher",
				"name":  name,
				"email": email,
			},
		})
	}
}
