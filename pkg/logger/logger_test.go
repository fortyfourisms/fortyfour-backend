package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ================================================================
// TEST: IsSensitiveKey
// ================================================================

func TestIsSensitiveKey_KunciSensitifTerdaftar(t *testing.T) {
	kunciSensitif := []string{
		"password",
		"token",
		"secret",
		"jwt_secret",
		"jwt_refresh_secret",
		"authorization",
		"cookie",
		"api_key",
		"access_token",
		"refresh_token",
		"otp",
		"old_password",
		"new_password",
		"db_password",
		"redis_password",
		"rabbitmq_password",
	}

	for _, key := range kunciSensitif {
		t.Run(key, func(t *testing.T) {
			assert.True(t, IsSensitiveKey(key),
				"key %q seharusnya terdeteksi sebagai sensitif", key)
		})
	}
}

func TestIsSensitiveKey_KunciTidakSensitif(t *testing.T) {
	kunciAman := []string{
		"username",
		"email",
		"name",
		"id",
		"created_at",
		"updated_at",
		"status",
		"role",
		"message",
		"data",
		"",
	}

	for _, key := range kunciAman {
		t.Run(key, func(t *testing.T) {
			assert.False(t, IsSensitiveKey(key),
				"key %q seharusnya tidak terdeteksi sebagai sensitif", key)
		})
	}
}

func TestIsSensitiveKey_CaseInsensitive(t *testing.T) {
	// IsSensitiveKey menggunakan strings.ToLower sebelum lookup,
	// sehingga harus konsisten untuk semua variasi huruf besar/kecil.
	cases := []struct {
		input    string
		expected bool
	}{
		{"PASSWORD", true},
		{"Password", true},
		{"pAsSwOrD", true},
		{"TOKEN", true},
		{"Token", true},
		{"SECRET", true},
		{"API_KEY", true},
		{"Api_Key", true},
		{"Authorization", true},
		{"AUTHORIZATION", true},
		{"USERNAME", false},
		{"Email", false},
		{"NAME", false},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			assert.Equal(t, c.expected, IsSensitiveKey(c.input),
				"IsSensitiveKey(%q) harus mengembalikan %v", c.input, c.expected)
		})
	}
}

func TestIsSensitiveKey_StringKosong(t *testing.T) {
	assert.False(t, IsSensitiveKey(""),
		"string kosong seharusnya tidak sensitif")
}

func TestIsSensitiveKey_KunciMiripTapiTidakSama(t *testing.T) {
	// Key yang mirip dengan kunci sensitif tapi tidak terdaftar
	kunciMirip := []string{
		"passwords",       // jamak dari password
		"tokens",          // jamak dari token
		"my_secret",       // mengandung secret tapi bukan persis
		"tokenize",        // prefix token tapi bukan persis
		"reset_password",  // mengandung password tapi bukan persis
		"secret_question", // mengandung secret tapi bukan persis
		"cookiename",      // mengandung cookie tapi bukan persis
	}

	for _, key := range kunciMirip {
		t.Run(key, func(t *testing.T) {
			assert.False(t, IsSensitiveKey(key),
				"key %q tidak terdaftar sebagai sensitif, seharusnya false", key)
		})
	}
}