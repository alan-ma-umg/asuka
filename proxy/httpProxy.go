package proxy

import (
	"github.com/chenset/asuka/helper"
	"net/url"
	"strings"
)

func HttpProxyHandler() (ssAddr []*AddrInfo) {
	for _, server := range helper.Env().HttpProxyServers {
		ss := &AddrInfo{
			Enable: server.Enable,
			//EnablePing: server.EnablePing,
			Interval: server.Interval,
			Name:     server.Name,
			//Group:      server.Group,
			Type:       strings.ToLower(server.Type),
			ServerAddr: strings.ToLower(server.Server + ":" + server.ServerPort),
			openChan:   make(chan bool),
			//closeChan:  make(chan bool, 1),
		}
		ssAddr = append(ssAddr, ss)
	}
	return
}

func HttpProxyParse(scheme string, str string) (servers []*AddrInfo) {
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

		serverAddr := ""
		if userInfoStr := urlAddr.User.String(); userInfoStr != "" {
			serverAddr += userInfoStr + "@"
		}

		servers = append(servers, &AddrInfo{
			Enable:   true,
			Interval: 1,
			Name:     urlAddr.Hostname(),
			//Group:      "new",
			Type:       strings.ToLower(urlAddr.Scheme),
			ServerAddr: strings.ToLower(serverAddr + urlAddr.Hostname() + ":" + urlAddr.Port()),
		})
	}

	return
}
