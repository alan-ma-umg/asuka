package project

import (
	"asuka/database"
	"asuka/helper"
	"asuka/spider"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Test struct {
	*Implement
	queueUrlLen      int64
	spiderChannel    chan *spider.Spider
	spiderLimit      map[string]int
	spiderLimitMutex sync.Mutex
}

func (my *Test) EntryUrl() []string {
	var links []string
	my.spiderChannel = make(chan *spider.Spider, 5)
	my.spiderLimit = make(map[string]int)

	for i := 0; i < 1000; i++ {
		links = append(links, "http://hk.flysay.com:88/")
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
func (my *Test) inQueue(spider *spider.Spider) {
	my.spiderLimitMutex.Lock()
	_, ok := my.spiderLimit[spider.Transport.S.ServerAddr]
	my.spiderLimitMutex.Unlock()

	if !ok {
		my.spiderChannel <- spider //stuck
		spider.Test = 10
		fmt.Println("in")
	}

	my.spiderLimitMutex.Lock()
	my.spiderLimit[spider.Transport.S.ServerAddr]++
	my.spiderLimitMutex.Unlock()
}

func (my *Test) release(spider *spider.Spider) {

	my.spiderLimitMutex.Lock()
	limit := my.spiderLimit[spider.Transport.S.ServerAddr]
	my.spiderLimitMutex.Unlock()

	if spider.Stop || limit >= 5 || spider.FailureLevel > 1 {
		<-my.spiderChannel
		fmt.Println("out")
		spider.Test = 20
		spider.Transport.Close()
		my.spiderLimitMutex.Lock()
		//spider.AddSleep(time.Duration(rand.Float64() * 10e9 * float64(my.spiderLimit[spider.Transport.S.ServerAddr])))
		delete(my.spiderLimit, spider.Transport.S.ServerAddr)
		my.spiderLimitMutex.Unlock()
	}
}

func (my *Test) ResponseAfter(spider *spider.Spider) {
	//my.release(spider)

	//if spider.FailureLevel > 1 {
	//	my.release(spider)
	//}
}
func (my *Test) Throttle(spider *spider.Spider) {
	//my.inQueue(spider)
	//my.release(spider)

	//spider.AddSleep(time.Duration(rand.Float64() * 60e9))

	//if times < 200 {
	//	times++
	//	spider.Transport.Close()
	//}

	//if spider.Transport.LoadRate(60) > 5 {
	//	spider.AddSleep((60 / 5) * 1e9)
	//}
}

func (my *Test) RequestBefore(spider *spider.Spider) {
	spider.Client().Timeout = time.Second * 2
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
