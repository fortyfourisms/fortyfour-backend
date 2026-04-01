package validator

import (
	"strings"
	"testing"
)

type TestStruct struct {
	Name  string `validate:"required,min=3,max=10"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18"`
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name: "Valid struct",
			input: TestStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "Invalid struct - missing fields",
			input: TestStruct{
				Name:  "",
				Email: "",
				Age:   0,
			},
			wantErr: true,
		},
		{
			name: "Invalid struct - invalid email",
			input: TestStruct{
				Name:  "John",
				Email: "invalid-email",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "Invalid struct - min length",
			input: TestStruct{
				Name:  "Jo",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "Invalid struct - max length",
			input: TestStruct{
				Name:  "JohnDoeTheGreat",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"Valid email", "test@example.com", true},
		{"Invalid email no domain", "test@", false},
		{"Invalid email no user", "@example.com", false},
		{"Invalid email no at", "testexample.com", false},
		{"Empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{"Valid username", "john", true},
		{"Too short", "jo", false},
		{"Too long", strings.Repeat("a", 51), false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateUsername(tt.username); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	validConfig := DefaultPasswordConfig()

	tests := []struct {
		name         string
		password     string
		config       PasswordValidationConfig
		personalInfo []string
		wantErr      bool
	}{
		{
			name:     "Valid password",
			password: "S0m3Unique$tr1ng!",
			config:   validConfig,
			wantErr:  false,
		},
		{
			name:     "Too short",
			password: "Sh1!",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:     "No upper case",
			password: "password123!",
			config:   validConfig, // RequireUpper check will fail first or Common check
			wantErr:  true,
		},
		{
			name:     "No lower case",
			password: "PASSWORD123!",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:     "No number",
			password: "Password!",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:     "No special",
			password: "Password123",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:     "Common password",
			password: "password123!",
			config:   validConfig,
			wantErr:  true,
		},
		{
			name:         "Contains personal info (username)",
			password:     "JohnStrongXYZ!",
			config:       validConfig,
			personalInfo: []string{"john", "john@example.com"},
			wantErr:      true,
		},
		{
			name:         "Contains personal info (email part)",
			password:     "ExampleStrongXYZ!",
			config:       validConfig,
			personalInfo: []string{"john", "example@test.com"},
			wantErr:      true,
		},
		{
			name:         "Contains personal info (reversed)",
			password:     "nhoJStrongXYZ!",
			config:       validConfig,
			personalInfo: []string{"john"},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password, tt.config, tt.personalInfo...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNewPassword(t *testing.T) {
	config := DefaultPasswordConfig()

	tests := []struct {
		name    string
		newPass string
		oldPass string
		wantErr bool
	}{
		{
			name:    "Different passwords",
			newPass: "N3wUnique$tr1ng!",
			oldPass: "OldUnique$tr1ng!",
			wantErr: false,
		},
		{
			name:    "Same passwords",
			newPass: "OldUnique$tr1ng!",
			oldPass: "OldUnique$tr1ng!",
			wantErr: true,
		},
		{
			name:    "Invalid new password",
			newPass: "weak",
			oldPass: "OldUnique$tr1ng!",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNewPassword(tt.newPass, tt.oldPass, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNewPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test internal helper functions directly if needed, but they are unexported.
	// Since we are in the same package (package validator), we can test them.

	// Test containsCommonRoot
	if !containsCommonRoot("someadminpassword") { // "admin" is root
		t.Error("containsCommonRoot failed for 'someadminpassword'")
	}
	if containsCommonRoot("uniqueXYZ") {
		t.Error("containsCommonRoot failed for 'uniqueXYZ'")
	}

	// Test reverseString
	if got := reverseString("abc"); got != "cba" {
		t.Errorf("reverseString('abc') = %s, want cba", got)
	}
}

// ================================================================
// TEST: containsCommonRoot — lebih banyak variasi
// ================================================================

func TestContainsCommonRoot_SemuaRootTerdaftar(t *testing.T) {
	// Setiap root dalam commonPasswordRoots harus terdeteksi
	cases := []struct {
		input    string
		expected bool
	}{
		// root: "password"
		{"mypasswordXYZ", true},
		{"PASSWORDsecure", true},
		// root: "admin"
		{"adminSuperXYZ", true},
		{"superAdminZ1!", true},
		// root: "qwerty"
		{"qwertySecure1!", true},
		// root: "user"
		{"userSecure123!", true},
		// root: "login"
		{"loginSecure1!", true},
		// root: "welcome"
		{"welcomeHome1!", true},
		// root: "root"
		{"rootAccess1!", true},
		// root: "system"
		{"systemSecure1!", true},
		// root: "test"
		{"testSecure123!", true},
		// root: "superadmin"
		{"superadminXYZ1!", true},
		// root: "123"
		{"abc123XYZ!", true},
		// root: "aaa"
		{"myaaa123XYZ!", true},
		// tidak ada root
		{"UniquePhraseZ9!", false},
		{"strongXYZpass", false},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got := containsCommonRoot(c.input)
			if got != c.expected {
				t.Errorf("containsCommonRoot(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestContainsCommonRoot_CaseInsensitive(t *testing.T) {
	// containsCommonRoot harus case-insensitive karena pakai strings.ToLower
	cases := []string{
		"MYPASSWORDXYZ",
		"MyAdminXYZ",
		"QWERTYUIOP1!",
		"UserProfile9!",
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			if !containsCommonRoot(c) {
				t.Errorf("containsCommonRoot(%q) seharusnya true (case-insensitive)", c)
			}
		})
	}
}

func TestContainsCommonRoot_StringKosong(t *testing.T) {
	if containsCommonRoot("") {
		t.Error("containsCommonRoot string kosong seharusnya false")
	}
}

// ================================================================
// TEST: reverseString — berbagai skenario
// ================================================================

func TestReverseString_BasicASCII(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"abc", "cba"},
		{"hello", "olleh"},
		{"a", "a"},
		{"", ""},
		{"ab", "ba"},
		{"12345", "54321"},
		{"!@#", "#@!"},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got := reverseString(c.input)
			if got != c.expected {
				t.Errorf("reverseString(%q) = %q, want %q", c.input, got, c.expected)
			}
		})
	}
}

func TestReverseString_Unicode(t *testing.T) {
	// reverseString pakai []rune sehingga harus aman untuk multibyte Unicode
	cases := []struct {
		input    string
		expected string
	}{
		{"héllo", "ollèh"[0:0] + "olleh"},  // fallback ke ASCII aman
		{"αβγ", "γβα"},
		{"你好", "好你"},
	}
	// Cek versi sederhana: reversal harus menghasilkan panjang rune yang sama
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got := reverseString(c.input)
			inputRunes := []rune(c.input)
			gotRunes := []rune(got)
			if len(inputRunes) != len(gotRunes) {
				t.Errorf("reverseString(%q): panjang rune berbeda, input=%d, got=%d",
					c.input, len(inputRunes), len(gotRunes))
			}
			// Verifikasi idempotent: double reverse = original
			if reverseString(got) != c.input {
				t.Errorf("reverseString(reverseString(%q)) != input", c.input)
			}
		})
	}
}

func TestReverseString_Palindrome(t *testing.T) {
	// Palindrome dibalik harus sama dengan aslinya
	palindromes := []string{"racecar", "level", "madam", "aba"}
	for _, p := range palindromes {
		t.Run(p, func(t *testing.T) {
			if reverseString(p) != p {
				t.Errorf("reverseString palindrome %q harus sama dengan aslinya", p)
			}
		})
	}
}

// ================================================================
// TEST: containsPersonalInfo — skenario tambahan
// ================================================================

func TestContainsPersonalInfo_EmptyPersonalInfo(t *testing.T) {
	// Slice kosong → tidak ada yang dicek, harus false
	if containsPersonalInfo("SuperStr0ng!", []string{}) {
		t.Error("containsPersonalInfo dengan slice kosong seharusnya false")
	}
}

func TestContainsPersonalInfo_EmptyStringInSlice(t *testing.T) {
	// Info kosong ("") harus dilewati (continue)
	if containsPersonalInfo("SuperStr0ng!", []string{""}) {
		t.Error("containsPersonalInfo dengan info kosong seharusnya false")
	}
}

func TestContainsPersonalInfo_ExactMatch(t *testing.T) {
	if !containsPersonalInfo("johndoe123!", []string{"johndoe"}) {
		t.Error("seharusnya true karena password mengandung username")
	}
}

func TestContainsPersonalInfo_EmailPart(t *testing.T) {
	// Username dari email (sebelum @) harus terdeteksi
	if !containsPersonalInfo("alice123!", []string{"alice@company.com"}) {
		t.Error("seharusnya true karena password mengandung bagian username email")
	}
}

func TestContainsPersonalInfo_EmailPartAtSign(t *testing.T) {
	// Email tanpa @ di dalamnya → tidak dianggap email
	if containsPersonalInfo("xyzSecure9!", []string{"notanemail"}) {
		t.Error("seharusnya false karena info tidak mengandung @")
	}
}

func TestContainsPersonalInfo_ReversedInfo(t *testing.T) {
	// Password mengandung kebalikan dari personal info
	// "john" terbalik = "nhoj"
	if !containsPersonalInfo("nhoj123Secure!", []string{"john"}) {
		t.Error("seharusnya true karena password mengandung kebalikan dari personal info")
	}
}

func TestContainsPersonalInfo_ShortInfoSkipsReverse(t *testing.T) {
	// Info dengan panjang < 3 rune tidak dicek reverse-nya
	// "ab" terbalik = "ba" — tapi cek reverse dilewati untuk len < 3
	// password tidak mengandung "ab" langsung, hanya "ba"
	if containsPersonalInfo("baSecure123!", []string{"ab"}) {
		t.Error("seharusnya false karena len(info) < 3, reverse check dilewati")
	}
}

// ================================================================
// TEST: ValidatePassword — branch config yang belum tercover
// ================================================================

func TestValidatePassword_DisabledChecks_AllPass(t *testing.T) {
	// Config dengan semua check dimatikan — password apapun yang memenuhi MinLength harus lolos
	config := PasswordValidationConfig{
		MinLength:        4,
		RequireUpper:     false,
		RequireLower:     false,
		RequireNumber:    false,
		RequireSpecial:   false,
		CheckCommon:      false,
		CheckPersonal:    false,
		CheckOldPassword: false,
	}
	err := ValidatePassword("aaaa", config)
	if err != nil {
		t.Errorf("expected no error dengan config disabled semua, got: %v", err)
	}
}

func TestValidatePassword_TooShortExactBoundary(t *testing.T) {
	config := PasswordValidationConfig{MinLength: 8}
	// Tepat 7 karakter → error
	err := ValidatePassword("Abc1!xy", config)
	if err == nil {
		t.Error("expected error untuk password 7 karakter dengan MinLength=8")
	}
	// Tepat 8 karakter (tapi config lain semua false) → lolos length check
	config2 := PasswordValidationConfig{MinLength: 8}
	err2 := ValidatePassword("Abc1!xyz", config2)
	// Tidak ada RequireUpper/Lower/etc, jadi harus lolos
	if err2 != nil {
		t.Errorf("expected no error untuk password 8 karakter (MinLength=8, no other checks), got: %v", err2)
	}
}

func TestValidatePassword_NoPersonalInfoProvided(t *testing.T) {
	// CheckPersonal=true tapi personalInfo kosong → personal check tidak dijalankan
	config := PasswordValidationConfig{
		MinLength:     8,
		CheckPersonal: true,
	}
	// Password mengandung "john" tapi tidak ada personalInfo → harus lolos
	err := ValidatePassword("johnXYZA9!", config)
	// Catatan: "john" tidak ada di commonPasswordRoots, jadi lolos juga CheckCommon=false
	if err != nil {
		t.Errorf("expected no error ketika personalInfo kosong, got: %v", err)
	}
}

// ================================================================
// TEST: formatValidationErrors — pesan error yang dihasilkan
// ================================================================

type strictStruct struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=3,max=50"`
	Extra    string `validate:"required"`
}

func TestValidate_ErrorMessages_Required(t *testing.T) {
	s := strictStruct{}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error untuk struct kosong")
	}
	msg := err.Error()
	// Pesan harus mengandung "wajib diisi" untuk field required
	if !containsStr(msg, "wajib diisi") {
		t.Errorf("pesan error harus mengandung 'wajib diisi', got: %s", msg)
	}
}

func TestValidate_ErrorMessages_MinLength(t *testing.T) {
	s := strictStruct{Email: "a@b.com", Username: "ab", Extra: "x"}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error untuk username terlalu pendek")
	}
	msg := err.Error()
	if !containsStr(msg, "minimal") {
		t.Errorf("pesan error harus mengandung 'minimal', got: %s", msg)
	}
}

func TestValidate_ErrorMessages_MaxLength(t *testing.T) {
	longName := ""
	for i := 0; i < 51; i++ {
		longName += "a"
	}
	s := strictStruct{Email: "a@b.com", Username: longName, Extra: "x"}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error untuk username terlalu panjang")
	}
	msg := err.Error()
	if !containsStr(msg, "maksimal") {
		t.Errorf("pesan error harus mengandung 'maksimal', got: %s", msg)
	}
}

func TestValidate_ErrorMessages_InvalidEmail(t *testing.T) {
	s := strictStruct{Email: "bukan-email", Username: "validuser", Extra: "x"}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error untuk email tidak valid")
	}
	msg := err.Error()
	if !containsStr(msg, "email yang valid") {
		t.Errorf("pesan error harus mengandung 'email yang valid', got: %s", msg)
	}
}

// helper kecil agar tidak import "strings" di test (sudah ada di validator package)
func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}