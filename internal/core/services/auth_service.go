package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	commons "github.com/demola234/defifundr/infrastructure/hash"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/demola234/defifundr/pkg/random"
	tokenMaker "github.com/demola234/defifundr/pkg/token_maker"
	"github.com/google/uuid"
)

type authService struct {
	userRepo    ports.UserRepository
	sessionRepo ports.SessionRepository
	oauthRepo   ports.OAuthRepository
	tokenMaker  tokenMaker.Maker
	config      config.Config
	logger      logging.Logger
}

// NewAuthService creates a new instance of authService.
func NewAuthService(
	userRepo ports.UserRepository,
	sessionRepo ports.SessionRepository,
	oauthRepo ports.OAuthRepository,
	tokenMaker tokenMaker.Maker,
	config config.Config,
	logger logging.Logger,
) ports.AuthService {
	// Create rate limiter for API operations

	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		oauthRepo:   oauthRepo,
		tokenMaker:  tokenMaker,
		config:      config,
		logger:      logger,
	}
}

// Login implements ports.AuthService.
func (a *authService) Login(ctx context.Context, email string, password string, userAgent string, clientIP string, provider, providerId string) (*domain.Session, *domain.User, error) {
	panic("unimplemented")
}

// Logout implements ports.AuthService.
func (a *authService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return a.sessionRepo.DeleteSession(ctx, sessionID)
}

// RefreshToken implements ports.AuthService.
func (a *authService) RefreshToken(ctx context.Context, refreshToken string, userAgent string, clientIP string) (*domain.Session, string, error) {
	// Get session by refresh token
	session, err := a.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	// Check if session exists
	if session == nil {
		return nil, "", errors.New("invalid refresh token")
	}

	// Check if session is blocked
	if session.IsBlocked {
		return nil, "", errors.New("session is blocked")
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		return nil, "", errors.New("refresh token has expired")
	}

	// Get the user
	user, err := a.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	// Generate a new access token
	accessToken, _, err := a.tokenMaker.CreateToken(
		user.Email,
		user.ID.String(),
		a.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create access token: %w", err)
	}

	// Optional: Refresh the session with a new refresh token
	// For security reasons, many implementations rotate refresh tokens
	if a.config.RotateRefreshTokens {
		// Delete the old session
		err = a.sessionRepo.DeleteSession(ctx, session.ID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to delete old session: %w", err)
		}

		// Create a new session

		newSessionID := domain.Session{
			ID:           uuid.New(),
			UserID:       user.ID,
			RefreshToken: random.RandomString(32),
			UserAgent:    userAgent,
			ClientIP:     clientIP,
			IsBlocked:    false,
			ExpiresAt:    time.Now().Add(a.config.RefreshTokenDuration),
			CreatedAt:    time.Now(),
		}

		session, err = a.sessionRepo.CreateSession(ctx, newSessionID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create new session: %w", err)
		}
	}

	return session, accessToken, nil
}

// RegisterUser implements the user registration process with Web3Auth integration
// RegisterUser implements the user registration process with Web3Auth integration
func (a *authService) RegisterUser(ctx context.Context, user domain.User, passwordStr string) (*domain.User, error) {
	a.logger.Info("Starting user registration process", map[string]interface{}{
		"email":    user.Email,
		"provider": user.AuthProvider,
	})

	// Step 1: Handle authentication based on provider
	if user.AuthProvider == "email" {
		// Email-based authentication requires a password
		if passwordStr == "" {
			a.logger.Error("Password required for email authentication", nil, map[string]interface{}{
				"email": user.Email,
			})
			return nil, errors.New("password is required for email authentication")
		}

		// Hash the password
		hashedPassword, err := commons.HashPassword(passwordStr)
		if err != nil {
			a.logger.Error("Failed to hash password", err, map[string]interface{}{
				"email": user.Email,
			})
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hashedPassword
	} else if user.AuthProvider != "" && user.WebAuthToken != "" {
		// For OAuth or Web3Auth, validate the token and fill user data
		claims, err := a.oauthRepo.ValidateWebAuthToken(ctx, user.WebAuthToken)
		if err != nil {
			a.logger.Error("Failed to validate WebAuth token", err, map[string]interface{}{
				"provider": user.AuthProvider,
			})
			return nil, fmt.Errorf("invalid authentication token: %w", err)
		}

		// Extract user information from OAuth claims
		if email, ok := claims["email"].(string); ok {
			user.Email = email
		}

		if name, ok := claims["name"].(string); ok {
			user.FirstName = strings.Split(name, " ")[0]
			user.LastName = strings.Join(strings.Split(name, " ")[1:], " ")
		}

		if profileImage, ok := claims["profileImage"].(string); ok {
			user.ProfilePicture = &profileImage
		}

		// Set provider ID (usually the email for Google OAuth)
		if verifierId, ok := claims["verifierId"].(string); ok {
			user.ProviderID = verifierId
		}

		// Refine provider information based on verifier
		if verifier, ok := claims["verifier"].(string); ok {
			if strings.Contains(verifier, "google") {
				user.AuthProvider = "google"
			} else if strings.Contains(verifier, "facebook") {
				user.AuthProvider = "facebook"
			} else if strings.Contains(verifier, "twitter") {
				user.AuthProvider = "twitter"
			}
			// Add more provider mappings as needed
		}

		// // Extract wallet information if available
		// if wallets, ok := claims["wallets"].([]interface{}); ok && len(wallets) > 0 {
		// 	if walletMap, ok := wallets[0].(map[string]interface{}); ok {
		// 		if publicKey, ok := walletMap["public_key"].(string); ok {
		// 			user.WalletPublicKey = publicKey
		// 		}
		// 	}
		// }

		// For OAuth users, no password is needed
		user.PasswordHash = ""
	} else {
		// If not email auth and no token provided, it's an error
		a.logger.Error("Missing authentication credentials", nil, map[string]interface{}{
			"provider": user.AuthProvider,
		})
		return nil, errors.New("missing authentication credentials")
	}

	// Step 2: Check if user with same email already exists
	existingUser, err := a.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		a.logger.Warn("Registration attempt for existing email", map[string]interface{}{
			"email": user.Email,
		})
		return nil, errors.New("email already registered")
	}

	// Step 3: Set default values if not provided
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// Step 4: Register the user in the database
	createdUser, err := a.userRepo.CreateUser(ctx, user)
	if err != nil {
		a.logger.Error("Failed to register user", err, map[string]interface{}{
			"email": user.Email,
		})
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return createdUser, nil
}

// RegisterAddressDetails implements ports.AuthService.
func (a *authService) RegisterAddressDetails(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	panic("unimplemented")
}

// RegisterBusiness implements ports.AuthService.
func (a *authService) RegisterBusiness(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	panic("unimplemented")
}

// RegisterBusinessDetails implements ports.AuthService.
func (a *authService) RegisterBusinessDetails(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	panic("unimplemented")
}

// RegisterPersonalDetails implements ports.AuthService.
func (a *authService) RegisterPersonalDetails(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	panic("unimplemented")
}

// ResetPassword implements ports.AuthService.
func (a *authService) ResetPassword(ctx context.Context, email string, code string, newPassword string) error {
	panic("unimplemented")
}

// SendPasswordResetEmail implements ports.AuthService.
func (a *authService) SendPasswordResetEmail(ctx context.Context, email string) error {
	panic("unimplemented")
}
