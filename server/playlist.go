package server

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrEmptyPlaylist = errors.New("Empty playlist")
	ErrHDNotFound    = errors.New("Files not found")
)

type Playlist struct {
	Name  string  `json:"name"`
	Items []*Item `json:"playlist"`
}

type Item struct {
	Comment    string `json:"comment"`
	File       string `json:"file"`
	StreamsEnd string `json:"streamsend"`
	GALabel    string `json:"galabel"`
}

func (i *Item) AvailableInHD() bool {
	index := strings.Index(i.Comment, "HD")
	if index >= 0 {
		return true
	}
	return false
}

func (i *Item) SwitchToHD(hdHostname string) error {
	u, err := url.Parse(i.File)
	if err != nil {
		return err
	}
	i.File = strings.Replace(i.File, u.Host, hdHostname, 1)
	i.File = strings.Replace(i.File, "/7f_", "/hd_", 1)
	return nil
}

func (p *Playlist) SwitchToHD(hdHostname string) error {
	var counter int
	if len(p.Items) == 0 {
		return ErrEmptyPlaylist
	}
	for id, item := range p.Items {
		if !item.AvailableInHD() {
			continue
		}
		if err := item.SwitchToHD(hdHostname); err != nil {
			continue
		}
		counter += 1
		p.Items[id] = item
	}
	if counter == 0 {
		return ErrHDNotFound
	}
	return nil
}
