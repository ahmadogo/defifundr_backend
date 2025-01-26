package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/demola234/defiraise/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	userLimiters = make(map[string]*rate.Limiter)
	limiterMu    sync.Mutex
)

// GetUserLimiter returns a rate limiter for the given user address.
func GetUserLimiter(userAddress string, requestsPerMinute int) *rate.Limiter {
	limiterMu.Lock()
	defer limiterMu.Unlock()

	if limiter, exists := userLimiters[userAddress]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), requestsPerMinute)
	userLimiters[userAddress] = limiter
	return limiter
}

func RateLimiterMiddleware(requestsPerMinute int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user address from the request (e.g., from headers, query params, or auth payload)
		userAddress := ctx.GetHeader("X-User-Address") // Example: Use a custom header
		if userAddress == "" {
			// Fallback to using the authenticated user's address
			if payload, exists := ctx.Get(authorizationPayloadKey); exists {
				userAddress = payload.(*token.Payload).Username // Assuming the payload contains the user's address
			} else {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user address is required"})
				return
			}
		}

		// Apply rate limiting
		limiter := GetUserLimiter(userAddress, requestsPerMinute)
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		ctx.Next()
	}
}
