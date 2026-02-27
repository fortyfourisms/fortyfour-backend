package services

import (
	"errors"
	"testing"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// MOCK CHAT REPOSITORY
// ============================================================

type mockChatRepo struct {
	SaveFn func(sessionID, userMsg, botMsg string) error
}

func (m *mockChatRepo) GetHistory(sessionID string) ([]dto.ChatHistory, error) {
	return nil, nil
}

func (m *mockChatRepo) Save(sessionID, userMsg, botMsg string) error {
	if m.SaveFn != nil {
		return m.SaveFn(sessionID, userMsg, botMsg)
	}
	return nil
}

// Compile-time check
var _ repository.ChatRepository = (*mockChatRepo)(nil)

// ============================================================
// HELPER: newTestChatService
// Membuat ChatService dengan sqlmock dan mock repo.
// gemini di-set nil — hanya boleh dipakai untuk test yang tidak
// menyentuh GenerateSQLQuery / FormatQueryResults.
// ============================================================

func newTestChatService(t *testing.T) (*ChatService, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	svc := &ChatService{
		repo:   &mockChatRepo{},
		gemini: &utils.GeminiClient{}, // struct kosong; tidak dipanggil di test ini
		db:     db,
	}
	return svc, mock
}

// ============================================================
// TestChatService_Repo
// ============================================================

func TestChatService_Repo(t *testing.T) {
	mockRepo := &mockChatRepo{}
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	svc := &ChatService{repo: mockRepo, gemini: &utils.GeminiClient{}, db: db}

	assert.Equal(t, mockRepo, svc.Repo())
}

// ============================================================
// TestChatService_GetDatabaseSchema
// ============================================================

func TestChatService_GetDatabaseSchema_Success(t *testing.T) {
	svc, mock := newTestChatService(t)

	rows := sqlmock.NewRows([]string{
		"TABLE_NAME", "COLUMN_NAME", "DATA_TYPE", "COLUMN_KEY", "IS_NULLABLE",
	}).
		AddRow("perusahaan", "nama_perusahaan", "varchar", "", "NO").
		AddRow("perusahaan", "email", "varchar", "", "YES").
		AddRow("perusahaan", "password", "varchar", "", "NO"). // sensitif → disaring
		AddRow("perusahaan", "id", "char", "PRI", "NO").       // sensitif → disaring
		AddRow("se", "nama_se", "varchar", "", "NO").
		AddRow("se", "kategori_se", "enum", "", "NO")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	schema, err := svc.GetDatabaseSchema()

	assert.NoError(t, err)
	assert.Contains(t, schema, "DATABASE SCHEMA:")
	assert.Contains(t, schema, "Table: perusahaan")
	assert.Contains(t, schema, "nama_perusahaan")
	assert.Contains(t, schema, "Table: se")
	assert.Contains(t, schema, "nama_se")

	// Kolom sensitif TIDAK boleh muncul
	assert.NotContains(t, schema, "  - password:")
	assert.NotContains(t, schema, "  - id:")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatService_GetDatabaseSchema_DBError(t *testing.T) {
	svc, mock := newTestChatService(t)

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("connection refused"))

	schema, err := svc.GetDatabaseSchema()

	assert.Error(t, err)
	assert.Empty(t, schema)
	assert.Contains(t, err.Error(), "connection refused")
}

func TestChatService_GetDatabaseSchema_EmptyDB(t *testing.T) {
	svc, mock := newTestChatService(t)

	rows := sqlmock.NewRows([]string{
		"TABLE_NAME", "COLUMN_NAME", "DATA_TYPE", "COLUMN_KEY", "IS_NULLABLE",
	})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	schema, err := svc.GetDatabaseSchema()

	assert.NoError(t, err)
	assert.Contains(t, schema, "DATABASE SCHEMA:")
}

func TestChatService_GetDatabaseSchema_TampilkanNotNull_DanPrimaryKey(t *testing.T) {
	svc, mock := newTestChatService(t)

	rows := sqlmock.NewRows([]string{
		"TABLE_NAME", "COLUMN_NAME", "DATA_TYPE", "COLUMN_KEY", "IS_NULLABLE",
	}).
		AddRow("sektor", "nama_sektor", "varchar", "", "NO").
		AddRow("sektor", "kode_sektor", "char", "PRI", "NO")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	schema, err := svc.GetDatabaseSchema()

	assert.NoError(t, err)
	assert.Contains(t, schema, "NOT NULL")
	assert.Contains(t, schema, "PRIMARY KEY")
}

// ============================================================
// TestChatService_ExecuteQuery
// ============================================================

func TestChatService_ExecuteQuery_Success(t *testing.T) {
	svc, mock := newTestChatService(t)

	rows := sqlmock.NewRows([]string{"nama_perusahaan", "email"}).
		AddRow("PT Test Indonesia", "pt@test.com").
		AddRow("CV Maju Bersama", "cv@maju.com")

	mock.ExpectQuery("SELECT nama_perusahaan, email FROM perusahaan").
		WillReturnRows(rows)

	results, err := svc.ExecuteQuery("SELECT nama_perusahaan, email FROM perusahaan")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "PT Test Indonesia", results[0]["nama_perusahaan"])
	assert.Equal(t, "pt@test.com", results[0]["email"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatService_ExecuteQuery_HapusKolomSensitif(t *testing.T) {
	svc, mock := newTestChatService(t)

	rows := sqlmock.NewRows([]string{"nama_perusahaan", "id", "password"}).
		AddRow("PT Aman", "uuid-123", "secret_pass")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := svc.ExecuteQuery("SELECT nama_perusahaan, id, password FROM perusahaan")

	assert.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "PT Aman", results[0]["nama_perusahaan"])

	_, hasID := results[0]["id"]
	_, hasPassword := results[0]["password"]
	assert.False(t, hasID, "kolom 'id' harus disaring")
	assert.False(t, hasPassword, "kolom 'password' harus disaring")
}

func TestChatService_ExecuteQuery_DBError(t *testing.T) {
	svc, mock := newTestChatService(t)

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("tabel tidak ada"))

	results, err := svc.ExecuteQuery("SELECT * FROM tabel_tidak_ada")

	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "query execution failed")
}

func TestChatService_ExecuteQuery_HasilKosong(t *testing.T) {
	svc, mock := newTestChatService(t)

	rows := sqlmock.NewRows([]string{"nama_perusahaan"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := svc.ExecuteQuery("SELECT nama_perusahaan FROM perusahaan WHERE 1=0")

	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Empty(t, results)
}

func TestChatService_ExecuteQuery_NilaiBytesDikonversiKeString(t *testing.T) {
	svc, mock := newTestChatService(t)

	// MySQL sering mengembalikan []byte untuk string columns
	rows := sqlmock.NewRows([]string{"nama_sektor"}).
		AddRow([]byte("Keuangan"))

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := svc.ExecuteQuery("SELECT nama_sektor FROM sektor")

	assert.NoError(t, err)
	require.Len(t, results, 1)
	// Harus berupa string, bukan []byte
	val, ok := results[0]["nama_sektor"].(string)
	assert.True(t, ok, "nilai harus dikonversi ke string")
	assert.Equal(t, "Keuangan", val)
}

// ============================================================
// TestChatService_FormatQueryResults_HasilKosong
// FormatQueryResults mengembalikan pesan statis jika tidak ada data,
// tanpa memanggil Gemini — aman untuk ditest.
// ============================================================

func TestChatService_FormatQueryResults_HasilKosong(t *testing.T) {
	svc, _ := newTestChatService(t)

	result, err := svc.FormatQueryResults("berapa total perusahaan?", []map[string]interface{}{})

	assert.NoError(t, err)
	assert.Equal(t, "Tidak ada data yang ditemukan.", result)
}

// ============================================================
// TestChatService_FilterSensitiveColumns
// Diakses via ExecuteQuery (filter diterapkan sebelum return).
// Test berikut memverifikasi semua kombinasi kolom sensitif.
// ============================================================

func TestChatService_FilterSensitiveColumns_SemuaKolomSensitif(t *testing.T) {
	svc, mock := newTestChatService(t)

	rows := sqlmock.NewRows([]string{"id", "password", "token", "api_key", "secret", "nama"}).
		AddRow("uuid", "pass", "tok", "key", "sec", "Trisna")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := svc.ExecuteQuery("SELECT * FROM users")

	assert.NoError(t, err)
	require.Len(t, results, 1)

	for _, sensitif := range []string{"id", "password", "token", "api_key", "secret"} {
		_, ada := results[0][sensitif]
		assert.False(t, ada, "kolom '%s' seharusnya disaring", sensitif)
	}
	assert.Equal(t, "Trisna", results[0]["nama"])
}

func TestChatService_FilterSensitiveColumns_CaseInsensitive(t *testing.T) {
	svc, mock := newTestChatService(t)

	// Kolom dengan huruf kapital — harus tetap disaring
	rows := sqlmock.NewRows([]string{"PASSWORD", "NAMA_PERUSAHAAN"}).
		AddRow("rahasia", "PT Besar")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := svc.ExecuteQuery("SELECT PASSWORD, NAMA_PERUSAHAAN FROM perusahaan")

	assert.NoError(t, err)
	require.Len(t, results, 1)

	_, adaPassword := results[0]["PASSWORD"]
	assert.False(t, adaPassword, "PASSWORD uppercase harus disaring")
	assert.Equal(t, "PT Besar", results[0]["NAMA_PERUSAHAAN"])
}

func TestChatService_FilterSensitiveColumns_TidakAdaSensitif(t *testing.T) {
	svc, mock := newTestChatService(t)

	rows := sqlmock.NewRows([]string{"nama_sektor", "total"}).
		AddRow("Keuangan", int64(5))

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := svc.ExecuteQuery("SELECT nama_sektor, total FROM sektor")

	assert.NoError(t, err)
	require.Len(t, results, 1)
	// Semua kolom harusnya tetap ada
	assert.Equal(t, "Keuangan", results[0]["nama_sektor"])
}

func TestChatService_FilterSensitiveColumns_MultipleRows(t *testing.T) {
	svc, mock := newTestChatService(t)

	rows := sqlmock.NewRows([]string{"id", "nama_perusahaan", "password"}).
		AddRow("uuid-1", "PT Satu", "pass1").
		AddRow("uuid-2", "PT Dua", "pass2").
		AddRow("uuid-3", "PT Tiga", "pass3")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := svc.ExecuteQuery("SELECT id, nama_perusahaan, password FROM perusahaan")

	assert.NoError(t, err)
	assert.Len(t, results, 3)

	for i, row := range results {
		_, hasID := row["id"]
		_, hasPassword := row["password"]
		assert.False(t, hasID, "row %d: kolom 'id' harus disaring", i)
		assert.False(t, hasPassword, "row %d: kolom 'password' harus disaring", i)
		assert.NotEmpty(t, row["nama_perusahaan"])
	}
}