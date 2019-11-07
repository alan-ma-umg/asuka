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
}

func (my *CDN) Name() string { return "CDN" }
func (my *CDN) Init(d *Dispatcher) {
	my.SpeedShowing = &SpeedShowing{}
	go func() {
		for {
			time.Sleep(20e9)
			queueUrlLen, _ := database.Redis().LLen(my.Name() + "_" + helper.Env().Redis.URLQueueKey).Result()
			if queueUrlLen < 1000 {
				d.queue.Enqueue(my.EntryUrl())
			}
		}
	}()
}

func (my *CDN) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, []string{
			"https://asuka.flysay.com/",
			"https://cdn.jsdelivr.net/gh/chenset/asuka/web/templates/static/asuka.css?i=" + time.Now().String() + strconv.Itoa(i), "https://cdn.jsdelivr.net/gh/chenset/asuka/web/templates/static/asuka.js?i=" + time.Now().String() + strconv.Itoa(i),
		}...)
	}

	return links
}
func (my *CDN) Throttle(spider *spider.Spider) {
	spider.AddSleep(time.Duration(rand.Float64() * 120e9))
}

func (my *CDN) RequestBefore(spider *spider.Spider) {
	//Referer
	spider.CurrentRequest().Header.Set("Referer", "https://pages.github.com/")
}

func (my *CDN) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	return
}

//func (my *CDN) ResponseAfter(spider *spider.Spider) {}
