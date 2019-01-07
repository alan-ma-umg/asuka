package spider

import (
	"bytes"
	"compress/gzip"
	"container/list"
	"errors"
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/proxy"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var linkRegex, _ = regexp.Compile("<a[^>]+href=\"([(\\.|h|/)][^\"]+)\"[^>]*>[^<]+</a>")
var imageRegex, _ = regexp.Compile("(?i)(http(s?):)([/|.|\\w|\\s|-])*\\.(?:jpg|gif|png)")

type RecentFetch struct {
	TransportName string
	StatusCode    int // http response status code
	Url           *url.URL
	ConsumeTime   time.Duration
	AddTime       time.Time
	ErrType       string
	ResponseSize  uint64
}

var RecentFetchList []*RecentFetch

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

	ConnectFail bool

	RequestStartTime time.Time
}

func New(t *proxy.Transport, j *cookiejar.Jar) *Spider {
	if j == nil {
		j, _ = cookiejar.New(nil)
	}
	c := &http.Client{Transport: t.T, Jar: j, Timeout: time.Second * 30}
	return &Spider{Transport: t, Client: c, RequestsMap: map[string]*http.Request{}, TimeList: list.New(), TimeLenLimit: 10}
}

func (spider *Spider) Throttle() {
	if spider.ConnectFail {
		time.Sleep(time.Minute)
		spider.ConnectFail = false
	}

	failureRate := spider.Transport.FailureRate(30)
	if failureRate > 0.3 {
		spider.ConnectFail = true
		time.Sleep(time.Minute)
		spider.ConnectFail = false
	}

	if failureRate > 0.5 {
		spider.ConnectFail = true
		time.Sleep(time.Hour)
		spider.ConnectFail = false
	}
}

// setRequest http.Request 是维持session会话的关键之一. 这里是在管理http.Request, 保证每个url能找到对应之前的http.Request
func (spider *Spider) SetRequest(url *url.URL, header *http.Header) *Spider {

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

		//Accept-Encoding: gzip
		if r.Header.Get("Accept-Encoding") == "" {
			r.Header.Set("Accept-Encoding", "gzip")
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
		spider.CurrentRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0."+strconv.FormatFloat(rand.Float64()*10000, 'f', 3, 64)+" Safari/537.36")
	}
	return spider
}

func (spider *Spider) Fetch(u *url.URL) (resp *http.Response, err error) {
	spider.SetRequest(u, nil)

	//time
	spider.RequestStartTime = time.Now()

	recentFetch := &RecentFetch{Url: u, AddTime: time.Now(), TransportName: spider.Transport.S.Name}
	RecentFetchList = append(RecentFetchList, recentFetch)
	spider.Transport.AddAccess(spider.CurrentRequest.URL.String())

	defer func() {
		if err != nil {
			recentFetch.ErrType = reflect.TypeOf(err).String()
		}

		recentFetch.ConsumeTime = time.Since(spider.RequestStartTime)

		recentFetchCount := 100
		if len(RecentFetchList) > recentFetchCount {
			RecentFetchList = RecentFetchList[len(RecentFetchList)-recentFetchCount:]
		}

		spider.TimeList.PushBack(time.Since(spider.RequestStartTime))
		if spider.TimeList.Len() > spider.TimeLenLimit {
			spider.TimeList.Remove(spider.TimeList.Front()) // FIFO
		}

		if r := recover(); r != nil {
			err = errors.New("spider.Fetch panic:" + fmt.Sprint(r))
		}
	}()

	resp, requestErr := spider.Client.Do(spider.CurrentRequest)
	if requestErr != nil {
		spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
		spider.requestErrorHandler(requestErr)
		return resp, requestErr
	}
	defer resp.Body.Close()

	//todo remove
	if !strings.Contains(resp.Header.Get("Content-type"), "text/html") {
		return resp, errors.New("Content-type:Content-type must be text/html, " + resp.Header.Get("Content-type") + " given")
	}

	//traffic count
	dump, err := httputil.DumpRequestOut(spider.CurrentRequest, true)
	if err == nil {
		spider.Transport.TrafficOut += uint64(len(dump))
	} else {
		spider.requestErrorHandler(err)
	}

	dump, err = httputil.DumpResponse(resp, true)
	if err == nil {
		recentFetch.ResponseSize = uint64(len(dump))
		spider.Transport.TrafficIn += recentFetch.ResponseSize
	} else {
		spider.responseErrorHandler(err)
	}

	recentFetch.StatusCode = resp.StatusCode

	//http status
	if resp.StatusCode != 200 {
		spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
		//log.Println("http status", resp.StatusCode, spider.CurrentRequest.URL.String())
	}

	//gzip decompression
	var reader io.ReadCloser
	switch strings.ToLower(resp.Header.Get("Content-Encoding")) {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		//*http.httpError todo  2019/01/05 14:58:53 Gzip Error:*http.httpError : read tcp 127.0.0.1:10281->127.0.0.1:10199: use of closed network connection (Client.Timeout exceeded while reading body)
		if err != nil {
			log.Println("Gzip Error:" + reflect.TypeOf(err).String() + " : " + err.Error())
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	res, err := ioutil.ReadAll(reader)
	if err != nil {
		spider.responseErrorHandler(err)
		return resp, err
	}

	spider.ResponseStr = string(res[:])
	spider.ResponseByte = res
	return resp, err
}

func (spider *Spider) requestErrorHandler(err error) {
	switch err.(type) {
	case *net.OpError:
		log.Println("Request *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
	case net.Error:
		return
		if !err.(net.Error).Timeout() && err != io.EOF && !strings.Contains(err.Error(), "nection was forcibly closed by the remote ho") && !strings.Contains(err.Error(), "EOF") && !strings.Contains(err.Error(), "no such host") && !strings.Contains(err.Error(), "nnection could be made because the target machine actively refus") {
			log.Println("Request net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
			//fmt.Errorf()
			//debug.PrintStack()
		}
	case *url.Error:
		log.Println("Request Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
	default:
		spider.ConnectFail = true
		println("Request Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
	}
}

func (spider *Spider) responseErrorHandler(err error) {
	switch err.(type) {
	case *net.OpError:
		return
		log.Println("Response *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
	case net.Error:
		if !err.(net.Error).Timeout() && err != io.EOF {
			log.Println("Response net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		}
	case *url.Error:
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
	default:
		if io.ErrUnexpectedEOF != err {
			log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		}
	}
}

func (spider *Spider) GetAvgTime() (t time.Duration) {
	count := 0
	cursor := spider.TimeList.Back()
	for {
		if cursor == nil {
			break
		}

		count++
		t += cursor.Value.(time.Duration)
		cursor = cursor.Prev()
	}

	if count == 0 {
		return
	}
	t /= time.Duration(count)
	return
}

func (spider *Spider) GetLinksByTokenizer() (res []*url.URL) {
	token := html.NewTokenizer(ioutil.NopCloser(bytes.NewBuffer(spider.ResponseByte)))
	for {
		switch next := token.Next(); next {
		case html.StartTagToken:
			token := token.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						value := strings.TrimSpace(attr.Val)
						if value == "" {
							continue
						}
						u, err := url.Parse(value)
						if err != nil {
							//log.Println("Url Parse: " + reflect.TypeOf(err).String() + " : " + value + " : " + err.Error() + " : from : " + spider.CurrentRequest.URL.String())
							continue
						}
						addUrl := spider.CurrentRequest.URL.ResolveReference(u)
						if addUrl.Scheme != "http" && addUrl.Scheme != "https" {
							continue
						}
						res = append(res, addUrl)
					}
				}
			}
		case html.ErrorToken:
			return
		}
	}
	return
}

func (spider *Spider) GetLinksByRegex() (arr []*url.URL) {
	for _, sub := range linkRegex.FindAllStringSubmatch(spider.ResponseStr, -1) {
		u, err := url.Parse(sub[1])
		if err != nil {
			continue
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
		time.Sleep(time.Second * 5)
		return
	}

	u, err := url.Parse(link)
	if err != nil {
		log.Println("URL parse filed ", link, err)
		return
	}

	ssArr := spider.Transport.S.ServerAddr
	if ssArr == "" {
		ssArr = "localhost"
	}

	spider.Transport.LoopCount++
	_, err = spider.Fetch(u)
	if err != nil {
		//todo register error handler
		//log.Println("Fetch Fial: "+reflect.TypeOf(err).String()+" : "+u.String(), err)
		return
	}

	for _, l := range spider.GetLinksByTokenizer() {
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
