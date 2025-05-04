package performance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/yourproject/auth"
	"golang.org/x/time/rate"
)

func TestLoginLoad(t *testing.T) {
	// Setup test server
	router := auth.SetupRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Rate limiter to control request rate
	limiter := rate.NewLimiter(rate.Every(time.Millisecond*10), 1)
	client := &http.Client{Timeout: 5 * time.Second}

	// Test parameters
	concurrentUsers := 100
	successCount := 0
	errorCount := 0

	var wg sync.WaitGroup
	wg.Add(concurrentUsers)

	startTime := time.Now()

	for i := 0; i < concurrentUsers; i++ {
		go func(id int) {
			defer wg.Done()
			
			// Wait for rate limiter
			limiter.Wait(context.Background())
			
			email := fmt.Sprintf("loaduser%d@example.com", id)
			password := fmt.Sprintf("password%d!", id)
			
			// Register user first
			registerPayload := map[string]string{
				"email":    email,
				"password": password,
				"name":     fmt.Sprintf("Load User %d", id),
			}
			body, _ := json.Marshal(registerPayload)
			_, _ = client.Post(testServer.URL+"/register", "application/json", bytes.NewBuffer(body))
			
			// Now test login
			loginPayload := map[string]string{
				"email":    email,
				"password": password,
			}
			body, _ = json.Marshal(loginPayload)
			
			resp, err := client.Post(testServer.URL+"/login", "application/json", bytes.NewBuffer(body))
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					successCount++
				} else {
					errorCount++
				}
				resp.Body.Close()
			} else {
				errorCount++
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	t.Logf("Completed %d login attempts in %v", concurrentUsers, duration)
	t.Logf("Success: %d, Errors: %d", successCount, errorCount)
	
	assert.True(t, float64(successCount)/float64(concurrentUsers) > 0.95, 
		"Success rate should be >95% under load")
	assert.True(t, duration < 5*time.Second, 
		"All requests should complete within 5 seconds")
}