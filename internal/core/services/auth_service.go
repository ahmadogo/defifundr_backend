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
	// // Get session by refresh token
	// session, err := a.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	// }

	// // Check if session exists
	// if session == nil {
	// 	return nil, errors.New("invalid refresh token")
	// }

	// // Check if session is blocked
	// if session.IsBlocked {
	// 	return nil, errors.New("session is blocked")
	// }

	// // Check if session has expired
	// if time.Now().After(session.ExpiresAt) {
	// 	return nil, errors.New("refresh token has expired")
	// }

	// // Get the user
	// user, err := a.userRepo.GetUserByID(ctx, session.UserID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get user: %w", err)
	// }

	// // Generate a new access token
	// accessToken, _, err := a.tokenMaker.CreateToken(
	// 	user.Email,
	// 	user.ID.String(),
	// 	a.config.AccessTokenDuration,
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create access token: %w", err)
	// }

	// // Optional: Refresh the session with a new refresh token
	// // For security reasons, many implementations rotate refresh tokens
	// if a.config.RotateRefreshTokens {
	// 	// Delete the old session
	// 	err = a.sessionRepo.DeleteSession(ctx, session.ID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to delete old session: %w", err)
	// 	}

	// 	// Create a new session

	// 	newSessionID := domain.Session{
	// 		ID:           uuid.New(),
	// 		UserID:       user.ID,
	// 		RefreshToken: random.RandomString(32),
	// 		UserAgent:    userAgent,
	// 		ClientIP:     clientIP,
	// 		IsBlocked:    false,
	// 		ExpiresAt:    time.Now().Add(a.config.RefreshTokenDuration),
	// 		CreatedAt:    time.Now(),
	// 	}

	// 	session, err = a.sessionRepo.CreateSession(ctx, newSessionID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to create new session: %w", err)
	// 	}
	// }

	// return session, accessToken, nil
	panic("unimplemented")
}

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
		// Now claims is a *Web3AuthClaims struct, not a map
		if claims.Email != "" {
			user.Email = claims.Email
		}

		if claims.Name != "" {
			nameParts := strings.Split(claims.Name, " ")
			user.FirstName = nameParts[0]
			if len(nameParts) > 1 {
				user.LastName = strings.Join(nameParts[1:], " ")
			}
		}

		if claims.ProfileImage != "" {
			profileImage := claims.ProfileImage
			user.ProfilePicture = &profileImage
		}

		// Set provider ID (usually the email for Google OAuth)
		if claims.VerifierID != "" {
			user.ProviderID = claims.VerifierID
		}

		// Refine provider information based on verifier
		if claims.Verifier != "" {
			if strings.Contains(claims.Verifier, "google") {
				user.AuthProvider = "google"
			} else if strings.Contains(claims.Verifier, "facebook") {
				user.AuthProvider = "facebook"
			} else if strings.Contains(claims.Verifier, "twitter") {
				user.AuthProvider = "twitter"
			}
			// Add more provider mappings as needed
		}

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

// RegisterBusiness implements ports.AuthService.
func (a *authService) RegisterBusiness(ctx context.Context, user domain.User) (*domain.User, error) {
	// Add Users business details
	// Update the user with business details
	panic("unimplemented")
}

// RegisterPersonalDetails implements ports.AuthService
func (a *authService) RegisterPersonalDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	a.logger.Info("Starting user personal details update process", map[string]interface{}{
		"user_id": user.ID,
	})

	// Get the existing user by ID
	existingUser, err := a.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		a.logger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Update only the personal details fields, keeping other fields as they are
	updatedUser := *existingUser
	updatedUser.Nationality = user.Nationality

	if user.AccountType != "" {
		updatedUser.AccountType = user.AccountType
	}

	if user.PersonalAccountType != "" {
		updatedUser.PersonalAccountType = user.PersonalAccountType
	}

	if user.PhoneNumber != nil {
		updatedUser.PhoneNumber = user.PhoneNumber
	}

	// Update the user with personal details
	result, err := a.userRepo.UpdateUserPersonalDetails(ctx, updatedUser)
	if err != nil {
		a.logger.Error("Failed to update user personal details", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to update user personal details: %w", err)
	}

	a.logger.Info("User personal details updated successfully", map[string]interface{}{
		"user_id": user.ID,
	})

	return result, nil
}

// RegisterBusinessDetails implements ports.AuthService
func (a *authService) RegisterBusinessDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	a.logger.Info("Starting business details update process", map[string]interface{}{
		"user_id": user.ID,
	})

	// Get the existing user by ID
	existingUser, err := a.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		a.logger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Update only the business details fields, keeping other fields as they are
	updatedUser := *existingUser

	if user.CompanyName != "" {
		updatedUser.CompanyName = user.CompanyName
	}

	if user.CompanyAddress != "" {
		updatedUser.CompanyAddress = user.CompanyAddress
	}

	if user.CompanyCity != "" {
		updatedUser.CompanyCity = user.CompanyCity
	}

	if user.CompanyPostalCode != "" {
		updatedUser.CompanyPostalCode = user.CompanyPostalCode
	}

	if user.CompanyCountry != "" {
		updatedUser.CompanyCountry = user.CompanyCountry
	}

	if user.CompanyWebsite != nil {
		updatedUser.CompanyWebsite = user.CompanyWebsite
	}

	if user.EmploymentType != nil {
		updatedUser.EmploymentType = user.EmploymentType
	}

	// Update the user with business details
	result, err := a.userRepo.UpdateUserBusinessDetails(ctx, updatedUser)
	if err != nil {
		a.logger.Error("Failed to update business details", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to update business details: %w", err)
	}

	a.logger.Info("Business details updated successfully", map[string]interface{}{
		"user_id": user.ID,
	})

	return result, nil
}

// RegisterAddressDetails implements ports.AuthService
func (a *authService) RegisterAddressDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	a.logger.Info("Starting address details update process", map[string]interface{}{
		"user_id": user.ID,
	})

	// Get the existing user by ID
	existingUser, err := a.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		a.logger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Update only the address details fields, keeping other fields as they are
	updatedUser := *existingUser

	if user.UserAddress != nil {
		updatedUser.UserAddress = user.UserAddress
	}

	if user.City != "" {
		updatedUser.City = user.City
	}

	if user.PostalCode != "" {
		updatedUser.PostalCode = user.PostalCode
	}

	// Update the user with address details
	result, err := a.userRepo.UpdateUserAddressDetails(ctx, updatedUser)
	if err != nil {
		a.logger.Error("Failed to update address details", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to update address details: %w", err)
	}

	a.logger.Info("Address details updated successfully", map[string]interface{}{
		"user_id": user.ID,
	})

	return result, nil
}

// GetUserByEmail implements ports.AuthService.
func (a *authService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByID implements ports.AuthService.
func (a *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (a *authService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := a.userRepo.CheckEmailExists(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}

	return exists, nil
}

// CreateSession creates a new session for the user
func (a *authService) CreateSession(ctx context.Context, userID uuid.UUID, userAgent, clientIP string, webOAuthClientID string, email string, login_type string) (*domain.Session, error) {
	a.logger.Info("Creating new session", map[string]interface{}{
		"user_id": userID,
		"ip":      clientIP,
	})

	// Generate a new access token
	accessToken, payload, err := a.tokenMaker.CreateToken(
		email,
		userID,
		a.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate Session ID
	session := domain.Session{
		ID:               uuid.New(),
		UserID:           userID,
		OAuthAccessToken: accessToken,
		UserAgent:        userAgent,
		ClientIP:         clientIP,
		IsBlocked:        false,
		ExpiresAt:        time.Now().Add(a.config.AccessTokenDuration),
		CreatedAt:        time.Now(),
		UserLoginType:    login_type,
	}

	a.logger.Info("Session created", map[string]interface{}{
		"session_id": session.ID,
		"token":      accessToken,
	})

	// Set OAuth fields with non-nil values
	oauthAccessToken := accessToken
	session.OAuthAccessToken = oauthAccessToken

	if webOAuthClientID != "" {
		session.WebOAuthClientID = &webOAuthClientID
	}

	userSession, err := a.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	a.logger.Info("Session created successfully", map[string]interface{}{
		"session_id": session.ID,
		"expires_at": payload.ExpiredAt,
	})

	return userSession, nil
}
