package testhelpers

import (
	"errors"
	"fortyfour-backend/internal/models"
	"time"
)

// ============================================================
// Mock User Repository
// ============================================================

type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(user *models.User) error {
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	user, exists := m.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) FindByID(id int) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Update(user *models.User) error {
	if _, exists := m.users[user.Username]; !exists {
		return errors.New("user not found")
	}
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) Delete(id int) error {
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
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts:  make(map[int]*models.Post),
		nextID: 1,
	}
}

func (m *MockPostRepository) Create(post *models.Post) error {
	post.ID = m.nextID
	m.nextID++
	post.CreatedAt = time.Now()
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) FindAll() ([]*models.Post, error) {
	posts := make([]*models.Post, 0, len(m.posts))
	for _, p := range m.posts {
		posts = append(posts, p)
	}
	return posts, nil
}

func (m *MockPostRepository) FindByID(id int) (*models.Post, error) {
	post, exists := m.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (m *MockPostRepository) FindByAuthorID(authorID int) ([]*models.Post, error) {
	posts := make([]*models.Post, 0)
	for _, p := range m.posts {
		if p.AuthorID == authorID {
			posts = append(posts, p)
		}
	}
	return posts, nil
}

func (m *MockPostRepository) Update(post *models.Post) error {
	if _, exists := m.posts[post.ID]; !exists {
		return errors.New("post not found")
	}
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Delete(id int) error {
	if _, exists := m.posts[id]; !exists {
		return errors.New("post not found")
	}
	delete(m.posts, id)
	return nil
}

// ============================================================
// Test Data Factories
// ============================================================

func CreateTestUser(id int, username, email string) *models.User {
	return &models.User{
		ID:       id,
		Username: username,
		Email:    email,
		Password: "hashedpassword",
	}
}

func CreateTestPost(id, authorID int, title, content string) *models.Post {
	return &models.Post{
		ID:        id,
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}
}
