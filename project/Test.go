package project

import (
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Test struct {
	*Implement
	queueUrlLen int64
	showStr     string
	speedMin    time.Duration
	speedTotal  time.Duration
	speedMax    time.Duration
}

func (my *Test) Init(d *Dispatcher) {
	go func() {
		for {
			time.Sleep(10e9)
			my.queueUrlLen, _ = database.Redis().LLen(my.Name() + "_" + helper.Env().Redis.URLQueueKey).Result()
		}
	}()

	my.speedMin = time.Hour
	my.showStr = "Waiting"
}

func (my *Test) Showing() string {
	return my.showStr
}

func (my *Test) Name() string {
	return "KEI"
}

func (my *Test) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "https://gg.flysay.com/forever/")
	}

	return links
}
func (my *Test) Throttle(spider *spider.Spider) {
	if spider.LoadRate(5) > 10 {
		spider.AddSleep(time.Duration(rand.Float64() * 1e9))
	}
}

func (my *Test) RequestBefore(spider *spider.Spider) {
	spider.SetRequestTimeout(time.Second * 10)
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *Test) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}

	return true, nil
}

func (my *Test) ResponseAfter(spider *spider.Spider) {
	if spider.CurrentResponse() != nil && spider.CurrentResponse().StatusCode == 200 {
		duration := spider.RequestEndTime.Sub(spider.RequestStartTime)
		if duration < my.speedMin {
			my.speedMin = duration
		}
		if duration > my.speedMax {
			my.speedMax = duration
		}

		my.speedTotal += duration
		if spider.GetAccessCount() > spider.GetFailureCount() {
			my.showStr = "MIN: " + my.speedMin.Truncate(time.Microsecond).String() + "  MAX: " + my.speedMax.Truncate(time.Microsecond).String() + "  AVG: " + (my.speedTotal / time.Duration(spider.GetAccessCount()-spider.GetFailureCount())).Truncate(time.Microsecond).String()
		}
	}

	//my.Implement.ResponseAfter(spider)
}

// queue
func (my *Test) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if my.queueUrlLen > 20000 {
		return
	}

	return l.String()
}
