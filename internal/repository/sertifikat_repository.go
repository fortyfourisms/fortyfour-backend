package repository

import (
	"database/sql"

	"fortyfour-backend/internal/models"
)

type SertifikatRepository struct {
	db *sql.DB
}

func NewSertifikatRepository(db *sql.DB) *SertifikatRepository {
	return &SertifikatRepository{db: db}
}

var _ SertifikatRepositoryInterface = (*SertifikatRepository)(nil)

func (r *SertifikatRepository) Create(s *models.Sertifikat) error {
	_, err := r.db.Exec(
		`INSERT INTO sertifikat (id, nomor_sertifikat, id_kelas, id_user, nama_peserta, nama_kelas, tanggal_terbit, pdf_path, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		s.ID, s.NomorSertifikat, s.IDKelas, s.IDUser, s.NamaPeserta, s.NamaKelas, s.TanggalTerbit, s.PDFPath,
	)
	return err
}

func (r *SertifikatRepository) FindByUserAndKelas(idUser, idKelas string) (*models.Sertifikat, error) {
	row := r.db.QueryRow(
		`SELECT id, nomor_sertifikat, id_kelas, id_user, nama_peserta, nama_kelas, tanggal_terbit, pdf_path, created_at
		 FROM sertifikat WHERE id_user=? AND id_kelas=?`,
		idUser, idKelas,
	)
	return scanSertifikat(row)
}

func (r *SertifikatRepository) FindByID(id string) (*models.Sertifikat, error) {
	row := r.db.QueryRow(
		`SELECT id, nomor_sertifikat, id_kelas, id_user, nama_peserta, nama_kelas, tanggal_terbit, pdf_path, created_at
		 FROM sertifikat WHERE id=?`, id,
	)
	return scanSertifikat(row)
}

func (r *SertifikatRepository) FindByUser(idUser string) ([]models.Sertifikat, error) {
	rows, err := r.db.Query(
		`SELECT id, nomor_sertifikat, id_kelas, id_user, nama_peserta, nama_kelas, tanggal_terbit, pdf_path, created_at
		 FROM sertifikat WHERE id_user=? ORDER BY tanggal_terbit DESC`, idUser,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Sertifikat
	for rows.Next() {
		var s models.Sertifikat
		var pdfPath sql.NullString
		if err := rows.Scan(
			&s.ID, &s.NomorSertifikat, &s.IDKelas, &s.IDUser,
			&s.NamaPeserta, &s.NamaKelas, &s.TanggalTerbit, &pdfPath, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		if pdfPath.Valid {
			s.PDFPath = &pdfPath.String
		}
		result = append(result, s)
	}
	return result, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanSertifikat(row *sql.Row) (*models.Sertifikat, error) {
	var s models.Sertifikat
	var pdfPath sql.NullString
	err := row.Scan(
		&s.ID, &s.NomorSertifikat, &s.IDKelas, &s.IDUser,
		&s.NamaPeserta, &s.NamaKelas, &s.TanggalTerbit, &pdfPath, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if pdfPath.Valid {
		s.PDFPath = &pdfPath.String
	}
	return &s, nil
}
