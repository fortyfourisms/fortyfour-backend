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
		)
		if err != nil {
			continue
		}

		if tanggal.Valid {
			i.Tanggal = tanggal.String
		}

		if targetNilai.Valid {
			i.TargetNilai = targetNilai.Float64
		}

		// Selalu hitung nilai_kematangan secara dinamis dari domain yang sudah terisi
		{
			var sum float64
			var count float64
			if idenNilai.Valid {
				sum += idenNilai.Float64
				count++
			}
			if protNilai.Valid {
				sum += protNilai.Float64
				count++
			}
			if detNilai.Valid {
				sum += detNilai.Float64
				count++
			}
			if gulihNilai.Valid {
				sum += gulihNilai.Float64
				count++
			}
			if count > 0 {
				i.NilaiKematangan = utils.RoundToTwo(sum / count)
			}
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
	)
	if err != nil {
		return nil, err
	}

	if tanggal.Valid {
		i.Tanggal = tanggal.String
	}

	if targetNilai.Valid {
		i.TargetNilai = targetNilai.Float64
	}

	// Selalu hitung nilai_kematangan secara dinamis dari domain yang sudah terisi
	{
		var sum float64
		var count float64
		if idenNilai.Valid {
			sum += idenNilai.Float64
			count++
		}
		if protNilai.Valid {
			sum += protNilai.Float64
			count++
		}
		if detNilai.Valid {
			sum += detNilai.Float64
			count++
		}
		if gulihNilai.Valid {
			sum += gulihNilai.Float64
			count++
		}
		if count > 0 {
			i.NilaiKematangan = utils.RoundToTwo(sum / count)
		}
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
	_, err := r.db.Exec(`DELETE FROM ikas WHERE id=?`, id)
	return err
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
func (r *IkasRepository) ParseExcelForImport(fileData []byte) (*dto.CreateIkasRequest, error) {
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()

	// Validasi jumlah sheet
	if len(sheets) < 2 {
		return nil, errors.New("file Excel tidak memiliki sheet ke-2")
	}
	if len(sheets) < 7 {
		return nil, errors.New("file Excel tidak memiliki sheet ke-7")
	}

	sheet2 := sheets[1] // Sheet ke-2 (index 1)
	sheet7 := sheets[6] // Sheet ke-7 (index 6)

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
		if val == "" {
			return 0, fmt.Errorf("cell %s kosong", cell)
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, fmt.Errorf("gagal parse cell %s: %v", cell, err)
		}
		return floatVal, nil
	}

	// ===== AMBIL DATA DARI SHEET 2 =====

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
		return nil, fmt.Errorf("perusahaan dengan nama '%s' tidak ditemukan di database", namaPerusahaan)
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

	// Tanggal dari D18
	tanggalStr, err := getCellString(sheet2, "D18")
	if err != nil {
		return nil, fmt.Errorf("error membaca tanggal (Sheet 2, D18): %v", err)
	}

	// Parse tanggal - support multiple format
	var tanggal string
	if tanggalStr != "" {
		// Coba parse sebagai Excel date number dulu
		if excelDate, err := strconv.ParseFloat(tanggalStr, 64); err == nil {
			// Convert Excel date to time.Time
			parsedTime, err := excelize.ExcelDateToTime(excelDate, false)
			if err == nil {
				tanggal = parsedTime.Format("2006-01-02") // Format MySQL
			}
		} else {
			// Bukan number, coba parse berbagai format string
			parsedTime, err := parseMultipleDateFormats(tanggalStr)
			if err == nil {
				tanggal = parsedTime.Format("2006-01-02") // Format MySQL
			} else {
				return nil, fmt.Errorf("format tanggal tidak valid (Sheet 2, D18): %s. Gunakan format DD-MM-YYYY, DD/MM/YYYY, atau YYYY-MM-DD", tanggalStr)
			}
		}
	} else {
		return nil, errors.New("tanggal (Sheet 2, D18) tidak boleh kosong")
	}
	// ===== AMBIL DATA DARI SHEET 7 =====

	// Target Nilai dari D4
	targetNilai, err := getCellFloat(sheet7, "D4")
	if err != nil {
		return nil, fmt.Errorf("error membaca target_nilai (Sheet 7, D4): %v", err)
	}

	// Construct CreateIkasRequest dengan semua data
	req := &dto.CreateIkasRequest{
		IDPerusahaan: idPerusahaan,
		Tanggal:      tanggal,
		Responden:    responden,
		Telepon:      telepon,
		Jabatan:      jabatan,
		TargetNilai:  targetNilai,
	}

	return req, nil
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
