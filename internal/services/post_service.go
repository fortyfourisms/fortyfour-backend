package services

import (
	"errors"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"time"
)

type PostService struct {
	postRepo repository.PostRepositoryInterface
}

func NewPostService(postRepo repository.PostRepositoryInterface) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (s *PostService) CreatePost(title, content, authorID string) (*models.Post, error) {
	if title == "" || content == "" {
		return nil, errors.New("title and content are required")
	}

	post := &models.Post{
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

func (s *PostService) UpdatePost(id int, title, content, userID string) (*models.Post, error) {
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

func (s *PostService) DeletePost(id int, userID string) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}

	if post.AuthorID != userID {
		return errors.New("unauthorized")
	}

	return s.postRepo.Delete(id)
}
