package api

import (
	"net/url"
	"strings"
)

type Item struct {
	ID         string `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Comment    string `json:"comment,omitempty"`
	Subtitle   string `json:"subtitle,omitempty"`
	Sub        string `json:"sub,omitempty"`
	Streamsend string `json:"streamsend,omitempty"`
	File       string `json:"file"`
	GALabel    string `json:"galabel"`
}

func (i *Item) AvailableInHD() bool {
	index := strings.Index(i.Title, "HD")
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

func (i *Item) RemoveHostnameFromSubtitleLink() error {
	if len(i.Sub) != 0 {
		u, err := url.Parse(i.Sub)
		if err != nil {
			return err
		}
		i.Sub = u.Path + "?" + u.RawQuery
	}
	if len(i.Subtitle) != 0 {
		u, err := url.Parse(i.Subtitle)
		if err != nil {
			return err
		}
		i.Subtitle = u.Path + "?" + u.RawQuery
	}
	return nil
}
