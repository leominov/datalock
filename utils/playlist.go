package utils

import (
	"net/url"
	"regexp"
	"strings"
)

const (
	playlistNameStandard = "Стандартный"
)

var (
	playlistRegexp     = regexp.MustCompile(`\/playls2\/([^/]+)/([^/]+)/([0-9]+)/list.xml`)
	playlistNameRegexp = regexp.MustCompile(`/trans([^/]+)/`)
)

func GetPlaylistNameByLink(link string) string {
	if !playlistNameRegexp.MatchString(link) {
		return playlistNameStandard
	}
	matches := playlistNameRegexp.FindStringSubmatch(link)
	if len(matches) > 0 {
		name := matches[1]
		if strings.Contains(name, "%") {
			name, _ = url.PathUnescape(name)
		}
		return name
	}
	return ""
}

func GetPlaylistLinksFromText(body []byte) map[string]string {
	result := make(map[string]string)
	matches := playlistRegexp.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return result
	}
	for _, match := range matches {
		link := match[0]
		name := GetPlaylistNameByLink(link)
		result[name] = link
	}
	return result
}
