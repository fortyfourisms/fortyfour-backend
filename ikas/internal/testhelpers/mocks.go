package testhelpers

import (
	"errors"
	"ikas/internal/dto"
	"ikas/internal/models"
	"sync"
)

// ============================================================
// Mock Identifikasi Repository
// ============================================================

type MockIdentifikasiRepository struct {
	identifikasis map[string]*models.Identifikasi
	mu            sync.RWMutex
}

func NewMockIdentifikasiRepository() *MockIdentifikasiRepository {
	return &MockIdentifikasiRepository{
		identifikasis: make(map[string]*models.Identifikasi),
	}
}

func (m *MockIdentifikasiRepository) Create(req dto.CreateIdentifikasiRequest, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.identifikasis[id] = &models.Identifikasi{
		ID:                id,
		NilaiIdentifikasi: req.NilaiIdentifikasi,
		NilaiSubdomain1:   req.NilaiSubdomain1,
		NilaiSubdomain2:   req.NilaiSubdomain2,
		NilaiSubdomain3:   req.NilaiSubdomain3,
		NilaiSubdomain4:   req.NilaiSubdomain4,
		NilaiSubdomain5:   req.NilaiSubdomain5,
	}
	return nil
}

func (m *MockIdentifikasiRepository) GetAll() ([]models.Identifikasi, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	identifikasis := make([]models.Identifikasi, 0, len(m.identifikasis))
	for _, i := range m.identifikasis {
		identifikasis = append(identifikasis, *i)
	}
	return identifikasis, nil
}

func (m *MockIdentifikasiRepository) GetByID(id string) (*models.Identifikasi, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	identifikasi, exists := m.identifikasis[id]
	if !exists {
		return nil, errors.New("identifikasi not found")
	}
	return identifikasi, nil
}

func (m *MockIdentifikasiRepository) Update(id string, identifikasi models.Identifikasi) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.identifikasis[id]; !exists {
		return errors.New("identifikasi not found")
	}
	m.identifikasis[id] = &identifikasi
	return nil
}

func (m *MockIdentifikasiRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.identifikasis[id]; !exists {
		return errors.New("identifikasi not found")
	}
	delete(m.identifikasis, id)
	return nil
}

// ============================================================
// Mock Deteksi Repository
// ============================================================

type MockDeteksiRepository struct {
	deteksis map[string]*models.Deteksi
	mu       sync.RWMutex
}

func NewMockDeteksiRepository() *MockDeteksiRepository {
	return &MockDeteksiRepository{
		deteksis: make(map[string]*models.Deteksi),
	}
}

func (m *MockDeteksiRepository) Create(req dto.CreateDeteksiRequest, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.deteksis[id] = &models.Deteksi{
		ID:              id,
		NilaiDeteksi:    req.NilaiDeteksi,
		NilaiSubdomain1: req.NilaiSubdomain1,
		NilaiSubdomain2: req.NilaiSubdomain2,
		NilaiSubdomain3: req.NilaiSubdomain3,
	}
	return nil
}

func (m *MockDeteksiRepository) GetAll() ([]models.Deteksi, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	deteksis := make([]models.Deteksi, 0, len(m.deteksis))
	for _, d := range m.deteksis {
		deteksis = append(deteksis, *d)
	}
	return deteksis, nil
}

func (m *MockDeteksiRepository) GetByID(id string) (*models.Deteksi, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	deteksi, exists := m.deteksis[id]
	if !exists {
		return nil, errors.New("deteksi not found")
	}
	return deteksi, nil
}

func (m *MockDeteksiRepository) Update(id string, deteksi models.Deteksi) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.deteksis[id]; !exists {
		return errors.New("deteksi not found")
	}
	m.deteksis[id] = &deteksi
	return nil
}

func (m *MockDeteksiRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.deteksis[id]; !exists {
		return errors.New("deteksi not found")
	}
	delete(m.deteksis, id)
	return nil
}

// ============================================================
// Mock Gulih Repository
// ============================================================

type MockGulihRepository struct {
	gulihs map[string]*models.Gulih
	mu     sync.RWMutex
}

func NewMockGulihRepository() *MockGulihRepository {
	return &MockGulihRepository{
		gulihs: make(map[string]*models.Gulih),
	}
}

func (m *MockGulihRepository) Create(req dto.CreateGulihRequest, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.gulihs[id] = &models.Gulih{
		ID:              id,
		NilaiGulih:      req.NilaiGulih,
		NilaiSubdomain1: req.NilaiSubdomain1,
		NilaiSubdomain2: req.NilaiSubdomain2,
		NilaiSubdomain3: req.NilaiSubdomain3,
		NilaiSubdomain4: req.NilaiSubdomain4,
	}
	return nil
}

func (m *MockGulihRepository) GetAll() ([]models.Gulih, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	gulihs := make([]models.Gulih, 0, len(m.gulihs))
	for _, g := range m.gulihs {
		gulihs = append(gulihs, *g)
	}
	return gulihs, nil
}

func (m *MockGulihRepository) GetByID(id string) (*models.Gulih, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	gulih, exists := m.gulihs[id]
	if !exists {
		return nil, errors.New("gulih not found")
	}
	return gulih, nil
}

func (m *MockGulihRepository) Update(id string, gulih models.Gulih) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.gulihs[id]; !exists {
		return errors.New("gulih not found")
	}
	m.gulihs[id] = &gulih
	return nil
}

func (m *MockGulihRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.gulihs[id]; !exists {
		return errors.New("gulih not found")
	}
	delete(m.gulihs, id)
	return nil
}

// ============================================================
// Mock Proteksi Repository
// ============================================================

type MockProteksiRepository struct {
	proteksis map[string]*models.Proteksi
	mu        sync.RWMutex
}

func NewMockProteksiRepository() *MockProteksiRepository {
	return &MockProteksiRepository{
		proteksis: make(map[string]*models.Proteksi),
	}
}

func (m *MockProteksiRepository) Create(req dto.CreateProteksiRequest, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.proteksis[id] = &models.Proteksi{
		ID:              id,
		NilaiProteksi:   req.NilaiProteksi,
		NilaiSubdomain1: req.NilaiSubdomain1,
		NilaiSubdomain2: req.NilaiSubdomain2,
		NilaiSubdomain3: req.NilaiSubdomain3,
		NilaiSubdomain4: req.NilaiSubdomain4,
		NilaiSubdomain5: req.NilaiSubdomain5,
		NilaiSubdomain6: req.NilaiSubdomain6,
	}
	return nil
}

func (m *MockProteksiRepository) GetAll() ([]models.Proteksi, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	proteksis := make([]models.Proteksi, 0, len(m.proteksis))
	for _, p := range m.proteksis {
		proteksis = append(proteksis, *p)
	}
	return proteksis, nil
}

func (m *MockProteksiRepository) GetByID(id string) (*models.Proteksi, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	proteksi, exists := m.proteksis[id]
	if !exists {
		return nil, errors.New("proteksi not found")
	}
	return proteksi, nil
}

func (m *MockProteksiRepository) Update(id string, proteksi models.Proteksi) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.proteksis[id]; !exists {
		return errors.New("proteksi not found")
	}
	m.proteksis[id] = &proteksi
	return nil
}

func (m *MockProteksiRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.proteksis[id]; !exists {
		return errors.New("proteksi not found")
	}
	delete(m.proteksis, id)
	return nil
}
