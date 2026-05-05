package repository

import (
	"database/sql"
	"fortyfour-backend/internal/models"
)

type BeritaRepository struct {
	db *sql.DB
}

func NewBeritaRepository(db *sql.DB) *BeritaRepository {
	return &BeritaRepository{db: db}
}

var _ BeritaRepositoryInterface = (*BeritaRepository)(nil)

func (r *BeritaRepository) Create(berita *models.Berita) error {
	query := `
		INSERT INTO berita (judul, deskripsi, tags, author_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`
	res, err := r.db.Exec(query, berita.Judul, berita.Deskripsi, berita.Tags, berita.AuthorID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	berita.ID = id
	return nil
}

func (r *BeritaRepository) FindAll() ([]models.Berita, error) {
	query := `
		SELECT b.id, b.judul, b.deskripsi, b.tags, b.author_id, b.created_at, b.updated_at,
		       u.username, u.display_name
		FROM berita b
		LEFT JOIN users u ON b.author_id = u.id
		ORDER BY b.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Berita
	for rows.Next() {
		var b models.Berita
		var author models.User
		var displayName sql.NullString
		var tags sql.NullString

		err := rows.Scan(
			&b.ID, &b.Judul, &b.Deskripsi, &tags, &b.AuthorID, &b.CreatedAt, &b.UpdatedAt,
			&author.Username, &displayName,
		)
		if err != nil {
			return nil, err
		}

		if tags.Valid {
			b.Tags = tags.String
		} else {
			b.Tags = "[]"
		}

		if displayName.Valid {
			tmp := displayName.String
			author.DisplayName = &tmp
		}
		b.Author = &author
		list = append(list, b)
	}
	return list, nil
}

func (r *BeritaRepository) FindByID(id int64) (*models.Berita, error) {
	query := `
		SELECT b.id, b.judul, b.deskripsi, b.tags, b.author_id, b.created_at, b.updated_at,
		       u.username, u.display_name
		FROM berita b
		LEFT JOIN users u ON b.author_id = u.id
		WHERE b.id = ?
	`
	row := r.db.QueryRow(query, id)

	var b models.Berita
	var author models.User
	var displayName sql.NullString
	var tags sql.NullString

	err := row.Scan(
		&b.ID, &b.Judul, &b.Deskripsi, &tags, &b.AuthorID, &b.CreatedAt, &b.UpdatedAt,
		&author.Username, &displayName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if tags.Valid {
		b.Tags = tags.String
	} else {
		b.Tags = "[]"
	}

	if displayName.Valid {
		tmp := displayName.String
		author.DisplayName = &tmp
	}
	b.Author = &author
	return &b, nil
}

func (r *BeritaRepository) Update(berita *models.Berita) error {
	query := `
		UPDATE berita
		SET judul = ?, deskripsi = ?, tags = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(query, berita.Judul, berita.Deskripsi, berita.Tags, berita.ID)
	return err
}

func (r *BeritaRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM berita WHERE id = ?`, id)
	return err
}
