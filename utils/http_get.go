package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func HttpGetRaw(url string, headers map[string]string) ([]byte, error) {
	var body []byte
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return body, err
	}
	req.Header.Del("Accept-Encoding")
	req.Header.Set("User-Agent", RandomUserAgent())
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("error getting url: %s (%s)", url, resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return b, nil
}

func HttpGet(url string) (string, error) {
	var body string
	b, err := HttpGetRaw(url, map[string]string{})
	if err != nil {
		return body, err
	}
	body = string(b)
	body = strings.TrimSpace(body)
	return body, nil
}
