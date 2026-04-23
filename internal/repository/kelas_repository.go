package repository

import (
	"database/sql"
	"fmt"
	"time"

	"fortyfour-backend/internal/models"
)

type KelasRepository struct {
	db *sql.DB
}

func NewKelasRepository(db *sql.DB) *KelasRepository {
	return &KelasRepository{db: db}
}

var _ KelasRepositoryInterface = (*KelasRepository)(nil)

const kelasSelectColumns = `id, judul, deskripsi, thumbnail, kategori, durasi_jp, penyelenggara, target_peserta, syarat_pendaftaran, informasi_umum, status, created_by, created_at, updated_at`

func (r *KelasRepository) Create(k *models.Kelas) error {
	_, err := r.db.Exec(
		`INSERT INTO kelas (id, judul, deskripsi, thumbnail, kategori, durasi_jp, penyelenggara, target_peserta, syarat_pendaftaran, informasi_umum, status, created_by, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		k.ID, k.Judul, k.Deskripsi, k.Thumbnail, k.Kategori, k.DurasiJP, k.Penyelenggara, k.TargetPeserta, k.SyaratPendaftaran, k.InformasiUmum, k.Status, k.CreatedBy,
	)
	return err
}

func (r *KelasRepository) FindByID(id string) (*models.Kelas, error) {
	row := r.db.QueryRow(
		`SELECT `+kelasSelectColumns+` FROM kelas WHERE id = ?`, id,
	)
	return scanKelas(row)
}

func (r *KelasRepository) FindAll(onlyPublished bool) ([]models.Kelas, error) {
	query := `SELECT ` + kelasSelectColumns + ` FROM kelas`
	if onlyPublished {
		query += ` WHERE status = 'published'`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Kelas
	for rows.Next() {
		k, err := scanKelasRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *k)
	}
	return result, nil
}

func (r *KelasRepository) Update(k *models.Kelas) error {
	_, err := r.db.Exec(
		`UPDATE kelas SET judul=?, deskripsi=?, thumbnail=?, kategori=?, durasi_jp=?, penyelenggara=?, target_peserta=?, syarat_pendaftaran=?, informasi_umum=?, status=?, updated_at=NOW() WHERE id=?`,
		k.Judul, k.Deskripsi, k.Thumbnail, k.Kategori, k.DurasiJP, k.Penyelenggara, k.TargetPeserta, k.SyaratPendaftaran, k.InformasiUmum, k.Status, k.ID,
	)
	return err
}

func (r *KelasRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM kelas WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("kelas dengan id %s tidak ditemukan", id)
	}
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanKelas(row *sql.Row) (*models.Kelas, error) {
	var k models.Kelas
	var deskripsi, thumbnail, kategori, penyelenggara, targetPeserta, syaratPendaftaran, informasiUmum sql.NullString
	var durasiJP sql.NullInt64
	err := row.Scan(
		&k.ID, &k.Judul, &deskripsi, &thumbnail,
		&kategori, &durasiJP, &penyelenggara, &targetPeserta, &syaratPendaftaran, &informasiUmum,
		&k.Status, &k.CreatedBy, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	assignNullString(&k.Deskripsi, deskripsi)
	assignNullString(&k.Thumbnail, thumbnail)
	assignNullString(&k.Kategori, kategori)
	assignNullString(&k.Penyelenggara, penyelenggara)
	assignNullString(&k.TargetPeserta, targetPeserta)
	assignNullString(&k.SyaratPendaftaran, syaratPendaftaran)
	assignNullString(&k.InformasiUmum, informasiUmum)
	assignNullInt(&k.DurasiJP, durasiJP)
	return &k, nil
}

func scanKelasRow(rows *sql.Rows) (*models.Kelas, error) {
	var k models.Kelas
	var deskripsi, thumbnail, kategori, penyelenggara, targetPeserta, syaratPendaftaran, informasiUmum sql.NullString
	var durasiJP sql.NullInt64
	err := rows.Scan(
		&k.ID, &k.Judul, &deskripsi, &thumbnail,
		&kategori, &durasiJP, &penyelenggara, &targetPeserta, &syaratPendaftaran, &informasiUmum,
		&k.Status, &k.CreatedBy, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	assignNullString(&k.Deskripsi, deskripsi)
	assignNullString(&k.Thumbnail, thumbnail)
	assignNullString(&k.Kategori, kategori)
	assignNullString(&k.Penyelenggara, penyelenggara)
	assignNullString(&k.TargetPeserta, targetPeserta)
	assignNullString(&k.SyaratPendaftaran, syaratPendaftaran)
	assignNullString(&k.InformasiUmum, informasiUmum)
	assignNullInt(&k.DurasiJP, durasiJP)
	return &k, nil
}

func assignNullString(target **string, val sql.NullString) {
	if val.Valid {
		*target = &val.String
	}
}

func assignNullInt(target **int, val sql.NullInt64) {
	if val.Valid {
		v := int(val.Int64)
		*target = &v
	}
}

// memastikan time.Time tidak zero ketika dipakai
var _ = time.Now
