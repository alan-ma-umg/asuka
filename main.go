package main

import (
	"fmt"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/project"
	"github.com/chenset/asuka/web"
	"log"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
	fmt.Println("http://127.0.0.1:" + strings.Split(helper.Env().WEBListen, ":")[len(strings.Split(helper.Env().WEBListen, ":"))-1])
	log.Println(web.Server([]*project.Dispatcher{
		project.New(&project.DouBan{}).Run(),
		//project.New(&project.Pixiv{}).Run(),
		//project.New(&project.Test{}).Run(),
		//project.New(&project.Test2{}).Run(),
		//project.New(&project.ZhiHu{}).Run(),
		//project.New(&project.JianShu{}).Run(),
		//project.New(&project.Www{}).Run(),

		//project.New(&project.Pixiv{}).CleanUp().Run(),
		//project.New(&project.DouBan{}).CleanUp().Run(),
		//project.New(&project.Test2{}).CleanUp().Run(),
		//project.New(&project.ZhiHu{}).CleanUp().Run(),
		//project.New(&project.JianShu{}).CleanUp().Run(),
		//project.New(&project.Www{}).CleanUp().Run(),
	}, helper.Env().WEBListen))
	helper.ExitHandleFunc()
}
