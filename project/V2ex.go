package project

import (
	"encoding/json"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type V2ex struct {
	*Implement
	*SpeedShowing
	itemRegex      *regexp.Regexp
	lastUpdateTime time.Time
}

func (my *V2ex) InitBloomFilterCapacity() uint { return 1000000 }
func (my *V2ex) Init(d *Dispatcher) {
	database.Redis().Del(my.Name() + "_" + helper.Env().Redis.URLQueueKey)
	go func() {
		for {
			time.Sleep(30e9)
			if database.Redis().LLen(my.Name()+"_"+helper.Env().Redis.URLQueueKey).Val() < 1000 {
				d.GetQueue().Enqueue(my.EntryUrl())
			}
		}
	}()
}

func (my *V2ex) resultRedisKeyName(cate string) string {
	return my.Name() + "_result" + "_" + cate
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
	spider.SetRequestTimeout(time.Minute)
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

	my.lastUpdateTime = time.Now()
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
	my.SpeedShowing.ResponseSuccess(spider)
	my.Implement.ResponseAfter(spider)
}

// queue
func (my *V2ex) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	return
}

func (my *V2ex) WEBSiteLoginRequired(w http.ResponseWriter, r *http.Request) bool {
	return false
}

func (my *V2ex) WEBSite(w http.ResponseWriter, r *http.Request) {
	if tab := r.URL.Query().Get("tab"); tab != "" {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		w.Header().Set("Content-type", "text/plain;charset=utf-8")
		if !my.lastUpdateTime.IsZero() {
			io.WriteString(w, my.lastUpdateTime.Format(time.Stamp)+"\n")
		}
		if m, err := database.Redis().LRange(my.resultRedisKeyName(strings.TrimSpace(strings.ToLower(tab))), 0, -1).Result(); err == nil {
			for i, v := range m {
				if limit > 0 && i >= limit {
					break
				}

				l := make(map[string]string)
				if err := json.Unmarshal([]byte(v), &l); err == nil {
					u := strings.Split(l["url"], "#")
					//https://www.v2ex.com/t/613922#reply19
					replay := ""
					if len(u) > 1 {
						replay = strings.TrimLeft(u[1], "replay")
					}
					io.WriteString(w, "\\"+strconv.Itoa(i+1)+" "+l["title"]+"["+replay+"]"+"\nhttps://v2ex.com"+u[0]+"\n")
				}
			}
		}
		return
	}

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
