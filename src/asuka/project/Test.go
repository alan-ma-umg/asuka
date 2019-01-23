package project

import (
	"asuka/spider"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Test struct {
}

func (my *Test) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "http://hk.flysay.com:888/")
	}

	return links
}

// frequency
func (my *Test) Throttle(spider *spider.Spider) {
	spider.AddSleep(time.Duration(rand.Float64() * 2e9))
}

func (my *Test) RequestBefore(spider *spider.Spider) {
	//Referer
	if spider.CurrentRequest() != nil && spider.CurrentRequest().Referer() == "" {
		spider.CurrentRequest().Header.Set("Referer", my.EntryUrl()[0])
	}

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

// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
// 一般用于采集数据的地方
func (my *Test) ResponseSuccess(spider *spider.Spider) {

}

// queue
func (my *Test) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	return l.String()
}

func (my *Test) ResponseAfter(spider *spider.Spider) {
}
