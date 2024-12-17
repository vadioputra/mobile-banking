package mocks

import (
	// "mobile-banking/internal/models"

	"github.com/stretchr/testify/mock"
)

// AccountRepository adalah mock untuk repository Account
type AccountRepository struct {
	mock.Mock
}

// // Create mocking method untuk membuat akun baru
// func (m *AccountRepository) Create(account *models.Account) error {
// 	args := m.Called(account)
// 	return args.Error(0)
// }

// FindByUserID mocking method untuk mencari akun berdasarkan ID user
// func (m *AccountRepository) FindByUserID(userID uint) ([]models.Account, error) {
// 	args := m.Called(userID)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).([]models.Account), args.Error(1)
// }

// // FindByAccountNumber mocking method untuk mencari akun berdasarkan nomor rekening
// func (m *AccountRepository) FindByAccountNumber(accountNumber string) (*models.Account, error) {
// 	args := m.Called(accountNumber)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*models.Account), args.Error(1)
// }

// UpdateBalance mocking method untuk memperbarui saldo akun
func (m *AccountRepository) UpdateBalance(accountID uint, newBalance float64) error {
	args := m.Called(accountID, newBalance)
	return args.Error(0)
}