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
	*url.URL
	Stop bool
}

type Transport struct {
	*helper.Counting
	S               *AddrInfo
	t               http.RoundTripper
	transportClosed bool

	RecentFewTimesResult []bool
}

func NewTransport(addr *AddrInfo) *Transport {
	return &Transport{S: addr, t: createHttpTransport(addr), Counting: &helper.Counting{}}
}

func createHttpTransport(SockInfo *AddrInfo) *http.Transport {
	t := &http.Transport{
		MaxIdleConnsPerHost:   2,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}

	switch SockInfo.Scheme {
	case "direct":
		t.Proxy = nil //disable system proxy
		t.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			//DualStack: true,
			FallbackDelay: time.Second,
		}).DialContext
	case "http", "https":
		t.Proxy = http.ProxyURL(SockInfo.URL) // with http proxy
		t.TLSHandshakeTimeout = time.Minute
		t.DialContext = (&net.Dialer{
			Timeout:   time.Minute,
			KeepAlive: time.Minute,
			//DualStack: true,
			FallbackDelay: time.Second,
		}).DialContext
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", SockInfo.Host, nil, proxy.Direct)
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
