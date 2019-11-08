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

type Forever struct {
	*Implement
	*SpeedShowing
	*SpiderThrottle
}

func (my *Forever) Name() string { return "Mai" }
func (my *Forever) Init(d *Dispatcher) {
	my.SpeedShowing = &SpeedShowing{}
	my.SpiderThrottle = &SpiderThrottle{}
	my.SetThrottleSpeed(.1)
	go func() {
		for {
			time.Sleep(30e9)
			if database.Redis().LLen(my.Name()+"_"+helper.Env().Redis.URLQueueKey).Val() > 100000 {
				d.GetQueue().CleanQueue()
				d.GetQueue().Enqueue(my.EntryUrl())
			}
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
func (my *Forever) ResponseAfter(spider *spider.Spider) {}
func (my *Forever) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "https://h.flysay.com/forever/")
	}

	return links
}
