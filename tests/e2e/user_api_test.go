package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://app:8080/api/v1"

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// generateUsername creates a unique username with a timestamp and random suffix
func generateUsername() string {
	rand.Seed(time.Now().UnixNano())
	timestamp := time.Now().Format("20060102150405")
	randomSuffix := rand.Intn(1000)
	return fmt.Sprintf("e2e_test_user_%s_%d", timestamp, randomSuffix)
}

// generateEmail creates a unique email using the generated username
func generateEmail(username string) string {
	return fmt.Sprintf("%s@test.com", username)
}

func TestUserRegistrationAndLogin(t *testing.T) {
	// Generate dynamic username and email
	dynamicUsername := generateUsername()
	dynamicEmail := generateEmail(dynamicUsername)

	// Registration test
	registerPayload := RegisterRequest{
		Username: dynamicUsername,
		Email:    dynamicEmail,
		Password: "testpassword123",
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(registerPayload)
	assert.NoError(t, err)

	// Send registration request
	resp, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check registration response
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Login test
	loginPayload := LoginRequest{
		Username: dynamicUsername,
		Password: "testpassword123",
	}

	// Convert login payload to JSON
	loginJSON, err := json.Marshal(loginPayload)
	assert.NoError(t, err)

	// Send login request
	loginResp, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewBuffer(loginJSON))
	assert.NoError(t, err)
	defer loginResp.Body.Close()

	// Check login response
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	// Decode login response
	var loginResponse map[string]string
	err = json.NewDecoder(loginResp.Body).Decode(&loginResponse)
	assert.NoError(t, err)

	// Verify token is returned
	assert.Contains(t, loginResponse, "token")
	assert.NotEmpty(t, loginResponse["token"])
}