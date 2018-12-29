package project

import (
	"goSpider/spider"
	"net/url"
)

type DouBan struct {
}

func (my *DouBan) EntryUrl() []string {
	return []string{
		"https://www.douban.com",
		"https://book.douban.com",
		"https://movie.douban.com",
		"https://www.zhihu.com/explore",
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
	//time.Sleep(2e9)
	//ch := make(chan int)
	//<-ch
}
