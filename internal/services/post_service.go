package services

import (
	"errors"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"time"
)

type PostService struct {
	postRepo  *repository.PostRepository
	idCounter int
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (s *PostService) CreatePost(title, content string, authorID int) (*models.Post, error) {
	if title == "" || content == "" {
		return nil, errors.New("title and content are required")
	}

	s.idCounter++
	post := &models.Post{
		ID:        s.idCounter,
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetAllPosts() ([]*models.Post, error) {
	return s.postRepo.FindAll()
}

func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	return s.postRepo.FindByID(id)
}

func (s *PostService) UpdatePost(id, userID int, title, content string) (*models.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != userID {
		return nil, errors.New("unauthorized")
	}

	post.Title = title
	post.Content = content

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(id, userID int) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}

	if post.AuthorID != userID {
		return errors.New("unauthorized")
	}

	return s.postRepo.Delete(id)
}
