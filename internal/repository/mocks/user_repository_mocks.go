package mocks

import (
	"mobile-banking/internal/models"

	"github.com/stretchr/testify/mock"
)

// UserRepository adalah mock untuk repository User
type UserRepository struct {
	mock.Mock
}

// Create mocking method untuk membuat user baru
func (m *UserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// FindByUsername mocking method untuk mencari user berdasarkan username
func (m *UserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// FindByEmail mocking method untuk mencari user berdasarkan email
func (m *UserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Update mocking method untuk memperbarui user
func (m *UserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Delete mocking method untuk menghapus user
func (m *UserRepository) Delete(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}