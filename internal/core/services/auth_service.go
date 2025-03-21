// // internal/core/services/auth_service.go
package services

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"github.com/demola234/defifundr/internal/core/domain"
// 	"github.com/demola234/defifundr/internal/core/ports"
// 	"github.com/google/uuid"
// 	"golang.org/x/crypto/bcrypt"
// )

// type authService struct {
// 	userRepo     ports.UserRepository
// 	sessionRepo  ports.SessionRepository
// 	otpRepo      ports.OTPRepository
	
// 	tokenExpiry  time.Duration
// }

// // NewAuthService creates a new authentication service with dependencies
// func NewAuthService(
// 	userRepo ports.UserRepository,
// 	sessionRepo ports.SessionRepository,
// 	otpRepo ports.OTPRepository,
// 	emailService ports.EmailService,
// 	jwtSecret []byte,
// 	tokenExpiry time.Duration,
// ) ports.AuthService {
// 	return &authService{
// 		userRepo:     userRepo,
// 		sessionRepo:  sessionRepo,
// 		otpRepo:      otpRepo,
// 		emailService: emailService,
// 		jwtSecret:    jwtSecret,
// 		tokenExpiry:  tokenExpiry,
// 	}
// }

// // Login authenticates a user and returns a session with JWT tokens
// func (s *authService) Login(ctx context.Context, email, password, userAgent, clientIP string) (*domain.Session, string, error) {
// 	// Get the user by email
// 	user, err := s.userRepo.GetUserByEmail(ctx, email)
// 	if err != nil {
// 		return nil, "", errors.New("invalid email or password")
// 	}

// 	// Validate password
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
// 		return nil, "", errors.New("invalid email or password")
// 	}

// 	// Generate JWT access token
// 	accessToken, err := s.generateJWT(user.ID)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	// Create refresh token
// 	refreshToken := s.generateRefreshToken()

// 	// Create session
// 	session := domain.Session{
// 		ID:            uuid.New(),
// 		UserID:        user.ID,
// 		RefreshToken:  refreshToken,
// 		UserAgent:     userAgent,
// 		UserLoginType: "email",
// 		MFAEnabled:    false, 
// 		ClientIP:      clientIP,
// 		IsBlocked:     false,
// 		ExpiresAt:     time.Now().Add(24 * 7 * time.Hour), // 1 week
// 		CreatedAt:     time.Now(),
// 	}

// 	createdSession, err := s.sessionRepo.CreateSession(ctx, session)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	return createdSession, accessToken, nil
// }

// // RegisterUser creates a new user and sends verification email
// func (s *authService) RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error) {
// 	// Check if email already exists
// 	existingUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
// 	if err == nil && existingUser != nil {
// 		return nil, errors.New("email already registered")
// 	}

// 	// Hash password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Set UUID and timestamps
// 	user.ID = uuid.New()
// 	user.CreatedAt = time.Now()
// 	user.UpdatedAt = time.Now()

// 	// Create user
// 	createdUser, err := s.userRepo.CreateUser(ctx, user, string(hashedPassword))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Generate and send verification OTP
// 	otp, err := s.GenerateOTP(ctx, user.ID, domain.OTPPurposeEmailVerification, user.Email)
// 	if err != nil {
// 		return createdUser, nil // Return user even if OTP generation fails
// 	}

// 	// Send verification email
// 	err = s.emailService.SendVerificationEmail(ctx, user.Email, user.FirstName, otp.OTPCode)
// 	if err != nil {
// 		return createdUser, nil // Return user even if email sending fails
// 	}

// 	return createdUser, nil
// }

// // VerifyEmail verifies a user's email using an OTP code
// func (s *authService) VerifyEmail(ctx context.Context, userID uuid.UUID, code string) error {
// 	return s.VerifyOTP(ctx, userID, domain.OTPPurposeEmailVerification, code)
// }

// // GenerateOTP creates a new OTP for the given purpose
// func (s *authService) GenerateOTP(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose, contactMethod string) (*domain.OTPVerification, error) {
// 	// Generate a random 6-digit code
// 	code := s.generateOTPCode()

// 	// Hash the OTP code
// 	hashedOTP, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create OTP record
// 	otp := domain.OTPVerification{
// 		ID:            uuid.New(),
// 		UserID:        userID,
// 		OTPCode:       code, 
// 		HashedOTP:     string(hashedOTP),
// 		Purpose:       purpose,
// 		ContactMethod: contactMethod,
// 		AttemptsMade:  0,
// 		MaxAttempts:   5,
// 		IsVerified:    false,
// 		CreatedAt:     time.Now(),
// 		ExpiresAt:     time.Now().Add(15 * time.Minute),
// 		IPAddress:     "", 
// 		UserAgent:     "", 
// 	}

// 	createdOTP, err := s.otpRepo.CreateOTP(ctx, otp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return createdOTP, nil
// }

// // VerifyOTP verifies an OTP code
// func (s *authService) VerifyOTP(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose, code string) error {
// 	// Get the most recent OTP for this user and purpose
// 	otp, err := s.otpRepo.GetOTPByUserIDAndPurpose(ctx, userID, purpose)
// 	if err != nil {
// 		return errors.New("invalid or expired verification code")
// 	}

// 	// Check if OTP is expired
// 	if time.Now().After(otp.ExpiresAt) {
// 		return errors.New("verification code has expired")
// 	}

// 	// Check if OTP is already verified
// 	if otp.IsVerified {
// 		return errors.New("code already verified")
// 	}

// 	// Check if max attempts exceeded
// 	if otp.AttemptsMade >= otp.MaxAttempts {
// 		return errors.New("maximum verification attempts exceeded")
// 	}

// 	// Verify the OTP code
// 	if err := bcrypt.CompareHashAndPassword([]byte(otp.HashedOTP), []byte(code)); err != nil {
// 		// Increment attempts
// 		_ = s.otpRepo.IncrementAttempts(ctx, otp.ID)
// 		return errors.New("invalid verification code")
// 	}

// 	// Mark as verified
// 	if err := s.otpRepo.VerifyOTP(ctx, otp.ID); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // RefreshToken refreshes an access token using a refresh token
// func (s *authService) RefreshToken(ctx context.Context, refreshToken, userAgent, clientIP string) (*domain.Session, string, error) {
// 	// Get session by refresh token
// 	session, err := s.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
// 	if err != nil {
// 		return nil, "", errors.New("invalid refresh token")
// 	}

// 	// Check if session is expired or blocked
// 	if time.Now().After(session.ExpiresAt) || session.IsBlocked {
// 		return nil, "", errors.New("session expired or invalid")
// 	}

// 	// Validate user agent and IP for security
// 	if session.UserAgent != userAgent || session.ClientIP != clientIP {
// 		// Potential token theft, block the session
// 		_ = s.sessionRepo.BlockSession(ctx, session.ID)
// 		return nil, "", errors.New("security validation failed")
// 	}

// 	// Generate new JWT
// 	accessToken, err := s.generateJWT(session.UserID)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	return session, accessToken, nil
// }

// // Logout invalidates a session
// func (s *authService) Logout(ctx context.Context, sessionID uuid.UUID) error {
// 	return s.sessionRepo.DeleteSession(ctx, sessionID)
// }

// // Helper methods
// func (s *authService) generateJWT(userID uuid.UUID) (string, error) {
// 	// Implementation using JWT library
// 	// This is a simplified placeholder
// 	return "jwt_token_" + userID.String(), nil
// }

// func (s *authService) generateRefreshToken() string {
// 	// Generate a secure random string
// 	// This is a simplified placeholder
// 	return uuid.New().String()
// }

// func (s *authService) generateOTPCode() string {
// 	// Generate a 6-digit code
// 	rand.Seed(time.Now().UnixNano())
// 	return fmt.Sprintf("%06d", rand.Intn(1000000))
// }
