package project

import (
	"encoding/json"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type V2ex struct {
	*Implement
	queueUrlLen int64
	showStr     string
	speedMin    time.Duration
	speedTotal  time.Duration
	speedMax    time.Duration
	itemRegex   *regexp.Regexp
}

func (my *V2ex) Init(d *Dispatcher) {

	database.Redis().Del(my.Name() + "_" + helper.Env().Redis.URLQueueKey)

	go func() {
		for {
			time.Sleep(10e9)
			my.queueUrlLen, _ = database.Redis().LLen(my.Name() + "_" + helper.Env().Redis.URLQueueKey).Result()
		}
	}()

	my.speedMin = time.Hour
	my.showStr = "Waiting"
}

func (my *V2ex) resultRedisKeyName(cate string) string {
	return my.Name() + "_result" + "_" + cate
}

func (my *V2ex) Showing() string {
	return my.showStr
}

func (my *V2ex) Name() string {
	return "Renge"
}

func (my *V2ex) EntryUrl() []string {
	var links []string

	for i := 0; i < 1000; i++ {
		links = append(links, "https://v2ex.com/?tab=hot")
		links = append(links, "https://v2ex.com/?tab=all")
	}

	return links
}

func (my *V2ex) Throttle(spider *spider.Spider) {
	spider.AddSleep(time.Duration(rand.Float64() * 600e9))
}

func (my *V2ex) RequestBefore(spider *spider.Spider) {
	spider.Client().Timeout = time.Minute
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *V2ex) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}

	return true, nil
}

// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
// 一般用于采集数据的地方
func (my *V2ex) ResponseSuccess(spider *spider.Spider) {
	if len(spider.ResponseByte) < 1000 {
		return
	}

	res := regexp.MustCompile("<a\\shref=\"([^\"]+)\"\\sclass=\"topic-link\">(.+)</a>").FindAllStringSubmatch(string(spider.ResponseByte), -1)

	var list []string

	for _, v := range res {
		if len(v) < 3 {
			return
		}

		if str, err := json.Marshal(map[string]string{
			"url":   v[1],
			"title": v[2],
		}); err == nil {
			list = append(list, string(str))
		}
	}

	if len(list) == 0 {
		return
	}

	cate := "all"
	if strings.Contains(strings.ToLower(spider.CurrentRequest().URL.String()), "hot") {
		cate = "hot"
	}
	database.Redis().Del(my.resultRedisKeyName(cate))
	if _, err := database.Redis().RPush(my.resultRedisKeyName(cate), list).Result(); err != nil {
		helper.SendTextToWXDoOnceDurationHour(my.Name() + " v2ex redis HMSet fail: " + err.Error())
		log.Println(err)
	}
}

func (my *V2ex) ResponseAfter(spider *spider.Spider) {
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

	my.Implement.ResponseAfter(spider)
}

// queue
func (my *V2ex) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {

	if my.queueUrlLen < 1000 {
		for _, l := range my.EntryUrl() {
			spider.GetQueue().Enqueue(l)
		}
	}

	return
}

func (my *V2ex) WEBSite(w http.ResponseWriter, r *http.Request) {

	var titlesHot []string
	var urlsHot []string

	if m, err := database.Redis().LRange(my.resultRedisKeyName("hot"), 0, -1).Result(); err == nil {
		for _, v := range m {
			l := make(map[string]string)
			if err := json.Unmarshal([]byte(v), &l); err == nil {
				titlesHot = append(titlesHot, l["title"])
				urlsHot = append(urlsHot, l["url"])
			}
		}
	}

	var titlesAll []string
	var urlsAll []string

	if m, err := database.Redis().LRange(my.resultRedisKeyName("all"), 0, -1).Result(); err == nil {
		for _, v := range m {
			l := make(map[string]string)
			if err := json.Unmarshal([]byte(v), &l); err == nil {
				titlesAll = append(titlesAll, l["title"])
				urlsAll = append(urlsAll, l["url"])
			}
		}
	}

	data := struct {
		ListTitleHot []string
		ListUrlHot   []string
		ListTitleAll []string
		ListUrlAll   []string
		ProjectName  string
	}{
		ListTitleHot: titlesHot,
		ListUrlHot:   urlsHot,
		ListTitleAll: titlesAll,
		ListUrlAll:   urlsAll,
		ProjectName:  my.Name(),
	}

	helper.GetTemplates().ExecuteTemplate(w, "v2ex.html", data)
}
