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
)
