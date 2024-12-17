package repository

import (
	"database/sql"
	"errors"
	"mobile-banking/internal/models"
	"fmt"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository{
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error{
	query := `
		INSERT INTO users (username, email, password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}
	fmt.Println("ini repository register")
	return nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE username = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query,username).Scan(
		&user.ID,
		&user.Username, 
		&user.Email, 
		&user.Password, 
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows{
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query,email).Scan(
		&user.ID,
		&user.Username, 
		&user.Email, 
		&user.Password, 
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows{
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2 
		WHERE id = $3
	`

	_, err := r.db.Exec(query, user.Username, user.Email, user.ID)
	return err
}