package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
	"mobile-banking-v3/models"
)


// Database connection
type Database struct {
	*gorm.DB
}


func InitDatabase() (*Database, error) {
    var err error
    db, err := gorm.Open(sqlite.Open("bank.db"), &gorm.Config{})
    if err != nil {
        return nil, err 
    }

    if err := db.AutoMigrate(&models.User{}, &models.Transaction{}); err != nil {
        return nil, err // Tangani error AutoMigrate
    }

    return &Database{DB: db}, nil 

}