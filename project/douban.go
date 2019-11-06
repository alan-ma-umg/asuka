package project

import (
	"bytes"
	"encoding/json"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/spider"
	"golang.org/x/net/html"
	"hash/crc32"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AsukaDouBan struct {
	Id        int64 `xorm:"pk autoincr"`
	DouBanId  int64
	Title     string   `xorm:"varchar(1024)"`
	Name      string   `xorm:"varchar(1024)"`
	Alias     []string `xorm:"json"`
	Date      int64
	DateStr   string
	Rating    float64
	Votes     int64
	Img       string   `xorm:"varchar(1024)"`
	Area      []string `xorm:"json"`
	Cate      string
	Genre     []string `xorm:"json"`
	Summary   string   `xorm:"varchar(10240)"`
	Author    []string `xorm:"json"`
	Director  []string `xorm:"json"`
	Actor     []string `xorm:"json"`
	Imdb      string   `xorm:"varchar(1024)"`
	Isbn      int64
	Data      map[string]interface{} `xorm:"json"`
	Url       string                 `xorm:"varchar(1024)"`
	UrlCrc32  int64                  `xorm:"index"`
	Version   int                    `xorm:"version"`
	UpdatedAt int                    `xorm:"updated"`
	CreatedAt int                    `xorm:"created"`
}

var isDouBanSubject = regexp.MustCompile(`douban.com/subject/[0-9]+/?$`).MatchString

type DouBan struct {
	*Implement
	lastRequestUrl  string
	dbSpeed         int
	dbSpeedNum      int
	lastInsertId    int64
	lastInsertError string
}

func (my *DouBan) InitBloomFilterCapacity() uint { return 10000000 }
func (my *DouBan) Name() string {
	return "Hitagi"
}

func (my *DouBan) Showing() (str string) {
	str = "ID: " + strconv.Itoa(int(my.lastInsertId)) + " : " + strconv.Itoa(my.dbSpeed) + "/s"
	if len(database.MysqlDelayInsertQueue) > 0 {
		str += "<span class='text-orange'> delay: " + strconv.Itoa(len(database.MysqlDelayInsertQueue)) + "</span>"
	}
	if my.lastInsertError != "" {
		str += " Error: " + my.lastInsertError
	}
	return
}

func (my *DouBan) Init(d *Dispatcher) {
	//create table
	err := database.Mysql().CreateTables(&AsukaDouBan{})
	if err != nil {
		d.StopTime = time.Now() //stop
		log.Println(err)
		return
	}
	database.Mysql().CreateIndexes(&AsukaDouBan{})

	database.Mysql().Table(&AsukaDouBan{}).Desc("id").Limit(1).Cols("id").Get(&my.lastInsertId)

	go func() {
		s := time.NewTicker(time.Second)
		dbSpeedPoint := my.dbSpeedNum
		for {
			<-s.C
			my.dbSpeed = int(my.dbSpeedNum - dbSpeedPoint)
			dbSpeedPoint = my.dbSpeedNum
		}
	}()
}

func (my *DouBan) EntryUrl() []string {
	var links []string
	for ii := 0; ii < 20; ii++ {
		//for i := 0; i <= 12000; i++ {
		//	if i%1000 == 0 {
		links = append(links, "https://book.douban.com/tag/")
		links = append(links, "https://movie.douban.com/tag/")
		links = append(links, "https://movie.douban.com/tag/2019")
		links = append(links, "https://movie.douban.com/tag/2018")
		links = append(links, "https://movie.douban.com/tag/2017")
		links = append(links, "https://movie.douban.com/tag/2016")
		links = append(links, "https://movie.douban.com/tag/2015")
		//}
		//links = append(links, "https://movie.douban.com/j/new_search_subjects?start="+strconv.Itoa(i))
		//}
	}
	return links
}

// frequency
func (my *DouBan) Throttle(spider *spider.Spider) {
	if spider.LoadRate(5) > 5.0 {
		spider.AddSleep(120e9)
	}

	spider.AddSleep(time.Duration(rand.Float64() * 50e9))
	//
	//if spider.FailureLevel > 1 {
	//	DouBanResetSpider(spider)
	//} else if rand.Intn(30) == 10 {
	//	DouBanResetSpider(spider)
	//}

	if spider.FailureLevel > 40 {
		spider.Delete = true
	}
}

func (my *DouBan) RequestBefore(spider *spider.Spider) {
	//accept
	if spider.CurrentRequest() != nil {
		spider.CurrentRequest().Header.Set("Accept", "text/html")
	}

	//Referer
	if spider.CurrentRequest() != nil && spider.CurrentRequest().Referer() == "" && my.lastRequestUrl != "" {
		spider.CurrentRequest().Header.Set("Referer", my.lastRequestUrl)
	}

	spider.SetRequestTimeout(time.Minute)
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *DouBan) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") && !strings.Contains(response.Header.Get("Content-type"), "application/json") {
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

func doubanParseDate(model *AsukaDouBan) {
	if t, err := time.Parse("2006-1-2", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006-1", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006年01月", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006年1月", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006年1月第一版", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006年01月第一版", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("20060102", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("200601", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006年", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006年1月2日", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006年01月02日", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006年1月第1版", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006.1", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006.01", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006/01/01", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006/01", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006/1/2", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("2006/1", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("06/01/02", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("060102", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if t, err := time.Parse("06/1/2", model.DateStr); err == nil {
		model.Date = t.Unix()
	} else if re := regexp.MustCompile(`[1-2][0-9]{3}`).FindStringSubmatch(model.DateStr); len(re) > 0 {
		if t, err := time.Parse("2006", re[0]); err == nil {
			model.Date = t.Unix()
		} else {
			//log.Println(err, model.Url, model.Title, model.DateStr)
		}
	} else {
		//log.Println(err, model.Url, model.Title, model.DateStr)
	}
}

func DouBanPageHtmlSecondly(n *html.Node, model *AsukaDouBan) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, model.Url, model.Title)
		}
	}()

	if n.Type == html.ElementNode {
		//date
		if model.DateStr == "" && n.Data == "span" && n.FirstChild != nil {
			for _, attr := range n.Attr {
				if attr.Val == "year" {
					model.DateStr = strings.TrimRight(strings.TrimLeft(strings.TrimSpace(n.FirstChild.Data), "("), ")")
					if model.DateStr != "" {
						doubanParseDate(model)
						//
						//if t, err := time.Parse("2006-1-2", model.DateStr); err == nil {
						//	model.Date = t.Unix()
						//} else if t, err := time.Parse("2006-1", model.DateStr); err == nil {
						//	model.Date = t.Unix()
						//} else if t, err := time.Parse("2006", model.DateStr); err == nil {
						//	model.Date = t.Unix()
						//} else {
						//	log.Println(err, model.Url, model.Title, model.DateStr)
						//}
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		DouBanPageHtmlSecondly(c, model)
	}
	return
}

func DouBanJsonUnmarshal(jsonStr []byte, model *AsukaDouBan) (err error) {
	for i, ch := range jsonStr {
		if !stateInString(ch) {
			jsonStr[i] = ' '
		}
	}

	err = json.Unmarshal(jsonStr, &model.Data)

	if err != nil && (strings.Contains(err.Error(), `in string escape code`) || strings.Contains(err.Error(), `in string literal`)) {
		for i, ch := range jsonStr {
			if ch == 92 { // \ = 92
				jsonStr[i] = '/'
			}
		}
		err = json.Unmarshal(jsonStr, &model.Data)
	}

	return
}

func DouBanPageHtml(n *html.Node, model *AsukaDouBan) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, model.Url, model.Title)
		}
	}()

	if n.Type == html.ElementNode {

		//title
		if model.Title == "" && n.Data == "title" && n.FirstChild != nil {
			model.Title = strings.TrimSpace(n.FirstChild.Data)
		}

		//data
		if model.Data == nil && n.Data == "script" && n.FirstChild != nil {
			for _, attr := range n.Attr {
				if attr.Val == "application/ld+json" {
					if err := DouBanJsonUnmarshal([]byte(n.FirstChild.Data), model); err != nil {
						log.Println(err, model.Url, model.Title, n.FirstChild.Data)
					} else {
						if v, ok := model.Data["author"]; ok {
							if v, ok := v.([]interface{}); ok {
								for _, v := range v {
									model.Author = append(model.Author, strings.TrimSpace(v.(map[string]interface{})["name"].(string)))
								}
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}

						if v, ok := model.Data["director"]; ok {
							if v, ok := v.([]interface{}); ok {
								for _, v := range v {
									model.Director = append(model.Director, strings.TrimSpace(v.(map[string]interface{})["name"].(string)))
								}
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}

						if v, ok := model.Data["actor"]; ok {
							if v, ok := v.([]interface{}); ok {
								for _, v := range v {
									model.Actor = append(model.Actor, strings.TrimSpace(v.(map[string]interface{})["name"].(string)))
								}
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}

						if v, ok := model.Data["genre"]; ok {
							if v, ok := v.([]interface{}); ok {
								for _, v := range v {
									model.Genre = append(model.Genre, strings.TrimSpace(v.(string)))
								}
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}
						if v, ok := model.Data["name"]; ok {
							if v, ok := v.(string); ok {
								model.Name = strings.TrimSpace(v)
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}
						if v, ok := model.Data["image"]; ok {
							if v, ok := v.(string); ok {
								model.Img = strings.TrimSpace(v)
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}

						if v, ok := model.Data["description"]; ok {
							if v, ok := v.(string); ok {
								model.Summary = strings.TrimSpace(v)
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}

						if v, ok := model.Data["aggregateRating"]; ok {
							if v, ok := v.(map[string]interface{}); ok && reflect.TypeOf(v["ratingCount"]).Kind() == reflect.String {
								if votes, err := strconv.Atoi(v["ratingCount"].(string)); err == nil {
									model.Votes = int64(votes)
								}
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}

							if v, ok := v.(map[string]interface{}); ok && reflect.TypeOf(v["ratingValue"]).Kind() == reflect.String {
								if f, err := strconv.ParseFloat(v["ratingValue"].(string), 64); err == nil {
									model.Rating = f
								}
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}

						if v, ok := model.Data["datePublished"]; ok {
							if v, ok := v.(string); ok {
								model.DateStr = strings.TrimSpace(v)
								if model.DateStr != "" {
									doubanParseDate(model)
									//if t, err := time.Parse("2006-1-2", model.DateStr); err == nil {
									//	model.Date = t.Unix()
									//} else if t, err := time.Parse("2006-1", model.DateStr); err == nil {
									//	model.Date = t.Unix()
									//} else if t, err := time.Parse("2006", model.DateStr); err == nil {
									//	model.Date = t.Unix()
									//} else {
									//	log.Println(err, model.Url, model.Title, model.DateStr)
									//}
								}
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}

						if v, ok := model.Data["isbn"]; ok {
							if v, ok := v.(string); ok {
								if i, err := strconv.ParseInt(v, 0, 64); err == nil {
									model.Isbn = int64(i)
								}
							} else {
								log.Println(err, model.Url, model.Title, n.FirstChild.Data)
							}
						}
					}
				}
			}
		}

		//IMDb
		if model.Imdb == "" && n.Data == "a" && n.FirstChild != nil {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.Contains(attr.Val, "www.imdb.com") {
					model.Imdb = strings.TrimSpace(n.FirstChild.Data)
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
		//if len(model.Alias) == 0 {
		if n.Data == "span" && n.FirstChild != nil && n.FirstChild.Data == "又名:" {
			if alias := strings.Split(n.NextSibling.Data, "/"); len(alias) > 0 {
				for _, v := range alias {
					model.Alias = append(model.Alias, strings.TrimSpace(v))
				}
			}
		}
		if n.Data == "span" && n.FirstChild != nil && n.FirstChild.Data == "副标题:" {
			if alias := strings.Split(n.NextSibling.Data, "/"); len(alias) > 0 {
				for _, v := range alias {
					model.Alias = append(model.Alias, strings.TrimSpace(v))
				}
			}
		}
		if n.Data == "span" && n.FirstChild != nil && n.FirstChild.Data == "原作名:" {
			if alias := strings.Split(n.NextSibling.Data, "/"); len(alias) > 0 {
				for _, v := range alias {
					model.Alias = append(model.Alias, strings.TrimSpace(v))
				}
			}
		}
		//}

		//date
		if model.DateStr == "" && n.Data == "span" && n.FirstChild != nil && n.FirstChild.Data == "出版年:" && n.NextSibling != nil && n.NextSibling.Data != "" {
			model.DateStr = strings.TrimSpace(n.NextSibling.Data)
			if model.DateStr != "" {
				doubanParseDate(model)
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
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		DouBanPageHtml(c, model)
	}
	return
}

// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
// 一般用于采集数据的地方
func (my *DouBan) ResponseSuccess(spider *spider.Spider) {
	//movie json handle
	if strings.HasPrefix(spider.CurrentRequest().URL.String(), "https://movie.douban.com/j/new_search_subjects") {
		decoder := json.NewDecoder(bytes.NewBuffer(spider.ResponseByte))
		movieJson := make(map[string]interface{})
		if err := decoder.Decode(&movieJson); err != nil {
			return
		}

		if data, ok := movieJson["data"]; ok {
			for _, item := range data.([]interface{}) {
				if m, ok := item.(map[string]interface{}); ok {
					if rawUrl, ok := m["url"]; ok {
						u, _ := url.Parse(rawUrl.(string))
						if enqueueUrl := my.EnqueueFilter(spider, u); enqueueUrl != "" {
							if exists, _ := spider.GetQueue().BlTestAndAddString(enqueueUrl); exists {
								continue
							}
							spider.GetQueue().Enqueue(strings.TrimSpace(enqueueUrl))
						}
					}
				}
			}
		}

		return
	}

	//subject
	my.lastRequestUrl = spider.CurrentRequest().URL.String()
	if !strings.Contains(my.lastRequestUrl, "douban.com/subject/") {
		return
	}

	node, err := html.Parse(ioutil.NopCloser(bytes.NewBuffer(spider.ResponseByte)))
	if err != nil {
		return
	}

	model := &AsukaDouBan{
		//DouBanID: int64(douBanId),
		Url:      spider.CurrentRequest().URL.String(),
		UrlCrc32: int64(crc32.ChecksumIEEE([]byte(spider.CurrentRequest().URL.String()))),
	}

	if paths := strings.Split(spider.CurrentRequest().URL.Path, "/"); len(paths) > 2 {
		model.DouBanId, _ = strconv.ParseInt(paths[2], 0, 64)
	}

	if strings.HasPrefix(spider.CurrentRequest().URL.String(), "https://movie.douban.com") {
		model.Cate = "电影"
	}
	if strings.HasPrefix(spider.CurrentRequest().URL.String(), "https://book.douban.com") {
		model.Cate = "图书"
	}

	//only douBan subject url
	if isDouBanSubject(strings.ToLower(spider.CurrentRequest().URL.String())) {
		DouBanPageHtml(node, model)
		if model.Title == "" {
			return
		}
		if model.DateStr == "" {
			DouBanPageHtmlSecondly(node, model)
		}
	}

	if model.Title == "页面不存在" && model.Name == "" {
		my.EnqueueForFailure(spider, nil, spider.CurrentRequest().URL.String(), spider.CurrentRequest().URL.String(), 3)
	}

	//clear
	model.Title = ""
	model.Data = make(map[string]interface{}, 1)

	//database
	existsModel := &AsukaDouBan{
		UrlCrc32: model.UrlCrc32,
		Url:      model.Url,
	}

	if ok, err := database.Mysql().Get(existsModel); err == nil {
		my.dbSpeedNum++
		if ok {
			//update
			model.Version = existsModel.Version
			if _, err = database.Mysql().Id(existsModel.Id).Update(model); err != nil {
				my.lastInsertError = time.Now().Format(time.RFC3339) + ":" + err.Error()
				log.Println(spider.CurrentRequest().URL.String(), err)
			}
		} else {
			//insert
			_, err = database.Mysql().Insert(model)
			my.lastInsertId = model.Id
			if err != nil {
				my.lastInsertError = time.Now().Format(time.RFC3339) + ":" + err.Error()
				database.MysqlDelayInsertTillSuccess(model)
				log.Println(spider.CurrentRequest().URL.String(), err)
			}
		}
	} else {
		my.lastInsertError = time.Now().Format(time.RFC3339) + ":" + err.Error()
		log.Println(spider.CurrentRequest().URL.String(), err)
	}
}

// queue
func (my *DouBan) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	if !strings.HasPrefix(strings.ToLower(l.String()), "https://movie.douban.com/subject") && !strings.HasPrefix(strings.ToLower(l.String()), "https://book.douban.com/subject") && !strings.HasPrefix(strings.ToLower(l.String()), "https://book.douban.com/tag") && !strings.HasPrefix(strings.ToLower(l.String()), "https://movie.douban.com/tag") {
		return
	}

	if strings.HasPrefix(strings.ToLower(l.String()), "https://book.douban.com/subject") && !isDouBanSubject(strings.ToLower(l.String())) {
		return
	}

	if strings.HasPrefix(strings.ToLower(l.String()), "https://movie.douban.com/subject") && !isDouBanSubject(strings.ToLower(l.String())) {
		return
	}

	//https://movie.douban.com/tag/%E6%97%A5%E8%AF%AD%E6%A0%87%E7%AD%BE?start=1
	query := make(url.Values)
	if l.Query().Get("start") != "" {
		query.Set("start", l.Query().Get("start"))
	}
	if l.Query().Get("type") != "" {
		query.Set("type", l.Query().Get("type"))
	}
	if query.Encode() != "" {
		return l.Scheme + "://" + l.Host + l.Path + "?" + query.Encode()
	}

	return l.Scheme + "://" + l.Host + l.Path
}

func (my *DouBan) WEBSite(w http.ResponseWriter, r *http.Request) {
	var result []*struct {
		//Id       int64
		DouBanId int64
		Name     string
		Date     int64
		//DateStr  string
		Rating float64
		Votes  int64
		Img    string
		Url    string
	}

	database.Mysql().Table("asuka_dou_ban").OrderBy("id desc").Limit(100).Find(&result)

	if byteJson, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-type", "application/json")
		w.Write(byteJson)
	}
}
