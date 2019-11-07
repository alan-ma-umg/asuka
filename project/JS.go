package project

import (
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"net/url"
	"time"
)

type JS struct {
	*Implement
	*SpeedShowing
}

func (my *JS) InitBloomFilterCapacity() uint { return 1000000 }
func (my *JS) Init(d *Dispatcher) {
	my.SpeedShowing = &SpeedShowing{}
	go func() {
		for {
			time.Sleep(20e9)
			queueUrlLen, _ := database.Redis().LLen(my.Name() + "_" + helper.Env().Redis.URLQueueKey).Result()
			if queueUrlLen < 10000 {
				d.queue.Enqueue(my.EntryUrl())
			}
		}
	}()
}

func (my *JS) Fetch(spider *spider.Spider, u *url.URL) (summary *spider.Summary, err error) {
	return spider.ChromeFetch(u)
}

func (my *JS) Name() string {
	return "JS"
}

func (my *JS) EntryUrl() []string {
	var links []string

	for i := 0; i < 10000; i++ {
		links = append(links, "https://asuka.flysay.com/")
	}

	return links
}

func (my *JS) Throttle(spider *spider.Spider) {
	spider.AddSleep(time.Minute * 3)
}

// queue
func (my *JS) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	return
}
