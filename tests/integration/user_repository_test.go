package integration

import (
	"mobile-banking/internal/models"
	"mobile-banking/internal/repository"
	"mobile-banking/pkg/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Integration(t *testing.T) {
	// Setup test database connection
	db, err := database.NewConnection("postgres://user:password@postgres:5432/mobile_banking?sslmode=disable")
	assert.NoError(t, err)
	defer db.Close()

	db.Migrate()

	// Initialize repository
	userRepo := repository.NewUserRepository(db.DB)

	// Test Create User
	user := &models.User{
		Username: "integration_test_user",
		Email:    "integration@test.com",
		Password: "testpassword",
	}

	// Create user
	err = userRepo.Create(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	// Find by username
	foundUser, err := userRepo.FindByUsername(user.Username)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Username, foundUser.Username)

	// Update user
	foundUser.Email = "updated@test.com"
	err = userRepo.Update(foundUser)
	assert.NoError(t, err)

	// Verify update
	updatedUser, err := userRepo.FindByUsername(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, "updated@test.com", updatedUser.Email)
}

