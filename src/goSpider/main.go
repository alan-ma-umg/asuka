package main

import (
	"fmt"
	"log"
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

	//database.Bl().ClearAll()                             //todo for test
	//database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	//c := &dispatcher.Dispatcher{}
	//c.Run(&project.Www{})
	//web.Server(c, ":666") // http://127.0.0.1:666

	fmt.Println(int(time.Duration(0.3e9/2).Seconds() * 1000))

}
