package main

import (
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"time"
	"goSpider/project"
	"goSpider/dispatcher"
)

func main() {
	mainStart := time.Now()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	database.Bl().ClearAll()                             //todo for test
	database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	c := &dispatcher.Dispatcher{}
	c.Run(&project.DouBan{})

	//c := &spider.Dispatcher{}

	//c.Run()
	//c.Run([]string{"https://www.douban.com/"}, func(s *spider.Spider, l *url.URL) bool {
	//	pass := false
	//	for _, white := range []string{"movie.douban.com", "book.douban.com"} {
	//		if strings.Contains(strings.ToLower(l.Hostname()), white) {
	//			pass = true
	//		}
	//	}
	//	return pass
	//})

}
