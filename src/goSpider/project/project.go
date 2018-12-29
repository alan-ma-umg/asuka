package project

import (
	"goSpider/spider"
	"net/url"
)

type Project interface {
	EntryUrl() []string

	// session
	NeedToLogin(spider *spider.Spider) bool
	IsLogin(spider *spider.Spider) bool
	Login(spider *spider.Spider)

	// queue
	EnqueueFilter(spider *spider.Spider, l *url.URL) bool

	// frequency
	NeedToPause(spider *spider.Spider) bool
	Throttle(spider *spider.Spider)
}
