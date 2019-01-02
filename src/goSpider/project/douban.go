package project

import (
	"goSpider/spider"
	"net/url"
)

type DouBan struct {
}

func (my *DouBan) EntryUrl() []string {
	return []string{
		"http://192.168.100.125:888/forever",
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
	}
}

func (my *DouBan) RequestBefore(spider *spider.Spider) {
}

func (my *DouBan) ResponseAfter(spider *spider.Spider) {
}

// queue
func (my *DouBan) EnqueueFilter(spider *spider.Spider, l *url.URL) bool {
	return true
}

// frequency
func (my *DouBan) Throttle(spider *spider.Spider) {
}
