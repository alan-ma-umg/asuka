package project

import (
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Death struct {
	*Implement
	*SpeedShowing
	*SpiderThrottle
	*Setting
}

func (my *Death) Name() string { return "Death" }
func (my *Death) Throttle(spider *spider.Spider) {
	spider.FailureLevel = 0
	spider.ResetSleep()
	my.SpiderThrottle.Throttle(spider)
}
func (my *Death) Init(d *Dispatcher) {
	my.SpeedShowing = &SpeedShowing{}
	my.SpiderThrottle = &SpiderThrottle{}
	my.Setting = &Setting{Option: &SettingOption{}}

	my.SetThrottleSpeed(1)
	go func() {
		for {
			time.Sleep(10e9)

			if database.Redis().LLen(my.Name()+"_"+helper.Env().Redis.URLQueueKey).Val() > 100000 {
				d.GetQueue().CleanQueue()
				d.GetQueue().Enqueue(my.EntryUrl())
			}

			if d.GetQueue().FailureLen() > 10000 {
				d.GetQueue().CleanFailure()
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
func (my *Death) ResponseAfter(spider *spider.Spider) {}
func (my *Death) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "http://127.0.0.1:"+strings.Split(helper.Env().WEBListen, ":")[len(strings.Split(helper.Env().WEBListen, ":"))-1]+"/forever/")
	}

	return links
}
