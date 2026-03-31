package services

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

// =========================
// MOCK REDIS
// Nama berbeda dari mockRedis (auth_service_test.go) dan
// mockTokenRedis (token_service_test.go) untuk menghindari konflik.
// =========================

type mockCacheRedis struct {
	store      map[string]string
	failSet    bool
	failGet    bool
	failDelete bool
	failScan   bool
	setCalled  int
	delCalled  int
}

func newMockCacheRedis() *mockCacheRedis {
	return &mockCacheRedis{store: make(map[string]string)}
}

func (m *mockCacheRedis) Set(key string, value interface{}, ttl time.Duration) error {
	m.setCalled++
	if m.failSet {
		return errors.New("redis set error")
	}
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be string")
	}
	m.store[key] = str
	return nil
}

func (m *mockCacheRedis) Get(key string) (string, error) {
	if m.failGet {
		return "", errors.New("redis get error")
	}
	val, ok := m.store[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}

func (m *mockCacheRedis) Delete(key string) error {
	m.delCalled++
	if m.failDelete {
		return errors.New("redis delete error")
	}
	delete(m.store, key)
	return nil
}

func (m *mockCacheRedis) Exists(key string) (bool, error) {
	_, ok := m.store[key]
	return ok, nil
}

func (m *mockCacheRedis) Close() error { return nil }

func (m *mockCacheRedis) Scan(pattern string) ([]string, error) {
	if m.failScan {
		return nil, errors.New("redis scan error")
	}
	// Tiru perilaku Redis SCAN: filter key berdasarkan pattern.
	// Mendukung wildcard suffix (*) seperti "user:*" dan exact match.
	prefix := strings.TrimSuffix(pattern, "*")
	keys := make([]string, 0, len(m.store))
	for k := range m.store {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

// =========================
// TEST: cacheGet
// =========================

func TestCacheGet_HitReturnsTrue(t *testing.T) {
	rc := newMockCacheRedis()
	rc.store["user:1"] = `{"id":"1","name":"alice"}`

	var result map[string]string
	got := cacheGet(rc, "user:1", &result)

	if !got {
		t.Error("expected true on cache hit")
	}
	if result["name"] != "alice" {
		t.Errorf("expected name='alice', got '%s'", result["name"])
	}
}

func TestCacheGet_MissReturnsFalse(t *testing.T) {
	rc := newMockCacheRedis()

	var result map[string]string
	got := cacheGet(rc, "user:nonexistent", &result)

	if got {
		t.Error("expected false on cache miss")
	}
}

func TestCacheGet_NilRedisReturnsFalse(t *testing.T) {
	var result map[string]string
	got := cacheGet(nil, "user:1", &result)

	if got {
		t.Error("expected false when redis is nil")
	}
}

func TestCacheGet_RedisGetErrorReturnsFalse(t *testing.T) {
	rc := newMockCacheRedis()
	rc.failGet = true

	var result map[string]string
	got := cacheGet(rc, "user:1", &result)

	if got {
		t.Error("expected false when redis Get returns error")
	}
}

func TestCacheGet_CorruptJSONReturnsFalse(t *testing.T) {
	rc := newMockCacheRedis()
	rc.store["user:1"] = `{not valid json`

	var result map[string]string
	got := cacheGet(rc, "user:1", &result)

	if got {
		t.Error("expected false when stored value is corrupt JSON")
	}
}

func TestCacheGet_PreservesTargetOnMiss(t *testing.T) {
	rc := newMockCacheRedis()

	result := map[string]string{"key": "original"}
	cacheGet(rc, "nonexistent", &result)

	if result["key"] != "original" {
		t.Error("target value should not be modified on cache miss")
	}
}

func TestCacheGet_SliceOfStructs(t *testing.T) {
	rc := newMockCacheRedis()

	type item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	rc.store["items:all"] = `[{"id":1,"name":"foo"},{"id":2,"name":"bar"}]`

	var items []item
	if !cacheGet(rc, "items:all", &items) {
		t.Fatal("expected true on cache hit for slice")
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if items[0].Name != "foo" {
		t.Errorf("expected first item name 'foo', got '%s'", items[0].Name)
	}
}

func TestCacheGet_NestedStruct(t *testing.T) {
	rc := newMockCacheRedis()

	type address struct {
		City string `json:"city"`
	}
	type user struct {
		ID      string  `json:"id"`
		Address address `json:"address"`
	}

	rc.store["user:nested"] = `{"id":"u1","address":{"city":"Jakarta"}}`

	var u user
	if !cacheGet(rc, "user:nested", &u) {
		t.Fatal("expected true for nested struct hit")
	}
	if u.Address.City != "Jakarta" {
		t.Errorf("expected city='Jakarta', got '%s'", u.Address.City)
	}
}

// =========================
// TEST: cacheSet
// =========================

func TestCacheSet_StoresValidJSON(t *testing.T) {
	rc := newMockCacheRedis()

	data := map[string]string{"id": "1", "name": "alice"}
	cacheSet(rc, "user:1", data, TTLDetail)

	raw, ok := rc.store["user:1"]
	if !ok {
		t.Fatal("expected key 'user:1' to be stored in redis")
	}

	var result map[string]string
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("stored value is not valid JSON: %v", err)
	}
	if result["name"] != "alice" {
		t.Errorf("expected name='alice', got '%s'", result["name"])
	}
}

func TestCacheSet_NilRedisIsNoOp(t *testing.T) {
	// Tidak boleh panic
	cacheSet(nil, "user:1", map[string]string{"x": "y"}, TTLDetail)
}

func TestCacheSet_RedisErrorIsGraceful(t *testing.T) {
	rc := newMockCacheRedis()
	rc.failSet = true

	// Tidak boleh panic (silent fail by design)
	cacheSet(rc, "user:1", map[string]string{"x": "y"}, TTLDetail)
}

func TestCacheSet_UnmarshalableDataDoesNotStore(t *testing.T) {
	rc := newMockCacheRedis()

	// Channel tidak bisa di-marshal ke JSON
	cacheSet(rc, "bad:key", make(chan int), TTLDetail)

	if _, ok := rc.store["bad:key"]; ok {
		t.Error("unmarshalable data should not be stored in redis")
	}
}

func TestCacheSet_OverwritesExistingKey(t *testing.T) {
	rc := newMockCacheRedis()
	rc.store["user:1"] = `{"name":"old"}`

	cacheSet(rc, "user:1", map[string]string{"name": "new"}, TTLDetail)

	var result map[string]string
	if err := json.Unmarshal([]byte(rc.store["user:1"]), &result); err != nil {
		t.Fatalf("invalid JSON after overwrite: %v", err)
	}
	if result["name"] != "new" {
		t.Errorf("expected name='new' after overwrite, got '%s'", result["name"])
	}
}

func TestCacheSet_CallsRedisSetOnce(t *testing.T) {
	rc := newMockCacheRedis()

	cacheSet(rc, "user:1", "value", TTLDetail)

	if rc.setCalled != 1 {
		t.Errorf("expected Set called once, got %d", rc.setCalled)
	}
}

func TestCacheSet_SliceData(t *testing.T) {
	rc := newMockCacheRedis()

	data := []int{1, 2, 3}
	cacheSet(rc, "numbers:all", data, TTLList)

	raw, ok := rc.store["numbers:all"]
	if !ok {
		t.Fatal("expected key to be stored")
	}

	var result []int
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("stored value not valid JSON: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d", len(result))
	}
}

// =========================
// TEST: cacheDelete
// =========================

func TestCacheDelete_RemovesKey(t *testing.T) {
	rc := newMockCacheRedis()
	rc.store["user:1"] = `{"id":"1"}`

	cacheDelete(rc, "user:1")

	if _, ok := rc.store["user:1"]; ok {
		t.Error("expected key to be deleted")
	}
}

func TestCacheDelete_NilRedisIsNoOp(t *testing.T) {
	// Tidak boleh panic
	cacheDelete(nil, "user:1")
}

func TestCacheDelete_RedisErrorIsGraceful(t *testing.T) {
	rc := newMockCacheRedis()
	rc.failDelete = true

	// Tidak boleh panic (silent fail by design)
	cacheDelete(rc, "user:1")
}

func TestCacheDelete_NonExistentKeyIsNoOp(t *testing.T) {
	rc := newMockCacheRedis()

	// Hapus key yang tidak ada — tidak boleh error/panic
	cacheDelete(rc, "nonexistent:key")
}

func TestCacheDelete_OnlyRemovesTargetKey(t *testing.T) {
	rc := newMockCacheRedis()
	rc.store["user:1"] = `{}`
	rc.store["user:2"] = `{}`

	cacheDelete(rc, "user:1")

	if _, ok := rc.store["user:1"]; ok {
		t.Error("user:1 should be deleted")
	}
	if _, ok := rc.store["user:2"]; !ok {
		t.Error("user:2 should NOT be affected")
	}
}

// =========================
// TEST: cacheDeletePattern
// =========================

func TestCacheDeletePattern_DeletesAllMatchingKeys(t *testing.T) {
	rc := newMockCacheRedis()
	rc.store["user:1"] = `{}`
	rc.store["user:2"] = `{}`
	rc.store["user:3"] = `{}`
	rc.store["role:1"] = `{}`

	cacheDeletePattern(rc, "user:*")

	for _, k := range []string{"user:1", "user:2", "user:3"} {
		if _, ok := rc.store[k]; ok {
			t.Errorf("key '%s' should have been deleted", k)
		}
	}
}

func TestCacheDeletePattern_DoesNotDeleteNonMatchingKeys(t *testing.T) {
	rc := newMockCacheRedis()
	rc.store["user:1"] = `{}`
	rc.store["role:1"] = `{}`

	cacheDeletePattern(rc, "user:*")

	if _, ok := rc.store["role:1"]; !ok {
		t.Error("key 'role:1' should NOT be deleted by 'user:*' pattern")
	}
}

func TestCacheDeletePattern_NilRedisIsNoOp(t *testing.T) {
	// Tidak boleh panic
	cacheDeletePattern(nil, "user:*")
}

func TestCacheDeletePattern_ScanErrorIsGraceful(t *testing.T) {
	rc := newMockCacheRedis()
	rc.store["user:1"] = `{}`
	rc.failScan = true

	// Tidak boleh panic (silent fail)
	cacheDeletePattern(rc, "user:*")

	// Key masih ada karena Scan gagal dan delete tidak dieksekusi
	if _, ok := rc.store["user:1"]; !ok {
		t.Error("key should remain untouched when Scan fails")
	}
}

func TestCacheDeletePattern_EmptyStoreIsNoOp(t *testing.T) {
	rc := newMockCacheRedis()

	// Store kosong — tidak boleh error/panic
	cacheDeletePattern(rc, "user:*")
}

func TestCacheDeletePattern_CallsDeleteForEachKey(t *testing.T) {
	rc := newMockCacheRedis()
	rc.store["user:1"] = `{}`
	rc.store["user:2"] = `{}`
	rc.store["user:3"] = `{}`

	cacheDeletePattern(rc, "user:*")

	if rc.delCalled != 3 {
		t.Errorf("expected Delete called 3 times, got %d", rc.delCalled)
	}
}

// =========================
// TEST: keyList
// =========================

func TestKeyList_ReturnsCorrectFormat(t *testing.T) {
	tests := []struct {
		resource string
		expected string
	}{
		{"user", "user:all"},
		{"role", "role:all"},
		{"sektor", "sektor:all"},
		{"perusahaan", "perusahaan:all"},
	}

	for _, tc := range tests {
		got := keyList(tc.resource)
		if got != tc.expected {
			t.Errorf("keyList(%q): want '%s', got '%s'", tc.resource, tc.expected, got)
		}
	}
}

func TestKeyList_EmptyResourceProducesColonAll(t *testing.T) {
	got := keyList("")
	if got != ":all" {
		t.Errorf("keyList(''): want ':all', got '%s'", got)
	}
}

// =========================
// TEST: keyDetail
// =========================

func TestKeyDetail_ReturnsCorrectFormat(t *testing.T) {
	tests := []struct {
		resource string
		id       string
		expected string
	}{
		{"user", "123", "user:123"},
		{"role", "abc-def", "role:abc-def"},
		{"sektor", "uuid-1234", "sektor:uuid-1234"},
		{"perusahaan", "comp-99", "perusahaan:comp-99"},
	}

	for _, tc := range tests {
		got := keyDetail(tc.resource, tc.id)
		if got != tc.expected {
			t.Errorf("keyDetail(%q, %q): want '%s', got '%s'", tc.resource, tc.id, tc.expected, got)
		}
	}
}

func TestKeyDetail_EmptyIDProducesResourceColon(t *testing.T) {
	got := keyDetail("user", "")
	if got != "user:" {
		t.Errorf("keyDetail('user', ''): want 'user:', got '%s'", got)
	}
}

// =========================
// TEST: TTL constants
// =========================

func TestTTLConstants_AreCorrect(t *testing.T) {
	if TTLList != 5*time.Minute {
		t.Errorf("TTLList: want 5m, got %v", TTLList)
	}
	if TTLDetail != 10*time.Minute {
		t.Errorf("TTLDetail: want 10m, got %v", TTLDetail)
	}
}

func TestTTLList_IsLessThanTTLDetail(t *testing.T) {
	if TTLList >= TTLDetail {
		t.Errorf("expected TTLList (%v) < TTLDetail (%v)", TTLList, TTLDetail)
	}
}

// =========================
// TEST: Integrasi
// =========================

func TestCacheUtil_SetThenGetHit(t *testing.T) {
	rc := newMockCacheRedis()

	type user struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	original := user{ID: "u-1", Name: "alice"}
	cacheSet(rc, keyDetail("user", "u-1"), original, TTLDetail)

	var retrieved user
	if !cacheGet(rc, keyDetail("user", "u-1"), &retrieved) {
		t.Fatal("expected cache hit after cacheSet")
	}
	if retrieved.ID != "u-1" || retrieved.Name != "alice" {
		t.Errorf("retrieved value mismatch: %+v", retrieved)
	}
}

func TestCacheUtil_SetThenDeleteThenGetMiss(t *testing.T) {
	rc := newMockCacheRedis()

	cacheSet(rc, "user:u-1", map[string]string{"id": "u-1"}, TTLDetail)
	cacheDelete(rc, "user:u-1")

	var result map[string]string
	if cacheGet(rc, "user:u-1", &result) {
		t.Error("expected cache miss after delete")
	}
}

func TestCacheUtil_PatternInvalidatesMultipleKeys(t *testing.T) {
	rc := newMockCacheRedis()

	// Simulasi: simpan list + 3 detail item
	cacheSet(rc, keyList("role"), []string{"admin", "viewer"}, TTLList)
	cacheSet(rc, keyDetail("role", "1"), "admin", TTLDetail)
	cacheSet(rc, keyDetail("role", "2"), "viewer", TTLDetail)

	// Invalidate semua
	cacheDeletePattern(rc, "role:*")

	// Semua harus miss
	var listResult []string
	if cacheGet(rc, keyList("role"), &listResult) {
		t.Error("expected list cache miss after pattern delete")
	}

	var detailResult string
	if cacheGet(rc, keyDetail("role", "1"), &detailResult) {
		t.Error("expected detail cache miss after pattern delete")
	}
}

func TestCacheUtil_NilRedisNeverPanicsAcrossAllFunctions(t *testing.T) {
	// Smoke test: semua fungsi dengan redis=nil tidak boleh panic
	var result map[string]string

	cacheGet(nil, "any:key", &result)
	cacheSet(nil, "any:key", "value", TTLDetail)
	cacheDelete(nil, "any:key")
	cacheDeletePattern(nil, "any:*")
}
