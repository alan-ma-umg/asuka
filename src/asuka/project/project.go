package project

import (
	"asuka/database"
	"asuka/helper"
	"asuka/proxy"
	"asuka/queue"
	"asuka/spider"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type IProject interface {
	// EntryUrl 万恶的起源
	// Firstly
	EntryUrl() []string

	// Throttle 控制抓取速度的一个地方
	// 使用spider.AddSleep()方法, 而不是time.Sleep().
	// Secondly
	Throttle(spider *spider.Spider)

	// RequestBefore http请求发起之前
	// Thirdly
	RequestBefore(spider *spider.Spider)

	// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
	// 一般用于根据Header过滤不想继续下载的response.content_type
	// Fourth
	DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error)

	// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
	// 一般用于采集数据的地方
	// Fifth
	ResponseSuccess(spider *spider.Spider)

	// EnqueueFilter HTTP完成并成功后, 从HTML中解析的每条URL都会经过这个筛选和处理. 空字符串则不入队列
	// Sixth
	EnqueueFilter(spider *spider.Spider, l *url.URL) string

	// ResponseAfter HTTP请求失败/成功之后
	// At Last
	ResponseAfter(spider *spider.Spider)
}

const RecentFetchCount = 50

type Dispatcher struct {
	IProject
	queue                *queue.Queue
	Spiders              []*spider.Spider
	Stop                 bool
	RecentFetchMutex     sync.Mutex
	RecentFetchLastIndex int64
	RecentFetchList      []*spider.Summary
}

func New(project IProject) *Dispatcher {
	d := &Dispatcher{IProject: project}

	d.queue = queue.NewQueue(d.GetProjectName())

	// kill signal handing
	helper.ExitHandleFuncSlice = append(helper.ExitHandleFuncSlice, func() {

		encStrSlice := make(map[string]interface{})
		for _, sp := range d.GetSpiders() {
			if sp.CurrentRequest() != nil && sp.CurrentRequest().URL != nil && sp.ResponseStr == "" {
				sp.Queue.Enqueue(sp.CurrentRequest().URL.String()) //check status & make improvement
				fmt.Println("enqueue " + sp.CurrentRequest().URL.String())
			}

			//gob
			gob.Register(d.IProject) //do register once
			encBuf := &bytes.Buffer{}
			enc := gob.NewEncoder(encBuf)
			err := enc.Encode(sp)
			if err != nil {
				log.Println(err)
				continue
			}
			encStrSlice[sp.Transport.S.ServerAddr] = encBuf.String()
		}
		//spider, write to redis
		database.Redis().Del(d.GetGOBKey())
		database.Redis().HMSet(d.GetGOBKey(), encStrSlice)

		//queue, write to file
		d.queue.BlSave()

		fmt.Println("Spiders status saved")
	})

	return d
}

func (my *Dispatcher) GetGOBKey() string {
	return my.GetProjectName() + "_gob"
}

func (my *Dispatcher) GetQueueKey() string {
	return my.queue.GetKey()
}

func (my *Dispatcher) GetProjectName() string {
	return strings.Split(reflect.TypeOf(my.IProject).String(), ".")[1]
}

func (my *Dispatcher) GetSpiders() []*spider.Spider {
	return my.Spiders
}

func (my *Dispatcher) InitSpider() []*spider.Spider {
	gobEnc, _ := database.Redis().HGetAll(my.GetGOBKey()).Result()

	for _, t := range my.InitTransport() {
		s := spider.New(t, my.queue)

		enable := s.Transport.S.Enable
		interval := s.Transport.S.Interval
		clientAddr := s.Transport.S.ClientAddr

		//recover from
		if recoverSpider, ok := gobEnc[s.Transport.S.ServerAddr]; ok {
			decBuf := &bytes.Buffer{}
			decBuf.WriteString(recoverSpider)
			dec := gob.NewDecoder(decBuf)
			err := dec.Decode(s)
			if err != nil {
				log.Println(err)
			}
			//s.ResetSleep()
		}

		s.Stop = !enable
		s.Transport.S.Enable = enable
		s.Transport.S.Interval = interval
		s.Transport.S.ClientAddr = clientAddr

		my.Spiders = append(my.Spiders, s)
	}
	return my.Spiders
}

func (my *Dispatcher) InitTransport() (transports []*proxy.Transport) {
	//append default transport
	dt, _ := proxy.NewTransport(&proxy.SsAddr{
		Name:     helper.Env().LocalTransport.Name,
		Enable:   helper.Env().LocalTransport.Enable,
		Interval: helper.Env().LocalTransport.Interval,
	})
	transports = append(transports, dt)

	var repeat []string
	for _, ssAddr := range proxy.SSLocalHandler() {
		if helper.Contains(repeat, ssAddr.ServerAddr) {
			log.Println("DUPLICATE: " + ssAddr.ServerAddr)
		}
		repeat = append(repeat, ssAddr.ServerAddr)

		<-ssAddr.OpenChan
		t, err := proxy.NewTransport(ssAddr)
		if err != nil {
			log.Println("proxy error: ", err)
			continue
		}
		transports = append(transports, t)
	}

	return
}

func (my *Dispatcher) Run() *Dispatcher {
	for _, l := range my.EntryUrl() {
		if !database.BlTestString(l) {
			my.queue.Enqueue(l)
		}
	}

	for _, s := range my.InitSpider() {
		go func(spider *spider.Spider) {
			for {
				for {
					if !my.Stop {
						break
					}
					time.Sleep(3e9)
				}
				Crawl(my, spider)
			}
		}(s)

		//ping
		go func(s *spider.Spider) {
			ipAddr, _ := lookIp(s.Transport.S.ServerAddr)
			for {
				if ipAddr == nil {
					time.Sleep(time.Minute)
					ipAddr, _ = lookIp(s.Transport.S.ServerAddr)
				} else {
					times := 3
					rtt, fail := helper.Ping(ipAddr, times)
					s.Transport.Ping = rtt
					s.Transport.PingFailureRate = float64(fail) / float64(times)
				}

				time.Sleep(7 * time.Second)
			}
		}(s)
	}

	return my
}

func (my *Dispatcher) CleanUp() *Dispatcher {
	for i := 0; i < 10; i++ {
		os.Remove(helper.Env().BloomFilterPath + my.GetProjectName() + "_enqueue_retry_" + strconv.Itoa(i) + ".db")
	}

	//database.Mysql().Exec("truncate asuka_dou_ban")
	database.Bl().ClearAll()
	database.Redis().Del(my.GetGOBKey())
	database.Redis().Del(my.GetQueueKey())
	return my
}

func Crawl(project *Dispatcher, spider *spider.Spider) {
	if project != nil {
		spider.RequestBefore = project.RequestBefore
		spider.DownloadFilter = project.DownloadFilter
		spider.ProjectThrottle = project.Throttle
	}
	spider.Throttle()

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

	defer func() {
		if project != nil {
			project.ResponseAfter(spider)
		}
	}()

	spider.Transport.LoopCount++

	_, summary, err := spider.Fetch(u)

	//recent fetch
	project.RecentFetchMutex.Lock()
	project.RecentFetchLastIndex++
	summary.Index = project.RecentFetchLastIndex
	summary.ConsumeTime = time.Since(spider.RequestStartTime).Truncate(time.Millisecond).String()
	project.RecentFetchList = append(project.RecentFetchList[helper.MaxInt(len(project.RecentFetchList)-RecentFetchCount, 0):], summary)
	project.RecentFetchMutex.Unlock()

	if err != nil {
		return
	}

	if project != nil {
		project.ResponseSuccess(spider)
	}

	for _, l := range spider.GetLinksByTokenizer() {
		enqueueUrl := ""
		if project != nil {
			enqueueUrl = project.EnqueueFilter(spider, l)
		} else {
			enqueueUrl = l.String()
		}

		if enqueueUrl != "" && database.BlTestAndAddString(enqueueUrl) {
			continue
		}

		if enqueueUrl != "" {
			spider.Queue.Enqueue(strings.TrimSpace(enqueueUrl))
		}
	}
}

func lookIp(addr string) (*net.IPAddr, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	return net.ResolveIPAddr("ip4:icmp", host)
}
