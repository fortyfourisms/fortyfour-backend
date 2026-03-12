package utils

import "strings"

func ExtractID(path, resource string) string {
	parts := strings.SplitN(path, "/"+resource, 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimPrefix(parts[1], "/")
}
