package proxy

import (
	"container/list"
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"math"
	"net"
	"net/http"
	"os"
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

	AccessList  []string
	FailureList []string

	accessCountList     *list.List
	failureCountList    *list.List
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
	instance := &Transport{S: ssAddr, T: t, accessCountList: list.New(), failureCountList: list.New(), LoopCount: 0}
	transportList = append(transportList, instance)
	return instance, nil
}

// AddAccess 每次调用请求时增加一次记录, 无论是否成功
func (transport *Transport) AddAccess(link string) {
	transport.AccessList = append(transport.AccessList, link)
}

// AddFailure 每次调用请求并失败时增加一次失败记录
func (transport *Transport) AddFailure(link string) {
	transport.FailureList = append(transport.FailureList, link)
}

func updateCountQueueLen(second int) {
	if CountQueueLen <= second/SecondInterval {
		CountQueueLen = second/SecondInterval + 2
	}
}

// LoadRate 获取指定秒数内的负载值.参数最小值SecondInterval秒
//todo 性能优化
func (transport *Transport) LoadRate(second int) float64 {
	updateCountQueueLen(second)
	rate := 0.0
	cursor := transport.accessCountList.Back()
	times := int(math.Ceil(float64(second) / SecondInterval))
	for i := 0; i < times; i++ {
		currentNum := 0
		prevNum := 0
		if cursor != nil {
			currentNum = cursor.Value.(int)
			cursor = cursor.Prev()
			if cursor != nil {
				prevNum = cursor.Value.(int)
			}
		}
		rate += float64(currentNum-prevNum) / SecondInterval
	}

	return rate / float64(times)
}

//AccessCount  获取指定秒数内的访问数/失败j数量.参数最小值SecondInterval秒
//todo 性能优化
func (transport *Transport) AccessCount(second int) (accessTimes, failureTimes int) {
	updateCountQueueLen(second)
	failureCursor := transport.failureCountList.Back()
	accessCursor := transport.accessCountList.Back()
	times := int(math.Ceil(float64(second) / SecondInterval))
	for i := 0; i < times; i++ {
		currentFailureNum := 0
		prevFailureNum := 0
		if failureCursor != nil {
			currentFailureNum = failureCursor.Value.(int)
			failureCursor = failureCursor.Prev()
			if failureCursor != nil {
				prevFailureNum = failureCursor.Value.(int)
			}
		}

		currentAccessNum := 0
		prevAccessNum := 0
		if accessCursor != nil {
			currentAccessNum = accessCursor.Value.(int)
			accessCursor = accessCursor.Prev()
			if accessCursor != nil {
				prevAccessNum = accessCursor.Value.(int)
			}
		}

		accessTimes += currentAccessNum - prevAccessNum
		failureTimes += currentFailureNum - prevFailureNum
	}

	return
}

func (transport *Transport) GetAccessCount() int {
	return transport.accessCountHistory + len(transport.AccessList)
}

func (transport *Transport) recordAccessCount() {
	transport.accessCountList.PushBack(transport.GetAccessCount())
	if transport.accessCountList.Len() > CountQueueLen {
		transport.accessCountList.Remove(transport.accessCountList.Front()) // FIFO
	}

	//todo lock
	listLen := len(transport.AccessList)
	limit := 1000
	if listLen > limit {
		transport.AccessList = transport.AccessList[0 : limit/2-1] //fixme !!!!!!!!!!!!!! [limit/2-1:]
		transport.accessCountHistory += listLen - limit/2
	}
}

func (transport *Transport) GetFailureCount() int {
	return transport.FailureCountHistory + len(transport.FailureList)
}

func (transport *Transport) recordFailureCount() {
	transport.failureCountList.PushBack(transport.GetFailureCount())
	if transport.failureCountList.Len() > CountQueueLen {
		transport.failureCountList.Remove(transport.failureCountList.Front()) // FIFO
	}

	//todo lock
	listLen := len(transport.FailureList)
	limit := 1000
	if listLen > limit {
		transport.FailureList = transport.FailureList[0 : limit/2-1] //fixme !!!!!!!!!!!!!! [limit/2-1:]
		transport.FailureCountHistory += listLen - limit/2
	}
}
