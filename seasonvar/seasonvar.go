package seasonvar

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/leominov/datalock/metrics"
	"github.com/leominov/datalock/utils"
)

const (
	Hostname         = "seasonvar.ru"
	SeriesLinkFormat = "http://%s%s"
)

var (
	seasonIDLinkRegexp      = regexp.MustCompile(`serial\-([0-9]+)\-`)
	seasonIDRegexp          = regexp.MustCompile(`data\-id\-season\=\"([0-9]+)\"`)
	serialIDRegexp          = regexp.MustCompile(`data\-id\-serial\=\"([0-9]+)\"`)
	seasonTitleRegexp       = regexp.MustCompile(`\<title\>([^<]+)\<\/title\>`)
	seasonKeywordsRegexp    = regexp.MustCompile(`\<meta\ name\=\"keywords\"\ content\=\"([^"]+)\"`)
	seasonDescriptionRegexp = regexp.MustCompile(`\<meta\ name\=\"description\"\ content\=\"([^"]+)\"`)

	BucketUsers = []byte("users")
	BucketMeta  = []byte("meta")
)

type Seasonvar struct {
	NodeName string
	Config   *Config
	DB       *bolt.DB
}

type SeasonMeta struct {
	Title       string
	ID          int
	Serial      int
	Keywords    string
	Description string
}

type User struct {
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	SecureMark string `json:"secure_mark"`
}

func New(config *Config) *Seasonvar {
	hostname, _ := os.Hostname()
	return &Seasonvar{
		NodeName: hostname,
		Config:   config,
	}
}

func (s *Seasonvar) Start() error {
	var err error
	s.DB, err = bolt.Open(path.Join(s.Config.DatabaseDir, "datalock.db"), 0600, nil)
	if err != nil {
		return err
	}
	return s.DB.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(BucketUsers); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(BucketMeta); err != nil {
			return err
		}
		return nil
	})
}

func (s *Seasonvar) Stop() error {
	return s.DB.Close()
}

func (s *Seasonvar) AbsoluteLink(link string) string {
	return fmt.Sprintf(SeriesLinkFormat, Hostname, link)
}

func (s *Seasonvar) GetCachedSeasonMeta(link string) (*SeasonMeta, error) {
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

func (s *Seasonvar) collectSeasonMeta(link string) (*SeasonMeta, error) {
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

func (s *Seasonvar) GetSeasonTitle(body string) (string, error) {
	title := seasonTitleRegexp.FindStringSubmatch(body)
	if len(title) < 1 {
		metrics.SeasonTitleErrorCount.Inc()
		return "", errors.New("season title not found")
	}
	return title[1], nil
}

func (s *Seasonvar) GetSeasonKeywords(body string) (string, error) {
	keywords := seasonKeywordsRegexp.FindStringSubmatch(body)
	if len(keywords) < 1 {
		metrics.SeasonKeywordsErrorCount.Inc()
		return "", errors.New("season keywords not found")
	}
	return keywords[1], nil
}

func (s *Seasonvar) GetSeasonDescription(body string) (string, error) {
	description := seasonDescriptionRegexp.FindStringSubmatch(body)
	if len(description) < 1 {
		metrics.SeasonDescriptionErrorCount.Inc()
		return "", errors.New("season description not found")
	}
	return description[1], nil
}

func (s *Seasonvar) GetSeasonID(body string) (int, error) {
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

func (s *Seasonvar) GetSerialID(body string) (int, error) {
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

func (s *Seasonvar) GetSeasonIDFromLink(link string) (int, error) {
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

func (s *Seasonvar) SetUser(u *User) error {
	u.SecureMark = utils.CleanText(u.SecureMark)
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketUsers)
		encoded, err := json.Marshal(u)
		if err != nil {
			return err
		}
		return b.Put([]byte(u.IP), encoded)
	})
}

func (s *Seasonvar) GetUser(ip string) (*User, error) {
	var u *User
	return u, s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketUsers)
		v := b.Get([]byte(ip))
		if len(v) == 0 {
			return errors.New("User not found")
		}
		if err := json.Unmarshal(v, &u); err != nil {
			return err
		}
		return nil
	})
}

func (s *Seasonvar) SetSeasonMeta(m *SeasonMeta) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketMeta)
		encoded, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return b.Put([]byte(strconv.Itoa(m.ID)), encoded)
	})
}

func (s *Seasonvar) GetSeasonMeta(id int) (*SeasonMeta, error) {
	var m *SeasonMeta
	return m, s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketMeta)
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

func (s *Seasonvar) CanShowHD(r *http.Request) bool {
	if coo, err := r.Cookie("hdq"); err == nil && coo.Value != "" {
		return true
	}
	return false
}
