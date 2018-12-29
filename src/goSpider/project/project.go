package project

import (
	"goSpider/spider"
	"net/url"
)

type Project interface {
	EntryUrl() []string

	// request before & after
	RequestBefore(spider *spider.Spider)
	ResponseAfter(spider *spider.Spider)

	// queue
	EnqueueFilter(spider *spider.Spider, l *url.URL) bool

	// frequency
	Throttle(spider *spider.Spider)
}
