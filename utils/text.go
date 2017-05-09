package utils

import (
	"regexp"
	"strings"
)

var (
	cleanRegexp = regexp.MustCompile(`[^a-z0-9]`)
)

func CleanText(str string) string {
	str = strings.TrimSpace(str)
	str = cleanRegexp.ReplaceAllString(str, "${1}")
	return str
}