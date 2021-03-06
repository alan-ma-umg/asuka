package spider

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/proxy"
	"github.com/chenset/asuka/queue"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var opts = []chromedp.ExecAllocatorOption{
	chromedp.NoFirstRun,
	chromedp.NoDefaultBrowserCheck,
	chromedp.Headless,

	chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36"),

	//chromedp.ProxyServer("socks5://127.0.0.1:1080"),
	// After Puppeteer's default behavior.
	chromedp.Flag("blink-settings", "imagesEnabled=false"),
	chromedp.Flag("disable-background-networking", true),
	chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
	chromedp.Flag("disable-background-timer-throttling", true),
	chromedp.Flag("disable-backgrounding-occluded-windows", true),
	chromedp.Flag("disable-breakpad", true),
	chromedp.Flag("disable-client-side-phishing-detection", true),
	chromedp.Flag("disable-default-apps", true),
	chromedp.Flag("disable-dev-shm-usage", true),
	chromedp.Flag("disable-extensions", true),
	chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
	chromedp.Flag("disable-hang-monitor", true),
	chromedp.Flag("disable-ipc-flooding-protection", true),
	chromedp.Flag("disable-popup-blocking", true),
	chromedp.Flag("disable-prompt-on-repost", true),
	chromedp.Flag("disable-renderer-backgrounding", true),
	chromedp.Flag("disable-sync", true),
	chromedp.Flag("force-color-profile", "srgb"),
	chromedp.Flag("metrics-recording-only", true),
	chromedp.Flag("safebrowsing-disable-auto-update", true),
	chromedp.Flag("enable-automation", true),
	chromedp.Flag("password-store", "basic"),
	chromedp.Flag("use-mock-keychain", true),
}

var allocCtx context.Context
var allocCtxCancel context.CancelFunc
var allocCtxDoOnce = &sync.Once{}

func resetAllocCtx() {
	if allocCtxCancel != nil {
		allocCtxCancel()
	}
	allocCtxDoOnce = new(sync.Once)
}

func getAllocCtx() context.Context {
	allocCtxDoOnce.Do(func() {
		//allocCtx, allocCtxCancel = chromedp.NewRemoteAllocator(context.Background(), "ws://10.0.0.2:9223/devtools/browser/ea74c89b-f787-4b1d-89d7-dace64d8b605")
		allocCtx, allocCtxCancel = chromedp.NewExecAllocator(context.Background(), opts...)
	})
	return allocCtx
}

const PeriodOfFailureSecond = 86400 / 2

type Summary struct {
	Index         int64
	TransportName string
	StatusCode    int // http response status code
	FindUrls      int //页面中找到的全部url
	NewUrls       int //页面中找到的全部新url, 新入队列的url数量
	RawUrl        string
	DnsLookupTime string
	TLSTime       string
	ConsumeTime   string
	AddTime       string
	ErrType       string
	TrafficInStr  string
	TrafficIn     uint64
	TrafficOut    uint64
	ContentType   string
}

type Spider struct {
	*helper.Counting

	transport    *proxy.Transport
	TransportUrl *url.URL
	client       *http.Client

	//requestsMap     map[string]*http.Request
	currentRequest  *http.Request
	currentResponse *http.Response

	//ResponseStr  string
	ResponseByte       []byte
	ResponseStatusCode int

	FailureLevel int

	startTime         time.Time
	RequestStartTime  time.Time
	RequestEndTime    time.Time
	dnsLookupDuration time.Duration
	tlsDuration       time.Duration
	requestTimeout    time.Duration
	Stop              bool
	Delete            bool
	sleepDuration     time.Duration
	GetQueue          func() *queue.Queue
	RequestBefore     func(spider *Spider)
	DownloadFilter    func(spider *Spider, response *http.Response) (bool, error)
	ProjectThrottle   func(spider *Spider)

	EnqueueForFailure func(spider *Spider, err error, rawUrl string, retryTimes int)

	//httpTrace                   *httptrace.ClientTrace
	RecentSeveralTimesResultCap int //改成方法, 让project可以灵活调用修改
	recentFewTimesResult        []bool
}

func New(transportUrl *url.URL, getQueue func() *queue.Queue) *Spider {
	return &Spider{TransportUrl: transportUrl, GetQueue: getQueue, startTime: time.Now(), Counting: &helper.Counting{}, RecentSeveralTimesResultCap: 5, requestTimeout: time.Minute}
}

// ResetSpider remove http.Client,http.Request,http.Client.Transport and http.Response than release memory
func (spider *Spider) ResetSpider() {
	spider.ResponseByte = nil
	spider.ResetResponse()
	spider.ResetClient()
	spider.ResetRequest()
}

func (spider *Spider) CurrentRequest() *http.Request {
	return spider.currentRequest
}

func (spider *Spider) CurrentResponse() *http.Response {
	return spider.currentResponse
}

func (spider *Spider) ResetResponse() {
	spider.currentResponse = nil
}

func (spider *Spider) SetRequestTimeout(duration time.Duration) {
	spider.requestTimeout = duration
}

func (spider *Spider) httpClient() *http.Client {
	if spider.client == nil || spider.client.Transport == nil || spider.transport == nil || spider.client.Transport.(*http.Transport) != spider.transport.Connect(spider.TransportUrl) {
		spider.ResetClient()

		// new one
		spider.transport = proxy.NewTransport(spider.TransportUrl)
		j, _ := cookiejar.New(nil)
		//todo 要支持无需重置spider而单独刷新cookie
		spider.client = &http.Client{Transport: spider.transport.Connect(spider.TransportUrl), Jar: j, Timeout: spider.requestTimeout}
	}

	return spider.client
}

func (spider *Spider) ResetClient() {
	if spider.client != nil {
		spider.client.CloseIdleConnections()
	}

	if spider.transport != nil {
		spider.transport.Close()
	}

	spider.client = nil
	spider.transport = nil
}

func (spider *Spider) AddSleep(duration time.Duration) {
	spider.sleepDuration += duration
}

func (spider *Spider) GetSleep() time.Duration {
	return spider.sleepDuration
}

func (spider *Spider) ResetSleep() {
	spider.sleepDuration = 0
}

func (spider *Spider) Throttle(dispatcherCallback func(spider *Spider)) {
	for {
		if !spider.Stop {
			break
		}

		spider.ResetSpider()
		time.Sleep(5e9)
	}

	if spider.FailureLevel > 0 {
		spider.AddSleep(time.Second)
	}
	if spider.FailureLevel > 1 {
		spider.AddSleep(time.Second * 30)
	}

	//Failure control
	if len(spider.recentFewTimesResult) >= spider.RecentSeveralTimesResultCap {
		spider.recentFewTimesResult = spider.recentFewTimesResult[len(spider.recentFewTimesResult)-spider.RecentSeveralTimesResultCap:]
		failCount := 0
		for _, v := range spider.recentFewTimesResult {
			if !v {
				failCount++
			}
		}
		if float64(failCount)/float64(spider.RecentSeveralTimesResultCap) >= 0.4 {
			spider.recentFewTimesResult = make([]bool, 0, spider.RecentSeveralTimesResultCap)

			accessCountAll, failureCountAll := spider.AccessCount(helper.MinInt(int(time.Since(spider.startTime).Seconds()), PeriodOfFailureSecond))
			failureRateAll := helper.SpiderFailureRate(accessCountAll, failureCountAll)
			if accessCountAll > 40 && failureRateAll > 95 {
				spider.FailureLevel = 100
				spider.AddSleep(time.Hour * 5)
			} else if accessCountAll > 40 && failureRateAll > 85 {
				spider.FailureLevel = 80
				spider.AddSleep(time.Hour)
			} else if accessCountAll > 30 && failureRateAll > 70 {
				spider.FailureLevel = 60
				spider.AddSleep(time.Minute * 30)
			} else if accessCountAll > 30 && failureRateAll > 60 {
				spider.FailureLevel = 40
				spider.AddSleep(time.Minute * 10)
			} else if spider.FailureLevel <= 20 {
				spider.FailureLevel = 20
				spider.AddSleep(time.Minute * 2)
			}
		}
	}

	spider.ProjectThrottle(spider)
	if dispatcherCallback != nil {
		dispatcherCallback(spider)
	}

	//exit check
	if spider.Delete {
		return
	}

	//go to sleep and reset sleep duration
	if duration := spider.GetSleep(); duration > 0 {
		time.Sleep(duration)
	}
	spider.ResetSleep()

	//reset failureLevel
	if spider.GetSleep() == 0 { //Maybe change by another goroutine when time.sleep
		spider.FailureLevel = 0
	}
}
func (spider *Spider) IsIdle() bool {
	if spider.ResponseByte == nil {
		if spider.client == nil || !spider.RequestEndTime.IsZero() || spider.RequestStartTime.IsZero() {
			return true
		}
	}

	return false
}

// setRequest http.Request 是维持session会话的关键之一. 这里是在管理http.Request, 保证每个url能找到对应之前的http.Request
func (spider *Spider) SetRequest(url *url.URL) *Spider {

	//todo requestsMap 还是需要做起来 !!!!!!!!!!!!!!!!

	//tld, err := helper.TldDomain(url)
	//if err != nil {
	//	tld = "DefaultRequest"
	//}
	//
	//r, ok := spider.requestsMap[tld]
	//if ok {
	//	r.URL = url
	//	spider.currentRequest = r
	//} else {
	//	r, err = http.NewRequest("GET", url.String(), nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	//Accept-Encoding: gzip
	//	if r.Header.Get("Accept-Encoding") == "" {
	//		r.Header.Set("Accept-Encoding", "gzip")
	//	}
	//
	//	spider.currentRequest = r
	//	//spider.requestsMap[tld] = r
	//}

	//spider.currentRequest.Close = true // prevents re-use of TCP connections between requests to the same hosts

	//if header != nil {
	//	for k := range *header {
	//		spider.currentRequest.Header.Set(k, header.Get(k))
	//	}
	//}

	if spider.currentRequest == nil {
		spider.currentRequest, _ = http.NewRequest("GET", url.String(), nil)
		if spider.currentRequest.Header.Get("Accept-Encoding") == "" {
			spider.currentRequest.Header.Set("Accept-Encoding", "gzip")
		}
		if spider.currentRequest.UserAgent() == "" {
			spider.currentRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/"+strconv.FormatFloat(rand.Float64()*10000, 'f', 3, 64)+" (KHTML, like Gecko) Chrome/77.0."+strconv.FormatFloat(rand.Float64()*10000, 'f', 3, 64)+" Safari/537.36")
		}

		// http trace
		spider.currentRequest = spider.currentRequest.WithContext(httptrace.WithClientTrace(spider.currentRequest.Context(), spider.getHttpTrace()))

	} else {
		spider.currentRequest.URL = url
	}

	//if r.Header.Get("Accept-Encoding") == "" {
	//	r.Header.Set("Accept-Encoding", "gzip")
	//}

	//spider.currentRequest = spider.currentRequest.WithContext(httptrace.WithClientTrace(spider.currentRequest.Context(), spider.httpTrace))
	return spider
}

func (spider *Spider) ResetRequest() {
	//spider.requestsMap = map[string]*http.Request{}
	spider.currentRequest = nil
}

func (spider *Spider) ChromeFetch(u *url.URL) (summary *Summary, err error) {
	spider.SetRequest(u) //setting spider.currentRequest
	//time
	spider.RequestStartTime = time.Now()
	spider.RequestEndTime = time.Time{} //empty

	spider.ResponseStatusCode = 0

	summary = &Summary{RawUrl: u.String(), AddTime: time.Now().Format("01-02 15:04:05"), TransportName: spider.TransportUrl.Host}

	defer func() {
		if err == nil {
			summary.StatusCode = 200
		}

		spider.ResponseStatusCode = summary.StatusCode
		spider.RequestEndTime = time.Now()
		//A few times result of http request
		spider.recentFewTimesResult = append(spider.recentFewTimesResult, spider.FailureLevel == 0)

		//recover
		if r := recover(); r != nil {
			err = errors.New("spider.Fetch panic:" + fmt.Sprint(r))
		}
	}()

	dom := ""
	ctx, ctxCancel := chromedp.NewContext(getAllocCtx())
	defer ctxCancel()
	ctx, timeoutCancel := context.WithTimeout(ctx, spider.requestTimeout)
	defer timeoutCancel()
	err = chromedp.Run(ctx,
		chromedp.Navigate(u.String()),
		chromedp.WaitReady("html"), //disable images
		chromedp.OuterHTML("html", &dom, chromedp.ByQuery),
	)

	if err != nil {
		if err == context.Canceled {
			resetAllocCtx()
		} else {
			log.Println(err)
		}

		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
		if spider.FailureLevel == 0 {
			spider.FailureLevel = 10
		}
		summary.ErrType = err.Error()
	}
	//log.Println(helper.ByteCountBinary(uint64(len(dom))))
	summary.TrafficInStr = helper.ByteCountBinary(uint64(len(dom)))
	summary.TrafficIn = uint64(len(dom))
	spider.ResponseByte = []byte(dom)
	return summary, err
}

func (spider *Spider) HttpFetch(u *url.URL) (summary *Summary, err error) {
	spider.SetRequest(u) //setting spider.currentRequest

	if spider.RequestBefore != nil {
		spider.RequestBefore(spider)
	}

	//time
	spider.RequestStartTime = time.Now()
	spider.RequestEndTime = time.Time{} //empty

	spider.ResponseStatusCode = 0

	summary = &Summary{RawUrl: spider.currentRequest.URL.String(), AddTime: time.Now().Format("01-02 15:04:05"), TransportName: spider.TransportUrl.Host}

	defer func() {
		spider.RequestEndTime = time.Now()
		spider.ResponseStatusCode = summary.StatusCode
		summary.DnsLookupTime = spider.dnsLookupDuration.Truncate(time.Millisecond).String()
		summary.TLSTime = spider.tlsDuration.Truncate(time.Millisecond).String()

		if spider.FailureLevel == 0 && summary.StatusCode != 0 && summary.StatusCode != 200 {
			spider.FailureLevel = 10
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 2)
		}

		//A few times result of http request
		spider.recentFewTimesResult = append(spider.recentFewTimesResult, spider.FailureLevel == 0)

		//spider.TimeSlice = append(spider.TimeSlice[helper.MaxInt(len(spider.TimeSlice)-spider.TimeLenLimit, 0):], time.Since(spider.RequestStartTime))

		//recover
		if r := recover(); r != nil {
			err = errors.New("spider.Fetch panic:" + fmt.Sprint(r))
		}
	}()

	//traffic out
	dump, err := httputil.DumpRequestOut(spider.currentRequest, true)
	summary.ErrType = spider.requestErrorHandler(err)
	summary.TrafficOut = uint64(len(dump))

	// HTTP request
	spider.currentResponse, err = spider.httpClient().Do(spider.currentRequest)
	if err != nil {
		summary.ErrType = spider.requestErrorHandler(err)
		return summary, err
	}
	defer spider.currentResponse.Body.Close()

	summary.StatusCode = spider.currentResponse.StatusCode

	if spider.DownloadFilter != nil {
		filter, err := spider.DownloadFilter(spider, spider.currentResponse)

		if err != nil || !filter {
			//traffic  response header only
			dump, _ = httputil.DumpResponse(spider.currentResponse, false)
			summary.TrafficInStr = helper.ByteCountBinary(uint64(len(dump)))
			summary.TrafficIn = uint64(len(dump))

			if err != nil {
				summary.ErrType = "project.Filtered: " + err.Error()
				return summary, errors.New(summary.ErrType)
			}

			if !filter {
				return summary, nil
			}
		}
	}

	resByte, err := ioutil.ReadAll(spider.currentResponse.Body)
	summary.ErrType = spider.responseErrorHandler(err)
	if err != nil {
		return summary, err
	}

	//traffic in
	dump, err = httputil.DumpResponse(spider.currentResponse, false)
	summary.ErrType = spider.responseErrorHandler(err)
	summary.TrafficInStr = helper.ByteCountBinary(uint64(len(dump) + len(resByte)))
	summary.TrafficIn = uint64(len(dump) + len(resByte))

	//gzip decompression
	reader := ioutil.NopCloser(bytes.NewBuffer(resByte))
	if strings.ToLower(spider.currentResponse.Header.Get("Content-Encoding")) == "gzip" {
		reader, err = gzip.NewReader(reader)
		summary.ErrType = spider.responseErrorHandler(err)
	}

	res, err := ioutil.ReadAll(reader)
	reader.Close()
	summary.ErrType = spider.responseErrorHandler(err)
	if err != nil {
		return summary, err
	}

	spider.ResponseByte = res
	return summary, err
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
		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 2)
		return "x509.SystemRootsError"
	case *x509.UnknownAuthorityError:
		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 2)
		return "x509.UnknownAuthorityError"
	case *x509.HostnameError:
		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 2)
		return "x509.HostnameError"
	case *net.DNSConfigError:
		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 2)
		return "x509.DNSConfigError"
	case *net.DNSError:
		return "net.DNSError"
	case *net.OpError:
		log.Println("Request *net.OpError  "+spider.TransportUrl.Host+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "net.OpError"
	case net.Error:
		if err.(net.Error).Timeout() {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 4)
			return "net.Timeout"
		}
		if io.EOF == err {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "io.EOF"
		}
		if io.ErrUnexpectedEOF == err {
			return "io.ErrUnexpectedEOF"
		}
		if strings.Contains(err.Error(), "transport connection broken") {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 2)
			return "connection broken"
		}
		if strings.Contains(err.Error(), "unexpected EOF") {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "unexpected EOF"
		}
		if strings.Contains(err.Error(), "x509: certificate") {
			return "x509: certificate"
		}
		if strings.Contains(err.Error(), "no such host") {
			return "no such host"
		}
		if strings.Contains(err.Error(), ": EOF") {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "other EOF"
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "reset by peer"
		}
		// Get ..... :read ...
		if strings.Contains(strings.ToLower(err.Error()), "proxyconnectt") {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 4)
			return "proxyconnectt failed"
		}
		if _, ok := err.(*url.Error); ok {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "net.Error => url.Error"
		}
		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
		//log.Println("Request net.Error  "+spider.Transport.S.Host+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "unknown"
	case *url.Error:
		log.Println("Request Error "+spider.TransportUrl.Host+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "url.Error"
	default:
		if strings.HasPrefix(err.Error(), "invalid URL") {
			return "invalid URL"
		}
		if strings.HasPrefix(err.Error(), "no Host in request URL http") {
			return "no Host"
		}

		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
		spider.FailureLevel = 10
		// 2019/10/19 19:15:47 spider.go:414: Request Error 182.23.2.100:49833 *errors.errorString:  net/http: invalid header field value "https://book.douban.com/tag/to?                                                         start=160&type=S\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf6\x05\x00\x00\x00\x00\x00\x00\xfa\x05\x00\x00\x00\x00\x00\x00\xfc\x05" for key Referer https://book.douban.com/subject/          26328539/
		log.Println("Request Error "+spider.TransportUrl.Host+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
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
		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
		//2019/01/25 15:19:03 spider.go:431: Response *net.OpError  jp-b.mitsuha-node.com *net.OpError:  local error: tls: bad record MAC https://book.douban.com/subject/1836097/
		//log.Println("Response *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "net.OpError"
	case net.Error:
		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 4)
		if err.(net.Error).Timeout() {
			return "net.Timeout"
		}
		log.Println("Response net.Error  "+spider.TransportUrl.Host+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "net.Error"
	case *url.Error:
		log.Println("Response Error "+spider.TransportUrl.Host+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "url.Error"
	case tls.RecordHeaderError:
		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
		log.Println("Response Error "+spider.TransportUrl.Host+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "tls.RecordHeaderError"
	case flate.CorruptInputError:
		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
		log.Println("Response Error "+spider.TransportUrl.Host+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "flate.CorruptInputError"
	default:
		if strings.HasPrefix(err.Error(), "malformed chunked encoding") {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "chunked encoding"
		}
		if strings.HasPrefix(err.Error(), "invalid URL") {
			return "invalid URL"
		}
		if strings.HasPrefix(err.Error(), "http: unexpected EOF reading trailer") {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "unexpected EOF reading trailer"
		}
		if strings.HasPrefix(err.Error(), "http:  reading trailer") {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "http.reading trailer"
		}
		if gzip.ErrHeader == err {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 2)
			return "gzip.ErrHeader"
		}
		if gzip.ErrChecksum == err {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 2)
			return "gzip.ErrChecksum"
		}
		if io.EOF == err {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "io.EOF"
		}
		if io.ErrUnexpectedEOF == err {
			spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
			return "io.ErrUnexpectedEOF"
		}

		spider.EnqueueForFailure(spider, err, spider.currentRequest.URL.String(), 3)
		log.Println("Response Error "+spider.TransportUrl.Host+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "unknown"
	}
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
						value := strings.Trim(strings.TrimSpace(attr.Val), "\x00")
						if value == "" {
							continue
						}
						u, err := url.Parse(value)
						if err != nil {
							continue
						}
						u.Fragment = "" //remove anchor
						addUrl := spider.currentRequest.URL.ResolveReference(u)
						addUrl.Fragment = "" //remove anchor
						if addUrl.Scheme != "http" && addUrl.Scheme != "https" {
							continue
						}

						if len(addUrl.Hostname()) < 4 || strings.Index(addUrl.Hostname(), ".") == -1 || !helper.OnlyDomainCharacter(addUrl.Hostname()) { //todo !!!!!!!!!!!!!!! do not use regex in helper.OnlyDomainCharacter !!!!!!!!!!!!!!!!!!!!!!
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
}

func (spider *Spider) getHttpTrace() *httptrace.ClientTrace {

	var DnsLookupTime time.Time
	var tlsTime time.Time
	return &httptrace.ClientTrace{
		// GetConn is called before a connection is created or
		// retrieved from an idle pool. The hostPort is the
		// "host:port" of the target or proxy. GetConn is called even
		// if there's already an idle cached connection available.
		GetConn: func(hostPort string) {

		},

		// GotConn is called after a successful connection is
		// obtained. There is no hook for failure to obtain a
		// connection; instead, use the error from
		// Transport.RoundTrip.
		GotConn: func(connInfo httptrace.GotConnInfo) {
			//log.Printf("Got Conn: %+v\n", connInfo)
		},

		// PutIdleConn is called when the connection is returned to
		// the idle pool. If err is nil, the connection was
		// successfully returned to the idle pool. If err is non-nil,
		// it describes why not. PutIdleConn is not called if
		// connection reuse is disabled via Transport.DisableKeepAlives.
		// PutIdleConn is called before the caller's Response.Body.Close
		// call returns.
		// For HTTP/2, this hook is not currently used.
		PutIdleConn: func(err error) {

		},

		// GotFirstResponseByte is called when the first byte of the response
		// headers is available.
		GotFirstResponseByte: func() {

		},

		// Got100Continue is called if the server replies with a "100
		// Continue" response.
		Got100Continue: func() {

		},

		// Got1xxResponse is called for each 1xx informational response header
		// returned before the final non-1xx response. Got1xxResponse is called
		// for "100 Continue" responses, even if Got100Continue is also defined.
		// If it returns an error, the client request is aborted with that error value.
		Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
			return nil
		}, // Go 1.11

		// DNSStart is called when a DNS lookup begins.
		DNSStart: func(httptrace.DNSStartInfo) {
			DnsLookupTime = time.Now()
		},

		// DNSDone is called when a DNS lookup ends.
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			spider.dnsLookupDuration = time.Since(DnsLookupTime)
			//log.Println(spider.CurrentRequest.URL.String(), dnsInfo)
		},

		// ConnectStart is called when a new connection's Dial begins.
		// If net.Dialer.DualStack (IPv6 "Happy Eyeballs") support is
		// enabled, this may be called multiple times.
		ConnectStart: func(network, addr string) {

		},

		// ConnectDone is called when a new connection's Dial
		// completes. The provided err indicates whether the
		// connection completedly successfully.
		// If net.Dialer.DualStack ("Happy Eyeballs") support is
		// enabled, this may be called multiple times.
		ConnectDone: func(network, addr string, err error) {

		},

		// TLSHandshakeStart is called when the TLS handshake is started. When
		// connecting to a HTTPS site via a HTTP proxy, the handshake happens after
		// the CONNECT request is processed by the proxy.
		TLSHandshakeStart: func() {
			tlsTime = time.Now()
		}, // Go 1.8

		// TLSHandshakeDone is called after the TLS handshake with either the
		// successful handshake's connection state, or a non-nil error on handshake
		// failure.
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			spider.tlsDuration = time.Since(tlsTime)
		}, // Go 1.8

		// WroteHeaderField is called after the Transport has written
		// each request header. At the time of this call the values
		// might be buffered and not yet written to the network.
		WroteHeaderField: func(key string, value []string) {

		}, // Go 1.11

		// WroteHeaders is called after the Transport has written
		// all request headers.
		WroteHeaders: func() {

		},

		// Wait100Continue is called if the Request specified
		// "Expected: 100-continue" and the Transport has written the
		// request headers but is waiting for "100 Continue" from the
		// server before writing the request body.
		Wait100Continue: func() {

		},

		// WroteRequest is called with the result of writing the
		// request and any body. It may be called multiple times
		// in the case of retried requests.
		WroteRequest: func(httptrace.WroteRequestInfo) {

		},
	}

}
