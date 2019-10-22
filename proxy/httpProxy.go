package proxy

import (
	"net/url"
	"strings"
)

//func HttpProxyHandler() (addr []*AddrInfo) {
//	for _, server := range helper.Env().HttpProxyServers {
//
//		addr = append(addr, &AddrInfo{
//			Enable: server.Enable,
//			//EnablePing: server.EnablePing,
//			Interval: server.Interval,
//			Name:     server.Name,
//			//Group:      server.Group,
//			Type:       strings.ToLower(server.Type),
//			ServerAddr: strings.ToLower(server.Server + ":" + server.ServerPort),
//			//openChan:   make(chan bool),
//			//closeChan:  make(chan bool, 1),
//		})
//	}
//	return
//}

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
