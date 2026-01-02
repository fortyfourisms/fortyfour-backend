package services

import (
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"testing"
)

func TestNewAuthService(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userRepo     repository.UserRepositoryInterface
		tokenService *TokenService
		want         *AuthService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAuthService(tt.userRepo, tt.tokenService)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		userRepo     repository.UserRepositoryInterface
		tokenService *TokenService
		// Named input parameters for target function.
		username  string
		password  string
		email     string
		roleID    *string
		idJabatan *string
		want      *models.User
		want2     *models.TokenPair
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAuthService(tt.userRepo, tt.tokenService)
			got, got2, gotErr := s.Register(tt.username, tt.password, tt.email, tt.roleID, tt.idJabatan)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Register() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Register() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Register() = %v, want %v", got, tt.want)
			}
			if true {
				t.Errorf("Register() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		userRepo     repository.UserRepositoryInterface
		tokenService *TokenService
		// Named input parameters for target function.
		username string
		password string
		want     *models.User
		want2    *models.TokenPair
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAuthService(tt.userRepo, tt.tokenService)
			got, got2, gotErr := s.Login(tt.username, tt.password)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Login() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Login() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
			if true {
				t.Errorf("Login() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		userRepo     repository.UserRepositoryInterface
		tokenService *TokenService
		// Named input parameters for target function.
		refreshToken string
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAuthService(tt.userRepo, tt.tokenService)
			gotErr := s.Logout(tt.refreshToken)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Logout() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Logout() succeeded unexpectedly")
			}
		})
	}
}
