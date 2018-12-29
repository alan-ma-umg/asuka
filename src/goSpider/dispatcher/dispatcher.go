package dispatcher

import (
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/project"
	"goSpider/proxy"
	"log"
	"sort"
	"sync"
	"goSpider/spider"
)

type Dispatcher struct {
	transportArr   []*proxy.Transport
	spiderArr      []*spider.Spider
	initSSOnce     sync.Once
	initSpiderOnce sync.Once
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

	chs := make(chan int, len(dispatcher.spiderArr))

	for i := 0; i < len(dispatcher.spiderArr); i++ {
		go func() {
			for {
				s := dispatcher.dispatcherSpider()
				if s == nil {
					log.Fatal("nil spider")
				}

				go func(s *spider.Spider) {
					project.Throttle(s)
					project.RequestBefore(s)
					s.Crawl(project.EnqueueFilter)
					project.ResponseAfter(s)
					chs <- 1
				}(s)
				<-chs
			}
		}()
	}

	stuck := make(chan int)
	<-stuck
}
