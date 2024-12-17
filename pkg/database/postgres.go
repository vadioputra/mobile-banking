package database

import (
	"database/sql"
	"fmt"
	"time"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func NewConnection(conn string)(*Database, error){	
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &Database{DB: db}, err
}

func (db *Database) Migrate() error {
	// Users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Accounts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			account_number VARCHAR(20) UNIQUE NOT NULL,
			balance DECIMAL(15,2) NOT NULL DEFAULT 0,
			account_type VARCHAR(20) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating accounts table: %v", err)
	}

	// Transactions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			from_account_id INTEGER REFERENCES accounts(id),
			to_account_id INTEGER REFERENCES accounts(id),
			amount DECIMAL(15,2) NOT NULL,
			transaction_type VARCHAR(20) NOT NULL,
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating transactions table: %v", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

