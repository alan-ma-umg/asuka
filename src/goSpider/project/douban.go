package project

import (
	"goSpider/spider"
	"net/url"
	"time"
	"net/http/cookiejar"
	"net/http"
)

type DouBan struct {
}

func (my *DouBan) EntryUrl() []string {
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

func (my *DouBan) RequestBefore(spider *spider.Spider) {
	//Referer
	if spider.CurrentRequest != nil && spider.CurrentRequest.Referer() == "" {
		spider.CurrentRequest.Header.Set("Referer", my.EntryUrl()[0])
	}

	spider.Client.Timeout = time.Second
}

func (my *DouBan) ResponseAfter(spider *spider.Spider) {
	//free the memory
	if len(spider.RequestsMap) > 10 {
		spider.Client.Jar, _ = cookiejar.New(nil)
		spider.RequestsMap = map[string]*http.Request{}
	}
}

// queue
func (my *DouBan) EnqueueFilter(spider *spider.Spider, l *url.URL) bool {
	return true
}

// frequency
func (my *DouBan) Throttle(spider *spider.Spider) {
	time.Sleep(5e9)
}
