package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate melakukan validasi terhadap struct
func Validate(data interface{}) error {
	err := validate.Struct(data)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

// formatValidationErrors mengubah validation errors menjadi pesan
func formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string

	for _, err := range errs {
		field := strings.ToLower(err.Field())

		switch err.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s wajib diisi", field))
		case "email":
			messages = append(messages, fmt.Sprintf("%s harus berupa email yang valid", field))
		case "min":
			messages = append(messages, fmt.Sprintf("%s minimal %s karakter", field, err.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s maksimal %s karakter", field, err.Param()))
		default:
			messages = append(messages, fmt.Sprintf("%s tidak valid", field))
		}
	}
	return fmt.Errorf("%s", strings.Join(messages, ", "))
}

// ValidateEmail melakukan validasi khusus untuk email
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	err := validate.Var(email, "required,email")
	return err == nil
}

// ValidateUsername melakukan validasi khusus untuk username
func ValidateUsername(username string) bool {
	if username == "" {
		return false
	}
	err := validate.Var(username, "required,min=3,max=50")
	return err == nil
}

// ==================== PASSWORD VALIDATION ====================

// List password yang mudah ditebak
var commonPasswords = map[string]bool{
	// 1–10 PASSWORD BASE
	"password123!":  true,
	"password@123":  true,
	"password123@":  true,
	"password2024!": true,
	"password2025!": true,
	"password#123":  true,
	"password!@#":   true,
	"password1@#":   true,
	"password12!":   true,
	"password321!":  true,

	// 11–20 P@SSWORD VARIANT
	"p@ssword1!":    true,
	"p@ssword123!":  true,
	"p@ssword@123":  true,
	"p@ssword2024!": true,
	"p@ssword2025!": true,
	"p@ssw0rd1!":    true,
	"p@ssw0rd123!":  true,
	"p@ssw0rd@123":  true,
	"p@ssw0rd2024!": true,
	"p@ssw0rd2025!": true,

	// 21–30 ADMIN BASE
	"admin123!":  true,
	"admin@123":  true,
	"admin123@":  true,
	"admin2024!": true,
	"admin2025!": true,
	"admin#123":  true,
	"admin!@#":   true,
	"admin1@#":   true,
	"admin321!":  true,
	"admin4321!": true,

	// 31–40 QWERTY / KEYBOARD
	"qwerty123!": true,
	"qwerty@123": true,
	"qwertyui1!": true,
	"qwertyui@1": true,
	"qwerty!@#":  true,
	"1qaz2wsx!":  true,
	"1qaz@wsx":   true,
	"1qaz2wsx@":  true,
	"1qaz!2wsx":  true,
	"qazwsx123!": true,

	// 41–50 ABC PATTERN
	"abc123!":    true,
	"abc@123!":   true,
	"abcd1234!":  true,
	"abc12345!":  true,
	"abc!@#123":  true,
	"abc123@#":   true,
	"abcdef1!":   true,
	"abcdef123!": true,
	"abc321!":    true,
	"abc4321!":   true,

	// 51–60 USER / LOGIN
	"user123!":   true,
	"user@123":   true,
	"user123@":   true,
	"user2024!":  true,
	"user2025!":  true,
	"login123!":  true,
	"login@123":  true,
	"login2024!": true,
	"login2025!": true,
	"user!@#123": true,

	// 61–70 WELCOME / DEFAULT
	"welcome123!":  true,
	"welcome@123":  true,
	"welcome123@":  true,
	"welcome2024!": true,
	"welcome2025!": true,
	"welcome!@#":   true,
	"welcome1@#":   true,
	"welcome321!":  true,
	"welcome4321!": true,
	"welcomeabc1!": true,

	// 71–80 TEST / DEV
	"test123!":    true,
	"test@123":    true,
	"test123@":    true,
	"testing123!": true,
	"testing@123": true,
	"test2024!":   true,
	"test2025!":   true,
	"test!@#123":  true,
	"tester123!":  true,
	"tester@123":  true,

	// 81–90 SUPER / ROOT
	"superadmin1!":   true,
	"superadmin123!": true,
	"superadmin@123": true,
	"root123!":       true,
	"root@123":       true,
	"root2024!":      true,
	"root2025!":      true,
	"system123!":     true,
	"system@123":     true,
	"system2025!":    true,

	// 91–100 GENERIC STRONG-LOOKING BUT WEAK
	"default123!":  true,
	"default@123":  true,
	"changeme123!": true,
	"changeme@123": true,
	"secure123!":   true,
	"secure@123":   true,
	"company123!":  true,
	"company@123":  true,
	"office123!":   true,
	"office@123":   true,
}

// Kata inti password umum (untuk deteksi substring)
var commonPasswordRoots = []string{
	"password",
	"admin",
	"qwerty",
	"user",
	"login",
	"welcome",
	"root",
	"system",
	"test",
	"superadmin",
	"123",
	"aaa",
}

type PasswordValidationConfig struct {
	MinLength        int
	RequireUpper     bool
	RequireLower     bool
	RequireNumber    bool
	RequireSpecial   bool
	CheckCommon      bool
	CheckPersonal    bool
	CheckOldPassword bool
}

// DefaultPasswordConfig returns the default password validation configuration
func DefaultPasswordConfig() PasswordValidationConfig {
	return PasswordValidationConfig{
		MinLength:        8,
		RequireUpper:     true,
		RequireLower:     true,
		RequireNumber:    true,
		RequireSpecial:   true,
		CheckCommon:      true,
		CheckPersonal:    true,
		CheckOldPassword: true,
	}
}

// ValidatePassword melakukan validasi password dengan berbagai kriteria keamanan
func ValidatePassword(password string, config PasswordValidationConfig, personalInfo ...string) error {
	password = strings.TrimSpace(password)

	// 1. Validasi panjang minimum
	if len(password) < config.MinLength {
		return fmt.Errorf("password minimal %d karakter", config.MinLength)
	}

	// 2. Validasi karakter huruf besar
	if config.RequireUpper && !containsUpperCase(password) {
		return errors.New("password harus mengandung minimal 1 huruf besar (A-Z)")
	}

	// 3. Validasi karakter huruf kecil
	if config.RequireLower && !containsLowerCase(password) {
		return errors.New("password harus mengandung minimal 1 huruf kecil (a-z)")
	}

	// 4. Validasi angka
	if config.RequireNumber && !containsNumber(password) {
		return errors.New("password harus mengandung minimal 1 angka (0-9)")
	}

	// 5. Validasi karakter khusus
	if config.RequireSpecial && !containsSpecialChar(password) {
		return errors.New("password harus mengandung minimal 1 karakter khusus (!@#$%^&*()_+-=[]{}|;:,.<>?)")
	}

	// 6. Validasi password umum yang mudah ditebak
	if config.CheckCommon {
		if isCommonPassword(password) || containsCommonRoot(password) {
			return errors.New("password terlalu umum dan mudah ditebak, gunakan password yang lebih unik")
		}
	}

	// 7. Validasi tidak mengandung data pribadi
	if config.CheckPersonal && len(personalInfo) > 0 {
		if containsPersonalInfo(password, personalInfo) {
			return errors.New("password tidak boleh mengandung informasi pribadi Anda")
		}
	}

	return nil
}

// ValidateNewPassword validates new password and compares with old password
func ValidateNewPassword(newPassword, oldPassword string, config PasswordValidationConfig, personalInfo ...string) error {
	// Validasi password baru
	if err := ValidatePassword(newPassword, config, personalInfo...); err != nil {
		return err
	}

	// Validasi tidak sama dengan password lama
	if config.CheckOldPassword && newPassword == oldPassword {
		return errors.New("password baru tidak boleh sama dengan password lama")
	}

	return nil
}

// containsUpperCase checks if string contains uppercase letter
func containsUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// containsLowerCase checks if string contains lowercase letter
func containsLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// containsNumber checks if string contains number
func containsNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// containsSpecialChar checks if string contains special character
func containsSpecialChar(s string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, r := range s {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}

// isCommonPassword checks if password is in common password list
func isCommonPassword(password string) bool {
	// Check exact match (case insensitive)
	return commonPasswords[strings.ToLower(password)]
}

// containsCommonRoot checks if password contains common root words
func containsCommonRoot(password string) bool {
	p := strings.ToLower(password)

	for _, root := range commonPasswordRoots {
		if strings.Contains(p, root) {
			return true
		}
	}

	return false
}

// containsPersonalInfo checks if password contains personal information
func containsPersonalInfo(password string, personalInfo []string) bool {
	passwordLower := strings.ToLower(password)

	for _, info := range personalInfo {
		if info == "" {
			continue
		}

		infoLower := strings.ToLower(info)

		// Check if password contains the personal info
		if strings.Contains(passwordLower, infoLower) {
			return true
		}

		// Check email username part (before @)
		if strings.Contains(info, "@") {
			emailParts := strings.Split(info, "@")
			if len(emailParts) > 0 && emailParts[0] != "" {
				if strings.Contains(passwordLower, strings.ToLower(emailParts[0])) {
					return true
				}
			}
		}

		// Check reversed personal info
		if len(infoLower) >= 3 && strings.Contains(passwordLower, reverseString(infoLower)) {
			return true
		}
	}

	return false
}

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
