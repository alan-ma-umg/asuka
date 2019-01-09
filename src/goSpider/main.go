package main

import (
	"fmt"
	"log"
	"time"
	"net/url"
	"goSpider/helper"
	"golang.org/x/net/publicsuffix"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	mainStart := time.Now()
	//pSt := profile.Start(profile.MemProfile)
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
		//pSt.Stop()
	}()

	//database.Bl().ClearAll()                             //todo for test
	//database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	//c := &dispatcher.Dispatcher{}
	//c.Run(&project.Www{})
	//web.Server(c, ":666") // http://127.0.0.1:666


	u,_:=url.Parse("http://www.b11a24ij16i.com.cngoods-6461.html/")//fixme
	fmt.Println(u.Hostname())
	fmt.Println(publicsuffix.PublicSuffix(u.Hostname()))
	fmt.Println(helper.TldDomain(u))
	fmt.Println(u.Hostname())

}

