package project

import (
	"asuka/database"
	"asuka/spider"
	"bytes"
	"encoding/json"
	"golang.org/x/net/html"
	"hash/crc32"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type AsukaDouBan struct {
	Id        int64 `xorm:"pk autoincr"`
	DouBanID  int64
	Title     string   `xorm:"varchar(1024)"`
	Name      string   `xorm:"varchar(1024)"`
	Alias     []string `xorm:"json"`
	Date      int64
	DateStr   string
	Rating    float64
	Votes     int64
	Img       string   `xorm:"varchar(1024)"`
	Area      []string `xorm:"json"`
	Cate      string //图书/电影
	Genre     []string `xorm:"json"`
	Summary   string   `xorm:"varchar(10240)"`
	Author    []string `xorm:"json"`
	Director  []string `xorm:"json"`
	Actor     []string `xorm:"json"`
	IMDB      string   `xorm:"varchar(1024)"`
	ISBN      int64
	Data      map[string]interface{} `xorm:"json"`
	Url       string                 `xorm:"varchar(1024)"`
	UrlCrc32  int64
	Version   int `xorm:"version"`
	UpdatedAt int `xorm:"updated"`
	CreatedAt int `xorm:"created"`
}

type DouBan struct {
	lastRequestUrl string
}

func (my *DouBan) EntryUrl() []string {
	err := database.Mysql().CreateTables(&AsukaDouBan{})
	if err != nil {
		panic(err)
	}
	return []string{
		"https://book.douban.com/tag/",
		"https://book.douban.com/tag/",
		"https://book.douban.com/tag/",
		"https://book.douban.com/tag/",
		"https://book.douban.com/tag/",
		"https://movie.douban.com/tag/",
		"https://movie.douban.com/tag/",
		"https://movie.douban.com/tag/",
		"https://movie.douban.com/tag/",
		"https://movie.douban.com/tag/",
	}
}

// frequency
func (my *DouBan) Throttle(spider *spider.Spider) {
	if spider.Transport.LoadRate(5) > 5.0 {
		spider.AddSleep(120e9)
	}

	spider.AddSleep(time.Duration(rand.Float64() * 20e9))

	if spider.FailureLevel > 1 {
		DouBanResetSpider(spider)
	} else if rand.Intn(30) == 10 {
		DouBanResetSpider(spider)
	}
}

func (my *DouBan) RequestBefore(spider *spider.Spider) {
	//accept
	if spider.CurrentRequest != nil {
		spider.CurrentRequest.Header.Set("Accept", "text/html")
	}

	//Referer
	if spider.CurrentRequest != nil && spider.CurrentRequest.Referer() == "" && my.lastRequestUrl != "" {
		spider.CurrentRequest.Header.Set("Referer", my.lastRequestUrl)
	}

	spider.Client.Timeout = 15 * time.Second
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *DouBan) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}
	if strings.ToLower(response.Header.Get("Content-Encoding")) != "gzip" {
		return false, nil
	}
	return true, nil
}

func stateInString(c byte) bool {
	if c == '"' {
		return true
	}
	if c == '\\' {
		return true
	}
	if c < 0x20 {
		return false
	}
	return true
}

func DouBanPageHtml(n *html.Node, model *AsukaDouBan) {
	if n.Type == html.ElementNode {

		//title
		if model.Title == "" && n.Data == "title" && n.FirstChild != nil {
			model.Title = strings.TrimSpace(n.FirstChild.Data)
		}

		//data
		if model.Data == nil && n.Data == "script" && n.FirstChild != nil {
			for _, attr := range n.Attr {
				if attr.Val == "application/ld+json" {

					jsonStr := []byte(n.FirstChild.Data)

					for i, ch := range jsonStr {
						if !stateInString(ch) {
							jsonStr[i] = ' '
						}
					}

					if err := json.Unmarshal(jsonStr, &model.Data); err != nil {
						log.Println(err, model.Url, model.Title, n.FirstChild.Data)
					} else {
						if v, ok := model.Data["author"]; ok {
							for _, v := range v.([]interface{}) {
								model.Author = append(model.Author, strings.TrimSpace(v.(map[string]interface{})["name"].(string)))
							}
						}

						if v, ok := model.Data["director"]; ok {
							for _, v := range v.([]interface{}) {
								model.Director = append(model.Director, strings.TrimSpace(v.(map[string]interface{})["name"].(string)))
							}
						}

						if v, ok := model.Data["actor"]; ok {
							for _, v := range v.([]interface{}) {
								model.Actor = append(model.Actor, strings.TrimSpace(v.(map[string]interface{})["name"].(string)))
							}
						}

						if v, ok := model.Data["genre"]; ok {
							for _, v := range v.([]interface{}) {
								model.Genre = append(model.Genre, strings.TrimSpace(v.(string)))
							}
						}
						if v, ok := model.Data["name"]; ok {
							model.Name = strings.TrimSpace(v.(string))
						}
						if v, ok := model.Data["image"]; ok {
							model.Img = strings.TrimSpace(v.(string))
						}

						if v, ok := model.Data["description"]; ok {
							model.Summary = strings.TrimSpace(v.(string))
						}

						if rating, ok := model.Data["aggregateRating"]; ok {
							if votes, err := strconv.Atoi(rating.(map[string]interface{})["ratingCount"].(string)); err == nil {
								model.Votes = int64(votes)
							}
							if f, err := strconv.ParseFloat(rating.(map[string]interface{})["ratingValue"].(string), 64); err == nil {
								model.Rating = f
							}
						}

						if v, ok := model.Data["datePublished"]; ok {
							model.DateStr = strings.TrimSpace(v.(string))
							if t, err := time.Parse("2006-1-2", model.DateStr); err == nil {
								model.Date = t.Unix()
							} else if t, err := time.Parse("2006-1", model.DateStr); err == nil {
								model.Date = t.Unix()
							} else if t, err := time.Parse("2006", model.DateStr); err == nil {
								model.Date = t.Unix()
							}
						}

						if v, ok := model.Data["isbn"]; ok {
							if i, err := strconv.ParseInt(v.(string), 0, 64); err == nil {
								model.ISBN = int64(i)
							}
						}
					}
				}
			}
		}

		//IMDb
		if model.IMDB == "" && n.Data == "a" && n.FirstChild != nil {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.Contains(attr.Val, "www.imdb.com") {
					model.IMDB = strings.TrimSpace(n.FirstChild.Data)
				}
			}
		}

		//img
		if model.Img == "" && n.Data == "img" && n.FirstChild == nil {
			for _, attr := range n.Attr {
				if attr.Val == "v:photo" {
					for _, attr := range n.Attr {
						if attr.Key == "src" {
							model.Img = attr.Val
						}
					}
				}
			}
		}

		//alias
		if len(model.Alias) == 0 && n.Data == "span" && n.FirstChild != nil && n.FirstChild.Data == "又名:" {
			if alias := strings.Split(n.NextSibling.Data, "/"); len(alias) > 0 {
				for _, v := range alias {
					model.Alias = append(model.Alias, strings.TrimSpace(v))
				}
			}
		}

		//alias
		if model.DateStr == "" && n.Data == "span" && n.FirstChild != nil && n.FirstChild.Data == "出版年:" && n.NextSibling != nil && n.NextSibling.Data != "" {
			model.DateStr = strings.TrimSpace(n.NextSibling.Data)
			if t, err := time.Parse("2006-1-2", model.DateStr); err == nil {
				model.Date = t.Unix()
			} else if t, err := time.Parse("2006-1", model.DateStr); err == nil {
				model.Date = t.Unix()
			} else if t, err := time.Parse("2006", model.DateStr); err == nil {
				model.Date = t.Unix()
			}
		}

		//votes
		if model.Votes == 0 && n.Data == "span" && n.FirstChild != nil {
			for _, attr := range n.Attr {
				if attr.Val == "v:votes" {
					if f, err := strconv.Atoi(strings.TrimSpace(n.FirstChild.Data)); err == nil {
						model.Votes = int64(f)
					}
				}
			}
		}

		//rating
		if model.Rating == 0 && n.Data == "strong" && n.FirstChild != nil {
			for _, attr := range n.Attr {
				if attr.Val == "v:average" {
					if f, err := strconv.ParseFloat(strings.TrimSpace(n.FirstChild.Data), 64); err == nil {
						model.Rating = f
					}
				}
			}
		}

		//area
		if len(model.Area) == 0 && n.Data == "span" && n.FirstChild != nil && n.FirstChild.Data == "制片国家/地区:" {
			if area := strings.Split(n.NextSibling.Data, "/"); len(area) > 0 {
				for _, v := range area {
					model.Area = append(model.Area, strings.TrimSpace(v))
				}
			}
		}

		//Summary
		if model.Summary == "" && n.Data == "div" && n.FirstChild != nil {
			for _, attr := range n.Attr {
				if attr.Val == "intro" {
					if strings.TrimSpace(n.FirstChild.Data) != "" {
						model.Summary += strings.TrimSpace(n.FirstChild.Data) + "</br>"
					}
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.FirstChild != nil {
							model.Summary += c.FirstChild.Data + "</br>"
						}
					}
				}
			}
		}

		//if n.Data == "strong" {
		//	for _, attr := range n.Attr {
		//		if attr.Key == "class" {
		//			//watch
		//			if *watch == "" && len(n.Attr) == 2 {
		//				*watch = n.Attr[1].Val
		//
		//				//view
		//			} else if *view == "" && len(n.Attr) == 2 {
		//				*view = n.Attr[1].Val
		//			}
		//		}
		//	}
		//}

		////watch && view
		//if n.Data == "strong" {
		//	for _, attr := range n.Attr {
		//		if attr.Key == "class" {
		//			//watch
		//			if *watch == "" && len(n.Attr) == 2 {
		//				*watch = n.Attr[1].Val
		//
		//				//view
		//			} else if *view == "" && len(n.Attr) == 2 {
		//				*view = n.Attr[1].Val
		//			}
		//		}
		//	}
		//}
		//
		////tag
		//if n.Data == "span" {
		//	for _, attr := range n.Attr {
		//		if attr.Key == "class" {
		//			if attr.Val == "Tag-content" {
		//				newTag := ""
		//				for tagN := n.FirstChild; tagN != nil; tagN = tagN.FirstChild {
		//					newTag = tagN.Data
		//				}
		//				*tag = append(*tag, strings.TrimSpace(newTag))
		//			}
		//		}
		//	}
		//
		//}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		DouBanPageHtml(c, model)
	}
	return
}

// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
// 一般用于采集数据的地方
func (my *DouBan) ResponseSuccess(spider *spider.Spider) {
	my.lastRequestUrl = spider.CurrentRequest.URL.String()
	node, err := html.Parse(ioutil.NopCloser(bytes.NewBuffer(spider.ResponseByte)))
	if err != nil {
		return
	}

	model := &AsukaDouBan{
		//DouBanID: int64(douBanId),
		Url:      spider.CurrentRequest.URL.String(),
		UrlCrc32: int64(crc32.ChecksumIEEE([]byte(spider.CurrentRequest.URL.String()))),
	}

	if paths := strings.Split(spider.CurrentRequest.URL.Path, "/"); len(paths) > 2 {
		model.DouBanID, _ = strconv.ParseInt(paths[2], 0, 64)
	}

	if strings.HasPrefix(spider.CurrentRequest.URL.String(), "https://movie.douban.com") {
		model.Cate = "电影"
	}
	if strings.HasPrefix(spider.CurrentRequest.URL.String(), "https://book.douban.com") {
		model.Cate = "图书"
	}

	DouBanPageHtml(node, model)
	if model.Title == "" {
		return
	}
	_, err = database.Mysql().Insert(model)
	if err != nil {
		log.Println(spider.CurrentRequest.URL.String(), err)
	}
}

// queue
func (my *DouBan) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if !strings.HasPrefix(strings.ToLower(l.String()), "https://movie.douban.com/subject") && !strings.HasPrefix(strings.ToLower(l.String()), "https://book.douban.com/subject") && !strings.HasPrefix(strings.ToLower(l.String()), "https://book.douban.com/tag") && !strings.HasPrefix(strings.ToLower(l.String()), "https://movie.douban.com/tag") {
		return
	}

	if strings.HasSuffix(strings.ToLower(l.String()), "buylinks") {
		return
	}

	return l.Scheme + "://" + l.Host + l.Path
}

func (my *DouBan) ResponseAfter(spider *spider.Spider) {

}

func DouBanResetSpider(spider *spider.Spider) {
	spider.RequestsMap = map[string]*http.Request{}
	spider.UpdateTransport()
}
