package routers

import (
	"github.com/demola234/defifundr/internal/adapters/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler, authMiddleware gin.HandlerFunc) {
	authRoutes := rg.Group("/auth")

	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
		authRoutes.POST("/verify-email", authHandler.VerifyEmail)
		authRoutes.POST("/resend-otp", authHandler.ResendOTP)

		// Protected routes (require authMiddleware)
		authRoutes.POST("/logout", authMiddleware, authHandler.Logout)
	
	}
}
