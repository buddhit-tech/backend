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

type StudentLoginRequest struct {
	StudentID string `json:"student_id"`
}

// StudentLogin - POST /student/login
func StudentLogin(db *sql.DB, mc *memcache.Client, devMode bool, otpTTL int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req StudentLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		var email string
		// Defensive check
		_ = db.QueryRow(`SELECT email FROM teachers WHERE teacher_id=$1`, req.StudentID).Err()

		if err := db.QueryRow(`SELECT email FROM students WHERE student_id=$1`, req.StudentID).Scan(&email); err != nil {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}

		otp := services.GenerateOTP()
		key := "student:" + req.StudentID

		if err := services.StoreOTP(mc, key, otp, otpTTL); err != nil {
			http.Error(w, "failed to store otp", http.StatusInternalServerError)
			return
		}

		services.SendOTP(devMode, email, otp)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "OTP sent (or logged in dev mode)",
			"uid":     req.StudentID,
		})
	}
}

// StudentVerifyOTP - POST /student/otp/verify
func StudentVerifyOTP(db *sql.DB, mc *memcache.Client, jwtSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.VerifyOTPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		key := "student:" + req.UID
		if !services.VerifyOTP(mc, key, req.OTP) {
			http.Error(w, "invalid or expired otp", http.StatusUnauthorized)
			return
		}

		var name, email string
		if err := db.QueryRow(`SELECT full_name, email FROM students WHERE student_id=$1`, req.UID).Scan(&name, &email); err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		token, err := services.CreateToken(jwtSecret, req.UID, "student", name, email, 24*time.Hour)
		if err != nil {
			http.Error(w, "failed to create token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "OTP verified",
			"token":   token,
			"user": map[string]string{
				"uid":   req.UID,
				"role":  "student",
				"name":  name,
				"email": email,
			},
		})
	}
}
