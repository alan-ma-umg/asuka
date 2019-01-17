package dispatcher

import (
	"asuka/database"
	"asuka/helper"
	"asuka/project"
	"asuka/proxy"
	"asuka/queue"
	"asuka/spider"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
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

		for _, ssAddr := range proxy.SSLocalHandler() {
			if !ssAddr.Enable {
				continue
			}

			for {
				if ssAddr.ClientAddr != "" {
					break
				}
				fmt.Println("Waiting for socks proxy")
				time.Sleep(time.Second / 10)
			}

			t, err := proxy.NewTransport(ssAddr)
			if err != nil {
				log.Println("proxy error: ", err)
				continue
			}
			dispatcher.transportArr = append(dispatcher.transportArr, t)
		}

		fmt.Println("Socks proxy is ready to go")
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
		go func(spider *spider.Spider) {
			for {

				Crawl(project, spider)
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

func Crawl(project project.Project, spider *spider.Spider) {
	if project != nil {
		spider.RequestBefore = project.RequestBefore
		spider.DownloadFilter = project.DownloadFilter
		spider.ProjectThrottle = project.Throttle
	}
	spider.Throttle()

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

	spider.Transport.LoopCount++

	_, err = spider.Fetch(u)
	if err != nil {
		return
	}

	if project != nil {
		project.ResponseSuccess(spider)
	}

	for _, l := range spider.GetLinksByTokenizer() {
		if database.BlTestAndAddString(l.String()) {
			continue
		}

		if project != nil && !project.EnqueueFilter(spider, l) {
			continue
		}

		spider.Queue.Enqueue(strings.TrimSpace(l.String()))
	}
}

func lookIp(addr string) (*net.IPAddr, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	return net.ResolveIPAddr("ip4:icmp", host)
}
