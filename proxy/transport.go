package proxy

import (
	"context"
	"github.com/chenset/asuka/helper"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type AddrInfo struct {
	Enable bool
	//EnablePing bool
	Interval float64
	Type     string
	Name     string
	//Group      string
	ServerAddr string
	//ClientAddr string
	//TrafficIn   uint64
	//TrafficOut  uint64
	//Connections int
	//listener  net.Listener
	openChan chan bool
	//closeChan chan bool
	//Status    int //0 init, 10 close, 20 socks connected|waiting, 30 remote established
}

type Transport struct {
	*helper.Counting
	S               *AddrInfo
	t               http.RoundTripper
	transportClosed bool

	//traffic size
	//TrafficIn  uint64
	//TrafficOut uint64

	//Ping            time.Duration
	//PingFailureRate float64

	RecentFewTimesResult []bool
}

func NewTransport(ssAddr *AddrInfo) (*Transport, error) {
	instance := &Transport{S: ssAddr, t: createHttpTransport(ssAddr), Counting: &helper.Counting{}}
	return instance, nil
}

func createHttpTransport(SockInfo *AddrInfo) *http.Transport {
	t := &http.Transport{
		MaxIdleConnsPerHost:   2,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	switch SockInfo.Type {
	case "direct":
		t.Proxy = nil //disable system proxy
		t.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
	case "http", "https":
		proxyURL, err := url.Parse(SockInfo.Type + "://" + SockInfo.ServerAddr)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		t.Proxy = http.ProxyURL(proxyURL) // with http proxy
		t.TLSHandshakeTimeout = time.Minute
		t.DialContext = (&net.Dialer{
			Timeout:   time.Minute,
			KeepAlive: time.Minute,
			DualStack: true,
		}).DialContext
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", SockInfo.ServerAddr, nil, proxy.Direct)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		t.Proxy = nil //disable system proxy
		t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}
	return t
}

func (transport *Transport) Close() {
	if !transport.transportClosed {
		transport.t.(*http.Transport).DisableKeepAlives = true
		transport.t.(*http.Transport).CloseIdleConnections()
		transport.transportClosed = true
	}
}

func (transport *Transport) Connect() *http.Transport {
	if transport.transportClosed {
		transport.t = createHttpTransport(transport.S)
		transport.transportClosed = false
	}

	return transport.t.(*http.Transport)
}
