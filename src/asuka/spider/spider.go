package spider

import (
	"asuka/helper"
	"asuka/proxy"
	"asuka/queue"
	"bytes"
	"compress/flate"
	"compress/gzip"
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
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const PeriodOfFailureSecond = 86400 / 2

type Summary struct {
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
	client    *http.Client
	Queue     *queue.Queue

	RequestsMap    map[string]*http.Request
	currentRequest *http.Request

	ResponseStr  string
	ResponseByte []byte

	TimeSlice    []time.Duration
	TimeLenLimit int

	FailureLevel int

	StartTime        time.Time
	RequestStartTime time.Time
	Stop             bool
	SleepDuration    time.Duration

	RequestBefore   func(spider *Spider)
	DownloadFilter  func(spider *Spider, response *http.Response) (bool, error)
	ProjectThrottle func(spider *Spider)

	httpTrace                   *httptrace.ClientTrace
	RecentSeveralTimesResultCap int

	Test int
}

func New(t *proxy.Transport, queue *queue.Queue) *Spider {
	spider := &Spider{Queue: queue, Transport: t, RequestsMap: map[string]*http.Request{}, TimeLenLimit: 10, StartTime: time.Now(), RecentSeveralTimesResultCap: 5}
	//spider.updateClient()
	spider.registerHttpTrace()
	return spider
}

func (spider *Spider) CurrentRequest() *http.Request {
	return spider.currentRequest
}

func (spider *Spider) Client() *http.Client {
	spider.setClient()
	return spider.client
}

func (spider *Spider) setClient() {
	if spider.client == nil || spider.client.Transport.(*http.Transport) != spider.Transport.Connect() {
		j, _ := cookiejar.New(nil)
		spider.client = &http.Client{Transport: spider.Transport.Connect(), Jar: j, Timeout: time.Second * 30}
		//fmt.Println("transport changed")
	}
}

func (spider *Spider) AddSleep(duration time.Duration) {
	spider.SleepDuration += duration
}

func (spider *Spider) GetSleep() time.Duration {
	return spider.SleepDuration
}

func (spider *Spider) ResetSleep() {
	spider.SleepDuration = 0
}

func (spider *Spider) Throttle(dispatcherCallback func(spider *Spider)) {
	if spider.Transport.S.Interval > .0 {
		spider.AddSleep(time.Duration(spider.Transport.S.Interval * 1e9))
	}
	for {
		if !spider.Stop {
			break
		}
		spider.Transport.Close()
		time.Sleep(3e9)
	}

	if spider.FailureLevel > 0 {
		spider.AddSleep(time.Second)
	}
	if spider.FailureLevel > 1 {
		spider.AddSleep(time.Second * 30)
	}

	//Failure control
	if len(spider.Transport.RecentFewTimesResult) >= spider.RecentSeveralTimesResultCap {
		spider.Transport.RecentFewTimesResult = spider.Transport.RecentFewTimesResult[len(spider.Transport.RecentFewTimesResult)-spider.RecentSeveralTimesResultCap:]
		failCount := 0
		for _, v := range spider.Transport.RecentFewTimesResult {
			if !v {
				failCount++
			}
		}
		if float64(failCount)/float64(spider.RecentSeveralTimesResultCap) >= 0.4 {
			spider.Transport.RecentFewTimesResult = make([]bool, 0, spider.RecentSeveralTimesResultCap)

			accessCountAll, failureCountAll := spider.Transport.AccessCount(helper.MinInt(int(time.Since(spider.StartTime).Seconds()), PeriodOfFailureSecond))
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

// setRequest http.Request 是维持session会话的关键之一. 这里是在管理http.Request, 保证每个url能找到对应之前的http.Request
func (spider *Spider) SetRequest(url *url.URL, header *http.Header) *Spider {

	tld, err := helper.TldDomain(url)
	if err != nil {
		tld = "DefaultRequest"
	}

	r, ok := spider.RequestsMap[tld]
	if ok {
		r.URL = url
		spider.currentRequest = r
	} else {
		r, err = http.NewRequest("GET", url.String(), nil)
		if err != nil {
			log.Fatal(err)
		}

		//Accept-Encoding: gzip
		if r.Header.Get("Accept-Encoding") == "" {
			r.Header.Set("Accept-Encoding", "gzip")
		}

		spider.currentRequest = r
		spider.RequestsMap[tld] = r
	}

	//spider.currentRequest.Close = true // prevents re-use of TCP connections between requests to the same hosts

	if header != nil {
		for k := range *header {
			spider.currentRequest.Header.Set(k, header.Get(k))
		}
	}

	if spider.currentRequest.UserAgent() == "" {
		spider.currentRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/"+strconv.FormatFloat(rand.Float64()*10000, 'f', 3, 64)+" (KHTML, like Gecko) Chrome/71.0."+strconv.FormatFloat(rand.Float64()*10000, 'f', 3, 64)+" Safari/537.36")
	}

	spider.currentRequest = spider.currentRequest.WithContext(httptrace.WithClientTrace(spider.currentRequest.Context(), spider.httpTrace))
	return spider
}

func (spider *Spider) Fetch(u *url.URL) (resp *http.Response, summary *Summary, err error) {
	spider.SetRequest(u, nil) //setting spider.currentRequest

	if spider.RequestBefore != nil {
		spider.RequestBefore(spider)
	}

	spider.ResponseStr = ""
	spider.ResponseByte = []byte{}

	//time
	spider.RequestStartTime = time.Now()

	summary = &Summary{RawUrl: spider.currentRequest.URL.String(), AddTime: time.Now().Format("01-02 15:04:05"), TransportName: spider.Transport.S.Name}

	spider.Transport.AddAccess(spider.currentRequest.URL.String())

	defer func() {
		if err != nil {
			spider.Transport.AddFailure(spider.currentRequest.URL.String())
		}

		if spider.FailureLevel == 0 && summary.StatusCode != 0 && summary.StatusCode != 200 {
			spider.FailureLevel = 10
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 2)
		}

		//A few times result of http request
		spider.Transport.RecentFewTimesResult = append(spider.Transport.RecentFewTimesResult, spider.FailureLevel == 0)

		spider.TimeSlice = append(spider.TimeSlice[helper.MaxInt(len(spider.TimeSlice)-spider.TimeLenLimit, 0):], time.Since(spider.RequestStartTime))

		//recover
		if r := recover(); r != nil {
			spider.Transport.AddFailure(spider.currentRequest.URL.String())
			err = errors.New("spider.Fetch panic:" + fmt.Sprint(r))
		}
	}()

	//traffic
	dump, err := httputil.DumpRequestOut(spider.currentRequest, true)
	summary.ErrType = spider.requestErrorHandler(err)
	spider.Transport.TrafficOut += uint64(len(dump))
	//for localhost
	if spider.Transport.S.Type == "" {
		spider.Transport.S.TrafficOut += uint64(len(dump))
	}

	resp, err = spider.client.Do(spider.currentRequest)
	if err != nil {
		summary.ErrType = spider.requestErrorHandler(err)
		return resp, summary, err
	}
	defer resp.Body.Close()
	summary.StatusCode = resp.StatusCode
	summary.ContentType = resp.Header.Get("Content-type")

	if spider.DownloadFilter != nil {
		filter, err := spider.DownloadFilter(spider, resp)

		if err != nil || !filter {
			//traffic  response header only
			dump, _ = httputil.DumpResponse(resp, false)
			summary.ResponseSize = helper.ByteCountBinary(uint64(len(dump)))
			spider.Transport.TrafficIn += uint64(len(dump))
			//for localhost
			if spider.Transport.S.Type == "" {
				spider.Transport.S.TrafficIn += uint64(len(dump))
			}

			if err != nil {
				summary.ErrType = "project.Filtered"
				return nil, summary, errors.New(summary.ErrType)
			}

			if !filter {
				return nil, summary, nil
			}
		}
	}

	resByte, err := ioutil.ReadAll(resp.Body)
	summary.ErrType = spider.responseErrorHandler(err)
	if err != nil {
		return resp, summary, err
	}

	//traffic
	dump, err = httputil.DumpResponse(resp, false)
	summary.ErrType = spider.responseErrorHandler(err)
	summary.ResponseSize = helper.ByteCountBinary(uint64(len(dump) + len(resByte)))
	spider.Transport.TrafficIn += uint64(len(dump) + len(resByte))
	//for localhost
	if spider.Transport.S.Type == "" {
		spider.Transport.S.TrafficIn += uint64(len(dump) + len(resByte))
	}

	//gzip decompression
	reader := ioutil.NopCloser(bytes.NewBuffer(resByte))
	defer reader.Close()
	if strings.ToLower(resp.Header.Get("Content-Encoding")) == "gzip" {
		reader, err = gzip.NewReader(reader)
		summary.ErrType = spider.responseErrorHandler(err)
	}

	res, err := ioutil.ReadAll(reader)
	summary.ErrType = spider.responseErrorHandler(err)
	if err != nil {
		return resp, summary, err
	}

	//http status
	if resp.StatusCode != 200 && err == nil {
		spider.Transport.AddFailure(spider.currentRequest.URL.String())
	}

	spider.ResponseStr = string(res[:])
	spider.ResponseByte = res
	return resp, summary, err
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
		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 2)
		return "x509.SystemRootsError"
	case *x509.UnknownAuthorityError:
		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 2)
		return "x509.UnknownAuthorityError"
	case *x509.HostnameError:
		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 2)
		return "x509.HostnameError"
	case *net.DNSConfigError:
		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 2)
		return "x509.DNSConfigError"
	case *net.DNSError:
		return "net.DNSError"
	case *net.OpError:
		log.Println("Request *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "net.OpError"
	case net.Error:
		if err.(net.Error).Timeout() {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 4)
			return "net.Timeout"
		}
		if io.EOF == err {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
			return "io.EOF"
		}
		if io.ErrUnexpectedEOF == err {
			return "io.ErrUnexpectedEOF"
		}
		if strings.Contains(err.Error(), "transport connection broken") {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 2)
			return "connection broken"
		}
		if strings.Contains(err.Error(), "unexpected EOF") {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
			return "unexpected EOF"
		}
		if strings.Contains(err.Error(), "x509: certificate") {
			return "x509: certificate"
		}
		if strings.Contains(err.Error(), "no such host") {
			return "no such host"
		}
		if strings.Contains(err.Error(), ": EOF") {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
			return "other EOF"
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
			return "reset by peer"
		}
		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
		// Get ..... :read ...
		//log.Println("Request net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "unknown"
	case *url.Error:
		log.Println("Request Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "url.Error"
	default:
		if strings.HasPrefix(err.Error(), "invalid URL") {
			return "invalid URL"
		}
		if strings.HasPrefix(err.Error(), "no Host in request URL http") {
			return "no Host"
		}

		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
		spider.FailureLevel = 10
		log.Println("Request Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
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
		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
		//2019/01/25 15:19:03 spider.go:431: Response *net.OpError  jp-b.mitsuha-node.com *net.OpError:  local error: tls: bad record MAC https://book.douban.com/subject/1836097/
		log.Println("Response *net.OpError  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "net.OpError"
	case net.Error:
		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 4)
		if err.(net.Error).Timeout() {
			return "net.Timeout"
		}
		log.Println("Response net.Error  "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "net.Error"
	case *url.Error:
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "url.Error"
	case tls.RecordHeaderError:
		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "tls.RecordHeaderError"
	case flate.CorruptInputError:
		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "flate.CorruptInputError"
	default:
		if strings.HasPrefix(err.Error(), "malformed chunked encoding") {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
			return "chunked encoding"
		}
		if strings.HasPrefix(err.Error(), "invalid URL") {
			return "invalid URL"
		}
		if strings.HasPrefix(err.Error(), "http: unexpected EOF reading trailer") {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
			return "unexpected EOF reading trailer"
		}
		if strings.HasPrefix(err.Error(), "http:  reading trailer") {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
			return "http.reading trailer"
		}
		if gzip.ErrHeader == err {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 2)
			return "gzip.ErrHeader"
		}
		if gzip.ErrChecksum == err {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 2)
			return "gzip.ErrChecksum"
		}
		if io.EOF == err {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
			return "io.EOF"
		}
		if io.ErrUnexpectedEOF == err {
			spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
			return "io.ErrUnexpectedEOF"
		}

		spider.Queue.EnqueueForFailure(spider.currentRequest.URL.String(), 3)
		log.Println("Response Error "+spider.Transport.S.Name+" "+reflect.TypeOf(err).String()+": ", err, spider.currentRequest.URL.String())
		return "unknown"
	}
}

func (spider *Spider) GetAvgTime() (t time.Duration) {
	for _, tt := range spider.TimeSlice {
		t += tt
	}

	if len(spider.TimeSlice) == 0 {
		return
	}
	t /= time.Duration(len(spider.TimeSlice))
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
						addUrl := spider.currentRequest.URL.ResolveReference(u)
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

func (spider *Spider) registerHttpTrace() {
	spider.httpTrace = &httptrace.ClientTrace{
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
		//Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
		//	return nil
		//}, // Go 1.11

		// DNSStart is called when a DNS lookup begins.
		DNSStart: func(httptrace.DNSStartInfo) {

		},

		// DNSDone is called when a DNS lookup ends.
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			//log.Println(spider.currentRequest.URL.String(), dnsInfo)
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

		}, // Go 1.8

		// TLSHandshakeDone is called after the TLS handshake with either the
		// successful handshake's connection state, or a non-nil error on handshake
		// failure.
		TLSHandshakeDone: func(tls.ConnectionState, error) {

		}, // Go 1.8

		// WroteHeaderField is called after the Transport has written
		// each request header. At the time of this call the values
		// might be buffered and not yet written to the network.
		//WroteHeaderField: func(key string, value []string) {
		//
		//}, // Go 1.11

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
