package domain

import (
	"time"

	"github.com/google/uuid"
)

type KYC struct {
	ID                   uuid.UUID `json:"id"`
	UserID               uuid.UUID `json:"user_id"`
	FaceVerification     bool      `json:"face_verification"`
	IdentityVerification bool      `json:"identity_verification"`
	VerificationType     string    `json:"verification_type"`
	VerificationNumber   string    `json:"verification_number"`
	VerificationStatus   string    `json:"verification_status"`
	UpdatedAt            time.Time `json:"updated_at"`
	CreatedAt            time.Time `json:"created_at"`
}
