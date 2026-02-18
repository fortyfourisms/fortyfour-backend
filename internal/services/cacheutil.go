package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"fortyfour-backend/pkg/cache"
)

// TTL constants — bisa disesuaikan sesuai kebutuhan
const (
	TTLList   = 5 * time.Minute  // untuk GetAll
	TTLDetail = 10 * time.Minute // untuk GetByID
)

// cacheGet mencoba ambil data dari Redis dan unmarshal ke target.
// Mengembalikan true jika cache hit, false jika miss atau error (fallback ke DB).
func cacheGet(rc cache.RedisInterface, key string, target interface{}) bool {
	if rc == nil {
		return false
	}
	raw, err := rc.Get(key)
	if err != nil {
		// Miss atau Redis down — fallback ke DB tanpa error
		return false
	}
	if err := json.Unmarshal([]byte(raw), target); err != nil {
		log.Printf("[cache] unmarshal error for key %s: %v", key, err)
		return false
	}
	return true
}

// cacheSet marshal data ke JSON dan simpan ke Redis.
// Jika Redis down, log warning tapi tidak return error (fallback tetap jalan).
func cacheSet(rc cache.RedisInterface, key string, data interface{}, ttl time.Duration) {
	if rc == nil {
		return
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("[cache] marshal error for key %s: %v", key, err)
		return
	}
	if err := rc.Set(key, string(b), ttl); err != nil {
		log.Printf("[cache] set error for key %s: %v", key, err)
	}
}

// cacheDelete hapus satu key dari Redis.
func cacheDelete(rc cache.RedisInterface, key string) {
	if rc == nil {
		return
	}
	if err := rc.Delete(key); err != nil {
		log.Printf("[cache] delete error for key %s: %v", key, err)
	}
}

// cacheDeletePattern hapus semua key yang cocok dengan pattern (menggunakan SCAN).
// Berguna untuk invalidate GetAll saat ada perubahan data.
func cacheDeletePattern(rc cache.RedisInterface, pattern string) {
	if rc == nil {
		return
	}
	keys, err := rc.Scan(pattern)
	if err != nil {
		log.Printf("[cache] scan error for pattern %s: %v", pattern, err)
		return
	}
	for _, k := range keys {
		cacheDelete(rc, k)
	}
}

// keyList membuat cache key untuk GetAll
func keyList(resource string) string {
	return fmt.Sprintf("%s:all", resource)
}

// keyDetail membuat cache key untuk GetByID
func keyDetail(resource, id string) string {
	return fmt.Sprintf("%s:%s", resource, id)
}