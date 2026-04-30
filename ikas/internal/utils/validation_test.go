package utils

import (
	"fmt"
	"testing"
)

func TestNormalizeInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"No extra spaces", "hello", "hello"},
		{"Leading/trailing spaces", "  hello  ", "hello"},
		{"Multiple internal spaces", "hello   world", "hello world"},
		{"Mixed spaces", "  hello   world  ", "hello world"},
		{"Tabs and newlines", "hello\t\nworld", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeInput(tt.input); got != tt.want {
				t.Errorf("NormalizeInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsSQLInjectionPattern(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Safe input", "hello world", false},
		{"Safe input with numbers", "user123", false},
		{"SQL comment", "admin' --", true},
		{"Multiple statements", "'; DROP TABLE users;", true},
		{"Single quote", "it's me", true},
		{"Double quote", "say \"hello\"", true},
		{"Stored procedure xp_", "xp_cmdshell", true},
		{"Stored procedure sp_", "sp_who", true},
		{"Exec command", "exec query", true},
		{"Execute command", "execute query", true},
		{"Drop command", "drop table", true},
		{"Insert command", "insert into", true},
		{"Delete command", "delete from", true},
		{"Update command", "update table", true},
		{"Union command", "union select", true},
		{"Select command", "select *", true},
		{"Create command", "create table", true},
		{"Alter command", "alter table", true},
		{"Shutdown command", "shutdown", true},
		{"Script tag", "<script>alert(1)</script>", true},
		{"Javascript", "javascript:void(0)", true},
		{"Onerror handler", "onerror=alert(1)", true},
		{"Onload handler", "onload=alert(1)", true},
		{"Uppercase pattern", "SELECT * FROM users", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsSQLInjectionPattern(tt.input); got != tt.want {
				t.Errorf("ContainsSQLInjectionPattern(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Valid alphanumeric", "User123", true},
		{"Valid with spaces and dots", "John Doe. Manager", true},
		{"Valid with parentheses", "Software (Engineering)", true},
		{"Valid with & and -", "R&D-Dept", true},
		{"Valid with comma and underscore", "City, State_Zip", true},
		{"Invalid character !", "hello!", false},
		{"Invalid character @", "user@domain", false},
		{"Invalid character #", "tag#1", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidInput(tt.input); got != tt.want {
				t.Errorf("IsValidInput(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{"Valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"Invalid UUID - too short", "550e8400-e29b-41d4-a716", false},
		{"Invalid UUID - no hyphens", "550e8400e29b41d4a716446655440000", false},
		{"Invalid UUID - non hex", "g50e8400-e29b-41d4-a716-446655440000", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidUUID(tt.id); got != tt.want {
				t.Errorf("IsValidUUID(%v) = %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}

func TestRoundFloat(t *testing.T) {
	tests := []struct {
		name      string
		val       float64
		precision int
		want      float64
	}{
		{"Round to 2 decimals", 3.14159, 2, 3.14},
		{"Round up", 3.145, 2, 3.15},
		{"Round to 0 decimals", 3.5, 0, 4.0},
		{"Round negative", -3.14159, 2, -3.14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundFloat(tt.val, tt.precision); got != tt.want {
				t.Errorf("RoundFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatFloat2Decimal(t *testing.T) {
	if got := FormatFloat2Decimal(3.14159); got != 3.14 {
		t.Errorf("FormatFloat2Decimal() = %v, want 3.14", got)
	}
}

func TestValidateFloat(t *testing.T) {
	tests := []struct {
		name          string
		val           float64
		fieldName     string
		allowNegative bool
		want          float64
		wantErr       bool
	}{
		{"Valid positive", 10.555, "Score", false, 10.56, false},
		{"Valid zero", 0.0, "Score", false, 0.0, false},
		{"Negative not allowed", -1.0, "Score", false, 0, true},
		{"Negative allowed", -10.555, "Score", true, -10.56, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateFloat(tt.val, tt.fieldName, tt.allowNegative)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateFloat() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && err.Error() != fmt.Sprintf("%s tidak boleh negatif", tt.fieldName) {
				t.Errorf("ValidateFloat() error message = %v, want %s tidak boleh negatif", err.Error(), tt.fieldName)
			}
		})
	}
}
