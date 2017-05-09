package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.1 Safari/603.1.30"
)

var (
	client http.Client
)

func HttpGetRaw(url string) ([]byte, error) {
	var body []byte
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return body, err
	}
	req.Header.Set("User-Agent", DefaultUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("error getting url: %s (%s)", url, resp.Status)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return b, nil
}

func HttpGet(url string) (string, error) {
	var body string
	b, err := HttpGetRaw(url)
	if err != nil {
		return body, err
	}
	body = string(b)
	body = strings.TrimSpace(body)
	return body, nil
}
