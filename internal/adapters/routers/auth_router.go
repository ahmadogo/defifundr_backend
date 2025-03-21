package routers

import (
	"github.com/demola234/defifundr/internal/adapters/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler, authMiddleware gin.HandlerFunc) {
	authRoutes := rg.Group("/auth")

	{
		authRoutes.POST("/register", authHandler.Register)
		// authRoutes.POST("/login", authHandler.Login)
		// authRoutes.POST("/verify", authHandler.VerifyUser)
		// authRoutes.POST("/resend-otp", authHandler.ResendOtp)

		// // Protected routes (require authMiddleware)
		// authRoutes.GET("/user", authMiddleware, authHandler.GetUser)
		// authRoutes.POST("/logout", authMiddleware, authHandler.Logout)
	}
}
