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
	GetByUUID(uuid string) (*dto.JawabanGulihResponse, error)
	GetByIkasID(ikasID string) ([]dto.JawabanGulihResponse, error)
	GetByIkasIDFromBuffer(ikasID string) ([]dto.JawabanGulihResponse, error)
	GetByPerusahaanID(perusahaanID string) ([]dto.JawabanGulihResponse, error)
	GetByPertanyaan(pertanyaanID int) ([]dto.JawabanGulihResponse, error)
	GetByPertanyaanAndPerusahaan(pertanyaanID int, perusahaanID string) ([]dto.JawabanGulihResponse, error)
	Update(uuid string, req dto.UpdateJawabanGulihRequest) error
	Delete(uuid string) error
	CheckPertanyaanExists(pertanyaanID int) (bool, error)
	CheckIkasExists(ikasID string) (bool, error)
	CheckDuplicate(ikasID string, pertanyaanID int, excludeID int) (bool, error)
	RecalculateGulih(ikasID string) error
	UpsertToBuffer(req dto.CreateJawabanGulihRequest) error
	GetBufferCount(ikasID string) (int, error)
	FlushBuffer(ikasID string) error
	CloneByIkasID(sourceID, targetID string) error
	GetIDByIkasAndPertanyaan(ikasID string, pertanyaanID int) (int, error)
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
		jg.uuid,
		jg.ikas_id, 
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

const jawabanGulihBufferSelectQuery = `
	SELECT 
		jg.id, 
		jg.uuid,
		jg.ikas_id, 
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
	FROM jawaban_gulih_buffer jg
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
		&item.UUID,
		&item.IkasID,
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
		(uuid, pertanyaan_gulih_id, ikas_id, jawaban_gulih, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		utils.GenerateUUID(),
		req.PertanyaanGulihID,
		req.IkasID,
		req.JawabanGulih,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	if err != nil {
		rollbar.Error(err)
		return 0, err
	}

	return res.LastInsertId()
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

func (r *JawabanGulihRepository) GetByUUID(uuid string) (*dto.JawabanGulihResponse, error) {
	query := jawabanGulihSelectQuery + ` WHERE jg.uuid = ?`
	row := r.db.QueryRow(query, uuid)
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

func (r *JawabanGulihRepository) GetByIkasID(ikasID string) ([]dto.JawabanGulihResponse, error) {
	query := jawabanGulihSelectQuery + ` WHERE jg.ikas_id = ? ORDER BY jg.created_at ASC`

	rows, err := r.db.Query(query, ikasID)
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

func (r *JawabanGulihRepository) GetByIkasIDFromBuffer(ikasID string) ([]dto.JawabanGulihResponse, error) {
	query := jawabanGulihBufferSelectQuery + ` WHERE jg.ikas_id = ? ORDER BY jg.created_at ASC`

	rows, err := r.db.Query(query, ikasID)
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

func (r *JawabanGulihRepository) GetByPerusahaanID(perusahaanID string) ([]dto.JawabanGulihResponse, error) {
	query := jawabanGulihSelectQuery + ` 
		JOIN ikas i ON jg.ikas_id = i.id 
		WHERE i.id_perusahaan = ? 
		ORDER BY jg.created_at ASC`

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

// GetByPertanyaanAndPerusahaan returns answers filtered by both question ID and company ID.
// This enforces multi-tenancy at the database level for non-admin users.
func (r *JawabanGulihRepository) GetByPertanyaanAndPerusahaan(pertanyaanID int, perusahaanID string) ([]dto.JawabanGulihResponse, error) {
	query := jawabanGulihSelectQuery + `
		JOIN ikas i ON jg.ikas_id = i.id
		WHERE jg.pertanyaan_gulih_id = ? AND i.id_perusahaan = ?
		ORDER BY jg.created_at ASC`

	rows, err := r.db.Query(query, pertanyaanID, perusahaanID)
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

func (r *JawabanGulihRepository) Update(uuid string, req dto.UpdateJawabanGulihRequest) error {
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

	query := "UPDATE jawaban_gulih SET " + strings.Join(updates, ", ") + " WHERE uuid = ?"
	args = append(args, uuid)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		rollbar.Error(err)
	}
	return err
}

func (r *JawabanGulihRepository) Delete(uuid string) error {
	_, err := r.db.Exec(`DELETE FROM jawaban_gulih WHERE uuid = ?`, uuid)
	if err != nil {
		rollbar.Error(err)
	}
	return err
}

func (r *JawabanGulihRepository) CheckPertanyaanExists(pertanyaanID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM pertanyaan_gulih WHERE id = ?`, pertanyaanID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanGulihRepository) CheckIkasExists(ikasID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM ikas WHERE id = ?`, ikasID).Scan(&count)
	if err != nil {
		rollbar.Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *JawabanGulihRepository) CheckDuplicate(ikasID string, pertanyaanID int, excludeID int) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM jawaban_gulih
		WHERE ikas_id = ? AND pertanyaan_gulih_id = ?`
	args := []interface{}{ikasID, pertanyaanID}

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

func (r *JawabanGulihRepository) RecalculateGulih(ikasID string) error {
	// Query rata-rata jawaban per kategori_id untuk assessment tertentu
	query := `
		SELECT k.id AS kategori_id, ROUND(AVG(jg.jawaban_gulih), 2) AS avg_nilai
		FROM jawaban_gulih jg
		JOIN pertanyaan_gulih pg ON jg.pertanyaan_gulih_id = pg.id
		JOIN sub_kategori sk ON pg.sub_kategori_id = sk.id
		JOIN kategori k ON sk.kategori_id = k.id
		WHERE jg.ikas_id = ? AND jg.jawaban_gulih IS NOT NULL
		GROUP BY k.id
		ORDER BY k.id`

	rows, err := r.db.Query(query, ikasID)
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
			(ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			nilai_gulih = VALUES(nilai_gulih),
			nilai_subdomain1 = VALUES(nilai_subdomain1),
			nilai_subdomain2 = VALUES(nilai_subdomain2),
			nilai_subdomain3 = VALUES(nilai_subdomain3),
			nilai_subdomain4 = VALUES(nilai_subdomain4)`

	_, err = r.db.Exec(upsertQuery,
		ikasID,
		nilaiGulih,
		subdomain[15], subdomain[16], subdomain[17], subdomain[18],
	)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	var gulihID int
	err = r.db.QueryRow(`SELECT id FROM gulih WHERE ikas_id = ?`, ikasID).Scan(&gulihID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	_, err = r.db.Exec(`UPDATE ikas SET id_gulih = ? WHERE id = ?`, gulihID, ikasID)
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
		SET i.nilai_kematangan = ROUND(
			COALESCE(iden.nilai_identifikasi, 0) * 0.25 + 
			COALESCE(prot.nilai_proteksi, 0) * 0.30 + 
			COALESCE(det.nilai_deteksi, 0) * 0.25 + 
			COALESCE(g.nilai_gulih, 0) * 0.20
		, 2)
		WHERE i.id = ? AND (
			iden.id IS NOT NULL OR prot.id IS NOT NULL OR det.id IS NOT NULL OR g.id IS NOT NULL
		)`

	_, err = r.db.Exec(updateKematanganQuery, ikasID)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}

func (r *JawabanGulihRepository) UpsertToBuffer(req dto.CreateJawabanGulihRequest) error {
	// 1. Cek dulu apakah data ini SUDAH ADA di tabel utama
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM jawaban_gulih WHERE ikas_id = ? AND pertanyaan_gulih_id = ?)`
	err := r.db.QueryRow(checkQuery, req.IkasID, req.PertanyaanGulihID).Scan(&exists)
	if err != nil {
		return err
	}

	// 2. Jika sudah ada di utama, jangan masukkan ke buffer.
	// Sebaliknya, pastikan di buffer bersih (cleanup sampah)
	if exists {
		_, _ = r.db.Exec(`DELETE FROM jawaban_gulih_buffer WHERE ikas_id = ? AND pertanyaan_gulih_id = ?`, req.IkasID, req.PertanyaanGulihID)
		return nil
	}

	// 3. Jika belum ada di utama, baru masukkan ke buffer seperti biasa
	query := `INSERT INTO jawaban_gulih_buffer 
		(uuid, pertanyaan_gulih_id, ikas_id, jawaban_gulih, evidence, validasi, keterangan)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		jawaban_gulih = VALUES(jawaban_gulih),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	_, err = r.db.Exec(query,
		utils.GenerateUUID(),
		req.PertanyaanGulihID,
		req.IkasID,
		req.JawabanGulih,
		req.Evidence,
		req.Validasi,
		req.Keterangan,
	)
	return err
}

func (r *JawabanGulihRepository) GetBufferCount(ikasID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM jawaban_gulih_buffer WHERE ikas_id = ?`
	err := r.db.QueryRow(query, ikasID).Scan(&count)
	return count, err
}

func (r *JawabanGulihRepository) FlushBuffer(ikasID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	moveQuery := `INSERT INTO jawaban_gulih 
		(uuid, pertanyaan_gulih_id, ikas_id, jawaban_gulih, evidence, validasi, keterangan)
		SELECT uuid, pertanyaan_gulih_id, ikas_id, jawaban_gulih, evidence, validasi, keterangan
		FROM jawaban_gulih_buffer WHERE ikas_id = ?
		ON DUPLICATE KEY UPDATE 
		jawaban_gulih = VALUES(jawaban_gulih),
		evidence = VALUES(evidence),
		validasi = VALUES(validasi),
		keterangan = VALUES(keterangan)`

	if _, err := tx.Exec(moveQuery, ikasID); err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM jawaban_gulih_buffer WHERE ikas_id = ?`, ikasID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *JawabanGulihRepository) CloneByIkasID(sourceID, targetID string) error {
	query := `
		INSERT INTO jawaban_gulih 
			(uuid, pertanyaan_gulih_id, ikas_id, jawaban_gulih, evidence, validasi, keterangan)
		SELECT 
			uuid, pertanyaan_gulih_id, ?, jawaban_gulih, evidence, validasi, keterangan
		FROM jawaban_gulih 
		WHERE ikas_id = ?`

	_, err := r.db.Exec(query, targetID, sourceID)
	if err != nil {
		rollbar.Error(err)
		return err
	}
	return nil
}

func (r *JawabanGulihRepository) GetIDByIkasAndPertanyaan(ikasID string, pertanyaanID int) (int, error) {
	var id int
	err := r.db.QueryRow(`SELECT id FROM jawaban_gulih WHERE ikas_id = ? AND pertanyaan_gulih_id = ?`, ikasID, pertanyaanID).Scan(&id)
	return id, err
}
