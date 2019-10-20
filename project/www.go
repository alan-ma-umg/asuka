package project

import (
	"bytes"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"github.com/willf/bloom"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

var tldFilter *bloom.BloomFilter
var tldFilterMutex = &sync.Mutex{}

type AsukaWww struct {
	Id        int                    `xorm:"pk autoincr"`
	Url       string                 `xorm:"varchar(1024)"`
	Data      map[string]interface{} `xorm:"json"`
	Version   int                    `xorm:"version"`
	UpdatedAt int                    `xorm:"updated"`
	CreatedAt int                    `xorm:"created"`
}

type Www struct {
	*Implement
	queueUrlLen int64
}

func (my *Www) Init(d *Dispatcher) {
	tldFilter = bloom.NewWithEstimates(5000000, 0.01)
	err := database.Mysql().CreateTables(&AsukaWww{})
	if err != nil {
		panic(err)
	}

	go func() {
		t := time.NewTicker(time.Second * 5)
		for {
			<-t.C
			my.queueUrlLen, _ = database.Redis().LLen(strings.Split(reflect.TypeOf(my).String(), ".")[1] + "_" + helper.Env().Redis.URLQueueKey).Result()
		}
	}()
}

func (my *Www) EntryUrl() []string {
	return []string{
		"https://www.douban.com/",
		"https://www.zhihu.com/explore",
		"https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=url&rsv_pq=91d6b9ef0003df6b&rsv_t=915bi46CZilwcCL7mzlzJoWjIX4rS87mPrBstd9AgYgORE4stRCZxzsFTjA&rqlang=cn&rsv_enter=1&rsv_sug3=4&rsv_sug1=4&rsv_sug7=101&rsv_sug2=0&inputT=373&rsv_sug4=1556",
	}
}

// frequency
func (my *Www) Throttle(spider *spider.Spider) {
	spider.RecentSeveralTimesResultCap = 20
	if spider.Transport.LoadRate(5) > 5.0 {
		spider.AddSleep(60e9)
	}
}

func (my *Www) RequestBefore(spider *spider.Spider) {
	if spider.CurrentRequest() != nil {
		spider.CurrentRequest().Header.Set("Accept", "text/html")
	}

	spider.Client().Timeout = 4 * time.Second
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *Www) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}
	if strings.ToLower(response.Header.Get("Content-Encoding")) != "gzip" {
		return false, nil
	}
	return true, nil
}

func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil {
			return n.FirstChild.Data
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
// 一般用于采集数据的地方
func (my *Www) ResponseSuccess(spider *spider.Spider) {
	title := ""
	node, err := html.Parse(ioutil.NopCloser(bytes.NewBuffer(spider.ResponseByte)))
	if err != nil {
		return
	}
	title = strings.TrimSpace(pageTitle(node))
	if title == "" {
		return
	}

	model := &AsukaWww{
		Url: spider.CurrentRequest().URL.String(),
		Data: map[string]interface{}{
			"title":  title,
			"server": spider.Transport.S.Host,
			"time":   time.Since(spider.RequestStartTime).String(),
		},
	}
	_, err = database.Mysql().Insert(model)
	if err != nil {
		database.MysqlDelayInsertTillSuccess(model)
		log.Println(spider.CurrentRequest().URL.String(), err)
	}
}

// queue
func (my *Www) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if my.queueUrlLen > 20000 {
		return
	}

	tld, err := helper.TldDomain(l)
	if err != nil {
		return
	}

	if strings.Contains(strings.ToLower(tld), "com.cn") {
		return
	}

	if strings.Contains(strings.ToLower(tld), "gov") {
		return
	}

	tldFilterMutex.Lock()
	defer tldFilterMutex.Unlock()
	if tldFilter.TestAndAddString(tld) {
		return
	}

	return l.String()
}

func (my *Www) ResponseAfter(spider *spider.Spider) {
	spider.ResetRequest()
	spider.Transport.Close()

	my.Implement.ResponseAfter(spider)
}
