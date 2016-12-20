package seasonvar

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	playlistLinkFormat = "http://datalock.ru/player/%d"
	seriesLinkFormat   = "http://seasonvar.ru%s"
)

var (
	linkRegexp              = regexp.MustCompile(`http\:\/\/seasonvar\.ru\/(.*)\.html`)
	seasonIDRegexp          = regexp.MustCompile(`data\-season\=\"([0-9]+)\"`)
	seasonTitleRegexp       = regexp.MustCompile(`\<title\>([^<]+)\<\/title\>`)
	seasonKeywordsRegexp    = regexp.MustCompile(`\<meta\ name\=\"keywords\"\ content\=\"([^"]+)\"`)
	seasonDescriptionRegexp = regexp.MustCompile(`\<meta\ name\=\"description\"\ content\=\"([^"]+)\"`)
)

type Seasonvar struct {
	Seasons map[string]*SeasonMeta
}

type SeasonMeta struct {
	Title       string
	ID          int
	Link        string
	Keywords    string
	Description string
}

func New() *Seasonvar {
	return &Seasonvar{
		Seasons: make(map[string]*SeasonMeta),
	}
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

func (s *Seasonvar) GetSeasonMeta(link string) (*SeasonMeta, bool, error) {
	if meta, ok := s.Seasons[link]; ok {
		return meta, true, nil
	}
	seasonMeta, err := s.collectSeasonMeta(link)
	if err != nil {
		return nil, false, err
	}
	s.Seasons[link] = seasonMeta
	return seasonMeta, false, nil
}

func (s *Seasonvar) collectSeasonMeta(link string) (*SeasonMeta, error) {
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

func (s *Seasonvar) GetSeasonID(body string) (int, error) {
	season := seasonIDRegexp.FindStringSubmatch(body)
	if len(season) < 1 {
		return 0, errors.New("season id not found")
	}
	i, err := strconv.Atoi(season[1])
	if err != nil {
		return 0, err
	}
	return i, nil
}
