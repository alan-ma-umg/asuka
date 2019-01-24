package main

import (
	"asuka/database"
	"asuka/helper"
	"asuka/project"
	"asuka/web"
	"fmt"
	"log"
	"os"
	"strconv"
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
	helper.ExitHandle()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	asuka()
}

func asuka() {
	projects := []*project.Dispatcher{
		project.New(&project.DouBan{}),
		project.New(&project.Test{}),
		//project.New(&project.ZhiHu{}),
		//project.New(&project.JianShu{}),
	}

	for _, p := range projects {
		//cleanUp(p) //todo !!!!!!!!!
		p.Run()
	}

	fmt.Println("Monitor: http://127.0.0.1:666")

	web.Server(projects, ":666") // http://127.0.0.1:666
}

func cleanUp(p *project.Dispatcher) {
	for i := 0; i < 10; i++ {
		os.Remove(helper.Env().BloomFilterPath + p.GetProjectName() + "_enqueue_retry_" + strconv.Itoa(i) + ".db")
	}
	//database.Mysql().Exec("truncate asuka_dou_ban")
	database.Bl().ClearAll()
	database.Redis().Del(p.GetGOBKey())
	database.Redis().Del(p.GetQueueKey())
}

//2019/01/24 11:56:11 h2_bundle.go:8723: protocol error: received *http.http2GoAwayFrame before a SETTINGS frame
