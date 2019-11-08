package project

import (
	"bytes"
	"encoding/gob"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/queue"
	"github.com/chenset/asuka/spider"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const RecentFetchCount = 50

type Dispatcher struct {
	IProject
	*helper.Counting
	queue            *queue.Queue
	spiders          []*spider.Spider //write this slice need to under spiderSliceMutex.lock
	spidersWaiting   []*spider.Spider //waiting for execute, write this slice need to under spiderSliceMutex.lock
	StartTime        time.Time
	StopTime         time.Time
	recentFetchMutex sync.Mutex
	spiderSliceMutex sync.Mutex
	//queueRetriesCapMutex sync.Mutex
	RecentFetchLastIndex int64
	tcpFilterErrorCount  int
	RecentFetchList      []*spider.Summary
	TrafficIn            uint64
	TrafficOut           uint64
	QueueRetries         []int
}

func New(project IProject, stopTime time.Time) *Dispatcher {
	d := &Dispatcher{IProject: project, StopTime: stopTime, Counting: &helper.Counting{}, StartTime: time.Now(), QueueRetries: make([]int, 1)}
	gob.Register(project)

	// kill signal handing
	helper.ExitHandleFuncSlice = append(helper.ExitHandleFuncSlice, func() {
		if r := recover(); r != nil {
			log.Println("Exit error")
			log.Println(r)
		}

		for _, sp := range d.GetSpiders() {
			if sp != nil && sp.CurrentRequest() != nil && sp.CurrentRequest().URL != nil && len(sp.ResponseByte) == 0 {
				sp.GetQueue().Enqueue(sp.CurrentRequest().URL.String()) //check status & make improvement
				//fmt.Println("enqueue " + sp.CurrentRequest().URL.String())
			}
		}

		//queue, write to file
		d.GetQueue().BlSave(true)

		//清空前获取
		GOBRedisKey := d.getGOBKey()
		//projectName := d.Name()

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
		//log.Println(projectName + " status saved")
	})

	return d
}

func (my *Dispatcher) IsStop() bool {
	if my.StopTime.IsZero() || time.Since(my.StopTime).Seconds() < 0 {
		return false
	}
	return true
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

func (my *Dispatcher) initProject() {
	defer database.Redis().Del(my.getGOBKey())

	// recover Dispatcher
	gobEnc, err := database.Redis().Get(my.getGOBKey()).Result()
	if err == nil && gobEnc != "" {
		decBuf := &bytes.Buffer{}
		decBuf.WriteString(gobEnc)
		gob.NewDecoder(decBuf).Decode(my)
	}

	//append default transport
	u, _ := url.Parse("direct://localhost")
	my.AddSpider(u)

	my.Init(my)

	rawUrls := my.EntryUrl()
	if int64(len(rawUrls)) > my.GetQueue().QueueLen() {
		my.GetQueue().Enqueue(rawUrls)
	}
}

func (my *Dispatcher) GetQueue() *queue.Queue {
	if my.queue == nil { //todo DoOnce in struct
		my.queue = queue.NewQueue(my.Name(), my.InitBloomFilterCapacity())
	}
	return my.queue
}

//AddSpider 加入的spider不是直接立即执行的, 会通过addSpidersWaiting添加到spidersWaiting中等待合适时机执行
func (my *Dispatcher) AddSpider(addr *url.URL) {
	my.spiderSliceMutex.Lock()
	defer my.spiderSliceMutex.Unlock()

	for _, oldSpider := range my.spiders {
		if oldSpider.TransportUrl.Host == addr.Host {
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

func (my *Dispatcher) EnqueueForFailure(spider *spider.Spider, err error, rawUrl string, retryTimes int) {
	go func() {
		success, tries := my.IProject.EnqueueForFailure(spider, err, rawUrl, rawUrl, retryTimes)
		if !success {
			my.QueueRetries[0]++
			return
		}

		if len(my.QueueRetries) <= tries {
			//my.queueRetriesCapMutex.Lock()
			my.QueueRetries = append(my.QueueRetries, make([]int, 1+tries-len(my.QueueRetries))...)
			//my.queueRetriesCapMutex.Unlock()
		}

		my.QueueRetries[tries]++
	}()
}

func (my *Dispatcher) runSpider(s *spider.Spider) {
	go func(spider *spider.Spider) {
		if my != nil {
			spider.RequestBefore = my.RequestBefore
			spider.DownloadFilter = my.DownloadFilter
			spider.ProjectThrottle = my.Throttle
			spider.EnqueueForFailure = my.EnqueueForFailure
		}

		for {
			if spider.Delete {
				my.RemoveSpider(spider) //上级调用也有锁, 这里也有所. 但是隔着一层Go Goroutine
				return
			}
			if my.IsStop() {
				spider.ResetSpider()
				my.addSpidersWaiting(spider, true) //上级调用也有锁, 这里也有所. 但是隔着一层Go Goroutine
				return
			}
			Crawl(my, spider, nil)
		}
	}(s)
}

func (my *Dispatcher) Run() *Dispatcher {
	go func() {
		my.initProject()
	}()

	//transport counter
	go func() {
		for {
			time.Sleep(time.Second * helper.SecondInterval)
			for _, s := range my.GetSpiders() {
				if s != nil {
					s.CountSliceCursor++
					s.RecordAccessSecondCount()
					s.RecordFailureSecondCount()
				}
			}
		}
	}()

	//project counter
	go func() {
		for {
			time.Sleep(time.Second * helper.SecondInterval)
			my.CountSliceCursor++
			my.RecordAccessSecondCount()
			my.RecordFailureSecondCount()
		}
	}()

	//tcp filter client mode
	if helper.Env().BloomFilterClient != "" {
		//empty tcpFilterErrorCount
		go func() {
			for {
				time.Sleep(time.Minute * 32)
				my.tcpFilterErrorCount = 0
			}
		}()
	} else {
		//release queue.BloomFilterInstance
		go func() {
			for {
				time.Sleep(time.Second * 500)

				//release queue.BloomFilterInstance
				//if all of spiders are idle, release queue.BloomFilterInstance after durations of stop
				if my.IsStop() && time.Since(my.StopTime).Seconds() > 300 {
					for _, s := range my.GetSpiders() {
						if s != nil && !s.IsIdle() {
							return //it's not idle
						}
					}

					my.GetQueue().ResetBloomFilterInstance()
				}
			}
		}()
	}

	//watching and run SpidersWaiting
	go func() {
		for {
			for {
				if !my.IsStop() {
					break
				}
				time.Sleep(time.Second * 5)
			}

			time.Sleep(2e9)

			my.spiderSliceMutex.Lock()
			for _, s := range my.getSpidersWaiting() {
				my.runSpider(s)
			}
			my.spidersWaiting = nil //nil slice
			my.spiderSliceMutex.Unlock()

			time.Sleep(10e9)
		}
	}()
	return my
}

func (my *Dispatcher) CleanUp() *Dispatcher {
	my.GetQueue().BlCleanUp()            //bloom filter & tcp bloom filter
	database.Redis().Del(my.getGOBKey()) //GOB
	my.GetQueue().CleanQueue()           //queue
	my.GetQueue().CleanFailure()         //queue failure
	my.QueueRetries = make([]int, 1)     //queue failure
	return my
}

func Crawl(project *Dispatcher, spider *spider.Spider, dispatcherCallback func(spider *spider.Spider)) {
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
			spider.ResponseByte = nil
		}
	}()

	project.AddAccess()
	spider.AddAccess()
	summary, err := project.Fetch(spider, u)

	if err != nil || summary.StatusCode != 200 {
		spider.AddFailure()
		project.AddFailure()
	} else {
		go func() {
			// remove from retries queue
			// async , wil broke the system if too fast
			project.GetQueue().DequeueForFailure(link)
		}()
	}

	if summary == nil {
		return
	}

	contentType := ""
	if spider.CurrentResponse() != nil {
		contentType = strings.ToLower(spider.CurrentResponse().Header.Get("Content-type"))
	} else {
		contentType = strings.ToLower(http.DetectContentType(spider.ResponseByte))
	}

	summary.ContentType = contentType

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

	if len(spider.ResponseByte) <= 10 || !strings.Contains(contentType, "html") {
		return
	}

	// async , wil broke the system if too fast
	//go func(urls []*url.URL) {

	var testUrls []string
	for _, l := range spider.GetLinksByTokenizer() {
		testUrl := ""
		if project != nil {
			testUrl = project.EnqueueFilter(spider, l)
		} else {
			testUrl = l.String()
		}
		if testUrl == "" {
			continue
		}
		testUrls = append(testUrls, strings.TrimSpace(testUrl))
		summary.FindUrls++
	}

	//bloom filter test
	exists, err := spider.GetQueue().BlTestAndAddStrings(testUrls)
	if err != nil || len(exists) != len(testUrls) {
		go func() {
			spider.GetQueue().Enqueue(link) //enqueue again
		}()
		log.Println("exists len: " + strconv.Itoa(len(exists)) + " testUrls len: " + strconv.Itoa(len(testUrls)))
		if err != nil {
			log.Println(err)
		}
		tcpFilterErrorHandle(project)
		return
	}

	//enqueue
	var enqueueUrls []string
	for i, b := range exists {
		if !b {
			summary.NewUrls++
			enqueueUrls = append(enqueueUrls, testUrls[i])
		}
	}
	spider.GetQueue().Enqueue(enqueueUrls)

	//for _, l := range urls {
	//	enqueueUrl := ""
	//	if project != nil {
	//		enqueueUrl = project.EnqueueFilter(spider, l)
	//	} else {
	//		enqueueUrl = l.String()
	//	}
	//
	//	if enqueueUrl == "" {
	//		continue
	//	}
	//
	//	summary.FindUrls++
	//	exists, err := spider.GetQueue().BlTestAndAddString(enqueueUrl)
	//	if err != nil {
	//		log.Println(err)
	//		tcpFilterErrorHandle(project)
	//		return //return and stop the project
	//	}
	//	if exists {
	//		continue
	//	}
	//	summary.NewUrls++
	//	spider.GetQueue().Enqueue(strings.TrimSpace(enqueueUrl))
	//}
	//}(spider.GetLinksByTokenizer())
}

func tcpFilterErrorHandle(project *Dispatcher) {
	project.tcpFilterErrorCount++
	if project.tcpFilterErrorCount > 5 {
		if !project.IsStop() {
			log.Println("too many failures, stop")
		}
		project.StopTime = time.Now()
	}
}
