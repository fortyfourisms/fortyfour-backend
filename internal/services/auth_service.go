package services

import (
	"errors"
	"strings"

	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/validator"

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
func (s *AuthService) Register(username, password, email string, roleID *string, idJabatan *string) (*models.User, *models.TokenPair, error) {
	// Trim spaces
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	email = strings.TrimSpace(email)

	// Validasi field tidak boleh kosong
	if username == "" {
		return nil, nil, errors.New("username wajib diisi")
	}
	if password == "" {
		return nil, nil, errors.New("password wajib diisi")
	}
	if email == "" {
		return nil, nil, errors.New("email wajib diisi")
	}

	// Validasi format email
	if !validator.ValidateEmail(email) {
		return nil, nil, errors.New("email tidak valid")
	}

	// Validasi format username
	if !validator.ValidateUsername(username) {
		return nil, nil, errors.New("username harus 3-50 karakter")
	}

	// Validasi panjang password
	if len(password) < 6 {
		return nil, nil, errors.New("password minimal 6 karakter")
	}

	// Check if username already exists
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return nil, nil, errors.New("username sudah digunakan")
	}

	// Check if email already exists
	exists, err := s.userRepo.EmailExists(email, nil)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, errors.New("email sudah digunakan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		Username:  username,
		Password:  string(hashedPassword),
		Email:     email,
		RoleID:    roleID,
		IDJabatan: idJabatan,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// Fetch user kembali untuk mendapatkan role_name
	user, err = s.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, nil, err
	}

	// Generate token pair with role
	tokens, err := s.tokenService.GenerateTokenPair(user.ID, user.Username, user.RoleName)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// Login authenticates a user and returns token pair
func (s *AuthService) Login(username, password string) (*models.User, *models.TokenPair, error) {
	// Trim spaces
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	// Validasi field tidak boleh kosong
	if username == "" {
		return nil, nil, errors.New("username wajib diisi")
	}
	if password == "" {
		return nil, nil, errors.New("password wajib diisi")
	}

	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, nil, errors.New("username atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, errors.New("username atau password salah")
	}

	// Generate token pair with role
	tokens, err := s.tokenService.GenerateTokenPair(user.ID, user.Username, user.RoleName)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// Logout revokes a single refresh token
func (s *AuthService) Logout(refreshToken string) error {
	// Trim spaces
	refreshToken = strings.TrimSpace(refreshToken)

	// Validasi token tidak boleh kosong
	if refreshToken == "" {
		return errors.New("refresh token wajib diisi")
	}

	return s.tokenService.RevokeRefreshToken(refreshToken)
}
