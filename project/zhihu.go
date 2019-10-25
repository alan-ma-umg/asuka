package project

import (
	"bytes"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"golang.org/x/net/html"
	"hash/crc32"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type AsukaZhiHu struct {
	Id        int64  `xorm:"pk autoincr"`
	Url       string `xorm:"varchar(1024)"`
	UrlCrc32  int64
	Title     string                 `xorm:"varchar(1024)"`
	Tag       []string               `xorm:"json"`
	Data      map[string]interface{} `xorm:"json"`
	Referer   string                 `xorm:"varchar(1024)"` //todo for test
	Cookie    string                 `xorm:"varchar(2048)"` //todo for test
	Version   int                    `xorm:"version"`
	UpdatedAt int                    `xorm:"updated"`
	CreatedAt int                    `xorm:"created"`
}

type ZhiHu struct {
	*Implement
	lastRequestUrl  string
	queueUrlLen     int64
	insertSpeed     int
	lastInsertId    int64
	lastInsertError string
}

func (my *ZhiHu) Name() (str string) {
	return "Miyamori"
}

func (my *ZhiHu) Showing() (str string) {
	str = "ID: " + strconv.Itoa(int(my.lastInsertId)) + " : " + strconv.Itoa(my.insertSpeed) + "/s"
	if len(database.MysqlDelayInsertQueue) > 0 {
		str += " delay: " + strconv.Itoa(len(database.MysqlDelayInsertQueue))
	}
	if my.lastInsertError != "" {
		str += " Error: " + my.lastInsertError
	}
	return
}
func (my *ZhiHu) Init(d *Dispatcher) {
	err := database.Mysql().CreateTables(&AsukaZhiHu{})
	if err != nil {
		panic(err)
	}

	go func() {
		s := time.NewTicker(time.Second)
		insertIdPoint := my.lastInsertId
		for {
			<-s.C
			my.insertSpeed = int(my.lastInsertId - insertIdPoint)
			insertIdPoint = my.lastInsertId
		}
	}()

	go func() {
		t := time.NewTicker(time.Second * 5)
		for {
			<-t.C
			my.queueUrlLen, _ = database.Redis().LLen(strings.Split(reflect.TypeOf(my).String(), ".")[1] + "_" + helper.Env().Redis.URLQueueKey).Result()
		}
	}()
}
func (my *ZhiHu) EntryUrl() []string {
	return []string{
		"https://www.zhihu.com/explore",
		"https://www.zhihu.com/explore",
		"https://www.zhihu.com/explore",
		"https://www.zhihu.com/explore",
		"https://www.zhihu.com/explore",
		"https://www.zhihu.com/explore",
		"https://www.zhihu.com/explore",
	}
}

// frequency
func (my *ZhiHu) Throttle(spider *spider.Spider) {
	if spider.LoadRate(5) > 5.0 {
		spider.AddSleep(120e9)
	}

	spider.AddSleep(time.Duration(rand.Float64() * 100e9))

	if spider.FailureLevel > 1 {
		zhiHuResetSpider(spider)
	} else if rand.Intn(30) == 10 {
		zhiHuResetSpider(spider)
	}
}

func (my *ZhiHu) RequestBefore(spider *spider.Spider) {
	//accept
	if spider.CurrentRequest() != nil {
		spider.CurrentRequest().Header.Set("Accept", "text/html")
	}

	//Referer
	if spider.CurrentRequest() != nil && spider.CurrentRequest().Referer() == "" && my.lastRequestUrl != "" {
		spider.CurrentRequest().Header.Set("Referer", my.lastRequestUrl)
	}

	spider.Client().Timeout = 10 * time.Second
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *ZhiHu) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}
	if strings.ToLower(response.Header.Get("Content-Encoding")) != "gzip" {
		return false, nil
	}
	return true, nil
}

func PageHtml(n *html.Node, title, watch, view *string, tag *[]string) {
	if n.Type == html.ElementNode {

		//title
		if *title == "" {
			if n.Data == "title" {
				if n.FirstChild != nil {
					*title = n.FirstChild.Data
				}
			}
		}

		//watch && view
		if n.Data == "strong" {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					//watch
					if *watch == "" && len(n.Attr) == 2 {
						*watch = n.Attr[1].Val

						//view
					} else if *view == "" && len(n.Attr) == 2 {
						*view = n.Attr[1].Val
					}
				}
			}
		}

		//tag
		if n.Data == "span" {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					if attr.Val == "Tag-content" {
						newTag := ""
						for tagN := n.FirstChild; tagN != nil; tagN = tagN.FirstChild {
							newTag = tagN.Data
						}
						*tag = append(*tag, strings.TrimSpace(newTag))
					}
				}
			}

		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		PageHtml(c, title, watch, view, tag)
	}

	return
}

// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
// 一般用于采集数据的地方
func (my *ZhiHu) ResponseSuccess(spider *spider.Spider) {
	my.lastRequestUrl = spider.CurrentRequest().URL.String()
	node, err := html.Parse(ioutil.NopCloser(bytes.NewBuffer(spider.ResponseByte)))
	if err != nil {
		return
	}

	var title, watch, view string
	var tag []string

	PageHtml(node, &title, &watch, &view, &tag)
	if title == "" {
		return
	}

	model := &AsukaZhiHu{
		Url:      spider.CurrentRequest().URL.String(),
		Referer:  spider.CurrentRequest().Referer(),                                                  //todo only test
		Cookie:   helper.TruncateStr([]rune(spider.CurrentRequest().Header.Get("cookie")), 2000, ""), //todo only test
		UrlCrc32: int64(crc32.ChecksumIEEE([]byte(spider.CurrentRequest().URL.String()))),
		Title:    title,
		Tag:      tag,
		Data: map[string]interface{}{
			"server": spider.TransportUrl.Host,
			"time":   time.Since(spider.RequestStartTime).String(),
			"watch":  watch,
			"view":   view,
		},
	}
	_, err = database.Mysql().Insert(model)
	my.lastInsertId = model.Id
	if err != nil {
		my.lastInsertError = time.Now().Format(time.RFC3339) + ":" + err.Error()
		database.MysqlDelayInsertTillSuccess(model)
		log.Println(spider.CurrentRequest().URL.String(), err)
	}
}

// queue
func (my *ZhiHu) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if my.queueUrlLen > 20000 {
		return
	}

	if !strings.HasPrefix(strings.ToLower(l.String()), "https://www.zhihu.com/people") && !strings.HasPrefix(strings.ToLower(l.String()), "https://www.zhihu.com/question") && !strings.HasPrefix(strings.ToLower(l.String()), "https://www.zhihu.com/collection") {
		return
	}

	return l.Scheme + "://" + l.Host + l.Path
}

func zhiHuResetSpider(spider *spider.Spider) {
	spider.ResetSpider()
}
