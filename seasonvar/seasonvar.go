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
	linkRegexp              = regexp.MustCompile(`http\:\/\/seasonvar\.ru\/(.*)\.html`)
	seasonIDRegexp          = regexp.MustCompile(`data\-season\=\"([0-9]+)\"`)
	seasonTitleRegexp       = regexp.MustCompile(`\<title\>([^<]+)\<\/title\>`)
	seasonKeywordsRegexp    = regexp.MustCompile(`\<meta\ name\=\"keywords\"\ content\=\"([^"]+)\"`)
	seasonDescriptionRegexp = regexp.MustCompile(`\<meta\ name\=\"description\"\ content\=\"([^"]+)\"`)
)

type Seasonvar struct{}

type SeasonMeta struct {
	Title       string
	ID          string
	Link        string
	Keywords    string
	Description string
}

func New() *Seasonvar {
	return &Seasonvar{}
}

func (s *Seasonvar) ValidateLink(link string) error {
	if linkRegexp.FindString(link) == "" {
		return errors.New("incorrect link format")
	}
	return nil
}

func (s *Seasonvar) AbsoluteLink(link string) string {
	return fmt.Sprintf(seriesLinkFormat, link)
}

func (s *Seasonvar) GetSeasonMeta(link string) (*SeasonMeta, error) {
	var sm *SeasonMeta
	body, err := httpGet(link)
	if err != nil {
		return nil, err
	}
	seasonID, err := s.GetSeasonID(body)
	if err != nil {
		return nil, err
	}
	seasonTitle, err := s.GetSeasonTitle(body)
	if err != nil {
		return nil, err
	}
	seasonKeywords, err := s.GetSeasonKeywords(body)
	if err != nil {
		return nil, err
	}
	seasonDescription, err := s.GetSeasonDescription(body)
	if err != nil {
		return nil, err
	}
	sm = &SeasonMeta{
		ID:          seasonID,
		Link:        fmt.Sprintf(playlistLinkFormat, seasonID),
		Title:       seasonTitle,
		Keywords:    seasonKeywords,
		Description: seasonDescription,
	}
	return sm, nil
}

func (s *Seasonvar) GetSeasonTitle(body string) (string, error) {
	title := seasonTitleRegexp.FindStringSubmatch(body)
	if len(title) < 1 {
		return "", errors.New("season title not found")
	}
	return title[1], nil
}

func (s *Seasonvar) GetSeasonKeywords(body string) (string, error) {
	keywords := seasonKeywordsRegexp.FindStringSubmatch(body)
	if len(keywords) < 1 {
		return "", errors.New("season keywords not found")
	}
	return keywords[1], nil
}

func (s *Seasonvar) GetSeasonDescription(body string) (string, error) {
	description := seasonDescriptionRegexp.FindStringSubmatch(body)
	if len(description) < 1 {
		return "", errors.New("season description not found")
	}
	return description[1], nil
}

func (s *Seasonvar) GetSeasonID(body string) (string, error) {
	season := seasonIDRegexp.FindStringSubmatch(body)
	if len(season) < 1 {
		return "", errors.New("season id not found")
	}
	return season[1], nil
}
