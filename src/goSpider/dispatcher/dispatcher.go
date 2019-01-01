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
	first.Transport.LoopCountCut++
	return first
}

func (dispatcher *Dispatcher) Run(project project.Project) {
	dispatcher.InitSpider()

	for _, l := range project.EntryUrl() {
		if !database.Bl().TestString(l) {
			database.AddUrlQueue(l)
		}
	}

	go func() {
		t := time.NewTicker(time.Minute)
		for {
			<-t.C
			min := 999999999999.0
			for _, s := range dispatcher.spiderArr {
				if s.Transport.LoopCountCut < min {
					min = s.Transport.LoopCountCut
				}
			}

			//todo lock
			for _, s := range dispatcher.spiderArr {
				s.Transport.LoopCountCut /= min
			}
		}
	}()

	spiderChs := make(map[string]chan *spider.Spider)
	for _, s := range dispatcher.spiderArr {
		spiderChs[s.Transport.S.ServerAddr] = make(chan *spider.Spider, 1)
	}

	go func() {
		for {
			s := dispatcher.dispatcherSpider()
			spiderChs[s.Transport.S.ServerAddr] <- s
		}
	}()

	time.Sleep(1e9)

	for _, s := range dispatcher.spiderArr {
		go func(sss *spider.Spider) {
			for {
				s := <-spiderChs[sss.Transport.S.ServerAddr]
				project.Throttle(s)
				project.RequestBefore(s)
				s.Crawl(project.EnqueueFilter)
				project.ResponseAfter(s)
			}
		}(s)
	}

	stuck := make(chan int)
	<-stuck
}
