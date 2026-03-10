package services

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/stretchr/testify/assert"
)

//
// ===============================
// MOCK PIC REPOSITORY
// ===============================
//

type mockPICRepository struct {
	CreateFn           func(req dto.CreatePICRequest, id string) error
	GetByIDFn          func(id string) (*dto.PICResponse, error)
	GetAllFn           func() ([]dto.PICResponse, error)
	GetByPerusahaanFn  func(idPerusahaan string) ([]dto.PICResponse, error)
	UpdateFn           func(id string, req dto.UpdatePICRequest) error
	DeleteFn           func(id string) error
}

func (m *mockPICRepository) Create(req dto.CreatePICRequest, id string) error {
	return m.CreateFn(req, id)
}

func (m *mockPICRepository) GetByID(id string) (*dto.PICResponse, error) {
	return m.GetByIDFn(id)
}

func (m *mockPICRepository) GetAll() ([]dto.PICResponse, error) {
	return m.GetAllFn()
}

func (m *mockPICRepository) Update(id string, req dto.UpdatePICRequest) error {
	return m.UpdateFn(id, req)
}

func (m *mockPICRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

func (m *mockPICRepository) GetByPerusahaan(idPerusahaan string) ([]dto.PICResponse, error) {
	if m.GetByPerusahaanFn != nil {
		return m.GetByPerusahaanFn(idPerusahaan)
	}
	return []dto.PICResponse{}, nil
}

// Compile-time check
var _ repository.PICRepositoryInterface = (*mockPICRepository)(nil)

//
// ===============================
// TEST CREATE
// ===============================
//

func TestPICService_Create_Success(t *testing.T) {
	nama := "John Doe"
	idPerusahaan := "uuid-perusahaan"

	repo := &mockPICRepository{
		CreateFn: func(req dto.CreatePICRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return &dto.PICResponse{
				ID:   id,
				Nama: nama,
				Perusahaan: &dto.PerusahaanInPIC{
					ID:             idPerusahaan,
					NamaPerusahaan: "PT Contoh",
				},
			}, nil
		},
	}

	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{
		Nama:         &nama,
		IDPerusahaan: &idPerusahaan,
	}

	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, nama, result.Nama)
	assert.NotNil(t, result.Perusahaan)
	assert.Equal(t, idPerusahaan, result.Perusahaan.ID)
}

func TestPICService_Create_ValidationError(t *testing.T) {
	repo := &mockPICRepository{}
	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST GET ALL
// ===============================
//

func TestPICService_GetAll_Success(t *testing.T) {
	repo := &mockPICRepository{
		GetAllFn: func() ([]dto.PICResponse, error) {
			return []dto.PICResponse{
				{ID: "1", Nama: "PIC 1"},
				{ID: "2", Nama: "PIC 2"},
			}, nil
		},
	}

	service := NewPICService(repo, nil)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

//
// ===============================
// TEST GET BY ID
// ===============================
//

func TestPICService_GetByID_Success(t *testing.T) {
	repo := &mockPICRepository{
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return &dto.PICResponse{
				ID:   id,
				Nama: "PIC Test",
				Perusahaan: &dto.PerusahaanInPIC{
					ID:             "uuid-perusahaan",
					NamaPerusahaan: "PT Test",
				},
			}, nil
		},
	}

	service := NewPICService(repo, nil)

	result, err := service.GetByID("uuid-test")

	assert.NoError(t, err)
	assert.Equal(t, "PIC Test", result.Nama)
	assert.NotNil(t, result.Perusahaan)
}

func TestPICService_GetByID_NotFound(t *testing.T) {
	repo := &mockPICRepository{
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}

	service := NewPICService(repo, nil)

	result, err := service.GetByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST UPDATE
// ===============================
//

func TestPICService_Update_Success(t *testing.T) {
	namaBaru := "Nama Baru"

	repo := &mockPICRepository{
		UpdateFn: func(id string, req dto.UpdatePICRequest) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return &dto.PICResponse{
				ID:   id,
				Nama: namaBaru,
			}, nil
		},
	}

	service := NewPICService(repo, nil)

	req := dto.UpdatePICRequest{
		Nama: &namaBaru,
	}

	result, err := service.Update("uuid-test", req)

	assert.NoError(t, err)
	assert.Equal(t, namaBaru, result.Nama)
}

//
// ===============================
// TEST DELETE
// ===============================
//

func TestPICService_Delete_Success(t *testing.T) {
	repo := &mockPICRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	service := NewPICService(repo, nil)

	err := service.Delete("uuid-test")

	assert.NoError(t, err)
}

//
// ===============================
// TEST CREATE — ERROR CASES
// ===============================
//

func TestPICService_Create_NamaNil_Error(t *testing.T) {
	idPerusahaan := "uuid-perusahaan"
	repo := &mockPICRepository{}
	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{
		IDPerusahaan: &idPerusahaan,
		// Nama: nil
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.EqualError(t, err, "nama wajib diisi")
	assert.Nil(t, result)
}

func TestPICService_Create_NamaKosong_Error(t *testing.T) {
	nama := "   "
	idPerusahaan := "uuid-perusahaan"
	repo := &mockPICRepository{}
	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{
		Nama:         &nama,
		IDPerusahaan: &idPerusahaan,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.EqualError(t, err, "nama wajib diisi")
	assert.Nil(t, result)
}

func TestPICService_Create_IDPerusahaanNil_Error(t *testing.T) {
	nama := "John Doe"
	repo := &mockPICRepository{}
	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{
		Nama: &nama,
		// IDPerusahaan: nil
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.EqualError(t, err, "id_perusahaan wajib diisi")
	assert.Nil(t, result)
}

func TestPICService_Create_IDPerusahaanKosong_Error(t *testing.T) {
	nama := "John Doe"
	idPerusahaan := ""
	repo := &mockPICRepository{}
	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{
		Nama:         &nama,
		IDPerusahaan: &idPerusahaan,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.EqualError(t, err, "id_perusahaan wajib diisi")
	assert.Nil(t, result)
}

func TestPICService_Create_RepoError(t *testing.T) {
	nama := "John Doe"
	idPerusahaan := "uuid-perusahaan"

	repo := &mockPICRepository{
		CreateFn: func(req dto.CreatePICRequest, id string) error {
			return errors.New("db error")
		},
	}
	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{
		Nama:         &nama,
		IDPerusahaan: &idPerusahaan,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
	assert.Nil(t, result)
}

func TestPICService_Create_GetByIDError(t *testing.T) {
	nama := "John Doe"
	idPerusahaan := "uuid-perusahaan"

	repo := &mockPICRepository{
		CreateFn: func(req dto.CreatePICRequest, id string) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return nil, errors.New("data tidak ditemukan setelah create")
		},
	}
	service := NewPICService(repo, nil)

	req := dto.CreatePICRequest{
		Nama:         &nama,
		IDPerusahaan: &idPerusahaan,
	}

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

//
// ===============================
// TEST GET ALL — ERROR & CACHE
// ===============================
//

func TestPICService_GetAll_Error(t *testing.T) {
	repo := &mockPICRepository{
		GetAllFn: func() ([]dto.PICResponse, error) {
			return nil, errors.New("db error")
		},
	}
	service := NewPICService(repo, nil)

	result, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPICService_GetAll_CacheHit_SkipRepo(t *testing.T) {
	// Siapkan Redis mock dengan data cache yang sudah ada
	rc := newMockRedisForPIC()
	expected := []dto.PICResponse{
		{ID: "cache-1", Nama: "PIC dari Cache"},
	}
	setCache(rc, keyList("pic"), expected)

	// Repo tidak dipanggil sama sekali jika cache hit
	repoCalled := false
	repo := &mockPICRepository{
		GetAllFn: func() ([]dto.PICResponse, error) {
			repoCalled = true
			return nil, errors.New("seharusnya tidak dipanggil")
		},
	}
	service := NewPICService(repo, rc)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.False(t, repoCalled, "repo tidak boleh dipanggil saat cache hit")
	assert.Len(t, result, 1)
	assert.Equal(t, "PIC dari Cache", result[0].Nama)
}

func TestPICService_GetAll_CacheMiss_HitRepo_ThenSetCache(t *testing.T) {
	rc := newMockRedisForPIC()

	repo := &mockPICRepository{
		GetAllFn: func() ([]dto.PICResponse, error) {
			return []dto.PICResponse{
				{ID: "db-1", Nama: "PIC dari DB"},
				{ID: "db-2", Nama: "PIC dari DB 2"},
			}, nil
		},
	}
	service := NewPICService(repo, rc)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// Setelah GetAll, cache harus ter-set
	exists, _ := rc.Exists(keyList("pic"))
	assert.True(t, exists, "data harus di-cache setelah GetAll")
}

//
// ===============================
// TEST GET BY ID — ERROR & CACHE
// ===============================
//

func TestPICService_GetByID_CacheHit_SkipRepo(t *testing.T) {
	rc := newMockRedisForPIC()
	expected := dto.PICResponse{ID: "pic-1", Nama: "PIC Cache"}
	setCache(rc, keyDetail("pic", "pic-1"), expected)

	repoCalled := false
	repo := &mockPICRepository{
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			repoCalled = true
			return nil, errors.New("seharusnya tidak dipanggil")
		},
	}
	service := NewPICService(repo, rc)

	result, err := service.GetByID("pic-1")

	assert.NoError(t, err)
	assert.False(t, repoCalled, "repo tidak boleh dipanggil saat cache hit")
	assert.Equal(t, "PIC Cache", result.Nama)
}

func TestPICService_GetByID_CacheMiss_HitRepo_ThenSetCache(t *testing.T) {
	rc := newMockRedisForPIC()

	repo := &mockPICRepository{
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return &dto.PICResponse{ID: id, Nama: "PIC dari DB"}, nil
		},
	}
	service := NewPICService(repo, rc)

	result, err := service.GetByID("pic-1")

	assert.NoError(t, err)
	assert.Equal(t, "PIC dari DB", result.Nama)

	// Cache harus ter-set setelah GetByID
	exists, _ := rc.Exists(keyDetail("pic", "pic-1"))
	assert.True(t, exists, "data harus di-cache setelah GetByID")
}

//
// ===============================
// TEST UPDATE — ERROR CASES
// ===============================
//

func TestPICService_Update_RepoError(t *testing.T) {
	namaBaru := "Nama Baru"

	repo := &mockPICRepository{
		UpdateFn: func(id string, req dto.UpdatePICRequest) error {
			return errors.New("update failed")
		},
	}
	service := NewPICService(repo, nil)

	result, err := service.Update("uuid-test", dto.UpdatePICRequest{Nama: &namaBaru})

	assert.Error(t, err)
	assert.EqualError(t, err, "update failed")
	assert.Nil(t, result)
}

func TestPICService_Update_GetByIDError(t *testing.T) {
	namaBaru := "Nama Baru"

	repo := &mockPICRepository{
		UpdateFn: func(id string, req dto.UpdatePICRequest) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return nil, errors.New("data tidak ditemukan")
		},
	}
	service := NewPICService(repo, nil)

	result, err := service.Update("uuid-test", dto.UpdatePICRequest{Nama: &namaBaru})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPICService_Update_InvalidatesCache(t *testing.T) {
	rc := newMockRedisForPIC()
	// Pre-populate cache
	setCache(rc, keyDetail("pic", "uuid-test"), dto.PICResponse{ID: "uuid-test", Nama: "Lama"})
	setCache(rc, keyList("pic"), []dto.PICResponse{{ID: "uuid-test"}})

	namaBaru := "Nama Baru"
	repo := &mockPICRepository{
		UpdateFn: func(id string, req dto.UpdatePICRequest) error {
			return nil
		},
		GetByIDFn: func(id string) (*dto.PICResponse, error) {
			return &dto.PICResponse{ID: id, Nama: namaBaru}, nil
		},
	}
	service := NewPICService(repo, rc)

	_, err := service.Update("uuid-test", dto.UpdatePICRequest{Nama: &namaBaru})

	assert.NoError(t, err)

	// Cache detail dan list harus sudah dihapus
	existsDetail, _ := rc.Exists(keyDetail("pic", "uuid-test"))
	existsList, _ := rc.Exists(keyList("pic"))
	assert.False(t, existsDetail, "cache detail harus dihapus setelah update")
	assert.False(t, existsList, "cache list harus dihapus setelah update")
}

//
// ===============================
// TEST DELETE — ERROR & CACHE
// ===============================
//

func TestPICService_Delete_Error(t *testing.T) {
	repo := &mockPICRepository{
		DeleteFn: func(id string) error {
			return errors.New("delete failed")
		},
	}
	service := NewPICService(repo, nil)

	err := service.Delete("uuid-test")

	assert.Error(t, err)
	assert.EqualError(t, err, "delete failed")
}

func TestPICService_Delete_InvalidatesCache(t *testing.T) {
	rc := newMockRedisForPIC()
	setCache(rc, keyDetail("pic", "uuid-test"), dto.PICResponse{ID: "uuid-test"})
	setCache(rc, keyList("pic"), []dto.PICResponse{{ID: "uuid-test"}})

	repo := &mockPICRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}
	service := NewPICService(repo, rc)

	err := service.Delete("uuid-test")

	assert.NoError(t, err)

	existsDetail, _ := rc.Exists(keyDetail("pic", "uuid-test"))
	existsList, _ := rc.Exists(keyList("pic"))
	assert.False(t, existsDetail, "cache detail harus dihapus setelah delete")
	assert.False(t, existsList, "cache list harus dihapus setelah delete")
}

//
// ===============================
// HELPERS UNTUK TEST CACHE
// ===============================
//

// newMockRedisForPIC membuat mock Redis client baru untuk test PIC
func newMockRedisForPIC() *testRedisClient {
	return &testRedisClient{data: make(map[string]string)}
}

// setCache menyimpan data JSON ke mock Redis
func setCache(rc *testRedisClient, key string, value interface{}) {
	b, _ := json.Marshal(value)
	rc.data[key] = string(b)
}

// testRedisClient adalah minimal Redis mock yang implement cache.RedisInterface
type testRedisClient struct {
	data map[string]string
}

func (r *testRedisClient) Set(key string, value interface{}, ttl time.Duration) error {
	if v, ok := value.(string); ok {
		r.data[key] = v
	}
	return nil
}
func (r *testRedisClient) Get(key string) (string, error) {
	v, ok := r.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}
func (r *testRedisClient) Delete(key string) error {
	delete(r.data, key)
	return nil
}
func (r *testRedisClient) Exists(key string) (bool, error) {
	_, ok := r.data[key]
	return ok, nil
}
func (r *testRedisClient) Scan(pattern string) ([]string, error) { return nil, nil }
func (r *testRedisClient) Close() error                          { return nil }
