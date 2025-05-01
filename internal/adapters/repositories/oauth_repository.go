package repositories

import (
	"context"
	"fmt"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"

	"github.com/MicahParks/keyfunc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

type OAuthRepository struct {
	store db.Queries
}

func NewOAuthRepository(store db.Queries) *OAuthRepository {
	return &OAuthRepository{store: store}
}



func (r *OAuthRepository) ValidateWebAuthToken(ctx context.Context, tokenString string) (*domain.Web3AuthClaims, error) {
	// Use the correct JWKS URL
	jwksURL := "https://api-auth.web3auth.io/jwks"
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
		// Add error handler for debugging
		RefreshErrorHandler: func(err error) {
			fmt.Printf("Error refreshing JWKS: %v\n", err)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %v", err)
	}

	// Configure the JWT parser to support ES256 algorithm
	parser := jwtv4.NewParser(jwtv4.WithValidMethods([]string{"ES256"}))

	// Parse the token using JWT v4 with structured claims
	claims := &domain.Web3AuthClaims{}
	token, err := parser.ParseWithClaims(tokenString, claims, jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Verify the issuer is Web3Auth
	if claims.Issuer != "https://api-auth.web3auth.io" {
		return nil, fmt.Errorf("invalid issuer: %v", claims.Issuer)
	}

	return claims, nil
}

// Update GetUserInfo to work with the new return type
func (r *OAuthRepository) GetUserInfo(ctx context.Context, token string) (*domain.User, error) {
	// Validate the token first
	claims, err := r.ValidateWebAuthToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Extract user info from claims
	email := claims.Email
	if email == "" {
		return nil, fmt.Errorf("email not found in token claims")
	}

	// Get additional user info from repository if needed
	user, err := r.GetUserInfo(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	// Combine claims with user data
	returnedUser := &domain.User{
		ID:                  user.ID,
		Email:               user.Email,
		AccountType:         user.AccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		PersonalAccountType: user.PersonalAccountType,
		Gender:              user.Gender,
		Nationality:         user.Nationality,
		ResidentialCountry:  user.ResidentialCountry,
		JobRole:             user.JobRole,
		ProfilePicture:      user.ProfilePicture,
		Address:             user.CompanyAddress,
		City:                user.CompanyCity,
		PostalCode:          user.CompanyPostalCode,
		AuthProvider:        user.AuthProvider,
		ProviderID:          user.ProviderID,
		CompanyWebsite:      user.CompanyWebsite,
		EmploymentType:      user.EmploymentType,
	}

	return returnedUser, nil
}
