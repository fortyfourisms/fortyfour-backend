package services

import (
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"ikas/internal/repository"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type RuangLingkupService struct {
	repo repository.RuangLingkupRepositoryInterface
}

func NewRuangLingkupService(repo repository.RuangLingkupRepositoryInterface) *RuangLingkupService {
	return &RuangLingkupService{repo: repo}
}

// Normalisasi input: trim dan hilangkan multiple spaces
func normalizeInput(input string) string {
	input = strings.TrimSpace(input)

	multipleSpaces := regexp.MustCompile(`\s+`)
	input = multipleSpaces.ReplaceAllString(input, " ")

	return input
}

// Validasi karakter berbahaya untuk SQL Injection
func containsSQLInjectionPattern(input string) bool {
	// Pattern umum SQL injection
	dangerousPatterns := []string{
		"--",         // SQL comment
		";",          // Multiple statements
		"'",          // String delimiter
		"\"",         // String delimiter
		"/*",         // Multi-line comment
		"*/",         // Multi-line comment
		"xp_",        // Extended stored procedures
		"sp_",        // Stored procedures
		"exec",       // Execute command
		"execute",    // Execute command
		"drop",       // Drop command
		"insert",     // Insert command
		"delete",     // Delete command
		"update",     // Update command
		"union",      // Union query
		"select",     // Select command
		"create",     // Create command
		"alter",      // Alter command
		"shutdown",   // Shutdown command
		"script",     // Script tag
		"javascript", // JavaScript
		"<script",    // Script tag
		"</script>",  // Script tag
		"onerror",    // Event handler
		"onload",     // Event handler
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}
	return false
}

// Validasi karakter yang diizinkan
func isValidInput(input string) bool {
	// Hanya izinkan: huruf, angka, spasi, dan beberapa karakter khusus umum
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,()&]+$`)
	return validPattern.MatchString(input)
}

// Validasi UUID format untuk mencegah injection via ID
func isValidUUID(id string) bool {
	uuidPattern := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	return uuidPattern.MatchString(id)
}

// Validasi untuk Create
func (s *RuangLingkupService) validateCreate(req *dto.CreateRuangLingkupRequest) error {
	// Normalisasi: trim whitespace + hilangkan multiple spaces
	req.NamaRuangLingkup = normalizeInput(req.NamaRuangLingkup)

	// NOT NULL: tidak boleh kosong
	if req.NamaRuangLingkup == "" {
		return errors.New("nama_ruang_lingkup tidak boleh kosong")
	}

	// Min karakter
	if len(req.NamaRuangLingkup) < 3 {
		return errors.New("nama_ruang_lingkup minimal 3 karakter")
	}

	// Max karakter
	if len(req.NamaRuangLingkup) > 50 {
		return errors.New("nama_ruang_lingkup maksimal 50 karakter")
	}

	// Validasi SQL Injection pattern (blacklist)
	if containsSQLInjectionPattern(req.NamaRuangLingkup) {
		return errors.New("nama_ruang_lingkup mengandung karakter yang tidak diizinkan")
	}

	// Validasi karakter yang diizinkan
	if !isValidInput(req.NamaRuangLingkup) {
		return errors.New("nama_ruang_lingkup hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
	}

	return nil
}

// Validasi untuk Update
func (s *RuangLingkupService) validateUpdate(req *dto.UpdateRuangLingkupRequest) error {
	// Jika field dikirim (bukan nil), lakukan validasi
	if req.NamaRuangLingkup != nil {
		// Normalisasi: trim whitespace + hilangkan multiple spaces
		normalized := normalizeInput(*req.NamaRuangLingkup)
		req.NamaRuangLingkup = &normalized

		// NOT NULL: tidak boleh string kosong
		if *req.NamaRuangLingkup == "" {
			return errors.New("nama_ruang_lingkup tidak boleh kosong")
		}

		// Min karakter
		if len(*req.NamaRuangLingkup) < 3 {
			return errors.New("nama_ruang_lingkup minimal 3 karakter")
		}

		// Max karakter
		if len(*req.NamaRuangLingkup) > 50 {
			return errors.New("nama_ruang_lingkup maksimal 50 karakter")
		}

		// Validasi SQL Injection pattern (blacklist)
		if containsSQLInjectionPattern(*req.NamaRuangLingkup) {
			return errors.New("nama_ruang_lingkup mengandung karakter yang tidak diizinkan")
		}

		// Validasi karakter yang diizinkan (whitelist - lebih ketat)
		if !isValidInput(*req.NamaRuangLingkup) {
			return errors.New("nama_ruang_lingkup hanya boleh mengandung huruf, angka, spasi, dan karakter -_.,()&")
		}
	}

	return nil
}

func (s *RuangLingkupService) Create(req dto.CreateRuangLingkupRequest) (*dto.RuangLingkupResponse, error) {
	// Validasi input
	if err := s.validateCreate(&req); err != nil {
		return nil, err
	}

	// Cek duplikasi data (case-insensitive, whitespace-trimmed)
	isDuplicate, err := s.repo.CheckDuplicateName(req.NamaRuangLingkup, "")
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("nama_ruang_lingkup sudah ada")
	}

	// Generate UUID
	newID := uuid.New().String()

	if err := s.repo.Create(req, newID); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Ambil data yang baru dibuat
	resp, err := s.repo.GetByID(newID)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return resp, nil
}

func (s *RuangLingkupService) GetAll() ([]dto.RuangLingkupResponse, error) {
	return s.repo.GetAll()
}

func (s *RuangLingkupService) GetByID(id string) (*dto.RuangLingkupResponse, error) {
	// Validasi format UUID untuk mencegah SQL injection via ID
	if !isValidUUID(id) {
		return nil, errors.New("format ID tidak valid")
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}
	return data, nil
}

func (s *RuangLingkupService) Update(id string, req dto.UpdateRuangLingkupRequest) (*dto.RuangLingkupResponse, error) {
	// Validasi format UUID untuk mencegah SQL injection via ID
	if !isValidUUID(id) {
		return nil, errors.New("format ID tidak valid")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		if err == sql.ErrNoRows {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	// Validasi input
	if err := s.validateUpdate(&req); err != nil {
		return nil, err
	}

	// Cek duplikasi nama
	if req.NamaRuangLingkup != nil {
		isDuplicate, err := s.repo.CheckDuplicateName(*req.NamaRuangLingkup, id)
		if err != nil {
			rollbar.Error(err)
			return nil, err
		}
		if isDuplicate {
			return nil, errors.New("nama_ruang_lingkup sudah ada")
		}
	}

	// Update
	if err := s.repo.Update(id, req); err != nil {
		rollbar.Error(err)
		return nil, err
	}

	// Ambil data terbaru
	updated, err := s.repo.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *RuangLingkupService) Delete(id string) error {
	// Validasi format UUID untuk mencegah SQL injection via ID
	if !isValidUUID(id) {
		return errors.New("format ID tidak valid")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	return s.repo.Delete(id)
}
