package project

import (
	"goSpider/spider"
	"net/url"
)

type DouBan struct {
}

func (my *DouBan) EntryUrl() []string {
	return []string{
		//"http://192.168.100.125:8188/forever",
		//"http://192.168.100.125:8188/forever",
		//"http://192.168.100.125:8188/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
		"http://192.168.100.125:888/forever",
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
	}
}

func (my *DouBan) RequestBefore(spider *spider.Spider) {
	//Referer
	if spider.CurrentRequest != nil && spider.CurrentRequest.Referer() == "" {
		spider.CurrentRequest.Header.Set("Referer", my.EntryUrl()[0])
	}
}

func (my *DouBan) ResponseAfter(spider *spider.Spider) {
}

// queue
func (my *DouBan) EnqueueFilter(spider *spider.Spider, l *url.URL) bool {
	return true
}

// frequency
func (my *DouBan) Throttle(spider *spider.Spider) {
	//time.Sleep(2e9)
}
