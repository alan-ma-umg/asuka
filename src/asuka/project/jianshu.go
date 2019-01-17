package project

import (
	"asuka/database"
	"asuka/helper"
	"asuka/spider"
	"bytes"
	"golang.org/x/net/html"
	"hash/crc32"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type AsukaJianShu struct {
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

type JianShu struct {
	lastRequestUrl string
}

func (my *JianShu) EntryUrl() []string {
	err := database.Mysql().CreateTables(&AsukaJianShu{})
	if err != nil {
		panic(err)
	}
	return []string{
		"https://www.jianshu.com/",
		"https://www.jianshu.com/",
		"https://www.jianshu.com/",
		"https://www.jianshu.com/",
		"https://www.jianshu.com/",
		"https://www.jianshu.com/",
	}
}

// frequency
func (my *JianShu) Throttle(spider *spider.Spider) {
	if spider.Transport.LoadRate(5) > 5.0 {
		spider.AddSleep(60e9)
	}

	spider.AddSleep(time.Duration(rand.Float64() * 60e9))
}

func (my *JianShu) RequestBefore(spider *spider.Spider) {
	//accept
	if spider.CurrentRequest != nil {
		spider.CurrentRequest.Header.Set("Accept", "text/html")
	}

	//Referer
	if spider.CurrentRequest != nil && spider.CurrentRequest.Referer() == "" && my.lastRequestUrl != "" {
		spider.CurrentRequest.Header.Set("Referer", my.lastRequestUrl)
	}

	spider.Client.Timeout = 10 * time.Second
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *JianShu) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}
	if strings.ToLower(response.Header.Get("Content-Encoding")) != "gzip" {
		return false, nil
	}
	return true, nil
}

func JianShuPageHtml(n *html.Node, title, watch, view *string, tag *[]string) {
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
func (my *JianShu) ResponseSuccess(spider *spider.Spider) {
	my.lastRequestUrl = spider.CurrentRequest.URL.String()
	node, err := html.Parse(ioutil.NopCloser(bytes.NewBuffer(spider.ResponseByte)))
	if err != nil {
		return
	}

	var title, watch, view string
	var tag []string

	JianShuPageHtml(node, &title, &watch, &view, &tag)
	if title == "" {
		return
	}

	//2019/01/14 16:41:03 zhihu.go:171: https://www.jianshu.com/nb/31338671 Error 1366: Incorrect string value: '\xF0\x9F\x92\x8E&\xF0...' for column 'title' at row 1
	_, err = database.Mysql().Insert(&AsukaJianShu{
		Url:      spider.CurrentRequest.URL.String(),
		Referer:  spider.CurrentRequest.Referer(),                                                  //todo only test
		Cookie:   helper.TruncateStr([]rune(spider.CurrentRequest.Header.Get("cookie")), 2000, ""), //todo only test
		UrlCrc32: int64(crc32.ChecksumIEEE([]byte(spider.CurrentRequest.URL.String()))),
		Title:    title,
		Tag:      tag,
		Data: map[string]interface{}{
			"server": spider.Transport.S.ServerAddr,
			"time":   time.Since(spider.RequestStartTime).String(),
			"watch":  watch,
			"view":   view,
		},
	})
	if err != nil {
		log.Println(spider.CurrentRequest.URL.String(), err)
	}
}

// queue
func (my *JianShu) EnqueueFilter(spider *spider.Spider, l *url.URL) bool {

	//tld, err := helper.TldDomain(l)
	//if err != nil {
	//	return false
	//}

	if !strings.HasPrefix(strings.ToLower(l.String()), "https://www.jianshu.com/p/") {
		return false
	}

	return true
}

func (my *JianShu) ResponseAfter(spider *spider.Spider) {
	spider.Transport.T.(*http.Transport).DisableKeepAlives = false
	if spider.FailureLevel > 10 {
		jianShuResetSpider(spider)
	} else if rand.Intn(30) == 10 {
		jianShuResetSpider(spider)
	}
}

func jianShuResetSpider(spider *spider.Spider) {
	spider.Transport.Reconnection()
	jar, _ := cookiejar.New(nil)
	spider.Client = &http.Client{Transport: spider.Transport.T, Jar: jar, Timeout: time.Second * 30}
	spider.RequestsMap = map[string]*http.Request{}
}
