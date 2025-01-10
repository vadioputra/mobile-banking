// main_test.go
package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
	"mobile-banking-v3/models"
)

// Setup test database
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatal(err)
    }
    
    // Migrate schemas
    db.AutoMigrate(&models.User{}, &models.Transaction{})
    return db
}

// Test User Registration
func TestCreateUser(t *testing.T) {
    // Setup
    testDB := setupTestDB(t)
    
    // Test data
    user := models.User{
        Username: "testuser",
        Password: "password123",
        Balance:  0,
    }
    
    // Create user
    err := testDB.Create(&user).Error
    
    // Assertions
    assert.NoError(t, err)
    assert.NotEqual(t, 0, user.ID)
    assert.Equal(t, "testuser", user.Username)
    assert.Equal(t, float64(0), user.Balance)
}

// Test Get User Balance
func TestGetBalance(t *testing.T) {
    // Setup
    testDB := setupTestDB(t)
    
    // Create test user
    user := models.User{
        Username: "balanceuser",
        Password: "pass123",
        Balance:  1000,
    }
    testDB.Create(&user)
    
    // Get user and check balance
    var foundUser models.User
    testDB.First(&foundUser, user.ID)
    
    // Assertions
    assert.Equal(t, float64(1000), foundUser.Balance)
}

// Test Transfer Money
func TestTransferMoney(t *testing.T) {
    // Setup
    testDB := setupTestDB(t)
    
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
    
    // Perform transfer
    tx := testDB.Begin()
    
    // Update balances
    transferAmount := float64(500)
    sender.Balance -= transferAmount
    recipient.Balance += transferAmount
    
    tx.Save(&sender)
    tx.Save(&recipient)
    
    // Create transaction record
    transaction := models.Transaction{
        FromUserID: sender.ID,
        ToUserID:   recipient.ID,
        Amount:     transferAmount,
        Type:       "transfer",
    }
    tx.Create(&transaction)
    
    tx.Commit()
    
    // Get updated users
    var updatedSender, updatedRecipient models.User
    testDB.First(&updatedSender, sender.ID)
    testDB.First(&updatedRecipient, recipient.ID)
    
    // Assertions
    assert.Equal(t, float64(500), updatedSender.Balance)
    assert.Equal(t, float64(500), updatedRecipient.Balance)
    
    // Check transaction record
    var savedTransaction models.Transaction
    testDB.First(&savedTransaction)
    assert.Equal(t, transferAmount, savedTransaction.Amount)
    assert.Equal(t, sender.ID, savedTransaction.FromUserID)
    assert.Equal(t, recipient.ID, savedTransaction.ToUserID)
}

// Test Insufficient Balance
func TestInsufficientBalance(t *testing.T) {
    // Setup
    testDB := setupTestDB(t)
    
    // Create sender with low balance
    sender := models.User{
        Username: "pooruser",
        Password: "pass123",
        Balance:  100,
    }
    recipient := models.User{
        Username: "recipient",
        Password: "pass123",
        Balance:  0,
    }
    testDB.Create(&sender)
    testDB.Create(&recipient)
    
    // Try to transfer more than balance
    transferAmount := float64(500)
    
    // Check if sender has enough balance
    canTransfer := sender.Balance >= transferAmount
    
    // Assertions
    assert.False(t, canTransfer)
    assert.Equal(t, float64(100), sender.Balance)  // Balance should remain unchanged
    assert.Equal(t, float64(0), recipient.Balance) // Recipient balance should remain unchanged
}