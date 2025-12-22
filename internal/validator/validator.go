package validator

import (
	"fmt"
	"strings"

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
