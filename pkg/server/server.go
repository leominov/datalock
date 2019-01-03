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
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/leominov/datalock/pkg/api"
	"github.com/leominov/datalock/pkg/backends"
	"github.com/leominov/datalock/pkg/metrics"
	"github.com/leominov/datalock/pkg/util/blacklist"
	"github.com/leominov/datalock/pkg/util/httpget"
	"github.com/leominov/datalock/pkg/util/useragent"
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
)

type Server struct {
	NodeName        string
	NodeList        []*Node
	Config          *Config
	Blacklist       *blacklist.Blacklist
	storeClient     backends.StoreClient
	reloadTemplates bool
}

type Node struct {
	NodeName string
	Hostname string
	mu       sync.Mutex
	Healthy  bool
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

func New(config *Config) (*Server, error) {
	hostname, _ := os.Hostname()
	storeClient, err := backends.New(config.StorageClientConfig)
	if err != nil {
		return nil, err
	}
	s := &Server{
		NodeName:    hostname,
		Config:      config,
		storeClient: storeClient,
	}
	err = s.LoadNodeList()
	if err != nil {
		return nil, err
	}
	if len(s.NodeList) > 0 {
		for _, node := range s.NodeList {
			log.Printf("Node %s is %s", node.NodeName, node.State())
		}
	}
	return s, nil
}

func BoolAsHit(hitCache bool) string {
	if hitCache {
		return "HIT"
	}
	return "MISS"
}

func (n *Node) State() string {
	if n.Healthy {
		return "healthy"
	}
	return "inactive"
}

func (s *Server) Run() {
	if len(s.NodeList) == 0 {
		return
	}
	log.Print("Starting node list heartbeat...")
	time.Sleep(5 * time.Second)
	for {
		for _, node := range s.NodeList {
			healthy := s.IsHealthyNode(node)
			if node.Healthy != healthy {
				node.mu.Lock()
				node.Healthy = healthy
				node.mu.Unlock()
				log.Printf("Node %s switch state to %s", node.NodeName, node.State())
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *Server) LoadNodeList() error {
	for _, nodeAddr := range s.Config.NodeList {
		node := s.GetNodeByAddr(nodeAddr)
		s.NodeList = append(s.NodeList, node)
	}
	return nil
}

func (s *Server) GetNodeByAddr(addr string) *Node {
	n := &Node{
		NodeName: addr,
		Hostname: addr,
	}
	n.mu.Lock()
	state := s.IsHealthyNode(n)
	n.mu.Unlock()
	n.Healthy = state
	return n
}

func (s *Server) IsHealthyNode(node *Node) bool {
	addr := fmt.Sprintf("http://%s%s", node.Hostname, s.Config.HealthzPath)
	cli := http.DefaultClient
	cli.Timeout = 5 * time.Second
	resp, err := cli.Get(addr)
	if err != nil {
		log.Printf("Error requesting node %s state: %s", node.NodeName, err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Node %s has incorrect status: %s", node.NodeName, resp.Status)
		return false
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading node %s state: %s", node.NodeName, err)
		return false
	}
	state := strings.TrimSpace(string(b))
	if state == "ok" {
		return true
	}
	log.Printf("Node %s has incorrect state: %s", node.NodeName, state)
	return false
}

func (s *Server) LoadBlacklist(path string) error {
	bl, err := blacklist.NewBlacklist(path)
	if err != nil {
		return err
	}
	s.Blacklist = bl
	return nil
}

func (s *Server) Stop() error {
	return s.storeClient.Close()
}

func (s *Server) AbsoluteLink(link string) string {
	return fmt.Sprintf(SeriesLinkFormat, s.Config.Hostname, link)
}

func (n *Node) AbsoluteLink(link string) string {
	return fmt.Sprintf(SeriesLinkFormat, n.Hostname, link)
}

func (s *Server) GetCachedSeasonMeta(link string) (*SeasonMeta, bool, error) {
	var seasonMeta *SeasonMeta
	var err error
	var hitCache bool
	seasonID, err := s.GetSeasonIDFromLink(link)
	if err != nil {
		return nil, false, err
	}
	seasonMeta, err = s.GetSeasonMeta(seasonID)
	if err != nil {
		hitCache = false
		seasonMeta, err = s.collectSeasonMeta(link)
		if err != nil {
			return nil, false, err
		}
	} else {
		hitCache = true
	}
	if err := s.SetSeasonMeta(seasonMeta); err != nil {
		return nil, hitCache, err
	}
	return seasonMeta, hitCache, nil
}

func (s *Server) SwitchSeriesLink(url string, isUserRequest bool) string {
	if isUserRequest {
		if s.Blacklist.IsAlias(url) {
			source := s.Blacklist.GetSource(url)
			return source
		}
	} else {
		if s.Blacklist.IsBlocked(url) {
			alias := s.Blacklist.GetAlias(url)
			return alias
		}
	}
	return url
}

func (s *Server) collectSeasonMeta(link string) (*SeasonMeta, error) {
	var seasonMeta *SeasonMeta
	metrics.HttpRequestsTotalCount.Inc()
	body, err := httpget.HttpGet(link)
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
	return s.storeClient.SetValue(strconv.Itoa(m.ID), &m)
}

func (s *Server) GetSeasonMeta(id int) (*SeasonMeta, error) {
	var m *SeasonMeta
	err := s.storeClient.GetValue(strconv.Itoa(id), &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func (s *Server) CanShowHD(r *http.Request) bool {
	if coo, err := r.Cookie("hdq"); err == nil && coo.Value != "" {
		return true
	}
	return false
}

func (s *Server) GetPlaylist(link string, hd, arrayResponse bool) (*api.Playlist, error) {
	b, err := httpget.HttpGetRaw(link, map[string]string{})
	if err != nil {
		return nil, err
	}
	playlist := new(api.Playlist)
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

func (s *Server) GetPlaylistsByLinks(links map[string]string, hd bool) ([]*api.Playlist, error) {
	playlists := []*api.Playlist{}
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
	req.Header.Set("User-Agent", useragent.Random())
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (s *Server) UpdateResponseBody(r *http.Response) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = r.Body.Close()
	if err != nil {
		return err
	}
	b = bytes.Replace(b, []byte(s.Config.Hostname), []byte(s.Config.PublicHostname), -1)
	for _, item := range s.Blacklist.List {
		b = bytes.Replace(b, []byte(item.Source[1:]), []byte(item.Alias[1:]), -1)
	}
	body := ioutil.NopCloser(bytes.NewReader(b))
	r.Body = body
	r.ContentLength = int64(len(b))
	r.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func (s *Server) MarkFlushTemplatesCache() {
	s.reloadTemplates = true
}

func (s *Server) FlushTemplatesCache() {
	if !s.reloadTemplates {
		return
	}
	log.Println("Clearing templates cache")
	s.reloadTemplates = false
	ParseTemplates(s.Config)
}
