package main

import (
	"asuka/helper"
	"asuka/project"
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
	//database.Mysql().Exec("truncate asuka_dou_ban")      //todo for test
	//database.Bl().ClearAll()                             //todo for test
	//database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	//c := &dispatcher.Dispatcher{}
	//c.Run([]project.Project{&project.Test{}}, queue.NewQueue())
	//fmt.Println("Monitor: http://127.0.0.1:666")

	//project.Dispatcher{}

	p := project.New(&project.Test{})
	p.Run()

	z := project.New(&project.ZhiHu{})
	z.Run()

	fmt.Println("Monitor: http://127.0.0.1:666")
	web.Server([]*project.Dispatcher{p, z}, ":666") // http://127.0.0.1:666
}
