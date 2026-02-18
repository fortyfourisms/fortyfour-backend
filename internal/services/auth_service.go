package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/validator"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
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

/*
------------------------------

	Helper: map models -> dto
	------------------------------
*/
func mapTokenPairToDTO(m *models.TokenPair) *dto.TokenPair {
	if m == nil {
		return nil
	}
	// Format ExpiresAt (time.Time) -> string for DTO
	return &dto.TokenPair{
		AccessToken:  m.AccessToken,
		RefreshToken: m.RefreshToken,
		ExpiresAt:    m.ExpiresAt.Format(time.RFC3339),
	}
}

/* ===================== Register ===================== */
// Register creates a new user and returns token pair DTO
func (s *AuthService) Register(
	req dto.RegisterRequest,
	perusahaanService PerusahaanServiceInterface,
) (*models.User, *dto.TokenPair, error) {
	// sanitize
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	email := strings.TrimSpace(req.Email)

	// validations
	if username == "" {
		return nil, nil, errors.New("username wajib diisi")
	}
	if password == "" {
		return nil, nil, errors.New("password wajib diisi")
	}
	if email == "" {
		return nil, nil, errors.New("email wajib diisi")
	}
	if !validator.ValidateEmail(email) {
		return nil, nil, errors.New("email tidak valid")
	}
	if !validator.ValidateUsername(username) {
		return nil, nil, errors.New("username harus 3-50 karakter")
	}

	// ===== VALIDASI PASSWORD =====
	personalInfo := []string{username, email}
	config := validator.DefaultPasswordConfig()

	// Validasi password dengan semua kriteria keamanan
	if err := validator.ValidatePassword(password, config, personalInfo...); err != nil {
		return nil, nil, err
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

	// Handle company creation/selection logic
	var idPerusahaan *string

	// Validate empty string for nama_perusahaan
	var namaPerusahaanTrimmed string
	if req.NamaPerusahaan != nil {
		namaPerusahaanTrimmed = strings.TrimSpace(*req.NamaPerusahaan)
	}

	var idPerusahaanTrimmed string
	if req.IDPerusahaan != nil {
		idPerusahaanTrimmed = strings.TrimSpace(*req.IDPerusahaan)
	}

	// Validate: Cannot provide both nama_perusahaan AND id_perusahaan
	if namaPerusahaanTrimmed != "" && idPerusahaanTrimmed != "" {
		return nil, nil, errors.New("tidak bisa mengisi nama_perusahaan dan id_perusahaan bersamaan")
	}

	// SCENARIO 1: User wants to CREATE new company (+ Add New Company clicked)
	if namaPerusahaanTrimmed != "" {
		// Create company with only nama_perusahaan (id_sub_sektor will be NULL)
		createReq := dto.CreatePerusahaanRequest{
			NamaPerusahaan: &namaPerusahaanTrimmed,
			// IDSubSektor is nil - will be NULL in database
		}

		perusahaan, err := perusahaanService.Create(createReq)
		if err != nil {
			return nil, nil, fmt.Errorf("gagal membuat perusahaan: %v", err)
		}
		idPerusahaan = &perusahaan.ID
	} else if idPerusahaanTrimmed != "" {
		// SCENARIO 2: User selects existing company from dropdown
		// Validate company exists
		_, err := perusahaanService.GetByID(idPerusahaanTrimmed)
		if err != nil {
			return nil, nil, errors.New("perusahaan tidak ditemukan")
		}
		idPerusahaan = &idPerusahaanTrimmed
	}
	// SCENARIO 3: Neither provided - user will select company later (id_perusahaan = nil)

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		Username:     username,
		Password:     string(hashed),
		Email:        email,
		RoleID:       req.RoleID,
		IDJabatan:    req.IDJabatan,
		IDPerusahaan: idPerusahaan,
	}

	if err := s.userRepo.Create(user); err != nil {
		// ROLLBACK: If company was created, delete it
		if namaPerusahaanTrimmed != "" && idPerusahaan != nil {
			perusahaanService.Delete(*idPerusahaan)
		}
		return nil, nil, err
	}

	// Fetch user kembali untuk mendapatkan role_name
	user, err = s.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, nil, err
	}

	// Generate token pair with role (returns *models.TokenPair)
	modelTokens, err := s.tokenService.GenerateTokenPair(user.ID, user.Username, user.RoleName)
	if err != nil {
		return nil, nil, err
	}

	return user, mapTokenPairToDTO(modelTokens), nil
}

/* ===================== Login ===================== */
// Login authenticates a user and returns user + token pair DTO (tokens == nil if MFA required)
func (s *AuthService) Login(username, password string) (*models.User, *dto.TokenPair, error) {
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

	// if MFA enabled — return user with nil tokens (handler will create pending token)
	if user.MFAEnabled {
		return user, nil, nil
	}

	// Generate token pair with role
	modelTokens, err := s.tokenService.GenerateTokenPair(user.ID, user.Username, user.RoleName)
	if err != nil {
		return nil, nil, err
	}

	return user, mapTokenPairToDTO(modelTokens), nil
}

/* ===================== Logout ===================== */
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

/* ===================== MFA Section ===================== */

// SetupMFA - generate TOTP provisioning URI and secret; save secret in temporary redis store
func (s *AuthService) SetupMFA(userID string) (string, string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", "", err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "FortyFour",
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", err
	}

	secret := key.Secret()
	redisKey := fmt.Sprintf("mfa_setup:%s", userID)

	if err := s.tokenService.redis.Set(redisKey, secret, 10*time.Minute); err != nil {
		return "", "", err
	}

	return key.URL(), secret, nil
}

// EnableMFA - verify first code and persist secret into users table
func (s *AuthService) EnableMFA(userID, code string) error {
	redisKey := fmt.Sprintf("mfa_setup:%s", userID)
	secret, err := s.tokenService.redis.Get(redisKey)
	if err != nil {
		return errors.New("mfa setup expired or not found")
	}

	if !totp.Validate(code, secret) {
		return errors.New("invalid mfa code")
	}

	if err := s.userRepo.SetMFA(userID, &secret, true); err != nil {
		return err
	}

	_ = s.tokenService.redis.Delete(redisKey)
	return nil
}

// CreateMFAPending - after password step create a short-lived pending token stored in redis
func (s *AuthService) CreateMFAPending(userID string) (string, error) {
	token := uuid.New().String()
	key := fmt.Sprintf("mfa_pending:%s", token)

	if err := s.tokenService.redis.Set(key, userID, 5*time.Minute); err != nil {
		return "", err
	}
	return token, nil
}

// VerifyMFA - verify pending mfa token + totp code and return user + tokens (dto)
func (s *AuthService) VerifyMFA(mfaToken, code string) (*models.User, *dto.TokenPair, error) {
	key := fmt.Sprintf("mfa_pending:%s", mfaToken)
	userID, err := s.tokenService.redis.Get(key)
	if err != nil {
		return nil, nil, errors.New("invalid or expired mfa token")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, err
	}

	if user.MFASecret == nil {
		return nil, nil, errors.New("mfa not configured")
	}

	if !totp.Validate(code, *user.MFASecret) {
		return nil, nil, errors.New("invalid mfa code")
	}

	modelTokens, err := s.tokenService.GenerateTokenPair(user.ID, user.Username, user.RoleName)
	if err != nil {
		return nil, nil, err
	}

	_ = s.tokenService.redis.Delete(key)
	return user, mapTokenPairToDTO(modelTokens), nil
}

/* ===================== Microsoft-style MFA Methods ===================== */

// Creates a token that allows user to setup MFA (Microsoft-style)
func (s *AuthService) CreateMFASetupToken(userID string) (string, error) {
	token := uuid.New().String()
	key := fmt.Sprintf("mfa_setup_token:%s", token)

	if err := s.tokenService.redis.Set(key, userID, 10*time.Minute); err != nil {
		return "", err
	}
	return token, nil
}

// Validates setup token and returns userID
func (s *AuthService) ValidateMFASetupToken(setupToken string) (string, error) {
	key := fmt.Sprintf("mfa_setup_token:%s", setupToken)
	userID, err := s.tokenService.redis.Get(key)
	if err != nil {
		return "", errors.New("invalid or expired setup token")
	}
	return userID, nil
}

// Verify code, enable MFA, and immediately return tokens (Microsoft-style)
func (s *AuthService) EnableMFAAndLogin(userID, code string) (*models.User, *dto.TokenPair, error) {
	redisKey := fmt.Sprintf("mfa_setup:%s", userID)
	secret, err := s.tokenService.redis.Get(redisKey)
	if err != nil {
		return nil, nil, errors.New("mfa setup expired or not found")
	}

	if !totp.Validate(code, secret) {
		return nil, nil, errors.New("invalid mfa code")
	}

	// Enable MFA
	if err := s.userRepo.SetMFA(userID, &secret, true); err != nil {
		return nil, nil, err
	}

	// Clean up setup data
	_ = s.tokenService.redis.Delete(redisKey)

	// Get updated user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, err
	}

	// Generate tokens immediately
	modelTokens, err := s.tokenService.GenerateTokenPair(user.ID, user.Username, user.RoleName)
	if err != nil {
		return nil, nil, err
	}

	return user, mapTokenPairToDTO(modelTokens), nil
}
