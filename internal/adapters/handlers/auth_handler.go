package http

import (
	"net/http"
	"time"

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

// RegisterRoutes registers all routes for the auth handler


// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email              string  `json:"email" binding:"required,email"`
		Password           string  `json:"password" binding:"required,min=8"`
		ProfilePicture     *string `json:"profile_picture"`
		AccountType        string  `json:"account_type" binding:"required"`
		Gender             *string `json:"gender"`
		PersonalAccountType string `json:"personal_account_type" binding:"required"`
		FirstName          string  `json:"first_name" binding:"required"`
		LastName           string  `json:"last_name" binding:"required"`
		Nationality        string  `json:"nationality" binding:"required"`
		ResidentialCountry *string `json:"residential_country"`
		JobRole            *string `json:"job_role"`
		CompanyWebsite     *string `json:"company_website"`
		EmploymentType     *string `json:"employment_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := domain.User{
		Email:              req.Email,
		ProfilePicture:     req.ProfilePicture,
		AccountType:        req.AccountType,
		Gender:             req.Gender,
		PersonalAccountType: req.PersonalAccountType,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Nationality:        req.Nationality,
		ResidentialCountry: req.ResidentialCountry,
		JobRole:            req.JobRole,
		CompanyWebsite:     req.CompanyWebsite,
		EmploymentType:     req.EmploymentType,
	}

	createdUser, err := h.authService.RegisterUser(c.Request.Context(), user, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please verify your email.",
		"user": gin.H{
			"id":         createdUser.ID,
			"email":      createdUser.Email,
			"first_name": createdUser.FirstName,
			"last_name":  createdUser.LastName,
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()

	session, accessToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password, userAgent, clientIP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set refresh token as HTTP-only cookie
	c.SetCookie(
		"refresh_token",
		session.RefreshToken,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",
		true,  // Secure
		true,  // HTTP-only
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   3600, // 1 hour in seconds
		"user_id":      session.UserID,
	})
}

// VerifyEmail handles email verification
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required,uuid"`
		Code   string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.authService.VerifyEmail(c.Request.Context(), userID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()

	session, accessToken, err := h.authService.RefreshToken(c.Request.Context(), refreshToken, userAgent, clientIP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   3600, // 1 hour in seconds
		"user_id":      session.UserID,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, _ := c.Get("session_id")
	if sessionID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID required"})
		return
	}

	err := h.authService.Logout(c.Request.Context(), sessionID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clear the refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ResendVerification resends the verification email
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required,uuid"`
		Email  string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	_, err = h.authService.GenerateOTP(c.Request.Context(), userID, domain.OTPPurposeEmailVerification, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email resent"})
}

// ForgotPassword initiates the password reset process
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implementation would look up the user and generate an OTP
	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link will be sent"})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		UserID      string `json:"user_id" binding:"required,uuid"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implementation would verify the OTP and update the password
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

// AuthMiddleware checks if the user is authenticated
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation would validate the JWT token
		// This is a simplified version
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Here you would extract and validate the JWT token
		// You would also fetch the session from the database

		// For demonstration, we're just setting a dummy session ID
		c.Set("session_id", uuid.New())
		c.Set("user_id", uuid.New())
		c.Next()
	}
}