package spider

import (
	"asuka/helper"
	"asuka/proxy"
	"asuka/queue"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"container/list"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
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
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	//Emergency error handling
	go func() {
		t := time.NewTicker(time.Second)
		for {
			<-t.C
			failCount := 0
			for _, spider := range spiderList {
				if spider.FailureLevel != 0 && len(spider.Transport.RecentFewTimesResultEmergency) >= RecentSeveralTimesResultCap {
					failCount++
				}
			}
			if float64(failCount)/float64(len(spiderList)) >= 0.35 {
				for _, spider := range spiderList {
					spider.Transport.RecentFewTimesResult = make([]bool, 0, RecentSeveralTimesResultCap)
					spider.Transport.RecentFewTimesResultEmergency = make([]bool, 0, RecentSeveralTimesResultCap)
					spider.FailureLevel = 150
					spider.ResetSleep()
					spider.AddSleep(time.Hour * 3)
				}
			}
		}
	}()
}

const RecentFetchCount = 100
const RecentSeveralTimesResultCap = 7

var RecentFetchMutex = &sync.Mutex{}
var RecentFetchLastIndex int64 = 0
var RecentFetchList []*RecentFetch
var spiderList []*Spider

type RecentFetch struct {
	Index         int64
	TransportName string
	StatusCode    int // http response status code
	RawUrl        string
	ConsumeTime   string
	AddTime       string
	ErrType       string
	ResponseSize  string
	ContentType   string
}

type Spider struct {
	Transport *proxy.Transport
	Client    *http.Client
	Queue     *queue.Queue

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
	sleepDuration    time.Duration
}

func New(t *proxy.Transport, j *cookiejar.Jar, queue *queue.Queue) *Spider {
	if j == nil {
		j, _ = cookiejar.New(nil)
	}
	c := &http.Client{Transport: t.T, Jar: j, Timeout: time.Second * 30}
	spider := &Spider{Queue: queue, Transport: t, Client: c, RequestsMap: map[string]*http.Request{}, TimeList: list.New(), TimeLenLimit: 10}
	spiderList = append(spiderList, spider)
	return spider
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

func (spider *Spider) Throttle() {
	if spider.Transport.S.Interval > .0 {
		spider.AddSleep(time.Duration(spider.Transport.S.Interval * 1e9))
	}
	for {
		if !spider.Stop {
			break
		}
		time.Sleep(3e9)
	}

	if spider.FailureLevel > 0 {
		spider.AddSleep(time.Second)
	}
	if spider.FailureLevel > 1 {
		spider.AddSleep(time.Minute)
	}

	//Failure control
	if len(spider.Transport.RecentFewTimesResult) >= RecentSeveralTimesResultCap {
		spider.Transport.RecentFewTimesResult = spider.Transport.RecentFewTimesResult[len(spider.Transport.RecentFewTimesResult)-RecentSeveralTimesResultCap:]
		failCount := 0
		for _, v := range spider.Transport.RecentFewTimesResult {
			if !v {
				failCount++
			}
		}
		if float64(failCount)/float64(RecentSeveralTimesResultCap) >= 0.4 {
			spider.Transport.RecentFewTimesResult = make([]bool, 0, RecentSeveralTimesResultCap)
			accessCountAll := spider.Transport.GetAccessCount()
			failureCountAll := spider.Transport.GetFailureCount()
			failureRateAll := helper.SpiderFailureRate(accessCountAll, failureCountAll)
			if accessCountAll > 40 && failureRateAll > 95 && spider.FailureLevel <= 100 {
				spider.FailureLevel = 100
				spider.AddSleep(time.Hour * 2)
			} else if accessCountAll > 40 && failureRateAll > 85 && spider.FailureLevel <= 80 {
				spider.FailureLevel = 80
				spider.AddSleep(time.Minute * 40)
			} else if accessCountAll > 30 && failureRateAll > 70 && spider.FailureLevel <= 60 {
				spider.FailureLevel = 60
				spider.AddSleep(time.Minute * 10)
			} else if accessCountAll > 30 && failureRateAll > 60 && spider.FailureLevel <= 30 {
				spider.FailureLevel = 40
				spider.AddSleep(time.Minute * 5)
			} else if spider.FailureLevel <= 20 {
				spider.FailureLevel = 20
				spider.AddSleep(time.Minute * 2)
			}
		}
	}

	//go to sleep and reset sleep duration
	if duration := spider.GetSleep(); duration > 0 {
		time.Sleep(duration)
	}
	spider.ResetSleep()
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
		r.URL = url
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

func (spider *Spider) Fetch(u *url.URL, requestBefore func(spider *Spider), downloadFilter func(spider *Spider, response *http.Response) (bool, error)) (resp *http.Response, err error) {
	spider.SetRequest(u, nil) //setting spider.CurrentRequest

	if requestBefore != nil {
		requestBefore(spider)
	}

	spider.ResponseStr = ""
	spider.ResponseByte = []byte{}

	//time
	spider.RequestStartTime = time.Now()

	recentFetch := &RecentFetch{RawUrl: spider.CurrentRequest.URL.String(), AddTime: time.Now().Format("01-02 15:04:05"), TransportName: spider.Transport.S.Name}

	spider.Transport.AddAccess(spider.CurrentRequest.URL.String())

	defer func() {
		if err != nil {
			spider.Transport.AddFailure(spider.CurrentRequest.URL.String())
		}

		//A few times result of http request
		spider.Transport.RecentFewTimesResult = append(spider.Transport.RecentFewTimesResult, spider.FailureLevel == 0)
		spider.Transport.RecentFewTimesResultEmergency = append(spider.Transport.RecentFewTimesResultEmergency, spider.FailureLevel == 0)

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
	recentFetch.ContentType = resp.Header.Get("Content-type")

	if downloadFilter != nil {
		filter, err := downloadFilter(spider, resp)

		if err != nil || !filter {
			//traffic  response header only
			dump, _ = httputil.DumpResponse(resp, false)
			recentFetch.ResponseSize = helper.ByteCountBinary(uint64(len(dump)))
			spider.Transport.TrafficIn += uint64(len(dump))

			if err != nil {
				recentFetch.ErrType = "project.Filtered"
				return nil, errors.New(recentFetch.ErrType)
			}

			if !filter {
				return nil, nil
			}
		}
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
		return "net.DNSError"
	case *net.OpError:
		log.Println("Request *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "net.OpError"
	case net.Error:
		if err.(net.Error).Timeout() {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 4)
			return "net.Timeout"
		}
		if io.EOF == err {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
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
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return "unexpected EOF"
		}
		if strings.Contains(err.Error(), "x509: certificate") {
			return "x509: certificate"
		}
		if strings.Contains(err.Error(), "no such host") {
			return "no such host"
		}
		if strings.Contains(err.Error(), ": EOF") {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return "other EOF"
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return "reset by peer"
		}
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
		// Get ..... :read ...
		//log.Println("Request net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
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

		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
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
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 1)
		log.Println("Response *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "net.OpError"
	case net.Error:
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 4)
		if err.(net.Error).Timeout() {
			return "net.Timeout"
		}
		log.Println("Response net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "net.Error"
	case *url.Error:
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "url.Error"
	case tls.RecordHeaderError:
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 1)
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "tls.RecordHeaderError"
	case flate.CorruptInputError:
		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 1)
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.CurrentRequest.URL.String())
		return "flate.CorruptInputError"
	default:
		if strings.HasPrefix(err.Error(), "malformed chunked encoding") {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 2)
			return "chunked encoding"
		}
		if strings.HasPrefix(err.Error(), "invalid URL") {
			return "invalid URL"
		}
		if strings.HasPrefix(err.Error(), "http: unexpected EOF reading trailer") {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return "unexpected EOF reading trailer"
		}
		if strings.HasPrefix(err.Error(), "http:  reading trailer") {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return "http.reading trailer"
		}
		if gzip.ErrHeader == err {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 1)
			return "gzip.ErrHeader"
		}
		if gzip.ErrChecksum == err {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 1)
			return "gzip.ErrChecksum"
		}
		if io.EOF == err {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return "io.EOF"
		}
		if io.ErrUnexpectedEOF == err {
			spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
			return "io.ErrUnexpectedEOF"
		}

		spider.Queue.EnqueueForFailure(spider.CurrentRequest.URL.String(), 3)
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
						u.Fragment = "" //remove anchor
						addUrl := spider.CurrentRequest.URL.ResolveReference(u)
						addUrl.Fragment = "" //remove anchor
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
}
