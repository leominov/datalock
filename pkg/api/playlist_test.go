package api

import (
	"strings"
	"testing"
)

func TestSwitchToHD(t *testing.T) {
	p := &Playlist{
		Items: []*Item{
			&Item{
				Title: "Episode 1 HD",
				File:  "http://sd/7f_episode_1.avi",
			},
			&Item{
				Title: "Episode 2 SD",
				File:  "http://sd/7f_episode_2.avi",
			},
			&Item{
				Title: "Episode 3 SD",
				File:  "http://[::1]:namedport",
			},
		},
	}
	p.SwitchToHD("hd")
	if p.Items[0].File != "http://hd/hd_episode_1.avi" {
		t.Errorf("Must be hd link")
	}
	if p.Items[1].File == "http://hd/hd_episode_2.avi" {
		t.Errorf("Must be sd link")
	}
	p = &Playlist{}
	err := p.SwitchToHD("hd")
	if err != ErrEmptyPlaylist {
		t.Errorf("Must be %v, but got %v", ErrEmptyPlaylist, err)
	}
	p = &Playlist{
		Items: []*Item{
			&Item{
				Title: "Episode 1 SD",
				File:  "http://sd/7f_episode_1.avi",
			},
			&Item{
				Title: "Episode 2 SD",
				File:  "http://sd/7f_episode_2.avi",
			},
		},
	}
	err = p.SwitchToHD("hd")
	if err != ErrHDNotFound {
		t.Errorf("Must be %v, but got %v", ErrHDNotFound, err)
	}
}

func TestDecodeLinks(t *testing.T) {
	tests := []struct {
		file string
		mp4  bool
	}{
		{
			file: "#2aHR0cDovL2RhdGEwOS1jZG4uZGF0YWxvY2sucnUvZmkybG0vMC83Zl9XZXN0d29ybGQuW1MwMUUwOV0uSEQ3MjAuRFVCLlt//b2xvbG8=xcXNzNDRdLmExLjI4LjExLjE2Lm1wNA==",
			mp4:  true,
		},
		{
			file: "#2aHRcDovL2!%$RhdGEwOS1jZG4uZGF0YWxvY2ucnUvZmkybG0vMC83Zl9XZXN0d29ybGQuW1MwMUUwOV0uSEQ3MjAuRFVCLlt//b2xvbG8=xcXNzNDRdLmExLjI4LjExLjE2Lm1wNA==",
			mp4:  false,
		},
		{
			file: "#2aHR//b2xvbG8=0cDovL2RhdGEwOS1jZG4uZGF0YWxvY2sucnUvZmkybG0vMC83Zl9XZXN0d29ybGQuW1MwMUUwOV0uSEQ3MjAuRFVCLltxcXNzNDRdLmExLjI4LjExLjE2Lm1wNA==",
			mp4:  true,
		},
		{
			file: "http://hostname/filename.mp4",
			mp4:  true,
		},
		{
			file: "http://hostname/filename",
			mp4:  false,
		},
	}
	for _, test := range tests {
		p := &Playlist{
			Items: []*Item{&Item{File: test.file}},
		}
		if err := p.DecodeLinks(); err != nil {
			t.Error(err)
		}
		ok := strings.HasSuffix(p.Items[0].File, ".mp4")
		if !ok && test.mp4 {
			t.Error("Must contains mp4 link")
		}
	}
}
