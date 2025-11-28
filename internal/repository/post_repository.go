package repository

import (
	"database/sql"
	"errors"
	"fortyfour-backend/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
	query := `INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, post.Title, post.Content, post.AuthorID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.ID = int(id)
	return nil
}

func (r *PostRepository) FindAll() ([]*models.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at 
	          FROM posts ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) FindByID(id int) (*models.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at 
	          FROM posts WHERE id = ?`

	post := &models.Post{}
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("post not found")
	}
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) FindByAuthorID(authorID int) ([]*models.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at 
	          FROM posts WHERE author_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	query := `UPDATE posts SET title = ?, content = ? WHERE id = ?`

	result, err := r.db.Exec(query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("post not found")
	}

	return nil
}

func (r *PostRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("post not found")
	}

	return nil
}
