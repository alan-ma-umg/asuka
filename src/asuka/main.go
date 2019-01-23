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
	p := project.New(&project.Test{})

	//cleanUp(p) //todo !!!!!!!!!

	p.Run()

	z := project.New(&project.ZhiHu{})
	//cleanUp(z) //todo !!!!!!!!!

	z.Run()
	fmt.Println("Monitor: http://127.0.0.1:666")
	projects := []*project.Dispatcher{p, z}

	web.Server(projects, ":666") // http://127.0.0.1:666
}

func cleanUp(p *project.Dispatcher) {
	for i := 0; i < 10; i++ {
		os.Remove(helper.Env().BloomFilterPath + p.GetProjectName() + "_enqueue_retry_" + strconv.Itoa(i) + ".db")
	}
	//database.Mysql().Exec("truncate asuka_dou_ban")
	database.Bl().ClearAll()
	database.Redis().Del("gob_" + p.GetProjectName())
	database.Redis().Del(p.GetQueueKey())
}
