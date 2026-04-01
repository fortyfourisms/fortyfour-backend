package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"
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

// ================================================================
// TEST: parseLevel
// ================================================================

func TestParseLevel_AllBranches(t *testing.T) {
	cases := []struct {
		input    string
		expected zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"DEBUG", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"INFO", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"WARN", zerolog.WarnLevel},
		{"warning", zerolog.WarnLevel},
		{"WARNING", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"ERROR", zerolog.ErrorLevel},
		{"fatal", zerolog.FatalLevel},
		{"FATAL", zerolog.FatalLevel},
		{"panic", zerolog.PanicLevel},
		{"PANIC", zerolog.PanicLevel},
		{"disabled", zerolog.Disabled},
		{"DISABLED", zerolog.Disabled},
		{"off", zerolog.Disabled},
		{"OFF", zerolog.Disabled},
		// default fallback → info
		{"", zerolog.InfoLevel},
		{"unknown", zerolog.InfoLevel},
		{"verbose", zerolog.InfoLevel},
		{"trace", zerolog.InfoLevel},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got := parseLevel(c.input)
			assert.Equal(t, c.expected, got,
				"parseLevel(%q) seharusnya menghasilkan %v", c.input, c.expected)
		})
	}
}

// ================================================================
// TEST: Init
// ================================================================

func TestInit_DevelopmentMode(t *testing.T) {
	// Memastikan Init tidak panik dan global level tereset
	assert.NotPanics(t, func() {
		Init("debug", "development")
	})

	// Setelah Init("debug"), global level harus DebugLevel
	assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
}

func TestInit_ProductionMode(t *testing.T) {
	assert.NotPanics(t, func() {
		Init("error", "production")
	})

	assert.Equal(t, zerolog.ErrorLevel, zerolog.GlobalLevel())
}

func TestInit_DefaultLevelOnUnknown(t *testing.T) {
	// Level tidak dikenal → fallback ke InfoLevel
	assert.NotPanics(t, func() {
		Init("unknown-level", "production")
	})

	assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
}

func TestInit_StagingEnvUseJSONLogger(t *testing.T) {
	// Environment selain "development" harus menggunakan JSON logger (else branch)
	assert.NotPanics(t, func() {
		Init("info", "staging")
	})
}

func TestInit_EmptyEnvironment(t *testing.T) {
	// Empty string environment → masuk else branch (JSON logger)
	assert.NotPanics(t, func() {
		Init("warn", "")
	})

	assert.Equal(t, zerolog.WarnLevel, zerolog.GlobalLevel())
}

// ================================================================
// TEST: Get
// ================================================================

func TestGet_ReturnsLogger(t *testing.T) {
	Init("info", "production")
	l := Get()
	// zerolog.Logger adalah struct; pastikan tidak zero value dari sudut pandang level
	// Get() harus mengembalikan logger yang sama dengan yang diset oleh Init
	assert.NotNil(t, l)
}

// ================================================================
// TEST: WithField — redaksi otomatis kunci sensitif
// ================================================================

func TestWithField_SensitiveKey_IsRedacted(t *testing.T) {
	var buf bytes.Buffer
	// Override global logger agar output ke buffer
	log = zerolog.New(&buf).With().Timestamp().Logger()

	event := WithField("password", "rahasia123")
	assert.NotNil(t, event, "WithField harus mengembalikan event non-nil")
	event.Send()

	output := buf.String()
	// Nilai asli tidak boleh muncul
	assert.NotContains(t, output, "rahasia123",
		"nilai sensitif tidak boleh ditulis ke log")
	// Redaksi harus muncul
	assert.Contains(t, output, "[REDACTED]",
		"teks [REDACTED] harus muncul untuk kunci sensitif")
}

func TestWithField_NonSensitiveKey_LogsValue(t *testing.T) {
	var buf bytes.Buffer
	log = zerolog.New(&buf).With().Logger()

	event := WithField("username", "johndoe")
	assert.NotNil(t, event)
	event.Send()

	output := buf.String()
	assert.Contains(t, output, "johndoe",
		"nilai non-sensitif harus ditulis ke log")
	assert.NotContains(t, output, "[REDACTED]")
}

func TestWithField_CaseInsensitiveSensitiveKey(t *testing.T) {
	var buf bytes.Buffer
	log = zerolog.New(&buf).With().Logger()

	event := WithField("TOKEN", "super-secret-token")
	assert.NotNil(t, event)
	event.Send()

	output := buf.String()
	assert.NotContains(t, output, "super-secret-token")
	assert.Contains(t, output, "[REDACTED]")
}

func TestWithField_IntValue_NonSensitive(t *testing.T) {
	var buf bytes.Buffer
	log = zerolog.New(&buf).With().Logger()

	event := WithField("user_id", 42)
	assert.NotNil(t, event)
	event.Send()

	output := buf.String()
	assert.Contains(t, output, "42")
}

// ================================================================
// TEST: logging functions (Info, Infof, Warn, Warnf, Error, Errorf,
//       Debug, Debugf) — verifikasi tidak panik dan output ke buffer
// ================================================================

func setupBufferLogger() *bytes.Buffer {
	var buf bytes.Buffer
	// Reset global level agar tidak memfilter message dari test Init sebelumnya
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log = zerolog.New(&buf).Level(zerolog.DebugLevel)
	return &buf
}

func TestInfo_WritesMessage(t *testing.T) {
	buf := setupBufferLogger()
	Info("pesan info test")
	assert.Contains(t, buf.String(), "pesan info test")
}

func TestInfof_WritesFormattedMessage(t *testing.T) {
	buf := setupBufferLogger()
	Infof("halo %s nomor %d", "dunia", 7)
	output := buf.String()
	assert.Contains(t, output, "halo dunia nomor 7")
}

func TestWarn_WritesMessage(t *testing.T) {
	buf := setupBufferLogger()
	Warn("pesan warn test")
	assert.Contains(t, buf.String(), "pesan warn test")
}

func TestWarnf_WritesFormattedMessage(t *testing.T) {
	buf := setupBufferLogger()
	Warnf("warn %s", "formatted")
	assert.Contains(t, buf.String(), "warn formatted")
}

func TestError_WritesMessageAndError(t *testing.T) {
	buf := setupBufferLogger()
	err := errors.New("sesuatu gagal")
	Error(err, "terjadi kesalahan")
	output := buf.String()
	assert.Contains(t, output, "terjadi kesalahan")
	assert.Contains(t, output, "sesuatu gagal")
}

func TestErrorf_WritesFormattedMessage(t *testing.T) {
	buf := setupBufferLogger()
	err := errors.New("db error")
	Errorf(err, "gagal query tabel %s", "users")
	output := buf.String()
	assert.Contains(t, output, "gagal query tabel users")
	assert.Contains(t, output, "db error")
}

func TestDebug_WritesMessage(t *testing.T) {
	buf := setupBufferLogger()
	Debug("pesan debug test")
	assert.Contains(t, buf.String(), "pesan debug test")
}

func TestDebugf_WritesFormattedMessage(t *testing.T) {
	buf := setupBufferLogger()
	Debugf("debug nilai=%d", 99)
	assert.Contains(t, buf.String(), "debug nilai=99")
}

// ================================================================
// TEST: level field dalam output JSON (sanity check format)
// ================================================================

func TestInfo_OutputIsValidJSON(t *testing.T) {
	var buf bytes.Buffer
	log = zerolog.New(&buf)
	Info("cek json")

	line := strings.TrimSpace(buf.String())
	var m map[string]interface{}
	err := json.Unmarshal([]byte(line), &m)
	assert.NoError(t, err, "output logger harus berupa JSON valid")
	assert.Equal(t, "info", m["level"])
	assert.Equal(t, "cek json", m["message"])
}

func TestError_OutputContainsLevelError(t *testing.T) {
	var buf bytes.Buffer
	log = zerolog.New(&buf)
	Error(errors.New("err"), "msg error")

	line := strings.TrimSpace(buf.String())
	var m map[string]interface{}
	err := json.Unmarshal([]byte(line), &m)
	assert.NoError(t, err)
	assert.Equal(t, "error", m["level"])
}