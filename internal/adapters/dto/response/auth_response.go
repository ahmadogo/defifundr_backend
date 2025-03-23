package response

import (
	"time"

	"github.com/google/uuid"
)

// SuccessResponse is a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is a generic error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// UserResponse represents the user data in responses
type UserResponse struct {
	ID                  uuid.UUID  `json:"id"`
	Email               string     `json:"email"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	AccountType         string     `json:"account_type"`
	PersonalAccountType string     `json:"personal_account_type,omitempty"`
	Nationality         string     `json:"nationality"`
	Gender              string     `json:"gender,omitempty"`
	ResidentialCountry  string     `json:"residential_country,omitempty"`
	JobRole             string     `json:"job_role,omitempty"`
	CompanyWebsite      string     `json:"company_website,omitempty"`
	EmploymentType      string     `json:"employment_type,omitempty"`
	ProfilePicture      string     `json:"profile_picture,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at,omitempty"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User           UserResponse `json:"user"`
	AccessToken    string       `json:"access_token,omitempty"`
	RefreshToken   string       `json:"refresh_token,omitempty"`
	SessionID      uuid.UUID    `json:"session_id"`
	ExpiresAt      time.Time    `json:"expires_at"`
	TokenType      string       `json:"token_type,omitempty"`
	EmailVerified  bool         `json:"email_verified,omitempty"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
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