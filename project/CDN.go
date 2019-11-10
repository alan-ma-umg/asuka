package project

import (
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

type CDN struct {
	*Implement
	*SpeedShowing
	*SpiderThrottle
}

func (my *CDN) Name() string { return "CDN" }
func (my *CDN) Init(d *Dispatcher) {
	my.SpeedShowing = &SpeedShowing{}
	my.SpiderThrottle = &SpiderThrottle{}
	my.SetThrottleSpeed(0.01)
	go func() {
		for {
			time.Sleep(20e9)
			if database.Redis().LLen(my.Name()+"_"+helper.Env().Redis.URLQueueKey).Val() < 1000 {
				d.queue.Enqueue(my.EntryUrl())
			}
		}
	}()
}

func (my *CDN) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, []string{
			"https://cdn.jsdelivr.net/gh/chenset/asuka@master/web/templates/static/asuka.css?i=" + strconv.Itoa(time.Now().Second()), "https://cdn.jsdelivr.net/gh/chenset/asuka@master/web/templates/static/asuka.js?i=" + strconv.Itoa(time.Now().Second()),
			"https://cdn.jsdelivr.net/gh/chenset/asuka@latest/web/templates/static/asuka.css?i=" + strconv.Itoa(time.Now().Second()), "https://cdn.jsdelivr.net/gh/chenset/asuka@latest/web/templates/static/asuka.js?i=" + strconv.Itoa(time.Now().Second()),
			"https://cdn.jsdelivr.net/gh/chenset/asuka/web/templates/static/asuka.css?i=" + strconv.Itoa(time.Now().Second()), "https://cdn.jsdelivr.net/gh/chenset/asuka/web/templates/static/asuka.js?i=" + strconv.Itoa(time.Now().Second()),
		}...)
	}

	return links
}

func (my *CDN) RequestBefore(spider *spider.Spider) {
	//Referer
	spider.CurrentRequest().Header.Set("Referer", "https://pages.github.com/"+strconv.Itoa(rand.Intn(1000000000)))
	spider.CurrentRequest().Header.Set("cache-control", "no-cache")
	spider.CurrentRequest().Header.Set("pragma", "no-cache")
}

func (my *CDN) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	return
}
