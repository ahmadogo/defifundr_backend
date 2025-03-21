// internal/utils/time.go
package utils

import (
	"sync"
	"time"
)

// For testing purposes, we can mock the current time
var (
	mockTime     *time.Time
	mockTimeLock sync.RWMutex
)

// GetCurrentTime returns the current time.
// If a mock time is set, it returns that instead.
// This allows for easier testing of time-dependent code.
func GetCurrentTime() time.Time {
	mockTimeLock.RLock()
	defer mockTimeLock.RUnlock()

	if mockTime != nil {
		return *mockTime
	}
	return time.Now().UTC()
}

// SetMockTime sets a mock time for testing.
// Pass nil to reset to normal time behavior.
func SetMockTime(t *time.Time) {
	mockTimeLock.Lock()
	defer mockTimeLock.Unlock()
	mockTime = t
}

// FormatTimeRFC3339 formats a time using RFC3339 format.
func FormatTimeRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTimeRFC3339 parses a string in RFC3339 format to time.Time.
func ParseTimeRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// AddDuration adds a duration to a time.
func AddDuration(t time.Time, d time.Duration) time.Time {
	return t.Add(d)
}

// OTPExpirationTime returns a time set to the standard OTP expiration duration from now.
// Typically used when generating new OTPs.
func OTPExpirationTime() time.Time {
	// 15 minutes is a common OTP expiration duration
	return GetCurrentTime().Add(15 * time.Minute)
}

// TokenExpirationTime returns a time set to the standard JWT token expiration from now.
func TokenExpirationTime() time.Time {
	// 1 hour is a common JWT token expiration duration
	return GetCurrentTime().Add(1 * time.Hour)
}

// RefreshTokenExpirationTime returns a time set to the standard refresh token expiration from now.
func RefreshTokenExpirationTime() time.Time {
	// 7 days is a common refresh token expiration duration
	return GetCurrentTime().Add(7 * 24 * time.Hour)
}
