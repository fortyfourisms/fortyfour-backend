package services

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo       *repository.UserRepository
	uploadPath string
}

func NewUserService(repo *repository.UserRepository, uploadPath string) *UserService {
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
	// Validate username exists
	exists, err := s.repo.UsernameExists(req.Username, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Validate email exists
	exists, err = s.repo.EmailExists(req.Email, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
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
		exists, err := s.repo.UsernameExists(*req.Username, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("username already exists")
		}
		user.Username = *req.Username
	}

	// Update email if provided
	if req.Email != nil {
		exists, err := s.repo.EmailExists(*req.Email, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	// Update role and jabatan
	if req.RoleID != nil {
		user.RoleID = req.RoleID
	}
	if req.IDJabatan != nil {
		user.IDJabatan = req.IDJabatan
	}

	if err := s.repo.Update(user); err != nil {
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

func (s *UserService) UpdatePassword(id string, req dto.UpdateUserPasswordRequest) error {
	// Get current password
	currentPassword, err := s.repo.GetPasswordByID(id)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
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
		FotoProfile: user.FotoProfile,
		Banner:      user.Banner,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}
}
