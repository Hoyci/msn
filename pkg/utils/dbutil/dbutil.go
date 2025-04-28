package dbutil

import (
	"regexp"
	"strings"
)

func ExtractFieldFromDetail(detail string) string {
	detail = strings.TrimSpace(detail)

	re := regexp.MustCompile(`(?i)Key\s*\(\s*(.*?)\s*\)\s*=`)
	matches := re.FindStringSubmatch(detail)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
