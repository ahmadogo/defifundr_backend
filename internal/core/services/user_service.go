package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	utils "github.com/demola234/defifundr/pkg/hash"
	jwtv5 "github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

type userService struct {
	userRepo ports.UserRepository
}

// UpdateKYC implements ports.UserService.
func (u *userService) UpdateKYC(ctx context.Context, kyc domain.KYC) error {
	panic("unimplemented")
}

// NewUserService creates a new instance of userService.
func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUserByID implements ports.UserService.
func (u *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with ID %s: %w", userID, err)
	}

	// Remove sensitive information
	user.Password = nil

	return user, nil
}

// UpdateUser implements ports.UserService.
func (u *userService) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	// Verify user exists
	existingUser, err := u.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with ID %s: %w", user.ID, err)
	}

	// Preserve fields that shouldn't be updated through this method
	user.Password = existingUser.Password
	user.CreatedAt = existingUser.CreatedAt
	user.UpdatedAt = time.Now()

	// Update the user
	updatedUser, err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Remove sensitive information
	updatedUser.Password = nil

	return updatedUser, nil
}

// UpdatePassword implements ports.UserService.
func (u *userService) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Get the user
	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user with ID %s: %w", userID, err)
	}

	// Verify old password
	if err := utils.CheckPassword(oldPassword, *user.Password); err != nil {
		return errors.New("incorrect old password")
	}

	// Validate new password (length, complexity, etc.)
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update the password
	err = u.userRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Helper function to validate password
func validatePassword(password string) error {
	// Implement password validation logic
	// For example: minimum length, required characters, etc.
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Add more validation as needed

	return nil
}

func verifyWeb3AuthToken(ctx context.Context, tokenString string) (jwtv5.MapClaims, error) {
	// Load JWKS from Web3Auth
	jwksURL := "https://authjs.web3auth.io/jwks"
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %v", err)
	}

	// Manually patch: wrap Keyfunc
	keyFunc := func(t *jwtv5.Token) (interface{}, error) {
		// Rebuild a v4 Token to use with keyfunc
		return jwks.Keyfunc((*jwtv5.Token)(t))
	}

	// Parse the token
	token, err := jwtv5.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwtv5.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func main() {
	ctx := context.Background()

	idToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ikc2LWlQN0NfNE5kcHBGd2lFdWNrTkNVX2V0RklWQkJ4eHQ0VUg5Y0I5RzAifQ.eyJpYXQiOjE3NDU3NzU5MjcsImlzcyI6Imh0dHBzOi8vYXV0aGpzLndlYjNhdXRoLmlvL2p3a3MiLCJhdWQiOiJkZW1vLXNmYS53ZWIzYXV0aC5pbyIsIndhbGxldHMiOlt7InB1YmxpY19rZXkiOiIweDAzOTVkMTliNzQ2MDI2NmQ5YTFhMWM1YjI3Yzk2NTkxNDlmZGNiNzYwZWZkODljYTUyNTYxOGZlZWY2NDRiMGZjYiIsInR5cGUiOiJ3ZWIzYXV0aF9hcHBfa2V5IiwiY3VydmUiOiJzZWNwMjU2azEifV0sImVtYWlsIjoia29sYXdvbGVvbHV3YXNlZ3VuNTY3QGdtYWlsLmNvbSIsIm5hbWUiOiJBZGVtb2xhIEtvbGF3b2xlIiwicHJvZmlsZUltYWdlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jSk9Vc0JGbTM3eExUUllsUlBpSDF0enVlbHo0Z1ZOWWN1UGdFOVFvTEpFOV9PbHBEckw9czk2LWMiLCJ2ZXJpZmllciI6InczYS1zZmEtd2ViLWdvb2dsZSIsInZlcmlmaWVySWQiOiJrb2xhd29sZW9sdXdhc2VndW41NjdAZ21haWwuY29tIiwidHlwZU9mTG9naW4iOiJqd3QiLCJvQXV0aElkVG9rZW4iOiJleUpoYkdjaU9pSlNVekkxTmlJc0ltdHBaQ0k2SWpJelpqZGhNelU0TXpjNU5tWTVOekV5T1dVMU5ERTRaamxpTWpFek5tWmpZekJoT1RZME5qSWlMQ0owZVhBaU9pSktWMVFpZlEuZXlKcGMzTWlPaUpvZEhSd2N6b3ZMMkZqWTI5MWJuUnpMbWR2YjJkc1pTNWpiMjBpTENKaGVuQWlPaUkxTVRreU1qZzVNVEU1TXprdFkzSnBNREZvTlRWc2MycGljMmxoTVdzM2JHdzJjWEJoYkhKMWN6YzFjSE11WVhCd2N5NW5iMjluYkdWMWMyVnlZMjl1ZEdWdWRDNWpiMjBpTENKaGRXUWlPaUkxTVRreU1qZzVNVEU1TXprdFkzSnBNREZvTlRWc2MycGljMmxoTVdzM2JHdzJjWEJoYkhKMWN6YzFjSE11WVhCd2N5NW5iMjluYkdWMWMyVnlZMjl1ZEdWdWRDNWpiMjBpTENKemRXSWlPaUl4TVRneU1EZzFOelV3T1RZMU5EWXpOVFF5TlRNaUxDSmxiV0ZwYkNJNkltdHZiR0YzYjJ4bGIyeDFkMkZ6WldkMWJqVTJOMEJuYldGcGJDNWpiMjBpTENKbGJXRnBiRjkyWlhKcFptbGxaQ0k2ZEhKMVpTd2libUptSWpveE56UTFOemMxTlRJMUxDSnVZVzFsSWpvaVFXUmxiVzlzWVNCTGIyeGhkMjlzWlNJc0luQnBZM1IxY21VaU9pSm9kSFJ3Y3pvdkwyeG9NeTVuYjI5bmJHVjFjMlZ5WTI5dWRHVnVkQzVqYjIwdllTOUJRMmM0YjJOS1QxVnpRa1p0TXpkNFRGUlNXV3hTVUdsSU1YUjZkV1ZzZWpSblZrNVpZM1ZRWjBVNVVXOU1Ta1U1WDA5c2NFUnlURDF6T1RZdFl5SXNJbWRwZG1WdVgyNWhiV1VpT2lKQlpHVnRiMnhoSWl3aVptRnRhV3g1WDI1aGJXVWlPaUpMYjJ4aGQyOXNaU0lzSW1saGRDSTZNVGMwTlRjM05UZ3lOU3dpWlhod0lqb3hOelExTnpjNU5ESTFMQ0pxZEdraU9pSTNOR05oTm1FNVl6WXdOalpqTmpFMlptRTRaVFl3TVdZMU1UVmhNRFJsWkdReFlUSm1OR0V4SW4wLmNsSmdycEFyWVJPUnRGcTBFYjVvTXQ2dHVCcXQ3RGVMQlhuRU51UnZzbk5FOEpGdXlCV1M0MUZadHc1SldoWDFqT0RKM0liQ29mVzZqUHdJU3hKeTdRb0ZrbnhhbFVaUGJ3SXNkR1NlSUpUQkJkNzlDYndyZWhHWXhqRGFrUUhPazBxR0xla2U4VGl2QUhoQlBuUjJnQkpkMW1vUUtHYmplbElZODFMN0JDZ2cwSFc1QXNJU0FMWGNnNDltcTlmdEVYZVFJTzRVQk44ZkNiUVhybHpKRHVHMFZFeks5aEpLZFJDQW1vVXdzV2sxeUluZzgtR2tIV0xfWkV0VzNCcU5QOFJIRHRxaXh1bWdZMVNMY3JQWjJ3dzJKejlSc1VxcTN1Qmg1Mk1VZEZrUVF3cWdpQkZ4cHpGaVlsTFhHaXhadDNYWUZmaTJmWGZIeS1DZER2RXYydyIsImV4cCI6MTc0NTg2MjMyN30.lWBdiSGwC8UwGt5hd2DKx-YRk761sdpAYQEluFx5iN0c8WSoNYwLggyDaPwj0ETwtokKwexMjCqVj4xW7lQ-Zw"

	claims, err := verifyWeb3AuthToken(ctx, idToken)
	if err != nil {
		fmt.Println("Web3Auth token verification failed:", err)
		return
	}

	fmt.Println("Verified Web3Auth token claims:", claims)
}
