package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CreateToken generates a JWT token for a user
func CreateToken(secret []byte, uid, role, name, email string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"uid":  uid,
		"role": role,
		"name": name,
		"email": email,
		"exp":  time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
