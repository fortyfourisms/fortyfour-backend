package testhelpers

import (
	"errors"
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

	if _, exists := m.users[user.Username]; !exists {
		return errors.New("user not found")
	}
	m.users[user.Username] = user
	return nil
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
