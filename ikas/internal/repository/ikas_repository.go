package repository

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"ikas/internal/dto"
	"ikas/internal/utils"
	"strconv"
	"strings"
	"time"

	"sort"

	"github.com/xuri/excelize/v2"
)

type IkasRepository struct {
	db *sql.DB
}

// ImportFromExcel implements IkasRepositoryInterface.
func (r *IkasRepository) ImportFromExcel(raw []byte) (*dto.IkasResponse, error) {
	panic("unimplemented")
}

func NewIkasRepository(db *sql.DB) *IkasRepository {
	return &IkasRepository{db: db}
}

// Update method Create untuk menerima nilai_kematangan (inisiasi IKAS)
func (r *IkasRepository) Create(req dto.CreateIkasRequest, id string, nilaiKematangan float64) error {
	query := `INSERT INTO ikas
		(id, id_perusahaan, id_identifikasi, id_proteksi, id_deteksi, id_gulih, tanggal, responden, telepon, jabatan,
		nilai_kematangan, target_nilai)
		VALUES (?, ?, NULL, NULL, NULL, NULL, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id,
		req.IDPerusahaan,
		req.Tanggal,
		req.Responden,
		req.Telepon,
		req.Jabatan,
		nilaiKematangan,
		req.TargetNilai,
	)

	return err
}

func (r *IkasRepository) GetAll() ([]dto.IkasResponse, error) {
	query := `
		SELECT 
			i.id,
			i.tanggal,
			i.responden,
			i.telepon,
			i.jabatan,
			i.nilai_kematangan,
			i.target_nilai,
			p.id,
			p.nama_perusahaan,
			iden.id,
			iden.nilai_identifikasi,
			iden.nilai_subdomain1,
			iden.nilai_subdomain2,
			iden.nilai_subdomain3,
			iden.nilai_subdomain4,
			iden.nilai_subdomain5,
			prot.id,
			prot.nilai_proteksi,
			prot.nilai_subdomain1,
			prot.nilai_subdomain2,
			prot.nilai_subdomain3,
			prot.nilai_subdomain4,
			prot.nilai_subdomain5,
			prot.nilai_subdomain6,
			det.id,
			det.nilai_deteksi,
			det.nilai_subdomain1,
			det.nilai_subdomain2,
			det.nilai_subdomain3,
			g.id,
			g.nilai_gulih,
			g.nilai_subdomain1,
			g.nilai_subdomain2,
			g.nilai_subdomain3,
			g.nilai_subdomain4,
			i.created_at,
			i.updated_at
		FROM ikas i
		LEFT JOIN perusahaan p ON i.id_perusahaan = p.id
		LEFT JOIN identifikasi iden ON i.id_identifikasi = iden.id
		LEFT JOIN proteksi prot ON i.id_proteksi = prot.id
		LEFT JOIN deteksi det ON i.id_deteksi = det.id
		LEFT JOIN gulih g ON i.id_gulih = g.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err) // pastikan error ini keluar
	}
	defer rows.Close()

	var result []dto.IkasResponse

	for rows.Next() {
		var i dto.IkasResponse
		var tanggal sql.NullString
		var nilaiKematangan, targetNilai sql.NullFloat64
		var perusahaanID, perusahaanNama sql.NullString
		var idenID sql.NullInt64
		var idenNilai, idenSub1, idenSub2, idenSub3, idenSub4, idenSub5 sql.NullFloat64
		var protID sql.NullInt64
		var protNilai, protSub1, protSub2, protSub3, protSub4, protSub5, protSub6 sql.NullFloat64
		var detID sql.NullInt64
		var detNilai, detSub1, detSub2, detSub3 sql.NullFloat64
		var gulihID sql.NullInt64
		var gulihNilai, gulihSub1, gulihSub2, gulihSub3, gulihSub4 sql.NullFloat64
		var createdAt, updatedAt sql.NullString

		err := rows.Scan(
			&i.ID,
			&tanggal,
			&i.Responden,
			&i.Telepon,
			&i.Jabatan,
			&nilaiKematangan,
			&targetNilai,
			&perusahaanID,
			&perusahaanNama,
			&idenID,
			&idenNilai,
			&idenSub1,
			&idenSub2,
			&idenSub3,
			&idenSub4,
			&idenSub5,
			&protID,
			&protNilai,
			&protSub1,
			&protSub2,
			&protSub3,
			&protSub4,
			&protSub5,
			&protSub6,
			&detID,
			&detNilai,
			&detSub1,
			&detSub2,
			&detSub3,
			&gulihID,
			&gulihNilai,
			&gulihSub1,
			&gulihSub2,
			&gulihSub3,
			&gulihSub4,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			continue
		}

		if tanggal.Valid {
			i.Tanggal = tanggal.String
		}

		if createdAt.Valid {
			i.CreatedAt = createdAt.String
		}

		if updatedAt.Valid {
			i.UpdatedAt = updatedAt.String
		}

		if targetNilai.Valid {
			i.TargetNilai = targetNilai.Float64
		}

		// Selalu hitung nilai_kematangan secara dinamis dari domain yang sudah terisi dengan bobot:
		// Identifikasi (25%), Proteksi (30%), Deteksi (25%), Gulih (20%)
		{
			i.NilaiKematangan = utils.RoundToTwo(
				idenNilai.Float64*0.25 +
					protNilai.Float64*0.30 +
					detNilai.Float64*0.25 +
					gulihNilai.Float64*0.20,
			)
		}

		// Set kategori kematangan keamanan siber
		i.KategoriKematanganKeamananSiber = utils.GetKategoriTingkatKematangan(i.NilaiKematangan)

		// Map perusahaan
		if perusahaanID.Valid && perusahaanNama.Valid {
			i.Perusahaan = &dto.PerusahaanInIkas{
				ID:             perusahaanID.String,
				NamaPerusahaan: perusahaanNama.String,
			}
		}

		// Map identifikasi
		if idenID.Valid {
			i.Identifikasi = &dto.IdentifikasiInIkas{
				ID:                              int(idenID.Int64),
				NilaiIdentifikasi:               idenNilai.Float64,
				KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(idenNilai.Float64),
				NilaiSubdomain1:                 idenSub1.Float64,
				NilaiSubdomain2:                 idenSub2.Float64,
				NilaiSubdomain3:                 idenSub3.Float64,
				NilaiSubdomain4:                 idenSub4.Float64,
				NilaiSubdomain5:                 idenSub5.Float64,
			}
		}

		// Map proteksi
		if protID.Valid {
			i.Proteksi = &dto.ProteksiInIkas{
				ID:                              int(protID.Int64),
				NilaiProteksi:                   protNilai.Float64,
				KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(protNilai.Float64),
				NilaiSubdomain1:                 protSub1.Float64,
				NilaiSubdomain2:                 protSub2.Float64,
				NilaiSubdomain3:                 protSub3.Float64,
				NilaiSubdomain4:                 protSub4.Float64,
				NilaiSubdomain5:                 protSub5.Float64,
				NilaiSubdomain6:                 protSub6.Float64,
			}
		}

		// Map deteksi
		if detID.Valid {
			i.Deteksi = &dto.DeteksiInIkas{
				ID:                              int(detID.Int64),
				NilaiDeteksi:                    detNilai.Float64,
				KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(detNilai.Float64),
				NilaiSubdomain1:                 detSub1.Float64,
				NilaiSubdomain2:                 detSub2.Float64,
				NilaiSubdomain3:                 detSub3.Float64,
			}
		}

		// Map gulih
		if gulihID.Valid {
			i.Gulih = &dto.GulihInIkas{
				ID:                              int(gulihID.Int64),
				NilaiGulih:                      gulihNilai.Float64,
				KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(gulihNilai.Float64),
				NilaiSubdomain1:                 gulihSub1.Float64,
				NilaiSubdomain2:                 gulihSub2.Float64,
				NilaiSubdomain3:                 gulihSub3.Float64,
				NilaiSubdomain4:                 gulihSub4.Float64,
			}
		}

		result = append(result, i)
	}

	return result, nil
}

func (r *IkasRepository) GetByPerusahaan(perusahaanID string) ([]dto.IkasResponse, error) {
	query := `
		SELECT 
			i.id,
			i.tanggal,
			i.responden,
			i.telepon,
			i.jabatan,
			i.nilai_kematangan,
			i.target_nilai,
			p.id,
			p.nama_perusahaan,
			iden.id,
			iden.nilai_identifikasi,
			iden.nilai_subdomain1,
			iden.nilai_subdomain2,
			iden.nilai_subdomain3,
			iden.nilai_subdomain4,
			iden.nilai_subdomain5,
			prot.id,
			prot.nilai_proteksi,
			prot.nilai_subdomain1,
			prot.nilai_subdomain2,
			prot.nilai_subdomain3,
			prot.nilai_subdomain4,
			prot.nilai_subdomain5,
			prot.nilai_subdomain6,
			det.id,
			det.nilai_deteksi,
			det.nilai_subdomain1,
			det.nilai_subdomain2,
			det.nilai_subdomain3,
			g.id,
			g.nilai_gulih,
			g.nilai_subdomain1,
			g.nilai_subdomain2,
			g.nilai_subdomain3,
			g.nilai_subdomain4,
			i.created_at,
			i.updated_at
		FROM ikas i
		LEFT JOIN perusahaan p ON i.id_perusahaan = p.id
		LEFT JOIN identifikasi iden ON i.id_identifikasi = iden.id
		LEFT JOIN proteksi prot ON i.id_proteksi = prot.id
		LEFT JOIN deteksi det ON i.id_deteksi = det.id
		LEFT JOIN gulih g ON i.id_gulih = g.id
		WHERE i.id_perusahaan = ?
	`

	rows, err := r.db.Query(query, perusahaanID)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var result []dto.IkasResponse

	for rows.Next() {
		var i dto.IkasResponse
		var tanggal sql.NullString
		var nilaiKematangan, targetNilai sql.NullFloat64
		var perusahaanID, perusahaanNama sql.NullString
		var idenID sql.NullInt64
		var idenNilai, idenSub1, idenSub2, idenSub3, idenSub4, idenSub5 sql.NullFloat64
		var protID sql.NullInt64
		var protNilai, protSub1, protSub2, protSub3, protSub4, protSub5, protSub6 sql.NullFloat64
		var detID sql.NullInt64
		var detNilai, detSub1, detSub2, detSub3 sql.NullFloat64
		var gulihID sql.NullInt64
		var gulihNilai, gulihSub1, gulihSub2, gulihSub3, gulihSub4 sql.NullFloat64
		var createdAt, updatedAt sql.NullString

		err := rows.Scan(
			&i.ID,
			&tanggal,
			&i.Responden,
			&i.Telepon,
			&i.Jabatan,
			&nilaiKematangan,
			&targetNilai,
			&perusahaanID,
			&perusahaanNama,
			&idenID,
			&idenNilai,
			&idenSub1,
			&idenSub2,
			&idenSub3,
			&idenSub4,
			&idenSub5,
			&protID,
			&protNilai,
			&protSub1,
			&protSub2,
			&protSub3,
			&protSub4,
			&protSub5,
			&protSub6,
			&detID,
			&detNilai,
			&detSub1,
			&detSub2,
			&detSub3,
			&gulihID,
			&gulihNilai,
			&gulihSub1,
			&gulihSub2,
			&gulihSub3,
			&gulihSub4,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			continue
		}

		if tanggal.Valid {
			i.Tanggal = tanggal.String
		}

		if createdAt.Valid {
			i.CreatedAt = createdAt.String
		}

		if updatedAt.Valid {
			i.UpdatedAt = updatedAt.String
		}

		if targetNilai.Valid {
			i.TargetNilai = targetNilai.Float64
		}

		{
			i.NilaiKematangan = utils.RoundToTwo(
				idenNilai.Float64*0.25 +
					protNilai.Float64*0.30 +
					detNilai.Float64*0.25 +
					gulihNilai.Float64*0.20,
			)
		}

		i.KategoriKematanganKeamananSiber = utils.GetKategoriTingkatKematangan(i.NilaiKematangan)

		if perusahaanID.Valid && perusahaanNama.Valid {
			i.Perusahaan = &dto.PerusahaanInIkas{
				ID:             perusahaanID.String,
				NamaPerusahaan: perusahaanNama.String,
			}
		}

		if idenID.Valid {
			i.Identifikasi = &dto.IdentifikasiInIkas{
				ID:                              int(idenID.Int64),
				NilaiIdentifikasi:               idenNilai.Float64,
				KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(idenNilai.Float64),
				NilaiSubdomain1:                 idenSub1.Float64,
				NilaiSubdomain2:                 idenSub2.Float64,
				NilaiSubdomain3:                 idenSub3.Float64,
				NilaiSubdomain4:                 idenSub4.Float64,
				NilaiSubdomain5:                 idenSub5.Float64,
			}
		}

		if protID.Valid {
			i.Proteksi = &dto.ProteksiInIkas{
				ID:                              int(protID.Int64),
				NilaiProteksi:                   protNilai.Float64,
				KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(protNilai.Float64),
				NilaiSubdomain1:                 protSub1.Float64,
				NilaiSubdomain2:                 protSub2.Float64,
				NilaiSubdomain3:                 protSub3.Float64,
				NilaiSubdomain4:                 protSub4.Float64,
				NilaiSubdomain5:                 protSub5.Float64,
				NilaiSubdomain6:                 protSub6.Float64,
			}
		}

		if detID.Valid {
			i.Deteksi = &dto.DeteksiInIkas{
				ID:                              int(detID.Int64),
				NilaiDeteksi:                    detNilai.Float64,
				KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(detNilai.Float64),
				NilaiSubdomain1:                 detSub1.Float64,
				NilaiSubdomain2:                 detSub2.Float64,
				NilaiSubdomain3:                 detSub3.Float64,
			}
		}

		if gulihID.Valid {
			i.Gulih = &dto.GulihInIkas{
				ID:                              int(gulihID.Int64),
				NilaiGulih:                      gulihNilai.Float64,
				KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(gulihNilai.Float64),
				NilaiSubdomain1:                 gulihSub1.Float64,
				NilaiSubdomain2:                 gulihSub2.Float64,
				NilaiSubdomain3:                 gulihSub3.Float64,
				NilaiSubdomain4:                 gulihSub4.Float64,
			}
		}

		result = append(result, i)
	}

	return result, nil
}

func (r *IkasRepository) GetByID(id string) (*dto.IkasResponse, error) {
	query := `
		SELECT 
			i.id,
			i.tanggal,
			i.responden,
			i.telepon,
			i.jabatan,
			i.nilai_kematangan,
			i.target_nilai,
			p.id,
			p.nama_perusahaan,
			iden.id,
			iden.nilai_identifikasi,
			iden.nilai_subdomain1,
			iden.nilai_subdomain2,
			iden.nilai_subdomain3,
			iden.nilai_subdomain4,
			iden.nilai_subdomain5,
			prot.id,
			prot.nilai_proteksi,
			prot.nilai_subdomain1,
			prot.nilai_subdomain2,
			prot.nilai_subdomain3,
			prot.nilai_subdomain4,
			prot.nilai_subdomain5,
			prot.nilai_subdomain6,
			det.id,
			det.nilai_deteksi,
			det.nilai_subdomain1,
			det.nilai_subdomain2,
			det.nilai_subdomain3,
			g.id,
			g.nilai_gulih,
			g.nilai_subdomain1,
			g.nilai_subdomain2,
			g.nilai_subdomain3,
			g.nilai_subdomain4,
			i.created_at,
			i.updated_at
		FROM ikas i
		LEFT JOIN perusahaan p ON i.id_perusahaan = p.id
		LEFT JOIN identifikasi iden ON i.id_identifikasi = iden.id
		LEFT JOIN proteksi prot ON i.id_proteksi = prot.id
		LEFT JOIN deteksi det ON i.id_deteksi = det.id
		LEFT JOIN gulih g ON i.id_gulih = g.id
		WHERE i.id = ?
	`

	row := r.db.QueryRow(query, id)

	var i dto.IkasResponse
	var tanggal sql.NullString
	var nilaiKematangan, targetNilai sql.NullFloat64
	var perusahaanID, perusahaanNama sql.NullString
	var idenID sql.NullInt64
	var idenNilai, idenSub1, idenSub2, idenSub3, idenSub4, idenSub5 sql.NullFloat64
	var protID sql.NullInt64
	var protNilai, protSub1, protSub2, protSub3, protSub4, protSub5, protSub6 sql.NullFloat64
	var detID sql.NullInt64
	var detNilai, detSub1, detSub2, detSub3 sql.NullFloat64
	var gulihID sql.NullInt64
	var gulihNilai, gulihSub1, gulihSub2, gulihSub3, gulihSub4 sql.NullFloat64
	var createdAt, updatedAt sql.NullString

	err := row.Scan(
		&i.ID,
		&tanggal,
		&i.Responden,
		&i.Telepon,
		&i.Jabatan,
		&nilaiKematangan,
		&targetNilai,
		&perusahaanID,
		&perusahaanNama,
		&idenID,
		&idenNilai,
		&idenSub1,
		&idenSub2,
		&idenSub3,
		&idenSub4,
		&idenSub5,
		&protID,
		&protNilai,
		&protSub1,
		&protSub2,
		&protSub3,
		&protSub4,
		&protSub5,
		&protSub6,
		&detID,
		&detNilai,
		&detSub1,
		&detSub2,
		&detSub3,
		&gulihID,
		&gulihNilai,
		&gulihSub1,
		&gulihSub2,
		&gulihSub3,
		&gulihSub4,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	if tanggal.Valid {
		i.Tanggal = tanggal.String
	}

	if createdAt.Valid {
		i.CreatedAt = createdAt.String
	}

	if updatedAt.Valid {
		i.UpdatedAt = updatedAt.String
	}

	if targetNilai.Valid {
		i.TargetNilai = targetNilai.Float64
	}

	// Selalu hitung nilai_kematangan secara dinamis dari domain yang sudah terisi dengan bobot:
	// Identifikasi (25%), Proteksi (30%), Deteksi (25%), Gulih (20%)
	{
		i.NilaiKematangan = utils.RoundToTwo(
			idenNilai.Float64*0.25 +
				protNilai.Float64*0.30 +
				detNilai.Float64*0.25 +
				gulihNilai.Float64*0.20,
		)
	}

	// Set kategori kematangan keamanan siber
	i.KategoriKematanganKeamananSiber = utils.GetKategoriTingkatKematangan(i.NilaiKematangan)

	// Map perusahaan
	if perusahaanID.Valid && perusahaanNama.Valid {
		i.Perusahaan = &dto.PerusahaanInIkas{
			ID:             perusahaanID.String,
			NamaPerusahaan: perusahaanNama.String,
		}
	}

	// Map identifikasi
	if idenID.Valid {
		i.Identifikasi = &dto.IdentifikasiInIkas{
			ID:                              int(idenID.Int64),
			NilaiIdentifikasi:               idenNilai.Float64,
			KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(idenNilai.Float64),
			NilaiSubdomain1:                 idenSub1.Float64,
			NilaiSubdomain2:                 idenSub2.Float64,
			NilaiSubdomain3:                 idenSub3.Float64,
			NilaiSubdomain4:                 idenSub4.Float64,
			NilaiSubdomain5:                 idenSub5.Float64,
		}
	}

	// Map proteksi
	if protID.Valid {
		i.Proteksi = &dto.ProteksiInIkas{
			ID:                              int(protID.Int64),
			NilaiProteksi:                   protNilai.Float64,
			KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(protNilai.Float64),
			NilaiSubdomain1:                 protSub1.Float64,
			NilaiSubdomain2:                 protSub2.Float64,
			NilaiSubdomain3:                 protSub3.Float64,
			NilaiSubdomain4:                 protSub4.Float64,
			NilaiSubdomain5:                 protSub5.Float64,
			NilaiSubdomain6:                 protSub6.Float64,
		}
	}

	// Map deteksi
	if detID.Valid {
		i.Deteksi = &dto.DeteksiInIkas{
			ID:                              int(detID.Int64),
			NilaiDeteksi:                    detNilai.Float64,
			KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(detNilai.Float64),
			NilaiSubdomain1:                 detSub1.Float64,
			NilaiSubdomain2:                 detSub2.Float64,
			NilaiSubdomain3:                 detSub3.Float64,
		}
	}

	// Map gulih
	if gulihID.Valid {
		i.Gulih = &dto.GulihInIkas{
			ID:                              int(gulihID.Int64),
			NilaiGulih:                      gulihNilai.Float64,
			KategoriTingkatKematanganDomain: utils.GetKategoriTingkatKematangan(gulihNilai.Float64),
			NilaiSubdomain1:                 gulihSub1.Float64,
			NilaiSubdomain2:                 gulihSub2.Float64,
			NilaiSubdomain3:                 gulihSub3.Float64,
			NilaiSubdomain4:                 gulihSub4.Float64,
		}
	}

	return &i, nil
}

func (r *IkasRepository) Update(id string, req dto.UpdateIkasRequest) error {
	query := "UPDATE ikas SET "
	args := []interface{}{}
	updates := []string{}

	if req.IDPerusahaan != nil {
		updates = append(updates, "id_perusahaan=?")
		args = append(args, *req.IDPerusahaan)
	}
	if req.Tanggal != nil {
		updates = append(updates, "tanggal=?")
		args = append(args, *req.Tanggal)
	}
	if req.Responden != nil {
		updates = append(updates, "responden=?")
		args = append(args, *req.Responden)
	}
	if req.Telepon != nil {
		updates = append(updates, "telepon=?")
		args = append(args, *req.Telepon)
	}
	if req.Jabatan != nil {
		updates = append(updates, "jabatan=?")
		args = append(args, *req.Jabatan)
	}
	if req.TargetNilai != nil {
		updates = append(updates, "target_nilai=?")
		args = append(args, *req.TargetNilai)
	}

	if len(updates) == 0 {
		return nil
	}

	query += strings.Join(updates, ", ")
	query += " WHERE id=?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *IkasRepository) Delete(id string) error {
	// 1. Ambil data ikas untuk mendapatkan ID domain yang terkait
	var idenID, protID, detID, gulIHID sql.NullInt64
	queryGet := `SELECT id_identifikasi, id_proteksi, id_deteksi, id_gulih FROM ikas WHERE id = ?`
	err := r.db.QueryRow(queryGet, id).Scan(&idenID, &protID, &detID, &gulIHID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 2. Delete buffer data by ikas_id
	tablesBuffer := []string{
		"jawaban_identifikasi_buffer",
		"jawaban_proteksi_buffer",
		"jawaban_deteksi_buffer",
		"jawaban_gulih_buffer",
	}
	for _, table := range tablesBuffer {
		if _, err := tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE ikas_id = ?`, table), id); err != nil {
			return fmt.Errorf("error deleting from %s: %v", table, err)
		}
	}

	// 3. Delete main answers by ikas_id
	tablesJawaban := []string{
		"jawaban_identifikasi",
		"jawaban_proteksi",
		"jawaban_deteksi",
		"jawaban_gulih",
	}
	for _, table := range tablesJawaban {
		if _, err := tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE ikas_id = ?`, table), id); err != nil {
			return fmt.Errorf("error deleting from %s: %v", table, err)
		}
	}

	// 4. Delete ikas (lebih dulu agar FK tidak bermasalah jika ada restriksi)
	if _, err := tx.Exec(`DELETE FROM ikas WHERE id = ?`, id); err != nil {
		return fmt.Errorf("error deleting from ikas: %v", err)
	}

	// 5. Delete domain data by specific ID (bukan perusahaan_id lagi)
	if idenID.Valid {
		if _, err := tx.Exec(`DELETE FROM identifikasi WHERE id = ?`, idenID.Int64); err != nil {
			return fmt.Errorf("error deleting from identifikasi: %v", err)
		}
	}
	if protID.Valid {
		if _, err := tx.Exec(`DELETE FROM proteksi WHERE id = ?`, protID.Int64); err != nil {
			return fmt.Errorf("error deleting from proteksi: %v", err)
		}
	}
	if detID.Valid {
		if _, err := tx.Exec(`DELETE FROM deteksi WHERE id = ?`, detID.Int64); err != nil {
			return fmt.Errorf("error deleting from deteksi: %v", err)
		}
	}
	if gulIHID.Valid {
		if _, err := tx.Exec(`DELETE FROM gulih WHERE id = ?`, gulIHID.Int64); err != nil {
			return fmt.Errorf("error deleting from gulih: %v", err)
		}
	}

	return tx.Commit()
}

func parseMultipleDateFormats(dateStr string) (time.Time, error) {
	// Daftar format yang didukung
	formats := []string{
		"02-01-2006", // DD-MM-YYYY
		"02/01/2006", // DD/MM/YYYY
		"02-01-06",   // DD-MM-YY
		"02/01/06",   // DD/MM/YY
		"2006-01-02", // YYYY-MM-DD (ISO)
		"01-02-2006", // MM-DD-YYYY
		"01/02/2006", // MM/DD/YYYY
		"01-02-06",   // MM-DD-YY
		"01/02/06",   // MM/DD/YY
		"2006/01/02", // YYYY/MM/DD
	}

	var lastErr error
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}

	return time.Time{}, lastErr
}

// ParseExcelForImport membaca file Excel dari sheet 2 dan sheet 7
// ParseExcelForImport membaca file Excel dan memetakan datanya ke dto.ParsedExcelData
func (r *IkasRepository) ParseExcelForImport(fileData []byte) (*dto.ParsedExcelData, error) {
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()

	// Validasi jumlah sheet (minimal 6 sheet: index 0-5)
	if len(sheets) < 6 {
		return nil, errors.New("file Excel tidak lengkap, minimal harus memiliki 6 sheet")
	}

	sheet2 := sheets[1] // Sheet ke-2 (index 1) - Info Dasar
	sheet3 := sheets[2] // Sheet ke-3 (index 2) - Identifikasi
	sheet4 := sheets[3] // Sheet ke-4 (index 3) - Proteksi
	sheet5 := sheets[4] // Sheet ke-5 (index 4) - Deteksi
	sheet6 := sheets[5] // Sheet ke-6 (index 5) - Gulih

	// Helper function untuk ambil nilai cell sebagai string
	getCellString := func(sheetName, cell string) (string, error) {
		val, err := f.GetCellValue(sheetName, cell)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(val), nil
	}

	// Helper function untuk ambil nilai cell sebagai float
	getCellFloat := func(sheetName, cell string) (float64, error) {
		val, err := f.GetCellValue(sheetName, cell)
		if err != nil {
			return 0, err
		}
		val = strings.TrimSpace(val)
		if val == "" || val == "N/A" {
			return 0, nil // Return 0 if empty or N/A, handled by business logic later
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			// Try parsing as percentage if applicable or handle other formats if needed
			return 0, fmt.Errorf("gagal parse cell %s di %s: %v", cell, sheetName, err)
		}
		return floatVal, nil
	}

	// ===== AMBIL DATA DASAR (SHEET 2) =====

	// Nama Perusahaan dari D4
	namaPerusahaan, err := getCellString(sheet2, "D4")
	if err != nil {
		return nil, fmt.Errorf("error membaca nama perusahaan (Sheet 2, D4): %v", err)
	}
	if namaPerusahaan == "" {
		return nil, errors.New("nama perusahaan (Sheet 2, D4) tidak boleh kosong")
	}

	// Cari ID Perusahaan berdasarkan nama
	idPerusahaan, err := r.FindPerusahaanByName(namaPerusahaan)
	if err != nil {
		return nil, fmt.Errorf("error mencari perusahaan: %v", err)
	}
	if idPerusahaan == "" {
		return nil, fmt.Errorf("perusahaan anda belum terdaftar: '%s'", namaPerusahaan)
	}

	// Tanggal dari D18 (Perlu diekstrak lebih awal untuk validasi tahun)
	tanggalStr, err := getCellString(sheet2, "D18")
	if err != nil {
		return nil, fmt.Errorf("error membaca tanggal (Sheet 2, D18): %v", err)
	}

	// Extract year from date for duplicate check
	var tahun int
	if tanggalStr != "" {
		if excelDate, err := strconv.ParseFloat(tanggalStr, 64); err == nil {
			parsedTime, err := excelize.ExcelDateToTime(excelDate, false)
			if err == nil {
				tahun = parsedTime.Year()
			}
		} else {
			parsedTime, err := parseMultipleDateFormats(tanggalStr)
			if err == nil {
				tahun = parsedTime.Year()
			}
		}
	}

	if tahun == 0 {
		tahun = time.Now().Year()
	}

	// VALIDASI: Cek apakah data IKAS untuk perusahaan ini sudah ada di tahun tersebut
	exists, err := r.CheckExistsByPerusahaanIDAndYear(idPerusahaan, tahun)
	if err != nil {
		return nil, fmt.Errorf("error validasi duplikasi data: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("Data IKAS untuk perusahaan '%s' pada tahun %d sudah ada. Anda tidak dapat melakukan import lagi.", namaPerusahaan, tahun)
	}

	// Telepon dari D10
	telepon, err := getCellString(sheet2, "D10")
	if err != nil {
		return nil, fmt.Errorf("error membaca telepon (Sheet 2, D10): %v", err)
	}

	// Responden dari D11
	responden, err := getCellString(sheet2, "D11")
	if err != nil {
		return nil, fmt.Errorf("error membaca responden (Sheet 2, D11): %v", err)
	}

	// Jabatan dari D12
	jabatan, err := getCellString(sheet2, "D12")
	if err != nil {
		return nil, fmt.Errorf("error membaca jabatan (Sheet 2, D12): %v", err)
	}

	// Target Nilai dari D15
	targetNilai, err := getCellFloat(sheet2, "D15")
	if err != nil {
		return nil, fmt.Errorf("error membaca target_nilai (Sheet 2, D15): %v", err)
	}

	var tanggal string
	if tanggalStr != "" {
		if excelDate, err := strconv.ParseFloat(tanggalStr, 64); err == nil {
			parsedTime, err := excelize.ExcelDateToTime(excelDate, false)
			if err == nil {
				tanggal = parsedTime.Format("2006-01-02")
			}
		} else {
			parsedTime, err := parseMultipleDateFormats(tanggalStr)
			if err == nil {
				tanggal = parsedTime.Format("2006-01-02")
			} else {
				return nil, fmt.Errorf("format tanggal tidak valid (Sheet 2, D18): %s", tanggalStr)
			}
		}
	} else {
		return nil, errors.New("tanggal (Sheet 2, D18) tidak boleh kosong")
	}

	result := &dto.ParsedExcelData{
		IkasRequest: dto.CreateIkasRequest{
			IDPerusahaan: idPerusahaan,
			Tanggal:      tanggal,
			Responden:    responden,
			Telepon:      telepon,
			Jabatan:      jabatan,
			TargetNilai:  targetNilai,
		},
	}

	// Helper to collect answers based on mappings
	collectAnswers := func(sheetName string, mappings map[string]int) ([]dto.ExcelSubdomainAnswer, error) {
		var answers []dto.ExcelSubdomainAnswer
		for cell, qID := range mappings {
			val, err := getCellFloat(sheetName, cell)
			if err != nil {
				// We might want to skip or handle errors differently, but for now let's be strict
				return nil, err
			}
			answers = append(answers, dto.ExcelSubdomainAnswer{
				PertanyaanID: qID,
				Jawaban:      val,
			})
		}

		// Sort manual berdasarkan PertanyaanID agar urut saat masuk ke DB
		sort.Slice(answers, func(i, j int) bool {
			return answers[i].PertanyaanID < answers[j].PertanyaanID
		})

		return answers, nil
	}

	// ===== SHEET 3: IDENTIFIKASI =====
	identifikasiMap := map[string]int{
		"D5": 1, "D6": 2, "D8": 3, "D10": 4, "D11": 5, "D14": 6, "D15": 7, "D16": 8, "D18": 9, "D19": 10,
		"D21": 11, "D23": 12, "D26": 13, "D27": 14, "D29": 15, "D31": 16, "D32": 17, "D34": 18, "D36": 19, "D37": 20,
		"D40": 21, "D42": 22, "D43": 23, "D45": 24, "D47": 25, "D48": 26, "D49": 27, "D50": 28, "D51": 29, "D52": 30,
		"D53": 31, "D55": 32, "D56": 33, "D58": 34, "D59": 35, "D61": 36, "D62": 37, "D63": 38, "D65": 39, "D68": 40,
		"D69": 41, "D71": 42, "D72": 43, "D74": 44, "D75": 45, "D77": 46, "D78": 47, "D79": 48, "D81": 49, "D82": 50,
		"D83": 51,
	}
	result.JawabanIdentifikasi, err = collectAnswers(sheet3, identifikasiMap)
	if err != nil {
		return nil, err
	}

	// ===== SHEET 4: PROTEKSI =====
	proteksiMap := map[string]int{
		"D5": 1, "D7": 2, "D8": 3, "D10": 4, "D11": 5, "D13": 6, "D16": 7, "D17": 8, "D18": 9, "D20": 10,
		"D21": 11, "D23": 12, "D24": 13, "D26": 14, "D28": 15, "D31": 16, "D32": 17, "D34": 18, "D35": 19, "D37": 20,
		"D38": 21, "D40": 22, "D41": 23, "D43": 24, "D44": 25, "D45": 26, "D46": 27, "D47": 28, "D49": 29, "D51": 30,
		"D54": 31, "D55": 32, "D56": 33, "D58": 34, "D60": 35, "D62": 36, "D65": 37, "D66": 38, "D68": 39, "D69": 40,
		"D71": 41, "D72": 42, "D74": 43, "D76": 44, "D77": 45, "D78": 46, "D79": 47, "D81": 48, "D83": 49, "D86": 50,
		"D87": 51, "D88": 52, "D89": 53, "D90": 54, "D91": 55, "D92": 56, "D94": 57, "D95": 58, "D96": 59, "D98": 60,
		"D99": 61, "D100": 62,
	}
	result.JawabanProteksi, err = collectAnswers(sheet4, proteksiMap)
	if err != nil {
		return nil, err
	}

	// ===== SHEET 5: DETEKSI =====
	deteksiMap := map[string]int{
		"D5": 1, "D6": 2, "D8": 3, "D9": 4, "D11": 5, "D13": 6, "D14": 7, "D17": 8, "D19": 9, "D20": 10,
		"D22": 11, "D24": 12, "D27": 13, "D28": 14, "D29": 15, "D31": 16, "D33": 17, "D35": 18, "D36": 19,
	}
	result.JawabanDeteksi, err = collectAnswers(sheet5, deteksiMap)
	if err != nil {
		return nil, err
	}

	// ===== SHEET 6: GULIH =====
	gulihMap := map[string]int{
		"D5": 1, "D6": 2, "D7": 3, "D8": 4, "D10": 5, "D11": 6, "D12": 7, "D13": 8, "D15": 9, "D17": 10,
		"D18": 11, "D19": 12, "D20": 13, "D22": 14, "D25": 15, "D26": 16, "D28": 17, "D30": 18, "D32": 19, "D33": 20,
		"D36": 21, "D37": 22, "D38": 23, "D39": 24, "D41": 25, "D43": 26, "D44": 27, "D46": 28, "D48": 29, "D50": 30,
		"D51": 31, "D52": 32, "D53": 33, "D54": 34, "D55": 35, "D57": 36, "D58": 37, "D59": 38, "D61": 39, "D63": 40,
		"D64": 41, "D66": 42, "D69": 43, "D71": 44, "D72": 45, "D74": 46, "D75": 47, "D77": 48, "D78": 49,
	}
	result.JawabanGulih, err = collectAnswers(sheet6, gulihMap)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindPerusahaanByName mencari ID perusahaan berdasarkan nama (case-insensitive, exact match)
func (r *IkasRepository) FindPerusahaanByName(namaPerusahaan string) (string, error) {
	var id string
	query := `SELECT id FROM perusahaan WHERE LOWER(TRIM(nama_perusahaan)) = LOWER(TRIM(?))`

	err := r.db.QueryRow(query, namaPerusahaan).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return id, nil
}

// CheckExistsByPerusahaanID mengecek apakah data IKAS untuk perusahaan tersebut sudah ada
func (r *IkasRepository) CheckExistsByPerusahaanID(id string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM ikas WHERE id_perusahaan = ?", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *IkasRepository) CheckExistsByPerusahaanIDAndYear(id string, year int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM ikas WHERE id_perusahaan = ? AND YEAR(tanggal) = ?", id, year).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *IkasRepository) GetIDByPerusahaanID(idPerusahaan string) (string, error) {
	var id string
	query := `SELECT id FROM ikas WHERE id_perusahaan = ?`
	err := r.db.QueryRow(query, idPerusahaan).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (r *IkasRepository) CheckOwnership(ikasID string, perusahaanID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM ikas WHERE id = ? AND id_perusahaan = ?`
	err := r.db.QueryRow(query, ikasID, perusahaanID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
