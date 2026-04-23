package repository

import (
	"database/sql"
	"errors"
	"testing"

	"ikas/internal/dto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func TestIkasRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewIkasRepository(db)
	req := dto.CreateIkasRequest{
		IDPerusahaan: "comp-1",
		Tanggal:      "2023-01-01",
		Responden:    "John Doe",
		Telepon:      "123456",
		Jabatan:      "Manager",
		TargetNilai:  4.5,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO ikas").
			WithArgs("ikas-1", req.IDPerusahaan, req.Tanggal, req.Responden, req.Telepon, req.Jabatan, 3.5, req.TargetNilai).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(req, "ikas-1", 3.5)
		assert.NoError(t, err)
	})

	t.Run("EmptyTanggal_Success", func(t *testing.T) {
		reqEmpty := req
		reqEmpty.Tanggal = ""
		mock.ExpectExec("INSERT INTO ikas").
			WithArgs("ikas-1", reqEmpty.IDPerusahaan, nil, reqEmpty.Responden, reqEmpty.Telepon, reqEmpty.Jabatan, 3.5, reqEmpty.TargetNilai).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(reqEmpty, "ikas-1", 3.5)
		assert.NoError(t, err)
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO ikas").WillReturnError(errors.New("db error"))
		err := repo.Create(req, "ikas-1", 3.5)
		assert.Error(t, err)
	})
}

func TestIkasRepository_GetQueries(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewIkasRepository(db)

	columns := []string{
		"id", "tanggal", "responden", "telepon", "jabatan", "nilai_kematangan", "target_nilai",
		"p_id", "p_nama",
		"iden_id", "iden_nilai", "iden_s1", "iden_s2", "iden_s3", "iden_s4", "iden_s5",
		"prot_id", "prot_nilai", "prot_s1", "prot_s2", "prot_s3", "prot_s4", "prot_s5", "prot_s6",
		"det_id", "det_nilai", "det_s1", "det_s2", "det_s3",
		"g_id", "g_nilai", "g_s1", "g_s2", "g_s3", "g_s4",
		"is_validated", "created_at", "updated_at",
	}

	t.Run("GetAll_Success", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow("ikas-1", "2023-01-01", "John", "123", "Manager", 0.0, 4.0,
				"comp-1", "Comp One",
				1, 4.0, 4, 4, 4, 4, 4,
				2, 4.0, 4, 4, 4, 4, 4, 4,
				3, 4.0, 4, 4, 4,
				4, 4.0, 4, 4, 4, 4,
				true, "2023-01-01", "2023-01-01").
			AddRow("ikas-2", nil, "Jane", "456", "Staff", 0.0, nil,
				nil, nil,
				nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil,
				false, "2023-01-02", nil)

		mock.ExpectQuery("SELECT (.+) FROM ikas i").WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "ikas-1", result[0].ID)
		assert.Equal(t, 4.0, result[0].NilaiKematangan) // Calculated: 4*0.25 + 4*0.3 + 4*0.25 + 4*0.2 = 4.0
		assert.Equal(t, "Comp One", result[0].Perusahaan.NamaPerusahaan)
		assert.True(t, result[0].IsValidated)

		assert.Equal(t, "ikas-2", result[1].ID)
		assert.False(t, result[1].IsValidated)
		assert.Nil(t, result[1].Perusahaan)
	})

	t.Run("GetByPerusahaan_Success", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow("ikas-1", "2023-01-01", "John", "123", "Manager", 0, 4.0,
				"comp-1", "Comp One",
				1, 4.0, 4, 4, 4, 4, 4,
				2, 4.0, 4, 4, 4, 4, 4, 4,
				3, 4.0, 4, 4, 4,
				4, 4.0, 4, 4, 4, 4,
				true, "2023-01-01", "2023-01-01")

		mock.ExpectQuery("SELECT (.+) FROM ikas i (.+) WHERE i.id_perusahaan = ?").
			WithArgs("comp-1").
			WillReturnRows(rows)

		result, err := repo.GetByPerusahaan("comp-1")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow("ikas-1", "2023-01-01", "John", "123", "Manager", 0, 4.0,
				"comp-1", "Comp One",
				1, 4.0, 4, 4, 4, 4, 4,
				2, 4.0, 4, 4, 4, 4, 4, 4,
				3, 4.0, 4, 4, 4,
				4, 4.0, 4, 4, 4, 4,
				true, "2023-01-01", "2023-01-01")

		mock.ExpectQuery("SELECT (.+) FROM ikas i (.+) WHERE i.id = ?").
			WithArgs("ikas-1").
			WillReturnRows(rows)

		result, err := repo.GetByID("ikas-1")
		assert.NoError(t, err)
		assert.Equal(t, "ikas-1", result.ID)
	})

	t.Run("GetQueries_Errors", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM ikas i").WillReturnError(errors.New("query error"))
		res1, err1 := repo.GetAll()
		assert.Error(t, err1)
		assert.Nil(t, res1)

		mock.ExpectQuery("SELECT (.+) FROM ikas i").WillReturnError(errors.New("query error"))
		res2, err2 := repo.GetByPerusahaan("c1")
		assert.Error(t, err2)
		assert.Nil(t, res2)

		mock.ExpectQuery("SELECT (.+) FROM ikas i").WillReturnError(sql.ErrNoRows)
		res3, err3 := repo.GetByID("i1")
		assert.Error(t, err3)
		assert.Nil(t, res3)
	})
}

func TestIkasRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewIkasRepository(db)

	t.Run("FullUpdate_Success", func(t *testing.T) {
		pID := "comp-1"
		date := "2023-02-01"
		resp := "Jane"
		telp := "999"
		jab := "VP"
		target := 5.0
		req := dto.UpdateIkasRequest{
			IDPerusahaan: &pID,
			Tanggal:      &date,
			Responden:    &resp,
			Telepon:      &telp,
			Jabatan:      &jab,
			TargetNilai:  &target,
		}

		mock.ExpectExec("UPDATE ikas SET id_perusahaan=\\?, tanggal=\\?, responden=\\?, telepon=\\?, jabatan=\\?, target_nilai=\\? WHERE id=\\?").
			WithArgs(pID, date, resp, telp, jab, target, "ikas-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update("ikas-1", req)
		assert.NoError(t, err)
	})

	t.Run("NoUpdates", func(t *testing.T) {
		err := repo.Update("ikas-1", dto.UpdateIkasRequest{})
		assert.NoError(t, err)
	})
}

func TestIkasRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewIkasRepository(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id_identifikasi, id_proteksi, id_deteksi, id_gulih FROM ikas WHERE id = ?").
			WithArgs("ikas-1").
			WillReturnRows(sqlmock.NewRows([]string{"id_iden", "id_prot", "id_det", "id_gulih"}).
				AddRow(1, 2, 3, 4))

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM jawaban_identifikasi_buffer WHERE ikas_id = ?").WithArgs("ikas-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM jawaban_proteksi_buffer WHERE ikas_id = ?").WithArgs("ikas-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM jawaban_deteksi_buffer WHERE ikas_id = ?").WithArgs("ikas-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM jawaban_gulih_buffer WHERE ikas_id = ?").WithArgs("ikas-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM jawaban_identifikasi WHERE ikas_id = ?").WithArgs("ikas-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM jawaban_proteksi WHERE ikas_id = ?").WithArgs("ikas-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM jawaban_deteksi WHERE ikas_id = ?").WithArgs("ikas-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM jawaban_gulih WHERE ikas_id = ?").WithArgs("ikas-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM ikas WHERE id = ?").WithArgs("ikas-1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM identifikasi WHERE id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM proteksi WHERE id = ?").WithArgs(2).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM deteksi WHERE id = ?").WithArgs(3).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM gulih WHERE id = ?").WithArgs(4).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete("ikas-1")
		assert.NoError(t, err)
	})

	t.Run("NotFound_Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM ikas").WillReturnError(sql.ErrNoRows)
		err := repo.Delete("non-existent")
		assert.NoError(t, err)
	})

	t.Run("TxError", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM ikas").WillReturnRows(sqlmock.NewRows([]string{"i", "p", "d", "g"}).AddRow(1, 1, 1, 1))
		mock.ExpectBegin().WillReturnError(errors.New("tx error"))
		err := repo.Delete("i1")
		assert.Error(t, err)
	})

	t.Run("ExecError_Rollback", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM ikas").WillReturnRows(sqlmock.NewRows([]string{"i", "p", "d", "g"}).AddRow(1, 1, 1, 1))
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM (.+)").WillReturnError(errors.New("exec error"))
		mock.ExpectRollback()
		err := repo.Delete("i1")
		assert.Error(t, err)
	})
}

func TestIkasRepository_Excel(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewIkasRepository(db)

	generateMockExcel := func() []byte {
		f := excelize.NewFile()
		// Sheet list currently has "Sheet1"
		f.NewSheet("Sheet2")
		f.NewSheet("Sheet3")
		f.NewSheet("Sheet4")
		f.NewSheet("Sheet5")
		f.NewSheet("Sheet6")

		// Sheet 2: Info Dasar
		f.SetCellValue("Sheet2", "D4", "Company Alpha")
		f.SetCellValue("Sheet2", "D10", "08123")
		f.SetCellValue("Sheet2", "D11", "John")
		f.SetCellValue("Sheet2", "D12", "CTO")
		f.SetCellValue("Sheet2", "D15", "4.0")
		f.SetCellValue("Sheet2", "D18", "01-01-2023")

		// Sheet 3: Identifikasi
		f.SetCellValue("Sheet3", "D5", "4.0")

		buffer, _ := f.WriteToBuffer()
		return buffer.Bytes()
	}

	t.Run("ParseExcel_Success", func(t *testing.T) {
		excelData := generateMockExcel()

		// Mock FindPerusahaanByName
		mock.ExpectQuery("SELECT id FROM perusahaan WHERE LOWER").
			WithArgs("Company Alpha").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("comp-alpha"))

		// Mock CheckExistsByPerusahaanIDAndYear
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ikas WHERE id_perusahaan = \\? AND YEAR").
			WithArgs("comp-alpha", 2023).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result, err := repo.ParseExcelForImport(excelData)
		assert.NoError(t, err)
		assert.Equal(t, "comp-alpha", result.IkasRequest.IDPerusahaan)
		assert.Equal(t, "2023-01-01", result.IkasRequest.Tanggal) 
	})

	t.Run("ParseExcel_CompanyNotFound", func(t *testing.T) {
		f := excelize.NewFile()
		f.NewSheet("Sheet2"); f.NewSheet("Sheet3"); f.NewSheet("Sheet4"); f.NewSheet("Sheet5"); f.NewSheet("Sheet6")
		f.SetCellValue("Sheet2", "D4", "Unknown")
		buf, _ := f.WriteToBuffer()

		mock.ExpectQuery("SELECT id FROM perusahaan").WillReturnRows(sqlmock.NewRows([]string{"id"}))

		res, err := repo.ParseExcelForImport(buf.Bytes())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "belum terdaftar")
		assert.Nil(t, res)
	})

	t.Run("ParseExcel_DuplicateData", func(t *testing.T) {
		excelData := generateMockExcel()
		mock.ExpectQuery("SELECT id FROM perusahaan").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c1"))
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(1))

		res, err := repo.ParseExcelForImport(excelData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sudah ada")
		assert.Nil(t, res)
	})

	t.Run("ParseExcel_FindPerusahaanError", func(t *testing.T) {
		excelData := generateMockExcel()
		mock.ExpectQuery("SELECT id FROM perusahaan").WillReturnError(errors.New("db error"))
		res, err := repo.ParseExcelForImport(excelData)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("ParseExcel_CheckExistsError", func(t *testing.T) {
		excelData := generateMockExcel()
		mock.ExpectQuery("SELECT id FROM perusahaan").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c1"))
		mock.ExpectQuery("SELECT COUNT").WillReturnError(errors.New("db error"))
		res, err := repo.ParseExcelForImport(excelData)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("ParseExcel_InvalidData", func(t *testing.T) {
		f := excelize.NewFile()
		f.NewSheet("Sheet2")
		f.NewSheet("Sheet3")
		f.NewSheet("Sheet4")
		f.NewSheet("Sheet5")
		f.NewSheet("Sheet6")
		f.SetCellValue("Sheet2", "D4", "Comp")
		f.SetCellValue("Sheet2", "D18", "invalid-date")
		buf, _ := f.WriteToBuffer()

		mock.ExpectQuery("SELECT id FROM perusahaan").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c1"))
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(0))

		res, err := repo.ParseExcelForImport(buf.Bytes())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "format tanggal")
		assert.Nil(t, res)
	})

	t.Run("ParseExcel_InvalidFloat", func(t *testing.T) {
		f := excelize.NewFile()
		f.NewSheet("Sheet2")
		f.NewSheet("Sheet3")
		f.NewSheet("Sheet4")
		f.NewSheet("Sheet5")
		f.NewSheet("Sheet6")
		f.SetCellValue("Sheet2", "D4", "Comp")
		f.SetCellValue("Sheet2", "D18", "01-01-2023")
		f.SetCellValue("Sheet2", "D15", "not-a-number")
		buf, _ := f.WriteToBuffer()

		mock.ExpectQuery("SELECT id FROM perusahaan").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c1"))
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(0))

		res, err := repo.ParseExcelForImport(buf.Bytes())
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("ParseExcel_MissingSheets", func(t *testing.T) {
		f := excelize.NewFile()
		buf, _ := f.WriteToBuffer()
		res, err := repo.ParseExcelForImport(buf.Bytes())
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("ImportFromExcel_Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			repo.ImportFromExcel(nil)
		})
	})
}

func TestIkasRepository_Utilities(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewIkasRepository(db)

	t.Run("CheckExistsByPerusahaanID", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").WithArgs("c1").WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(1))
		exists, err := repo.CheckExistsByPerusahaanID("c1")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CheckExists_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").WillReturnError(errors.New("db error"))
		exists, err := repo.CheckExistsByPerusahaanID("c1")
		assert.Error(t, err)
		assert.False(t, exists)
	})

	t.Run("CheckOwnership", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").WithArgs("i1", "p1").WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(1))
		owned, err := repo.CheckOwnership("i1", "p1")
		assert.NoError(t, err)
		assert.True(t, owned)
	})

	t.Run("CheckOwnership_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").WillReturnError(errors.New("db error"))
		owned, err := repo.CheckOwnership("i1", "p1")
		assert.Error(t, err)
		assert.False(t, owned)
	})

	t.Run("IsLocked", func(t *testing.T) {
		mock.ExpectQuery("SELECT is_validated").WithArgs("i1").WillReturnRows(sqlmock.NewRows([]string{"v"}).AddRow(true))
		locked, err := repo.IsLocked("i1")
		assert.NoError(t, err)
		assert.True(t, locked)
	})

	t.Run("IsLocked_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT is_validated").WillReturnError(errors.New("db error"))
		locked, err := repo.IsLocked("i1")
		assert.Error(t, err)
		assert.False(t, locked)
	})

	t.Run("GetLatestByPerusahaan", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, tanggal").WithArgs("c1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "t", "r", "te", "j", "nk", "tn", "iv"}).
				AddRow("i1", "2023-01-01", "R", "T", "J", 4.0, 5.0, true))
		res, err := repo.GetLatestByPerusahaan("c1")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "i1", res.ID)
	})

	t.Run("GetLatestByPerusahaan_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, tanggal").WillReturnError(errors.New("db error"))
		res, err := repo.GetLatestByPerusahaan("c1")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("GetIDByPerusahaanID", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM ikas").WithArgs("p1").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("i1"))
		id, err := repo.GetIDByPerusahaanID("p1")
		assert.NoError(t, err)
		assert.Equal(t, "i1", id)
	})

	t.Run("UpdateValidationStatus", func(t *testing.T) {
		mock.ExpectExec("UPDATE ikas SET is_validated").WithArgs(true, "i1").WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.UpdateValidationStatus("i1", true)
		assert.NoError(t, err)
	})

	t.Run("CreateInitial", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO ikas").WithArgs("i2", "2023-01-01", "i1").WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.CreateInitial("i1", "i2", "2023-01-01")
		assert.NoError(t, err)
	})

	t.Run("UpdateDomainLinks", func(t *testing.T) {
		mock.ExpectExec("UPDATE ikas SET id_identifikasi").WithArgs("iden", "prot", "det", "gul", "ikas").WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.UpdateDomainLinks("ikas", "iden", "prot", "det", "gul")
		assert.NoError(t, err)
	})
}

func TestParseMultipleDateFormats(t *testing.T) {
	// Small trick to test the private function if needed, but here it is in the same package
	t.Run("VariousFormats", func(t *testing.T) {
		formats := []string{
			"02-01-2023",
			"02/01/2023",
			"02-01-23",
			"02/01/23",
			"2023-01-02",
			"2023/01/02",
		}
		for _, f := range formats {
			_, err := parseMultipleDateFormats(f)
			assert.NoError(t, err, "Failed to parse: "+f)
		}
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		_, err := parseMultipleDateFormats("invalid")
		assert.Error(t, err)
	})
}

