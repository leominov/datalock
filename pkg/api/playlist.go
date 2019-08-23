package api

import (
	"errors"
)

var (
	ErrEmptyPlaylist = errors.New("Empty playlist")
	ErrHDNotFound    = errors.New("Files not found")
)

type Playlist struct {
	Name  string  `json:"name"`
	Items []*Item `json:"playlist"`
}

func (p *Playlist) DecodeLinks() error {
	for _, item := range p.Items {
		item.DecodeFileURL()
	}
	return nil
}

func (p *Playlist) UpdateSubtitleLinks() error {
	for _, item := range p.Items {
		if err := item.RemoveHostnameFromSubtitleLink(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Playlist) SwitchToHD(hdHostname string) error {
	var counter int
	if len(p.Items) == 0 {
		return ErrEmptyPlaylist
	}
	for _, item := range p.Items {
		if !item.AvailableInHD() {
			continue
		}
		if err := item.SwitchToHD(hdHostname); err != nil {
			continue
		}
		counter++
	}
	if counter == 0 {
		return ErrHDNotFound
	}
	return nil
}
