package project

import (
	"goSpider/spider"
	"net/url"
)

type DouBan struct {
}

func (douBan *DouBan) EntryUrl() []string {
	return []string{}
}

// session
func (douBan *DouBan) NeedToLogin(spider *spider.Spider) bool {
	return false
}
func (douBan *DouBan) IsLogin(spider *spider.Spider) bool {
	return false
}
func (douBan *DouBan) Login(spider *spider.Spider) {
}

// queue
func (douBan *DouBan) EnqueueFilter(spider *spider.Spider, l *url.URL) bool {
	return true
}

// frequency
func (douBan *DouBan) NeedToPause(spider *spider.Spider) bool {
	return false
}
func (douBan *DouBan) Throttle(spider *spider.Spider) {

}
