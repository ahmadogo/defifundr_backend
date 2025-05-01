package request

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// RegisterRequest represents the user registration request
type RegisterUserRequest struct {
	Email        string `json:"email" binding:"omitempty"`
	Password     string `json:"password,omitempty" binding:"omitempty,min=8"`
	FirstName    string `json:"first_name" binding:"omitempty"`
	LastName     string `json:"last_name" binding:"omitempty"`
	Provider     string `json:"provider" binding:"omitempty"`
	ProviderID   string `json:"provider_id" binding:"omitempty"`
	WebAuthToken string `json:"web_auth_token" binding:"required"`
}

// Validate validates the register request
func (r *RegisterUserRequest) Validate() error {
	// Validate email
	if !isValidEmail(r.Email) {
		return errors.New("invalid email format")
	}

	// Validate password
	if err := validatePassword(r.Password); err != nil {
		return err
	}

	// Additional validations as needed
	return nil
}

// LoginRequest represents the user login request
type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
}

// Validate validates the login request
func (r *LoginRequest) Validate() error {
	// Validate email
	if !isValidEmail(r.Email) {
		return errors.New("invalid email format")
	}

	// Validate password is not empty
	if strings.TrimSpace(r.Password) == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}

// UpdateProfileRequest represents the profile update request
type UpdateProfileRequest struct {
	FirstName          string `json:"first_name" binding:"required"`
	LastName           string `json:"last_name" binding:"required"`
	Nationality        string `json:"nationality" binding:"required"`
	Gender             string `json:"gender"`
	ResidentialCountry string `json:"residential_country"`
	JobRole            string `json:"job_role"`
	CompanyWebsite     string `json:"company_website"`
	EmploymentType     string `json:"employment_type"`
}

// Validate validates the profile update request
func (r *UpdateProfileRequest) Validate() error {
	// Validate first name and last name are not empty
	if strings.TrimSpace(r.FirstName) == "" || strings.TrimSpace(r.LastName) == "" {
		return errors.New("first name and last name cannot be empty")
	}

	// Validate nationality is not empty
	if strings.TrimSpace(r.Nationality) == "" {
		return errors.New("nationality cannot be empty")
	}

	return nil
}

// ChangePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Validate validates the change password request
func (r *ChangePasswordRequest) Validate() error {
	// Validate old password is not empty
	if strings.TrimSpace(r.OldPassword) == "" {
		return errors.New("old password cannot be empty")
	}

	// Validate new password
	if err := validatePassword(r.NewPassword); err != nil {
		return err
	}

	// Ensure old and new passwords are different
	if r.OldPassword == r.NewPassword {
		return errors.New("new password must be different from old password")
	}

	return nil
}

// VerifyEmailRequest represents the email verification request
type VerifyEmailRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	OTPCode string `json:"otp_code" binding:"required"`
}

// Validate validates the email verification request
func (r *VerifyEmailRequest) Validate() error {
	// Validate user ID is not empty
	if strings.TrimSpace(r.UserID) == "" {
		return errors.New("user ID cannot be empty")
	}

	// Validate OTP code is not empty
	if strings.TrimSpace(r.OTPCode) == "" {
		return errors.New("OTP code cannot be empty")
	}

	return nil
}

// ResendOTPRequest represents the resend OTP request
type ResendOTPRequest struct {
	UserID        string `json:"user_id" binding:"required"`
	Purpose       string `json:"purpose" binding:"required"`
	ContactMethod string `json:"contact_method" binding:"required"`
}

// Validate validates the resend OTP request
func (r *ResendOTPRequest) Validate() error {
	// Validate user ID is not empty
	if strings.TrimSpace(r.UserID) == "" {
		return errors.New("user ID cannot be empty")
	}

	// Validate purpose is valid
	if !isValidOTPPurpose(r.Purpose) {
		return errors.New("invalid OTP purpose")
	}

	// Validate contact method
	if !isValidContactMethod(r.ContactMethod) {
		return errors.New("invalid contact method")
	}

	return nil
}

// UpdateKYCRequest represents the KYC update request
type UpdateKYCRequest struct {
	IDType            string    `json:"id_type" binding:"required"`
	IDNumber          string    `json:"id_number" binding:"required"`
	IDIssuingCountry  string    `json:"id_issuing_country" binding:"required"`
	IDExpiryDate      time.Time `json:"id_expiry_date" binding:"required"`
	IDFrontImage      string    `json:"id_front_image" binding:"required"`
	IDBackImage       string    `json:"id_back_image"`
	SelfieImage       string    `json:"selfie_image" binding:"required"`
	AddressProofImage string    `json:"address_proof_image"`
}

// Validate validates the KYC update request
func (r *UpdateKYCRequest) Validate() error {
	// Validate ID type is valid
	if !isValidIDType(r.IDType) {
		return errors.New("invalid ID type")
	}

	// Validate ID number is not empty
	if strings.TrimSpace(r.IDNumber) == "" {
		return errors.New("ID number cannot be empty")
	}

	// Validate ID issuing country is not empty
	if strings.TrimSpace(r.IDIssuingCountry) == "" {
		return errors.New("ID issuing country cannot be empty")
	}

	// Validate ID expiry date is in the future
	if r.IDExpiryDate.Before(time.Now()) {
		return errors.New("ID expiry date must be in the future")
	}

	// Validate required images are provided
	if strings.TrimSpace(r.IDFrontImage) == "" {
		return errors.New("ID front image is required")
	}

	if strings.TrimSpace(r.SelfieImage) == "" {
		return errors.New("selfie image is required")
	}

	return nil
}

// Helper functions for validation

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validatePassword checks if the password meets security requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	if !hasUpperCase || !hasLowerCase || !hasNumber || !hasSpecialChar {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}

// isValidAccountType checks if the account type is valid
func isValidAccountType(accountType string) bool {
	validTypes := []string{"personal", "business"}
	for _, validType := range validTypes {
		if accountType == validType {
			return true
		}
	}
	return false
}

// isValidOTPPurpose checks if the OTP purpose is valid
func isValidOTPPurpose(purpose string) bool {
	validPurposes := []string{"email_verification", "password_reset", "two_factor_auth"}
	for _, validPurpose := range validPurposes {
		if purpose == validPurpose {
			return true
		}
	}
	return false
}

// isValidContactMethod checks if the contact method is valid
func isValidContactMethod(method string) bool {
	validMethods := []string{"email", "phone"}
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}

// isValidIDType checks if the ID type is valid
func isValidIDType(idType string) bool {
	validTypes := []string{"passport", "national_id", "drivers_license", "residence_permit"}
	for _, validType := range validTypes {
		if idType == validType {
			return true
		}
	}
	return false
}
