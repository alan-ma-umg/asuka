package proxy

import (
	"asuka/helper"
	"context"
	"golang.org/x/net/proxy"
	"log"
	"math"
	"net"
	"net/http"
	"sync"
	"time"
)

const SecondInterval = 1
const MinuteInterval = 10

var CountQueueLen = 0 //Dynamic changes

func init() {
	//save
	go func() {
		t := time.NewTicker(time.Second * SecondInterval)
		for {
			<-t.C
			for _, t := range transportList {
				t.recordAccessCount()
				t.recordFailureCount()

			}
		}
	}()
}

var transportList []*Transport

type Transport struct {
	S *SsAddr
	T http.RoundTripper

	countSliceMutex         sync.RWMutex
	accessCountSecondSlice        []int
	failureCountSecondSlice       []int
	accessCountMinuteSlice  []int
	failureCountMinuteSlice []int

	accessCountHistory  int
	FailureCountHistory int

	LoopCount int

	//traffic size
	TrafficIn  uint64
	TrafficOut uint64

	Ping            time.Duration
	PingFailureRate float64

	RecentFewTimesResult          []bool
	RecentFewTimesResultEmergency []bool
}

func NewTransport(ssAddr *SsAddr) (*Transport, error) {
	instance := &Transport{S: ssAddr, T: createHttpTransport(ssAddr), LoopCount: 0}
	transportList = append(transportList, instance)
	return instance, nil
}

func createHttpTransport(SockInfo *SsAddr) *http.Transport {
	t := &http.Transport{
		Proxy:                 nil, //disable system proxy
		MaxIdleConnsPerHost:   2,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if SockInfo.ServerAddr == "" { //no socks5 proxy

		t.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
	} else { //use socks5 proxy
		dialer, err := proxy.SOCKS5("tcp", SockInfo.ClientAddr, nil, proxy.Direct)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}
	return t
}

// AddAccess 每次调用请求时增加一次记录, 无论是否成功
func (transport *Transport) AddAccess(link string) {
	transport.accessCountHistory++
}

// AddFailure 每次调用请求并失败时增加一次失败记录
func (transport *Transport) AddFailure(link string) {
	transport.FailureCountHistory++
}

func (transport *Transport) GetAccessCount() int {
	return transport.accessCountHistory
}

func (transport *Transport) GetFailureCount() int {
	return transport.FailureCountHistory
}

func updateCountQueueLen(second int) {
	if CountQueueLen <= second/SecondInterval {
		CountQueueLen = second/SecondInterval + 2
	}
}

// LoadRate 获取指定秒数内的负载值.参数最小值SecondInterval秒
func (transport *Transport) LoadRate(second int) float64 {
	//Read lock
	defer func() {
		transport.countSliceMutex.RUnlock()
	}()
	transport.countSliceMutex.RLock()

	updateCountQueueLen(second)
	rate := 0.0

	sliceLen := len(transport.accessCountSecondSlice)

	times := int(math.Ceil(float64(second) / SecondInterval))
	for i := 0; i < times; i++ {
		currentNum := 0
		prevNum := 0
		if i < sliceLen {
			currentNum = transport.accessCountSecondSlice[sliceLen-i-1]
			if i+1 < sliceLen {
				prevNum = transport.accessCountSecondSlice[sliceLen-i-2]
			}
		}
		rate += float64(currentNum-prevNum) / SecondInterval
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

	updateCountQueueLen(second)

	accessSliceLen := len(transport.accessCountSecondSlice)
	failureSliceLen := len(transport.failureCountSecondSlice)

	times := int(math.Ceil(float64(second) / SecondInterval))
	for i := 0; i < times; i++ {
		currentFailureNum := 0
		prevFailureNum := 0
		if i < failureSliceLen {
			currentFailureNum = transport.failureCountSecondSlice[failureSliceLen-i-1]
			if i+1 < failureSliceLen {
				prevFailureNum = transport.failureCountSecondSlice[failureSliceLen-i-2]
			}
		}

		currentAccessNum := 0
		prevAccessNum := 0
		if i < accessSliceLen {
			currentAccessNum = transport.accessCountSecondSlice[accessSliceLen-i-1]
			if i+1 < accessSliceLen {
				prevAccessNum = transport.accessCountSecondSlice[accessSliceLen-i-2]
			}
		}

		accessTimes += currentAccessNum - prevAccessNum
		failureTimes += currentFailureNum - prevFailureNum
	}

	return
}

func (transport *Transport) recordAccessCount() {
	//Write lock
	defer func() {
		transport.countSliceMutex.Unlock()
	}()
	transport.countSliceMutex.Lock()

	//slice fifo
	transport.accessCountSecondSlice = append(transport.accessCountSecondSlice[helper.MaxInt(len(transport.accessCountSecondSlice)-CountQueueLen, 0):], transport.GetAccessCount())
}

func (transport *Transport) recordFailureCount() {
	//Write lock
	defer func() {
		transport.countSliceMutex.Unlock()
	}()
	transport.countSliceMutex.Lock()

	//slice fifo
	transport.failureCountSecondSlice = append(transport.failureCountSecondSlice[helper.MaxInt(len(transport.failureCountSecondSlice)-CountQueueLen, 0):], transport.GetFailureCount())
}

func (transport *Transport) Reconnect() {
	if transport.S.ServerAddr != "" {
		transport.S.Listener.Close()
		transport.S.CloseChan <- true
		transport.S.ClientAddr = ""
		transport.T.(*http.Transport).CloseIdleConnections()
		<-transport.S.OpenChan
	}
	transport.T = createHttpTransport(transport.S)
}
