package main

import (
	"asuka/dispatcher"
	"asuka/helper"
	"asuka/project"
	"asuka/queue"
	"asuka/web"
	"fmt"
	"log"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(os.Args) < 2 {
		log.Fatal("Example:/path/to/asuka /path/to/env.json")
	}
	helper.PathToEnvFile = os.Args[1]
}

func main() {
	mainStart := time.Now()
	//pSt := profile.Start(profile.MemProfile)
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
		//pSt.Stop()
	}()

	//todo for test
	//for i := 0; i < 10; i++ {
	//	os.Remove(helper.Env().BloomFilterPath + "enqueue_retry_" + strconv.Itoa(i) + ".db")
	//}
	//database.Mysql().Exec("truncate asuka_jian_shu")     //todo for test
	//database.Bl().ClearAll()                             //todo for test
	//database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	c := &dispatcher.Dispatcher{}
	c.Run(&project.JianShu{}, queue.NewQueue())
	web.Server(c, ":666") // http://127.0.0.1:666
}
