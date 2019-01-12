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

const RecentFetchCount = 10

var RecentFetchMutex = &sync.Mutex{}

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

	recentFetch := &RecentFetch{Url: u, AddTime: time.Now(), TransportName: spider.Transport.S.Name}
	RecentFetchList = append(RecentFetchList, recentFetch)
	spider.Transport.AddAccess(spider.CurrentRequest.URL.String())

	defer func() {
		if err != nil {
			recentFetch.ErrType = reflect.TypeOf(err).String()
			spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
		}

		spider.TimeList.PushBack(time.Since(spider.RequestStartTime))
		if spider.TimeList.Len() > spider.TimeLenLimit {
			spider.TimeList.Remove(spider.TimeList.Front()) // FIFO
		}

		recentFetch.ConsumeTime = time.Since(spider.RequestStartTime)

		//recent fetch
		RecentFetchMutex.Lock()
		if len(RecentFetchList) > RecentFetchCount {
			RecentFetchList = RecentFetchList[len(RecentFetchList)-RecentFetchCount:]
		}
		RecentFetchMutex.Unlock()

		if r := recover(); r != nil {
			spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
			err = errors.New("spider.Fetch panic:" + fmt.Sprint(r))
		}
	}()

	//traffic
	dump, err := httputil.DumpRequestOut(spider.CurrentRequest, true)
	spider.requestErrorHandler(err)
	spider.Transport.TrafficOut += uint64(len(dump))

	resp, err = spider.Client.Do(spider.CurrentRequest)
	if err != nil {
		spider.requestErrorHandler(err)
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
	spider.responseErrorHandler(err)
	if err != nil {
		return resp, err
	}

	//traffic
	dump, err = httputil.DumpResponse(resp, false)
	spider.responseErrorHandler(err)
	recentFetch.ResponseSize = uint64(len(dump) + len(resByte))
	spider.Transport.TrafficIn += recentFetch.ResponseSize

	//gzip decompression
	reader := ioutil.NopCloser(bytes.NewBuffer(resByte))
	defer reader.Close()
	if strings.ToLower(resp.Header.Get("Content-Encoding")) == "gzip" {
		reader, err = gzip.NewReader(reader)
		spider.responseErrorHandler(err)
	}

	res, err := ioutil.ReadAll(reader)
	spider.responseErrorHandler(err)
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

func (spider *Spider) requestErrorHandler(err error) {
	if err == nil {
		return
	}

	if spider.FailureLevel == 0 {
		spider.FailureLevel = 1
	}

	switch err.(type) {
	case x509.SystemRootsError:
		return
	case *x509.SystemRootsError:
		return
	case x509.UnknownAuthorityError:
		return
	case *x509.UnknownAuthorityError:
		return
	case x509.HostnameError:
		return
	case *x509.HostnameError:
		return
	case *net.DNSConfigError:
		//An error reading the machine's DNS configuration.
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 2)
		return
	case *net.DNSError:
		return //DNSError represents a DNS lookup error
	case *net.OpError:
		log.Println("Request *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
	case net.Error:
		if err.(net.Error).Timeout() {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return
		}
		if io.EOF == err {
			return
		}
		if io.ErrUnexpectedEOF == err {
			return
		}
		if strings.Contains(err.Error(), "transport connection broken") {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 2)
			return
		}
		if strings.Contains(err.Error(), "unexpected EOF") {
			return
		}
		if strings.Contains(err.Error(), "x509: certificate") {
			return
		}
		if strings.Contains(err.Error(), "no such host") {
			return
		}
		if strings.Contains(err.Error(), ": EOF") {
			return
		}

		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)

		//return
		//if !err.(net.Error).Timeout() && err != io.EOF && !strings.Contains(err.Error(), "nection was forcibly closed by the remote ho") && !strings.Contains(err.Error(), "EOF") && !strings.Contains(err.Error(), "no such host") && !strings.Contains(err.Error(), "nnection could be made because the target machine actively refus") {
		log.Println("Request net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		//fmt.Errorf()
		//debug.PrintStack()
		//}
	case *url.Error:
		log.Println("Request Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
	default:
		if strings.HasPrefix(err.Error(), "invalid URL") {
			return
		}

		spider.FailureLevel = 10
		//2019/01/09 11:00:19 spider.go:262: Request Error jp-4.mitsuha-node.com *errors.errorString:  http: no Host in request URL http:
		log.Println("Request Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
	}
}

func (spider *Spider) responseErrorHandler(err error) {
	if err == nil {
		return
	}

	if spider.FailureLevel == 0 {
		spider.FailureLevel = 1
	}

	switch err.(type) {
	case *net.OpError:
		log.Println("Response *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return
	case net.Error:
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
		if err.(net.Error).Timeout() {
			return
		}

		//if !err.(net.Error).Timeout() && err != io.EOF {
		log.Println("Response net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		//}
	case *url.Error:
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
	case tls.RecordHeaderError:
		//2019/01/10 19:43:18 spider.go:306: Response Error hk-4.mitsuha-node.com tls.RecordHeaderError:  tls: received record with version fb2d when expecting version 303 http://gmail.com
		//log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return
	case flate.CorruptInputError:
		//2019/01/10 19:58:51 spider.go:306: Response Error us-7.mitsuha-node.com flate.CorruptInputError:  flate: corrupt input before offset 2446 http://www.woquzhe.com
		//2019/01/10 19:45:38 spider.go:306: Response Error hk-6.mitsuha-node.com flate.CorruptInputError:  flate: corrupt input before offset 8407 http://top.200.net/
		return
	default:
		if strings.HasPrefix(err.Error(), "malformed chunked encoding") {
			//2019/01/10 19:44:39 spider.go:306: Response Error sg-1.mitsuha-node.com *errors.errorString:  malformed chunked encoding http://wxip.me
			return
		}
		if strings.HasPrefix(err.Error(), "invalid URL") {
			return
		}
		if strings.HasPrefix(err.Error(), "http: unexpected EOF reading trailer") {
			return
		}
		if strings.HasPrefix(err.Error(), "http:  reading trailer") {
			//fmt.Println("http: unexpected EOF reading trailer !!!!!!!!!!!")
			return
		}
		if gzip.ErrHeader == err {
			//fmt.Println(err)
			return
		}
		if gzip.ErrChecksum == err {
			//fmt.Println(err)
			return
		}
		if io.EOF == err {
			return
		}
		if io.ErrUnexpectedEOF == err {
			//2019/01/09 11:37:09 spider.go:281: Response Error hk-8.mitsuha-node.com *errors.errorString:  gzip: invalid checksum http://www.fcx110.com/
			//2019/01/08 21:05:45 spider.go:252: Response Error hk-1a.mitsuha-node.com *errors.errorString:  gzip: invalid header http://www.s80.cc
			//2019/01/08 19:23:55 spider.go:251: Response Error jp-2.mitsuha-node.com *errors.errorString:  malformed chunked encoding http://www.jygedu.net/
			//2019/01/09 11:11:36 spider.go:281: Response Error jp-b.mitsuha-node.com flate.CorruptInputError:  flate: corrupt input before offset 2168 http://www.wyzc.com/
			//log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
			return
		}

		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
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
