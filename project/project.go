package project

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/proxy"
	"github.com/chenset/asuka/queue"
	"github.com/chenset/asuka/spider"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

type IProject interface {
	// Init DoOnce func
	Init()

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

	//Showing 在web监控上展示信息
	Showing() string
	Name() string
}

type Implement struct {
}

func (my *Implement) Init() {
}
func (my *Implement) Showing() string {
	return "Have a nice day !"
}
func (my *Implement) ResponseSuccess(spider *spider.Spider) {
}
func (my *Implement) ResponseAfter(spider *spider.Spider) {
}
func (my *Implement) Name() string {
	return ""
}

const RecentFetchCount = 50

type Dispatcher struct {
	IProject
	*helper.Counting
	queue                *queue.Queue
	spiders              []*spider.Spider //make public for GOB
	Stop                 bool
	recentFetchMutex     sync.Mutex
	spiderSliceMutex     sync.Mutex
	RecentFetchLastIndex int64
	RecentFetchList      []*spider.Summary
	//traffic size
	TrafficIn  uint64
	TrafficOut uint64
}

func New(project IProject) *Dispatcher {
	d := &Dispatcher{IProject: project, Counting: &helper.Counting{}}
	gob.Register(project)
	d.queue = queue.NewQueue(d.Name())
	d.Init()

	// kill signal handing
	helper.ExitHandleFuncSlice = append(helper.ExitHandleFuncSlice, func() {
		if r := recover(); r != nil {
			fmt.Println("Exit error")
			fmt.Println(r)
		}

		for _, sp := range d.GetSpiders() {
			if sp.CurrentRequest() != nil && sp.CurrentRequest().URL != nil && sp.ResponseStr == "" {
				sp.Queue.Enqueue(sp.CurrentRequest().URL.String()) //check status & make improvement
				//fmt.Println("enqueue " + sp.CurrentRequest().URL.String())
			}
		}

		//gob
		encBuf := &bytes.Buffer{}
		if err := gob.NewEncoder(encBuf).Encode(d); err != nil {
			log.Println(err)
		} else {
			//spider, write to redis
			database.Redis().Del(d.getGOBKey())
			database.Redis().Set(d.getGOBKey(), encBuf.String(), 0)
		}

		//queue, write to file
		d.queue.BlSave()

		fmt.Println(d.Name() + " status saved")
	})

	return d
}

func (my *Dispatcher) getGOBKey() string {
	return my.Name() + "_gob"
}

func (my *Dispatcher) GetQueueKey() string {
	return my.queue.GetKey()
}

func (my *Dispatcher) Name() string {
	if name := my.IProject.Name(); name != "" {
		return name
	}

	return strings.Split(reflect.TypeOf(my.IProject).String(), ".")[1]
}

func (my *Dispatcher) GetSpiders() []*spider.Spider {
	return my.spiders
}

func (my *Dispatcher) initSpider() {
	defer database.Redis().Del(my.getGOBKey())

	gobEnc, err := database.Redis().Get(my.getGOBKey()).Result()
	recoverSpiders := make(map[string]*spider.Spider)
	if err == nil && gobEnc != "" {
		decBuf := &bytes.Buffer{}
		decBuf.WriteString(gobEnc)
		if err = gob.NewDecoder(decBuf).Decode(my); err != nil {
			log.Println(err)
		} else {
			for _, item := range my.spiders {
				recoverSpiders[item.Transport.S.ServerAddr] = item
			}
		}
		my.spiders = []*spider.Spider{}
	}

	for _, t := range my.initTransport() {
		s := spider.New(t, my.queue)

		name := s.Transport.S.Name
		enable := s.Transport.S.Enable
		interval := s.Transport.S.Interval
		clientAddr := s.Transport.S.ClientAddr

		//recover from
		if recoverSpider, ok := recoverSpiders[s.Transport.S.ServerAddr]; ok {
			encBuf := &bytes.Buffer{}
			if err = gob.NewEncoder(encBuf).Encode(recoverSpider); err != nil || gob.NewDecoder(encBuf).Decode(s) != nil {
				log.Println(err)
			}
		}

		s.Stop = !enable
		s.Transport.S.Name = name
		s.Transport.S.Enable = enable
		s.Transport.S.Interval = interval
		s.Transport.S.ClientAddr = clientAddr

		my.spiders = append(my.spiders, s)
	}
}

func (my *Dispatcher) initTransport() (transports []*proxy.Transport) {
	//append default transport
	dt, _ := proxy.NewTransport(&proxy.SsAddr{
		Name:       helper.Env().LocalTransport.Name,
		Group:      helper.Env().LocalTransport.Group,
		Enable:     helper.Env().LocalTransport.Enable,
		EnablePing: false,
		Interval:   helper.Env().LocalTransport.Interval,
		Type:       "local",
	})
	transports = append(transports, dt)
	var repeat []string

	for _, ssAddr := range proxy.HttpProxyHandler() {
		if helper.Contains(repeat, ssAddr.ServerAddr) {
			//log.Println("DUPLICATE: " + ssAddr.ServerAddr)
			continue
		}
		repeat = append(repeat, ssAddr.ServerAddr)

		t, err := proxy.NewTransport(ssAddr)
		if err != nil {
			log.Println("proxy error: ", err)
			continue
		}
		transports = append(transports, t)
	}

	for _, ssAddr := range proxy.SSLocalHandler() {
		if helper.Contains(repeat, ssAddr.ServerAddr) {
			log.Println("DUPLICATE: " + ssAddr.ServerAddr)
		}
		repeat = append(repeat, ssAddr.ServerAddr)

		t, err := proxy.NewTransport(ssAddr)
		if err != nil {
			log.Println("proxy error: ", err)
			continue
		}
		transports = append(transports, t)
	}

	return
}

func (my *Dispatcher) AddSpider(ssAddr *proxy.SsAddr) {
	my.spiderSliceMutex.Lock()
	defer my.spiderSliceMutex.Unlock()

	for _, oldSpider := range my.spiders {
		if oldSpider.Transport.S.ServerAddr == ssAddr.ServerAddr {
			return
		}
	}

	t, _ := proxy.NewTransport(ssAddr)
	s := spider.New(t, my.queue)
	my.spiders = append([]*spider.Spider{s}, my.spiders...)
	my.runSpider(s)
}

func (my *Dispatcher) runSpider(s *spider.Spider) {
	go func(spider *spider.Spider) {
		defer my.RemoveSpider(spider)
		for {
			for {
				if spider.Delete {
					return
				}
				if !my.Stop {
					break
				}
				spider.Transport.Close()
				time.Sleep(3e9)
			}
			Crawl(my, spider, nil)
		}
	}(s)
}

func (my *Dispatcher) Run() *Dispatcher {
	my.initSpider()

	for _, l := range my.EntryUrl() {
		if !my.queue.BlTestString(l) {
			my.queue.Enqueue(l)
		}
	}

	for _, s := range my.GetSpiders() {
		my.runSpider(s)
	}

	go func() {
		t := time.NewTicker(time.Second * helper.SecondInterval)
		defer t.Stop()
		for {
			<-t.C
			for _, s := range my.GetSpiders() {
				if s != nil {
					s.Transport.CountSliceCursor++
					s.Transport.RecordAccessSecondCount()
					s.Transport.RecordFailureSecondCount()
				}
			}
		}
	}()

	go func() {
		t := time.NewTicker(time.Second * helper.SecondInterval)
		defer t.Stop()
		for {
			<-t.C
			my.CountSliceCursor++
			my.RecordAccessSecondCount()
			my.RecordFailureSecondCount()
		}
	}()

	return my
}

func (my *Dispatcher) RemoveSpider(s *spider.Spider) {
	my.spiderSliceMutex.Lock()
	defer my.spiderSliceMutex.Unlock()
	var newSpiders []*spider.Spider
	for _, e := range my.GetSpiders() {
		if e != s {
			newSpiders = append(newSpiders, e)
		}
	}

	my.spiders = newSpiders
}

//func (my *Dispatcher) SearchSpider(serverName string) *spider.Spider {
//	for _, e := range my.GetSpiders() {
//		if e.Transport.S.Name == serverName {
//			return e
//		}
//	}
//	return nil
//}

func (my *Dispatcher) CleanUp() *Dispatcher {
	//database.Mysql().Exec("truncate asuka_dou_ban")
	my.queue.BlCleanUp()
	database.Redis().Del(my.getGOBKey())
	database.Redis().Del(my.GetQueueKey())
	return my
}

func Crawl(project *Dispatcher, spider *spider.Spider, dispatcherCallback func(spider *spider.Spider)) {
	if project != nil {
		spider.RequestBefore = project.RequestBefore
		spider.DownloadFilter = project.DownloadFilter
		spider.ProjectThrottle = project.Throttle
	}
	spider.Throttle(dispatcherCallback)

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

	project.AddAccess()
	_, summary, err := spider.Fetch(u)
	if err != nil || summary.StatusCode != 200 {
		project.AddFailure()
	}

	//IO
	project.TrafficIn += summary.TrafficIn
	project.TrafficOut += summary.TrafficOut

	//recent fetch
	project.recentFetchMutex.Lock()
	project.RecentFetchLastIndex++
	summary.Index = project.RecentFetchLastIndex
	summary.ConsumeTime = time.Since(spider.RequestStartTime).Truncate(time.Millisecond).String()
	project.RecentFetchList = append(project.RecentFetchList[helper.MaxInt(len(project.RecentFetchList)-RecentFetchCount, 0):], summary)
	project.recentFetchMutex.Unlock()

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

		if enqueueUrl == "" {
			continue
		}

		if spider.Queue.BlTestAndAddString(enqueueUrl) {
			continue
		}

		spider.Queue.Enqueue(strings.TrimSpace(enqueueUrl))
	}
}
