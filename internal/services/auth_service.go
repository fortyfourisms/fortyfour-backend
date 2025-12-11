// services/auth_service.go
package services

import (
	"errors"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     repository.UserRepositoryInterface
	tokenService *TokenService
}

func NewAuthService(userRepo repository.UserRepositoryInterface, tokenService *TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Register creates a new user and returns token pair
func (s *AuthService) Register(username, password, email string, idJabatan ...string) (*models.User, *models.TokenPair, error) {
	if username == "" || password == "" || email == "" {
		return nil, nil, errors.New("all fields are required")
	}

	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return nil, nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	if len(idJabatan) > 0 {
		user.IDJabatan = &idJabatan[0]
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// Generate token pair
	tokens, err := s.tokenService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// Login authenticates a user and returns token pair
func (s *AuthService) Login(username, password string) (*models.User, *models.TokenPair, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	tokens, err := s.tokenService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// Logout revokes a single refresh token
func (s *AuthService) Logout(refreshToken string) error {
	return s.tokenService.RevokeRefreshToken(refreshToken)
}
