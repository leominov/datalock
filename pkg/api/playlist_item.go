package api

import (
	"encoding/base64"
	"net/url"
	"strings"
)

type Item struct {
	ID         string  `json:"id,omitempty"`
	Title      string  `json:"title,omitempty"`
	Comment    string  `json:"comment,omitempty"`
	Subtitle   string  `json:"subtitle,omitempty"`
	Sub        string  `json:"sub,omitempty"`
	Streamsend string  `json:"streamsend,omitempty"`
	File       string  `json:"file"`
	GALabel    string  `json:"galabel"`
	Folder     []*Item `json:"folder,omitempty"`
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

func (i *Item) DecodeFileURL() {
	if !strings.HasPrefix(i.File, "#2") {
		return
	}
	link := i.File
	link = strings.Replace(link, "#2", "", -1)
	link = strings.Replace(link, "//b2xvbG8=", "", -1)
	decoded, err := base64.StdEncoding.DecodeString(link)
	if err != nil {
		return
	}
	link = string(decoded)
	if !strings.HasPrefix(link, "http") {
		return
	}
	i.File = link
}
