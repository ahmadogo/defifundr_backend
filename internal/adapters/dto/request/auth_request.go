package request

// RegisterRequest represents the request body for the register endpoint
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest represents the request body for the login endpoint
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// VerifyEmailRequest represents the request body for the verify email endpoint
type VerifyEmailRequest struct {
	OTPCode string `json:"otp_code" binding:"required"`
}

// RefreshTokenRequest represents the request body for the refresh token endpoint
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
