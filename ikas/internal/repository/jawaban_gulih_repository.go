package repository

import (
	"database/sql"
	"ikas/internal/dto"
	"ikas/internal/utils"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type JawabanGulihRepositoryInterface interface {
	Create(req dto.CreateJawabanGulihRequest) (int64, error)
	GetAll() ([]dto.JawabanGulihResponse, error)
	GetByID(id int) (*dto.JawabanGulihResponse, error)
	GetByPerusahaan(perusahaanID string) ([]dto.JawabanGulihResponse, error)
	GetByPertanyaan(pertanyaanID int) ([]dto.JawabanGulihResponse, error)
	Update(id int, req dto.UpdateJawabanGulihRequest) error
	Delete(id int) error
	CheckPertanyaanExists(id int) (bool, error)
	CheckPerusahaanExists(id string) (bool, error)
	CheckDuplicate(perusahaanID string, pertanyaanID int, excludeID int) (bool, error)
	RecalculateGulih(perusahaanID string) error
	UpsertToBuffer(req dto.CreateJawabanGulihRequest) error
	GetBufferCount(perusahaanID string) (int, error)
	FlushBuffer(perusahaanID string) error
}

type JawabanGulihRepository struct {
	db *sql.DB
}

func NewJawabanGulihRepository(db *sql.DB) *JawabanGulihRepository {
	return &JawabanGulihRepository{db: db}
}

const jawabanGulihSelectQuery = `
	SELECT 
		jg.id, 
		jg.perusahaan_id, 
		jg.jawaban_gulih, 
		jg.evidence, 
		jg.validasi, 
		jg.keterangan, 
		jg.created_at, 
		jg.updated_at,
		pg.id, 
		pg.pertanyaan_gulih, 
		sk.id, sk.nama_sub_kategori,
		k.id, k.nama_kategori,
		d.id, d.nama_domain
	FROM jawaban_gulih jg
	JOIN pertanyaan_gulih pg ON jg.pertanyaan_gulih_id = pg.id
	JOIN sub_kategori sk ON pg.sub_kategori_id = sk.id
	JOIN kategori k ON sk.kategori_id = k.id
	JOIN domain d ON k.domain_id = d.id`

func scanJawabanGulih(row interface {
	Scan(dest ...any) error
}) (dto.JawabanGulihResponse, error) {
	var item dto.JawabanGulihResponse
	err := row.Scan(
		&item.ID,
		&item.PerusahaanID,
		&item.JawabanGulih,
		&item.Evidence,
		&item.Validasi,
		&item.Keterangan,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.PertanyaanGulih.ID,
		&item.PertanyaanGulih.PertanyaanGulih,
		&item.PertanyaanGulih.SubKategori.ID,
		&item.PertanyaanGulih.SubKategori.NamaSubKategori,
		&item.PertanyaanGulih.SubKategori.Kategori.ID,
		&item.PertanyaanGulih.SubKategori.Kategori.NamaKategori,
		&item.PertanyaanGulih.SubKategori.Kategori.Domain.ID,
		&item.PertanyaanGulih.SubKategori.Kategori.Domain.NamaDomain,
	)
	return item, err
}

func (r *JawabanGulihRepository) Create(req dto.CreateJawabanGulihRequest) (int64, error) {
	query := `INSERT INTO jawaban_gulih 
		(pertanyaan_gulih_id, perusahaan_id, jawaban_gulih, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		req.PertanyaanGulihID,
		req.PerusahaanID,
		req.JawabanGulih,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	if err != nil {
		rollbar.Error(err)
		return 0, err
	}

	return result.LastInsertId()
}

func (r *JawabanGulihRepository) GetAll() ([]dto.JawabanGulihResponse, error) {
	query := jawabanGulihSelectQuery + ` ORDER BY jg.created_at ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var results []dto.JawabanGulihResponse
	for rows.Next() {
		item, err := scanJawabanGulih(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		results = append(results, item)
	}

	return results, nil
}

func (r *JawabanGulihRepository) GetByID(id int) (*dto.JawabanGulihResponse, error) {
	query := jawabanGulihSelectQuery + ` WHERE jg.id = ?`
	row := r.db.QueryRow(query, id)
	item, err := scanJawabanGulih(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		rollbar.Error(err)
		return nil, err
	}
	return &item, nil
}

func (r *JawabanGulihRepository) GetByPerusahaan(perusahaanID string) ([]dto.JawabanGulihResponse, error) {
	query := jawabanGulihSelectQuery + ` WHERE jg.perusahaan_id = ? ORDER BY jg.created_at ASC`
	rows, err := r.db.Query(query, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var results []dto.JawabanGulihResponse
	for rows.Next() {
		item, err := scanJawabanGulih(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		results = append(results, item)
	}

	return results, nil
}

func (r *JawabanGulihRepository) GetByPertanyaan(pertanyaanID int) ([]dto.JawabanGulihResponse, error) {
	query := jawabanGulihSelectQuery + ` WHERE jg.pertanyaan_gulih_id = ? ORDER BY jg.created_at ASC`
	rows, err := r.db.Query(query, pertanyaanID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	defer rows.Close()

	var results []dto.JawabanGulihResponse
	for rows.Next() {
		item, err := scanJawabanGulih(rows)
		if err != nil {
			rollbar.Error(err)
			continue
		}
		results = append(results, item)
	}

	return results, nil
}

func (r *JawabanGulihRepository) Update(id int, req dto.UpdateJawabanGulihRequest) error {
	var updates []string
	var args []interface{}

	if req.JawabanGulih != nil {
		updates = append(updates, "jawaban_gulih = ?")
		args = append(args, req.JawabanGulih)
	}
	if req.Evidence != nil {
		updates = append(updates, "evidence = ?")
		args = append(args, req.Evidence)
	}
	if req.Validasi != nil {
		updates = append(updates, "validasi = ?")
		args = append(args, req.Validasi)
	}
	if req.Keterangan != nil {
		updates = append(updates, "keterangan = ?")
		args = append(args, req.Keterangan)
	}

	if len(updates) == 0 {
		return nil
	}

	query := "UPDATE jawaban_gulih SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		rollbar.Error(err)
	}
	return err
}

func (r *JawabanGulihRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM jawaban_gulih WHERE id = ?`, id)
	if err != nil {
		rollbar.Error(err)
	}
	return err
}

func (r *JawabanGulihRepository) CheckPertanyaanExists(id int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM pertanyaan_gulih WHERE id = ?`, id).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanGulihRepository) CheckPerusahaanExists(id string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM perusahaan WHERE id = ?`, id).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanGulihRepository) CheckDuplicate(perusahaanID string, pertanyaanID int, excludeID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM jawaban_gulih WHERE perusahaan_id = ? AND pertanyaan_gulih_id = ?`
	args := []interface{}{perusahaanID, pertanyaanID}

	if excludeID != 0 {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanGulihRepository) RecalculateGulih(perusahaanID string) error {
	query := `
		SELECT k.id AS kategori_id, ROUND(AVG(jg.jawaban_gulih), 2) AS avg_nilai
		FROM jawaban_gulih jg
		JOIN pertanyaan_gulih pg ON jg.pertanyaan_gulih_id = pg.id
		JOIN sub_kategori sk ON pg.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		WHERE jg.perusahaan_id = ? AND jg.jawaban_gulih IS NOT NULL
		GROUP BY k.id
		ORDER BY k.id`

	rows, err := r.db.Query(query, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	defer rows.Close()

	// Map kategori_id ke nilai subdomain (default 0)
	// kategori 15→sub1, 16→sub2, 17→sub3, 18→sub4
	subdomain := map[int]float64{15: 0, 16: 0, 17: 0, 18: 0}

	for rows.Next() {
		var kategoriID int
		var avgNilai float64
		if err := rows.Scan(&kategoriID, &avgNilai); err != nil {
			rollbar.Error(err)
			continue
		}
		if kategoriID >= 15 && kategoriID <= 18 {
			subdomain[kategoriID] = avgNilai
		}
	}

	// Hitung rata-rata keseluruhan (nilai_gulih)
	nilaiGulih := utils.RoundToTwo((subdomain[15] + subdomain[16] + subdomain[17] + subdomain[18]) / 4.0)

	upsertQuery := `
		INSERT INTO gulih 
			(perusahaan_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			nilai_gulih = VALUES(nilai_gulih),
			nilai_subdomain1 = VALUES(nilai_subdomain1),
			nilai_subdomain2 = VALUES(nilai_subdomain2),
			nilai_subdomain3 = VALUES(nilai_subdomain3),
			nilai_subdomain4 = VALUES(nilai_subdomain4)`

	_, err = r.db.Exec(upsertQuery,
		perusahaanID,
		nilaiGulih,
		subdomain[15], subdomain[16], subdomain[17], subdomain[18],
	)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	var gulihID int
	err = r.db.QueryRow(`SELECT id FROM gulih WHERE perusahaan_id = ?`, perusahaanID).Scan(&gulihID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	_, err = r.db.Exec(`UPDATE ikas SET id_gulih = ? WHERE id_perusahaan = ?`, gulihID, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	updateKematanganQuery := `
		UPDATE ikas i
		LEFT JOIN identifikasi iden ON i.id_identifikasi = iden.id
		LEFT JOIN proteksi prot ON i.id_proteksi = prot.id
		LEFT JOIN deteksi det ON i.id_deteksi = det.id
		LEFT JOIN gulih g ON i.id_gulih = g.id
		SET i.nilai_kematangan = ROUND((
			COALESCE(iden.nilai_identifikasi, 0) + 
			COALESCE(prot.nilai_proteksi, 0) + 
			COALESCE(det.nilai_deteksi, 0) + 
			COALESCE(g.nilai_gulih, 0)
		) / (
			(CASE WHEN iden.id IS NOT NULL THEN 1 ELSE 0 END) +
			(CASE WHEN prot.id IS NOT NULL THEN 1 ELSE 0 END) +
			(CASE WHEN det.id IS NOT NULL THEN 1 ELSE 0 END) +
			(CASE WHEN g.id IS NOT NULL THEN 1 ELSE 0 END)
		), 2)
		WHERE i.id_perusahaan = ? AND (
			iden.id IS NOT NULL OR prot.id IS NOT NULL OR det.id IS NOT NULL OR g.id IS NOT NULL
		)`

	_, err = r.db.Exec(updateKematanganQuery, perusahaanID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *JawabanGulihRepository) UpsertToBuffer(req dto.CreateJawabanGulihRequest) error {
	query := `INSERT INTO jawaban_gulih_buffer 
		(pertanyaan_gulih_id, perusahaan_id, jawaban_gulih, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		jawaban_gulih = VALUES(jawaban_gulih),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	_, err := r.db.Exec(query,
		req.PertanyaanGulihID,
		req.PerusahaanID,
		req.JawabanGulih,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	return err
}

func (r *JawabanGulihRepository) GetBufferCount(perusahaanID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM jawaban_gulih_buffer WHERE perusahaan_id = ?`
	err := r.db.QueryRow(query, perusahaanID).Scan(&count)
	return count, err
}

func (r *JawabanGulihRepository) FlushBuffer(perusahaanID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	moveQuery := `INSERT INTO jawaban_gulih 
		(pertanyaan_gulih_id, perusahaan_id, jawaban_gulih, evidence, validasi, keterangan)
		SELECT pertanyaan_gulih_id, perusahaan_id, jawaban_gulih, evidence, validasi, keterangan
		FROM jawaban_gulih_buffer WHERE perusahaan_id = ?
		ON DUPLICATE KEY UPDATE 
		jawaban_gulih = VALUES(jawaban_gulih),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	if _, err := tx.Exec(moveQuery, perusahaanID); err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM jawaban_gulih_buffer WHERE perusahaan_id = ?`, perusahaanID); err != nil {
		return err
	}

	return tx.Commit()
}
