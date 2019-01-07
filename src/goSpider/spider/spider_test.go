package spider

import (
	"fmt"
	"goSpider/proxy"
	"testing"
	"time"
	"net/http"
	"log"
	"io"
	"golang.org/x/net/html"
	"net/url"
	"io/ioutil"
	"bytes"
	"regexp"
)

func TestSpider_Crawl(t *testing.T) {
	var transportArr []*proxy.Transport
	for _, ssAddr := range proxy.SSLocalHandler() {
		fmt.Println(ssAddr)
		t, err := proxy.NewTransport(ssAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		transportArr = append(transportArr, t)
	}

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

func TestSpider_Regex(t *testing.T) {
	//fmt.Println(123123)
	//log.Println(435345345)

	baseUrl, _ := url.Parse("http://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=1&ch=&tn=baiduerr&bar=&wd=1")
	resp, err := http.Get(baseUrl.String())

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	st := time.Now()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(time.Since(st))

	//cp := ioutil.NopCloser(bytes.NewBuffer(b))
	//
	//st = time.Now()
	//for _, sub := range getLinks(cp) {
	//	_, err := url.Parse(sub)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	(baseUrl.ResolveReference(baseUrl).String())
	//}
	//
	//fmt.Println(time.Since(st))
	//st = time.Now()
	//cp = ioutil.NopCloser(bytes.NewBuffer(b))
	//fmt.Println(getLinks(cp))
	//fmt.Println(time.Since(st))

	bodyStr := string(b)
	st = time.Now()
	getLinksByRegex(bodyStr)
	fmt.Println(time.Since(st))

	st = time.Now()
	getLinksByTokenizer(ioutil.NopCloser(bytes.NewBuffer(b)))
	fmt.Println(time.Since(st))

	//res, _ := ioutil.ReadAll(resp.Body)
	//
	//var linkRegex, _ = regexp.Compile("<a[^>]+href=\"([(\\.|h|/)][^\"]+)\"[^>]*>")
	//
	//for _, sub := range linkRegex.FindAllStringSubmatch(string(res[:]), -1) {
	//	u, err := url.Parse(sub[1])
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	//arr = append(arr, spider.CurrentRequest.URL.ResolveReference(u))
	//	fmt.Println(u.String())
	//}
}

var llinkRegex, _ = regexp.Compile("<a[^>]+href=\"([(\\.|h|/)][^\"]+)\"[^>]*>")

func getLinksByRegex(str string) (res []string) {
	for _, sub := range llinkRegex.FindAllStringSubmatch(str, -1) {
		//u, err := url.Parse(sub[1])
		//if err != nil {
		//	log.Fatal(err)
		//}
		res = append(res, sub[1])
	}

	return
}

func getLinksByTokenizer(body io.Reader) (res []string) {
	token := html.NewTokenizer(body)
	for {
		switch next := token.Next(); next {
		case html.StartTagToken:
			token := token.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {


						res = append(res, attr.Val)
					}
				}
			}
		case html.ErrorToken:
			return
		}
	}
	return
}
