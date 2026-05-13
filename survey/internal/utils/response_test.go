package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TEST RESPOND JSON
func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"message": "success",
	}

	RespondJSON(w, http.StatusOK, data)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Error("expected application/json header")
	}

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body["message"] != "success" {
		t.Error("expected message to be success")
	}
}

// TEST RESPOND SUCCESS
func TestRespondSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := "some-data"

	RespondSuccess(w, http.StatusOK, "OK", data)

	res := w.Result()
	defer res.Body.Close()

	var body JSONResponse
	json.NewDecoder(res.Body).Decode(&body)

	if body.Status != "success" {
		t.Error("expected status success")
	}
	if body.Message != "OK" {
		t.Error("expected message OK")
	}
	if body.Data != data {
		t.Error("data mismatch")
	}
}

// TEST RESPOND ERROR
func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()

	RespondError(w, http.StatusBadRequest, "bad request")

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", res.StatusCode)
	}

	var body JSONResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body.Status != "error" {
		t.Error("expected status error")
	}

	if body.Message != "bad request" {
		t.Error("expected error message")
	}
}
