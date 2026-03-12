package repository

import (
	"testing"

	"fortyfour-backend/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// TestNewInMemoryChatRepo
// ============================================================

func TestNewInMemoryChatRepo(t *testing.T) {
	repo := NewInMemoryChatRepo()

	require.NotNil(t, repo)
	assert.NotNil(t, repo.data)
	assert.Empty(t, repo.data)
}

// ============================================================
// TestInMemoryChatRepo_Save
// ============================================================

func TestInMemoryChatRepo_Save_Success(t *testing.T) {
	repo := NewInMemoryChatRepo()

	err := repo.Save("session-1", "apa itu perusahaan?", "Perusahaan adalah...")

	assert.NoError(t, err)
	assert.Len(t, repo.data["session-1"], 1)
	assert.Equal(t, "apa itu perusahaan?", repo.data["session-1"][0].User)
	assert.Equal(t, "Perusahaan adalah...", repo.data["session-1"][0].Bot)
}

func TestInMemoryChatRepo_Save_MultipleMessages_SameSession(t *testing.T) {
	repo := NewInMemoryChatRepo()

	_ = repo.Save("session-1", "pertanyaan 1", "jawaban 1")
	_ = repo.Save("session-1", "pertanyaan 2", "jawaban 2")
	_ = repo.Save("session-1", "pertanyaan 3", "jawaban 3")

	assert.Len(t, repo.data["session-1"], 3)
	assert.Equal(t, "pertanyaan 1", repo.data["session-1"][0].User)
	assert.Equal(t, "pertanyaan 2", repo.data["session-1"][1].User)
	assert.Equal(t, "pertanyaan 3", repo.data["session-1"][2].User)
}

func TestInMemoryChatRepo_Save_MultipleSessions(t *testing.T) {
	repo := NewInMemoryChatRepo()

	_ = repo.Save("session-A", "tanya A", "jawab A")
	_ = repo.Save("session-B", "tanya B", "jawab B")

	assert.Len(t, repo.data, 2)
	assert.Len(t, repo.data["session-A"], 1)
	assert.Len(t, repo.data["session-B"], 1)
}

func TestInMemoryChatRepo_Save_OrderPreserved(t *testing.T) {
	repo := NewInMemoryChatRepo()

	messages := []struct{ user, bot string }{
		{"q1", "a1"},
		{"q2", "a2"},
		{"q3", "a3"},
		{"q4", "a4"},
		{"q5", "a5"},
	}

	for _, m := range messages {
		_ = repo.Save("sess", m.user, m.bot)
	}

	history := repo.data["sess"]
	require.Len(t, history, 5)
	for i, m := range messages {
		assert.Equal(t, m.user, history[i].User)
		assert.Equal(t, m.bot, history[i].Bot)
	}
}

// ============================================================
// TestInMemoryChatRepo_GetHistory
// ============================================================

func TestInMemoryChatRepo_GetHistory_Success(t *testing.T) {
	repo := NewInMemoryChatRepo()

	_ = repo.Save("session-1", "pertanyaan 1", "jawaban 1")
	_ = repo.Save("session-1", "pertanyaan 2", "jawaban 2")

	history, err := repo.GetHistory("session-1")

	assert.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, "pertanyaan 1", history[0].User)
	assert.Equal(t, "jawaban 1", history[0].Bot)
	assert.Equal(t, "pertanyaan 2", history[1].User)
}

func TestInMemoryChatRepo_GetHistory_SessionTidakAda_ReturnEmpty(t *testing.T) {
	repo := NewInMemoryChatRepo()

	history, err := repo.GetHistory("session-tidak-ada")

	// Mengembalikan slice kosong bukan error
	assert.NoError(t, err)
	assert.Empty(t, history)
}

func TestInMemoryChatRepo_GetHistory_ReturnsCorrectType(t *testing.T) {
	repo := NewInMemoryChatRepo()
	_ = repo.Save("s1", "q", "a")

	history, err := repo.GetHistory("s1")

	assert.NoError(t, err)
	assert.IsType(t, []dto.ChatHistory{}, history)
}

func TestInMemoryChatRepo_GetHistory_TidakMemengaruhiSession_Lain(t *testing.T) {
	repo := NewInMemoryChatRepo()

	_ = repo.Save("session-A", "tanya A", "jawab A")
	_ = repo.Save("session-B", "tanya B1", "jawab B1")
	_ = repo.Save("session-B", "tanya B2", "jawab B2")

	historyA, _ := repo.GetHistory("session-A")
	historyB, _ := repo.GetHistory("session-B")

	assert.Len(t, historyA, 1)
	assert.Len(t, historyB, 2)
}

// ============================================================
// TestInMemoryChatRepo_DeleteSession
// ============================================================

func TestInMemoryChatRepo_DeleteSession_Success(t *testing.T) {
	repo := NewInMemoryChatRepo()

	_ = repo.Save("session-1", "tanya", "jawab")
	require.Len(t, repo.data["session-1"], 1)

	err := repo.DeleteSession("session-1")

	assert.NoError(t, err)
	_, exists := repo.data["session-1"]
	assert.False(t, exists)
}

func TestInMemoryChatRepo_DeleteSession_NotFound(t *testing.T) {
	repo := NewInMemoryChatRepo()

	err := repo.DeleteSession("session-tidak-ada")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session-tidak-ada")
	assert.Contains(t, err.Error(), "not found")
}

func TestInMemoryChatRepo_DeleteSession_HanyaHapusTarget(t *testing.T) {
	repo := NewInMemoryChatRepo()

	_ = repo.Save("session-A", "tanya A", "jawab A")
	_ = repo.Save("session-B", "tanya B", "jawab B")

	err := repo.DeleteSession("session-A")

	assert.NoError(t, err)

	_, existsA := repo.data["session-A"]
	assert.False(t, existsA, "session-A harus terhapus")

	_, existsB := repo.data["session-B"]
	assert.True(t, existsB, "session-B tidak boleh ikut terhapus")
}

func TestInMemoryChatRepo_DeleteSession_IdempotentError(t *testing.T) {
	repo := NewInMemoryChatRepo()

	_ = repo.Save("session-1", "q", "a")
	_ = repo.DeleteSession("session-1")

	// Hapus lagi → error not found
	err := repo.DeleteSession("session-1")
	assert.Error(t, err)
}

func TestInMemoryChatRepo_DeleteSession_LaluSimpanUlang(t *testing.T) {
	repo := NewInMemoryChatRepo()

	_ = repo.Save("session-1", "lama", "lama")
	_ = repo.DeleteSession("session-1")

	// Simpan ulang dengan session ID yang sama
	err := repo.Save("session-1", "baru", "baru")
	assert.NoError(t, err)
	assert.Len(t, repo.data["session-1"], 1)
	assert.Equal(t, "baru", repo.data["session-1"][0].User)
}

// ============================================================
// TestInMemoryChatRepo_FullLifecycle
// ============================================================

func TestInMemoryChatRepo_FullLifecycle(t *testing.T) {
	repo := NewInMemoryChatRepo()

	// Simpan beberapa pesan
	assert.NoError(t, repo.Save("sess", "q1", "a1"))
	assert.NoError(t, repo.Save("sess", "q2", "a2"))

	// GetHistory
	history, err := repo.GetHistory("sess")
	assert.NoError(t, err)
	assert.Len(t, history, 2)

	// Delete
	assert.NoError(t, repo.DeleteSession("sess"))

	// GetHistory setelah delete → kosong
	history, err = repo.GetHistory("sess")
	assert.NoError(t, err)
	assert.Empty(t, history)

	// Delete lagi → error not found
	assert.Error(t, repo.DeleteSession("sess"))
}
