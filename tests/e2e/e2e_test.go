package e2e_test

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

const (
    baseURL = "http://app:8080"
    timeout = 10 * time.Second
)

type TestClient struct {
    client  *http.Client
    baseURL string
}

func NewTestClient() *TestClient {
    return &TestClient{
        client: &http.Client{
            Timeout: timeout,
        },
        baseURL: baseURL,
    }
}

func (tc *TestClient) makeRequest(method, path string, body interface{}, token string) (*http.Response, error) {
    var bodyReader io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        bodyReader = bytes.NewBuffer(jsonBody)
    }

    req, err := http.NewRequest(method, fmt.Sprintf("%s%s", tc.baseURL, path), bodyReader)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    if token != "" {
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
    }

    return tc.client.Do(req)
}

func TestE2EBankingFlow(t *testing.T) {
    client := NewTestClient()

    // Test 1: Register Users
    t.Run("Register Users", func(t *testing.T) {
        users := []map[string]string{
            {"username": "john_doe", "password": "pass123"},
            {"username": "jane_doe", "password": "pass456"},
        }

        for _, user := range users {
            resp, err := client.makeRequest("POST", "/register", user, "")
            assert.NoError(t, err)
            assert.Equal(t, http.StatusCreated, resp.StatusCode)
            resp.Body.Close()
        }
    })

    // Test 2: Login and Get Tokens
    var tokens []string
    t.Run("Login Users", func(t *testing.T) {
        users := []map[string]string{
            {"username": "john_doe", "password": "pass123"},
            {"username": "jane_doe", "password": "pass456"},
        }

        for _, user := range users {
            resp, err := client.makeRequest("POST", "/login", user, "")
            assert.NoError(t, err)
            assert.Equal(t, http.StatusOK, resp.StatusCode)

            var result map[string]string
            err = json.NewDecoder(resp.Body).Decode(&result)
            assert.NoError(t, err)
            assert.NotEmpty(t, result["token"])
            tokens = append(tokens, result["token"])
            resp.Body.Close()
        }
    })

    // Test 3: Check Initial Balances
    t.Run("Check Initial Balances", func(t *testing.T) {
        for _, token := range tokens {
            resp, err := client.makeRequest("GET", "/balance", nil, token)
            assert.NoError(t, err)
            assert.Equal(t, http.StatusOK, resp.StatusCode)

            var result map[string]float64
            err = json.NewDecoder(resp.Body).Decode(&result)
            assert.NoError(t, err)
            assert.Equal(t, float64(0), result["balance"])
            resp.Body.Close()
        }
    })

    // Test 4: Make Deposit
    t.Run("Make Deposit", func(t *testing.T) {
        // Deposit for john_doe
        depositData := map[string]interface{}{
            "amount": 1000.0,
        }

        resp, err := client.makeRequest("POST", "/deposit", depositData, tokens[0])
        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var result map[string]interface{}
        err = json.NewDecoder(resp.Body).Decode(&result)
        assert.NoError(t, err)
        assert.Equal(t, float64(1000), result["balance"])
        resp.Body.Close()

        // Verify balance after deposit
        resp, err = client.makeRequest("GET", "/balance", nil, tokens[0])
        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var balance map[string]float64
        err = json.NewDecoder(resp.Body).Decode(&balance)
        assert.NoError(t, err)
        assert.Equal(t, float64(1000), balance["balance"])
        resp.Body.Close()
    })

    // Test 5: Test Invalid Deposit
    t.Run("Test Invalid Deposit", func(t *testing.T) {
        // Test negative amount
        negativeDeposit := map[string]interface{}{
            "amount": -100.0,
        }

        resp, err := client.makeRequest("POST", "/deposit", negativeDeposit, tokens[0])
        assert.NoError(t, err)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
        resp.Body.Close()

        // Test zero amount
        zeroDeposit := map[string]interface{}{
            "amount": 0.0,
        }

        resp, err = client.makeRequest("POST", "/deposit", zeroDeposit, tokens[0])
        assert.NoError(t, err)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
        resp.Body.Close()
    })

    // Test 6: Make Transfer
    t.Run("Make Transfer", func(t *testing.T) {
        transferData := map[string]interface{}{
            "to_username": "jane_doe",
            "amount":     500.0,
        }

        resp, err := client.makeRequest("POST", "/transfer", transferData, tokens[0])
        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        resp.Body.Close()
    })

    // Test 7: Verify Final Balances
    t.Run("Verify Final Balances", func(t *testing.T) {
        // Check john_doe's balance (should be 500 after deposit 1000 and transfer 500)
        resp, err := client.makeRequest("GET", "/balance", nil, tokens[0])
        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var balance1 map[string]float64
        err = json.NewDecoder(resp.Body).Decode(&balance1)
        assert.NoError(t, err)
        assert.Equal(t, float64(500), balance1["balance"])
        resp.Body.Close()

        // Check jane_doe's balance (should be 500 after receiving transfer)
        resp, err = client.makeRequest("GET", "/balance", nil, tokens[1])
        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var balance2 map[string]float64
        err = json.NewDecoder(resp.Body).Decode(&balance2)
        assert.NoError(t, err)
        assert.Equal(t, float64(500), balance2["balance"])
        resp.Body.Close()
    })

    // Test 8: Test Invalid Transfer
    t.Run("Test Invalid Transfer", func(t *testing.T) {
        transferData := map[string]interface{}{
            "to_username": "non_existent_user",
            "amount":     100.0,
        }

        resp, err := client.makeRequest("POST", "/transfer", transferData, tokens[0])
        assert.NoError(t, err)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
        resp.Body.Close()
    })

    // Test 9: Test Insufficient Balance
    t.Run("Test Insufficient Balance", func(t *testing.T) {
        transferData := map[string]interface{}{
            "to_username": "jane_doe",
            "amount":     1000000.0, // Amount larger than balance
        }

        resp, err := client.makeRequest("POST", "/transfer", transferData, tokens[0])
        assert.NoError(t, err)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
        resp.Body.Close()
    })
}