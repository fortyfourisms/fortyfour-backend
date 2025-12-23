package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"strings"
)

type IkasRepository struct {
	db *sql.DB
}

func NewIkasRepository(db *sql.DB) *IkasRepository {
	return &IkasRepository{db: db}
}

func (r *IkasRepository) Create(req dto.CreateIkasRequest, id string) error {
	query := `INSERT INTO ikas
				(id, id_perusahaan, tanggal, responden, telepon, jabatan,
				nilai_kematangan, target_nilai, id_identifikasi, id_proteksi,
				id_deteksi, id_gulih)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id,
		req.IDPerusahaan,
		req.Tanggal,
		req.Responden,
		req.Telepon,
		req.Jabatan,
		req.NilaiKematangan,
		req.TargetNilai,
		req.IDIdentifikasi,
		req.IDProteksi,
		req.IDDeteksi,
		req.IDGulih,
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
			g.nilai_subdomain4
		FROM ikas i
		LEFT JOIN perusahaan p ON i.id_perusahaan = p.id
		LEFT JOIN identifikasi iden ON i.id_identifikasi = iden.id
		LEFT JOIN proteksi prot ON i.id_proteksi = prot.id
		LEFT JOIN deteksi det ON i.id_deteksi = det.id
		LEFT JOIN gulih g ON i.id_gulih = g.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.IkasResponse

	for rows.Next() {
		var i dto.IkasResponse
		var perusahaanID, perusahaanNama sql.NullString
		var idenID sql.NullString
		var idenNilai, idenSub1, idenSub2, idenSub3, idenSub4, idenSub5 sql.NullFloat64
		var protID sql.NullString
		var protNilai, protSub1, protSub2, protSub3, protSub4, protSub5, protSub6 sql.NullFloat64
		var detID sql.NullString
		var detNilai, detSub1, detSub2, detSub3 sql.NullFloat64
		var gulihID sql.NullString
		var gulihNilai, gulihSub1, gulihSub2, gulihSub3, gulihSub4 sql.NullFloat64

		err := rows.Scan(
			&i.ID,
			&i.Tanggal,
			&i.Responden,
			&i.Telepon,
			&i.Jabatan,
			&i.NilaiKematangan,
			&i.TargetNilai,
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
		)
		if err != nil {
			continue
		}

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
				ID:                idenID.String,
				NilaiIdentifikasi: idenNilai.Float64,
				// NilaiSubdomain1:   idenSub1.Float64,
				// NilaiSubdomain2:   idenSub2.Float64,
				// NilaiSubdomain3:   idenSub3.Float64,
				// NilaiSubdomain4:   idenSub4.Float64,
				// NilaiSubdomain5:   idenSub5.Float64,
			}
		}

		// Map proteksi
		if protID.Valid {
			i.Proteksi = &dto.ProteksiInIkas{
				ID:            protID.String,
				NilaiProteksi: protNilai.Float64,
				// NilaiSubdomain1: protSub1.Float64,
				// NilaiSubdomain2: protSub2.Float64,
				// NilaiSubdomain3: protSub3.Float64,
				// NilaiSubdomain4: protSub4.Float64,
				// NilaiSubdomain5: protSub5.Float64,
				// NilaiSubdomain6: protSub6.Float64,
			}
		}

		// Map deteksi
		if detID.Valid {
			i.Deteksi = &dto.DeteksiInIkas{
				ID:           detID.String,
				NilaiDeteksi: detNilai.Float64,
				// NilaiSubdomain1: detSub1.Float64,
				// NilaiSubdomain2: detSub2.Float64,
				// NilaiSubdomain3: detSub3.Float64,
			}
		}

		// Map gulih
		if gulihID.Valid {
			i.Gulih = &dto.GulihInIkas{
				ID:         gulihID.String,
				NilaiGulih: gulihNilai.Float64,
				// NilaiSubdomain1: gulihSub1.Float64,
				// NilaiSubdomain2: gulihSub2.Float64,
				// NilaiSubdomain3: gulihSub3.Float64,
				// NilaiSubdomain4: gulihSub4.Float64,
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
			g.nilai_subdomain4
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
	var perusahaanID, perusahaanNama sql.NullString
	var idenID sql.NullString
	var idenNilai, idenSub1, idenSub2, idenSub3, idenSub4, idenSub5 sql.NullFloat64
	var protID sql.NullString
	var protNilai, protSub1, protSub2, protSub3, protSub4, protSub5, protSub6 sql.NullFloat64
	var detID sql.NullString
	var detNilai, detSub1, detSub2, detSub3 sql.NullFloat64
	var gulihID sql.NullString
	var gulihNilai, gulihSub1, gulihSub2, gulihSub3, gulihSub4 sql.NullFloat64

	err := row.Scan(
		&i.ID,
		&i.Tanggal,
		&i.Responden,
		&i.Telepon,
		&i.Jabatan,
		&i.NilaiKematangan,
		&i.TargetNilai,
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
	)
	if err != nil {
		return nil, err
	}

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
			ID:                idenID.String,
			NilaiIdentifikasi: idenNilai.Float64,
			// NilaiSubdomain1:   idenSub1.Float64,
			// NilaiSubdomain2:   idenSub2.Float64,
			// NilaiSubdomain3:   idenSub3.Float64,
			// NilaiSubdomain4:   idenSub4.Float64,
			// NilaiSubdomain5:   idenSub5.Float64,
		}
	}

	// Map proteksi
	if protID.Valid {
		i.Proteksi = &dto.ProteksiInIkas{
			ID:            protID.String,
			NilaiProteksi: protNilai.Float64,
			// NilaiSubdomain1: protSub1.Float64,
			// NilaiSubdomain2: protSub2.Float64,
			// NilaiSubdomain3: protSub3.Float64,
			// NilaiSubdomain4: protSub4.Float64,
			// NilaiSubdomain5: protSub5.Float64,
			// NilaiSubdomain6: protSub6.Float64,
		}
	}

	// Map deteksi
	if detID.Valid {
		i.Deteksi = &dto.DeteksiInIkas{
			ID:           detID.String,
			NilaiDeteksi: detNilai.Float64,
			// NilaiSubdomain1: detSub1.Float64,
			// NilaiSubdomain2: detSub2.Float64,
			// NilaiSubdomain3: detSub3.Float64,
		}
	}

	// Map gulih
	if gulihID.Valid {
		i.Gulih = &dto.GulihInIkas{
			ID:         gulihID.String,
			NilaiGulih: gulihNilai.Float64,
			// NilaiSubdomain1: gulihSub1.Float64,
			// NilaiSubdomain2: gulihSub2.Float64,
			// NilaiSubdomain3: gulihSub3.Float64,
			// NilaiSubdomain4: gulihSub4.Float64,
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
	if req.NilaiKematangan != nil {
		updates = append(updates, "nilai_kematangan=?")
		args = append(args, *req.NilaiKematangan)
	}
	if req.TargetNilai != nil {
		updates = append(updates, "target_nilai=?")
		args = append(args, *req.TargetNilai)
	}
	if req.IDIdentifikasi != nil {
		updates = append(updates, "id_identifikasi=?")
		args = append(args, *req.IDIdentifikasi)
	}
	if req.IDProteksi != nil {
		updates = append(updates, "id_proteksi=?")
		args = append(args, *req.IDProteksi)
	}
	if req.IDDeteksi != nil {
		updates = append(updates, "id_deteksi=?")
		args = append(args, *req.IDDeteksi)
	}
	if req.IDGulih != nil {
		updates = append(updates, "id_gulih=?")
		args = append(args, *req.IDGulih)
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
	_, err := r.db.Exec(`DELETE FROM ikas WHERE id=?`, id)
	return err
}
