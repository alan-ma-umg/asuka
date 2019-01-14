package main

import (
	"fmt"
	"goSpider/database"
	"goSpider/dispatcher"
	"goSpider/helper"
	"goSpider/project"
	"goSpider/queue"
	"goSpider/web"
	"log"
	"os"
	"strconv"
	"time"
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

	//todo for test
	for i := 0; i < 10; i++ {
		os.Remove(helper.Env().BloomFilterPath + "enqueue_retry_" + strconv.Itoa(i) + ".db")
	}
	database.Mysql().Exec("truncate asuka_zhi_hu")       //todo for test
	database.Bl().ClearAll()                             //todo for test
	database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	c := &dispatcher.Dispatcher{}
	c.Run(&project.ZhiHu{}, queue.NewQueue())
	web.Server(c, ":666") // http://127.0.0.1:666
}
