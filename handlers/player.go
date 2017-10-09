package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

var (
	allowHdReplacer   = strings.NewReplacer("swichHDno", "swichHD")
	prerollCodeRegexp = regexp.MustCompile(`\<script\ type\=\"text\/javascript\"\>var.*\<\/script\>`)
)

func playerRewriteBody(resp *http.Response) (err error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	b = prerollCodeRegexp.ReplaceAll(b, nil)
	b = []byte(allowHdReplacer.Replace(string(b)))
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func PlayerHandler(s *server.Server) http.Handler {
	u, _ := url.Parse(s.AbsoluteLink("/"))
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	reverseProxy.ModifyResponse = playerRewriteBody
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = u.Hostname()
		r.Header.Set("User-Agent", utils.RandomUserAgent())
		r.Header.Del("Accept-Encoding")
		reverseProxy.ServeHTTP(w, r)
	})
}
