// internal/core/ports/input_ports.go
package ports

import (
	"context"

	"github.com/demola234/defifundr/internal/core/domain"
	emailEnums "github.com/demola234/defifundr/pkg/utils"
	"github.com/google/uuid"
)

// AuthService defines the use cases for authentication
type AuthService interface {
	// User authentication
	Login(ctx context.Context, email, password, userAgent, clientIP string) (*domain.Session, *domain.User, error)
	RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error)
	RegisterBusiness(ctx context.Context, user domain.User, password string) (*domain.User, error)
	RegisterPersonalDetails(ctx context.Context, user domain.User, password string) (*domain.User, error)
	RegisterAddressDetails(ctx context.Context, user domain.User, password string) (*domain.User, error)
	RegisterBusinessDetails(ctx context.Context, user domain.User, password string) (*domain.User, error)

	// Forgot password
	SendPasswordResetEmail(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, code, newPassword string) error

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

// WaitlistService defines the use cases for the waitlist feature
type WaitlistService interface {
	JoinWaitlist(ctx context.Context, email, fullName, referralSource string) (*domain.WaitlistEntry, error)
	GetWaitlistPosition(ctx context.Context, id uuid.UUID) (int, error)
	GetWaitlistStats(ctx context.Context) (map[string]interface{}, error)
	ListWaitlist(ctx context.Context, page, pageSize int, filters map[string]string) ([]domain.WaitlistEntry, int64, error)
	ExportWaitlist(ctx context.Context) ([]byte, error)
}

// EmailService defines methods for sending application emails
type EmailSender interface {
	SendEmail(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}) error
	SendEmailWithAttachment(ctx context.Context, recipient string, subject string, templateName string,
		data map[string]interface{}, attachments []emailEnums.EmailAttachment) error
	QueueEmail(ctx context.Context, recipient string, subject string, templateName string,
		data map[string]interface{}, priority emailEnums.EmailPriority) (string, error)
}

// EmailService defines methods for sending application emails
type EmailService interface {
	SendWaitlistConfirmation(ctx context.Context, email, name, referralCode string, position int) error
	SendWaitlistInvitation(ctx context.Context, email, name string, inviteLink string) error
	SendBatchUpdate(ctx context.Context, emails []string, subject, message string) error
}

type OAuthService interface {
	ValidateWebAuthToken(ctx context.Context, tokenString string) (map[string]interface{}, error)
	GetUserInfo(ctx context.Context, token string) (map[string]interface{}, error)
}
