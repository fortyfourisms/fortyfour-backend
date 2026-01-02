package utils

import (
	"testing"
)

func TestTimeNowUnix(t *testing.T) {
	result := TimeNowUnix()
	if result <= 0 {
		t.Error("expected positive unix timestamp")
	}
}

