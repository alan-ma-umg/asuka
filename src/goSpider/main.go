package main

import (
	"fmt"
	"time"
	"goSpider/web"
	"goSpider/dispatcher"
	"goSpider/project"
)

func main() {
	mainStart := time.Now()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	//database.Bl().ClearAll()                             //todo for test
	//database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	c := &dispatcher.Dispatcher{}
	go func() {
		web.Server(c, ":888") // http://127.0.0.1:888
	}()
	c.Run(&project.DouBan{})
}
