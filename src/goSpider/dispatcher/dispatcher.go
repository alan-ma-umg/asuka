package dispatcher

import (
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/project"
	"goSpider/proxy"
	"sync"
	"goSpider/spider"
	"sort"
	"time"
)

type Dispatcher struct {
	transportArr   []*proxy.Transport
	spiderArr      []*spider.Spider
	initSSOnce     sync.Once
	initSpiderOnce sync.Once
}

func (dispatcher *Dispatcher) GetSpiders() []*spider.Spider {
	return dispatcher.spiderArr
}

func (dispatcher *Dispatcher) InitSpider() []*spider.Spider {
	dispatcher.initSpiderOnce.Do(func() {
		for _, t := range dispatcher.InitTransport() {
			dispatcher.spiderArr = append(dispatcher.spiderArr, spider.New(t, nil))
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

func (dispatcher *Dispatcher) dispatcherSpider() *spider.Spider {
	if len(dispatcher.spiderArr) == 0 {
		return nil
	}

	//rand.Seed(time.Now().Unix())
	sort.SliceStable(dispatcher.spiderArr, func(i, j int) bool {
		//rand.Seed(time.Now().Unix())
		//return rand.Intn(5) == 1
		return dispatcher.spiderArr[i].Transport.LoadBalanceRate() < dispatcher.spiderArr[j].Transport.LoadBalanceRate()
	})

	first := dispatcher.spiderArr[0]
	first.Transport.LoopCount++
	return first
}

func (dispatcher *Dispatcher) Run(project project.Project) {
	dispatcher.InitSpider()

	for _, l := range project.EntryUrl() {
		if !database.Bl().TestString(l) {
			database.AddUrlQueue(l)
		}
	}

	spiderChs := make(map[string]chan *spider.Spider)
	for _, s := range dispatcher.spiderArr {
		spiderChs[s.Transport.S.ServerAddr] = make(chan *spider.Spider, 1)
	}

	go func() {
		for {
			s := dispatcher.dispatcherSpider()
			//fmt.Println(1)
			spiderChs[s.Transport.S.ServerAddr] <- s
			//fmt.Println(2)
		}
	}()

	time.Sleep(1e9)

	for _, s := range dispatcher.spiderArr {

		//go func(sss *spider.Spider) {
		//	for {
		//		s := dispatcher.dispatcherSpider()
		//		spiderChs[s.Transport.S.ServerAddr] <- s
		//	}
		//}(s)

		go func(sss *spider.Spider) {
			for {
				//fmt.Println(3)
				s:=<-spiderChs[sss.Transport.S.ServerAddr]
				//fmt.Println(4)
				project.Throttle(s)
				project.RequestBefore(s)
				s.Crawl(project.EnqueueFilter)
				project.ResponseAfter(s)
			}
		}(s)
	}
	//
	//chs := make(chan int, len(dispatcher.spiderArr))
	//
	//for i := 0; i < len(dispatcher.spiderArr); i++ {
	//	go func() {
	//		for {
	//			s := dispatcher.dispatcherSpider()
	//			if s == nil {
	//				log.Fatal("nil spider")
	//			}
	//
	//			go func(s *spider.Spider) {
	//				project.Throttle(s)
	//				project.RequestBefore(s)
	//				s.Crawl(project.EnqueueFilter)
	//				project.ResponseAfter(s)
	//				chs <- 1
	//			}(s)
	//			<-chs
	//		}
	//	}()
	//}

	stuck := make(chan int)
	<-stuck
}
