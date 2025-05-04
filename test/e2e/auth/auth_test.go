package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/yourproject/auth" // replace with your actual auth package path
	"github.com/yourproject/models" // replace with your actual models package path
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

var (
	testServer *httptest.Server
	client     *http.Client
)

func TestMain(m *testing.M) {
	// Setup
	router := auth.SetupRouter() // Assuming you have a router setup
	testServer = httptest.NewServer(router)
	client = &http.Client{
		Timeout: 10 * time.Second,
	}

	// Run tests
	code := m.Run()

	// Teardown
	testServer.Close()
	
	os.Exit(code)
}

// Helper function to register a test user
func registerUser(t *testing.T, email, password string) {
	payload := map[string]string{
		"email":    email,
		"password": password,
		"name":     "Test User",
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	resp, err := client.Post(testServer.URL+"/register", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUserRegistration(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful registration",
			payload: map[string]string{
				"email":    "test@example.com",
				"password": "securePassword123!",
				"name":     "Test User",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid email",
			payload: map[string]string{
				"email":    "notanemail",
				"password": "password",
				"name":     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid email",
		},
		{
			name: "Weak password",
			payload: map[string]string{
				"email":    "weakpass@example.com",
				"password": "123",
				"name":     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "password too weak",
		},
		{
			name: "Duplicate email",
			payload: map[string]string{
				"email":    "duplicate@example.com",
				"password": "password123!",
				"name":     "Test User",
			},
			expectedStatus: http.StatusCreated,
		},
	}

	// Run the duplicate email test separately
	t.Run(tests[3].name, func(t *testing.T) {
		body, err := json.Marshal(tests[3].payload)
		require.NoError(t, err)

		// First registration should succeed
		resp, err := client.Post(testServer.URL+"/register", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Second registration should fail
		resp, err = client.Post(testServer.URL+"/register", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	// Run other tests
	for _, tt := range tests[:3] {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			resp, err := client.Post(testServer.URL+"/register", "application/json", bytes.NewBuffer(body))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	// First register a test user
	testEmail := "login_test@example.com"
	testPassword := "password123!"
	registerUser(t, testEmail, testPassword)

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "Successful login",
			payload: map[string]string{
				"email":    testEmail,
				"password": testPassword,
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "Wrong password",
			payload: map[string]string{
				"email":    testEmail,
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name: "Nonexistent user",
			payload: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "password",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			resp, err := client.Post(testServer.URL+"/login", "application/json", bytes.NewBuffer(body))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectToken {
				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.NotEmpty(t, response["token"], "Expected auth token in response")
				assert.NotEmpty(t, response["expires"], "Expected expiration time in response")
			}
		})
	}
}

func TestPasswordResetFlow(t *testing.T) {
	email := "reset_test@example.com"
	registerUser(t, email, "oldPassword123!")

	// Test reset request
	t.Run("Request password reset", func(t *testing.T) {
		payload := map[string]string{"email": email}
		body, _ := json.Marshal(payload)
		
		resp, err := client.Post(testServer.URL+"/password/reset/request", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "Reset email sent", response["message"])
	})

	// Test reset with invalid token
	t.Run("Reset with invalid token", func(t *testing.T) {
		payload := map[string]string{
			"token":    "invalid-token",
			"password": "newPassword123!",
		}
		body, _ := json.Marshal(payload)
		
		resp, err := client.Post(testServer.URL+"/password/reset/confirm", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid token")
	})

	// Note: Testing with a valid token would require either:
	// 1. Mocking the email service to capture the token
	// 2. Adding a test endpoint to generate a valid token
	// This is implementation-specific
}

func TestAuthenticationEdgeCases(t *testing.T) {
	t.Run("Concurrent login attempts", func(t *testing.T) {
		email := "concurrent_test@example.com"
		registerUser(t, email, "password123!")
		
		var wg sync.WaitGroup
		attempts := 5
		successCount := 0
		
		for i := 0; i < attempts; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				payload := map[string]string{
					"email":    email,
					"password": "password123!",
				}
				body, _ := json.Marshal(payload)
				
				resp, err := client.Post(testServer.URL+"/login", "application/json", bytes.NewBuffer(body))
				if err == nil && resp.StatusCode == http.StatusOK {
					successCount++
					resp.Body.Close()
				}
			}()
		}
		
		wg.Wait()
		assert.Equal(t, 1, successCount, "Only one login should succeed")
	})
	
	t.Run("Brute force protection", func(t *testing.T) {
		email := "bruteforce_test@example.com"
		registerUser(t, email, "strongPassword123!")
		
		// Make multiple failed attempts
		for i := 0; i < 6; i++ {
			payload := map[string]string{
				"email":    email,
				"password": "wrongpassword",
			}
			body, _ := json.Marshal(payload)
			
			resp, err := client.Post(testServer.URL+"/login", "application/json", bytes.NewBuffer(body))
			require.NoError(t, err)
			resp.Body.Close()
			
			if i < 5 {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			} else {
				// After 5 attempts, should be locked
				assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
			}
		}
		
		// Verify legitimate login still works after cooldown
		t.Run("Legitimate login after cooldown", func(t *testing.T) {
			// In a real test, you'd wait for the cooldown period
			// For testing, we'll just verify the account isn't permanently locked
			payload := map[string]string{
				"email":    email,
				"password": "strongPassword123!",
			}
			body, _ := json.Marshal(payload)
			
			resp, err := client.Post(testServer.URL+"/login", "application/json", bytes.NewBuffer(body))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("Malformed JSON input", func(t *testing.T) {
		resp, err := client.Post(testServer.URL+"/login", "application/json", bytes.NewBuffer([]byte("{invalid json")))
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid JSON")
	})
}