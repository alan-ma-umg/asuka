package dispatcher

import (
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/project"
	"goSpider/proxy"
	"goSpider/queue"
	"goSpider/spider"
	"net"
	"sync"
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

func (dispatcher *Dispatcher) InitSpider(queue *queue.Queue) []*spider.Spider {
	dispatcher.initSpiderOnce.Do(func() {
		for _, t := range dispatcher.InitTransport() {
			s := spider.New(t, nil, queue)
			dispatcher.spiderArr = append(dispatcher.spiderArr, s)
		}
	})
	return dispatcher.spiderArr
}

func (dispatcher *Dispatcher) InitTransport() []*proxy.Transport {
	dispatcher.initSSOnce.Do(func() {
		if helper.Env().LocalTransport.Enable {
			//append default transport
			dt, _ := proxy.NewTransport(&proxy.SsAddr{})
			dt.S.Name = helper.Env().LocalTransport.Name
			dt.S.Enable = helper.Env().LocalTransport.Enable
			dt.S.Interval = helper.Env().LocalTransport.Interval
			dispatcher.transportArr = append(dispatcher.transportArr, dt)
		}

		//todo 可用性维护
		for _, ssAddr := range proxy.SSLocalHandler() {
			if !ssAddr.Enable {
				continue
			}

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

func (dispatcher *Dispatcher) Run(project project.Project, queue *queue.Queue) {

	for _, l := range project.EntryUrl() {
		if !database.BlTestString(l) {
			queue.Enqueue(l)
		}
	}

	for _, s := range dispatcher.InitSpider(queue) {
		go func(s *spider.Spider) {
			for {
				s.Throttle()
				project.Throttle(s)
				project.RequestBefore(s)
				s.Crawl(project.EnqueueFilter)
				project.ResponseAfter(s)
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

				time.Sleep(7)
			}
		}(s)
	}
}

func lookIp(addr string) (*net.IPAddr, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	return net.ResolveIPAddr("ip4:icmp", host)
}
