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
	queueUrlLen int64

	DefaultShowingEnable bool //todo 通过组合方式加入进来
	DefaultShowing       string
	DefaultSpeedCount    uint
	DefaultSpeedTotal    time.Duration
	DefaultSpeedMin      time.Duration
	DefaultSpeedAvgDiv   time.Duration
	DefaultSpeedMax      time.Duration
}

func (my *Death) Showing() string { my.DefaultShowingEnable = true; return my.DefaultShowing }
func (my *Death) Name() string    { return "Death" }
func (my *Death) Init(d *Dispatcher) {
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

func (my *Death) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "http://127.0.0.1:"+strings.Split(helper.Env().WEBListen, ":")[len(strings.Split(helper.Env().WEBListen, ":"))-1]+"/forever/")
	}

	return links
}

func (my *Death) ResponseAfter(spider *spider.Spider) {
}

// queue
func (my *Death) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if my.queueUrlLen > 100000 {
		return
	}

	return l.String()
}

func (my *Death) ResponseSuccess(spider *spider.Spider) {
	if my.DefaultShowingEnable {
		if my.DefaultShowing == "" {
			my.DefaultSpeedMin = time.Hour
		}
		duration := spider.RequestEndTime.Sub(spider.RequestStartTime)
		if duration < my.DefaultSpeedMin {
			my.DefaultSpeedMin = duration
		}
		if duration > my.DefaultSpeedMax {

			my.DefaultSpeedMax = duration
		}
		if my.DefaultSpeedAvgDiv == 0 {
			my.DefaultSpeedAvgDiv = duration
		}
		my.DefaultSpeedAvgDiv = (my.DefaultSpeedAvgDiv + duration) / 2

		my.DefaultSpeedTotal += duration
		my.DefaultSpeedCount++
		my.DefaultShowing = "MIN: " + my.DefaultSpeedMin.Truncate(time.Microsecond).String() + "  MAX: " + my.DefaultSpeedMax.Truncate(time.Microsecond).String() + "  AVG: " + (my.DefaultSpeedTotal / time.Duration(my.DefaultSpeedCount)).Truncate(time.Microsecond).String() + " / " + my.DefaultSpeedAvgDiv.Truncate(time.Microsecond).String()
	}
}
