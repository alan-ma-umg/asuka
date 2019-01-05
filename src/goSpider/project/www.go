package project

import (
	"github.com/willf/bloom"
	"goSpider/helper"
	"goSpider/spider"
	"net/url"
	"strings"
	"time"
	"net/http"
)

var tldFilter = bloom.NewWithEstimates(10000000, 0.001)

type Www struct {
}

func (my *Www) EntryUrl() []string {
	return []string{
		//"https://www.douban.com/",
		//"https://www.zhihu.com/explore",
		"https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=url&rsv_pq=91d6b9ef0003df6b&rsv_t=915bi46CZilwcCL7mzlzJoWjIX4rS87mPrBstd9AgYgORE4stRCZxzsFTjA&rqlang=cn&rsv_enter=1&rsv_sug3=4&rsv_sug1=4&rsv_sug7=101&rsv_sug2=0&inputT=373&rsv_sug4=1556",
	}
}

func (my *Www) RequestBefore(spider *spider.Spider) {
	if spider.CurrentRequest != nil {
		spider.CurrentRequest.Header.Set("Accept", "text/html")
	}

	spider.Client.Timeout = 5 * time.Second
}

func (my *Www) ResponseAfter(spider *spider.Spider) {
	if len(spider.RequestsMap) > 10 {
		spider.RequestsMap = map[string]*http.Request{}
	}
}

// queue
func (my *Www) EnqueueFilter(spider *spider.Spider, l *url.URL) bool {

	tld, err := helper.TldDomain(l.String())
	if err != nil {
		return false
	}

	if strings.Contains(strings.ToLower(tld), "gov") {
		return false
	}

	if tldFilter.TestAndAddString(tld) {
		return false
	}

	return true
}

// frequency
func (my *Www) Throttle(spider *spider.Spider) {
	//time.Sleep(2e9)
}
