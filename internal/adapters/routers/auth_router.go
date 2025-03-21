package handlers

import "github.com/gin-gonic/gin"

// AuthHandler struct
type AuthHandler struct {
	AuthService AuthService
}

func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/verify-email", h.VerifyEmail)
		auth.POST("/refresh-token", h.RefreshToken)
		auth.POST("/logout", h.AuthMiddleware(), h.Logout)
		auth.POST("/resend-verification", h.ResendVerification)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
	}
}
