package proxy

import (
	"net/url"
	"strings"
)

func Socks5ProxyParse(str string) (servers []*AddrInfo) {
	str = strings.Replace(str, "\r\n", "\n", len(str))
	str = strings.Replace(str, "\r", "\n", len(str))
	for _, line := range strings.Split(strings.TrimSpace(str), "\n") {
		line = strings.ToLower(strings.TrimSpace(line))
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "socks5") {
			line = "socks5://" + line
		}

		urlAddr, err := url.Parse(line)
		if err != nil || urlAddr.Port() == "" {
			continue
		}
		servers = append(servers, &AddrInfo{URL: urlAddr})
	}

	return
}
