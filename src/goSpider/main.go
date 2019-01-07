package main

import (
	"fmt"
	"goSpider/database"
	"goSpider/dispatcher"
	"goSpider/helper"
	"goSpider/project"
	"goSpider/web"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mainStart := time.Now()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()
	database.Bl().ClearAll()                             //todo for test
	database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	c := &dispatcher.Dispatcher{}
	c.Run(&project.Www{})
	web.Server(c, ":666") // http://127.0.0.1:666
}
