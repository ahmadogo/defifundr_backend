package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/internal/adapters/dto/request"
	"github.com/demola234/defifundr/internal/adapters/dto/response"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	token "github.com/demola234/defifundr/pkg/token_maker"
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
// @Tags authentication
// @Accept json
// @Produce json
// @Param register body request.RegisterUserRequest true "User registration data"
// @Success 201 {object} response.SuccessResponse "Successfully registered"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 409 {object} response.ErrorResponse "User already exists"
// @Failure 429 {object} response.ErrorResponse "Too many requests"
// @Router /auth/register [post]
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
		errorMessage := fmt.Sprintf("Failed to register user: %s", err.Error())

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

	// If Provider is "email", then ProviderID is set to the email
	if req.Provider == "email" {
		createdUser.ProviderID = req.Email
	}

	// Create session and generate access token
	session, err := h.authService.CreateSession(ctx, createdUser.ID, ctx.Request.UserAgent(), ctx.ClientIP(), req.WebAuthToken, createdUser.Email, "registration")
	if err != nil {
		reqLogger.Error("Failed to generate access token", err, map[string]interface{}{
			"user_id": createdUser.ID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to generate access token",
		})
		return
	}

	// Log successful registration
	reqLogger.Info("User registered successfully", map[string]interface{}{
		"user_id": createdUser.ID.String(),
		"email":   createdUser.Email,
	})

	// Create the session response

	sessionResponse := response.SessionResponse{
		ID:            session.ID,
		UserID:        createdUser.ID,
		UserAgent:     session.UserAgent,
		ClientIP:      session.ClientIP,
		AccessToken:   session.OAuthAccessToken,
		UserLoginType: req.Provider,
		ExpiresAt:     session.ExpiresAt,
		CreatedAt:     session.CreatedAt,
	}

	// Return success response with user and token information
	ctx.JSON(http.StatusCreated, response.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data: response.LoginResponse{
			User: response.UserResponse{
				ID:         createdUser.ID.String(),
				Email:      createdUser.Email,
				FirstName:  createdUser.FirstName,
				LastName:   createdUser.LastName,
				Provider:   createdUser.AuthProvider,
				ProviderID: createdUser.ProviderID,
				CreatedAt:  createdUser.CreatedAt,
				UpdatedAt:  createdUser.UpdatedAt,
			},
			AccessToken: sessionResponse,
			ExpiresAt:   session.ExpiresAt,
			SessionID:   session.ID,
		},
	})
}

// RegisterUserPersonalDetails handles updating a user's personal details
// @Summary Update user personal details
// @Description Update personal details for a registered user
// @Tags authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Param personalDetails body request.RegisterPersonalDetailsRequest true "User personal details"
// @Success 200 {object} response.SuccessResponse "Successfully updated personal details"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/register/user/personal-details [post]
func (h *AuthHandler) RegisterUserPersonalDetails(ctx *gin.Context) {
	// Extract request correlation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing register user personal details request")

	// Get userID from authorization payload in context
	authPayload, exists := ctx.Get("authorization_payload")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	user := authPayload.(*token.Payload)
	// Parse request body
	var req request.RegisterPersonalDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		reqLogger.Error("Invalid personal details", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}


	// Create user domain model with updated fields
	userDetails := domain.User{
		ID:                  user.UserID,
		Nationality:         req.Nationality,
		AccountType:         req.AccountType,
		PersonalAccountType: req.PersonalAccountType,
	}

	// Add optional fields if provided
	if req.PhoneNumber != "" {
		userDetails.PhoneNumber = &req.PhoneNumber
	}

	// Update user through service
	updatedUser, err := h.authService.RegisterPersonalDetails(ctx, userDetails)
	if err != nil {
		reqLogger.Error("Failed to update personal details", err, map[string]interface{}{
			"user_id": user.UserID,
			"error":   err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to update personal details",
		})
		return
	}

	// Return success response
	reqLogger.Info("Personal details updated successfully", map[string]interface{}{
		"user_id": user.UserID,
	})
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Personal details updated successfully",
		Data: response.UserResponse{
			ID:         updatedUser.ID.String(),
			Email:      updatedUser.Email,
			FirstName:  updatedUser.FirstName,
			LastName:   updatedUser.LastName,
			Provider:   updatedUser.AuthProvider,
			ProviderID: updatedUser.ProviderID,
			CreatedAt:  updatedUser.CreatedAt,
			UpdatedAt:  updatedUser.UpdatedAt,
		},
	})
}

// RegisterUserAddressDetails handles updating a user's address details
// @Summary Update user address details
// @Description Update address details for a registered user
// @Tags authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Param addressDetails body request.RegisterAddressDetailsRequest true "User address details"
// @Success 200 {object} response.SuccessResponse "Successfully updated address details"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/register/user/address-details [post]
func (h *AuthHandler) RegisterUserAddressDetails(ctx *gin.Context) {
	// Extract request correlation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing register user address details request")

	// Get userID from authorization payload in context
	authPayload, exists := ctx.Get("authorization_payload")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	user := authPayload.(*token.Payload)

	// Parse request body
	var req request.RegisterAddressDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		reqLogger.Error("Invalid address details", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	

	// Create user domain model with the updated address details
	userDetails := domain.User{
		ID:          user.UserID,
		UserAddress: &req.AddressLine1,
		City:        req.City,
		PostalCode:  req.PostalCode,
	}

	// Update user through service
	updatedUser, err := h.authService.RegisterAddressDetails(ctx, userDetails)
	if err != nil {
		reqLogger.Error("Failed to update address details", err, map[string]interface{}{
			"user_id": user.UserID,
			"error":   err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to update address details",
		})
		return
	}

	// Return success response
	reqLogger.Info("Address details updated successfully", map[string]interface{}{
		"user_id": user.UserID,
	})
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Address details updated successfully",
		Data: response.UserResponse{
			ID:         updatedUser.ID.String(),
			Email:      updatedUser.Email,
			FirstName:  updatedUser.FirstName,
			LastName:   updatedUser.LastName,
			Provider:   updatedUser.AuthProvider,
			ProviderID: updatedUser.ProviderID,
			CreatedAt:  updatedUser.CreatedAt,
			UpdatedAt:  updatedUser.UpdatedAt,
		},
	})
}

// RegisterBusinessDetails handles updating a user's business details
// @Summary Update business details
// @Description Update business details for a registered user
// @Tags authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Param businessDetails body request.RegisterBusinessDetailsRequest true "Business details"
// @Success 200 {object} response.SuccessResponse "Successfully updated business details"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/register/business/business-details [post]
func (h *AuthHandler) RegisterBusinessDetails(ctx *gin.Context) {
	// Extract request correlation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing register business details request")

	// Get userID from authorization payload in context
	authPayload, exists := ctx.Get("authorization_payload")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	user := authPayload.(*token.Payload)

	// Parse request body
	var req request.RegisterBusinessDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Create user domain model with the updated business details
	companyWebsite := req.CompanyWebsite
	employmentType := req.CompanyType

	userDetails := domain.User{
		ID:                user.UserID,
		CompanyName:       req.CompanyName,
		CompanyAddress:    req.CompanyAddress,
		CompanyCity:       req.CompanyCity,
		CompanyPostalCode: req.CompanyPostalCode,
		CompanyCountry:    req.CompanyCountry,
	}

	// Add optional fields if provided
	if companyWebsite != "" {
		userDetails.CompanyWebsite = &companyWebsite
	}

	if employmentType != "" {
		userDetails.EmploymentType = &employmentType
	}

	// Update user through service
	updatedUser, err := h.authService.RegisterBusinessDetails(ctx, userDetails)
	if err != nil {
		reqLogger.Error("Failed to update business details", err, map[string]interface{}{
			"user_id": user.UserID,
			"error":   err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to update business details",
		})
		return
	}

	// Return success response
	reqLogger.Info("Business details updated successfully", map[string]interface{}{
		"user_id": user.UserID,
	})
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Business details updated successfully",
		Data: response.UserResponse{
			ID:         updatedUser.ID.String(),
			Email:      updatedUser.Email,
			FirstName:  updatedUser.FirstName,
			LastName:   updatedUser.LastName,
			Provider:   updatedUser.AuthProvider,
			ProviderID: updatedUser.ProviderID,
			CreatedAt:  updatedUser.CreatedAt,
			UpdatedAt:  updatedUser.UpdatedAt,
		},
	})
}

// CheckEmailExists handles checking if an email already exists in the database
// @Summary Check if email exists
// @Description Check if an email address is already registered
// @Tags authentication
// @Accept json
// @Produce json
// @Param email body request.CheckEmailRequest true "Email to check"
// @Success 200 {object} response.SuccessResponse "Email check result"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/register/user/check-email [post]
func (h *AuthHandler) CheckEmailExists(ctx *gin.Context) {
	// Extract request correlation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing check email exists request")

	var req request.CheckEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		reqLogger.Error("Invalid email", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Check if email exists
	exists, err := h.authService.CheckEmailExists(ctx, req.Email)
	if err != nil {
		reqLogger.Error("Failed to check email", err, map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to check email",
		})
		return
	}

	// Log result
	reqLogger.Info("Email check completed", map[string]interface{}{
		"email":  req.Email,
		"exists": exists,
	})

	// Return result
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Email check completed",
		Data: map[string]interface{}{
			"email":  req.Email,
			"exists": exists,
		},
	})
}
