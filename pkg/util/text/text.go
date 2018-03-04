package text

import (
	"encoding/base64"
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

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func Base64Decode(s string) string {
	r, _ := base64.StdEncoding.DecodeString(s)
	return string(r)
}
