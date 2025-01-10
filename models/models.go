package models

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `json:"username" gorm:"unique"`
    Password string `json:"password"`
    Balance  float64 `json:"balance"`
}

type Transaction struct {
    gorm.Model
    FromUserID uint    `json:"from_user_id"`
    ToUserID   uint    `json:"to_user_id"`
    Amount     float64 `json:"amount"`
    Type       string  `json:"type"` // "transfer", "deposit", "withdrawal"
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type TransferRequest struct {
    ToUsername string  `json:"to_username"`
    Amount     float64 `json:"amount"`
}

type DepositRequest struct {
    Amount float64 `json:"amount" validate:"required,gt=0"`
}