package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
)

type seRepository struct {
	db *sql.DB
}

func NewSERepository(db *sql.DB) SERepositoryInterface {
	return &seRepository{db: db}
}

/* ================= CREATE ================= */

func (r *seRepository) Create(
	req dto.CreateSERequest,
	id string,
	totalBobot int,
	kategori string,
) error {
	_, err := r.db.Exec(`
		INSERT INTO se (
			id,
			id_perusahaan,
			id_sub_sektor,
			id_csirt,
			nilai_investasi,
			anggaran_operasional,
			kepatuhan_peraturan,
			teknik_kriptografi,
			jumlah_pengguna,
			data_pribadi,
			klasifikasi_data,
			kekritisan_proses,
			dampak_kegagalan,
			potensi_kerugian_dan_dampak_negatif,
			nama_se,
			ip_se,
			as_number_se,
			pengelola_se,
			fitur_se,
			total_bobot,
			kategori_se
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`,
		id,
		req.IDPerusahaan,
		req.IDSubSektor,
		req.IDCsirt,
		req.NilaiInvestasi,
		req.AnggaranOperasional,
		req.KepatuhanPeraturan,
		req.TeknikKriptografi,
		req.JumlahPengguna,
		req.DataPribadi,
		req.KlasifikasiData,
		req.KekritisanProses,
		req.DampakKegagalan,
		req.PotensiKerugiandanDampakNegatif,
		req.NamaSE,
		req.IpSE,
		req.AsNumberSE,
		req.PengelolaSE,
		req.FiturSE,
		totalBobot,
		kategori,
	)

	return err
}

/* ================= GET ALL ================= */

func (r *seRepository) GetAll() ([]dto.SEResponse, error) {
	rows, err := r.db.Query(`
		SELECT
			se.id,
			se.id_perusahaan,
			se.id_sub_sektor,
			se.id_csirt,
			se.nilai_investasi,
			se.anggaran_operasional,
			se.kepatuhan_peraturan,
			se.teknik_kriptografi,
			se.jumlah_pengguna,
			se.data_pribadi,
			se.klasifikasi_data,
			se.kekritisan_proses,
			se.dampak_kegagalan,
			se.potensi_kerugian_dan_dampak_negatif,
			se.nama_se,
			se.ip_se,
			se.as_number_se,
			se.pengelola_se,
			se.fitur_se,
			se.total_bobot,
			se.kategori_se,
			se.created_at,
			se.updated_at,

			p.id,
			p.nama_perusahaan,

			COALESCE(ss.id, ''),
			COALESCE(ss.nama_sub_sektor, ''),
			COALESCE(s.id, ''),
			COALESCE(s.nama_sektor, ''),

			COALESCE(c.id, ''),
			COALESCE(c.nama_csirt, '')
		FROM se
		JOIN perusahaan p ON se.id_perusahaan = p.id
		LEFT JOIN sub_sektor ss ON se.id_sub_sektor = ss.id
		LEFT JOIN sektor s ON ss.id_sektor = s.id
		LEFT JOIN csirt c ON se.id_csirt = c.id
		ORDER BY se.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.SEResponse

	for rows.Next() {
		var se dto.SEResponse
		se.Perusahaan = &dto.PerusahaanMiniResponse{}
		se.SubSektor = &dto.SubSektorMiniResponse{}
		se.Csirt = &dto.CsirtMiniResponse{}

		err := rows.Scan(
			&se.ID,
			&se.IDPerusahaan,
			&se.IDSubSektor,
			&se.IDCsirt,
			&se.NilaiInvestasi,
			&se.AnggaranOperasional,
			&se.KepatuhanPeraturan,
			&se.TeknikKriptografi,
			&se.JumlahPengguna,
			&se.DataPribadi,
			&se.KlasifikasiData,
			&se.KekritisanProses,
			&se.DampakKegagalan,
			&se.PotensiKerugiandanDampakNegatif,
			&se.NamaSE,
			&se.IpSE,
			&se.AsNumberSE,
			&se.PengelolaSE,
			&se.FiturSE,
			&se.TotalBobot,
			&se.KategoriSE,
			&se.CreatedAt,
			&se.UpdatedAt,

			&se.Perusahaan.ID,
			&se.Perusahaan.NamaPerusahaan,

			&se.SubSektor.ID,
			&se.SubSektor.NamaSubSektor,
			&se.SubSektor.IDSektor,
			&se.SubSektor.NamaSektor,

			&se.Csirt.ID,
			&se.Csirt.NamaCsirt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, se)
	}

	return result, nil
}

/* ================= GET BY PERUSAHAAN ================= */

func (r *seRepository) GetByPerusahaan(idPerusahaan string) ([]dto.SEResponse, error) {
	rows, err := r.db.Query(`
		SELECT
			se.id,
			se.id_perusahaan,
			se.id_sub_sektor,
			se.id_csirt,
			se.nilai_investasi,
			se.anggaran_operasional,
			se.kepatuhan_peraturan,
			se.teknik_kriptografi,
			se.jumlah_pengguna,
			se.data_pribadi,
			se.klasifikasi_data,
			se.kekritisan_proses,
			se.dampak_kegagalan,
			se.potensi_kerugian_dan_dampak_negatif,
			se.nama_se,
			se.ip_se,
			se.as_number_se,
			se.pengelola_se,
			se.fitur_se,
			se.total_bobot,
			se.kategori_se,
			se.created_at,
			se.updated_at,

			p.id,
			p.nama_perusahaan,

			COALESCE(ss.id, ''),
			COALESCE(ss.nama_sub_sektor, ''),
			COALESCE(s.id, ''),
			COALESCE(s.nama_sektor, ''),

			COALESCE(c.id, ''),
			COALESCE(c.nama_csirt, '')
		FROM se
		JOIN perusahaan p ON se.id_perusahaan = p.id
		LEFT JOIN sub_sektor ss ON se.id_sub_sektor = ss.id
		LEFT JOIN sektor s ON ss.id_sektor = s.id
		LEFT JOIN csirt c ON se.id_csirt = c.id
		WHERE se.id_perusahaan = ?
		ORDER BY se.created_at DESC
	`, idPerusahaan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.SEResponse

	for rows.Next() {
		var se dto.SEResponse
		se.Perusahaan = &dto.PerusahaanMiniResponse{}
		se.SubSektor = &dto.SubSektorMiniResponse{}
		se.Csirt = &dto.CsirtMiniResponse{}

		err := rows.Scan(
			&se.ID,
			&se.IDPerusahaan,
			&se.IDSubSektor,
			&se.IDCsirt,
			&se.NilaiInvestasi,
			&se.AnggaranOperasional,
			&se.KepatuhanPeraturan,
			&se.TeknikKriptografi,
			&se.JumlahPengguna,
			&se.DataPribadi,
			&se.KlasifikasiData,
			&se.KekritisanProses,
			&se.DampakKegagalan,
			&se.PotensiKerugiandanDampakNegatif,
			&se.NamaSE,
			&se.IpSE,
			&se.AsNumberSE,
			&se.PengelolaSE,
			&se.FiturSE,
			&se.TotalBobot,
			&se.KategoriSE,
			&se.CreatedAt,
			&se.UpdatedAt,

			&se.Perusahaan.ID,
			&se.Perusahaan.NamaPerusahaan,

			&se.SubSektor.ID,
			&se.SubSektor.NamaSubSektor,
			&se.SubSektor.IDSektor,
			&se.SubSektor.NamaSektor,

			&se.Csirt.ID,
			&se.Csirt.NamaCsirt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, se)
	}

	return result, nil
}

/* ================= GET BY ID ================= */

func (r *seRepository) GetByID(id string) (*dto.SEResponse, error) {
	row := r.db.QueryRow(`
		SELECT
			se.id,
			se.id_perusahaan,
			se.id_sub_sektor,
			se.id_csirt,
			se.nilai_investasi,
			se.anggaran_operasional,
			se.kepatuhan_peraturan,
			se.teknik_kriptografi,
			se.jumlah_pengguna,
			se.data_pribadi,
			se.klasifikasi_data,
			se.kekritisan_proses,
			se.dampak_kegagalan,
			se.potensi_kerugian_dan_dampak_negatif,
			se.nama_se,
			se.ip_se,
			se.as_number_se,
			se.pengelola_se,
			se.fitur_se,
			se.total_bobot,
			se.kategori_se,
			se.created_at,
			se.updated_at,

			p.id,
			p.nama_perusahaan,

			COALESCE(ss.id, ''),
			COALESCE(ss.nama_sub_sektor, ''),
			COALESCE(s.id, ''),
			COALESCE(s.nama_sektor, ''),

			COALESCE(c.id, ''),
			COALESCE(c.nama_csirt, '')
		FROM se
		JOIN perusahaan p ON se.id_perusahaan = p.id
		LEFT JOIN sub_sektor ss ON se.id_sub_sektor = ss.id
		LEFT JOIN sektor s ON ss.id_sektor = s.id
		LEFT JOIN csirt c ON se.id_csirt = c.id
		WHERE se.id = ?
	`, id)

	var se dto.SEResponse
	se.Perusahaan = &dto.PerusahaanMiniResponse{}
	se.SubSektor = &dto.SubSektorMiniResponse{}
	se.Csirt = &dto.CsirtMiniResponse{}

	err := row.Scan(
		&se.ID,
		&se.IDPerusahaan,
		&se.IDSubSektor,
		&se.IDCsirt,
		&se.NilaiInvestasi,
		&se.AnggaranOperasional,
		&se.KepatuhanPeraturan,
		&se.TeknikKriptografi,
		&se.JumlahPengguna,
		&se.DataPribadi,
		&se.KlasifikasiData,
		&se.KekritisanProses,
		&se.DampakKegagalan,
		&se.PotensiKerugiandanDampakNegatif,
		&se.NamaSE,
		&se.IpSE,
		&se.AsNumberSE,
		&se.PengelolaSE,
		&se.FiturSE,
		&se.TotalBobot,
		&se.KategoriSE,
		&se.CreatedAt,
		&se.UpdatedAt,

		&se.Perusahaan.ID,
		&se.Perusahaan.NamaPerusahaan,

		&se.SubSektor.ID,
		&se.SubSektor.NamaSubSektor,
		&se.SubSektor.IDSektor,
		&se.SubSektor.NamaSektor,

		&se.Csirt.ID,
		&se.Csirt.NamaCsirt,
	)
	if err != nil {
		return nil, err
	}

	return &se, nil
}

/* ================= UPDATE ================= */

func (r *seRepository) Update(
	id string,
	req dto.UpdateSERequest,
	totalBobot int,
	kategori string,
) error {
	query := `UPDATE se SET `
	params := []interface{}{}
	updates := []string{}

	// Update karakteristik instansi
	if req.NilaiInvestasi != nil {
		updates = append(updates, "nilai_investasi = ?")
		params = append(params, *req.NilaiInvestasi)
	}
	if req.AnggaranOperasional != nil {
		updates = append(updates, "anggaran_operasional = ?")
		params = append(params, *req.AnggaranOperasional)
	}
	if req.KepatuhanPeraturan != nil {
		updates = append(updates, "kepatuhan_peraturan = ?")
		params = append(params, *req.KepatuhanPeraturan)
	}
	if req.TeknikKriptografi != nil {
		updates = append(updates, "teknik_kriptografi = ?")
		params = append(params, *req.TeknikKriptografi)
	}
	if req.JumlahPengguna != nil {
		updates = append(updates, "jumlah_pengguna = ?")
		params = append(params, *req.JumlahPengguna)
	}
	if req.DataPribadi != nil {
		updates = append(updates, "data_pribadi = ?")
		params = append(params, *req.DataPribadi)
	}
	if req.KlasifikasiData != nil {
		updates = append(updates, "klasifikasi_data = ?")
		params = append(params, *req.KlasifikasiData)
	}
	if req.KekritisanProses != nil {
		updates = append(updates, "kekritisan_proses = ?")
		params = append(params, *req.KekritisanProses)
	}
	if req.DampakKegagalan != nil {
		updates = append(updates, "dampak_kegagalan = ?")
		params = append(params, *req.DampakKegagalan)
	}
	if req.PotensiKerugiandanDampakNegatif != nil {
		updates = append(updates, "potensi_kerugian_dan_dampak_negatif = ?")
		params = append(params, *req.PotensiKerugiandanDampakNegatif)
	}

	// Update informasi SE
	if req.IDPerusahaan != nil {
		updates = append(updates, "id_perusahaan = ?")
		params = append(params, *req.IDPerusahaan)
	}
	if req.IDSubSektor != nil {
		updates = append(updates, "id_sub_sektor = ?")
		params = append(params, *req.IDSubSektor)
	}
	if req.IDCsirt != nil {
		updates = append(updates, "id_csirt = ?")
		params = append(params, *req.IDCsirt)
	}
	if req.NamaSE != nil {
		updates = append(updates, "nama_se = ?")
		params = append(params, *req.NamaSE)
	}
	if req.IpSE != nil {
		updates = append(updates, "ip_se = ?")
		params = append(params, *req.IpSE)
	}
	if req.AsNumberSE != nil {
		updates = append(updates, "as_number_se = ?")
		params = append(params, *req.AsNumberSE)
	}
	if req.PengelolaSE != nil {
		updates = append(updates, "pengelola_se = ?")
		params = append(params, *req.PengelolaSE)
	}
	if req.FiturSE != nil {
		updates = append(updates, "fitur_se = ?")
		params = append(params, *req.FiturSE)
	}

	// Update hasil kalkulasi
	updates = append(updates, "total_bobot = ?")
	params = append(params, totalBobot)
	updates = append(updates, "kategori_se = ?")
	params = append(params, kategori)
	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")

	// Build query
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = ?"
	params = append(params, id)

	_, err := r.db.Exec(query, params...)
	return err
}

/* ================= DELETE ================= */

func (r *seRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM se WHERE id = ?`, id)
	return err
}