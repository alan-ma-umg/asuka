package project

import (
	"bytes"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/spider"
	"golang.org/x/net/html"
	"hash/crc32"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AsukaZhiHu struct {
	Id        int64  `xorm:"pk autoincr"`
	Url       string `xorm:"varchar(1024)"`
	Referer   string `xorm:"varchar(1024)"` //todo for test
	UrlCrc32  int64
	Title     string                 `xorm:"varchar(1024)"`
	Tag       []string               `xorm:"json"`
	Data      map[string]interface{} `xorm:"json"`
	Version   int                    `xorm:"version"`
	UpdatedAt int                    `xorm:"updated"`
	CreatedAt int                    `xorm:"created"`
}

func init() {
	err := database.Mysql().CreateTables(&AsukaZhiHu{})
	if err != nil {
		panic(err)
	}
}

type ZhiHu struct {
	lastRequestUrl string
}

func (my *ZhiHu) EntryUrl() []string {
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
func (my *ZhiHu) Throttle(spider *spider.Spider) {
	if spider.Transport.LoadRate(5) > 5.0 {
		spider.AddSleep(60e9)
	}

	spider.AddSleep(time.Duration(rand.Float64() * 50e9))
}

func (my *ZhiHu) RequestBefore(spider *spider.Spider) {
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
	my.lastRequestUrl = spider.CurrentRequest.URL.String()
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
	_, err = database.Mysql().Insert(&AsukaZhiHu{
		Url:      spider.CurrentRequest.URL.String(),
		Referer:  spider.CurrentRequest.Referer(), //todo only test
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
func (my *ZhiHu) EnqueueFilter(spider *spider.Spider, l *url.URL) bool {

	tld, err := helper.TldDomain(l)
	if err != nil {
		return false
	}

	if !strings.Contains(strings.ToLower(tld), "jianshu.com") {
		return false
	}

	return true
}

func (my *ZhiHu) ResponseAfter(spider *spider.Spider) {
	//free the memory
	//if len(spider.RequestsMap) > 10 {
	//	spider.Client.Jar, _ = cookiejar.New(nil)
	//	spider.RequestsMap = map[string]*http.Request{}
	//}
}
