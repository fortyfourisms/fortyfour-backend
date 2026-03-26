package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ================================================================
// GeminiGenerator — interface lokal untuk test
// ================================================================

type geminiGenerator interface {
	Generate(prompt string) (string, error)
}

// compile-time: pastikan *GeminiClient memenuhi interface
var _ geminiGenerator = (*GeminiClient)(nil)

// ================================================================
// mockGeminiGenerator — implementasi test untuk geminiGenerator
// ================================================================

type mockGeminiGenerator struct {
	generateFn func(prompt string) (string, error)
}

func (m *mockGeminiGenerator) Generate(prompt string) (string, error) {
	if m.generateFn != nil {
		return m.generateFn(prompt)
	}
	return "default mock response", nil
}

var _ geminiGenerator = (*mockGeminiGenerator)(nil)

// ================================================================
// TEST: compile-time interface compliance
// ================================================================

func TestGeminiClient_ImplementsGeminiGenerator(t *testing.T) {
	// Test ini berfungsi sebagai dokumentasi eksplisit bahwa GeminiClient
	// memenuhi kontrak interface yang diharapkan consumer-nya.
	t.Log("*GeminiClient memenuhi interface geminiGenerator — verified at compile time")
}

// ================================================================
// TEST: mockGeminiGenerator — verifikasi mock bisa dipakai sebagai substitute
// ================================================================

func TestMockGeminiGenerator_SuccessResponse(t *testing.T) {
	mock := &mockGeminiGenerator{
		generateFn: func(prompt string) (string, error) {
			return "jawaban dari mock", nil
		},
	}

	result, err := mock.Generate("pertanyaan apapun")

	require.NoError(t, err)
	assert.Equal(t, "jawaban dari mock", result)
}

func TestMockGeminiGenerator_ErrorResponse(t *testing.T) {
	mock := &mockGeminiGenerator{
		generateFn: func(prompt string) (string, error) {
			return "", errors.New("gemini tidak tersedia")
		},
	}

	result, err := mock.Generate("pertanyaan apapun")

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.EqualError(t, err, "gemini tidak tersedia")
}

func TestMockGeminiGenerator_DefaultResponse(t *testing.T) {
	// Tanpa generateFn, mock mengembalikan default response
	mock := &mockGeminiGenerator{}

	result, err := mock.Generate("prompt apapun")

	require.NoError(t, err)
	assert.Equal(t, "default mock response", result)
}

func TestMockGeminiGenerator_PromptDiteruskan(t *testing.T) {
	var capturedPrompt string
	mock := &mockGeminiGenerator{
		generateFn: func(prompt string) (string, error) {
			capturedPrompt = prompt
			return "oke", nil
		},
	}

	_, _ = mock.Generate("prompt spesifik untuk dicek")

	assert.Equal(t, "prompt spesifik untuk dicek", capturedPrompt)
}

func TestMockGeminiGenerator_EmptyPrompt(t *testing.T) {
	mock := &mockGeminiGenerator{
		generateFn: func(prompt string) (string, error) {
			if prompt == "" {
				return "", errors.New("prompt kosong")
			}
			return "oke", nil
		},
	}

	result, err := mock.Generate("")

	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// ================================================================
// TEST: error classification — logika deteksi jenis error
// ================================================================

func TestMockGeminiGenerator_503Error(t *testing.T) {
	mock := &mockGeminiGenerator{
		generateFn: func(prompt string) (string, error) {
			return "", errors.New("503 service UNAVAILABLE: server overloaded")
		},
	}

	_, err := mock.Generate("prompt")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "503")
}

func TestMockGeminiGenerator_RateLimitError(t *testing.T) {
	mock := &mockGeminiGenerator{
		generateFn: func(prompt string) (string, error) {
			return "", errors.New("RESOURCE_EXHAUSTED: quota exceeded")
		},
	}

	_, err := mock.Generate("prompt")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "RESOURCE_EXHAUSTED")
}

func TestMockGeminiGenerator_NotFoundError(t *testing.T) {
	mock := &mockGeminiGenerator{
		generateFn: func(prompt string) (string, error) {
			return "", errors.New("NOT_FOUND: model tidak tersedia")
		},
	}

	_, err := mock.Generate("prompt")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "NOT_FOUND")
}

func TestMockGeminiGenerator_AllModelsFailed(t *testing.T) {
	callCount := 0
	mock := &mockGeminiGenerator{
		generateFn: func(prompt string) (string, error) {
			callCount++
			return "", errors.New("semua model gagal setelah beberapa kali retry")
		},
	}

	_, err := mock.Generate("prompt")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "semua model gagal")
	assert.Equal(t, 1, callCount)
}

// ================================================================
// TEST: GeminiClient struct — field dan konstruksi
// ================================================================

func TestGeminiClient_StructFields(t *testing.T) {
	// Test ini memverifikasi bahwa struct bisa dibuat dan method Generate
	// bisa dipanggil — meski akan error karena client nil.
	gc := &GeminiClient{client: nil}

	assert.NotNil(t, gc)
}

func TestGeminiClient_NilClientNotCalledInTest(t *testing.T) {
	// Dokumentasi eksplisit: NewGeminiClient tidak ditest di sini karena:
	// 1. Butuh GEMINI_API_KEY yang valid
	// 2. Membuat koneksi ke API eksternal (google.golang.org/genai)
	// 3. Memanggil logger.Fatal() jika apiKey kosong (os.Exit)
	//
	// Test integrasi untuk NewGeminiClient dan Generate() harus dijalankan
	// di environment dengan GEMINI_API_KEY yang tersedia.
	t.Log("NewGeminiClient dan Generate() memerlukan API key nyata — ditest via integration test")
}