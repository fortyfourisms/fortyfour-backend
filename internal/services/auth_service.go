package services

import (
	"errors"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	idCounter int
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(username, password, email string) (*models.User, string, error) {
	// Validate
	if username == "" || password == "" || email == "" {
		return nil, "", errors.New("all fields are required")
	}

	// Check if user exists
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return nil, "", errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	s.idCounter++
	user := &models.User{
		ID:       s.idCounter,
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	// Find user
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
