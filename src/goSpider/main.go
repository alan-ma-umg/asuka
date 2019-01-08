package main

import (
	"fmt"
	"log"
	"time"
	"goSpider/dispatcher"
	"goSpider/project"
	"goSpider/web"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	mainStart := time.Now()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()
	//database.Bl().ClearAll()                             //todo for test
	//database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	c := &dispatcher.Dispatcher{}
	c.Run(&project.Www{})
	web.Server(c, ":666") // http://127.0.0.1:666
}
