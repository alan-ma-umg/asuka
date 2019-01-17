package proxy

import (
	"asuka/helper"
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"math"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const SecondInterval = 1

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

	countSliceMutex   sync.RWMutex
	accessCountSlice  []int
	failureCountSlice []int

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
	t := http.DefaultTransport
	if ssAddr == nil {
		ssAddr = &SsAddr{}
	}

	if ssAddr.ServerAddr != "" {
		//socks5 proxy
		dialer, err := proxy.SOCKS5("tcp", ssAddr.ClientAddr, nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			return nil, err
		}

		//http transport
		t = &http.Transport{
			//MaxIdleConnsPerHost: 2,
			MaxIdleConns:        10,
			IdleConnTimeout:     20 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,

			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
	}
	instance := &Transport{S: ssAddr, T: t, LoopCount: 0}
	transportList = append(transportList, instance)
	return instance, nil
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

	sliceLen := len(transport.accessCountSlice)

	times := int(math.Ceil(float64(second) / SecondInterval))
	for i := 0; i < times; i++ {
		currentNum := 0
		prevNum := 0
		if i < sliceLen {
			currentNum = transport.accessCountSlice[sliceLen-i-1]
			if i+1 < sliceLen {
				prevNum = transport.accessCountSlice[sliceLen-i-2]
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

	accessSliceLen := len(transport.accessCountSlice)
	failureSliceLen := len(transport.failureCountSlice)

	times := int(math.Ceil(float64(second) / SecondInterval))
	for i := 0; i < times; i++ {
		currentFailureNum := 0
		prevFailureNum := 0
		if i < failureSliceLen {
			currentFailureNum = transport.failureCountSlice[failureSliceLen-i-1]
			if i+1 < failureSliceLen {
				prevFailureNum = transport.failureCountSlice[failureSliceLen-i-2]
			}
		}

		currentAccessNum := 0
		prevAccessNum := 0
		if i < accessSliceLen {
			currentAccessNum = transport.accessCountSlice[accessSliceLen-i-1]
			if i+1 < accessSliceLen {
				prevAccessNum = transport.accessCountSlice[accessSliceLen-i-2]
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
	transport.accessCountSlice = append(transport.accessCountSlice[helper.MaxInt(len(transport.accessCountSlice)-CountQueueLen, 0):], transport.GetAccessCount())
}

func (transport *Transport) recordFailureCount() {
	//Write lock
	defer func() {
		transport.countSliceMutex.Unlock()
	}()
	transport.countSliceMutex.Lock()

	//slice fifo
	transport.failureCountSlice = append(transport.failureCountSlice[helper.MaxInt(len(transport.failureCountSlice)-CountQueueLen, 0):], transport.GetFailureCount())
}

func (transport *Transport) Reconnect() {
	transport.T.(*http.Transport).DisableKeepAlives = false
	transport.T.(*http.Transport).CloseIdleConnections()
	transport.S.CloseChan <- true
	transport.S.Listener.Close()
}
