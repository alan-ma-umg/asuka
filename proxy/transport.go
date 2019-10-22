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
	u                    *url.URL
	t                    http.RoundTripper
	transportClosed      bool
	RecentFewTimesResult []bool
}

func NewTransport(u *url.URL) *Transport {
	return &Transport{u: u, t: createHttpTransport(u), Counting: &helper.Counting{}}
}

func createHttpTransport(u *url.URL) *http.Transport {
	t := &http.Transport{
		MaxIdleConnsPerHost:   2,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}

	switch u.Scheme {
	case "direct":
		t.Proxy = nil //disable system proxy
		t.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			//DualStack: true,
			FallbackDelay: time.Second,
		}).DialContext
	case "http", "https":
		t.Proxy = http.ProxyURL(u) // with http proxy
		t.TLSHandshakeTimeout = time.Minute
		t.DialContext = (&net.Dialer{
			Timeout:   time.Minute,
			KeepAlive: time.Minute,
			//DualStack: true,
			FallbackDelay: time.Second,
		}).DialContext
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", u.Host, nil, proxy.Direct)
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
		transport.t = createHttpTransport(transport.u)
		transport.transportClosed = false
	}

	return transport.t.(*http.Transport)
}
