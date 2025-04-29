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

type oauthService struct {
	userRepo ports.UserRepository
	config   config.Config
}

// NewOAuthService creates a new instance of oauthService.
func NewOAuthService(
	userRepo ports.UserRepository,
	config config.Config,
) ports.OAuthService {
	return &oauthService{
		userRepo: userRepo,
		config:   config,
	}
}

func (o *oauthService) ValidateWebAuthToken(ctx context.Context, tokenString string) (map[string]interface{}, error) {
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

func (o *oauthService) GetUserInfo(ctx context.Context, token string) (map[string]interface{}, error) {
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
