package handlers

import (
	"bytes"
	"encoding/json"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestPostHandler_CreatePost_Success(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := services.NewPostService(repo)
	handler := NewPostHandler(service)

	reqBody := map[string]string{
		"title":   "Test Post",
		"content": "Test Content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/posts/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	handler.CreatePost(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var post models.Post
	json.NewDecoder(w.Body).Decode(&post)

	if post.Title != "Test Post" {
		t.Errorf("expected title 'Test Post', got '%s'", post.Title)
	}
}

func TestPostHandler_GetPosts_Success(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := services.NewPostService(repo)
	handler := NewPostHandler(service)

	service.CreatePost("Post 1", "Content 1", 1)
	service.CreatePost("Post 2", "Content 2", 1)

	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	w := httptest.NewRecorder()

	handler.GetPosts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var posts []*models.Post
	json.NewDecoder(w.Body).Decode(&posts)

	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}
}

func TestPostHandler_UpdatePost_Success(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := services.NewPostService(repo)
	handler := NewPostHandler(service)

	created, _ := service.CreatePost("Original", "Original Content", 1)

	reqBody := map[string]string{
		"title":   "Updated",
		"content": "Updated Content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/posts/update?id="+strconv.Itoa(created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	handler.UpdatePost(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
