package spider

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"testing"
	"time"
)

func TestSpider_Crawl(t *testing.T) {
	fmt.Println(int(math.Ceil((time.Duration(0) + 1).Seconds())))
}

func TestSpider_Fetch(tt *testing.T) {

}

func TestSpider_GetImageLinks(t *testing.T) {
}

func TestSpider_GetLinks(t *testing.T) {
	u, _ := url.Parse("https://www.zhihu.com/question/37362725/answer/152869802#ds987gf/fgd45")
	//u.Fragment = ""
	fmt.Println(u.String())
	log.Println(u)

	uu, _ := url.Parse("https://www.zhihu.com/terms#sec-zhihu-bean")
	newUU := uu.ResolveReference(u)
	newUU.Fragment = ""
	fmt.Println(newUU.String())
	log.Println(newUU)
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
