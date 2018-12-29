package spider

import (
	"fmt"
	"goSpider/proxy"
	"testing"
	"time"
)

func TestSpider_Crawl(t *testing.T) {
	var transportArr []*proxy.Transport
	for _, ssAddr := range proxy.SsLocalHandler() {
		fmt.Println(ssAddr)
		t, err := proxy.NewTransport(ssAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		transportArr = append(transportArr, t)
	}
	//link, _ := url.Parse("http://ip.flysay.com")

	for _, t := range transportArr {
		go func(t *proxy.Transport) {
			//r, _ := http.NewRequest("GET", link.String(), nil)
			New(t, nil).Crawl(nil)
		}(t)
	}

	time.Sleep(time.Hour * 24)

	//u, _ := url.Parse("http://www.baidu.com")
	//sp := New(u)
	//sp.Fetch()
}

func TestSpider_Fetch(tt *testing.T) {

}

func TestSpider_GetImageLinks(t *testing.T) {
}

func TestSpider_GetLinks(t *testing.T) {
}
