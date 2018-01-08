package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/leominov/datalock/metrics"
	"github.com/leominov/datalock/utils"
)

const (
	SeriesLinkFormat = "http://%s%s"
)

var (
	seasonIDLinkRegexp      = regexp.MustCompile(`serial\-([0-9]+)\-`)
	seasonIDRegexp          = regexp.MustCompile(`data\-id\-season\=\"([0-9]+)\"`)
	serialIDRegexp          = regexp.MustCompile(`data\-id\-serial\=\"([0-9]+)\"`)
	seasonTitleRegexp       = regexp.MustCompile(`\<title\>([^<]+)\<\/title\>`)
	seasonKeywordsRegexp    = regexp.MustCompile(`\<meta\ name\=\"keywords\"\ content\=\"([^"]+)\"`)
	seasonDescriptionRegexp = regexp.MustCompile(`\<meta\ name\=\"description\"\ content\=\"([^"]+)\"`)

	MetaBucket = []byte("meta")
)

type Server struct {
	NodeName string
	Config   *Config
	DB       *bolt.DB
}

type SeasonMeta struct {
	Title       string `json:"title"`
	ID          int    `json:"id"`
	Serial      int    `json:"serial"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
}

type User struct {
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	SecureMark string `json:"secure_mark"`
}

func New(config *Config) *Server {
	hostname, _ := os.Hostname()
	return &Server{
		NodeName: hostname,
		Config:   config,
	}
}

func (s *Server) Start() error {
	var err error
	s.DB, err = bolt.Open(path.Join(s.Config.DatabaseDir, "datalock.db"), 0600, nil)
	if err != nil {
		return err
	}
	return s.DB.Update(func(tx *bolt.Tx) error {
		// Always create Meta bucket.
		if _, err := tx.CreateBucketIfNotExists(MetaBucket); err != nil {
			return err
		}
		return nil
	})
}

func (s *Server) Stop() error {
	return s.DB.Close()
}

func (s *Server) AbsoluteLink(link string) string {
	return fmt.Sprintf(SeriesLinkFormat, s.Config.Hostname, link)
}

func (s *Server) GetCachedSeasonMeta(link string) (*SeasonMeta, error) {
	var seasonMeta *SeasonMeta
	var err error
	seasonID, err := s.GetSeasonIDFromLink(link)
	if err != nil {
		return nil, err
	}
	seasonMeta, err = s.GetSeasonMeta(seasonID)
	if err != nil {
		seasonMeta, err = s.collectSeasonMeta(link)
		if err != nil {
			return nil, err
		}
	}
	if err := s.SetSeasonMeta(seasonMeta); err != nil {
		return nil, err
	}
	return seasonMeta, nil
}

func (s *Server) collectSeasonMeta(link string) (*SeasonMeta, error) {
	var seasonMeta *SeasonMeta
	metrics.HttpRequestsTotalCount.Inc()
	body, err := utils.HttpGet(link)
	if err != nil {
		metrics.HttpRequestsErrorCount.Inc()
		return nil, err
	}
	seasonID, err := s.GetSeasonID(body)
	if err != nil {
		return nil, err
	}
	serialID, err := s.GetSerialID(body)
	if err != nil {
		return nil, err
	}
	seasonTitle, err := s.GetSeasonTitle(body)
	if err != nil {
		seasonTitle = ""
	}
	seasonKeywords, err := s.GetSeasonKeywords(body)
	if err != nil {
		seasonKeywords = ""
	}
	seasonDescription, err := s.GetSeasonDescription(body)
	if err != nil {
		seasonDescription = ""
	}
	seasonMeta = &SeasonMeta{
		ID:          seasonID,
		Serial:      serialID,
		Title:       seasonTitle,
		Keywords:    seasonKeywords,
		Description: seasonDescription,
	}
	return seasonMeta, nil
}

func (s *Server) GetSeasonTitle(body string) (string, error) {
	title := seasonTitleRegexp.FindStringSubmatch(body)
	if len(title) < 1 {
		metrics.SeasonTitleErrorCount.Inc()
		return "", errors.New("season title not found")
	}
	return title[1], nil
}

func (s *Server) GetSeasonKeywords(body string) (string, error) {
	keywords := seasonKeywordsRegexp.FindStringSubmatch(body)
	if len(keywords) < 1 {
		metrics.SeasonKeywordsErrorCount.Inc()
		return "", errors.New("season keywords not found")
	}
	return keywords[1], nil
}

func (s *Server) GetSeasonDescription(body string) (string, error) {
	description := seasonDescriptionRegexp.FindStringSubmatch(body)
	if len(description) < 1 {
		metrics.SeasonDescriptionErrorCount.Inc()
		return "", errors.New("season description not found")
	}
	return description[1], nil
}

func (s *Server) GetSeasonID(body string) (int, error) {
	season := seasonIDRegexp.FindStringSubmatch(body)
	if len(season) < 1 {
		metrics.SeasonIDErrorCount.Inc()
		return 0, errors.New("season id not found")
	}
	i, err := strconv.Atoi(season[1])
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (s *Server) GetSerialID(body string) (int, error) {
	serial := serialIDRegexp.FindStringSubmatch(body)
	if len(serial) < 1 {
		metrics.SerialIDErrorCount.Inc()
		return 0, errors.New("serial id not found")
	}
	i, err := strconv.Atoi(serial[1])
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (s *Server) GetSeasonIDFromLink(link string) (int, error) {
	season := seasonIDLinkRegexp.FindStringSubmatch(link)
	if len(season) < 1 {
		return 0, errors.New("season id not found")
	}
	i, err := strconv.Atoi(season[1])
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (s *Server) GetUser(ip string) *User {
	return &User{
		IP:         "127.0.0.1",
		SecureMark: "0",
	}
}

func (s *Server) SetSeasonMeta(m *SeasonMeta) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(MetaBucket)
		encoded, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return b.Put([]byte(strconv.Itoa(m.ID)), encoded)
	})
}

func (s *Server) GetSeasonMeta(id int) (*SeasonMeta, error) {
	var m *SeasonMeta
	return m, s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(MetaBucket)
		v := b.Get([]byte(strconv.Itoa(id)))
		if len(v) == 0 {
			return errors.New("Meta not found")
		}
		if err := json.Unmarshal(v, &m); err != nil {
			return err
		}
		return nil
	})
}

func (s *Server) CanShowHD(r *http.Request) bool {
	if coo, err := r.Cookie("hdq"); err == nil && coo.Value != "" {
		return true
	}
	return false
}

func (s *Server) GetPlaylist(link string, hd, arrayResponse bool) (*Playlist, error) {
	b, err := utils.HttpGetRaw(link, map[string]string{})
	if err != nil {
		return nil, err
	}
	playlist := new(Playlist)
	if arrayResponse {
		if err := json.Unmarshal(b, &playlist.Items); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal(b, &playlist); err != nil {
			return nil, err
		}
	}
	if hd {
		// Nothing change if switching was failed
		playlist.SwitchToHD(s.Config.HdHostname)
	}
	if err := playlist.UpdateSubtitleLinks(); err != nil {
		log.Println(err)
	}
	return playlist, nil
}

func (s *Server) GetPlaylistsByLinks(links map[string]string, hd bool) ([]*Playlist, error) {
	playlists := []*Playlist{}
	for name, link := range links {
		linkAbs := s.AbsoluteLink(link)
		playlist, err := s.GetPlaylist(linkAbs, hd, true)
		if err != nil {
			return nil, err
		}
		if len(playlist.Items) == 0 {
			continue
		}
		playlist.Name = name
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (s *Server) FixReferer(req *http.Request) {
	refererRaw := req.Header.Get("Referer")
	if len(refererRaw) == 0 {
		return
	}
	refererUrl, err := url.Parse(refererRaw)
	if err != nil {
		req.Header.Del("Referer")
		return
	}
	refererUrl.Host = s.Config.Hostname
	req.Header.Set("Referer", refererUrl.String())
}

func (s *Server) NewPlaylistRequest(form url.Values) (*http.Request, error) {
	form.Add("type", "html5")
	req, err := http.NewRequest("POST", s.AbsoluteLink("/player.php"), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", utils.RandomUserAgent())
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (s *Server) UpdateHostnameResponseBody(r *http.Response) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = r.Body.Close()
	if err != nil {
		return err
	}
	b = bytes.Replace(b, []byte(s.Config.Hostname), []byte(s.Config.PublicHostname), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	r.Body = body
	r.ContentLength = int64(len(b))
	r.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}
