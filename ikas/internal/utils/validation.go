package utils

import (
	"regexp"
	"strings"
)

// Normalisasi input: trim dan hilangkan multiple spaces
func NormalizeInput(input string) string {
	input = strings.TrimSpace(input)

	multipleSpaces := regexp.MustCompile(`\s+`)
	input = multipleSpaces.ReplaceAllString(input, " ")

	return input
}

// Validasi karakter berbahaya untuk SQL Injection
func ContainsSQLInjectionPattern(input string) bool {
	// Pattern umum SQL injection
	dangerousPatterns := []string{
		"--",         // SQL comment
		";",          // Multiple statements
		"'",          // String delimiter
		"\"",         // String delimiter
		"/*",         // Multi-line comment
		"*/",         // Multi-line comment
		"xp_",        // Extended stored procedures
		"sp_",        // Stored procedures
		"exec",       // Execute command
		"execute",    // Execute command
		"drop",       // Drop command
		"insert",     // Insert command
		"delete",     // Delete command
		"update",     // Update command
		"union",      // Union query
		"select",     // Select command
		"create",     // Create command
		"alter",      // Alter command
		"shutdown",   // Shutdown command
		"script",     // Script tag
		"javascript", // JavaScript
		"<script",    // Script tag
		"</script>",  // Script tag
		"onerror",    // Event handler
		"onload",     // Event handler
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}
	return false
}

// Validasi karakter yang diizinkan
func IsValidInput(input string) bool {
	// Hanya izinkan: huruf, angka, spasi, dan beberapa karakter khusus umum
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,()&]+$`)
	return validPattern.MatchString(input)
}

// Validasi UUID format untuk mencegah injection via ID
func IsValidUUID(id string) bool {
	uuidPattern := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	return uuidPattern.MatchString(id)
}
