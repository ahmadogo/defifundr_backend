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

func (r *OAuthRepository) ValidateWebAuthToken(ctx context.Context, tokenString string) (map[string]interface{}, error) {

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

func (r *OAuthRepository) GetUserInfo(ctx context.Context, token string) (*domain.User, error) {
	// Validate the token first
	claims, err := r.ValidateWebAuthToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Extract user info from claims
	email, _ := claims["email"].(string)
	if email == "" {
		return nil, fmt.Errorf("email not found in token claims")
	}

	// Get additional user info from repository if needed
	user, err := r.GetUserByEmail(ctx, email)
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
		Gender:              &user.Gender.String,
		Nationality:         user.Nationality,
		ResidentialCountry:  &user.ResidentialCountry.String,
		JobRole:             &user.JobRole.String,
		ProfilePicture:      &user.ProfilePicture.String,
		Address:             user.CompanyAddress.String,
		City:                user.CompanyCity.String,
		PostalCode:          user.CompanyPostalCode.String,
		AuthProvider:        user.AuthProvider.String,
		ProviderID:          user.ProviderID.String,
		CompanyWebsite:      &user.CompanyWebsite.String,
		EmploymentType:      &user.EmploymentType.String,
	}

	return returnedUser, nil
}

func (r *OAuthRepository) GetUserByEmail(ctx context.Context, email string) (*db.Users, error) {
	user, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return &user, nil
}
