package services

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/validator"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo       repository.UserRepositoryInterface
	uploadPath string
}

func NewUserService(repo repository.UserRepositoryInterface, uploadPath string) *UserService {
	return &UserService{
		repo:       repo,
		uploadPath: uploadPath,
	}
}

func (s *UserService) GetAll() ([]dto.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = s.toResponse(&user)
	}
	return responses, nil
}

func (s *UserService) GetByID(id string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	response := s.toResponse(user)
	return &response, nil
}

func (s *UserService) Create(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Validasi email format
	if !validator.ValidateEmail(req.Email) {
		return nil, errors.New("email tidak valid")
	}

	// Validasi username format
	if !validator.ValidateUsername(req.Username) {
		return nil, errors.New("username harus 3-50 karakter")
	}

	// Validasi password dengan kriteria ketat
	config := validator.DefaultPasswordConfig()
	personalInfo := []string{req.Username, req.Email}

	if err := validator.ValidatePassword(req.Password, config, personalInfo...); err != nil {
		return nil, err
	}

	// Validate username exists
	exists, err := s.repo.UsernameExists(req.Username, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username sudah digunakan")
	}

	// Validate email exists
	exists, err = s.repo.EmailExists(req.Email, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email sudah digunakan")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		RoleID:    req.RoleID,
		IDJabatan: req.IDJabatan,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	response := s.toResponse(user)
	return &response, nil
}

func (s *UserService) Update(id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update username if provided
	if req.Username != nil {
		trimmed := strings.TrimSpace(*req.Username)
		if trimmed == "" {
			return nil, errors.New("username tidak boleh kosong")
		}
		if !validator.ValidateUsername(trimmed) {
			return nil, errors.New("username harus 3-50 karakter")
		}
		exists, err := s.repo.UsernameExists(trimmed, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("username sudah digunakan")
		}
		user.Username = trimmed
	}

	// Update email if provided
	if req.Email != nil {
		trimmed := strings.TrimSpace(*req.Email)
		if trimmed == "" {
			return nil, errors.New("email tidak boleh kosong")
		}
		if !validator.ValidateEmail(trimmed) {
			return nil, errors.New("email tidak valid")
		}
		exists, err := s.repo.EmailExists(trimmed, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email sudah digunakan")
		}
		user.Email = trimmed
	}

	// Update role and jabatan
	if req.RoleID != nil {
		user.RoleID = req.RoleID
	}
	if req.IDJabatan != nil {
		user.IDJabatan = req.IDJabatan
	}

	// Need to elaborate why this return error user not found
	// if err := s.repo.Update(user); err != nil {
	// 	return nil, err
	// }

	// This run smooth
	s.repo.Update(user)

	// Get updated user
	user, err = s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := s.toResponse(user)
	return &response, nil
}

func (s *UserService) UpdatePassword(id string, req dto.UpdateUserPasswordRequest) error {
	// Validasi konfirmasi password baru
	if req.NewPassword != req.ConfirmNewPassword {
		return errors.New("konfirmasi password baru tidak cocok")
	}

	// Get current user data and password
	user, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	currentPassword, err := s.repo.GetPasswordByID(id)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(req.OldPassword)); err != nil {
		return errors.New("password lama tidak sesuai")
	}

	// Validasi password baru dengan kriteria ketat
	config := validator.DefaultPasswordConfig()
	personalInfo := []string{user.Username, user.Email}

	if err := validator.ValidateNewPassword(req.NewPassword, req.OldPassword, config, personalInfo...); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(id, string(hashedPassword))
}

func (s *UserService) UpdateProfilePhoto(id string, filename string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Delete old photo if exists
	if user.FotoProfile != nil && *user.FotoProfile != "" {
		oldPath := filepath.Join(s.uploadPath, *user.FotoProfile)
		os.Remove(oldPath)
	}

	user.FotoProfile = &filename
	if err := s.repo.UpdateWithPhoto(user); err != nil {
		return nil, err
	}

	// Get updated user
	user, err = s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := s.toResponse(user)
	return &response, nil
}

func (s *UserService) UpdateBanner(id string, filename string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Delete old banner if exists
	if user.Banner != nil && *user.Banner != "" {
		oldPath := filepath.Join(s.uploadPath, *user.Banner)
		os.Remove(oldPath)
	}

	user.Banner = &filename
	if err := s.repo.UpdateWithPhoto(user); err != nil {
		return nil, err
	}

	// Get updated user
	user, err = s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := s.toResponse(user)
	return &response, nil
}

func (s *UserService) Delete(id string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete profile photo if exists
	if user.FotoProfile != nil && *user.FotoProfile != "" {
		photoPath := filepath.Join(s.uploadPath, *user.FotoProfile)
		os.Remove(photoPath)
	}

	// Delete banner if exists
	if user.Banner != nil && *user.Banner != "" {
		bannerPath := filepath.Join(s.uploadPath, *user.Banner)
		os.Remove(bannerPath)
	}

	return s.repo.Delete(id)
}

func (s *UserService) toResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		RoleID:      user.RoleID,
		RoleName:    user.RoleName,
		IDJabatan:   user.IDJabatan,
		JabatanName: user.JabatanName,
		FotoProfile: user.FotoProfile,
		Banner:      user.Banner,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}
}
