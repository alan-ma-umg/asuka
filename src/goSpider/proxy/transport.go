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
const countQueueLen = 2000
const TLSHandshakeTimeout = 10e9

//time url
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

	accessCountList  *list.List
	failureCountList *list.List

	loopCount int
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
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
			TLSHandshakeTimeout: TLSHandshakeTimeout,
		}
	}

	instance := &Transport{S: ssAddr, T: t, accessCountList: list.New(), failureCountList: list.New()}
	transportList = append(transportList, instance)
	return instance, nil
}

// AddAccess 每次调用请求时增加一次记录, 无论是否成功
func (transport *Transport) AddAccess(link string) {
	transport.loopCount++
	transport.AccessList = append(transport.AccessList, link)
}

// AddFailure 每次调用请求并失败时增加一次失败记录
func (transport *Transport) AddFailure(link string) {
	transport.FailureList = append(transport.FailureList, link)
}

func (transport *Transport) LoadBalanceRate() float64 {
	//todo 考虑失败率
	rate := float64(transport.loopCount) * transport.LoadRate(60) / float64(transport.S.Weight)
	transport.loopCount = 1
	return rate
}

// LoadRate 获取指定秒数内的负载值.参数最小值5秒, 最大取值 countQueueLen*SecondInterval-1
func (transport *Transport) LoadRate(second int) float64 {
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

// FailureRate 获取指定秒数内的失败率.参数最小值5秒, 最大取值 countQueueLen*SecondInterval-1
func (transport *Transport) FailureRate(second int) float64 {
	rate := 0.0
	failureCursor := transport.failureCountList.Back()
	accessCursor := transport.accessCountList.Back()
	times := int(math.Ceil(float64(second) / SecondInterval))
	for i := 0; i < times; i++ {
		currentFailureNum := 0
		if failureCursor != nil {
			currentFailureNum = failureCursor.Value.(int)
			failureCursor = failureCursor.Prev()
		}

		currentAccessNum := 0
		if accessCursor != nil {
			currentAccessNum = accessCursor.Value.(int)
			accessCursor = accessCursor.Prev()
		}

		if currentAccessNum != 0 {
			rate += float64(currentFailureNum/SecondInterval) / float64(currentAccessNum/SecondInterval)
		}
	}

	return rate / float64(times)
}

func (transport *Transport) recordAccessCount() {
	transport.accessCountList.PushBack(len(transport.AccessList))
	if transport.accessCountList.Len() > countQueueLen {
		transport.accessCountList.Remove(transport.accessCountList.Front()) // FIFO
	}
}

func (transport *Transport) recordFailureCount() {

	transport.failureCountList.PushBack(len(transport.FailureList))
	if transport.failureCountList.Len() > countQueueLen {
		transport.failureCountList.Remove(transport.failureCountList.Front()) // FIFO
	}
}
