package spider

import (
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/proxy"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var linkRegex, _ = regexp.Compile("<a[^>]+href=\"([(\\.|h|/)][^\"]+)\"[^>]*>[^<]+</a>")
var imageRegex, _ = regexp.Compile("(?i)(http(s?):)([/|.|\\w|\\s|-])*\\.(?:jpg|gif|png)")

type Spider struct {
	Transport      *proxy.Transport
	Client         *http.Client
	CurrentRequest *http.Request
	RequestsMap    map[string]*http.Request
	BodyStr        string
	BodyByte       []byte
}

func New(t *proxy.Transport, j *cookiejar.Jar) *Spider {
	if j == nil {
		j, _ = cookiejar.New(nil)
	}
	c := &http.Client{Transport: t.T, Jar: j}
	return &Spider{Transport: t, Client: c, RequestsMap: map[string]*http.Request{}}
}

func (spider *Spider) LoginWithHeader(url *url.URL, header *http.Header) {
	spider.setRequest(url, header).Fetch(url)
}

// setRequest http.Request 是维持session会话的关键之一. 这里是在管理http.Request, 保证每个url能找到对应之前的http.Request
func (spider *Spider) setRequest(url *url.URL, header *http.Header) *Spider {

	tld, err := helper.TldDomain(url.String())
	if err != nil {
		tld = url.String()
	}

	r, ok := spider.RequestsMap[tld]
	if ok {
		spider.CurrentRequest = r
	} else {
		r, err = http.NewRequest("GET", url.String(), nil)
		if err != nil {
			log.Fatal(err)
		}
		spider.CurrentRequest = r
		spider.RequestsMap[tld] = r
	}

	if header != nil {
		for k := range *header {
			spider.CurrentRequest.Header.Set(k, header.Get(k))
		}
	}
	return spider
}

func (spider *Spider) Fetch(url *url.URL) bool {
	spider.setRequest(url, nil)

	spider.Transport.AddAccess(spider.CurrentRequest.URL.String())
	res, httpCode, err := requestUrl(spider.Client, spider.CurrentRequest)
	if err != nil {
		spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
		fmt.Println("Request Error ", err)
		return false
	}
	if httpCode != 200 {
		spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
		fmt.Println("http status" + strconv.Itoa(httpCode))
		return false
	}
	spider.BodyStr = string(res[:])
	spider.BodyByte = res
	return true
}

func (spider *Spider) GetLinks() (arr []*url.URL) {
	for _, sub := range linkRegex.FindAllStringSubmatch(spider.BodyStr, -1) {
		u, err := url.Parse(sub[1])
		if err != nil {
			panic(err)
		}
		arr = append(arr, spider.CurrentRequest.URL.ResolveReference(u))
	}

	return arr
}

func (spider *Spider) GetImageLinks() (arr []*url.URL) {
	for _, sub := range imageRegex.FindAllStringSubmatch(spider.BodyStr, -1) {
		u, err := url.Parse(sub[0])
		if err != nil {
			panic(err)
		}
		arr = append(arr, spider.CurrentRequest.URL.ResolveReference(u))
	}

	return arr
}

func (spider *Spider) Crawl(filter func(spider *Spider, l *url.URL) bool) {
	link, err := database.PopUrlQueue()
	if err != nil {
		fmt.Println("queue is empty: ", err)
		time.Sleep(time.Second * 10)
		return
	}

	u, err := url.Parse(link)
	if err != nil {
		fmt.Println("URL parse filed ", link, err)
		return
	}

	ssArr := spider.Transport.S.ServerAddr
	if ssArr == "" {
		ssArr = "localhost"
	}
	fmt.Println(ssArr, u.String()) //todo

	spider.Fetch(u)
	fmt.Println(strings.Contains(spider.BodyStr, "Golang1"))

	for _, l := range spider.GetLinks() {
		if filter != nil && !filter(spider, l) {
			continue
		}
		//pass := false
		//for _, white := range hostWhiteList {
		//	if strings.Contains(strings.ToLower(l.Hostname()), white) {
		//		pass = true
		//	}
		//}
		//if !pass {
		//	continue
		//}

		if database.Bl().TestAndAddString(l.String()) {
			continue
		}

		database.AddUrlQueue(strings.TrimSpace(l.String()))
	}
}

func requestUrl(httpClient *http.Client, r *http.Request) (res []byte, httpCode int, err error) {
	//todo auto redirect

	// create a request
	//req, err := http.NewRequest("GET", u.String(), nil)
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "can't create request:", err)
	//	os.Exit(2)
	//}

	//request

	resp, err := httpClient.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// read response body
	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return res, resp.StatusCode, nil
}

func requireCookie() *cookiejar.Jar {
	c, _ := cookiejar.New(nil)
	return c
}

func DownloadImage(sp *Spider) {
	for _, u := range sp.GetImageLinks() {
		if database.Bl().TestString(u.String()) {
			continue
		}
		database.Bl().AddString(u.String())

		go func(u *url.URL) {
			//r, _ := http.NewRequest("GET", u.String(), nil)
			imgSp := New(nil, nil)
			imgSp.Fetch(u)

			filename := filepath.Base(u.String())
			if filename == "" {
				filename = strconv.Itoa(rand.Int())
			}

			savePath := helper.WorkspacePath() + filename
			if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
				log.Fatal(err)
			}
			outFile, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			outFile.Write(imgSp.BodyByte)
			outFile.Close()
		}(u)
	}
}
