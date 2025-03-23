package domain

import (
	"time"

	"github.com/google/uuid"
)

type OTPPurpose string

const (
	OTPPurposeEmailVerification OTPPurpose = "email_verification"
	OTPPurposePasswordReset     OTPPurpose = "password_reset"
	OTPPurposePhoneVerification OTPPurpose = "phone_verification"
	OTPPurposeAccountRecovery   OTPPurpose = "account_recovery"
	OTPPurposeTwoFactorAuth     OTPPurpose = "two_factor_auth"
	OTPPurposeLoginConfirmation OTPPurpose = "login_confirmation"
)

type OTPVerification struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"user_id"`
	OTPCode       string      `json:"-"`
	HashedOTP     string      `json:"-"`
	Purpose       OTPPurpose  `json:"purpose"`
	ContactMethod string      `json:"contact_method"`
	AttemptsMade  int         `json:"attempts_made"`
	MaxAttempts   int         `json:"max_attempts"`
	IsVerified    bool        `json:"is_verified"`
	CreatedAt     time.Time   `json:"created_at"`
	ExpiresAt     time.Time   `json:"expires_at"`
	VerifiedAt    *time.Time  `json:"verified_at,omitempty"`
	IPAddress     string      `json:"ip_address"`
	UserAgent     string      `json:"user_agent"`
	DeviceID      *uuid.UUID  `json:"device_id,omitempty"`
}
