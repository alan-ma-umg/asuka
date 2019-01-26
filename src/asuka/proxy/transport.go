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
const MinuteInterval = 40 * 60
const CountQueueSecondCap = MinuteInterval * 2

var CountQueueMinuteCap = 1 //initial value, will dynamic changes

func init() {
	go func() {
		s := time.NewTicker(time.Second * SecondInterval)
		for {
			<-s.C
			for _, t := range transportList {
				t.CountSliceCursor++
				t.recordAccessSecondCount()
				t.recordFailureSecondCount()
			}
		}
	}()

	//ping
	go func() {
		for {
			time.Sleep(10e9)
			serverAddrMap := make(map[string][]*Transport)
			for _, t := range transportList {
				serverAddrMap[t.S.ServerAddr] = append(serverAddrMap[t.S.ServerAddr], t)
			}

			for host, transports := range serverAddrMap {
				go func(host string, transports []*Transport) {
					ipAddr, _ := lookIp(host)
					if ipAddr == nil {
						return
					}
					times := 3
					rtt, fail := helper.Ping(ipAddr, times)
					for _, t := range transports {
						if t.Ping == 0 {
							t.Ping = rtt
							t.PingFailureRate = float64(fail) / float64(times)
						} else {
							t.Ping = (rtt + t.Ping) / 2
							t.PingFailureRate = ((float64(fail) / float64(times)) + t.PingFailureRate) / 2
						}
					}
				}(host, transports)
			}
			time.Sleep(time.Minute)
		}
	}()
}

func lookIp(addr string) (*net.IPAddr, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	return net.ResolveIPAddr("ip4:icmp", host)
}

var transportList []*Transport

type Transport struct {
	S *SsAddr
	t http.RoundTripper

	countSliceMutex         sync.RWMutex
	CountSliceCursor        int
	AccessCountSecondSlice  []uint32
	FailureCountSecondSlice []uint32
	AccessCountMinuteSlice  []uint32
	FailureCountMinuteSlice []uint32

	AccessCountHistory  int
	FailureCountHistory int

	LoopCount int

	//traffic size
	TrafficIn  uint64
	TrafficOut uint64

	Ping            time.Duration
	PingFailureRate float64

	RecentFewTimesResult []bool
}

func NewTransport(ssAddr *SsAddr) (*Transport, error) {
	instance := &Transport{S: ssAddr, t: createHttpTransport(ssAddr), LoopCount: 0}
	transportList = append(transportList, instance)
	return instance, nil
}

func createHttpTransport(SockInfo *SsAddr) *http.Transport {
	t := &http.Transport{
		Proxy:                 nil, //disable system proxy
		MaxIdleConnsPerHost:   2,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if SockInfo.ServerAddr == "" { //no socks5 proxy

		t.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
	} else { //use socks5 proxy
		SockInfo.WaitUntilConnected() //waiting
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
	transport.AccessCountHistory++
}

// AddFailure 每次调用请求并失败时增加一次失败记录
func (transport *Transport) AddFailure(link string) {
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

func (transport *Transport) recordAccessSecondCount() {
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

func (transport *Transport) recordFailureSecondCount() {
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

func (transport *Transport) Reconnect() {
	if transport.S.ServerAddr != "" {
		transport.S.Close()
		transport.t.(*http.Transport).CloseIdleConnections()
	}
	transport.t = createHttpTransport(transport.S)
}

func (transport *Transport) GetHttpTransport() *http.Transport {
	return transport.t.(*http.Transport)
}
