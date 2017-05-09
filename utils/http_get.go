package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpGetRaw(url string) ([]byte, error) {
	var body []byte
	resp, err := http.Get(url)
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
