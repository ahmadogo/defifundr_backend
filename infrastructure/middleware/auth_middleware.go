package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	newError "github.com/demola234/defifundr/pkg/app_errors"
	token "github.com/demola234/defifundr/pkg/token_maker"

	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader     = "authorization"
	authorizationBearer     = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a Gin middleware for authentication
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip authentication for swagger endpoints
		if strings.HasPrefix(ctx.Request.URL.Path, "/swagger") {
			ctx.Next()
			return
		}

		authHeader := ctx.GetHeader(authorizationHeader)
		if len(authHeader) == 0 {
			err := errors.New("authorization header not found")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newError.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		stringSplit := strings.Fields(authHeader)
		if len(stringSplit) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newError.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		authType := strings.ToLower(stringSplit[0])
		if authType != authorizationBearer {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newError.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		accessToken := stringSplit[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newError.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// AuthMiddlewareWithLogger creates a Gin middleware for authentication with logging
func AuthMiddlewareWithLogger(tokenMaker token.Maker, logger logging.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		// Skip authentication for swagger and other public endpoints
		if isPublicRoute(path) {
			logger.Info("Skipping auth for public route", map[string]interface{}{
				"path": path,
			})
			ctx.Next()
			return
		}

		authHeader := ctx.GetHeader(authorizationHeader)
		logger.Info("Auth header received", map[string]interface{}{
			"header_present": len(authHeader) > 0,
			"path":           path,
		})

		if len(authHeader) == 0 {
			err := errors.New("authorization header not found")
			logger.Error("Auth failed: No authorization header", err, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newError.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		stringSplit := strings.Fields(authHeader)
		if len(stringSplit) < 2 {
			err := errors.New("invalid authorization header format")
			logger.Error("Auth failed: Invalid header format", err, map[string]interface{}{
				"header_parts": len(stringSplit),
			})
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newError.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		authType := strings.ToLower(stringSplit[0])
		if authType != authorizationBearer {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			logger.Error("Auth failed: Unsupported auth type", err, map[string]interface{}{
				"auth_type": authType,
			})
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newError.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		accessToken := stringSplit[1]
		logger.Info("Verifying token", map[string]interface{}{
			"token_fragment": accessToken[:5] + "..." + accessToken[len(accessToken)-5:],
		})

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			logger.Error("Auth failed: Token verification failed", err, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newError.ErrorResponse(err, http.StatusUnauthorized))
			return
		}

		logger.Info("Authentication successful", map[string]interface{}{
			"user_id": payload.UserID,
		})

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// isPublicRoute determines if a route should skip authentication
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/swagger",
		"/health",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/verify-email",
		"/api/v1/auth/refresh-token",
		"/api/v1/auth/forgot-password",
		"/api/v1/auth/reset-password",
		"/api/v1/auth/google",
		"/api/v1/auth/google/callback",
		"/api/v1/waitlist/join",
	}

	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}

	return false
}
