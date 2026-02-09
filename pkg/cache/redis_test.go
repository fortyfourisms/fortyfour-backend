package cache_test

import (
	"fortyfour-backend/pkg/cache"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMiniRedis creates a fake Redis server for testing
func setupMiniRedis(t *testing.T) (*miniredis.Miniredis, cache.RedisConfig) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	cfg := cache.RedisConfig{
		Host:     "localhost",
		Port:     mr.Port(),
		Password: "",
		DB:       0,
	}

	return mr, cfg
}

func TestNewRedisClient(t *testing.T) {
	t.Run("Successful connection", func(t *testing.T) {
		mr, cfg := setupMiniRedis(t)
		defer mr.Close()

		client, err := cache.NewRedisClient(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, client)

		// Clean up
		if client != nil {
			client.Close()
		}
	})

	t.Run("Failed connection - invalid host", func(t *testing.T) {
		cfg := cache.RedisConfig{
			Host:     "invalid-host-12345",
			Port:     "6379",
			Password: "",
			DB:       0,
		}

		client, err := cache.NewRedisClient(cfg)

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "failed to connect to Redis")
	})

	t.Run("Failed connection - invalid port", func(t *testing.T) {
		cfg := cache.RedisConfig{
			Host:     "localhost",
			Port:     "99999",
			Password: "",
			DB:       0,
		}

		client, err := cache.NewRedisClient(cfg)

		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("Connection with password", func(t *testing.T) {
		mr, cfg := setupMiniRedis(t)
		defer mr.Close()

		// Set password on mini redis
		mr.RequireAuth("testpass")
		cfg.Password = "testpass"

		client, err := cache.NewRedisClient(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, client)

		if client != nil {
			client.Close()
		}
	})

	t.Run("Wrong password", func(t *testing.T) {
		mr, cfg := setupMiniRedis(t)
		defer mr.Close()

		mr.RequireAuth("correctpass")
		cfg.Password = "wrongpass"

		client, err := cache.NewRedisClient(cfg)

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestRedisClient_Set(t *testing.T) {
	mr, cfg := setupMiniRedis(t)
	defer mr.Close()

	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	tests := []struct {
		name       string
		key        string
		value      interface{}
		expiration time.Duration
		wantErr    bool
	}{
		{
			name:       "Set string value",
			key:        "test:string",
			value:      "hello world",
			expiration: 5 * time.Minute,
			wantErr:    false,
		},
		{
			name:       "Set integer value",
			key:        "test:int",
			value:      42,
			expiration: 1 * time.Hour,
			wantErr:    false,
		},
		{
			name:       "Set with no expiration",
			key:        "test:noexp",
			value:      "permanent",
			expiration: 0,
			wantErr:    false,
		},
		{
			name:       "Set empty string",
			key:        "test:empty",
			value:      "",
			expiration: 1 * time.Minute,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Set(tt.key, tt.value, tt.expiration)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the value was set
				exists, err := client.Exists(tt.key)
				assert.NoError(t, err)
				assert.True(t, exists)
			}
		})
	}
}

func TestRedisClient_Get(t *testing.T) {
	mr, cfg := setupMiniRedis(t)
	defer mr.Close()

	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	t.Run("Get existing key", func(t *testing.T) {
		key := "test:get:existing"
		expectedValue := "test value"

		// Set value first
		err := client.Set(key, expectedValue, 5*time.Minute)
		require.NoError(t, err)

		// Get value
		value, err := client.Get(key)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("Get non-existing key", func(t *testing.T) {
		key := "test:get:nonexistent"

		value, err := client.Get(key)

		assert.Error(t, err)
		assert.Equal(t, "", value)
		assert.Equal(t, redis.Nil, err)
	})

	t.Run("Get expired key", func(t *testing.T) {
		key := "test:get:expired"

		// Set with very short expiration
		err := client.Set(key, "will expire", 1*time.Millisecond)
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(10 * time.Millisecond)

		// Fast forward time in miniredis
		mr.FastForward(2 * time.Millisecond)

		value, err := client.Get(key)

		assert.Error(t, err)
		assert.Equal(t, "", value)
	})
}

func TestRedisClient_Delete(t *testing.T) {
	mr, cfg := setupMiniRedis(t)
	defer mr.Close()

	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	t.Run("Delete existing key", func(t *testing.T) {
		key := "test:delete:existing"

		// Set value first
		err := client.Set(key, "to be deleted", 5*time.Minute)
		require.NoError(t, err)

		// Delete
		err = client.Delete(key)
		assert.NoError(t, err)

		// Verify deletion
		exists, err := client.Exists(key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Delete non-existing key", func(t *testing.T) {
		key := "test:delete:nonexistent"

		// Delete non-existing key (should not error)
		err := client.Delete(key)
		assert.NoError(t, err)
	})

	t.Run("Delete multiple times", func(t *testing.T) {
		key := "test:delete:multiple"

		// Set value
		err := client.Set(key, "test", 5*time.Minute)
		require.NoError(t, err)

		// Delete twice
		err = client.Delete(key)
		assert.NoError(t, err)

		err = client.Delete(key)
		assert.NoError(t, err)
	})
}

func TestRedisClient_Exists(t *testing.T) {
	mr, cfg := setupMiniRedis(t)
	defer mr.Close()

	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	t.Run("Check existing key", func(t *testing.T) {
		key := "test:exists:yes"

		// Set value
		err := client.Set(key, "exists", 5*time.Minute)
		require.NoError(t, err)

		// Check existence
		exists, err := client.Exists(key)

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Check non-existing key", func(t *testing.T) {
		key := "test:exists:no"

		exists, err := client.Exists(key)

		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Check after deletion", func(t *testing.T) {
		key := "test:exists:deleted"

		// Set and delete
		err := client.Set(key, "will be deleted", 5*time.Minute)
		require.NoError(t, err)

		err = client.Delete(key)
		require.NoError(t, err)

		// Check existence
		exists, err := client.Exists(key)

		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestRedisClient_Close(t *testing.T) {
	mr, cfg := setupMiniRedis(t)
	defer mr.Close()

	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err)

	t.Run("Close connection", func(t *testing.T) {
		err := client.Close()
		assert.NoError(t, err)
	})

	t.Run("Operations after close should fail", func(t *testing.T) {
		// Try to set after closing
		err := client.Set("test:after:close", "value", 1*time.Minute)
		assert.Error(t, err)
	})
}

func TestRedisClient_ComplexScenarios(t *testing.T) {
	mr, cfg := setupMiniRedis(t)
	defer mr.Close()

	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	t.Run("Set, Get, Update, Delete workflow", func(t *testing.T) {
		key := "test:workflow"

		// Set initial value
		err := client.Set(key, "initial", 5*time.Minute)
		assert.NoError(t, err)

		// Get value
		value, err := client.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, "initial", value)

		// Update value
		err = client.Set(key, "updated", 5*time.Minute)
		assert.NoError(t, err)

		// Get updated value
		value, err = client.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, "updated", value)

		// Delete
		err = client.Delete(key)
		assert.NoError(t, err)

		// Verify deletion
		exists, err := client.Exists(key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Multiple keys operations", func(t *testing.T) {
		keys := []string{"key1", "key2", "key3"}

		// Set multiple keys
		for i, key := range keys {
			err := client.Set(key, i, 5*time.Minute)
			assert.NoError(t, err)
		}

		// Check all exist
		for _, key := range keys {
			exists, err := client.Exists(key)
			assert.NoError(t, err)
			assert.True(t, exists)
		}

		// Delete all
		for _, key := range keys {
			err := client.Delete(key)
			assert.NoError(t, err)
		}
	})
}

func TestRedisClient_Expiration(t *testing.T) {
	mr, cfg := setupMiniRedis(t)
	defer mr.Close()

	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	t.Run("Key expires correctly", func(t *testing.T) {
		key := "test:expiration"

		// Set with 100ms expiration
		err := client.Set(key, "expires soon", 100*time.Millisecond)
		assert.NoError(t, err)

		// Should exist immediately
		exists, err := client.Exists(key)
		assert.NoError(t, err)
		assert.True(t, exists)

		// Fast forward time
		mr.FastForward(150 * time.Millisecond)

		// Should not exist after expiration
		exists, err = client.Exists(key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestRedisClient_Scan(t *testing.T) {
	mr, cfg := setupMiniRedis(t)
	defer mr.Close()

	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	t.Run("Scan matching keys", func(t *testing.T) {
		// Setup keys
		keysToCreate := []string{"test:scan:1", "test:scan:2", "other:key"}
		for _, k := range keysToCreate {
			err := client.Set(k, "value", 5*time.Minute)
			require.NoError(t, err)
		}

		// Scan for test:scan:*
		keys, err := client.Scan("test:scan:*")
		assert.NoError(t, err)

		// Note: Miniredis might return scan results in any order or implementation specific
		// Scan documentation says it returns keys.

		// We expect 2 keys
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "test:scan:1")
		assert.Contains(t, keys, "test:scan:2")
	})

	t.Run("Scan no matching keys", func(t *testing.T) {
		keys, err := client.Scan("nonexistent:*")
		assert.NoError(t, err)
		assert.Empty(t, keys)
	})
}

// Test interface compliance
func TestRedisClient_ImplementsInterface(t *testing.T) {
	mr, cfg := setupMiniRedis(t)
	defer mr.Close()

	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// Verify it implements RedisInterface
	var _ cache.RedisInterface = client
}

// Benchmark tests
func BenchmarkRedisClient_Set(b *testing.B) {
	mr, _ := miniredis.Run()
	defer mr.Close()

	cfg := cache.RedisConfig{
		Host: "localhost",
		Port: mr.Port(),
		DB:   0,
	}

	client, _ := cache.NewRedisClient(cfg)
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Set("benchmark:key", "value", 5*time.Minute)
	}
}

func BenchmarkRedisClient_Get(b *testing.B) {
	mr, _ := miniredis.Run()
	defer mr.Close()

	cfg := cache.RedisConfig{
		Host: "localhost",
		Port: mr.Port(),
		DB:   0,
	}

	client, _ := cache.NewRedisClient(cfg)
	defer client.Close()

	// Setup
	client.Set("benchmark:key", "value", 5*time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Get("benchmark:key")
	}
}
