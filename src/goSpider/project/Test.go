package project

import (
	"goSpider/spider"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Test struct {
}

func (my *Test) EntryUrl() []string {
	return []string{
		//"http://192.168.100.125:8188/forever",
		//"http://192.168.100.125:8188/forever",
		//"http://192.168.100.125:8188/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"http://192.168.100.125:888/forever",
		//"https://book.douban.com",
		//"https://movie.douban.com",
		//"https://www.zhihu.com/explore",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		//"http://10.0.0.180:888/forever/",
		"http://hk.flysay.com:888/",
		"http://hk.flysay.com:888",
		"http://hk.flysay.com:888",
		"http://hk.flysay.com:888",
		"http://hk.flysay.com:888",
		"http://hk.flysay.com:888",
		"http://hk.flysay.com:888",
		"http://hk.flysay.com:888",
		"http://hk.flysay.com:888",
	}
}

func (my *Test) RequestBefore(spider *spider.Spider) {
	//Referer
	if spider.CurrentRequest != nil && spider.CurrentRequest.Referer() == "" {
		spider.CurrentRequest.Header.Set("Referer", my.EntryUrl()[0])
	}

	spider.Client.Timeout = time.Second
}

func (my *Test) ResponseAfter(spider *spider.Spider) {
	//free the memory
	if len(spider.RequestsMap) > 10 {
		spider.Client.Jar, _ = cookiejar.New(nil)
		spider.RequestsMap = map[string]*http.Request{}
	}
}

// queue
func (my *Test) EnqueueFilter(spider *spider.Spider, l *url.URL) bool {
	return true
}

// frequency
func (my *Test) Throttle(spider *spider.Spider) {
	//time.Sleep(10e9)
}
