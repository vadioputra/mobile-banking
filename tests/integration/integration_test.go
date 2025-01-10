// integration_test.go
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
	"github.com/golang-jwt/jwt/v4"
	"mobile-banking-v3/models"
	"mobile-banking-v3/handlers"
	"mobile-banking-v3/pkg"
	
)

// Setup test app and database
func setupTestApp() (*fiber.App, *gorm.DB) {
    // Setup database
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&models.User{}, &models.Transaction{})
    
    // Create Fiber app
    app := fiber.New()
    
    // Setup routes
    app.Post("/register", func(c *fiber.Ctx) error {
		return handlers.RegisterHandler(c, db) 
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		return handlers.LoginHandler(c, db) 
	})
	
    // Protected routes
	app.Get("/balance", pkg.JwtMiddleware, func(c *fiber.Ctx) error {
		return handlers.GetBalanceHandler(c, db) 
	})
	app.Post("/transfer", pkg.JwtMiddleware, func(c *fiber.Ctx) error {
		return handlers.TransferHandler(c, db) 
	})
    
    return app, db
}

// Test Register and Login Flow
func TestRegisterAndLoginIntegration(t *testing.T) {
    app, _ := setupTestApp()
    
    // Test Register
    registerData := map[string]string{
        "username": "testuser",
        "password": "testpass",
    }
    jsonData, _ := json.Marshal(registerData)
    
    req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
    
    // Test Login
    loginReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
    loginReq.Header.Set("Content-Type", "application/json")
    loginResp, err := app.Test(loginReq)
    
    assert.NoError(t, err)
    assert.Equal(t, 200, loginResp.StatusCode)
    
    // Check if token is returned
    var result map[string]string
    json.NewDecoder(loginResp.Body).Decode(&result)
    assert.NotEmpty(t, result["token"])
}

// Test Balance Check Integration
func TestBalanceCheckIntegration(t *testing.T) {
    app, testDB := setupTestApp()
    
    // Create test user
    user := models.User{
        Username: "balanceuser",
        Password: "pass123",
        Balance:  1000,
    }
    testDB.Create(&user)
    
    // Generate token for the user
    token := generateTestToken(user.ID)
    
    // Test balance check
    req := httptest.NewRequest("GET", "/balance", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Verify balance
    var balanceResp map[string]float64
    json.NewDecoder(resp.Body).Decode(&balanceResp)
    assert.Equal(t, float64(1000), balanceResp["balance"])
}

// Test Transfer Integration
func TestTransferIntegration(t *testing.T) {
    app, testDB := setupTestApp()
    
    // Create sender and recipient
    sender := models.User{
        Username: "sender",
        Password: "pass123",
        Balance:  1000,
    }
    recipient := models.User{
        Username: "recipient",
        Password: "pass123",
        Balance:  0,
    }
    testDB.Create(&sender)
    testDB.Create(&recipient)
    
    // Generate token for sender
    token := generateTestToken(sender.ID)
    
    // Prepare transfer request
    transferData := map[string]interface{}{
        "to_username": "recipient",
        "amount":      500,
    }
    jsonData, _ := json.Marshal(transferData)
    
    // Test transfer
    req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Verify balances after transfer
    var senderUpdated, recipientUpdated models.User
    testDB.First(&senderUpdated, sender.ID)
    testDB.First(&recipientUpdated, recipient.ID)
    
    assert.Equal(t, float64(500), senderUpdated.Balance)
    assert.Equal(t, float64(500), recipientUpdated.Balance)
}

// Helper function to generate test JWT token
func generateTestToken(userID uint) string {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = userID
    t, _ := token.SignedString([]byte("your-secret-key"))
    return t
}