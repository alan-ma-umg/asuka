package proxy

import (
	"asuka/helper"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestHTTPProxyTransport(t *testing.T) {
	//creating the proxyURL
	proxyStr := ""
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	//creating the URL to be loaded through the proxy
	urlStr := "https://www.douban.com/"
	uurl, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", uurl.String(), nil)
	if err != nil {
		log.Println(err)
	}

	//calling the URL
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	//getting the response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	log.Println(string(data))
}

func TestTransport(t *testing.T) {
	t1, _ := NewTransport(&SsAddr{})

	//t1.LoadRate(86400 * 2)
	//log.Println(86400 * 2)
	CountQueueMinuteCap = 86400 * 3 * 200 / (60 * 10)
	for i := 0; i < CountQueueSecondCap; i++ {
		t1.AddAccess()
		t1.AddAccess()
		t1.AddAccess()
		t1.AddFailure()
		//t1.AddFailure("sfdsfsdf")

		//t1.AddAccess("sfdsfsdf")
		//if rand.Intn(2) == 2 {
		//	t1.AddFailure("sfdsfsdf")
		//	t1.AddAccess("sfdsfsdf")
		//}

		//if rand.Intn(3) == 2 {
		//	t1.AddAccess("sfdsfsdf")
		//	t1.AddAccess("sfdsfsdf")
		//	t1.AddFailure("sfdsfsdf")
		//}

		//time.Sleep(0.1e9)
		//fmt.Println(time.Since(s))
		//fmt.Println("Load: ", t1.LoadRate(5))
		//fmt.Println("Load: ", t1.LoadRate(5))
		//fmt.Println("Fail: ", t1.FailureRate(5))
		//fmt.Println("Load: ", t1.LoadRate(120))
		//fmt.Println("Fail: ", t1.FailureRate(120))
		//fmt.Println("Load: ", t1.LoadRate(300))
		//fmt.Println("Fail: ", t1.FailureRate(300))
		//fmt.Println("Load: ", t1.LoadRate(600))
		//fmt.Println("Fail: ", t1.FailureRate(600))

		//helper.PrintMemUsage()
		//t1.countSliceCursor++
		t1.recordAccessSecondCount()
		t1.recordFailureSecondCount()
	}
	//t1.countSliceCursor++
	//t1.recordAccessSecondCount()
	//t1.recordFailureSecondCount()
	helper.PrintMemUsage()
	//time.Sleep(2e9)
	//time.Sleep(2e9)
	s := time.Now()
	for i := 0; i < 10000; i++ {
		t1.LoadRate(5)
		t1.LoadRate(60)
		t1.LoadRate(900)
		t1.LoadRate(1800)
		t1.LoadRate(200000)
		helper.SpiderFailureRate(t1.AccessCount(30 * 60))
	}
	fmt.Println(time.Since(s))
	fmt.Println("Load: ", t1.LoadRate(30*60+10))
	fmt.Println("Load: ", t1.LoadRate(60))
	fmt.Println("Load: ", t1.LoadRate(5))
	fmt.Println("Load: ", t1.LoadRate(60*10))
	fmt.Println("Load: ", t1.LoadRate(60*10))
	fmt.Println("Load: ", t1.LoadRate(899))
	fmt.Println("Load: ", t1.LoadRate(900))
	fmt.Println("Load: ", t1.LoadRate(901))
	fmt.Println("Load: ", t1.LoadRate(910))
	fmt.Println("Load: ", t1.LoadRate(900))
	fmt.Println("Load: ", t1.LoadRate(1200))
	fmt.Println("Load: ", t1.LoadRate(1500))
	fmt.Println("Load: ", t1.LoadRate(1790))
	fmt.Println("Load: ", t1.LoadRate(1795))
	fmt.Println("Load: ", t1.LoadRate(1798))
	fmt.Println("Load: ", t1.LoadRate(1799))
	fmt.Println("Load: ", t1.LoadRate(1800))
	fmt.Println("Load: ", t1.LoadRate(1801))
	fmt.Println("Load: ", t1.LoadRate(1802))
	fmt.Println("Load: ", t1.LoadRate(CountQueueSecondCap-2))
	fmt.Println("Load: ", t1.LoadRate(CountQueueSecondCap-1))
	fmt.Println("Load: ", t1.LoadRate(CountQueueSecondCap))
	fmt.Println("Load: ", t1.LoadRate(CountQueueSecondCap+1))
	fmt.Println("Load: ", t1.LoadRate(CountQueueSecondCap+2))
	fmt.Println("Load: ", t1.LoadRate(3600))

	fmt.Println(helper.SpiderFailureRate(t1.AccessCount(30 * 60)))
	//log.Println(len(t1.accessCountSecondSlice))
	//log.Println(len(t1.accessCountMinuteSlice))
	log.Println(t1.AccessCount(CountQueueSecondCap + 1))

	startTime := time.Now()

	time.Sleep(1e9)

	fmt.Println(int(time.Since(startTime).Seconds()))

}
