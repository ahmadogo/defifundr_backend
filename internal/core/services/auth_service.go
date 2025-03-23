package services

import (
	"context"
	"time"

	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/demola234/defifundr/pkg/random"
	"github.com/google/uuid"
)

type authService struct {
	user_repo ports.UserRepository
	otp_repo  ports.OTPRepository
}

// NewUserUsecase creates a new instance of userUsecase.
func NewAuthService(user_repo ports.UserRepository, otp_repo ports.OTPRepository) ports.AuthService {
	return &authService{user_repo: user_repo, otp_repo: otp_repo}
}

// GenerateOTP implements ports.AuthService.
func (a *authService) GenerateOTP(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose, contactMethod string) (*domain.OTPVerification, error) {

	otpCode := random.RandomOtp()

	hashedOtp := ""
	// Generate a random 6-digit code
	otp, err := a.otp_repo.CreateOTP(ctx, domain.OTPVerification{
		ID:            uuid.New(),
		UserID:        userID,
		OTPCode:       otpCode,
		HashedOTP:     hashedOtp,
		Purpose:       purpose,
		ContactMethod: contactMethod,
		ExpiresAt:     time.Now().Add(time.Minute * 5),
		AttemptsMade:  int(1),
		IsVerified:    false,
	})

	if err != nil {
		return nil, err
	}

	// Send OTP to user
	// a.emailService.SendVerificationEmail(ctx, user.Email, user.FirstName, otp.OTPCode)

	return otp, nil
}

// Login implements ports.AuthService.
func (a *authService) Login(ctx context.Context, email string, password string, userAgent string, clientIP string) (*domain.Session, string, error) {
	panic("unimplemented")
}

// Logout implements ports.AuthService.
func (a *authService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	panic("unimplemented")
}

// RefreshToken implements ports.AuthService.
func (a *authService) RefreshToken(ctx context.Context, refreshToken string, userAgent string, clientIP string) (*domain.Session, string, error) {
	panic("unimplemented")
}

// RegisterUser implements ports.AuthService.
func (a *authService) RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	userDetails, err := a.user_repo.CreateUser(ctx, domain.User{
		ID:                  uuid.New(),
		Email:               user.Email,
		Password:            &password,
		ProfilePicture:      user.ProfilePicture,
		AccountType:         user.AccountType,
		Gender:              user.Gender,
		PersonalAccountType: user.PersonalAccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Nationality:         user.Nationality,
		ResidentialCountry:  user.ResidentialCountry,
		JobRole:             user.JobRole,
		CompanyWebsite:      user.CompanyWebsite,
		EmploymentType:      user.EmploymentType,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	})

	if err != nil {
		return nil, err
	}

	_, err = a.GenerateOTP(ctx, userDetails.ID, domain.OTPPurposeEmailVerification, userDetails.Email)
	if err != nil {
		return userDetails, nil
	}

	// Send OTP to user
	// a.emailService.SendVerificationEmail(ctx, user.Email, user.FirstName, otp.OTPCode)

	return userDetails, nil
}

// VerifyEmail implements ports.AuthService.
func (a *authService) VerifyEmail(ctx context.Context, userID uuid.UUID, code string) error {
	err := a.VerifyOTP(ctx, userID, domain.OTPPurposeEmailVerification, code)
	if err != nil {
		return err
	}

	return nil
}

// VerifyOTP implements ports.AuthService.
func (a *authService) VerifyOTP(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose, code string) error {
	err := a.otp_repo.VerifyOTP(ctx, userID, code)
	if err != nil {
		return err
	}

	return nil
}
