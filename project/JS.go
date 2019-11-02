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

type JS struct {
	*Implement
	queueUrlLen int64
	showStr     string
	speedMin    time.Duration
	speedTotal  time.Duration
	speedMax    time.Duration
}

func (my *JS) InitBloomFilterCapacity() uint { return 1000000 }
func (my *JS) Init(d *Dispatcher) {
	go func() {
		for {
			time.Sleep(2e9)
			my.queueUrlLen, _ = database.Redis().LLen(my.Name() + "_" + helper.Env().Redis.URLQueueKey).Result()
		}
	}()

	//for _, s := range d.GetSpiders() {
	//	d.RemoveSpider(s)
	//}

	//for i := 0; i < helper.MaxInt(2, runtime.NumCPU()); i++ {
	//	uu, _ := url.Parse("direct://thread-" + strconv.Itoa(i))
	//	d.AddSpider(uu)
	//}

	my.speedMin = time.Hour
	my.showStr = "Waiting"
}

func (my *JS) Showing() string {
	return my.showStr
}

func (my *JS) Name() string {
	return "JS"
}

func (my *JS) EntryUrl() []string {
	var links []string

	for i := 0; i < 10000; i++ {
		links = append(links, "http://127.0.0.1:"+strings.Split(helper.Env().WEBListen, ":")[len(strings.Split(helper.Env().WEBListen, ":"))-1]+"/forever/")
	}

	return links
}
func (my *JS) Throttle(spider *spider.Spider) {
	spider.ResetSleep()
	if spider.LoadRate(5) > 10 {
		spider.AddSleep(time.Duration(rand.Float64() * 1e9))
	}

	spider.AddSleep(time.Second * 3)
}

func (my *JS) RequestBefore(spider *spider.Spider) {
	//spider.Client().Timeout = time.Second * 10
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *JS) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	//if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
	//	return false, nil
	//}
	//
	//return true, nil
	return true, nil
}

func (my *JS) ResponseAfter(spider *spider.Spider) {
	//if spider.CurrentResponse() != nil && spider.CurrentResponse().StatusCode == 200 {
	//	duration := spider.RequestEndTime.Sub(spider.RequestStartTime)
	//	if duration < my.speedMin {
	//		my.speedMin = duration
	//	}
	//	if duration > my.speedMax {
	//		my.speedMax = duration
	//	}
	//
	//	my.speedTotal += duration
	//	if spider.GetAccessCount() > spider.GetFailureCount() {
	//		my.showStr = "MIN: " + my.speedMin.Truncate(time.Microsecond).String() + "  MAX: " + my.speedMax.Truncate(time.Microsecond).String() + "  AVG: " + (my.speedTotal / time.Duration(spider.GetAccessCount()-spider.GetFailureCount())).Truncate(time.Microsecond).String()
	//	}
	//}

	//my.Implement.ResponseAfter(spider)
}

// queue
func (my *JS) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if my.queueUrlLen > 100000 {
		return
	}

	return l.String()
}
