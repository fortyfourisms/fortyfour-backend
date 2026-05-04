package cache

import "time"

// DefaultCacheExpiration specifies the default TTL for cached items.
const DefaultCacheExpiration = 24 * time.Hour

// Cache Keys for IKAS Service
const (
	CacheKeyIkasRecords            = "ikas:records"
	CacheKeyPertanyaanIdentifikasi = "ikas:questions:identifikasi"
	CacheKeyPertanyaanProteksi     = "ikas:questions:proteksi"
	CacheKeyPertanyaanDeteksi      = "ikas:questions:deteksi"
	CacheKeyPertanyaanGulih        = "ikas:questions:gulih"

	// Cache Keys for Master Data
	CacheKeyRuangLingkup = "ikas:master:ruanglingkup"
	CacheKeyDomain       = "ikas:master:domain"
	CacheKeyKategori     = "ikas:master:kategori"
	CacheKeySubKategori  = "ikas:master:subkategori"

	// Cache Key Prefixes for Jawaban (per ikas_id)
	// Usage: fmt.Sprintf(CacheKeyPrefixJawabanIdentifikasi+"%s", ikasID)
	CacheKeyPrefixJawabanIdentifikasi = "ikas:jawaban:identifikasi:"
	CacheKeyPrefixJawabanProteksi     = "ikas:jawaban:proteksi:"
	CacheKeyPrefixJawabanDeteksi      = "ikas:jawaban:deteksi:"
	CacheKeyPrefixJawabanGulih        = "ikas:jawaban:gulih:"

	// Cache Key Prefixes for Domain (per ikas_id)
	// Usage: fmt.Sprintf(CacheKeyPrefixIdentifikasi+"%s", ikasID)
	CacheKeyPrefixIdentifikasi = "ikas:domain:identifikasi:"
	CacheKeyPrefixProteksi     = "ikas:domain:proteksi:"
	CacheKeyPrefixDeteksi      = "ikas:domain:deteksi:"
	CacheKeyPrefixGulih        = "ikas:domain:gulih:"

	// Cache Key Prefixes for Audit Logs (paginated, keyed by page, limit, and optional ikas_id)
	// Usage: fmt.Sprintf(CacheKeyPrefixAuditLogs+"%d:%d", page, limit)
	//        fmt.Sprintf(CacheKeyPrefixAuditLogsByIkas+"%s:%d:%d", ikasID, page, limit)
	CacheKeyPrefixAuditLogs       = "ikas:audit_logs:all:"
	CacheKeyPrefixAuditLogsByIkas = "ikas:audit_logs:ikas:"

	// AuditLogsCacheExpiration is shorter than default to keep audit data fresher.
	AuditLogsCacheExpiration = 5 * time.Minute
)

// GetJawabanIdentifikasiKey returns the cache key for Jawaban Identifikasi
func GetJawabanIdentifikasiKey(ikasID string) string {
	return CacheKeyPrefixJawabanIdentifikasi + ikasID
}

// GetJawabanProteksiKey returns the cache key for Jawaban Proteksi
func GetJawabanProteksiKey(ikasID string) string {
	return CacheKeyPrefixJawabanProteksi + ikasID
}

// GetJawabanDeteksiKey returns the cache key for Jawaban Deteksi
func GetJawabanDeteksiKey(ikasID string) string {
	return CacheKeyPrefixJawabanDeteksi + ikasID
}

// GetJawabanGulihKey returns the cache key for Jawaban Gulih
func GetJawabanGulihKey(ikasID string) string {
	return CacheKeyPrefixJawabanGulih + ikasID
}

// GetIdentifikasiKey returns the cache key for Domain Identifikasi
func GetIdentifikasiKey(ikasID string) string {
	return CacheKeyPrefixIdentifikasi + ikasID
}

// GetProteksiKey returns the cache key for Domain Proteksi
func GetProteksiKey(ikasID string) string {
	return CacheKeyPrefixProteksi + ikasID
}

// GetDeteksiKey returns the cache key for Domain Deteksi
func GetDeteksiKey(ikasID string) string {
	return CacheKeyPrefixDeteksi + ikasID
}

// GetGulihKey returns the cache key for Domain Gulih
func GetGulihKey(ikasID string) string {
	return CacheKeyPrefixGulih + ikasID
}
