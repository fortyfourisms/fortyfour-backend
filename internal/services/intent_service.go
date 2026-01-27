package services

import "strings"

func DetectIntent(message string) string {
	msg := strings.ToLower(message)

	switch {
	case strings.Contains(msg, "perusahaan terbaru"):
		return "latest_perusahaan"
	default:
		return "unknown"
	}
}
