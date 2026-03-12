package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// Minimal Redis mock untuk cacheutil test
// ============================================================

type cacheTestRedis struct {
	data    map[string]string
	setErr  error
	getErr  error
	delErr  error
	scanErr error
}

func newCacheTestRedis() *cacheTestRedis {
	return &cacheTestRedis{data: make(map[string]string)}
}

func (r *cacheTestRedis) Set(key string, value interface{}, ttl time.Duration) error {
	if r.setErr != nil {
		return r.setErr
	}
	if v, ok := value.(string); ok {
		r.data[key] = v
	}
	return nil
}

func (r *cacheTestRedis) Get(key string) (string, error) {
	if r.getErr != nil {
		return "", r.getErr
	}
	v, ok := r.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}

func (r *cacheTestRedis) Delete(key string) error {
	if r.delErr != nil {
		return r.delErr
	}
	delete(r.data, key)
	return nil
}

func (r *cacheTestRedis) Exists(key string) (bool, error) {
	_, ok := r.data[key]
	return ok, nil
}

func (r *cacheTestRedis) Scan(pattern string) ([]string, error) {
	if r.scanErr != nil {
		return nil, r.scanErr
	}
	// Simple prefix matching: "resource:*" → match all keys with prefix "resource:"
	prefix := strings.TrimSuffix(pattern, "*")
	var keys []string
	for k := range r.data {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (r *cacheTestRedis) Close() error { return nil }

// ============================================================
// TEST cacheGet
// ============================================================

func TestCacheGet_NilRedis_ReturnsFalse(t *testing.T) {
	var target map[string]string
	hit := cacheGet(nil, "some-key", &target)
	assert.False(t, hit, "nil Redis harus return false")
}

func TestCacheGet_KeyNotFound_ReturnsFalse(t *testing.T) {
	rc := newCacheTestRedis()
	var target map[string]string
	hit := cacheGet(rc, "tidak-ada", &target)
	assert.False(t, hit, "key tidak ada harus return false")
}

func TestCacheGet_ValidJSON_ReturnsTrue(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["user:1"] = `{"id":"1","name":"Alice"}`

	var target map[string]string
	hit := cacheGet(rc, "user:1", &target)

	assert.True(t, hit)
	assert.Equal(t, "1", target["id"])
	assert.Equal(t, "Alice", target["name"])
}

func TestCacheGet_InvalidJSON_ReturnsFalse(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["bad-key"] = `ini bukan json{`

	var target map[string]string
	hit := cacheGet(rc, "bad-key", &target)

	assert.False(t, hit, "JSON rusak harus return false")
}

func TestCacheGet_RedisGetError_ReturnsFalse(t *testing.T) {
	rc := newCacheTestRedis()
	rc.getErr = errors.New("redis connection timeout")

	var target map[string]string
	hit := cacheGet(rc, "any-key", &target)

	assert.False(t, hit, "Redis error harus return false tanpa panic")
}

func TestCacheGet_Struct_UnmarshalCorrect(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	rc := newCacheTestRedis()
	rc.data["person:1"] = `{"name":"Budi","age":30}`

	var p Person
	hit := cacheGet(rc, "person:1", &p)

	assert.True(t, hit)
	assert.Equal(t, "Budi", p.Name)
	assert.Equal(t, 30, p.Age)
}

func TestCacheGet_Slice_UnmarshalCorrect(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["items:all"] = `[{"id":"a"},{"id":"b"}]`

	var result []map[string]string
	hit := cacheGet(rc, "items:all", &result)

	assert.True(t, hit)
	assert.Len(t, result, 2)
	assert.Equal(t, "a", result[0]["id"])
}

// ============================================================
// TEST cacheSet
// ============================================================

func TestCacheSet_NilRedis_NoPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		cacheSet(nil, "key", map[string]string{"x": "y"}, TTLList)
	})
}

func TestCacheSet_StoresJSON(t *testing.T) {
	rc := newCacheTestRedis()
	data := map[string]string{"hello": "world"}

	cacheSet(rc, "test-key", data, TTLDetail)

	stored, ok := rc.data["test-key"]
	assert.True(t, ok, "key harus tersimpan")
	assert.Contains(t, stored, "hello")
	assert.Contains(t, stored, "world")
}

func TestCacheSet_SetError_NoPanic(t *testing.T) {
	rc := newCacheTestRedis()
	rc.setErr = errors.New("redis full")

	assert.NotPanics(t, func() {
		cacheSet(rc, "key", map[string]string{"a": "b"}, TTLList)
	})
}

func TestCacheSet_Struct(t *testing.T) {
	type Item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	rc := newCacheTestRedis()
	cacheSet(rc, "item:1", Item{ID: "1", Name: "Test"}, TTLDetail)

	stored, ok := rc.data["item:1"]
	assert.True(t, ok)
	assert.Contains(t, stored, `"id":"1"`)
	assert.Contains(t, stored, `"name":"Test"`)
}

func TestCacheSet_Slice(t *testing.T) {
	rc := newCacheTestRedis()
	items := []string{"a", "b", "c"}

	cacheSet(rc, "list:all", items, TTLList)

	stored, ok := rc.data["list:all"]
	assert.True(t, ok)
	assert.Contains(t, stored, "a")
}

// ============================================================
// TEST cacheDelete
// ============================================================

func TestCacheDelete_NilRedis_NoPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		cacheDelete(nil, "some-key")
	})
}

func TestCacheDelete_RemovesKey(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["to-delete"] = "value"

	cacheDelete(rc, "to-delete")

	_, exists := rc.data["to-delete"]
	assert.False(t, exists, "key harus sudah dihapus")
}

func TestCacheDelete_KeyNotExist_NoPanic(t *testing.T) {
	rc := newCacheTestRedis()

	assert.NotPanics(t, func() {
		cacheDelete(rc, "tidak-ada")
	})
}

func TestCacheDelete_DeleteError_NoPanic(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["key"] = "value"
	rc.delErr = errors.New("redis delete error")

	assert.NotPanics(t, func() {
		cacheDelete(rc, "key")
	})
}

func TestCacheDelete_OnlyDeletesTargetKey(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["key-1"] = "a"
	rc.data["key-2"] = "b"
	rc.data["key-3"] = "c"

	cacheDelete(rc, "key-2")

	assert.Equal(t, 2, len(rc.data), "hanya key-2 yang terhapus")
	_, k1 := rc.data["key-1"]
	_, k3 := rc.data["key-3"]
	assert.True(t, k1)
	assert.True(t, k3)
}

// ============================================================
// TEST cacheDeletePattern
// ============================================================

func TestCacheDeletePattern_NilRedis_NoPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		cacheDeletePattern(nil, "resource:*")
	})
}

func TestCacheDeletePattern_DeletesMatchingKeys(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["perusahaan:1"] = "a"
	rc.data["perusahaan:2"] = "b"
	rc.data["sektor:1"] = "c"

	cacheDeletePattern(rc, "perusahaan:*")

	assert.Equal(t, 1, len(rc.data), "hanya key perusahaan yang terhapus")
	_, sektorMasih := rc.data["sektor:1"]
	assert.True(t, sektorMasih, "key sektor:1 harus tetap ada")
}

func TestCacheDeletePattern_NoMatch_NoDelete(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["role:1"] = "a"
	rc.data["role:2"] = "b"

	cacheDeletePattern(rc, "user:*")

	assert.Equal(t, 2, len(rc.data), "tidak ada yang terhapus jika tidak ada yang cocok")
}

func TestCacheDeletePattern_ScanError_NoPanic(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["key:1"] = "a"
	rc.scanErr = errors.New("scan failed")

	assert.NotPanics(t, func() {
		cacheDeletePattern(rc, "key:*")
	})

	// Key tidak terhapus karena scan gagal
	assert.Equal(t, 1, len(rc.data))
}

func TestCacheDeletePattern_DeletesAll_WhenPatternMatchAll(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["user:1"] = "a"
	rc.data["user:2"] = "b"
	rc.data["user:abc"] = "c"

	cacheDeletePattern(rc, "user:*")

	assert.Equal(t, 0, len(rc.data))
}

// ============================================================
// TEST keyList
// ============================================================

func TestKeyList_Format(t *testing.T) {
	assert.Equal(t, "perusahaan:all", keyList("perusahaan"))
	assert.Equal(t, "role:all", keyList("role"))
	assert.Equal(t, "user:all", keyList("user"))
}

func TestKeyList_EmptyResource(t *testing.T) {
	result := keyList("")
	assert.Equal(t, ":all", result)
}

// ============================================================
// TEST keyDetail
// ============================================================

func TestKeyDetail_Format(t *testing.T) {
	assert.Equal(t, "perusahaan:uuid-123", keyDetail("perusahaan", "uuid-123"))
	assert.Equal(t, "role:abc-def", keyDetail("role", "abc-def"))
}

func TestKeyDetail_DifferentIDsSameName_ProduceDifferentKeys(t *testing.T) {
	key1 := keyDetail("user", "id-1")
	key2 := keyDetail("user", "id-2")
	assert.NotEqual(t, key1, key2)
}

func TestKeyDetail_DifferentResourcesSameID_ProduceDifferentKeys(t *testing.T) {
	key1 := keyDetail("role", "123")
	key2 := keyDetail("user", "123")
	assert.NotEqual(t, key1, key2)
}

// ============================================================
// TEST TTL constants
// ============================================================

func TestTTLConstants(t *testing.T) {
	assert.Equal(t, 5*time.Minute, TTLList, "TTLList harus 5 menit")
	assert.Equal(t, 10*time.Minute, TTLDetail, "TTLDetail harus 10 menit")
	assert.Greater(t, TTLDetail, TTLList, "TTLDetail harus lebih lama dari TTLList")
}

// ============================================================
// TEST integrasi: cacheSet → cacheGet
// ============================================================

func TestCacheSetThenGet_Roundtrip(t *testing.T) {
	type User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	rc := newCacheTestRedis()
	user := User{ID: "u-1", Name: "Alice"}

	cacheSet(rc, "user:u-1", user, TTLDetail)

	var got User
	hit := cacheGet(rc, "user:u-1", &got)

	assert.True(t, hit)
	assert.Equal(t, "u-1", got.ID)
	assert.Equal(t, "Alice", got.Name)
}

func TestCacheDeleteThenGet_ReturnsFalse(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["key:1"] = `{"id":"1"}`

	cacheDelete(rc, "key:1")

	var target map[string]string
	hit := cacheGet(rc, "key:1", &target)
	assert.False(t, hit, "setelah delete, get harus miss")
}

func TestCacheDeletePatternThenGet_ReturnsFalse(t *testing.T) {
	rc := newCacheTestRedis()
	rc.data["product:1"] = `{"id":"1"}`
	rc.data["product:2"] = `{"id":"2"}`

	cacheDeletePattern(rc, "product:*")

	var t1, t2 map[string]string
	assert.False(t, cacheGet(rc, "product:1", &t1))
	assert.False(t, cacheGet(rc, "product:2", &t2))
}
