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

	if body["success"] != true {
		t.Error("expected success to be true")
	}

	respData, ok := body["data"].(map[string]interface{})
	if !ok || respData["message"] != "success" {
		t.Error("invalid response data")
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

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body["success"] != false {
		t.Error("expected success to be false")
	}

	if body["message"] != "bad request" {
		t.Error("expected error message")
	}
}
