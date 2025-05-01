package routers

import (
	"github.com/demola234/defifundr/internal/adapters/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler, authMiddleware gin.HandlerFunc) {
	authRoutes := rg.Group("/auth")

	{
		authRoutes.POST("/register", authHandler.RegisterUser)

		authRoutes.POST("/register/user/check-email", authHandler.CheckEmailExists)
		authRoutes.POST("/register/user/personal-details", authMiddleware, authHandler.RegisterUserPersonalDetails)
		authRoutes.POST("/register/user/address-details", authMiddleware, authHandler.RegisterUserAddressDetails)
		authRoutes.POST("/register/business/business-details", authMiddleware, authHandler.RegisterBusinessDetails)

		// authRoutes.POST("/login", authHandler.Login)

		// authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		// authRoutes.POST("/reset-password", authHandler.ResetPassword)

		// Protected routes (require authMiddleware)
		// authRoutes.POST("/logout", authMiddleware, authHandler.Logout)

	}
}
