package proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func HttpProxyParse(scheme string, str string) (urls []*url.URL) {
	str = strings.Replace(str, "\r\n", "\n", len(str))
	str = strings.Replace(str, "\r", "\n", len(str))
	for _, line := range strings.Split(strings.TrimSpace(str), "\n") {
		line = strings.ToLower(strings.TrimSpace(line))
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "http") {
			line = scheme + "://" + line
		}

		urlAddr, err := url.Parse(line)
		if err != nil || urlAddr.Port() == "" {
			continue
		}

		urls = append(urls, urlAddr)
	}

	return
}
func ScyllaJson(jsonStr string) (urls []*url.URL) {
	var m map[string][]struct {
		Ip       string
		Port     int
		Is_https bool
		Is_valid bool
	}
	json.Unmarshal([]byte(jsonStr), &m)

	if proxies, ok := m["proxies"]; ok {
		for _, item := range proxies {
			if !item.Is_valid {
				continue
			}
			scheme := "http"
			if item.Is_https {
				scheme = "https"
			}
			urls = append(urls, HttpProxyParse(scheme, item.Ip+":"+strconv.Itoa(item.Port))...)
		}
	}

	return
}

func ScyllaApi(url string) (urls []*url.URL) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	return ScyllaJson(string(body))
}
