package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func RegisterTestUser(t *testing.T, serverURL, email, password string) string {
	payload := map[string]string{
		"email":    email,
		"password": password,
		"name":     "Test User",
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	resp, err := http.Post(serverURL+"/register", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		return response["id"] // assuming your API returns user ID
	}
	return ""
}

func GetAuthToken(t *testing.T, serverURL, email, password string) string {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	resp, err := http.Post(serverURL+"/login", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	return response["token"]
}