package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// ====== Request payloads ======
type StudentLoginRequest struct {
	StudentID string `json:"student_id"`
}

type VerifyStudentOTPRequest struct {
	UID string `json:"uid"`
	OTP string `json:"otp"`
}

// ====== Combined Auth Handler ======
func StudentAuthHandler(db *sql.DB, mc *memcache.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Step 1: Decode request
		var loginReq StudentLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Step 2: Generate UID + OTP
		uid := generateUID()
		otp := "123456" // Static OTP for now, replace with random later

		// Store OTP in Memcache with TTL 5 mins
		if err := mc.Set(&memcache.Item{
			Key:        "student_otp_" + uid,
			Value:      []byte(otp),
			Expiration: 300,
		}); err != nil {
			http.Error(w, "failed to store otp", http.StatusInternalServerError)
			return
		}

		// Response with UID (to be used for verification)
		json.NewEncoder(w).Encode(map[string]string{
			"uid": uid,
			"msg": "OTP sent (mocked)",
		})
	}
}

func VerifyStudentOTPHandler(db *sql.DB, mc *memcache.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req VerifyStudentOTPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Step 1: Get OTP from memcache
		item, err := mc.Get("student_otp_" + req.UID)
		if err != nil {
			http.Error(w, "otp expired or not found", http.StatusUnauthorized)
			return
		}

		// Step 2: Compare
		if string(item.Value) != req.OTP {
			http.Error(w, "invalid otp", http.StatusUnauthorized)
			return
		}

		// Step 3: Create session token
		sessionToken := GenerateSessionToken()

		// Store session in memcache with TTL 1 hour
		if err := mc.Set(&memcache.Item{
			Key:        "student_session_" + req.UID,
			Value:      []byte(sessionToken),
			Expiration: int32(time.Hour.Seconds()),
		}); err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		// Response
		json.NewEncoder(w).Encode(map[string]string{
			"session_token": sessionToken,
		})
	}
}

// ===== Helpers =====
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
