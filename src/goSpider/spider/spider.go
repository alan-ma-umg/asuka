package spider

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"container/list"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/proxy"
	"goSpider/queue"
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
	"sync"
	"time"
)

var linkRegex, _ = regexp.Compile("<a[^>]+href=\"([(\\.|h|/)][^\"]+)\"[^>]*>[^<]+</a>")
var imageRegex, _ = regexp.Compile("(?i)(http(s?):)([/|.|\\w|\\s|-])*\\.(?:jpg|gif|png)")

const RecentFetchCount = 100

var RecentFetchMutex = &sync.Mutex{}
var RecentFetchLastIndex int64 = 0

type RecentFetch struct {
	Index         int64
	TransportName string
	StatusCode    int // http response status code
	RawUrl        string
	ConsumeTime   string
	AddTime       string
	ErrType       string
	ResponseSize  string
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

	FailureLevel int

	RequestStartTime time.Time
	Stop             bool

	Queue *queue.Queue
}

func New(t *proxy.Transport, j *cookiejar.Jar, queue *queue.Queue) *Spider {
	if j == nil {
		j, _ = cookiejar.New(nil)
	}
	c := &http.Client{Transport: t.T, Jar: j, Timeout: time.Second * 30}
	return &Spider{Queue: queue, Transport: t, Client: c, RequestsMap: map[string]*http.Request{}, TimeList: list.New(), TimeLenLimit: 10}
}

func (spider *Spider) Throttle() {
	if spider.Transport.S.Interval > .0 {
		time.Sleep(time.Second * time.Duration(spider.Transport.S.Interval))
	}
	for {
		//todo make improvement
		if !spider.Stop {
			break
		}
		time.Sleep(2e9)
	}

	if spider.FailureLevel > 0 {
		time.Sleep(time.Second)
	}
	if spider.FailureLevel > 1 {
		time.Sleep(time.Minute)
	}

	accessCount, failureCount := spider.Transport.AccessCount(60)                      //fixme 速度低于每分钟7次这里永远不会发生
	if accessCount > 7 && helper.SpiderFailureRate(accessCount, failureCount) > 50.0 { //fixme 速度低于每分钟7次这里永远不会发生
		spider.FailureLevel = 100
		accessCountAll := spider.Transport.GetAccessCount()
		failureCountAll := spider.Transport.GetFailureCount()
		failureRateAll := helper.SpiderFailureRate(accessCountAll, failureCountAll)
		if accessCountAll > 40 && failureRateAll > 95 {
			spider.FailureLevel = 100
			time.Sleep(time.Hour * 2)
		} else if accessCountAll > 40 && failureRateAll > 85 {
			spider.FailureLevel = 80
			time.Sleep(time.Minute * 40)
		} else if accessCountAll > 30 && failureRateAll > 70 {
			spider.FailureLevel = 60
			time.Sleep(time.Minute * 10)
		} else if accessCountAll > 30 && failureRateAll > 60 {
			spider.FailureLevel = 40
			time.Sleep(time.Minute * 5)
		} else {
			spider.FailureLevel = 20
			time.Sleep(time.Minute * 2)
		}
	}

	spider.FailureLevel = 0
}

// setRequest http.Request 是维持session会话的关键之一. 这里是在管理http.Request, 保证每个url能找到对应之前的http.Request
func (spider *Spider) SetRequest(url *url.URL, header *http.Header) *Spider {

	tld, err := helper.TldDomain(url)
	if err != nil {
		tld = "DefaultRequest"
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

	//spider.CurrentRequest.Close = true // prevents re-use of TCP connections between requests to the same hosts

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

	spider.ResponseStr = ""
	spider.ResponseByte = []byte{}

	//time
	spider.RequestStartTime = time.Now()

	recentFetch := &RecentFetch{RawUrl: u.String(), AddTime: time.Now().Format("01-02 15:04:05"), TransportName: spider.Transport.S.Name}

	spider.Transport.AddAccess(spider.CurrentRequest.URL.String())

	defer func() {
		if err != nil {
			spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
		}

		spider.TimeList.PushBack(time.Since(spider.RequestStartTime))
		if spider.TimeList.Len() > spider.TimeLenLimit {
			spider.TimeList.Remove(spider.TimeList.Front()) // FIFO
		}

		//recent fetch
		RecentFetchLastIndex++
		recentFetch.Index = RecentFetchLastIndex
		recentFetch.ConsumeTime = time.Since(spider.RequestStartTime).Truncate(time.Millisecond).String()
		RecentFetchList = append(RecentFetchList, recentFetch)

		RecentFetchMutex.Lock()
		if len(RecentFetchList) > RecentFetchCount {
			RecentFetchList = RecentFetchList[len(RecentFetchList)-RecentFetchCount:]
		}
		RecentFetchMutex.Unlock()

		//recover
		if r := recover(); r != nil {
			spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
			err = errors.New("spider.Fetch panic:" + fmt.Sprint(r))
		}
	}()

	//traffic
	dump, err := httputil.DumpRequestOut(spider.CurrentRequest, true)
	recentFetch.ErrType = spider.requestErrorHandler(err)
	spider.Transport.TrafficOut += uint64(len(dump))

	resp, err = spider.Client.Do(spider.CurrentRequest)
	if err != nil {
		recentFetch.ErrType = spider.requestErrorHandler(err)
		return resp, err
	}
	defer resp.Body.Close()
	recentFetch.StatusCode = resp.StatusCode

	//todo remove
	if !strings.Contains(resp.Header.Get("Content-type"), "text/html") {
		//return resp, errors.New("Content-type:Content-type must be text/html, " + resp.Header.Get("Content-type") + " given")
	}
	//todo remove
	if strings.ToLower(resp.Header.Get("Content-Encoding")) != "gzip" {
		//return resp, errors.New("Content-Encoding must be gzip , " + resp.Header.Get("Content-Encoding") + " given")
	}

	resByte, err := ioutil.ReadAll(resp.Body)
	recentFetch.ErrType = spider.responseErrorHandler(err)
	if err != nil {
		return resp, err
	}

	//traffic
	dump, err = httputil.DumpResponse(resp, false)
	recentFetch.ErrType = spider.responseErrorHandler(err)
	recentFetch.ResponseSize = helper.ByteCountBinary(uint64(len(dump) + len(resByte)))
	spider.Transport.TrafficIn += uint64(len(dump) + len(resByte))

	//gzip decompression
	reader := ioutil.NopCloser(bytes.NewBuffer(resByte))
	defer reader.Close()
	if strings.ToLower(resp.Header.Get("Content-Encoding")) == "gzip" {
		reader, err = gzip.NewReader(reader)
		recentFetch.ErrType = spider.responseErrorHandler(err)
	}

	res, err := ioutil.ReadAll(reader)
	recentFetch.ErrType = spider.responseErrorHandler(err)
	if err != nil {
		return resp, err
	}

	//http status
	if resp.StatusCode != 200 && err == nil {
		spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
	}

	spider.ResponseStr = string(res[:])
	spider.ResponseByte = res
	return resp, err
}

func (spider *Spider) requestErrorHandler(err error) string {
	if err == nil {
		return ""
	}

	if spider.FailureLevel == 0 {
		spider.FailureLevel = 1
	}

	switch err.(type) {
	case *x509.SystemRootsError:
		return "x509.SystemRootsError"
	case *x509.UnknownAuthorityError:
		return "x509.UnknownAuthorityError"
	case *x509.HostnameError:
		return "x509.HostnameError"
	case *net.DNSConfigError:
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 2)
		return "x509.DNSConfigError"
	case *net.DNSError:
		return "x509.DNSError"
	case *net.OpError:
		log.Println("Request *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "net.OpError"
	case net.Error:
		if err.(net.Error).Timeout() {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return "net.Timeout"
		}
		if io.EOF == err {
			return "io.EOF"
		}
		if io.ErrUnexpectedEOF == err {
			return "io.ErrUnexpectedEOF"
		}
		if strings.Contains(err.Error(), "transport connection broken") {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 2)
			return "connection broken"
		}
		if strings.Contains(err.Error(), "unexpected EOF") {
			return "unexpected EOF"
		}
		if strings.Contains(err.Error(), "x509: certificate") {
			return "x509: certificate"
		}
		if strings.Contains(err.Error(), "no such host") {
			return "no such host"
		}
		if strings.Contains(err.Error(), ": EOF") {
			return "other EOF"
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return "reset by peer"
		}
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
		log.Println("Request net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "unknown"
	case *url.Error:
		log.Println("Request Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "url.Error"
	default:
		if strings.HasPrefix(err.Error(), "invalid URL") {
			return "invalid URL"
		}
		if strings.HasPrefix(err.Error(), "no Host in request URL http") {
			return "no Host"
		}
		spider.FailureLevel = 10
		log.Println("Request Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "unknown"
	}
}

func (spider *Spider) responseErrorHandler(err error) string {
	if err == nil {
		return ""
	}

	if spider.FailureLevel == 0 {
		spider.FailureLevel = 1
	}

	switch err.(type) {
	case *net.OpError:
		log.Println("Response *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "net.OpError"
	case net.Error:
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
		if err.(net.Error).Timeout() {
			return "net.Timeout"
		}
		log.Println("Response net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "net.Error"
	case *url.Error:
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "url.Error"
	case tls.RecordHeaderError:
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "tls.RecordHeaderError"
	case flate.CorruptInputError:
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "flate.CorruptInputError"
	default:
		if strings.HasPrefix(err.Error(), "malformed chunked encoding") {
			return "chunked encoding"
		}
		if strings.HasPrefix(err.Error(), "invalid URL") {
			return "invalid URL"
		}
		if strings.HasPrefix(err.Error(), "http: unexpected EOF reading trailer") {
			return "unexpected EOF reading trailer"
		}
		if strings.HasPrefix(err.Error(), "http:  reading trailer") {
			return "http.reading trailer"
		}
		if gzip.ErrHeader == err {
			return "gzip.ErrHeader"
		}
		if gzip.ErrChecksum == err {
			return "gzip.ErrChecksum"
		}
		if io.EOF == err {
			return "io.EOF"
		}
		if io.ErrUnexpectedEOF == err {
			return "io.ErrUnexpectedEOF"
		}
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "unknown"
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
							continue
						}
						addUrl := spider.CurrentRequest.URL.ResolveReference(u)
						if addUrl.Scheme != "http" && addUrl.Scheme != "https" {
							continue
						}

						if len(addUrl.Hostname()) < 4 || strings.Index(addUrl.Hostname(), ".") == -1 || !helper.OnlyDomainCharacter(addUrl.Hostname()) {
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
	link, err := spider.Queue.Dequeue()
	if err != nil {
		time.Sleep(time.Second * 5)
		return
	}

	u, err := url.Parse(link)
	if err != nil {
		log.Println("URL parse failed ", link, err)
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

		if database.BlTestAndAddString(l.String()) {
			continue
		}

		spider.Queue.Enqueue(strings.TrimSpace(l.String()))
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
