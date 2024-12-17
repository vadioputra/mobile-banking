package unit

import (
	"testing"
	"mobile-banking/internal/models"
	"mobile-banking/internal/service"
	"mobile-banking/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Register(t *testing.T) {
	// Create mock repository
	mockUserRepo := new(mocks.UserRepository)
	userService := service.NewUserService(mockUserRepo)

	// Test data
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Setup expectations
	// Mock FindByUsername to return nil (user doesn't exist)
	mockUserRepo.On("FindByUsername", user.Username).Return(nil, nil)
	
	// Mock FindByEmail to return nil (email doesn't exist)
	mockUserRepo.On("FindByEmail", user.Email).Return(nil, nil)
	
	// Mock Create method to simulate successful user creation
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Execute
	createdUser, err := userService.Register(user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, user.Username, createdUser.Username)
	assert.Equal(t, user.Email, createdUser.Email)
	
	// Verify that expected methods were called
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_RegisterWithExistingUsername(t *testing.T) {
	// Create mock repository
	mockUserRepo := new(mocks.UserRepository)
	userService := service.NewUserService(mockUserRepo)

	// Test data
	user := &models.User{
		Username: "existinguser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Existing user found by username
	existingUser := &models.User{
		Username: user.Username,
	}

	// Setup expectations
	mockUserRepo.On("FindByUsername", user.Username).Return(existingUser, nil)

	// Execute
	createdUser, err := userService.Register(user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdUser)
	assert.EqualError(t, err, "username already exists")
	
	// Verify method calls
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_RegisterWithExistingEmail(t *testing.T) {
	// Create mock repository
	mockUserRepo := new(mocks.UserRepository)
	userService := service.NewUserService(mockUserRepo)

	// Test data
	user := &models.User{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	// Existing user found by email
	existingUser := &models.User{
		Email: user.Email,
	}

	// Setup expectations
	mockUserRepo.On("FindByUsername", user.Username).Return(nil, nil)
	mockUserRepo.On("FindByEmail", user.Email).Return(existingUser, nil)

	// Execute
	createdUser, err := userService.Register(user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdUser)
	assert.EqualError(t, err, "email already exists")
	
	// Verify method calls
	mockUserRepo.AssertExpectations(t)
}