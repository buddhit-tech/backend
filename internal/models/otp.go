package models

// VerifyOTPRequest is used for OTP verification for both teachers and students
type VerifyOTPRequest struct {
	UID string `json:"uid"`
	OTP string `json:"otp"`
}
