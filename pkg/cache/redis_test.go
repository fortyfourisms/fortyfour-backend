package cache_test

import (
	"fortyfour-backend/pkg/cache"
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     cache.RedisConfig
		want    *cache.RedisClient
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := cache.NewRedisClient(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewRedisClient() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewRedisClient() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("NewRedisClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_Set(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cfg cache.RedisConfig
		// Named input parameters for target function.
		key        string
		value      interface{}
		expiration time.Duration
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := cache.NewRedisClient(tt.cfg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := r.Set(tt.key, tt.value, tt.expiration)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Set() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Set() succeeded unexpectedly")
			}
		})
	}
}

func TestRedisClient_Get(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cfg cache.RedisConfig
		// Named input parameters for target function.
		key     string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := cache.NewRedisClient(tt.cfg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := r.Get(tt.key)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Get() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Get() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_Delete(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cfg cache.RedisConfig
		// Named input parameters for target function.
		key     string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := cache.NewRedisClient(tt.cfg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := r.Delete(tt.key)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Delete() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Delete() succeeded unexpectedly")
			}
		})
	}
}

func TestRedisClient_Exists(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cfg cache.RedisConfig
		// Named input parameters for target function.
		key     string
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := cache.NewRedisClient(tt.cfg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := r.Exists(tt.key)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Exists() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Exists() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_Close(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cfg     cache.RedisConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := cache.NewRedisClient(tt.cfg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := r.Close()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Close() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Close() succeeded unexpectedly")
			}
		})
	}
}
