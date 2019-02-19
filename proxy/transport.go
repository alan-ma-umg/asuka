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

type Transport struct {
	*helper.Counting
	S               *SsAddr
	t               http.RoundTripper
	transportClosed bool

	//traffic size
	//TrafficIn  uint64
	//TrafficOut uint64

	Ping            time.Duration
	PingFailureRate float64

	RecentFewTimesResult []bool
}

func NewTransport(ssAddr *SsAddr) (*Transport, error) {
	instance := &Transport{S: ssAddr, t: createHttpTransport(ssAddr), Counting: &helper.Counting{}}
	return instance, nil
}

func createHttpTransport(SockInfo *SsAddr) *http.Transport {
	t := &http.Transport{
		MaxIdleConnsPerHost:   2,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	switch SockInfo.Type {
	case "local":
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
	case "ss", "ssr":
		SockInfo.WaitUntilConnected() //waiting
		dialer, err := proxy.SOCKS5("tcp", SockInfo.ClientAddr, nil, proxy.Direct)
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
		if transport.S.Type == "ss" || transport.S.Type == "ssr" {
			transport.S.Close()
		}

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
