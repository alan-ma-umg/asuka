package proxy

import (
	"context"
	"github.com/chenset/asuka/helper"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Transport struct {
	t               http.RoundTripper
	transportClosed bool
}

func NewTransport(u *url.URL) *Transport {
	return &Transport{t: createHttpTransport(u)}
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
		t.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			separator := strings.LastIndex(address, ":")
			return (&net.Dialer{
				Timeout:   time.Minute,
				KeepAlive: time.Minute,
				//DualStack: true,
				FallbackDelay: time.Second,
			}).DialContext(ctx, network, helper.GetDNSCache().Lookup(address[:separator])+address[separator:])
		}
	case "http", "https":
		t.Proxy = http.ProxyURL(u) // with http proxy
		t.TLSHandshakeTimeout = time.Minute

		t.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			separator := strings.LastIndex(address, ":")
			return (&net.Dialer{
				Timeout:   time.Minute,
				KeepAlive: time.Minute,
				//DualStack: true,
				FallbackDelay: time.Second,
			}).DialContext(ctx, network, helper.GetDNSCache().Lookup(address[:separator])+address[separator:])
		}
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", u.Host, nil, proxy.Direct)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		t.Proxy = nil //disable system proxy
		t.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			separator := strings.LastIndex(address, ":")
			return dialer.Dial(network, helper.GetDNSCache().Lookup(address[:separator])+address[separator:])
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

func (transport *Transport) Connect(u *url.URL) *http.Transport {
	if transport.transportClosed {
		transport.t = createHttpTransport(u)
		transport.transportClosed = false
	}

	return transport.t.(*http.Transport)
}
