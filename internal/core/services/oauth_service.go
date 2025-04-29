package services

import (
	"context"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/demola234/defifundr/config"

	jwtv4 "github.com/golang-jwt/jwt/v4"

	"github.com/demola234/defifundr/internal/core/ports"
)

// Ensure oauthService implements ports.OAuthService
var _ ports.OAuthService = (*OAuthServiceImpl)(nil)

// OAuthServiceImpl implements the OAuthService interface
type OAuthServiceImpl struct {
	userRepo ports.UserRepository
	config   config.Config

	// For testing
	ValidateWebAuthTokenFunc func(ctx context.Context, tokenString string) (map[string]interface{}, error)
}

// NewOAuthService creates a new instance of OAuthServiceImpl.
func NewOAuthService(
	userRepo ports.UserRepository,
	config config.Config,
) ports.OAuthService {
	return &OAuthServiceImpl{
		userRepo: userRepo,
		config:   config,
	}
}

func (o *OAuthServiceImpl) ValidateWebAuthToken(ctx context.Context, tokenString string) (map[string]interface{}, error) {
	// If using a test mock function, use it instead
	if o.ValidateWebAuthTokenFunc != nil {
		return o.ValidateWebAuthTokenFunc(ctx, tokenString)
	}

	// Standard implementation
	jwksURL := "https://authjs.web3auth.io/jwks"
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %v", err)
	}

	// First parse with JWT v4 (which is what keyfunc expects)
	v4Token, err := jwtv4.Parse(tokenString, func(token *jwtv4.Token) (interface{}, error) {
		return jwks.Keyfunc(token)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !v4Token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	v4Claims, ok := v4Token.Claims.(jwtv4.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// Convert to map[string]interface{} to return
	claims := map[string]interface{}(v4Claims)
	return claims, nil
}

func (o *OAuthServiceImpl) GetUserInfo(ctx context.Context, token string) (map[string]interface{}, error) {
	// Validate the token first
	claims, err := o.ValidateWebAuthToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Extract user info from claims
	email, _ := claims["email"].(string)
	if email == "" {
		return nil, fmt.Errorf("email not found in token claims")
	}

	// Get additional user info from repository if needed
	user, err := o.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	// Combine claims with user data
	userInfo := map[string]interface{}{
		"user_id":    user.ID,
		"email":      email,
		"name":       claims["name"],
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	return userInfo, nil
}

func (o *OAuthServiceImpl) GetUserInfoByEmail(ctx context.Context, email string) (map[string]interface{}, error) {
	// Get user from repository
	user, err := o.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	// Build user info map
	userInfo := map[string]interface{}{
		"user_id":    user.ID,
		"email":      user.Email,
		"name":       user.FirstName + " " + user.LastName,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	return userInfo, nil
}
