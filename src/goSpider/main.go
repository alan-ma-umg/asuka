package main

import (
	"fmt"
	"time"
	"goSpider/dispatcher"
	"goSpider/project"
	"goSpider/database"
	"goSpider/helper"
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
}
