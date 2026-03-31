package services

import (
	"errors"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/cache"

	"github.com/google/uuid"
)

type SEService interface {
	Create(req dto.CreateSERequest) (*dto.SEResponse, error)
	GetAll() ([]dto.SEResponse, error)
	GetByID(id string) (*dto.SEResponse, error)
	GetByPerusahaan(idPerusahaan string) ([]dto.SEResponse, error)
	Update(id string, req dto.UpdateSERequest) (*dto.SEResponse, error)
	Delete(id string) error
}

type seService struct {
	repo repository.SERepositoryInterface
	rc   cache.RedisInterface
}

func NewSEService(repo repository.SERepositoryInterface, rc cache.RedisInterface) SEService {
	return &seService{repo: repo, rc: rc}
}

/* =======================
   CREATE
======================= */

func (s *seService) Create(req dto.CreateSERequest) (*dto.SEResponse, error) {
	// Validasi ID Perusahaan
	if strings.TrimSpace(req.IDPerusahaan) == "" {
		return nil, errors.New("id_perusahaan wajib diisi")
	}

	// Validasi informasi SE
	if strings.TrimSpace(req.NamaSE) == "" {
		return nil, errors.New("nama_se wajib diisi")
	}
	if strings.TrimSpace(req.IpSE) == "" {
		return nil, errors.New("ip_se wajib diisi")
	}
	if strings.TrimSpace(req.AsNumberSE) == "" {
		return nil, errors.New("as_number_se wajib diisi")
	}
	if strings.TrimSpace(req.PengelolaSE) == "" {
		return nil, errors.New("pengelola_se wajib diisi")
	}

	// Hitung total bobot dan kategori
	totalBobot, err := hitungTotalBobotCreate(req)
	if err != nil {
		return nil, err
	}

	kategori := hitungKategoriSE(totalBobot)
	if kategori == "" {
		return nil, errors.New("total bobot tidak valid untuk kategorisasi")
	}

	id := uuid.NewString()

	if err := s.repo.Create(req, id, totalBobot, kategori); err != nil {
		return nil, err
	}

	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, keyDetail("se", id), result, TTLDetail)
	cacheDelete(s.rc, keyList("se"))
	cacheDelete(s.rc, "se:perusahaan:"+req.IDPerusahaan)

	return result, nil
}

/* =======================
   READ
======================= */

func (s *seService) GetAll() ([]dto.SEResponse, error) {
	key := keyList("se")
	var result []dto.SEResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	result, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, result, TTLList)
	return result, nil
}

func (s *seService) GetByID(id string) (*dto.SEResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id wajib diisi")
	}
	key := keyDetail("se", id)
	var result dto.SEResponse
	if cacheGet(s.rc, key, &result) {
		return &result, nil
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, data, TTLDetail)
	return data, nil
}

func (s *seService) GetByPerusahaan(idPerusahaan string) ([]dto.SEResponse, error) {
	if strings.TrimSpace(idPerusahaan) == "" {
		return nil, errors.New("id_perusahaan wajib diisi")
	}
	key := "se:perusahaan:" + idPerusahaan
	var result []dto.SEResponse
	if cacheGet(s.rc, key, &result) {
		return result, nil
	}

	data, err := s.repo.GetByPerusahaan(idPerusahaan)
	if err != nil {
		return nil, err
	}

	cacheSet(s.rc, key, data, TTLList)
	return data, nil
}

/* =======================
   UPDATE
======================= */

func (s *seService) Update(id string, req dto.UpdateSERequest) (*dto.SEResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id wajib diisi")
	}

	// Get existing data untuk kalkulasi ulang
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("data tidak ditemukan")
	}

	// Hitung total bobot berdasarkan data yang di-update
	totalBobot, err := hitungTotalBobotUpdate(req, existing)
	if err != nil {
		return nil, err
	}

	kategori := hitungKategoriSE(totalBobot)
	if kategori == "" {
		return nil, errors.New("total bobot tidak valid untuk kategorisasi")
	}

	if err := s.repo.Update(id, req, totalBobot, kategori); err != nil {
		return nil, err
	}

	cacheDelete(s.rc, keyDetail("se", id))
	cacheDelete(s.rc, keyList("se"))
	cacheDelete(s.rc, "se:perusahaan:"+existing.IDPerusahaan)

	existing.TotalBobot = totalBobot
	existing.KategoriSE = kategori

	if req.NamaSE != nil {
		existing.NamaSE = *req.NamaSE
	}
	if req.FiturSE != nil {
		existing.FiturSE = *req.FiturSE
	}
	if req.IpSE != nil {
		existing.IpSE = *req.IpSE
	}
	if req.AsNumberSE != nil {
		existing.AsNumberSE = *req.AsNumberSE
	}
	if req.PengelolaSE != nil {
		existing.PengelolaSE = *req.PengelolaSE
	}

	return existing, nil
}

/* =======================
   DELETE
======================= */

func (s *seService) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id wajib diisi")
	}

	// Ambil data dulu untuk invalidate cache per perusahaan
	existing, _ := s.repo.GetByID(id)

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	cacheDelete(s.rc, keyDetail("se", id))
	cacheDelete(s.rc, keyList("se"))
	if existing != nil {
		cacheDelete(s.rc, "se:perusahaan:"+existing.IDPerusahaan)
	}

	return nil
}

/* ======================================================
   HELPER FUNCTIONS (A/B/C → BOBOT)
====================================================== */

func jawabanKeBobot(jawaban string) (int, error) {
	switch strings.ToUpper(strings.TrimSpace(jawaban)) {
	case "A":
		return 5, nil
	case "B":
		return 2, nil
	case "C":
		return 1, nil
	default:
		return 0, errors.New("jawaban harus A, B, atau C")
	}
}

func hitungTotalBobotCreate(req dto.CreateSERequest) (int, error) {
	karakteristik := []string{
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
	}

	total := 0
	for i, k := range karakteristik {
		bobot, err := jawabanKeBobot(k)
		if err != nil {
			return 0, errors.New("karakteristik " + string(rune(i+1)) + ": " + err.Error())
		}
		total += bobot
	}
	return total, nil
}

func hitungTotalBobotUpdate(req dto.UpdateSERequest, existing *dto.SEResponse) (int, error) {
	// Gunakan nilai existing jika tidak di-update
	karakteristik := []string{
		getStringValue(req.NilaiInvestasi, existing.NilaiInvestasi),
		getStringValue(req.AnggaranOperasional, existing.AnggaranOperasional),
		getStringValue(req.KepatuhanPeraturan, existing.KepatuhanPeraturan),
		getStringValue(req.TeknikKriptografi, existing.TeknikKriptografi),
		getStringValue(req.JumlahPengguna, existing.JumlahPengguna),
		getStringValue(req.DataPribadi, existing.DataPribadi),
		getStringValue(req.KlasifikasiData, existing.KlasifikasiData),
		getStringValue(req.KekritisanProses, existing.KekritisanProses),
		getStringValue(req.DampakKegagalan, existing.DampakKegagalan),
		getStringValue(req.PotensiKerugiandanDampakNegatif, existing.PotensiKerugiandanDampakNegatif),
	}

	total := 0
	for i, k := range karakteristik {
		bobot, err := jawabanKeBobot(k)
		if err != nil {
			return 0, errors.New("karakteristik " + string(rune(i+1)) + ": " + err.Error())
		}
		total += bobot
	}
	return total, nil
}

func getStringValue(newValue *string, existingValue string) string {
	if newValue != nil {
		return *newValue
	}
	return existingValue
}

/* =======================
   KATEGORISASI SE
======================= */

func hitungKategoriSE(totalBobot int) string {
	switch {
	case totalBobot >= 35 && totalBobot <= 50:
		return "Strategis"
	case totalBobot >= 16 && totalBobot <= 34:
		return "Tinggi"
	case totalBobot >= 10 && totalBobot <= 15:
		return "Rendah"
	default:
		return ""
	}
}
