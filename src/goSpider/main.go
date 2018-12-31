package main

import (
	"fmt"
	"time"
	"goSpider/dispatcher"
	"goSpider/project"
	"goSpider/database"
	"goSpider/helper"
	"goSpider/web"
)

func main() {
	mainStart := time.Now()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	database.Bl().ClearAll()                             //todo for test
	database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	go func() {
		web.Forever() // http://127.0.0.1:888
	}()

	c := &dispatcher.Dispatcher{}
	go func() {
		web.Monitor(c) // http://127.0.0.1:88/monitor
	}()
	c.Run(&project.DouBan{})
}
