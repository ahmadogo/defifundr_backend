package response

import (
	"time"

	"github.com/google/uuid"
)

// SuccessResponse is a generic success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UserResponse represents the user data in responses
type UserResponse struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Provider   string     `json:"provider"`
	ProviderID string     `json:"provider_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type SessionResponse struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"-"`
	UserAgent        string    `json:"user_agent"`
	WebOAuthClientID *string   `json:"web_oauth_client_id,omitempty"`
	OAuthAccessToken *string   `json:"oauth_access_token,omitempty"`
	OAuthIDToken     *string   `json:"oauth_id_token,omitempty"`
	UserLoginType    string    `json:"user_login_type"`
	MFAEnabled       bool      `json:"mfa_enabled"`
	ClientIP         string    `json:"client_ip"`
	IsBlocked        bool      `json:"is_blocked"`
	ExpiresAt        time.Time `json:"expires_at"`
	CreatedAt        time.Time `json:"created_at"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User          UserResponse    `json:"user"`
	AccessToken   SessionResponse `json:"access_token,omitempty"`
	RefreshToken  string          `json:"refresh_token,omitempty"`
	SessionID     uuid.UUID       `json:"session_id"`
	ExpiresAt     time.Time       `json:"expires_at"`
	TokenType     string          `json:"token_type,omitempty"`
	EmailVerified bool            `json:"email_verified,omitempty"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// OTPResponse represents an OTP verification response
type OTPResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	Purpose      string    `json:"purpose"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsVerified   bool      `json:"is_verified"`
	AttemptsMade int       `json:"attempts_made"`
	MaxAttempts  int       `json:"max_attempts"`
}

// KYCResponse represents a KYC verification response
type KYCResponse struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	VerificationStatus string    `json:"verification_status"`
	SubmittedAt        time.Time `json:"submitted_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// PageResponse represents a paginated response
type PageResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	Items      interface{} `json:"items"`
}

// ProfileCompletionResponse represents a profile completion status response
type ProfileCompletionResponse struct {
	CompletionPercentage int      `json:"completion_percentage"`
	MissingFields        []string `json:"missing_fields,omitempty"`
	RequiredActions      []string `json:"required_actions,omitempty"`
}
