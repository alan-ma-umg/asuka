package project

import (
	"asuka/spider"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Test struct {
}

func (my *Test) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "http://hk.flysay.com:888/")
	}

	return links
}

var times = 0

// frequency
func (my *Test) Throttle(spider *spider.Spider) {
	spider.AddSleep(time.Duration(rand.Float64() * 2e9))

	if times < 10 {
		times++
		spider.UpdateTransport()
	}
}

func (my *Test) RequestBefore(spider *spider.Spider) {
	//Referer
	if spider.CurrentRequest != nil && spider.CurrentRequest.Referer() == "" {
		spider.CurrentRequest.Header.Set("Referer", my.EntryUrl()[0])
	}

	spider.Client().Timeout = time.Second * 10
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *Test) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}

	return true, nil
}

// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
// 一般用于采集数据的地方
func (my *Test) ResponseSuccess(spider *spider.Spider) {

}

// queue
func (my *Test) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	return l.String()
}

func (my *Test) ResponseAfter(spider *spider.Spider) {
	//spider.Transport.T.(*http.Transport).DisableKeepAlives = false
	//spider.Transport.T.(*http.Transport).CloseIdleConnections()
	//spider.Transport.S.CloseChan <- true
	//spider.Transport.S.Listener.Close()
	//if rand.Intn(10) == 2 {
	//	spider.Transport.T.(*http.Transport).DisableKeepAlives = true
	//} else {
	//	spider.Transport.T.(*http.Transport).DisableKeepAlives = false
	//}
	//spider.Transport.T.(*http.Transport).CloseIdleConnections()
	//dialer, err := proxy.SOCKS5("tcp", spider.Transport.S.ClientAddr, nil, proxy.Direct)
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
	//	return
	//}
	//spider.Transport.T
	//http transport
	//t := &http.Transport{
	//MaxIdleConnsPerHost: 2,
	//MaxIdleConns:        10,
	//IdleConnTimeout:     20 * time.Second,
	//TLSHandshakeTimeout: 10 * time.Second,
	//
	//DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
	//	return dialer.Dial(network, addr)
	//},
	//}
	//
	//spider.Transport.T = nil
	//spider.Transport.T = t

	//free the memory
	//if len(spider.RequestsMap) > 10 {
	//	spider.Client.Jar, _ = cookiejar.New(nil)
	//	spider.RequestsMap = map[string]*http.Request{}
	//}
}
