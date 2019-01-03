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
			s := spider.New(t, nil)
			dispatcher.spiderArr = append(dispatcher.spiderArr, s)
		}
	})
	return dispatcher.spiderArr
}

func (dispatcher *Dispatcher) InitTransport() []*proxy.Transport {
	dispatcher.initSSOnce.Do(func() {
		if helper.Env().LocalTransportEnable {
			//append default transport
			dt, _ := proxy.NewTransport(&proxy.SsAddr{Weight: helper.Env().LocalTransportWeight})
			dispatcher.transportArr = append(dispatcher.transportArr, dt)
		}

		//todo 可用性维护
		for _, ssAddr := range proxy.SsLocalHandler() {
			t, err := proxy.NewTransport(ssAddr)
			if err != nil {
				fmt.Println("proxy error: ", err)
				continue
			}
			dispatcher.transportArr = append(dispatcher.transportArr, t)
		}
	})

	return dispatcher.transportArr
}

func (dispatcher *Dispatcher) Run(project project.Project) {
	for _, l := range project.EntryUrl() {
		if !database.Bl().TestString(l) {
			database.AddUrlQueue(l)
		}
	}

	for _, s := range dispatcher.InitSpider() {
		go func(s *spider.Spider) {
			for {
				s.Throttle()
				project.Throttle(s)
				project.RequestBefore(s)
				s.Crawl(project.EnqueueFilter)
				project.ResponseAfter(s)
			}
		}(s)
	}
}
