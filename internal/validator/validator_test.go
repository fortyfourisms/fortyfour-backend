package validator

import (
	"strings"
	"testing"
)

type TestStruct struct {
	Name  string `validate:"required,min=3,max=10"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18"`
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name: "Valid struct",
			input: TestStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "Invalid struct - missing fields",
			input: TestStruct{
				Name:  "",
				Email: "",
				Age:   0,
			},
			wantErr: true,
		},
		{
			name: "Invalid struct - invalid email",
			input: TestStruct{
				Name:  "John",
				Email: "invalid-email",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "Invalid struct - min length",
			input: TestStruct{
				Name:  "Jo",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "Invalid struct - max length",
			input: TestStruct{
				Name:  "JohnDoeTheGreat",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"Valid email", "test@example.com", true},
		{"Invalid email no domain", "test@", false},
		{"Invalid email no user", "@example.com", false},
		{"Invalid email no at", "testexample.com", false},
		{"Empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{"Valid username", "john", true},
		{"Too short", "jo", false},
		{"Too long", strings.Repeat("a", 51), false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateUsername(tt.username); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	validConfig := DefaultPasswordConfig()

	tests := []struct {
		name         string
		password     string
		config       PasswordValidationConfig
		personalInfo []string
		wantErr      bool
	}{
		{
			name:     "Valid password",
			password: "S0m3Unique$tr1ng!",
			config:   validConfig,
			wantErr:  false,
		},
		{
			name:     "Too short",
			password: "Sh1!",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:     "No upper case",
			password: "password123!",
			config:   validConfig, // RequireUpper check will fail first or Common check
			wantErr:  true,
		},
		{
			name:     "No lower case",
			password: "PASSWORD123!",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:     "No number",
			password: "Password!",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:     "No special",
			password: "Password123",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:     "Common password",
			password: "password123!",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:         "Contains personal info (username)",
			password:     "JohnStrongXYZ!",
			config:       validConfig,
			personalInfo: []string{"john", "john@example.com"},
			wantErr:      true,
		},
		{
			name:         "Contains personal info (email part)",
			password:     "ExampleStrongXYZ!",
			config:       validConfig,
			personalInfo: []string{"john", "example@test.com"},
			wantErr:      true,
		},
		{
			name:         "Contains personal info (reversed)",
			password:     "nhoJStrongXYZ!",
			config:       validConfig,
			personalInfo: []string{"john"},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password, tt.config, tt.personalInfo...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNewPassword(t *testing.T) {
	config := DefaultPasswordConfig()

	tests := []struct {
		name    string
		newPass string
		oldPass string
		wantErr bool
	}{
		{
			name:    "Different passwords",
			newPass: "N3wUnique$tr1ng!",
			oldPass: "OldUnique$tr1ng!",
			wantErr: false,
		},
		{
			name:    "Same passwords",
			newPass: "OldUnique$tr1ng!",
			oldPass: "OldUnique$tr1ng!",
			wantErr: true,
		},
		{
			name:    "Invalid new password",
			newPass: "weak",
			oldPass: "OldUnique$tr1ng!",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNewPassword(tt.newPass, tt.oldPass, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNewPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test internal helper functions directly if needed, but they are unexported.
	// Since we are in the same package (package validator), we can test them.

	// Test containsCommonRoot
	if !containsCommonRoot("someadminpassword") { // "admin" is root
		t.Error("containsCommonRoot failed for 'someadminpassword'")
	}
	if containsCommonRoot("uniqueXYZ") {
		t.Error("containsCommonRoot failed for 'uniqueXYZ'")
	}

	// Test reverseString
	if got := reverseString("abc"); got != "cba" {
		t.Errorf("reverseString('abc') = %s, want cba", got)
	}
}
