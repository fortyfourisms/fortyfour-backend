package testhelpers

import (
	"context"
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/pkg/cache"
	"sync"
	"time"
)

// ============================================================
// Mock Redis Client for Testing
// ============================================================

type MockRedisClient struct {
	data map[string]string
	mu   sync.RWMutex
}

// Ensure MockRedisClient implements cache.RedisInterface
var _ cache.RedisInterface = (*MockRedisClient)(nil)

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

func (m *MockRedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Convert value to string
	strValue, ok := value.(string)
	if !ok {
		// Handle []byte conversion
		if bytes, ok := value.([]byte); ok {
			strValue = string(bytes)
		} else {
			return errors.New("invalid value type")
		}
	}

	m.data[key] = strValue
	return nil
}

func (m *MockRedisClient) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.data[key]
	if !exists {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (m *MockRedisClient) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return nil
}

func (m *MockRedisClient) Exists(key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.data[key]
	return exists, nil
}

func (m *MockRedisClient) Close() error {
	return nil
}

// Helper to clear all data (useful for test cleanup)
func (m *MockRedisClient) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]string)
}

// ============================================================
// Mock User Repository
// ============================================================

type MockUserRepository struct {
	users map[string]*models.User
	mu    sync.RWMutex
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Update(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for uname, u := range m.users {
		if u.ID == user.ID {
			if uname != user.Username {
				delete(m.users, uname)
			}
			m.users[user.Username] = user
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for username, user := range m.users {
		if user.ID == id {
			delete(m.users, username)
			return nil
		}
	}
	return errors.New("user not found")
}

// SetMFA - mock implementation to satisfy interface and tests
func (m *MockUserRepository) SetMFA(userID string, secret *string, enabled bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.users {
		if u.ID == userID {
			u.MFAEnabled = enabled
			u.MFASecret = secret
			return nil
		}
	}
	return errors.New("user not found")
}

// ============================================================
// Mock Post Repository
// ============================================================

type MockPostRepository struct {
	posts  map[int]*models.Post
	nextID int
	mu     sync.RWMutex
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts:  make(map[int]*models.Post),
		nextID: 1,
	}
}

func (m *MockPostRepository) Create(post *models.Post) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	post.ID = m.nextID
	m.nextID++
	post.CreatedAt = time.Now()
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) FindAll() ([]*models.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	posts := make([]*models.Post, 0, len(m.posts))
	for _, p := range m.posts {
		posts = append(posts, p)
	}
	return posts, nil
}

func (m *MockPostRepository) FindByID(id int) (*models.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	post, exists := m.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (m *MockPostRepository) FindByAuthorID(authorID string) ([]*models.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	posts := make([]*models.Post, 0)
	for _, p := range m.posts {
		if p.AuthorID == authorID {
			posts = append(posts, p)
		}
	}
	return posts, nil
}

func (m *MockPostRepository) Update(post *models.Post) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.posts[post.ID]; !exists {
		return errors.New("post not found")
	}
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Delete(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.posts[id]; !exists {
		return errors.New("post not found")
	}
	delete(m.posts, id)
	return nil
}

func (m *MockUserRepository) EmailExists(email string, excludeID *string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Email == email {
			if excludeID != nil && user.ID == *excludeID {
				continue
			}
			return true, nil
		}
	}
	return false, nil
}

func (m *MockUserRepository) UsernameExists(username string, excludeID *string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[username]
	if !exists {
		return false, nil
	}

	if excludeID != nil && user.ID == *excludeID {
		return false, nil
	}

	return true, nil
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

func (m *MockUserRepository) GetPasswordByID(id string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.ID == id {
			return user.Password, nil
		}
	}
	return "", errors.New("user not found")
}

func (m *MockUserRepository) UpdateWithPhoto(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find user by ID
	for _, u := range m.users {
		if u.ID == user.ID {
			// Update all fields
			u.Username = user.Username
			u.Email = user.Email
			u.RoleID = user.RoleID
			u.IDJabatan = user.IDJabatan
			u.FotoProfile = user.FotoProfile
			u.Banner = user.Banner
			u.UpdatedAt = user.UpdatedAt
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) UpdatePassword(id, hashedPassword string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, user := range m.users {
		if user.ID == id {
			user.Password = hashedPassword
			return nil
		}
	}
	return errors.New("user not found")
}

// ============================================================
// Test Data Factories
// ============================================================

func CreateTestUser(id, username, email string) *models.User {
	return &models.User{
		ID:       id,
		Username: username,
		Email:    email,
		Password: "hashedpassword",
	}
}

func CreateTestPost(id int, authorID string, title, content string) *models.Post {
	return &models.Post{
		ID:        id,
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}
}

// ============================================================
// Mock SSE Service
// ============================================================

type MockSSEService struct{}

func NewMockSSEService() *MockSSEService {
	return &MockSSEService{}
}

func (m *MockSSEService) NotifyCreate(resource string, data interface{}, userID string) {}
func (m *MockSSEService) NotifyUpdate(resource string, data interface{}, userID string) {}
func (m *MockSSEService) NotifyDelete(resource string, id interface{}, userID string)   {}

// ============================================================
// Mock Role Repository
// ============================================================

type MockRoleRepository struct {
	roles map[string]*models.Role
	mu    sync.RWMutex
}

func NewMockRoleRepository() *MockRoleRepository {
	return &MockRoleRepository{
		roles: make(map[string]*models.Role),
	}
}

func (m *MockRoleRepository) Create(ctx context.Context, role *models.Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.roles[role.ID] = role
	return nil
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id string) (*models.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	role, exists := m.roles[id]
	if !exists {
		return nil, errors.New("role not found")
	}
	return role, nil
}

func (m *MockRoleRepository) GetAll(ctx context.Context) ([]*models.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	roles := make([]*models.Role, 0, len(m.roles))
	for _, role := range m.roles {
		roles = append(roles, role)
	}
	return roles, nil
}

func (m *MockRoleRepository) Update(ctx context.Context, role *models.Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.roles[role.ID]; !exists {
		return errors.New("role not found")
	}
	m.roles[role.ID] = role
	return nil
}

func (m *MockRoleRepository) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.roles[id]; !exists {
		return errors.New("role not found")
	}
	delete(m.roles, id)
	return nil
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, role := range m.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, nil // Return nil, nil if not found (not an error)
}

// ============================================================
// Mock Jabatan Repository
// ============================================================

type MockJabatanRepository struct {
	jabatans map[string]*dto.JabatanResponse
	mu       sync.RWMutex
}

func NewMockJabatanRepository() *MockJabatanRepository {
	return &MockJabatanRepository{
		jabatans: make(map[string]*dto.JabatanResponse),
	}
}

func (m *MockJabatanRepository) Create(req dto.CreateJabatanRequest, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	nama := ""
	if req.NamaJabatan != nil {
		nama = *req.NamaJabatan
	}

	m.jabatans[id] = &dto.JabatanResponse{
		ID:          id,
		NamaJabatan: nama,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	return nil
}

func (m *MockJabatanRepository) GetAll() ([]dto.JabatanResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jabatans := make([]dto.JabatanResponse, 0, len(m.jabatans))
	for _, j := range m.jabatans {
		jabatans = append(jabatans, *j)
	}
	return jabatans, nil
}

func (m *MockJabatanRepository) GetByID(id string) (*dto.JabatanResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jabatan, exists := m.jabatans[id]
	if !exists {
		return nil, errors.New("jabatan not found")
	}
	return jabatan, nil
}

func (m *MockJabatanRepository) Update(id string, jabatan dto.JabatanResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.jabatans[id]; !exists {
		return errors.New("jabatan not found")
	}
	m.jabatans[id] = &jabatan
	return nil
}

func (m *MockJabatanRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.jabatans[id]; !exists {
		return errors.New("jabatan not found")
	}
	delete(m.jabatans, id)
	return nil
}

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

// ============================================================
// Mock PIC Repository
// ============================================================

type MockPICRepository struct {
	pics map[string]*dto.PICResponse
	mu   sync.RWMutex
}

func NewMockPICRepository() *MockPICRepository {
	return &MockPICRepository{
		pics: make(map[string]*dto.PICResponse),
	}
}

func (m *MockPICRepository) Create(req dto.CreatePICRequest, id string) error {
m.mu.Lock()
defer m.mu.Unlock()

nama := ""
if req.Nama != nil {
nama = *req.Nama
}
telepon := ""
if req.Telepon != nil {
telepon = *req.Telepon
}

m.pics[id] = &dto.PICResponse{
ID:        id,
Nama:      nama,
Telepon:   telepon,
CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
}
return nil
}

func (m *MockPICRepository) GetAll() ([]dto.PICResponse, error) {
m.mu.RLock()
defer m.mu.RUnlock()

pics := make([]dto.PICResponse, 0, len(m.pics))
for _, p := range m.pics {
pics = append(pics, *p)
}
return pics, nil
}

func (m *MockPICRepository) GetByID(id string) (*dto.PICResponse, error) {
m.mu.RLock()
defer m.mu.RUnlock()

pic, exists := m.pics[id]
if !exists {
return nil, errors.New("pic not found")
}
return pic, nil
}

func (m *MockPICRepository) Update(id string, req dto.UpdatePICRequest) error {
m.mu.Lock()
defer m.mu.Unlock()

if _, exists := m.pics[id]; !exists {
return errors.New("pic not found")
}
if req.Nama != nil {
m.pics[id].Nama = *req.Nama
}
if req.Telepon != nil {
m.pics[id].Telepon = *req.Telepon
}
m.pics[id].UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
return nil
}

func (m *MockPICRepository) Delete(id string) error {
m.mu.Lock()
defer m.mu.Unlock()

if _, exists := m.pics[id]; !exists {
return errors.New("pic not found")
}
delete(m.pics, id)
return nil
}

// ============================================================
// Mock Perusahaan Repository
// ============================================================

type MockPerusahaanRepository struct {
	perusahaans map[string]*dto.PerusahaanResponse
	mu          sync.RWMutex
}

func NewMockPerusahaanRepository() *MockPerusahaanRepository {
return &MockPerusahaanRepository{
perusahaans: make(map[string]*dto.PerusahaanResponse),
}
}

func (m *MockPerusahaanRepository) Create(req dto.CreatePerusahaanRequest, id string) error {
m.mu.Lock()
defer m.mu.Unlock()

nama := ""
if req.NamaPerusahaan != nil {
nama = *req.NamaPerusahaan
}
photo := ""
if req.Photo != nil {
photo = *req.Photo
}

m.perusahaans[id] = &dto.PerusahaanResponse{
ID:             id,
Photo:          photo,
NamaPerusahaan: nama,
CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
UpdatedAt:      time.Now().Format("2006-01-02 15:04:05"),
}
return nil
}

func (m *MockPerusahaanRepository) GetAll() ([]dto.PerusahaanResponse, error) {
m.mu.RLock()
defer m.mu.RUnlock()

perusahaans := make([]dto.PerusahaanResponse, 0, len(m.perusahaans))
for _, p := range m.perusahaans {
perusahaans = append(perusahaans, *p)
}
return perusahaans, nil
}

func (m *MockPerusahaanRepository) GetByID(id string) (*dto.PerusahaanResponse, error) {
m.mu.RLock()
defer m.mu.RUnlock()

perusahaan, exists := m.perusahaans[id]
if !exists {
return nil, errors.New("perusahaan not found")
}
return perusahaan, nil
}

func (m *MockPerusahaanRepository) Update(id string, perusahaan dto.PerusahaanResponse) error {
m.mu.Lock()
defer m.mu.Unlock()

if _, exists := m.perusahaans[id]; !exists {
return errors.New("perusahaan not found")
}
m.perusahaans[id] = &perusahaan
return nil
}

func (m *MockPerusahaanRepository) Delete(id string) error {
m.mu.Lock()
defer m.mu.Unlock()

if _, exists := m.perusahaans[id]; !exists {
return errors.New("perusahaan not found")
}
delete(m.perusahaans, id)
return nil
}
