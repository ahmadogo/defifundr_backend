package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/demola234/defiraise/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRateLimiterMiddleware(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name         string
		requests     int                                                               // Number of requests to simulate
		userAddress  string                                                            // User address to simulate
		expectedCode int                                                               // Expected HTTP status code
		setupAuth    func(t *testing.T, request *http.Request, tokenMaker token.Maker) // Optional: Setup auth for authenticated routes
	}{
		{
			name:         "WithinRateLimit",
			requests:     5, // Allow 5 requests per minute
			userAddress:  "user1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "ExceedRateLimit",
			requests:     15, // Exceed the limit of 10 requests per minute
			userAddress:  "user2",
			expectedCode: http.StatusTooManyRequests,
		},
		{
			name:         "DifferentUsers",
			requests:     5, // User3 and User4 should not affect each other's limits
			userAddress:  "user3",
			expectedCode: http.StatusOK,
		},
	}

	// Create a test server with the rate limiter middleware
	server := newTestServer(t, nil)
	rateLimit := 10 // Allow 10 requests per minute

	// Define a test endpoint
	testPath := "/rate-limit-test"
	server.router.GET(
		testPath,
		RateLimiterMiddleware(rateLimit), // Apply rate limiter middleware
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
		},
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.requests; i++ {
				// Create a new request
				recorder := httptest.NewRecorder()
				request, err := http.NewRequest(http.MethodGet, testPath, nil)
				require.NoError(t, err)

				// Set the user address in the request header
				request.Header.Set("X-User-Address", tc.userAddress)

				// Serve the request
				server.router.ServeHTTP(recorder, request)

				// Check the response for the last request
				if i == tc.requests-1 {
					require.Equal(t, tc.expectedCode, recorder.Code)
				}
			}
		})
	}
}
