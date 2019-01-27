package project

import (
	"asuka/database"
	"asuka/helper"
	"asuka/spider"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Test struct {
	*Implement
	queueUrlLen int64
}

func (my *Test) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "http://hk.flysay.com:888/")
	}

	go func() {
		t := time.NewTicker(time.Second * 5)
		for {
			<-t.C
			my.queueUrlLen, _ = database.Redis().LLen(strings.Split(reflect.TypeOf(my).String(), ".")[1] + "_" + helper.Env().Redis.URLQueueKey).Result()
		}
	}()

	return links
}

// frequency
var times int

func (my *Test) Throttle(spider *spider.Spider) {
	spider.AddSleep(time.Duration(rand.Float64() * 1e9))

	if times < 200 {
		times++
		spider.UpdateTransport()
	}
}

func (my *Test) RequestBefore(spider *spider.Spider) {
	spider.Client().Timeout = time.Second * 10
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *Test) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}

	return true, nil
}

// queue
func (my *Test) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if my.queueUrlLen > 20000 {
		return
	}

	return l.String()
}
