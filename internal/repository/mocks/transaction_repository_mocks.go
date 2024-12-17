package mocks

import (
	// "mobile-banking/internal/models"

	"github.com/stretchr/testify/mock"
)

// TransactionRepository adalah mock untuk repository Transaction
type TransactionRepository struct {
	mock.Mock
}

// // Create mocking method untuk membuat transaksi baru
// func (m *TransactionRepository) Create(transaction *models.Transaction) error {
// 	args := m.Called(transaction)
// 	return args.Error(0)
// }

// // FindByUserID mocking method untuk mencari transaksi berdasarkan ID user
// func (m *TransactionRepository) FindByUserID(userID uint, limit, offset int) ([]models.Transaction, error) {
// 	args := m.Called(userID, limit, offset)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).([]models.Transaction), args.Error(1)
// }

// // FindByAccountID mocking method untuk mencari transaksi berdasarkan ID akun
// func (m *TransactionRepository) FindByAccountID(accountID uint, limit, offset int) ([]models.Transaction, error) {
// 	args := m.Called(accountID, limit, offset)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).([]models.Transaction), args.Error(1)
// }

// // GetTotalTransactionAmount mocking method untuk mendapatkan total jumlah transaksi
// func (m *TransactionRepository) GetTotalTransactionAmount(userID uint, transactionType models.TransactionType) (float64, error) {
// 	args := m.Called(userID, transactionType)
// 	return args.Get(0).(float64), args.Error(1)
// }