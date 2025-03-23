// internal/core/ports/input_ports.go
package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/demola234/defifundr/internal/core/domain"
)

// AuthService defines the use cases for authentication
type AuthService interface {
	// User authentication
	Login(ctx context.Context, email, password, userAgent, clientIP string) (*domain.Session, *domain.User, error)
	RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error)
	VerifyEmail(ctx context.Context, userID uuid.UUID, code string) error
	
	// OTP operations
	GenerateOTP(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose, contactMethod string) (*domain.OTPVerification, error)
	VerifyOTP(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose, code string) error
	
	// Session management
	RefreshToken(ctx context.Context, refreshToken, userAgent, clientIP string) (*domain.Session, string, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
}

// UserService defines the use cases for user operations
type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	UpdateKYC(ctx context.Context, kyc domain.KYC) error
}
