package service

import (
	"errors"
	"mobile-banking/internal/models"
	"mobile-banking/internal/repository"
	"mobile-banking/pkg/utils"
	"time"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface{
	Register(user *models.User) (*models.UserDTO, error)
	Login(username, password string) (string, error)
}

type userService struct{
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService{
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(user *models.User) (*models.UserDTO, error){
	
	existingUser, _ := s.userRepo.FindByUsername(user.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}
	fmt.Println("ini service register")
	existingEmail, _ := s.userRepo.FindByEmail(user.Email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &models.UserDTO{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
	}, nil
}

func (s *userService) Login(username, password string) (string, error){
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if !CheckPasswordHash(password, user.Password){
		return "", errors.New("invalid username or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username)

	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func HashPassword(password string) (string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}