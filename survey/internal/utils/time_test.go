package utils

import (
	"testing"
	"time"
)

// TEST FORMAT ISO
func TestFormatISO(t *testing.T) {
	tm := time.Date(2026, 4, 29, 10, 30, 0, 0, time.UTC)

	got := FormatISO(tm)
	want := "2026-04-29T10:30:00Z"

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

// TEST NOW FUNCTION
func TestNow(t *testing.T) {
	before := time.Now()
	got := Now()
	after := time.Now()

	// cek apakah waktu berada di range yang masuk akal
	if got.Before(before) || got.After(after) {
		t.Error("Now() returned invalid time")
	}
}