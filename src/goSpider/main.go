package main

import (
	"fmt"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/spider"
	"net/url"
	"strings"
	"time"
)

func main() {
	mainStart := time.Now()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	database.Bl().ClearAll()                             //todo for test
	database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	c := &spider.Dispatcher{}
	c.Run([]string{"https://www.douban.com/"}, func(s *spider.Spider, l *url.URL) bool {
		pass := false
		for _, white := range []string{"movie.douban.com", "book.douban.com"} {
			if strings.Contains(strings.ToLower(l.Hostname()), white) {
				pass = true
			}
		}
		return pass
	})

}
