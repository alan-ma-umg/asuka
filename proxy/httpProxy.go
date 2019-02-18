package proxy

import (
	"github.com/chenset/asuka/helper"
	"strings"
)

func HttpProxyHandler() (ssAddr []*SsAddr) {
	for _, server := range helper.Env().HttpProxyServers {
		ss := &SsAddr{
			Enable:     server.Enable,
			EnablePing: server.EnablePing,
			Interval:   server.Interval,
			Name:       server.Name,
			Group:      server.Group,
			Type:       strings.ToLower(server.Type),
			ServerAddr: strings.ToLower(server.Server + ":" + server.ServerPort),
			openChan:   make(chan bool),
			closeChan:  make(chan bool, 1),
		}
		ssAddr = append(ssAddr, ss)
	}
	return
}
