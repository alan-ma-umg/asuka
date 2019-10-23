package project

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/queue"
	"github.com/chenset/asuka/spider"
	"io"
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
	Init(my *Dispatcher)

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

	// EnqueueForFailure 请求或者响应失败时重新入失败队列, 可以修改这里修改加入失败队列的实现
	EnqueueForFailure(spider *spider.Spider, err error, rawUrl string, retryTimes int)

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
	WEBSite(w http.ResponseWriter, r *http.Request)
}

type Implement struct{}

func (my *Implement) Init() {}
func (my *Implement) Showing() string {
	return "Have a nice day !"
}

// EnqueueForFailure 请求或者响应失败时重新入失败队列, 可以修改这里修改加入失败队列的实现
func (my *Implement) EnqueueForFailure(spider *spider.Spider, err error, rawUrl string, retryTimes int) {
	spider.GetQueue().EnqueueForFailure(rawUrl, retryTimes)
}
func (my *Implement) ResponseSuccess(spider *spider.Spider) {}

// ResponseAfter HTTP请求失败/成功之后
// At Last
func (my *Implement) ResponseAfter(spider *spider.Spider) {
	spider.ResetSpider() //现在是每请求一次, 就重置一次. 请求代理也会重新连接
}

func (my *Implement) Name() string {
	return ""
}
func (my *Implement) WEBSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=UTF-8")
	io.WriteString(w, my.Name())
}

const RecentFetchCount = 50

type Dispatcher struct {
	IProject
	*helper.Counting
	queue                *queue.Queue
	spiders              []*spider.Spider //write this slice need to under spiderSliceMutex.lock
	spidersWaiting       []*spider.Spider //waiting for execute, write this slice need to under spiderSliceMutex.lock
	Stop                 bool
	recentFetchMutex     sync.Mutex
	spiderSliceMutex     sync.Mutex
	RecentFetchLastIndex int64
	RecentFetchList      []*spider.Summary
	TrafficIn            uint64
	TrafficOut           uint64
	StartTime            time.Time
}

func New(project IProject, stop bool) *Dispatcher {
	d := &Dispatcher{IProject: project, Stop: stop, Counting: &helper.Counting{}, StartTime: time.Now()}
	gob.Register(project)

	// kill signal handing
	helper.ExitHandleFuncSlice = append(helper.ExitHandleFuncSlice, func() {
		if r := recover(); r != nil {
			fmt.Println("Exit error")
			fmt.Println(r)
		}

		for _, sp := range d.GetSpiders() {
			if sp != nil && sp.CurrentRequest() != nil && sp.CurrentRequest().URL != nil && len(sp.ResponseByte) == 0 {
				sp.GetQueue().Enqueue(sp.CurrentRequest().URL.String()) //check status & make improvement
				//fmt.Println("enqueue " + sp.CurrentRequest().URL.String())
			}
		}

		//queue, write to file
		d.GetQueue().BlSave()

		//清空前获取
		GOBRedisKey := d.getGOBKey()
		projectName := d.Name()

		//gob
		d.RecentFetchLastIndex = 0 //序列化前清空
		d.RecentFetchList = nil    //序列化前清空
		//d.IProject = nil           //这里提前清空容易导致其他地方还未退出时读取到空指针
		encBuf := &bytes.Buffer{}
		if err := gob.NewEncoder(encBuf).Encode(d); err != nil {
			log.Println(err)
		} else {
			//spider, write to redis
			database.Redis().Del(GOBRedisKey)
			database.Redis().Set(GOBRedisKey, encBuf.String(), 0)
		}

		fmt.Println(projectName + " status saved")
	})

	return d
}

func (my *Dispatcher) getGOBKey() string {
	return my.Name() + "_gob"
}

func (my *Dispatcher) GetQueueKey() string {
	return my.GetQueue().GetKey()
}

func (my *Dispatcher) Name() string {
	if name := my.IProject.Name(); name != "" {
		return name
	}

	return strings.Split(reflect.TypeOf(my.IProject).String(), ".")[1]
}

//GetSpiders need to check each item it's not nil when foreach this return value
func (my *Dispatcher) GetSpiders() []*spider.Spider {
	return my.spiders
}

func (my *Dispatcher) getSpidersWaiting() []*spider.Spider {
	return my.spidersWaiting
}

func (my *Dispatcher) initSpider() { //todo 这个需要重构
	defer database.Redis().Del(my.getGOBKey())

	// recover Dispatcher
	gobEnc, err := database.Redis().Get(my.getGOBKey()).Result()
	//recoverSpiders := make(map[string]*spider.Spider)
	if err == nil && gobEnc != "" {
		decBuf := &bytes.Buffer{}
		decBuf.WriteString(gobEnc)
		gob.NewDecoder(decBuf).Decode(my)
		//log.Println(err)
		//} else {
		//	for _, item := range my.spiders {
		//		recoverSpiders[item.Transport.S.Host] = item
		//	}
		//}
		//my.spiders = []*spider.Spider{}
	}

	//append default transport
	u, _ := url.Parse("direct://localhost")
	my.AddSpider(u)
	//s := spider.New(u, my.queue)
	//my.spiders = append(my.spiders, s)

	//for _, t := range my.initTransport() {

	//name := s.Transport.S.Name
	//enable := s.Transport.S.Enable
	//interval := s.Transport.S.Interval
	//clientAddr := s.Transport.S.ClientAddr

	//recover from
	//if recoverSpider, ok := recoverSpiders[s.Transport.S.Host]; ok {
	//	encBuf := &bytes.Buffer{}
	//	if err = gob.NewEncoder(encBuf).Encode(recoverSpider); err != nil || gob.NewDecoder(encBuf).Decode(s) != nil {
	//		log.Println(err)
	//	}
	//}

	//s.Stop = t.S.Stop
	//s.Transport.S.Name = name
	//s.Transport.S.Enable = enable
	//s.Transport.S.Interval = interval
	//s.Transport.S.ClientAddr = clientAddr

	//}
}

//func (my *Dispatcher) initTransport() (transports []*proxy.Transport) {
//append default transport
//u, _ := url.Parse("direct://localhost")
//dt := proxy.NewTransport(&proxy.AddrInfo{URL: u})
//dt.S.Stop = !helper.Env().LocalTransport
//return append(transports, dt)
//}

func (my *Dispatcher) GetQueue() *queue.Queue {
	if my.queue == nil {
		my.queue = queue.NewQueue(my.Name())
	}
	return my.queue
}

//AddSpider 加入的spider不是直接立即执行的, 会通过addSpidersWaiting添加到spidersWaiting中等待合适时机执行
func (my *Dispatcher) AddSpider(addr *url.URL) {
	my.spiderSliceMutex.Lock()
	defer my.spiderSliceMutex.Unlock()

	for _, oldSpider := range my.spiders {
		if oldSpider.Transport.U.Host == addr.Host {
			return
		}
	}

	s := spider.New(addr, my.GetQueue)
	//my.spiders = append([]*spider.Spider{s}, my.spiders...) // 为了让localhost在最前
	my.spiders = append(my.spiders, s)

	my.addSpidersWaiting(s, false) //调用处已经有锁,不用再次检查
	return
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

	var newSpidersWaiting []*spider.Spider
	for _, e := range my.getSpidersWaiting() {
		if e != s {
			newSpidersWaiting = append(newSpidersWaiting, e)
		}
	}
	my.spidersWaiting = newSpidersWaiting

	s.ResetSpider()
	s = nil
}

//addSpidersWaiting 添加待执行的spider
func (my *Dispatcher) addSpidersWaiting(s *spider.Spider, checkLock bool) {
	if checkLock {
		my.spiderSliceMutex.Lock()
	}
	my.spidersWaiting = append(my.spidersWaiting, s)
	if checkLock {
		my.spiderSliceMutex.Unlock()
	}
}

func (my *Dispatcher) loopRunSpidersWaiting() {
	go func() {
		for {
			for {
				if !my.Stop {
					break
				}
				time.Sleep(time.Second * 5)
			}

			my.spiderSliceMutex.Lock()
			for _, s := range my.getSpidersWaiting() {
				my.runSpider(s)
			}
			my.spidersWaiting = nil //nil slice
			my.spiderSliceMutex.Unlock()

			time.Sleep(10e9)
		}
	}()
}

func (my *Dispatcher) runSpider(s *spider.Spider) {
	go func(spider *spider.Spider) {
		for {
			if spider.Delete {
				my.RemoveSpider(spider)
				return
			}
			if my.Stop {
				spider.ResetSpider()
				my.addSpidersWaiting(spider, true) //上级调用也有锁, 这里也有所. 但是隔着一层Go Goroutine
				return
			}
			Crawl(my, spider, nil)
		}
	}(s)
}

func (my *Dispatcher) Run() *Dispatcher {
	my.initSpider()
	my.Init(my)

	for _, l := range my.EntryUrl() {
		//if !my.GetQueue().BlTestString(l) {
		my.GetQueue().Enqueue(l)
		//}
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

	my.loopRunSpidersWaiting()
	return my
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
	my.GetQueue().BlRemoveFile()
	database.Redis().Del(my.getGOBKey())
	database.Redis().Del(my.GetQueueKey())
	return my
}

func Crawl(project *Dispatcher, spider *spider.Spider, dispatcherCallback func(spider *spider.Spider)) {
	if project != nil {
		spider.RequestBefore = project.RequestBefore
		spider.DownloadFilter = project.DownloadFilter
		spider.ProjectThrottle = project.Throttle
		spider.EnqueueForFailure = project.EnqueueForFailure
	}
	spider.Throttle(dispatcherCallback)

	link, err := spider.GetQueue().Dequeue()
	if err != nil {
		time.Sleep(time.Second * 5)
		return
	}

	u, err := url.Parse(link)
	if err != nil {
		//log.Println("URL parse failed ", link, err)
		return
	}

	defer func() {
		if project != nil {
			project.ResponseAfter(spider)
		}
	}()

	project.AddAccess()
	summary, err := spider.Fetch(u)
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

	//非图片才抓取links, todo 还需要处理为更准确的类型
	if !strings.Contains(spider.CurrentResponse().Header.Get("Content-type"), "image") {
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

			if spider.GetQueue().BlTestAndAddString(enqueueUrl) {
				continue
			}

			spider.GetQueue().Enqueue(strings.TrimSpace(enqueueUrl))
		}
	}
}
