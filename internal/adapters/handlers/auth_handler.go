package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/internal/adapters/dto/request"
	"github.com/demola234/defifundr/internal/adapters/dto/response"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService ports.AuthService
	logger      logging.Logger
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService ports.AuthService, logger logging.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// RegisterUser handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags Business-User
// @Accept json
// @Produce json
// @Param register body request.RegisterUserRequest true "User registration data"
// @Success 201 {object} response.UserResponse "Successfully registered"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 409 {object} response.ErrorResponse "User already exists"
// @Failure 429 {object} response.ErrorResponse "Too many requests"
// @Router /auth/register/user [post]
func (h *AuthHandler) RegisterUser(ctx *gin.Context) {
	// Extract request correlation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing register user request")

	var req request.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid registration request", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid registration request",
		})
		return
	}

	// Basic validation
	if req.Provider != "" && req.Provider != "email" && req.WebAuthToken == "" {
		reqLogger.Error("Missing web auth token for provider", nil, map[string]interface{}{
			"provider": req.Provider,
			"email":    req.Email,
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Web auth token is required for provider authentication",
		})
		return
	}

	// Create user domain model from request
	user := domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		AuthProvider: req.Provider,
		ProviderID:   req.ProviderID,
		WebAuthToken: req.WebAuthToken,
		// Set default values for required fields
		AccountType:         "personal", // Default account type
		PersonalAccountType: "user",     // Default personal account type
		Nationality:         "unknown",  // This will need to be updated later
	}

	// Register the user
	createdUser, err := h.authService.RegisterUser(ctx, user, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := fmt.Sprintf("Failed to register user: %w", err)

		// Map specific errors to appropriate status codes
		if strings.Contains(strings.ToLower(err.Error()), "already registered") ||
			strings.Contains(strings.ToLower(err.Error()), "already exists") {
			statusCode = http.StatusConflict
			errorMessage = "Email already registered"
		} else if strings.Contains(strings.ToLower(err.Error()), "invalid") ||
			strings.Contains(strings.ToLower(err.Error()), "required") {
			statusCode = http.StatusBadRequest
			errorMessage = err.Error()
		}

		reqLogger.Error("Failed to register user", err, map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})

		ctx.JSON(statusCode, response.ErrorResponse{
			Success: false,
			Message: errorMessage,
		})
		return
	}

	// Log successful registration
	reqLogger.Info("User registered successfully", map[string]interface{}{
		"user_id": createdUser.ID.String(),
		"email":   createdUser.Email,
	})

	// Return success response
	ctx.JSON(http.StatusCreated, response.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data: response.UserResponse{
			ID:        createdUser.ID.String(),
			Email:     createdUser.Email,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			CreatedAt: createdUser.CreatedAt,
			UpdatedAt: createdUser.UpdatedAt,
		},
	})
}

// Helper functions

// validateRegistrationRequest validates the registration request
func validateRegistrationRequest(req request.RegisterUserRequest) error {
	// Basic validation logic
	if req.Email == "" {
		return errors.New("email is required")
	}

	// If using email auth, password is required
	if (req.Provider == "" || req.Provider == "email") && req.Password == "" {
		return errors.New("password is required for email authentication")
	}

	if req.FirstName == "" {
		return errors.New("first name is required")
	}

	if req.LastName == "" {
		return errors.New("last name is required")
	}

	return nil
}

// mapRegisterRequestToUser maps registration request to user domain model
func mapRegisterRequestToUser(req request.RegisterUserRequest) domain.User {
	return domain.User{
		ID:         uuid.New(),
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		ProviderID: req.ProviderID,
	}
}

// getStringPtrValue safely extracts string value from pointer
func getStringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {
	// Extract request co-relation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing verify email request")
}

// RegisterBusiness
func (h *AuthHandler) RegisterBusiness(ctx *gin.Context) {
	// Extract request co-relation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing register business request")
}

// RegisterUserPersonalDetails
func (h *AuthHandler) RegisterUserPersonalDetails(ctx *gin.Context) {
	// Extract request co-relation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing register user personal details_request")
}

// RegisterUserAddressDetails
func (h *AuthHandler) RegisterUserAddressDetails(ctx *gin.Context) {
	// Extract request co-relation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing register user address details request")
}

// RegisterBusinessDetails
func (h *AuthHandler) RegisterBusinessDetails(ctx *gin.Context) {
	// Extract request co-relation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing register business details request")
}
