package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/demola234/defifundr/internal/adapters/dto/request"
	"github.com/demola234/defifundr/internal/adapters/dto/response"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/demola234/defifundr/pkg/app_errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService ports.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param register body request.RegisterRequest true "User registration data"
// @Success 201 {object} response.UserResponse "Successfully registered"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 409 {object} response.ErrorResponse "User already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req request.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   app_errors.ErrInvalidRequest.Error(),
			Details: err.Error(),
		})
		return
	}

	// Validate request data
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   app_errors.ErrInvalidRequest.Error(),
			Details: err.Error(),
		})
		return
	}

	// Create domain User object from request
	user := domain.User{
		Email:               req.Email,
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		AccountType:         req.AccountType,
		PersonalAccountType: req.PersonalAccountType,
		Nationality:         req.Nationality,
		Gender:              &req.Gender,
		ResidentialCountry:  &req.ResidentialCountry,
		JobRole:             &req.JobRole,
		CompanyWebsite:      &req.CompanyWebsite,
		EmploymentType:      &req.EmploymentType,
	}

	// Register the user
	registeredUser, err := h.authService.RegisterUser(ctx, user, req.Password)
	if err != nil {
		errResponse := response.ErrorResponse{
			Error: app_errors.ErrInternalServer.Error(),
		}

		if app_errors.IsAppError(err) {
			appErr := err.(*app_errors.AppError)
			errResponse.Error = appErr.Error()
			
			if appErr.ErrorType == app_errors.ErrorTypeConflict {
				ctx.JSON(http.StatusConflict, errResponse)
				return
			}
			
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	// Create response DTO
	userResponse := response.UserResponse{
		ID:                  registeredUser.ID,
		Email:               registeredUser.Email,
		FirstName:           registeredUser.FirstName,
		LastName:            registeredUser.LastName,
		AccountType:         registeredUser.AccountType,
		PersonalAccountType: registeredUser.PersonalAccountType,
		Nationality:         registeredUser.Nationality,
		CreatedAt:           registeredUser.CreatedAt,
	}

	if registeredUser.Gender != nil {
		userResponse.Gender = *registeredUser.Gender
	}
	
	if registeredUser.ResidentialCountry != nil {
		userResponse.ResidentialCountry = *registeredUser.ResidentialCountry
	}
	
	if registeredUser.JobRole != nil {
		userResponse.JobRole = *registeredUser.JobRole
	}
	
	if registeredUser.CompanyWebsite != nil {
		userResponse.CompanyWebsite = *registeredUser.CompanyWebsite
	}
	
	if registeredUser.EmploymentType != nil {
		userResponse.EmploymentType = *registeredUser.EmploymentType
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse{
		Message: "User registered successfully",
		Data:    userResponse,
	})
}

// Login godoc
// @Summary Login a user
// @Description Authenticate a user and generate access token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body request.LoginRequest true "User login data"
// @Success 200 {object} response.LoginResponse "Successfully logged in"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req request.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   app_errors.ErrInvalidRequest.Error(),
			Details: err.Error(),
		})
		return
	}

	// Get user agent and client IP
	userAgent := ctx.Request.UserAgent()
	clientIP := ctx.ClientIP()

	// Authenticate user
	session, user, err := h.authService.Login(ctx, req.Email, req.Password, userAgent, clientIP)
	if err != nil {
		errResponse := response.ErrorResponse{
			Error: app_errors.ErrInternalServer.Error(),
		}

		if app_errors.IsAppError(err) {
			appErr := err.(*app_errors.AppError)
			errResponse.Error = appErr.Error()
			
			if appErr.ErrorType == app_errors.ErrorTypeUnauthorized {
				ctx.JSON(http.StatusUnauthorized, errResponse)
				return
			}
			
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	// Set refresh token as HTTP-only cookie
	ctx.SetCookie(
		"refresh_token",
		session.RefreshToken,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	// Create login response
	loginResponse := response.LoginResponse{
		User: response.UserResponse{
			ID:                  user.ID,
			Email:               user.Email,
			FirstName:           user.FirstName,
			LastName:            user.LastName,
			AccountType:         user.AccountType,
			PersonalAccountType: user.PersonalAccountType,
			Nationality:         user.Nationality,
		},
		SessionID: session.ID,
		ExpiresAt: session.ExpiresAt,
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Login successful",
		Data:    loginResponse,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.TokenResponse "New access token"
// @Failure 401 {object} response.ErrorResponse "Invalid refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Refresh token required",
		})
		return
	}

	// Get user agent and client IP
	userAgent := ctx.Request.UserAgent()
	clientIP := ctx.ClientIP()

	// Refresh token
	session, accessToken, err := h.authService.RefreshToken(ctx, refreshToken, userAgent, clientIP)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Invalid refresh token",
		})
		return
	}

	// Set new refresh token cookie if token was rotated
	ctx.SetCookie(
		"refresh_token",
		session.RefreshToken,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Token refreshed",
		Data: response.TokenResponse{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresAt:   session.ExpiresAt,
		},
	})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate user session
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Successfully logged out"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	// Get session ID from context (set by auth middleware)
	sessionID, exists := ctx.Get("session_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	// Convert session ID to UUID
	sessionUUID, ok := sessionID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Invalid session ID",
		})
		return
	}

	// Logout user
	err := h.authService.Logout(ctx, sessionUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to logout",
		})
		return
	}

	// Clear refresh token cookie
	ctx.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Logged out successfully",
	})
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Verify user email using OTP
// @Tags auth
// @Accept json
// @Produce json
// @Param verification body request.VerifyEmailRequest true "Email verification data"
// @Success 200 {object} response.SuccessResponse "Email verified successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Router /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {
	var req request.VerifyEmailRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   app_errors.ErrInvalidRequest.Error(),
			Details: err.Error(),
		})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	// Verify email
	err = h.authService.VerifyEmail(ctx, userID, req.OTPCode)
	if err != nil {
		errResponse := response.ErrorResponse{
			Error: "Failed to verify email",
		}

		if app_errors.IsAppError(err) {
			appErr := err.(*app_errors.AppError)
			errResponse.Error = appErr.Error()
		}

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Email verified successfully",
	})
}

// ResendOTP godoc
// @Summary Resend OTP
// @Description Resend OTP for verification
// @Tags auth
// @Accept json
// @Produce json
// @Param resend body request.ResendOTPRequest true "Resend OTP data"
// @Success 200 {object} response.SuccessResponse "OTP sent successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Router /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(ctx *gin.Context) {
	var req request.ResendOTPRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   app_errors.ErrInvalidRequest.Error(),
			Details: err.Error(),
		})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	// Determine OTP purpose
	var purpose domain.OTPPurpose
	switch req.Purpose {
	case "email_verification":
		purpose = domain.OTPPurposeEmailVerification
	case "password_reset":
		purpose = domain.OTPPurposePasswordReset
	default:
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid OTP purpose",
		})
		return
	}

	// Generate new OTP
	otp, err := h.authService.GenerateOTP(ctx, userID, purpose, req.ContactMethod)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: fmt.Sprintf("Failed to generate OTP: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "OTP sent successfully",
		Data: gin.H{
			"expires_at": otp.ExpiresAt,
		},
	})
}