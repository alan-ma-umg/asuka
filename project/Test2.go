package project

import (
	"asuka/database"
	"asuka/helper"
	"asuka/spider"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Test2 struct {
	*Implement
	queueUrlLen int64
}

//var test2BaseUrl = "http://z.flysay.com:888/"
var test2BaseUrl = "http://hk.flysay.com:88/"

func (my *Test2) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, test2BaseUrl)
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
func (my *Test2) Throttle(spider *spider.Spider) {
	//spider.AddSleep(time.Duration(rand.Float64() * 1e9))
}

func (my *Test2) RequestBefore(spider *spider.Spider) {
	spider.Client().Timeout = time.Minute
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *Test2) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}

	return true, nil
}

// queue
func (my *Test2) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if my.queueUrlLen > 20000 {
		return
	}

	if !strings.HasPrefix(strings.ToLower(l.String()), test2BaseUrl) {
		return
	}

	return l.String()
}