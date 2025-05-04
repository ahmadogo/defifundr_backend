package mobile

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourproject/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMobileAuthFlow(t *testing.T) {
	router := auth.SetupRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	client := &http.Client{}

	// Simulate mobile headers
	t.Run("Mobile user agent", func(t *testing.T) {
		payload := map[string]string{
			"email":    "mobile_test@example.com",
			"password": "mobilePassword123!",
			"name":     "Mobile User",
		}
		body, _ := json.Marshal(payload)
		
		req, _ := http.NewRequest("POST", testServer.URL+"/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15")
		
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		
		// Check if response is optimized for mobile
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		
		assert.NotContains(t, response, "complex_layout", 
			"Response should not contain complex layout for mobile")
	})
	
	t.Run("Mobile OAuth flow", func(t *testing.T) {
		// Test mobile-optimized OAuth redirect
		req, _ := http.NewRequest("GET", testServer.URL+"/oauth/mobile", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; Mobile)")
		
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusFound, resp.StatusCode)
		redirectURL := resp.Header.Get("Location")
		assert.Contains(t, redirectURL, "mobile=true", 
			"OAuth redirect should be mobile-optimized")
	})
}