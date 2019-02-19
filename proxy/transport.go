package proxy

import (
	"context"
	"github.com/chenset/asuka/helper"
	"golang.org/x/net/proxy"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const SecondInterval = 1
const MinuteInterval = 40 * 60
const CountQueueSecondCap = MinuteInterval * 2

var CountQueueMinuteCap = 1 //initial value, will dynamic changes

type Transport struct {
	S               *SsAddr
	t               http.RoundTripper
	transportClosed bool

	countSliceMutex         sync.RWMutex
	CountSliceCursor        int
	AccessCountSecondSlice  []uint32
	FailureCountSecondSlice []uint32
	AccessCountMinuteSlice  []uint32
	FailureCountMinuteSlice []uint32

	AccessCountHistory  int
	FailureCountHistory int

	//traffic size
	TrafficIn  uint64
	TrafficOut uint64

	Ping            time.Duration
	PingFailureRate float64

	RecentFewTimesResult []bool
}

func NewTransport(ssAddr *SsAddr) (*Transport, error) {
	instance := &Transport{S: ssAddr, t: createHttpTransport(ssAddr)}
	return instance, nil
}

func createHttpTransport(SockInfo *SsAddr) *http.Transport {
	t := &http.Transport{
		MaxIdleConnsPerHost:   2,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	switch SockInfo.Type {
	case "local":
		t.Proxy = nil //disable system proxy
		t.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
	case "http", "https":
		proxyURL, err := url.Parse(SockInfo.Type + "://" + SockInfo.ServerAddr)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		t.Proxy = http.ProxyURL(proxyURL) // with http proxy
		t.TLSHandshakeTimeout = time.Minute
		t.DialContext = (&net.Dialer{
			Timeout:   time.Minute,
			KeepAlive: time.Minute,
			DualStack: true,
		}).DialContext
	case "ss", "ssr":
		SockInfo.WaitUntilConnected() //waiting
		dialer, err := proxy.SOCKS5("tcp", SockInfo.ClientAddr, nil, proxy.Direct)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		t.Proxy = nil //disable system proxy
		t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}
	return t
}

// AddAccess 每次调用请求时增加一次记录, 无论是否成功
func (transport *Transport) AddAccess() {
	transport.AccessCountHistory++
}

// AddFailure 每次调用请求并失败时增加一次失败记录
func (transport *Transport) AddFailure() {
	transport.FailureCountHistory++
}

func (transport *Transport) GetAccessCount() int {
	return transport.AccessCountHistory
}

func (transport *Transport) GetFailureCount() int {
	return transport.FailureCountHistory
}

//updateCountQueueCap
func updateCountQueueCap(second int) {
	if CountQueueMinuteCap <= second/MinuteInterval {
		CountQueueMinuteCap = second/MinuteInterval + 1
	}
}

// LoadRate 获取指定秒数内的负载值.参数最小值SecondInterval秒
func (transport *Transport) LoadRate(second int) (rate float64) {
	//Read lock
	defer func() {
		transport.countSliceMutex.RUnlock()
	}()
	transport.countSliceMutex.RLock()

	updateCountQueueCap(second)

	sliceLen := len(transport.AccessCountSecondSlice)
	if sliceLen == 0 || second == 0 {
		return
	}

	times := int(math.Ceil(float64(second) / SecondInterval))

	//SecondInterval
	if times <= CountQueueSecondCap {
		return float64(transport.AccessCountSecondSlice[sliceLen-1]-transport.AccessCountSecondSlice[helper.MaxInt(sliceLen-times-1, 0)]) / float64(times)
	}

	minuteSliceLen := len(transport.AccessCountMinuteSlice)
	minSecond := helper.MinInt(second, minuteSliceLen*MinuteInterval+transport.CountSliceCursor%(MinuteInterval))
	realTimeSecond := minSecond%(MinuteInterval) + MinuteInterval
	rate += float64(transport.AccessCountSecondSlice[sliceLen-1] - transport.AccessCountSecondSlice[helper.MaxInt(sliceLen-realTimeSecond-1, 0)])
	minSecond -= realTimeSecond
	if minSecond > 0 {
		rate += float64(transport.AccessCountMinuteSlice[minuteSliceLen-1] - transport.AccessCountMinuteSlice[minuteSliceLen-minSecond/(MinuteInterval)-1])
	}
	return rate / float64(times)
}

//AccessCount  获取指定秒数内的访问数/失败j数量.参数最小值SecondInterval秒
func (transport *Transport) AccessCount(second int) (accessTimes, failureTimes int) {
	//Read lock
	defer func() {
		transport.countSliceMutex.RUnlock()
	}()
	transport.countSliceMutex.RLock()

	updateCountQueueCap(second)

	accessSliceLen := len(transport.AccessCountSecondSlice)
	failureSliceLen := len(transport.FailureCountSecondSlice)

	if accessSliceLen+failureSliceLen == 0 {
		return
	}

	times := int(math.Ceil(float64(second) / SecondInterval))
	if times == 0 {
		return
	}

	if times <= CountQueueSecondCap {
		if accessSliceLen != 0 {
			accessTimes = int(transport.AccessCountSecondSlice[accessSliceLen-1] - transport.AccessCountSecondSlice[helper.MaxInt(accessSliceLen-times-1, 0)])
		}

		if failureSliceLen != 0 {
			failureTimes = int(transport.FailureCountSecondSlice[failureSliceLen-1] - transport.FailureCountSecondSlice[helper.MaxInt(failureSliceLen-times-1, 0)])
		}
		return
	}

	minuteSliceLen := len(transport.AccessCountMinuteSlice) //len(transport.FailureCountMinuteSlice)/
	minSecond := helper.MinInt(second, minuteSliceLen*MinuteInterval+transport.CountSliceCursor%(MinuteInterval))
	realTimeSecond := minSecond%(MinuteInterval) + MinuteInterval
	minSecond -= realTimeSecond

	accessTimes += int(transport.AccessCountSecondSlice[accessSliceLen-1] - transport.AccessCountSecondSlice[helper.MaxInt(accessSliceLen-realTimeSecond-1, 0)])
	if minSecond > 0 {
		accessTimes += int(transport.AccessCountMinuteSlice[minuteSliceLen-1] - transport.AccessCountMinuteSlice[minuteSliceLen-minSecond/(MinuteInterval)-1])
	}

	failureTimes += int(transport.FailureCountSecondSlice[accessSliceLen-1] - transport.FailureCountSecondSlice[helper.MaxInt(accessSliceLen-realTimeSecond-1, 0)])
	if minSecond > 0 {
		failureTimes += int(transport.FailureCountMinuteSlice[minuteSliceLen-1] - transport.FailureCountMinuteSlice[minuteSliceLen-minSecond/(MinuteInterval)-1])
	}

	return
}

func (transport *Transport) RecordAccessSecondCount() {
	//Write lock
	defer func() {
		transport.countSliceMutex.Unlock()
	}()
	transport.countSliceMutex.Lock()
	//slice fifo
	transport.AccessCountSecondSlice = append(transport.AccessCountSecondSlice[helper.MaxInt(len(transport.AccessCountSecondSlice)-CountQueueSecondCap, 0):], uint32(transport.GetAccessCount()))

	if transport.CountSliceCursor%MinuteInterval == 0 {
		transport.AccessCountMinuteSlice = append(transport.AccessCountMinuteSlice[helper.MaxInt(len(transport.AccessCountMinuteSlice)-CountQueueMinuteCap, 0):], transport.AccessCountSecondSlice[len(transport.AccessCountSecondSlice)-MinuteInterval])
	}
}

func (transport *Transport) RecordFailureSecondCount() {
	//Write lock
	defer func() {
		transport.countSliceMutex.Unlock()
	}()
	transport.countSliceMutex.Lock()

	//slice fifo
	transport.FailureCountSecondSlice = append(transport.FailureCountSecondSlice[helper.MaxInt(len(transport.FailureCountSecondSlice)-CountQueueSecondCap, 0):], uint32(transport.GetFailureCount()))

	if transport.CountSliceCursor%MinuteInterval == 0 {
		transport.FailureCountMinuteSlice = append(transport.FailureCountMinuteSlice[helper.MaxInt(len(transport.FailureCountMinuteSlice)-CountQueueMinuteCap, 0):], transport.FailureCountSecondSlice[len(transport.FailureCountSecondSlice)-MinuteInterval])
	}
}

func (transport *Transport) Close() {
	if !transport.transportClosed {
		if transport.S.Type == "ss" || transport.S.Type == "ssr" {
			transport.S.Close()
		}

		transport.t.(*http.Transport).DisableKeepAlives = true
		transport.t.(*http.Transport).CloseIdleConnections()
		transport.transportClosed = true
	}
}

func (transport *Transport) Connect() *http.Transport {
	if transport.transportClosed {
		transport.t = createHttpTransport(transport.S)
		transport.transportClosed = false
	}

	return transport.t.(*http.Transport)
}
