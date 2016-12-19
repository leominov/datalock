package seasonvar

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	playlistLinkFormat = "http://datalock.ru/player/%s"
	seriesLinkFormat   = "http://seasonvar.ru%s"
)

var (
	linkRegexp   = regexp.MustCompile(`http\:\/\/seasonvar\.ru\/(.*)\.html`)
	seasonRegexp = regexp.MustCompile(`data\-season\=\"([0-9]+)\"`)
)

type Seasonvar struct{}

func New() *Seasonvar {
	return &Seasonvar{}
}

func (s *Seasonvar) ValidateLink(link string) error {
	if linkRegexp.FindString(link) == "" {
		return errors.New("incorrect url format")
	}
	return nil
}

func (s *Seasonvar) AbsoluteLink(link string) string {
	return fmt.Sprintf(seriesLinkFormat, link)
}

func (s *Seasonvar) GetSeasonLink(url string) (string, error) {
	body, err := httpGet(url)
	if err != nil {
		return "", err
	}
	season := seasonRegexp.FindStringSubmatch(body)
	if len(season) < 1 {
		return "", errors.New("season id not found")
	}
	return fmt.Sprintf(playlistLinkFormat, season[1]), nil
}
