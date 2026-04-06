package utils

import (
	"testing"
	"time"
)

func TestTimeNowUnix(t *testing.T) {
	result := TimeNowUnix()
	if result <= 0 {
		t.Error("expected positive unix timestamp")
	}
}

func TestTimeNowUnix_IsMonotonicallyIncreasing(t *testing.T) {
	first := TimeNowUnix()
	second := TimeNowUnix()

	if second < first {
		t.Errorf("second call (%d) seharusnya >= first call (%d)", second, first)
	}
}

func TestTimeNowUnix_IsCloseToSystemTime(t *testing.T) {
	// Nilai yang dikembalikan harus dalam rentang ±5 detik dari waktu sistem
	before := time.Now().Unix()
	result := TimeNowUnix()
	after := time.Now().Unix()

	if result < before-5 || result > after+5 {
		t.Errorf("TimeNowUnix() = %d, tidak dalam rentang [%d, %d]", result, before-5, after+5)
	}
}

func TestTimeNowUnix_IsAfterYear2020(t *testing.T) {
	// Unix timestamp 1 Jan 2020 00:00:00 UTC = 1577836800
	const year2020Unix = int64(1577836800)
	result := TimeNowUnix()

	if result < year2020Unix {
		t.Errorf("TimeNowUnix() = %d, seharusnya lebih besar dari epoch 2020 (%d)", result, year2020Unix)
	}
}
