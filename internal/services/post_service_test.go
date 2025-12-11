package services

import (
	"fortyfour-backend/internal/testhelpers"
	"testing"
)

func TestPostService_CreatePost_Success(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	post, err := service.CreatePost("Test Title", "Test Content", "1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if post == nil {
		t.Fatal("expected post to be created")
	}

	if post.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got '%s'", post.Title)
	}

	if post.Content != "Test Content" {
		t.Errorf("expected content 'Test Content', got '%s'", post.Content)
	}

	if post.AuthorID != "1" {
		t.Errorf("expected author_id 1, got %d", post.AuthorID)
	}

	if post.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestPostService_CreatePost_EmptyFields(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	testCases := []struct {
		name    string
		title   string
		content string
	}{
		{"empty title", "", "content"},
		{"empty content", "title", ""},
		{"both empty", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreatePost(tc.title, tc.content, "1")

			if err == nil {
				t.Fatal("expected error for empty fields")
			}

			if err.Error() != "title and content are required" {
				t.Errorf("expected 'title and content are required', got '%s'", err.Error())
			}
		})
	}
}

func TestPostService_GetAllPosts(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	service.CreatePost("Post 1", "Content 1", "1")
	service.CreatePost("Post 2", "Content 2", "1")
	service.CreatePost("Post 3", "Content 3", "2")

	posts, err := service.GetAllPosts()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(posts) != 3 {
		t.Errorf("expected 3 posts, got %d", len(posts))
	}
}

func TestPostService_GetPostByID_Success(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	created, _ := service.CreatePost("Test Post", "Test Content", "1")

	post, err := service.GetPostByID(created.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if post.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, post.ID)
	}

	if post.Title != "Test Post" {
		t.Errorf("expected title 'Test Post', got '%s'", post.Title)
	}
}

func TestPostService_GetPostByID_NotFound(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	_, err := service.GetPostByID(999)

	if err == nil {
		t.Fatal("expected error for non-existent post")
	}

	if err.Error() != "post not found" {
		t.Errorf("expected 'post not found', got '%s'", err.Error())
	}
}

func TestPostService_UpdatePost_Success(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	created, _ := service.CreatePost("Original Title", "Original Content", "1")

	updated, err := service.UpdatePost(created.ID, "Updated Title", "Updated Content", "1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", updated.Title)
	}

	if updated.Content != "Updated Content" {
		t.Errorf("expected content 'Updated Content', got '%s'", updated.Content)
	}
}

func TestPostService_UpdatePost_Unauthorized(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	created, _ := service.CreatePost("Test Post", "Test Content", "1")

	_, err := service.UpdatePost(created.ID, "Updated Title", "Updated Content", "2")

	if err == nil {
		t.Fatal("expected error for unauthorized update")
	}

	if err.Error() != "unauthorized" {
		t.Errorf("expected 'unauthorized', got '%s'", err.Error())
	}
}

func TestPostService_DeletePost_Success(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	created, _ := service.CreatePost("Test Post", "Test Content", "1")

	err := service.DeletePost(created.ID, "1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.GetPostByID(created.ID)
	if err == nil {
		t.Error("expected post to be deleted")
	}
}

func TestPostService_DeletePost_Unauthorized(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	created, _ := service.CreatePost("Test Post", "Test Content", "1")

	err := service.DeletePost(created.ID, "2")

	if err == nil {
		t.Fatal("expected error for unauthorized delete")
	}

	if err.Error() != "unauthorized" {
		t.Errorf("expected 'unauthorized', got '%s'", err.Error())
	}
}

func TestPostService_DeletePost_NotFound(t *testing.T) {
	repo := testhelpers.NewMockPostRepository()
	service := NewPostService(repo)

	err := service.DeletePost(999, "1")

	if err == nil {
		t.Fatal("expected error for non-existent post")
	}

	if err.Error() != "post not found" {
		t.Errorf("expected 'post not found', got '%s'", err.Error())
	}
}
