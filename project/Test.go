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
}

func (my *Test) Name() string { return "Kei" }
func (my *Test) Init(d *Dispatcher) {
	go func() {
		for {
			time.Sleep(10e9)
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
func (my *Test) ResponseAfter(spider *spider.Spider) {}
func (my *Test) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "https://gg.flysay.com/forever/")
	}

	return links
}
