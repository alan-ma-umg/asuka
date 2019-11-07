package project

import (
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"net/url"
	"runtime"
	"strconv"
	"time"
)

type Test struct {
	*Implement
	*SpeedShowing
	*SpiderThrottle
	queueUrlLen int64
}

func (my *Test) Name() string { return "Kei" }
func (my *Test) Init(d *Dispatcher) {
	my.SpeedShowing = &SpeedShowing{}
	my.SpiderThrottle = &SpiderThrottle{}
	my.SetThrottleSpeed(.1)
	go func() {
		for {
			time.Sleep(5e9)
			my.queueUrlLen, _ = database.Redis().LLen(my.Name() + "_" + helper.Env().Redis.URLQueueKey).Result()
		}
	}()

	for _, s := range d.GetSpiders() {
		d.RemoveSpider(s)
	}

	for i := 0; i < helper.MaxInt(10, runtime.NumCPU()*10); i++ {
		uu, _ := url.Parse("direct://thread-" + strconv.Itoa(i))
		d.AddSpider(uu)
	}
}
func (my *Test) ResponseAfter(spider *spider.Spider) {}
func (my *Test) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "https://gg.flysay.com/forever/")
	}

	return links
}

// queue
func (my *Test) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if my.queueUrlLen > 100000 {
		return
	}
	return l.String()
}
