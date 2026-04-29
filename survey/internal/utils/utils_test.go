package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TEST VALUE OR NULL
func TestValueOrNull(t *testing.T) {
	val := "hello"
	res := ValueOrNull(&val)

	if res != "hello" {
		t.Errorf("expected hello, got %v", res)
	}

	resNil := ValueOrNull(nil)
	if resNil != nil {
		t.Errorf("expected nil, got %v", resNil)
	}
}

// TEST VALUE OR EMPTY
func TestValueOrEmpty(t *testing.T) {
	val := "world"
	res := ValueOrEmpty(&val)

	if res != "world" {
		t.Errorf("expected world, got %s", res)
	}

	resNil := ValueOrEmpty(nil)
	if resNil != "" {
		t.Errorf("expected empty string, got %s", resNil)
	}
}

// TEST INT OR NULL
func TestIntOrNull(t *testing.T) {
	v := 10
	res := IntOrNull(&v)

	if res != 10 {
		t.Errorf("expected 10, got %v", res)
	}

	resNil := IntOrNull(nil)
	if resNil != nil {
		t.Errorf("expected nil, got %v", resNil)
	}
}

// TEST BOOL OR NULL
func TestBoolOrNull(t *testing.T) {
	v := true
	res := BoolOrNull(&v)

	if res != true {
		t.Errorf("expected true, got %v", res)
	}

	resNil := BoolOrNull(nil)
	if resNil != nil {
		t.Errorf("expected nil, got %v", resNil)
	}
}

// TEST ADAPT HANDLER
func TestAdaptHandler(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := AdaptHandler(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}
}