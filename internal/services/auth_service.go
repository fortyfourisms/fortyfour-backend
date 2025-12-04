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
	idCounter    int
}

func NewAuthService(userRepo repository.UserRepositoryInterface, tokenService *TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) Register(username, password, email string) (*models.User, *models.TokenPair, error) {
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

	s.idCounter++
	user := &models.User{
		ID:       s.idCounter,
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
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

func (s *AuthService) Login(username, password string) (*models.User, *models.TokenPair, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Generate token pair
	tokens, err := s.tokenService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.tokenService.RevokeRefreshToken(refreshToken)
}
