package dispatcher

import (
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/project"
	"goSpider/proxy"
	"goSpider/spider"
	"sync"
)

type Dispatcher struct {
	transportArr    []*proxy.Transport
	spiderArr       []*spider.Spider
	originSpiderArr []*spider.Spider
	initSSOnce      sync.Once
	initSpiderOnce  sync.Once
}

func (dispatcher *Dispatcher) GetSpiders() []*spider.Spider {
	return dispatcher.originSpiderArr
}

func (dispatcher *Dispatcher) InitSpider() []*spider.Spider {
	dispatcher.initSpiderOnce.Do(func() {
		for _, t := range dispatcher.InitTransport() {
			s := spider.New(t, nil)
			dispatcher.spiderArr = append(dispatcher.spiderArr, s)
			dispatcher.originSpiderArr = append(dispatcher.originSpiderArr, s)
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
				fmt.Println("proxy error: ", err)
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

func (dispatcher *Dispatcher) Run(project project.Project) {
	dispatcher.InitSpider()

	for _, l := range project.EntryUrl() {
		if !database.Bl().TestString(l) {
			database.AddUrlQueue(l)
		}
	}

	for _, s := range dispatcher.GetSpiders() {
		go func(s *spider.Spider) {
			for {
				project.Throttle(s)
				project.RequestBefore(s)
				s.Crawl(project.EnqueueFilter)
				project.ResponseAfter(s)
			}
		}(s)
	}
}
