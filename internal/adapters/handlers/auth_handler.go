package handlers

import (
	"github.com/demola234/defifundr/internal/adapters/dto/request"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/demola234/defifundr/pkg/app_errors"
	"github.com/gin-gonic/gin"
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
// @Success 201 {object} map[string]interface{} "Successfully registered"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Router /auth/register [post]
func (s *AuthHandler) Register(ctx *gin.Context) {
	var req *request.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		err := app_errors.ErrInvalidRequest.Error()
		ctx.JSON(400, gin.H{"error": err})
		return
	}

	// Implement the registration logic here
	user, err := s.authService.RegisterUser(ctx, domain.User{
		Email:    req.Email,
		Password: &req.Password,
	}, req.Email)
	if err != nil {
		// Handle the error
		err := app_errors.ErrInvalidRequest.Error()
		ctx.JSON(400, gin.H{"error": err})
		return

	}

	// Return the user
	ctx.JSON(200, user)
}
