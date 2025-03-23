package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// Secret key used for HMAC-based hashing (should be loaded from configuration)
var hmacSecret = []byte("your-secret-key-use-env-variable-in-production")

// HashPassword creates a bcrypt hash of a password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password matches the hashed password
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// HashData creates an HMAC SHA-256 hash of the input data
// This is suitable for hashing OTP codes and other non-password data
func HashData(data string) string {
	h := hmac.New(sha256.New, hmacSecret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SetHMACSecret allows changing the HMAC secret key (useful for testing or configuration)
func SetHMACSecret(secret []byte) {
	hmacSecret = secret
}

// GenerateRandomString generates a random string of the specified length
// Useful for creating secure tokens or temporary passwords
func GenerateRandomString(length int) (string, error) {
	// Implementation would go here
	// This could use crypto/rand to generate secure random bytes
	// and then encode them as hex or base64
	return "", nil
}