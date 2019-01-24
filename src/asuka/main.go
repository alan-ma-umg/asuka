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
	helper.ExitHandle()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	asuka()
}

func asuka() {
	fmt.Println("Monitor: http://127.0.0.1:666")

	web.Server([]*project.Dispatcher{
		project.New(&project.DouBan{}).Run(),
		project.New(&project.Test{}).Run(),
		//project.New(&project.ZhiHu{}).Run(),
		//project.New(&project.JianShu{}).Run(),
		//project.New(&project.Www{}).Run(),

		//project.New(&project.Test{}).CleanUp().Run(),
		//project.New(&project.ZhiHu{}).CleanUp().Run(),
		//project.New(&project.JianShu{}).CleanUp().Run(),
	}, ":666") // http://127.0.0.1:666
}

//2019/01/24 11:56:11 h2_bundle.go:8723: protocol error: received *http.http2GoAwayFrame before a SETTINGS frame
