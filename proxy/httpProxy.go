package proxy

import (
	"github.com/chenset/asuka/helper"
	"net/url"
	"strings"
)

func HttpProxyHandler() (ssAddr []*SsAddr) {
	for _, server := range helper.Env().HttpProxyServers {
		ss := &SsAddr{
			Enable: server.Enable,
			//EnablePing: server.EnablePing,
			Interval: server.Interval,
			Name:     server.Name,
			//Group:      server.Group,
			Type:       strings.ToLower(server.Type),
			ServerAddr: strings.ToLower(server.Server + ":" + server.ServerPort),
			openChan:   make(chan bool),
			closeChan:  make(chan bool, 1),
		}
		ssAddr = append(ssAddr, ss)
	}
	return
}

func HttpProxyParse(str string) (servers []*SsAddr) {
	str = strings.Replace(str, "\r\n", "\n", len(str))
	str = strings.Replace(str, "\r", "\n", len(str))
	for _, line := range strings.Split(strings.TrimSpace(str), "\n") {
		line = strings.ToLower(strings.TrimSpace(line))
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "http") {
			line = "http://" + line
		}

		urlAddr, err := url.Parse(line)
		if err != nil || urlAddr.Port() == "" {
			continue
		}

		serverAddr := ""
		if userInfoStr := urlAddr.User.String(); userInfoStr != "" {
			serverAddr += userInfoStr + "@"
		}

		servers = append(servers, &SsAddr{
			Enable:   true,
			Interval: 1,
			Name:     urlAddr.Hostname(),
			//Group:      "new",
			Type:       strings.ToLower(urlAddr.Scheme),
			ServerAddr: strings.ToLower(serverAddr + urlAddr.Hostname() + ":" + urlAddr.Port()),
		})

		//servers = append(servers, &HttpProxyServer{
		//	Enable:     true,
		//	EnablePing: true,
		//	Interval:   0,
		//	Name:       ipAddr.Hostname(),
		//	Group:      "httpProxy",
		//	Server:     ipAddr.Port(),
		//	ServerPort: ipAddr.Hostname(),
		//	Type:       "http",
		//})
	}

	return
	//b, _ := json.Marshal(servers)
	//return string(b)
}
