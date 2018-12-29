package spider

import (
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/project"
	"goSpider/proxy"
	"log"
	"sort"
	"sync"
	"time"
)

type Dispatcher struct {
	transportArr   []*proxy.Transport
	spiderArr      []*Spider
	initSSOnce     sync.Once
	initSpiderOnce sync.Once
}

func (dispatcher *Dispatcher) InitSpider() []*Spider {
	dispatcher.initSpiderOnce.Do(func() {
		for _, t := range dispatcher.InitTransport() {
			dispatcher.spiderArr = append(dispatcher.spiderArr, New(t, nil))
		}
	})
	return dispatcher.spiderArr
}

func (dispatcher *Dispatcher) InitTransport() []*proxy.Transport {
	dispatcher.initSSOnce.Do(func() {

		//todo 可用性维护
		for _, ssAddr := range proxy.SsLocalHandler() {
			t, err := proxy.NewTransport(ssAddr)
			if err != nil {
				fmt.Println(err)
				continue
			}
			dispatcher.transportArr = append(dispatcher.transportArr, t)
		}

		if helper.Env().LocalTransportEnable {
			//append default transport
			dt, _ := proxy.NewTransport(&proxy.SsAddr{Weight: helper.Env().LocalTransportWeight})
			dispatcher.transportArr = append(dispatcher.transportArr, dt)
		}
	})

	return dispatcher.transportArr
}

func (dispatcher *Dispatcher) dispatcherSpider() *Spider {
	if len(dispatcher.spiderArr) == 0 {
		return nil
	}

	sort.SliceStable(dispatcher.spiderArr, func(i, j int) bool {
		return dispatcher.spiderArr[i].Transport.LoadBalanceRate() < dispatcher.spiderArr[j].Transport.LoadBalanceRate()
	})

	return dispatcher.spiderArr[0]
}

func (dispatcher *Dispatcher) Run(project project.Project) {
	dispatcher.InitSpider()

	for _, l := range project.EntryUrl() {
		if !database.Bl().TestString(l) {
			database.AddUrlQueue(l)
		}
	}
	database.UrlQueueSave()

	chs := make(chan time.Duration, len(dispatcher.spiderArr))
	for {
		s := dispatcher.dispatcherSpider()
		if s == nil {
			log.Fatal("nil spider")
		}

		if project.NeedToPause(s) {
			time.Sleep(10e9)
			continue
		}

		if project.NeedToLogin(s) && !project.IsLogin(s) {
			project.Login(s)
		}

		go func(s *Spider) {
			st := time.Now()
			s.Crawl(project.EnqueueFilter)
			et := time.Since(st)

			project.Throttle(s)

			chs <- et
		}(s)
		fmt.Println(<-chs)

	}
}

//func (dispatcher *Dispatcher) Run(entryLinks []string, filter func(spider *Spider, l *url.URL) bool) {
//	dispatcher.InitSpider()
//
//	for _, l := range entryLinks {
//		if !database.Bl().TestString(l) {
//			database.AddUrlQueue(l)
//		}
//	}
//	database.UrlQueueSave()
//
//	chs := make(chan time.Duration, len(dispatcher.spiderArr))
//	for {
//		s := dispatcher.dispatcherSpider()
//		if s == nil {
//			log.Fatal("nil spider")
//		}
//		go func(s *Spider) {
//			st := time.Now()
//			s.Crawl(filter)
//			chs <- time.Since(st)
//		}(s)
//		fmt.Println(<-chs)
//	}
//}
