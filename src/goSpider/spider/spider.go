package spider

import (
	"container/list"
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var linkRegex, _ = regexp.Compile("<a[^>]+href=\"([(\\.|h|/)][^\"]+)\"[^>]*>[^<]+</a>")
var imageRegex, _ = regexp.Compile("(?i)(http(s?):)([/|.|\\w|\\s|-])*\\.(?:jpg|gif|png)")

type Spider struct {
	Transport *proxy.Transport
	Client    *http.Client

	RequestsMap     map[string]*http.Request
	CurrentRequest  *http.Request
	CurrentResponse *http.Response

	ResponseStr  string
	ResponseByte []byte

	TimeList     *list.List
	TimeLenLimit int
}

func New(t *proxy.Transport, j *cookiejar.Jar) *Spider {
	if j == nil {
		j, _ = cookiejar.New(nil)
	}
	c := &http.Client{Transport: t.T, Jar: j}
	return &Spider{Transport: t, Client: c, RequestsMap: map[string]*http.Request{}, TimeList: list.New(), TimeLenLimit: 10}
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

	if spider.CurrentRequest.UserAgent() == "" {
		spider.CurrentRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	}
	return spider
}

func (spider *Spider) Fetch(url *url.URL) (*http.Response, error) {
	spider.setRequest(url, nil)

	spider.Transport.AddAccess(spider.CurrentRequest.URL.String())

	st := time.Now()
	defer func() {
		spider.TimeList.PushBack(time.Since(st))
		if spider.TimeList.Len() > spider.TimeLenLimit {
			spider.TimeList.Remove(spider.TimeList.Front()) // FIFO
		}
	}()

	resp, err := spider.Client.Do(spider.CurrentRequest)

	//res, httpCode, err := requestUrl(spider.Client, spider.CurrentRequest)
	if err != nil {
		spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
		fmt.Println("Request Error ", err)
		return resp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
		fmt.Println("http status", resp.StatusCode)
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	spider.ResponseStr = string(res[:])
	spider.ResponseByte = res
	return resp, err
}

func (spider *Spider) GetAvgTime() (t time.Duration) {
	//var all time.Duration
	count := spider.TimeList.Len()
	for i := 0; i < count; i++ {
		cursor := spider.TimeList.Back()
		if cursor == nil {
			break
		}

		t += cursor.Value.(time.Duration)
	}

	if count == 0 {
		return
	}
	t /= time.Duration(count)

	return
}

func (spider *Spider) GetLinks() (arr []*url.URL) {
	for _, sub := range linkRegex.FindAllStringSubmatch(spider.ResponseStr, -1) {
		u, err := url.Parse(sub[1])
		if err != nil {
			panic(err)
		}
		arr = append(arr, spider.CurrentRequest.URL.ResolveReference(u))
	}

	return arr
}

func (spider *Spider) GetImageLinks() (arr []*url.URL) {
	for _, sub := range imageRegex.FindAllStringSubmatch(spider.ResponseStr, -1) {
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
		//fmt.Println("queue is empty: ", err)
		time.Sleep(time.Second * 5)
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

	//st := time.Now()
	_, err = spider.Fetch(u)
	//fmt.Println(ssArr, time.Since(st), u.String()) //todo
	if err != nil {
		fmt.Println(u.String(), err)
		return
	}

	for _, l := range spider.GetLinks() {
		if filter != nil && !filter(spider, l) {
			continue
		}

		if database.Bl().TestAndAddString(l.String()) {
			continue
		}

		database.AddUrlQueue(strings.TrimSpace(l.String()))
	}
}

//func DownloadImage(sp *Spider) {
//	for _, u := range sp.GetImageLinks() {
//		if database.Bl().TestString(u.String()) {
//			continue
//		}
//		database.Bl().AddString(u.String())
//
//		go func(u *url.URL) {
//			//r, _ := http.NewRequest("GET", u.String(), nil)
//			imgSp := New(nil, nil)
//			imgSp.Fetch(u)
//
//			filename := filepath.Base(u.String())
//			if filename == "" {
//				filename = strconv.Itoa(rand.Int())
//			}
//
//			savePath := helper.WorkspacePath() + filename
//			if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
//				log.Fatal(err)
//			}
//			outFile, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
//			if err != nil {
//				log.Fatal(err)
//			}
//			outFile.Write(imgSp.BodyByte)
//			outFile.Close()
//		}(u)
//	}
//}
