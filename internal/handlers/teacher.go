package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// Local struct (specific to teachers)
type TeacherAuthRequest struct {
	TeacherID string `json:"teacher_id,omitempty"`
	UID       string `json:"uid,omitempty"`
	OTP       string `json:"otp,omitempty"`
}

type TeacherAuthResponse struct {
	Message string `json:"message"`
	UID     string `json:"uid,omitempty"`
}

// Combined Teacher Auth Handler (Login + OTP verify)
func TeacherAuthHandler(db *sql.DB, mc *memcache.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TeacherAuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Step 1: If only TeacherID is provided → Generate OTP
		if req.TeacherID != "" && req.OTP == "" {
			uid := "teacher-" + req.TeacherID + "-" + time.Now().Format("20060102150405")
			otp := "123456" // For now, static OTP (replace with random generator)

			// Save OTP in memcache (with TTL 5 min)
			err := mc.Set(&memcache.Item{
				Key:        uid,
				Value:      []byte(otp),
				Expiration: 300,
			})
			if err != nil {
				http.Error(w, "Failed to save OTP", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(TeacherAuthResponse{
				Message: "OTP sent successfully",
				UID:     uid,
			})
			return
		}

		// Step 2: If UID + OTP are provided → Verify OTP
		if req.UID != "" && req.OTP != "" {
			item, err := mc.Get(req.UID)
			if err != nil {
				http.Error(w, "OTP expired or invalid", http.StatusUnauthorized)
				return
			}

			if string(item.Value) != req.OTP {
				http.Error(w, "Invalid OTP", http.StatusUnauthorized)
				return
			}

			json.NewEncoder(w).Encode(TeacherAuthResponse{
				Message: "OTP verified successfully",
				UID:     req.UID,
			})
			return
		}

		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}
}
